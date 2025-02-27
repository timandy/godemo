package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/timandy/routinex/tools/consts"
	"github.com/timandy/routinex/tools/exec"
)

var goToolDir = exec.RunCmdOutput([]string{"go", "env", "GOTOOLDIR"})

func TestHelp(t *testing.T) {
	args := []string{"routinex", "-h", "/demo", "-p", "ttt", "go", "version"}
	os.Args = args
	main()
}

func TestOtherCmd(t *testing.T) {
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", "git", "version"}
	os.Args = args
	main()
}

func TestOtherCmdHelp(t *testing.T) {
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", "git", "-h"}
	os.Args = args
	main()
}

func TestCompileCmdHelp(t *testing.T) {
	compilePath := filepath.Join(goToolDir, consts.CompileName)
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", compilePath, "-h"}
	os.Args = args
	// expect exit 2
	// main()
}
