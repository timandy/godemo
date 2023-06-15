package instrument

import (
	"go/parser"
	"go/token"

	"github.com/timandy/routiner/instrument/api"
	"github.com/timandy/routiner/instrument/instruments"
	"github.com/timandy/routiner/tools/flag"
	"github.com/timandy/routiner/tools/stringutil"
)

var defaults = []api.Instrument{instruments.NewRuntimeInstrument()}

func Execute(args []string) []string {
	// resolve options
	options := resolveCompileOptions(args)
	if options == nil {
		return args
	}
	// exec instruments and return new args
	return execute(options.Clone())
}

func execute(options *api.CompileOptions) []string {
	for _, ins := range defaults {
		asmHdrIdx := stringutil.LastIndexOf(options.Args, "-asmhdr")
		if asmHdrIdx == -1 {
			return options.Args
		}
		execute0(ins, options, asmHdrIdx)
	}
	return options.Args
}

func execute0(ins api.Instrument, options *api.CompileOptions, asmHdrIdx int) {
	// define result
	result := api.NewInstrumentResult()
	// proc args after exec
	defer func() {
		for idx, path := range result.ReplaceFiles {
			options.Args[idx] = path
		}
		options.Args = append(options.Args, result.ExtraFiles...)
	}()
	// verify this ins can handle the package
	if !ins.PreHandlePackage(options, result) {
		return
	}
	for idx, length := asmHdrIdx+1, len(options.Args); idx < length; idx++ {
		path := options.Args[idx]
		if !ins.PreHandleFile(path, idx, options, result) {
			continue
		}
		// parse the ast file
		fset := token.NewFileSet()
		af, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		if !ins.HandleFile(path, idx, fset, af, options, result) {
			continue
		}
		ins.PostHandleFile(path, idx, fset, af, options, result)
	}
	ins.PostHandlePackage(options, result)
}

func resolveCompileOptions(args []string) *api.CompileOptions {
	options := resolveCompileOptions0(args)
	if options != nil {
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
