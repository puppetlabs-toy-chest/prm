package set

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/puppetlabs/prm/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createSetPuppetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "puppet <VERSION>",
		Short: "Sets the Puppet runtime to the specified version",
		Long:  `Sets the Puppet runtime to the specified version`,
		RunE:  setPuppetVersion,
	}

	return tmp
}

func setPuppetVersion(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("only a single Puppet version can be set")
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a Puppet version after 'set puppet'")
	}

	puppetSemVer, err := semver.NewVersion(args[0])
	if err != nil {
		return fmt.Errorf("'%s' is not a semantic (x.y.z) Puppet version: %s", args[0], err)
	}

	viper.Set(prm.PuppetVerCfgKey, puppetSemVer.String())
	utils.WriteConfig()

	return err
}

// TODO: (GH-26) Consume a list of available Puppet versions to faciliate tab completion
// on command line
