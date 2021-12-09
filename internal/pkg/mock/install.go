package mock

import (
	"fmt"
	"path/filepath"
)

type PctInstaller struct {
	ExpectedToolPkg   string
	ExpectedTargetDir string
	ExpectedGitUri    string
}

func (p *PctInstaller) Install(templatePkg, targetDir string, force bool) (string, error) {
	if templatePkg != p.ExpectedToolPkg {
		return "", fmt.Errorf("templatePkg (%v) did not match expected value (%v)", templatePkg, p.ExpectedToolPkg)
	}

	if targetDir != p.ExpectedTargetDir {
		return "", fmt.Errorf("targetDir (%v) did not match expected value (%v)", targetDir, p.ExpectedTargetDir)
	}

	return filepath.Clean("/unit/test/path"), nil
}

func (p *PctInstaller) InstallClone(gitUri, targetDir, tempDir string, force bool) (string, error) {
	if gitUri != p.ExpectedGitUri {
		return "", fmt.Errorf("gitUri (%v) did not match expected value (%v)", gitUri, p.ExpectedGitUri)
	}

	if tempDir == "" {
		return "", fmt.Errorf("tempDir was an empty string")
	}

	if targetDir != p.ExpectedTargetDir {
		return "", fmt.Errorf("targetDir (%v) did not match expected value (%v)", targetDir, p.ExpectedTargetDir)
	}

	return filepath.Clean("/unit/test/path"), nil
}
