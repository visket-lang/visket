package codegen

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/ast"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/llir/llvm/ir/types"
)

func (c *CodeGen) llvmType(t *ast.Type) types.Type {
	typ, ok := c.context.findType(t.Token.Literal)
	if !ok {
		errors.ErrorExit(fmt.Sprintf("%s | unknown type %s", t.Token.Pos, t.Token.Literal))
	}

	if t.IsArray {
		typ = types.NewArray(t.Len, typ)
	}

	return typ
}

type Struct struct {
	Name    string
	Members []*Member
	Type    *types.StructType
}

type Member struct {
	Name string
	Id   int
	Type types.Type
}

func (s *Struct) findMember(name string) int {
	for _, m := range s.Members {
		if m.Name == name {
			return m.Id
		}
	}

	return -1
}
