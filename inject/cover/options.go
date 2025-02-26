package cover

import (
	"path/filepath"
	"strings"

	"github.com/timandy/routiner/tools/json"
	"github.com/timandy/routiner/tools/os"
	"github.com/timandy/routiner/tools/slices"
)

type CoverOptions struct {
	PkgCfg  string   `name:"" shorthand:"pkgcfg" usage:"enable full-package instrumentation mode using params from specified config file"`
	Package string   // package
	Debug   bool     // debug mode enabled or not
	Verbose bool     // verbose mode enabled or not
	Args    []string // remain args exclude the options of current program
}

func (c *CoverOptions) IsDebug() bool {
	return c.Debug
}

func (c *CoverOptions) IsVerbose() bool {
	return c.Verbose
}

func (c *CoverOptions) GetArgs() []string {
	return c.Args
}

func (c *CoverOptions) SetArgs(args []string) {
	c.Args = args
}

func (c *CoverOptions) GetPackage() string {
	return c.Package
}

func (c *CoverOptions) GetWorkDir() string {
	return ""
}

func (c *CoverOptions) ReadConfig() {
	if !os.IsFile(c.PkgCfg) {
		return
	}
	bytes := os.ReadFile(c.PkgCfg)
	m := json.Unmarshal[map[string]any](bytes)
	if m == nil {
		return
	}
	pkg, ok := m["PkgPath"]
	if !ok {
		return
	}
	c.Package = pkg.(string)
}

func (c *CoverOptions) IsValid(execName string) bool {
	cmd := filepath.Base(execName)
	if ext := filepath.Ext(cmd); ext != "" {
		cmd = strings.TrimSuffix(cmd, ext)
	}
	return cmd == "cover" && c.PkgCfg != ""
}

func (c *CoverOptions) Clone() *CoverOptions {
	return &CoverOptions{
		PkgCfg:  c.PkgCfg,
		Package: c.Package,
		Debug:   c.Debug,
		Verbose: c.Verbose,
		Args:    slices.Clone(c.Args),
	}
}
