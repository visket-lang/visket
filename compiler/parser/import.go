package parser

import (
	"github.com/arata-nvm/visket/compiler/lexer"
	"os"
	"path"
)

func (p *Parser) importFile(filename string) bool {
	filePath, ok := p.findFile(filename + ".sl")
	if !ok {
		return false
	}

	if _, ok := p.importedFiles[filePath]; ok {
		p.nextToken()
		return true
	}
	p.importedFiles[filePath] = true

	lex, err := lexer.NewFromFile(filePath)
	if err != nil {
		return false
	}

	p.l = append(p.l, lex)
	p.nextToken()
	return true
}

func (p *Parser) findFile(filename string) (string, bool) {
	dir, _ := path.Split(p.l[len(p.l)-1].Filename())
	filePath := path.Join(dir, filename)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, true
	}

	return "", false
}
