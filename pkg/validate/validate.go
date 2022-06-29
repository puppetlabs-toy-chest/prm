//nolint:structcheck,unused
package validate

import (
	"errors"
	"fmt"
	"github.com/puppetlabs/prm/pkg/backend"
	"github.com/puppetlabs/prm/pkg/backend/docker"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/utils"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

var (
	toolLogOutputPaths map[string]string // Key = toolName, Value = logFilePath, stores each tool's log file path
)

type Validator struct {
	Backend        backend.BackendI
	AFS            *afero.Afero
	DirectoryPaths backend.DirectoryPaths
	RunningConfig  config.Config
}

func (v *Validator) Validate(toolsInfo []backend.ToolInfo, workerCount int, settings backend.OutputSettings) error {
	if status := v.Backend.Status(); !status.IsAvailable {
		return docker.ErrDockerNotRunning
	}

	if len(toolsInfo) == 0 {
		return fmt.Errorf("no tools provided for validation")
	}
	toolLogOutputPaths = make(map[string]string)

	tasks := v.createTasks(toolsInfo)

	pool := utils.CreateWorkerPool(tasks, workerCount)
	pool.Run()

	err := v.outputResults(tasks, settings)
	return err
}

func (v Validator) taskFunc(tool backend.ToolInfo) func() backend.ValidationOutput {
	return func() backend.ValidationOutput {
		toolName := tool.Tool.Cfg.Plugin.Id
		log.Info().Msgf("Validating with the %s tool", toolName)
		output := backend.ValidationOutput{Err: nil, ExitCode: 0}

		err := v.Backend.GetTool(tool.Tool, v.RunningConfig)
		if err != nil {
			log.Error().Msgf("Failed to validate with tool: %s/%s", tool.Tool.Cfg.Plugin.Author, tool.Tool.Cfg.Plugin.Id)
			output = backend.ValidationOutput{Err: err, ExitCode: backend.VALIDATION_ERROR}
			return output
		}

		exitCode, stdout, err := v.Backend.Validate(tool, v.RunningConfig, v.DirectoryPaths)
		if err != nil {
			output = backend.ValidationOutput{Err: err, ExitCode: exitCode, Stdout: stdout}
			return output
		}
		output.Stdout = stdout
		return output
	}
}

func (v *Validator) outputResults(tasks []*utils.Task[backend.ValidationOutput], settings backend.OutputSettings) error {
	err := v.writeOutputLogs(tasks, settings)
	if err != nil {
		return err
	}

	tableContents := createTableContents(tasks, settings.ResultsView)
	headers := []string{"Tool Name", "Validation Exit Code"}
	if settings.ResultsView == "file" {
		headers = append(headers, "File Location")
	}
	renderTable(headers, tableContents)

	if errorCount := getErrorCount(tasks); errorCount > 0 {
		return errors.New(getErrorMessage(errorCount))
	}

	return nil
}

func writeOutputToTerminal(tasks []*utils.Task[backend.ValidationOutput]) {
	for _, task := range tasks {
		output := task.Output
		if output.Err == nil {
			continue
		}

		var errText string
		if output.Err.Error() != "" {
			errText = output.Err.Error()
		} else {
			errText = output.Stdout
		}
		errText = cleanOutput(errText)

		log.Error().Msgf(fmt.Sprintf("%s:\n%s", task.Name, cleanOutput(errText)))
	}
}

func cleanOutput(text string) string {
	exp := regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	text = exp.ReplaceAllString(text, "")
	text = strings.TrimPrefix(text, "\n") // Trim prefix newline if it exists
	text = strings.TrimSuffix(text, "\n")
	text = strings.ReplaceAll(text, "/code/", "")
	return text
}

func renderTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.AppendBulk(data)
	fmt.Println()
	table.Render()
}

func (v *Validator) createTasks(toolsInfo []backend.ToolInfo) []*utils.Task[backend.ValidationOutput] {
	tasks := make([]*utils.Task[backend.ValidationOutput], len(toolsInfo))
	for i, info := range toolsInfo {
		tasks[i] = utils.CreateTask[backend.ValidationOutput](info.Tool.Cfg.Plugin.Id, v.taskFunc(info), backend.ValidationOutput{})
	}
	return tasks
}

func (v *Validator) checkAndCreateDir(dir string) error {
	_, err := v.AFS.Stat(dir)
	if os.IsNotExist(err) {
		err = v.AFS.MkdirAll(dir, 0750)
		return err
	} else if err != nil {
		return err
	}

	return nil
}

func createLogFilePath(outputDir string, toolId string) string {
	timeNow := time.Now()
	fileName := fmt.Sprintf("%v_%v_%v_%v_%v-%v-%v.log", toolId, timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())
	fullPath := path.Join(outputDir, fileName)
	toolLogOutputPaths[toolId] = fullPath

	return fullPath
}

func (v *Validator) writeOutputToFile(tasks []*utils.Task[backend.ValidationOutput], outputDir string) error {
	for _, task := range tasks {
		err := v.checkAndCreateDir(outputDir)
		if err != nil {
			return err
		}

		filePath := createLogFilePath(outputDir, task.Name)
		log.Debug().Msgf("output filepath: %v", filePath)

		file, err := v.AFS.Create(filePath)
		if err != nil {
			return err
		}

		err = writeStringToFile(file, task.Output)
		if err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			log.Error().Msgf("Error closing file: %s", err)
		}
	}

	return nil
}

func writeStringToFile(file afero.File, output backend.ValidationOutput) error {
	errText := ""
	// Remove ANSI formatting from output strings
	if output.Err != nil {
		errText = cleanOutput(output.Err.Error())
	}
	stdout := cleanOutput(output.Stdout)

	_, err := file.WriteString(fmt.Sprintf("%s\n%s", stdout, errText))
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) writeOutputLogs(tasks []*utils.Task[backend.ValidationOutput], settings backend.OutputSettings) (err error) {
	if settings.ResultsView == "terminal" {
		writeOutputToTerminal(tasks)
		return nil
	}

	if settings.ResultsView == "file" {
		err := v.writeOutputToFile(tasks, settings.OutputDir)
		return err
	}

	return fmt.Errorf("invalid --resultsView flag specified")
}

func getErrorCount(tasks []*utils.Task[backend.ValidationOutput]) (count int) {
	for _, task := range tasks {
		output := task.Output
		if output.Err != nil {
			count++
		}
	}
	return count
}

func createTableContents(tasks []*utils.Task[backend.ValidationOutput], resultsView string) (tableContents [][]string) {
	for _, task := range tasks {
		output := task.Output
		if resultsView == "file" { // Will also include the path to each
			outputPath := toolLogOutputPaths[task.Name]
			// Shortens the output file path so table doesn't become unreadable as a result of long file paths
			if shortOutputDir := strings.Split(outputPath, ".prm-validate"); len(shortOutputDir) == 2 {
				outputPath = fmt.Sprint(".prm-validate", shortOutputDir[1])
			}
			tableContents = append(tableContents, []string{task.Name, fmt.Sprintf("%d", output.ExitCode), outputPath})
		} else {
			tableContents = append(tableContents, []string{task.Name, fmt.Sprintf("%d", output.ExitCode)})
		}
	}
	return tableContents
}

func getErrorMessage(count int) string {
	spelling := "errors"
	if count == 1 {
		spelling = "error"
	}

	return fmt.Sprintf("Validation returned %d %s", count, spelling)
}
