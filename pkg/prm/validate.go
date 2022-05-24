//nolint:structcheck,unused
package prm

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type ValidateExitCode int64

const (
	VALIDATION_PASS ValidateExitCode = iota
	VALIDATION_FAILED
	VALIDATION_ERROR
)

var (
	toolLogOutputPaths map[string]string // Key = toolName, Value = logFilePath, stores each tool's log file path
)

func (p *Prm) Validate(toolsInfo []ToolInfo, workerCount int, settings OutputSettings) error {
	if len(toolsInfo) == 0 {
		return fmt.Errorf("no tools provided for validation")
	}
	toolLogOutputPaths = make(map[string]string)

	tasks := p.createTasks(toolsInfo)

	pool := CreateWorkerPool(tasks, workerCount)
	pool.Run()

	err := p.outputResults(tasks, settings)
	return err
}

func (p Prm) taskFunc(tool ToolInfo) func() ValidationOutput {
	return func() ValidationOutput {
		toolName := tool.Tool.Cfg.Plugin.Id
		log.Info().Msgf("Validating with the %s tool", toolName)
		output := ValidationOutput{err: nil, exitCode: 0}

		err := p.Backend.GetTool(tool.Tool, p.RunningConfig)
		if err != nil {
			log.Error().Msgf("Failed to validate with tool: %s/%s", tool.Tool.Cfg.Plugin.Author, tool.Tool.Cfg.Plugin.Id)
			output = ValidationOutput{err: err, exitCode: VALIDATION_ERROR}
			return output
		}

		exitCode, stdout, err := p.Backend.Validate(tool, p.RunningConfig, DirectoryPaths{codeDir: p.CodeDir, cacheDir: p.CacheDir})
		if err != nil {
			output = ValidationOutput{err: err, exitCode: exitCode, stdout: stdout}
			return output
		}
		output.stdout = stdout
		return output
	}
}

func (p *Prm) outputResults(tasks []*Task[ValidationOutput], settings OutputSettings) error {
	err := p.writeOutputLogs(tasks, settings)
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

func writeOutputToTerminal(tasks []*Task[ValidationOutput]) {
	for _, task := range tasks {
		output := task.Output
		if output.err == nil {
			continue
		}

		var errText string
		if output.err.Error() != "" {
			errText = output.err.Error()
		} else {
			errText = output.stdout
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

func (p *Prm) createTasks(toolsInfo []ToolInfo) []*Task[ValidationOutput] {
	tasks := make([]*Task[ValidationOutput], len(toolsInfo))
	for i, info := range toolsInfo {
		tasks[i] = CreateTask[ValidationOutput](info.Tool.Cfg.Plugin.Id, p.taskFunc(info), ValidationOutput{})
	}
	return tasks
}

func (p *Prm) checkAndCreateDir(dir string) error {
	_, err := p.AFS.Stat(dir)
	if os.IsNotExist(err) {
		err = p.AFS.MkdirAll(dir, 0750)
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

func (p *Prm) writeOutputToFile(tasks []*Task[ValidationOutput], outputDir string) error {
	for _, task := range tasks {
		err := p.checkAndCreateDir(outputDir)
		if err != nil {
			return err
		}

		filePath := createLogFilePath(outputDir, task.Name)
		log.Debug().Msgf("output filepath: %v", filePath)

		file, err := p.AFS.Create(filePath)
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

func writeStringToFile(file afero.File, output ValidationOutput) error {
	errText := ""
	// Remove ANSI formatting from output strings
	if output.err != nil {
		errText = cleanOutput(output.err.Error())
	}
	stdout := cleanOutput(output.stdout)

	_, err := file.WriteString(fmt.Sprintf("%s\n%s", stdout, errText))
	if err != nil {
		return err
	}

	return nil
}

func (p *Prm) writeOutputLogs(tasks []*Task[ValidationOutput], settings OutputSettings) (err error) {
	if settings.ResultsView == "terminal" {
		writeOutputToTerminal(tasks)
		return nil
	}

	if settings.ResultsView == "file" {
		err := p.writeOutputToFile(tasks, settings.OutputDir)
		return err
	}

	return fmt.Errorf("invalid --resultsView flag specified")
}

func getErrorCount(tasks []*Task[ValidationOutput]) (count int) {
	for _, task := range tasks {
		output := task.Output
		if output.err != nil {
			count++
		}
	}
	return count
}

func createTableContents(tasks []*Task[ValidationOutput], resultsView string) (tableContents [][]string) {
	for _, task := range tasks {
		output := task.Output
		if resultsView == "file" { // Will also include the path to each
			outputPath := toolLogOutputPaths[task.Name]
			// Shortens the output file path so table doesn't become unreadable as a result of long file paths
			if shortOutputDir := strings.Split(outputPath, ".prm-validate"); len(shortOutputDir) == 2 {
				outputPath = fmt.Sprint(".prm-validate", shortOutputDir[1])
			}
			tableContents = append(tableContents, []string{task.Name, fmt.Sprintf("%d", output.exitCode), outputPath})
		} else {
			tableContents = append(tableContents, []string{task.Name, fmt.Sprintf("%d", output.exitCode)})
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
