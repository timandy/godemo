package main

import (
	"os"

	"github.com/timandy/routiner/instrument"
	"github.com/timandy/routiner/tools/exec"
	"github.com/timandy/routiner/tools/flag"
)

func main() {
	// parse routine compile options
	args := os.Args[1:]
	routineCompileOptions := &RoutineCompileOptions{}
	flagSet := flag.ParseStruct(routineCompileOptions, os.Args[0], args)
	// print usage
	if routineCompileOptions.Help {
		flagSet.Usage()
		return
	}
	// no remained args, do nothing and return
	remainArgs := flagSet.Args()
	// exists remained args, run the cmd finally
	defer func() {
		exec.RunCmd(remainArgs)
	}()
	// exec inject
	remainArgs = instrument.Execute(remainArgs)
}

type RoutineCompileOptions struct {
	Help bool `name:"help" shorthand:"h" usage:"show this message and exit"`
}
