package build

import (
	"github.com/arata-nvm/visket/compiler"
	"io/ioutil"
	"os"
	"path"
)

func EmitLLVM(filePath, outputPath string, doOptimize bool) error {
	llPath, err := GenLl(filePath, doOptimize)
	if err != nil {
		return err
	}

	// 出力先のデフォルト値を設定する
	if outputPath == "" {
		outputPath = getFileNameWithoutExt(filePath) + ".ll"
	}

	// 出力先にファイルをコピーする
	data, err := ioutil.ReadFile(llPath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outputPath, data, 0666)
	return err
}

func Build(filePath, outputPath string, doOptimize bool) error {
	llPath, err := GenLl(filePath, doOptimize)
	if err != nil {
		return err
	}

	// 出力先のデフォルト値を設定する
	if outputPath == "" {
		outputPath = getFileNameWithoutExt(filePath)
	}

	// 実行可能ファイルにコンパイルする
	err = buildLlFile(llPath, outputPath, doOptimize)
	return err
}

func GenLl(filePath string, doOptimize bool) (string, error) {
	// .slファイルをコンパイルする
	c := compiler.New()
	c.Compile(filePath).ShowExit(false)
	if doOptimize {
		c.Optimize()
	}

	// 一時ディレクトリを作成する
	tmpDir, err := ioutil.TempDir("", "visket")
	if err != nil {
		return "", err
	}

	// コンパイル結果をファイルに出力する
	llPath := path.Join(tmpDir, getFileNameWithoutExt(filePath)+".ll")
	err = ioutil.WriteFile(llPath, []byte(c.GenIR()), 0666)
	if err != nil {
		return "", err
	}

	// インクルードしたファイルをコンパイル、ファイルに出力する
	var llPaths []string
	llPaths = append(llPaths, llPath)
	for _, includedFile := range c.IncludeFiles() {
		// ファイルの存在確認
		if _, err = os.Stat(includedFile); err != nil {
			return "", err
		}

		llPath = path.Join(tmpDir, getFileNameWithoutExt(includedFile)+".c.ll")
		err = buildIncludedFile(includedFile, llPath, doOptimize)
		if err != nil {
			return "", err
		}

		llPaths = append(llPaths, llPath)
	}

	// コンパイルされたファイルをリンクする
	llPath = path.Join(tmpDir, "main.ll")
	err = linkLlFiles(llPaths, llPath)
	if err != nil {
		return "", err
	}

	return llPath, nil
}
