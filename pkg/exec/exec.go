//nolint:structcheck,unused
package exec

import (
	"fmt"
	"github.com/puppetlabs/prm/pkg/backend"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/puppetlabs/prm/pkg/tool"
	"github.com/rs/zerolog/log"
)

type ExecExitCode int64

const (
	EXEC_SUCCESS ExecExitCode = iota
	EXEC_FAILURE
	EXEC_ERROR
)

type Exec struct {
	Prm           *prm.Prm
	RunningConfig config.Config
}

// Executes a tool with the given arguments, against the codeDir.
func (e *Exec) Exec(tool *tool.Tool, args []string) error {
	status := e.Prm.Backend.Status()
	if !status.IsAvailable {
		return fmt.Errorf("backend is unavailable: %s", status.StatusMessage)
	}

	// is the tool available?
	err := p.Backend.GetTool(tool, p.RunningConfig)
	if err != nil {
		log.Error().Msgf("Failed to exec tool: %s/%s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	}

	// the tool is available so execute against it
	exit, err := p.Backend.Exec(tool, args, p.RunningConfig, backend.DirectoryPaths{codeDir: p.CodeDir, cacheDir: p.CacheDir})
	if err != nil {
		log.Error().Msgf("Error executing tool %s/%s: %s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, err.Error())
		return err
	}

	switch exit {
	case tool.SUCCESS:
		log.Info().Msgf("Tool %s/%s executed successfully", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
	case tool.FAILURE:
		log.Error().Msgf("Tool %s/%s failed to execute", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	case tool.TOOL_ERROR:
		log.Error().Msgf("Tool %s/%s encountered an error", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	case tool.TOOL_NOT_FOUND:
		log.Error().Msgf("Tool %s/%s not found", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	default:
		log.Info().Msgf("Tool %s/%s exited with code %d", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, exit)
	}

	return nil
}
