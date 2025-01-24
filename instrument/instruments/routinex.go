package instruments

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/timandy/routiner/instrument/api"
	"github.com/timandy/routiner/tools/log"
)

type RoutineXInstrument struct {
}

func NewRoutineXInstrument() api.Instrument {
	return &RoutineXInstrument{}
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInstrument) PreHandlePackage(options *api.CompileOptions, result *api.InstrumentResult) bool {
	return options.Package == "github.com/timandy/routine" || options.Package == "github.com/timandy/routine/g"
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInstrument) PreHandleFile(path string, idx int, options *api.CompileOptions, result *api.InstrumentResult) bool {
	return true
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInstrument) HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *api.CompileOptions, result *api.InstrumentResult) bool {
	for _, comment := range af.Comments {
		for _, c := range comment.List {
			if r.hasTag(c) {
				return true
			}
		}
	}
	return false
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInstrument) PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *api.CompileOptions, result *api.InstrumentResult) {
	srcDir := filepath.Dir(path)
	srcShortName := filepath.Base(path)
	srcExtName := filepath.Ext(srcShortName)
	destShortName := srcShortName[:len(srcShortName)-len(srcExtName)] + "_link.go"
	destPath := filepath.Join(srcDir, destShortName)
	result.ReplaceFiles[idx] = destPath
	if options.Debug || options.Verbose {
		log.Infof("replace source '%v' with '%v'", srcShortName, destShortName)
	}
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInstrument) PostHandlePackage(options *api.CompileOptions, result *api.InstrumentResult) {
}

// hasTag 是否有 !routinex 编译标记
func (r *RoutineXInstrument) hasTag(comment *ast.Comment) bool {
	return comment != nil && strings.TrimSpace(comment.Text) == "//go:build !routinex"
}
