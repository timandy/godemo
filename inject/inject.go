package inject

import (
	"go/parser"
	"go/token"
	"strings"

	"github.com/timandy/routiner/inject/api"
	"github.com/timandy/routiner/inject/injectors"
	"github.com/timandy/routiner/tools/flag"
	"github.com/timandy/routiner/tools/log"
	"github.com/timandy/routiner/tools/opt"
	"github.com/timandy/routiner/tools/slices"
	"github.com/timandy/routiner/tools/stringutil"
)

var defaults = []api.Injector{injectors.NewRuntimeInjector(), injectors.NewRoutineXInjector()}

func Execute(args []string, app *opt.AppOptions) []string {
	// resolve options
	options := resolveCompileOptions(args, app)
	if options == nil {
		return args
	}
	if options.Debug {
		log.PrintArg("workdir", options.WorkDir())
	}
	// exec injectors and return new args
	return execute(options.Clone())
}

func execute(options *api.CompileOptions) []string {
	for _, injector := range defaults {
		asmHdrIdx := stringutil.LastIndexOf(options.Args, "-asmhdr")
		if asmHdrIdx == -1 {
			return options.Args
		}
		execute0(injector, options, asmHdrIdx)
	}
	return options.Args
}

func execute0(injector api.Injector, options *api.CompileOptions, asmHdrIdx int) {
	// define result
	result := api.NewInjectResult()
	// proc args after exec
	defer func() {
		for idx, path := range result.ReplaceFiles {
			options.Args[idx] = path
		}
		args := append(options.Args, result.ExtraFiles...)
		args = slices.Filter(args, func(str string) bool { return str != "" })
		options.Args = args
	}()
	// verify this injector can handle the package
	if !injector.PreHandlePackage(options, result) {
		return
	}
	for idx, length := asmHdrIdx+1, len(options.Args); idx < length; idx++ {
		path := options.Args[idx]
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		if !injector.PreHandleFile(path, idx, options, result) {
			continue
		}
		// parse the ast file
		fset := token.NewFileSet()
		af, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		if !injector.HandleFile(path, idx, fset, af, options, result) {
			continue
		}
		injector.PostHandleFile(path, idx, fset, af, options, result)
	}
	injector.PostHandlePackage(options, result)
}

func resolveCompileOptions(args []string, app *opt.AppOptions) *api.CompileOptions {
	options := resolveCompileOptions0(args)
	if options != nil {
		options.Debug = app.Debug
		options.Verbose = app.Verbose
		options.Args = args
	}
	return options
}

func resolveCompileOptions0(args []string) *api.CompileOptions {
	if len(args) == 0 {
		return nil
	}
	options := &api.CompileOptions{}
	flagSet := flag.ParseStruct(options, args[0], args[1:])
	if options.IsValid(flagSet.Name()) {
		return options
	}
	remainArgs := flagSet.Args()
	return resolveCompileOptions0(remainArgs)
}
