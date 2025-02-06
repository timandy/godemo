package injectors

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/timandy/routiner/inject/api"
	"github.com/timandy/routiner/tools/log"
)

type RoutineXInjector struct {
}

func NewRoutineXInjector() api.Injector {
	return &RoutineXInjector{}
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PreHandlePackage(options *api.CompileOptions, result *api.InjectResult) bool {
	return options.Package == "github.com/timandy/routine" || options.Package == "github.com/timandy/routine/g"
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) PreHandleFile(path string, idx int, options *api.CompileOptions, result *api.InjectResult) bool {
	return true
}

//goland:noinspection GoUnusedParameter
func (r *RoutineXInjector) HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *api.CompileOptions, result *api.InjectResult) bool {
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
func (r *RoutineXInjector) PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *api.CompileOptions, result *api.InjectResult) {
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
func (r *RoutineXInjector) PostHandlePackage(options *api.CompileOptions, result *api.InjectResult) {
}

// hasTag 是否有 !routinex 编译标记
func (r *RoutineXInjector) hasTag(comment *ast.Comment) bool {
	return comment != nil && strings.TrimSpace(comment.Text) == "//go:build !routinex"
}
