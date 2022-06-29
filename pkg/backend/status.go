package backend

import (
	"encoding/json"
	"fmt"
	"github.com/puppetlabs/prm/pkg/config"
	"strings"

	"github.com/Masterminds/semver"
)

type Status struct {
	PuppetVersion *semver.Version
	Backend       config.BackendType
	BackendStatus BackendStatus
}

func GetStatus(cfg config.Config) Status {
	status := Status{
		PuppetVersion: cfg.PuppetVersion,
		Backend:       cfg.Backend,
		BackendStatus: back,
	}

	return status
}

func FormatStatus(status Status, outputType string) (statusMessage string, err error) {
	if outputType == "json" {
		jsonBytes, err := json.Marshal(status)
		if err == nil {
			statusMessage = string(jsonBytes)
		}
	} else {
		var messageLines strings.Builder
		messageLines.WriteString(fmt.Sprintf("> Puppet version: %s\n", status.PuppetVersion))
		if status.BackendStatus.IsAvailable {
			messageLines.WriteString(fmt.Sprintf("> Backend: %s (running)\n", status.Backend))
		} else {
			messageLines.WriteString(fmt.Sprintf("> Backend: %s (error)\n", status.Backend))
			messageLines.WriteString(fmt.Sprintf("> %s\n", status.BackendStatus.StatusMessage))
		}
		statusMessage = messageLines.String()
	}
	return statusMessage, err
}
