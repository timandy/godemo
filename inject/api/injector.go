package api

import (
	"go/ast"
	"go/token"
)

type Injector interface {
	PreHandlePackage(options *CompileOptions, result *InjectResult) bool

	PreHandleFile(path string, idx int, options *CompileOptions, result *InjectResult) bool

	HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *CompileOptions, result *InjectResult) bool

	PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *CompileOptions, result *InjectResult)

	PostHandlePackage(options *CompileOptions, result *InjectResult)
}
