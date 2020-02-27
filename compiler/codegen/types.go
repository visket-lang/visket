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
		errors.ErrorExit(fmt.Sprintf("%s | unknown type %s", t.Token.Pos, t.String()))
	}

	if t.IsArray {
		typ = types.NewArray(t.Len, typ)
	}

	return typ
}
