package exec

import (
	"os"
	"os/exec"
)

func RunCmd(args []string) {
	if len(args) == 0 {
		return
	}
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	startCmd(cmd)
	statusCode := waitCmd(cmd)
	os.Exit(statusCode)
}

func startCmd(cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

func waitCmd(cmd *exec.Cmd) int {
	// wait exit and return exit code
	err := cmd.Wait()
	if err == nil {
		return 0
	}
	if status, ok := err.(*exec.ExitError); ok {
		return status.ExitCode()
	}
	return 1
}
