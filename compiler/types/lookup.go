package types

var (
	VOID  = NewSlVoid()
	BOOL  = NewSlBool()
	INT   = NewSlInt()
	FLOAT = NewSlFloat()
)

var nameToType = map[string]SlType{
	"void":  VOID,
	"bool":  BOOL,
	"int":   INT,
	"float": FLOAT,
}

func LookupType(name string) (SlType, bool) {
	typ, ok := nameToType[name]
	return typ, ok
}
