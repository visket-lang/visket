package build

import (
	"os"
	"os/exec"
)

func buildLlFile(path, outputPath string, doOptimize bool) error {
	clangArgs := []string{
		"-Wno-override-module",
		"-lm",
		path,
		"-o", outputPath,
	}

	if doOptimize {
		clangArgs = append(clangArgs, "-O3")
		clangArgs = append(clangArgs, "-flto")
	} else {
		clangArgs = append(clangArgs, "-O0")
	}

	cmd := exec.Command("clang", clangArgs...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func buildIncludedFile(includedFile, outputPath string, doOptimize bool) error {
	args := []string{
		"-S",
		"-emit-llvm",
		"-o", outputPath,
		includedFile,
	}

	if doOptimize {
		args = append(args, "-O3")
	} else {
		args = append(args, "-O0")
	}

	cmd := exec.Command("clang", args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func linkLlFiles(paths []string, outputPath string) error {
	args := []string{
		"-S",
		"-o", outputPath,
	}

	args = append(args, paths...)

	cmd := exec.Command("llvm-link", args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
