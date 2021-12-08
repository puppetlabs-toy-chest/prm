package prm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
)

type Status struct {
	PuppetVersion *semver.Version
	Backend       BackendType
	BackendStatus BackendStatus
}

func (p *Prm) GetStatus() (status Status) {
	status.PuppetVersion = p.RunningConfig.PuppetVersion
	status.Backend = p.RunningConfig.Backend
	status.BackendStatus = p.Backend.Status()

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
			messageLines.WriteString(fmt.Sprintf("> %s\n", status.BackendStatus.StatusMsg))
		}
		statusMessage = messageLines.String()
	}
	return statusMessage, err
}
