package goutils

import (
	"fmt"
	"os/exec"
	"time"
)

// GetBuildInfo will return git build information as a string. No GOVVV required.
func GetBuildInfo() (str string, err error) {
	var (
		buildDate,
		gitCommit,
		gitBranch,
		gitState string
		out []byte
		cmd *exec.Cmd
	)

	// Manually populate govvv gitBranch
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err = cmd.CombinedOutput()
	if err != nil {
		gitBranch = "unknown"
	}
	gitBranch = string(out)
	gitBranch, _ = StringToAlphaNumeric(gitBranch)

	// Manually populate govvv gitCommit
	cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err = cmd.CombinedOutput()
	if err != nil {
		gitCommit = "0000000"
	}
	gitCommit = string(out)
	gitCommit, _ = StringToAlphaNumeric(gitCommit)

	// Manually populate govvv gitState
	cmd = exec.Command("git", "diff", "--stat")
	out, err = cmd.CombinedOutput()
	if err != nil {
		gitState = "Unknown"
	}
	gitState = string(out)
	if len(gitState) > 0 {
		gitState = "Dirty"
	} else {
		gitState = "Clean"
	}

	// Manually populate govvv buildDate
	dt := time.Now()
	buildDate = dt.Format("01-02-2006 15:04:05")

	return fmt.Sprintf("%s-%s %s %s", gitBranch, gitCommit, gitState, buildDate), nil
}
