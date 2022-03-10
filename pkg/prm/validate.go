//nolint:structcheck,unused
package prm

import "github.com/rs/zerolog/log"

type ValidateExitCode int64

const (
	VALIDATION_PASS ValidateExitCode = iota
	VALIDATION_FAILED
	VALIDATION_ERROR
)

// Validate allows a lits of tool names to be executed against
// the codeDir.
//
// Tools can be empty, in which case we expect that a local
// configuration file (validate.yml) will contain a list of
// tools to run.
func (p *Prm) Validate(tool *Tool, outputSettings OutputSettings) error {

	// is the tool available?
	err := p.Backend.GetTool(tool, p.RunningConfig)
	if err != nil {
		log.Error().Msgf("Failed to validate with tool: %s/%s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	}

	// the tool is available so execute against it
	exit, err := p.Backend.Validate(tool, p.RunningConfig, DirectoryPaths{codeDir: p.CodeDir, cacheDir: p.CacheDir}, outputSettings)

	switch exit {
	case VALIDATION_PASS:
		log.Info().Msgf("Tool %s/%s validated successfully", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		log.Info().Msg("PASS")
	case VALIDATION_FAILED:
		log.Error().Msgf("Tool %s/%s validation returned at least one failure", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		log.Info().Msg("FAIL")
	case VALIDATION_ERROR:
		log.Error().Msgf("Tool %s/%s encountered errored during validation %s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, err)
		log.Info().Msg("ERROR")
	default:
		log.Info().Msgf("Tool %s/%s exited with code %d", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, exit)
	}

	return err
}
