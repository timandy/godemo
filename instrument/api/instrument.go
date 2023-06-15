package api

import (
	"go/ast"
	"go/token"
)

type Instrument interface {
	PreHandlePackage(options *CompileOptions, result *InstrumentResult) bool

	PreHandleFile(path string, idx int, options *CompileOptions, result *InstrumentResult) bool

	HandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *CompileOptions, result *InstrumentResult) bool

	PostHandleFile(path string, idx int, fset *token.FileSet, af *ast.File, options *CompileOptions, result *InstrumentResult)

	PostHandlePackage(options *CompileOptions, result *InstrumentResult)
}
