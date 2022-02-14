package config_processor

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type ConfigProcessor struct {
	AFS *afero.Afero
}

func (p *ConfigProcessor) ProcessConfig(sourceDir, targetDir string, force bool) (string, error) {
	// Read config to determine tool properties
	info, err := p.readConfig(filepath.Join(sourceDir, "prm-config.yml"))
	if err != nil {
		return "", fmt.Errorf("Invalid config: %v", err.Error())
	}

	// Create namespaced directory and move contents of temp folder to it
	namespacedPath, err := p.setupToolNamespace(targetDir, info, sourceDir, force)
	if err != nil {
		return "", fmt.Errorf("Unable to install in namespace: %v", err.Error())
	}
	return namespacedPath, nil
}

func (p *ConfigProcessor) CheckConfig(configFile string) error {
	info, err := p.readConfig(configFile)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("The following attributes are missing in %s:\n", configFile)
	orig := msg
	// These parts are essential for build and deployment.

	if info.Plugin.Id == "" {
		msg = msg + "  * id\n"
	}
	if info.Plugin.Author == "" {
		msg = msg + "  * author\n"
	}
	if info.Plugin.Version == "" {
		msg = msg + "  * version\n"
	}
	if info.Plugin.UpstreamProjUrl == "" {
		msg = msg + "  * upstream project url\n"
	}
	if info.Plugin.Display == "" {
		msg = msg + "  * display name\n"
	}
	if msg != orig {
		return fmt.Errorf(msg)
	}

	return nil
}

func (p *ConfigProcessor) readConfig(configFile string) (info prm.ToolConfigInfo, err error) {
	fileBytes, err := p.AFS.ReadFile(configFile)
	if err != nil {
		return info, err
	}

	// use viper to parse the config as it knows how to work with mapstructure squash
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(fileBytes))
	if err != nil {
		return info, err
	}

	err = viper.Unmarshal(&info)
	if err != nil {
		return info, err
	}

	return info, err
}

func (p *ConfigProcessor) setupToolNamespace(targetDir string, info prm.ToolConfigInfo, untarPath string, force bool) (string, error) {
	// author/id/version
	toolPath := filepath.Join(targetDir, info.Plugin.Author, info.Plugin.Id)

	err := p.AFS.MkdirAll(toolPath, 0750)
	if err != nil {
		return "", err
	}

	namespacePath := filepath.Join(targetDir, info.Plugin.Author, info.Plugin.Id, info.Plugin.Version)

	// finally move to the full path
	err = p.AFS.Rename(untarPath, namespacePath)
	if err != nil {
		// if a tool already exists
		if !force {
			// error unless forced
			return "", fmt.Errorf("Tool already installed (%s)", namespacePath)
		} else {
			// remove the exiting tool
			err = p.AFS.RemoveAll(namespacePath)
			if err != nil {
				return "", fmt.Errorf("Unable to overwrite existing tool: %v", err)
			}
			// perform the move again
			err = p.AFS.Rename(untarPath, namespacePath)
			if err != nil {
				return "", fmt.Errorf("Unable to force install: %v", err)
			}
		}
	}

	return namespacePath, err
}
