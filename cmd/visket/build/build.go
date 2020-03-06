package build

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler"
	"github.com/arata-nvm/visket/compiler/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
)

func EmitLLVM(filename, outputPath string, optimize bool) error {
	defer onPanicked()
	fmt.Printf("Compiling %s\n", filename)
	c := compiler.New()
	c.Compile(filename).ShowExit()
	if optimize {
		fmt.Println("Optimizing")
		c.Optimize()
	}
	fmt.Println("Building")
	compiled := c.GenIR()

	if outputPath == "" {
		outputPath = getFileNameWithoutExt(filename) + ".ll"
	}

	err := ioutil.WriteFile(outputPath, []byte(compiled), 0666)
	if err != nil {
		return err
	}

	fmt.Println("Finished")
	return nil
}

func Build(filename, outputPath string, optimize bool) error {
	defer onPanicked()
	fmt.Printf("Compiling %s\n", filename)
	c := compiler.New()
	c.Compile(filename).ShowExit()
	if optimize {
		fmt.Println("Optimizing")
		c.Optimize()
	}
	fmt.Println("Building")
	compiled := c.GenIR()

	tmpDir, err := ioutil.TempDir("", "visket")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tmpDir+"/main.ll", []byte(compiled), 0666)
	if err != nil {
		return err
	}

	if outputPath == "" {
		outputPath = getFileNameWithoutExt(filename)
	}

	clangArgs := []string{
		"-Wno-override-module",
		tmpDir + "/main.ll",
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

	os.RemoveAll(tmpDir)

	fmt.Println("Finished")
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
