package main

import (
	"context"

	"github.com/puppetlabs/pdkgo/pkg/telemetry"
	"github.com/puppetlabs/prm/cmd/root"
	appver "github.com/puppetlabs/prm/cmd/version"
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

	var rootCmd = root.CreateRootCommand()

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
