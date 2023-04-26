package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"
	LABEL   = "LABEL"
	INT     = "INT"
	STRING  = "STRING"
	COMMA   = "COMMA"

	ADD = "ADD"
	AND = "AND"
	DEC = "DEC"
	DIV = "DIV"
	INC = "INC"
	MUL = "MUL"
	OR  = "OR"
	SUB = "SUB"
	XOR = "XOR"

	CALL  = "CALL"
	JMP   = "JMP"
	JMPNZ = "JMPNZ"
	JMPZ  = "JMPZ"
	RET   = "RET"

	PUSH = "PUSH"
	POP  = "POP"

	IS_STRING  = "IS_STRING"
	IS_INTEGER = "IS_INTEGER"
	STRING2INT = "STRING2INT"
	INT2STRING = "INT2STRING"

	CMP = "CMP"

	STORE = "STORE"

	PRINT_INT = "PRINT_INT"
	PRINT_STR = "PRINT_STR"

	PEEK = "PEEK"
	POKE = "POKE"

	CONCAT = "CONCAT"
	DATA   = "DATA"
	DB     = "DB"
	EXIT   = "EXIT"
	MEMCPY = "MEMCPY"
	NOP    = "NOP"
	RANDOM = "RANDOM"
	SYSTEM = "SYSTEM"
	TRAP   = "TRAP"
)

var keywords = map[string]Type{
	// compare
	"cmp": CMP,

	// types
	"is_integer": IS_INTEGER,
	"is_string":  IS_STRING,
	"int2string": INT2STRING,
	"string2int": STRING2INT,

	// store
	"store": STORE,

	// print
	"print_int": PRINT_INT,
	"print_str": PRINT_STR,

	// math
	"add": ADD,
	"and": AND,
	"dec": DEC,
	"div": DIV,
	"inc": INC,
	"mul": MUL,
	"or":  OR,
	"sub": SUB,
	"xor": XOR,

	// control-flow
	"call":  CALL,
	"jmp":   JMP,
	"jmpnz": JMPNZ,
	"jmpz":  JMPZ,
	"ret":   RET,

	// stack
	"push": PUSH,
	"pop":  POP,

	// memory
	"peek": PEEK,
	"poke": POKE,

	// misc
	"concat": CONCAT,
	"DATA":   DATA,
	"DB":     DB,
	"exit":   EXIT,
	"memcpy": MEMCPY,
	"nop":    NOP,
	"random": RANDOM,
	"system": SYSTEM,
	"int":    TRAP,
}

type Type string

type Token struct {
	Type    Type
	Literal string
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) Type {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
