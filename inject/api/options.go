package api

import (
	"path/filepath"
	"strings"

	"github.com/timandy/routiner/tools/slices"
)

type CompileOptions struct {
	Package string   `name:"package" shorthand:"p" usage:"set expected package import path"`
	Output  string   `name:"output" shorthand:"o" usage:"write output to file"`
	Debug   bool     // debug mode enabled or not
	Verbose bool     // verbose mode enabled or not
	Args    []string // remain args exclude the options of current program
}

func (c CompileOptions) IsValid(execName string) bool {
	cmd := filepath.Base(execName)
	if ext := filepath.Ext(cmd); ext != "" {
		cmd = strings.TrimSuffix(cmd, ext)
	}
	return cmd == "compile" && c.Package != "" && c.Output != ""
}

func (c CompileOptions) WorkDir() string {
	return filepath.Dir(c.Output)
}

func (c CompileOptions) Clone() *CompileOptions {
	return &CompileOptions{
		Package: c.Package,
		Output:  c.Output,
		Debug:   c.Debug,
		Verbose: c.Verbose,
		Args:    slices.Clone(c.Args),
	}
}
