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

type ValidateYmlContent struct {
	Tools []string `yaml:"tools"`
}

// checkGroups takes a slice of tool names and iterates through each
// checking against a map of toolGroups. If a toolGroup name is found
// the toolGroup is expanded and the list of tools is updated.
func (*Prm) checkGroups(tools []string) []string {
	for index, toolName := range tools {
		if toolGroup, ok := ToolGroups[toolName]; ok {
			// remove the group from the list
			tools = append(tools[:index], tools[index+1:]...)
			// add the expanded toolgroup to the list
			tools = append(tools, toolGroup...)
		}
	}

	// remove duplicates
	allKeys := make(map[string]bool)
	clean := []string{}
	for _, item := range tools {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			clean = append(clean, item)
		}
	}

	return clean
}

// Check if the code being executed against contains a yalidate.yml. Read in the
// list of tool names from validate.yml into a list. Pass the list of
// tool names to flattenToolList to expand out any groups. Then return
// the complete list.
func (p *Prm) CheckLocalConfig() ([]string, error) {
	// check if validate.yml exits in the codeDir
	validateFile := filepath.Join(p.CodeDir, "validate.yml")
	if _, err := p.AFS.Stat(validateFile); err != nil {
		log.Error().Msgf("validate.yml not found in %s", p.CodeDir)
		return []string{}, err
	}

	// read in validate.yml
	contents, err := p.AFS.ReadFile(validateFile)
	if err != nil {
		log.Error().Msgf("Error reading validate.yml: %s", err)
		return []string{}, err
	}

	// parse validate.yml to our temporary struct
	var userList ValidateYmlContent
	err = yaml.Unmarshal(contents, &userList)
	if err != nil {
		log.Error().Msgf("validate.yml is not formated correctly: %s", err)
		return []string{}, err
	}

	return p.checkGroups(userList.Tools), nil
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
func (p *Prm) List(toolPath string, toolName string) {
	log.Debug().Msgf("Searching %+v for tool configs", toolPath)
	// Triple glob to match author/id/version/ToolConfigFileName
	matches, _ := p.IOFS.Glob(toolPath + "/**/**/**/" + ToolConfigFileName)

	var tmpls []ToolConfig
	for _, file := range matches {
		log.Debug().Msgf("Found: %+v", file)
		i := p.readToolConfig(file)
		if i.Cfg.Plugin != nil {
			i.Cfg.Path = filepath.Dir(file)
			tmpls = append(tmpls, i.Cfg)
		}
	}

	if toolName != "" {
		log.Debug().Msgf("Filtering for: %s", toolName)
		tmpls = p.FilterFiles(tmpls, func(f ToolConfig) bool { return f.Plugin.Id == toolName })
	}

	tmpls = p.filterNewestVersions(tmpls)

	// cache for use with the rest of the program
	// this is a seperate cache from the one used by the CLI
	p.createToolCache(tmpls)
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
			log.Warn().Msgf("Could not locate any tools at %+v", viper.GetString("toolpath"))
		} else if count == 1 {
			stringBuilder := &strings.Builder{}
			for key := range tools {
				stringBuilder.WriteString(fmt.Sprintf("DisplayName:     %v\n", tools[key].Cfg.Plugin.Display))
				stringBuilder.WriteString(fmt.Sprintf("Author:          %v\n", tools[key].Cfg.Plugin.Author))
				stringBuilder.WriteString(fmt.Sprintf("Name:            %v\n", tools[key].Cfg.Plugin.Id))
				stringBuilder.WriteString(fmt.Sprintf("Project_URL:     %v\n", tools[key].Cfg.Plugin.UpstreamProjUrl))
				stringBuilder.WriteString(fmt.Sprintf("Version:         %v\n", tools[key].Cfg.Plugin.Version))
			}
			output = stringBuilder.String()
		} else {
			stringBuilder := &strings.Builder{}
			table := tablewriter.NewWriter(stringBuilder)
			table.SetHeader([]string{"DisplayName", "Author", "Name", "Project_URL", "Version"})
			table.SetBorder(false)
			for key := range tools {
				table.Append([]string{tools[key].Cfg.Plugin.Display, tools[key].Cfg.Plugin.Author, tools[key].Cfg.Plugin.Id, tools[key].Cfg.Plugin.UpstreamProjUrl, tools[key].Cfg.Plugin.Version})
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
