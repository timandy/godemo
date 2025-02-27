package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	args := []string{"routinex", "-h", "/demo", "-p", "ttt", "go", "version"}
	os.Args = args
	main()
}

func TestOtherCmd(t *testing.T) {
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", "git", "version"}
	os.Args = args
	main()
}

func TestOtherCmdHelp(t *testing.T) {
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", "git", "-h"}
	os.Args = args
	main()
}

func TestCompileCmdHelp(t *testing.T) {
	goToolDir := getGoToolDir()
	compilePath := path.Join(goToolDir, "compile.exe")
	args := []string{"routinex", "-p", "/demo", "-p", "ttt", compilePath, "-h"}
	os.Args = args
	// expect exit 2
	// main()
}

func getGoToolDir() string {
	cmd := exec.Command("go", "env")
	out := bytes.Buffer{}
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	//运行命令
	err := cmd.Run()
	if err != nil {
		fmt.Println("go env error:", err)
		return ""
	}
	//获取输出
	outStr := out.String()
	rows := strings.Split(outStr, "\n")
	for _, row := range rows {
		skvArray := strings.Split(row, "=")
		if len(skvArray) != 2 {
			continue
		}
		skArray := skvArray[0]
		sk := strings.Split(skArray, " ")
		if len(sk) != 2 {
			continue
		}
		k := sk[1]
		if !strings.EqualFold(k, "GOTOOLDIR") {
			continue
		}
		return skvArray[1]
	}
	return ""
}
