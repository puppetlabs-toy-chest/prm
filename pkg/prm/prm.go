// nolint:structcheck,unused
package prm

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/go-version"
	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ToolConfigName     = "prm-config"
	ToolConfigFileName = "prm-config.yml"
)

type Prm struct {
	AFS           *afero.Afero
	IOFS          *afero.IOFS
	RunningConfig Config
	CodeDir       string
	CacheDir      string
	Cache         map[string]*Tool
	Backend       BackendI
}

type PuppetVersion struct {
	version semver.Version
}

type Group struct {
	ID    string     `yaml:"id"`
	Tools []ToolInst `yaml:"tools"`
}

type ValidateYmlContent struct {
	Groups []Group `yaml:"groups"`
}

func (p *Prm) getValidateFilePath() (string, error) {
	validateFile := filepath.Join(p.CodeDir, "validate.yml")
	if _, err := p.AFS.Stat(validateFile); err != nil {
		log.Info().Msgf("Reference the 'prm help exec' help section for exec command usage.")
		return "", err
	}
	return validateFile, nil
}

func (p *Prm) getGroupsFromFile(validateFile string) ([]Group, error) {
	contentBytes, err := p.AFS.ReadFile(validateFile)
	if err != nil {
		log.Error().Msgf("Error reading validate.yml: %s", err)
		return []Group{}, err
	}

	var contentStruct ValidateYmlContent
	err = yaml.Unmarshal(contentBytes, &contentStruct)
	if err != nil {
		log.Error().Msgf("validate.yml is not formatted correctly: %s", err)
		return []Group{}, err
	}

	return contentStruct.Groups, nil
}

func checkDuplicateToolsInGroups(tools []ToolInst) error {
	toolNames := make(map[string]bool)

	for _, tool := range tools {
		if toolNames[tool.Name] {
			return fmt.Errorf("duplicate tool '%s' found. Validation groups cannot contain duplicate tools", tool.Name)
		}
		toolNames[tool.Name] = true
	}

	return nil
}

func getSelectedGroup(groups []Group, selectedGroupID string) (Group, error) {
	if selectedGroupID == "" && len(groups) > 0 {
		if selectedGroupID == "" {
			log.Warn().Msgf("No group specified. Defaulting to the '%s' tool group", groups[0].ID)
		}
		selectedGroupID = groups[0].ID
	}

	for _, group := range groups {
		if group.ID == selectedGroupID {
			err := checkDuplicateToolsInGroups(group.Tools)
			if err != nil {
				return Group{}, err
			}
			log.Info().Msgf("Found tool group: %v ", group.ID)
			return group, nil
		}
	}

	return Group{}, fmt.Errorf("specified tool group '%s' not found", selectedGroupID)
}

func (p *Prm) GetValidationGroupFromFile(selectedGroupID string) (Group, error) {
	// check if validate.yml exits in the codeDir
	validateFile, err := p.getValidateFilePath()
	if err != nil {
		return Group{}, err
	}

	groups, err := p.getGroupsFromFile(validateFile)
	if err != nil {
		return Group{}, err
	}

	return getSelectedGroup(groups, selectedGroupID)
}

// Check to see if the requested tool can be found installed.
// If installed read the tool configuration and return
func (p *Prm) IsToolAvailable(tool string) (*Tool, bool) {

	if p.Cache[tool] != nil {
		return p.Cache[tool], true
	}

	return nil, false
}

// Check to see if the tool is ready to execute
func (p *Prm) IsToolReady(tool *Tool) bool {
	err := p.Backend.GetTool(tool, p.RunningConfig)
	return err == nil
}

// What version of Puppet is requested by the user
func (*Prm) getPuppetVersion() PuppetVersion {
	return PuppetVersion{}
}

func (p *Prm) readToolConfig(configFile string) Tool {
	file, err := p.AFS.ReadFile(configFile)
	if err != nil {
		log.Error().Msgf("unable to read tool config, %s", configFile)
	}

	var tool Tool

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(file))
	if err != nil {
		log.Error().Msgf("unable to read tool config, %s: %s", configFile, err.Error())
		return Tool{}
	}
	err = viper.Unmarshal(&tool.Cfg)

	if err != nil {
		log.Error().Msgf("unable to parse tool config, %s", configFile)
		return Tool{}
	}

	return tool
}

// List lists all templates in a given path and parses their configuration. Does
// not return any errors from parsing invalid templates, but returns them as
// debug log events
func (p *Prm) List(toolPath string, toolName string, onlyValidators bool) error {
	log.Debug().Msgf("Searching %+v for tool configs", toolPath)
	// Triple glob to match author/id/version/ToolConfigFileName
	matches, _ := p.IOFS.Glob(toolPath + "/**/**/**/" + ToolConfigFileName)

	var tmpls []ToolConfig
	for _, file := range matches {
		log.Debug().Msgf("Found: %+v", file)
		i := p.readToolConfig(file)
		if i.Cfg.Plugin != nil {
			if onlyValidators && !i.Cfg.Common.CanValidate {
				log.Debug().Msgf("Not a validator: %+v", file)
				continue
			}
			i.Cfg.Path = filepath.Dir(file)
			tmpls = append(tmpls, i.Cfg)
		}
	}

	if len(tmpls) == 0 {
		if onlyValidators {
			return fmt.Errorf("no validators found in %+v", toolPath)
		}
		return fmt.Errorf("no tools found in %+v", toolPath)
	}

	if toolName != "" {
		log.Debug().Msgf("Filtering for: %s", toolName)
		tmpls = p.FilterFiles(tmpls, func(f ToolConfig) bool { return f.Plugin.Id == toolName })
	}

	tmpls = p.filterNewestVersions(tmpls)

	// cache for use with the rest of the program
	// this is a seperate cache from the one used by the CLI
	p.createToolCache(tmpls)

	return nil
}

func (p *Prm) filterNewestVersions(tt []ToolConfig) (ret []ToolConfig) {
	for _, t := range tt {
		id := t.Plugin.Id
		author := t.Plugin.Author
		// Look for tools with the same author and id
		tools := p.FilterFiles(tt, func(f ToolConfig) bool { return f.Plugin.Id == id && f.Plugin.Author == author })
		if len(tools) > 1 {
			// If the author/id template has 2+ entries, that's multiple versions
			// check first to see if the return list already has an entry for this template
			if len(p.FilterFiles(ret, func(f ToolConfig) bool { return f.Plugin.Id == id && f.Plugin.Author == author })) == 0 {
				// turn the version strings into version objects for sorting and comparison
				versionsRaw := []string{}
				for _, t := range tools {
					versionsRaw = append(versionsRaw, t.Plugin.Version)
				}
				versions := make([]*version.Version, len(versionsRaw))
				for i, raw := range versionsRaw {
					v, _ := version.NewVersion(raw)
					versions[i] = v
				}
				sort.Sort(version.Collection(versions))
				// select the latest version
				highestVersion := versions[len(versions)-1]
				highestVersionTemplate := p.FilterFiles(tools, func(f ToolConfig) bool {
					actualVersion, _ := version.NewVersion(f.Plugin.Version)
					return actualVersion.Equal(highestVersion)
				})
				ret = append(ret, highestVersionTemplate[0])
			}
		} else {
			// If the author/id template only has 1 entry, it's already the latest version
			ret = append(ret, t)
		}
	}

	return ret
}

func (p *Prm) FilterFiles(ss []ToolConfig, test func(ToolConfig) bool) (ret []ToolConfig) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func (p *Prm) createToolCache(tmpls []ToolConfig) {
	// initialise the cache
	p.Cache = make(map[string]*Tool)
	// Iterate through the list of tool configs and
	// add them to the map
	for _, t := range tmpls {
		name := t.Plugin.Author + "/" + t.Plugin.Id
		tool := Tool{
			Cfg: t,
		}
		p.Cache[name] = &tool
	}
}

// FormatTools formats one or more templates to display on the console in
// table format or json format.
func (*Prm) FormatTools(tools map[string]*Tool, jsonOutput string) (string, error) {
	output := ""
	switch jsonOutput {
	case "table":
		count := len(tools)
		if count < 1 {
			return "", fmt.Errorf("could not locate any tools at %+v", viper.GetString("toolpath"))
		} else if count == 1 {
			stringBuilder := &strings.Builder{}
			for _, value := range tools {
				stringBuilder.WriteString(fmt.Sprintf("DisplayName:     %v\n", value.Cfg.Plugin.Display))
				stringBuilder.WriteString(fmt.Sprintf("Author:          %v\n", value.Cfg.Plugin.Author))
				stringBuilder.WriteString(fmt.Sprintf("Name:            %v\n", value.Cfg.Plugin.Id))
				stringBuilder.WriteString(fmt.Sprintf("Project_URL:     %v\n", value.Cfg.Plugin.UpstreamProjUrl))
				stringBuilder.WriteString(fmt.Sprintf("Version:         %v\n", value.Cfg.Plugin.Version))
			}
			output = stringBuilder.String()
		} else {
			stringBuilder := &strings.Builder{}
			table := tablewriter.NewWriter(stringBuilder)
			table.SetHeader([]string{"DisplayName", "Author", "Name", "Project_URL", "Version"})
			table.SetBorder(false)
			sortedTools := sortTools(tools)
			for _, value := range sortedTools {
				table.Append([]string{value.Cfg.Plugin.Display, value.Cfg.Plugin.Author, value.Cfg.Plugin.Id, value.Cfg.Plugin.UpstreamProjUrl, value.Cfg.Plugin.Version})
			}
			table.Render()
			output = stringBuilder.String()
		}
	case "json":
		j := jsoniter.ConfigFastest
		// This can't actually error because it's always getting a valid data struct;
		// if there are problems building the data struct for the template, we error
		// at that point instead.
		prettyJSON, _ := j.MarshalIndent(&tools, "", "  ")
		output = string(prettyJSON)
	}
	return output, nil
}

func sortTools(tools map[string]*Tool) []*Tool {
	var sortedTools []*Tool
	for _, tool := range tools {
		sortedTools = append(sortedTools, tool)
	}

	sort.Slice(sortedTools, func(i, j int) bool {
		return sortedTools[i].Cfg.Plugin.Display < sortedTools[j].Cfg.Plugin.Display
	})

	return sortedTools
}
