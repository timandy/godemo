package instruments

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/timandy/routiner/instrument/api"
	"github.com/timandy/routiner/tools/astutil"
	"github.com/timandy/routiner/tools/log"
	"github.com/timandy/routiner/tools/os"
	"github.com/timandy/routiner/tools/stringutil"
)

type RuntimeInstrument struct {
}

func NewRuntimeInstrument() api.Instrument {
	return &RuntimeInstrument{}
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInstrument) PreHandlePackage(options *api.CompileOptions, result *api.InstrumentResult) bool {
	return options.Package == "runtime"
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInstrument) PreHandleFile(path string, idx int, options *api.CompileOptions, result *api.InstrumentResult) bool {
	return strings.HasSuffix(path, "runtime2.go") || strings.HasSuffix(path, "proc.go")
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInstrument) HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *api.CompileOptions, result *api.InstrumentResult) bool {
	handled := false
	ast.Inspect(af, func(node ast.Node) bool {
		if r.handleNode(node, options) {
			handled = true
			return false
		}
		return true
	})
	return handled
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInstrument) PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *api.CompileOptions, result *api.InstrumentResult) {
	srcShortName := filepath.Base(path)
	destPath := filepath.Join(options.WorkDir(), srcShortName)
	astutil.SaveAs(destPath, fset, af)
	result.ReplaceFiles[idx] = destPath
}

//goland:noinspection GoUnusedParameter
func (r *RuntimeInstrument) PostHandlePackage(options *api.CompileOptions, result *api.InstrumentResult) {
	code := stringutil.ExecuteTemplate(`package runtime

import _ "unsafe"

//go:nosplit
//go:linkname getg0
func getg0() interface{} {
	return g0
}

//go:nosplit
//go:linkname getgp
func getgp() *g {
	return getg()
}
`, nil)
	// save file
	destShortName := "runtime_routine.go"
	destPath := filepath.Join(options.WorkDir(), destShortName)
	os.WriteFile(destPath, code)
	result.ExtraFiles = append(result.ExtraFiles, destPath)
	if options.Debug || options.Verbose {
		log.Info("create function 'runtime.getg0' and 'runtime.getgp'")
	}
}

func (r *RuntimeInstrument) handleNode(node ast.Node, options *api.CompileOptions) bool {
	switch n := node.(type) {
	case *ast.TypeSpec:
		ident := n.Name
		if ident == nil || ident.Name != "g" {
			return false
		}
		st, isSt := n.Type.(*ast.StructType)
		if !isSt {
			return false
		}
		fields := st.Fields
		if fields == nil {
			return false
		}
		fieldList := fields.List
		if len(fieldList) == 0 {
			return false
		}
		threadLocalsField := astutil.CreateField("threadLocals", "unsafe.Pointer")
		inheritableThreadLocalsField := astutil.CreateField("inheritableThreadLocals", "unsafe.Pointer")
		fields.List = append(fieldList, threadLocalsField, inheritableThreadLocalsField)
		if options.Debug || options.Verbose {
			log.Info("enhance struct 'runtime.g' add field 'threadLocals unsafe.Pointer'")
			log.Info("enhance struct 'runtime.g' add field 'inheritableThreadLocals unsafe.Pointer'")
		}
		return true
	case *ast.FuncDecl:
		// check name
		ident := n.Name
		if ident == nil || ident.Name != "goexit0" {
			return false
		}
		// check type not nil
		funcType := n.Type
		if funcType == nil {
			return false
		}
		// check no results
		results := funcType.Results
		if results != nil && len(results.List) > 0 {
			return false
		}
		// check only one params
		params := funcType.Params
		if params == nil || len(params.List) != 1 {
			return false
		}
		//
		body := n.Body
		if body == nil {
			return false
		}
		list := body.List
		if len(list) == 0 {
			return false
		}
		x, index := astutil.IndexAssignTimerNil(list)
		if index == -1 {
			return false
		}
		threadLocalsStmt := astutil.CreateAssignNilStmt(x, "threadLocals")
		inheritableThreadLocalsStmt := astutil.CreateAssignNilStmt(x, "inheritableThreadLocals")
		body.List = append(list[:index+1], append([]ast.Stmt{threadLocalsStmt, inheritableThreadLocalsStmt}, list[index+1:]...)...)
		if options.Debug || options.Verbose {
			log.Info("enhance function 'runtime.goexit0' add statement 'gp.threadLocals = nil'")
			log.Info("enhance function 'runtime.goexit0' add statement 'gp.inheritableThreadLocals = nil'")
		}
		return true
	default:
		return false
	}
}
