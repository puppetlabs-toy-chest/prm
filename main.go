package main

import (
	"context"
	"net/http"

	cmd_build "github.com/puppetlabs/pdkgo/cmd/build"
	"github.com/puppetlabs/pdkgo/pkg/build"
	"github.com/puppetlabs/pdkgo/pkg/exec_runner"
	"github.com/puppetlabs/pdkgo/pkg/gzip"
	"github.com/puppetlabs/pdkgo/pkg/install"
	"github.com/puppetlabs/pdkgo/pkg/tar"
	"github.com/puppetlabs/pdkgo/pkg/telemetry"
	"github.com/puppetlabs/prm/cmd/exec"
	"github.com/puppetlabs/prm/cmd/explain"
	"github.com/puppetlabs/prm/cmd/get"
	cmd_install "github.com/puppetlabs/prm/cmd/install"
	"github.com/puppetlabs/prm/cmd/root"
	"github.com/puppetlabs/prm/cmd/set"
	"github.com/puppetlabs/prm/cmd/status"
	"github.com/puppetlabs/prm/cmd/validate"
	appver "github.com/puppetlabs/prm/cmd/version"
	"github.com/puppetlabs/prm/internal/pkg/config_processor"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/puppetlabs/prm/pkg/utils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	version           = "dev"
	commit            = "none"
	date              = "unknown"
	honeycomb_api_key = "not_set"
	honeycomb_dataset = "not_set"
)

func main() {
	// Telemetry must be initialized before anything else;
	// If the telemetry build tag was not passed, this is all null ops
	ctx, traceProvider, parentSpan := telemetry.Start(context.Background(), honeycomb_api_key, honeycomb_dataset, "prm")

	// Create PRM context
	fs := afero.NewOsFs() // configure afero to use real filesystem
	prmApi := &prm.Prm{
		AFS:  &afero.Afero{Fs: fs},
		IOFS: &afero.IOFS{Fs: fs},
	}

	var rootCmd = root.CreateRootCommand(prmApi)

	// Get the command called and its arguments;
	// The arguments are only necessary if we want to
	// hand them off as an attribute to the parent span:
	// do we? Otherwise we just need the calledCommand
	calledCommand, calledCommandArguments := root.GetCalledCommand(rootCmd)
	telemetry.AddStringSpanAttribute(parentSpan, "arguments", calledCommandArguments)

	var verCmd = appver.CreateVersionCommand(version, date, commit)
	v := appver.Format(version, date, commit)
	rootCmd.Version = v
	rootCmd.SetVersionTemplate(v)
	rootCmd.AddCommand(verCmd)

	// set command
	sc := set.SetCommand{Utils: &utils.Utils{}}
	rootCmd.AddCommand(sc.CreateSetCommand())

	// get command
	rootCmd.AddCommand(get.CreateGetCommand(prmApi))

	// exec command
	rootCmd.AddCommand(exec.CreateCommand(prmApi))

	// validate command
	rootCmd.AddCommand(validate.CreateCommand(prmApi))

	// status command
	rootCmd.AddCommand(status.CreateStatusCommand(prmApi))

	// build
	buildCmd := cmd_build.BuildCommand{
		ProjectType: "tool",
		Builder: &build.Builder{
			Tar:  &tar.Tar{AFS: prmApi.AFS},
			Gzip: &gzip.Gzip{AFS: prmApi.AFS},
			AFS:  prmApi.AFS,
			ConfigProcessor: &config_processor.ConfigProcessor{
				AFS: prmApi.AFS,
			},
			ConfigFile: "prm-config.yml",
		},
	}
	rootCmd.AddCommand(buildCmd.CreateCommand())

	// install command
	installCmd := cmd_install.InstallCommand{
		PrmInstaller: &install.Installer{
			Tar:        &tar.Tar{AFS: prmApi.AFS},
			Gunzip:     &gzip.Gunzip{AFS: prmApi.AFS},
			AFS:        prmApi.AFS,
			IOFS:       prmApi.IOFS,
			HTTPClient: &http.Client{},
			Exec:       &exec_runner.Exec{},
			ConfigProcessor: &config_processor.ConfigProcessor{
				AFS: prmApi.AFS,
			},
		},
		AFS: prmApi.AFS,
	}
	rootCmd.AddCommand(installCmd.CreateCommand())

	// explain
	rootCmd.AddCommand(explain.CreateCommand())

	// initialize
	cobra.OnInitialize(root.InitLogger, root.InitConfig)

	// instrument & execute called command
	ctx, childSpan := telemetry.NewSpan(ctx, calledCommand)
	err := rootCmd.ExecuteContext(ctx)
	telemetry.RecordSpanError(childSpan, err)
	telemetry.EndSpan(childSpan)

	// Send all events
	telemetry.ShutDown(ctx, traceProvider, parentSpan)

	// Handle exiting with/out errors.
	cobra.CheckErr(err)
}
