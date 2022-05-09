package proto

type Instruction byte

const INS_BINARY_START = INS_BINARY_ADD - 1
const INS_ATTR_ASSIGN_START = INS_ATTR_ASSIGN - 1

const (
	INS_UNARY_NOT   byte = iota // ~
	INS_UNARY_NEG               // -
	INS_UNARY_LNOT              // !
	INS_BINARY_ADD              // +
	INS_BINARY_SUB              // -
	INS_BINARY_MUL              // *
	INS_BINARY_DIV              // /
	INS_BINARY_MOD              // %
	INS_BINARY_AND              // &
	INS_BINARY_XOR              // ^
	INS_BINARY_OR               // |
	INS_BINARY_IDIV             // //
	INS_BINARY_SHR              // >>
	INS_BINARY_SHL              // <<
	INS_BINARY_LE               // <=
	INS_BINARY_GE               // >=
	INS_BINARY_LT               // <
	INS_BINARY_GT               // >
	INS_BINARY_EQ               // ==
	INS_BINARY_NE               // !=
	INS_BINARY_LAND             // &&
	INS_BINARY_LOR              // ||
	INS_BINARY_ATTR             // []
	INS_LOAD_CONST
	INS_LOAD_NAME
	INS_STORE_NAME
	INS_PUSH_NAME_NIL
	INS_PUSH_NAME
	INS_RESIZE_NAMETABLE
	INS_POP_TOP
	INS_STOP
	INS_SLICE_NEW
	INS_SLICE_APPEND
	INS_MAP_NEW
	INS_ATTR_ASSIGN       // =
	INS_ATTR_ASSIGN_ADDEQ // +=
	INS_ATTR_ASSIGN_SUBEQ // -=
	INS_ATTR_ASSIGN_MULEQ // *=
	INS_ATTR_ASSIGN_DIVEQ // /=
	INS_ATTR_ASSIGN_MODEQ // %=
	INS_ATTR_ASSIGN_ANDEQ // &=
	INS_ATTR_ASSIGN_XOREQ // ^=
	INS_ATTR_ASSIGN_OREQ  // |=
	INS_ATTR_ACCESS
	INS_JUMP_REL
	INS_JUMP_ABS
	INS_JUMP_IF
	INS_JUMP_CASE
	INS_ROT_TWO
)