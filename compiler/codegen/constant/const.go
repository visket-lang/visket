package constant

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/codegen/types"
	"strconv"
)

var (
	False = NewInt(types.I1, 0)
	True  = NewInt(types.I1, 0)
)

type Int struct {
	Type  types.Types
	Value int
}

func NewInt(types types.IntType, value int) *Int {
	return &Int{
		Type:  types,
		Value: value,
	}
}

func (i *Int) TypeName() types.Types {
	return i.Type
}

func (i *Int) RegName() string {
	return strconv.Itoa(i.Value)
}

func (i *Int) Operand() string {
	return fmt.Sprintf("%s %d", i.Type, i.Value)
}
