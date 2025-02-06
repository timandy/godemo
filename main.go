package main

import (
	"os"

	"github.com/timandy/routiner/inject"
	"github.com/timandy/routiner/tools/exec"
	"github.com/timandy/routiner/tools/flag"
	"github.com/timandy/routiner/tools/log"
	"github.com/timandy/routiner/tools/opt"
)

func main() {
	// parse routine compile options
	args := os.Args[1:]
	appOpt := &opt.AppOptions{}
	flagSet := flag.ParseStruct(appOpt, os.Args[0], args)
	if appOpt.Debug {
		log.PrintArgs("entry", os.Args)
	}
	// print usage
	if appOpt.Help {
		flagSet.SortFlags = false
		flag.PrintUsage(flagSet)
		return
	}
	// no remained args, do nothing and return
	remainArgs := flagSet.Args()
	if appOpt.Debug {
		log.PrintArgs("before", remainArgs)
	}
	// exists remained args, run the cmd finally
	defer func() {
		if appOpt.Debug {
			log.PrintArgs("after", remainArgs)
		}
		exec.RunCmd(remainArgs)
	}()
	// exec inject
	remainArgs = inject.Execute(remainArgs, appOpt)
}
