package config_processor

import (
	"bytes"
	"fmt"
	"github.com/puppetlabs/prm/pkg/tool"

	"github.com/puppetlabs/pct/pkg/config_processor"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type ConfigProcessor struct {
	AFS *afero.Afero
}

func (p *ConfigProcessor) GetConfigMetadata(configFile string) (metadata config_processor.ConfigMetadata, err error) {
	configInfo, err := p.ReadConfig(configFile)
	if err != nil {
		return metadata, err
	}

	err = p.CheckConfig(configFile)
	if err != nil {
		return metadata, err
	}

	metadata = config_processor.ConfigMetadata{
		Author:  configInfo.Plugin.Author,
		Id:      configInfo.Plugin.Id,
		Version: configInfo.Plugin.Version,
	}
	return metadata, nil
}

func (p *ConfigProcessor) CheckConfig(configFile string) error {
	info, err := p.ReadConfig(configFile)
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
	if msg != orig {
		return fmt.Errorf(msg)
	}

	return nil
}

func (p *ConfigProcessor) ReadConfig(configFile string) (info tool.ToolConfigInfo, err error) {
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
