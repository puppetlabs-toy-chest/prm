package install

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/puppetlabs/pdkgo/pkg/install"
	"github.com/puppetlabs/pdkgo/pkg/telemetry"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type InstallCommand struct {
	ToolPkgPath  string
	InstallPath  string
	Force        bool
	PrmInstaller install.InstallerI
	GitUri       string
	AFS          *afero.Afero
}

type InstallCommandI interface {
	CreateCommand() *cobra.Command
}

func (ic *InstallCommand) CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "install <tool.tar.gz|uri> [flags]",
		Short:   "Installs a tool (in tar.gz format)",
		Long:    `Installs a tool (in tar.gz format) to the default or specified tool path`,
		PreRunE: ic.preExecute,
		RunE:    ic.executeInstall,
	}
	tmp.Flags().StringVar(&ic.InstallPath, "toolpath", "", "location of installed tools")
	err := viper.BindPFlag("toolpath", tmp.Flags().Lookup("toolpath"))
	tmp.Flags().BoolVarP(&ic.Force, "force", "f", false, "Forces the install of a tool without error, if it already exists. ")
	tmp.Flags().StringVar(&ic.GitUri, "git-uri", "", "Installs a tool package from a remote git repository.")

	cobra.CheckErr(err)

	return tmp
}

func (ic *InstallCommand) executeInstall(cmd *cobra.Command, args []string) error {
	_, span := telemetry.NewSpan(cmd.Context(), "install")
	defer telemetry.EndSpan(span)
	telemetry.AddStringSpanAttribute(span, "name", "install")

	toolInstallationPath := ""
	var err error = nil
	if ic.GitUri != "" { // For cloning a tool
		// Create temp folder
		tempDir, dirErr := ic.AFS.TempDir("", "")
		defer func() {
			dirErr := ic.AFS.Remove(tempDir)
			if dirErr != nil {
				log.Error().Msgf("Failed to remove temp dir: %v", dirErr)
			}
		}()
		if dirErr != nil {
			return fmt.Errorf("Could not create tempdir to clone tool to: %v", err)
		}
		toolInstallationPath, err = ic.PrmInstaller.InstallClone(ic.GitUri, ic.InstallPath, tempDir, ic.Force)
	} else { // For downloading and/or locally installing a tool
		toolInstallationPath, err = ic.PrmInstaller.Install(ic.ToolPkgPath, ic.InstallPath, ic.Force)
	}

	if err != nil {
		return err
	}
	log.Info().Msgf("Tool installed to %v", toolInstallationPath)
	return nil
}

func (ic *InstallCommand) setInstallPath() error {
	if ic.InstallPath == "" {
		defaultToolPath := viper.GetString(prm.ToolPathCfgKey)
		if defaultToolPath == "" {
			return fmt.Errorf("Could not determine location to install tool") //: %v", err)
		}
		ic.InstallPath = defaultToolPath
	}
	return nil
}

func (ic *InstallCommand) preExecute(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		if ic.GitUri != "" {
			return ic.setInstallPath()
		}
		return fmt.Errorf("Path to tool package (tar.gz) should be first argument")
	}

	if len(args) == 1 {
		ic.ToolPkgPath = args[0]
		return ic.setInstallPath()
	}

	if len(args) > 1 {
		return fmt.Errorf("Incorrect number of arguments; path to tool package (tar.gz) should be first argument")
	}

	return nil
}
