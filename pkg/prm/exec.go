//nolint:structcheck,unused
package prm

import (
	"github.com/rs/zerolog/log"
)

type ExecExitCode int64

const (
	EXEC_SUCCESS ExecExitCode = iota
	EXEC_FAILURE
	EXEC_ERROR
)

// Executes a tool with the given arguments, against the codeDir.
func (p *Prm) Exec(tool *Tool, args []string) error {
	if status := p.Backend.Status(); !status.IsAvailable {
		return ErrDockerNotRunning
	}

	// is the tool available?
	err := p.Backend.GetTool(tool)
	if err != nil {
		log.Error().Msgf("Failed to exec tool: %s/%s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	}

	// the tool is available so execute against it
	exit, err := p.Backend.Exec(tool, args, DirectoryPaths{codeDir: p.CodeDir, cacheDir: p.CacheDir})
	if err != nil {
		log.Error().Msgf("Error executing tool %s/%s: %s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, err.Error())
		return err
	}

	switch exit {
	case SUCCESS:
		log.Info().Msgf("Tool %s/%s executed successfully", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
	case FAILURE:
		log.Error().Msgf("Tool %s/%s failed to execute", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	case TOOL_ERROR:
		log.Error().Msgf("Tool %s/%s encountered an error", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	case TOOL_NOT_FOUND:
		log.Error().Msgf("Tool %s/%s not found", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	default:
		log.Info().Msgf("Tool %s/%s exited with code %d", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, exit)
	}

	return nil
}
