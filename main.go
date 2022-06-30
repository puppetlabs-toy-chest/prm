package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/puppetlabs/pct/pkg/telemetry"
	"github.com/puppetlabs/prm/cmd/root"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	honeycomb_api_key = "not_set"
	honeycomb_dataset = "not_set"
)

func main() {
	// Telemetry must be initialized before anything else;
	// If the telemetry build tag was not passed, this is all null ops
	ctx, traceProvider, parentSpan := telemetry.Start(context.Background(), honeycomb_api_key, honeycomb_dataset, "prm")

	// Get the command called and its arguments;
	// The arguments are only necessary if we want to
	// hand them off as an attribute to the parent span:
	// do we? Otherwise we just need the calledCommand
	calledCommand, calledCommandArguments := root.GetCalledCommand(rootCmd)
	telemetry.AddStringSpanAttribute(parentSpan, "arguments", calledCommandArguments)

	// initialize
	cobra.OnInitialize(logger.InitLogger, config.InitConfig)

	// instrument & execute called command
	ctx, childSpan := telemetry.NewSpan(ctx, calledCommand)
	err := rootCmd.ExecuteContext(ctx)
	telemetry.RecordSpanError(childSpan, err)
	telemetry.EndSpan(childSpan)

	// Send all events
	telemetry.ShutDown(ctx, traceProvider, parentSpan)

	// Handle exiting with/out errors.
	//cobra.CheckErr(err)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}
