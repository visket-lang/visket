package build

import (
	"fmt"
	"github.com/arata-nvm/visket/compiler/errors"
	"path/filepath"
	"runtime/debug"
)

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
