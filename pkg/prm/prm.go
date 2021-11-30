// nolint:structcheck,unused
package prm

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/go-version"
	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	ToolConfigName     = "prm-config"
	ToolConfigFileName = "prm-config.yml"
)

type Prm struct {
	AFS     *afero.Afero
	IOFS    *afero.IOFS
	codeDir string
	cache   []toolCache
}

type toolCache struct {
	toolName string
	tool     *Tool
}
type PuppetVersion struct {
	version semver.Version
}

// Given a list of tool names, check if these are groups, and return
// an expanded list containing all the toolNames
func (*Prm) checkGroups(tools []string) []string {
	// TODO
	return []string{}
}

// Look within codeDir for a "validate.yml" containing
// a list of tools and/or tool groups that should be run against
// code within codeDir.
func (*Prm) checkLocalConfig() []string {
	// TODO
	return []string{}
}

// Check to see if the requested tool can be found installed.
// If installed read the tool configuration and return
func (*Prm) isToolAvailable(tool string) (Tool, bool) {
	return Tool{}, false
}

// Check to see if the tool is ready to execute
func (*Prm) isToolReady(tool *Tool) bool {
	return false
}

// save traversing to the filesystem
func (*Prm) cacheTool(tool *Tool) error {
	// TODO
	return nil
}

// What version of Puppet is requested by the user
func (*Prm) getPuppetVersion() PuppetVersion {
	return PuppetVersion{}
}

func (p *Prm) Get(toolDirPath string) (Tool, error) {
	file := filepath.Join(toolDirPath, ToolConfigFileName)
	_, err := p.AFS.Stat(file)
	if os.IsNotExist(err) {
		return Tool{}, fmt.Errorf("Couldn't find an installed tool at '%s'", toolDirPath)
	}
	i := p.readToolConfig(file)
	if reflect.DeepEqual(i, Tool{}) {
		return Tool{}, fmt.Errorf("Couldn't parse tool config at '%s'", file)
	}
	return i, nil
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
func (p *Prm) List(toolPath string, toolName string) []ToolConfig {
	log.Debug().Msgf("Searching %+v for tool configs", toolPath)
	// Triple glob to match author/id/version/ToolConfigFileName
	matches, _ := p.IOFS.Glob(toolPath + "/**/**/**/" + ToolConfigFileName)

	var tmpls []ToolConfig
	for _, file := range matches {
		log.Debug().Msgf("Found: %+v", file)
		i := p.readToolConfig(file)
		if i.Cfg.Plugin != nil {
			tmpls = append(tmpls, i.Cfg)
		}
	}

	if toolName != "" {
		log.Debug().Msgf("Filtering for: %s", toolName)
		tmpls = p.FilterFiles(tmpls, func(f ToolConfig) bool { return f.Plugin.Id == toolName })
	}

	tmpls = p.filterNewestVersions(tmpls)

	return tmpls
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

// FormatTools formats one or more templates to display on the console in
// table format or json format.
func (*Prm) FormatTools(tools []ToolConfig, jsonOutput string) (string, error) {
	output := ""
	switch jsonOutput {
	case "table":
		count := len(tools)
		if count < 1 {
			log.Warn().Msgf("Could not locate any tools at %+v", viper.GetString("toolpath"))
		} else if count == 1 {
			stringBuilder := &strings.Builder{}
			stringBuilder.WriteString(fmt.Sprintf("DisplayName:     %v\n", tools[0].Plugin.Display))
			stringBuilder.WriteString(fmt.Sprintf("Author:          %v\n", tools[0].Plugin.Author))
			stringBuilder.WriteString(fmt.Sprintf("Name:            %v\n", tools[0].Plugin.Id))
			stringBuilder.WriteString(fmt.Sprintf("ProjectURL:      %v\n", tools[0].Plugin.UpstreamProjUrl))
			stringBuilder.WriteString(fmt.Sprintf("Version:         %v\n", tools[0].Plugin.Version))
			output = stringBuilder.String()
		} else {
			stringBuilder := &strings.Builder{}
			table := tablewriter.NewWriter(stringBuilder)
			table.SetHeader([]string{"DisplayName", "Author", "Name", "ProjectURL", "Version"})
			table.SetBorder(false)
			for _, v := range tools {
				table.Append([]string{v.Plugin.Display, v.Plugin.Author, v.Plugin.Id, v.Plugin.UpstreamProjUrl, v.Plugin.Version})
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
