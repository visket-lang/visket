package build

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler"
	"github.com/arata-nvm/visket/compiler/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
)

func EmitLLVM(filename, outputPath string, optimize bool) error {
	defer onPanicked()

	c := compiler.New()
	c.Compile(filename).ShowExit()
	if optimize {
		c.Optimize()
	}
	compiled := c.GenIR()

	if outputPath == "" {
		outputPath = getFileNameWithoutExt(filename) + ".ll"
	}

	err := ioutil.WriteFile(outputPath, []byte(compiled), 0666)
	if err != nil {
		return err
	}

	return nil
}

func Build(filename, outputPath string, optimize bool) error {
	defer onPanicked()

	tmpDir, err := ioutil.TempDir("", "visket")
	if err != nil {
		return err
	}

	llFilePath := path.Join(tmpDir, "/main.ll")

	err = EmitLLVM(filename, llFilePath, optimize)
	if err != nil {
		return err
	}

	if outputPath == "" {
		outputPath = getFileNameWithoutExt(filename)
	}

	clangArgs := []string{
		"-Wno-override-module",
		llFilePath,
		"-o", outputPath,
	}

	if optimize {
		clangArgs = append(clangArgs, "-O3")
	}

	cmd := exec.Command("clang", clangArgs...)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}

	return nil
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func onPanicked() {
	if err := recover(); err != nil {
		errors.Error("failed compiling")
		errors.Error(fmt.Sprintf("%+v", err))
		errors.ErrorExit(string(debug.Stack()))
	}
}
