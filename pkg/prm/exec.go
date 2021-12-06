//nolint:structcheck,unused
package prm

import "github.com/rs/zerolog/log"

type ExecExitCode int64

const (
	EXEC_SUCCESS ExecExitCode = iota
	EXEC_FAILURE
	EXEC_ERROR
)

// Executes a tool with the given arguments, against the codeDir.
func (p *Prm) Exec(tool *Tool, args []string) error {

	var backend BackendI

	switch RunningConfig.Backend {
	case DOCKER:
		backend = &Docker{}
	default:
		backend = &Docker{}
	}

	exit, err := backend.Exec(tool, args)

	if err != nil {
		log.Error().Msgf("Error executing tool %s/%s: %s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, err.Error())
		return err
	}

	switch exit {
	case SUCCESS:
		log.Info().Msgf("Tool %s/%s executed successfully", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		break
	case FAILURE:
		log.Error().Msgf("Tool %s/%s failed to execute", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		break
	case TOOL_ERROR:
		log.Error().Msgf("Tool %s/%s encountered an error", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		break
	case TOOL_NOT_FOUND:
		log.Error().Msgf("Tool %s/%s not found", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		break
	default:
		log.Info().Msgf("Tool %s/%s exited with code %d", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, exit)
		break
	}

	return nil
}
