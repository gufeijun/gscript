package ast

import . "gscript/complier/lexer"

// do not change the order of following const
const (
	ASIGN_OP_START  = TOKEN_ASIGN_START + iota
	ASIGN_OP_ASSIGN // =
	ASIGN_OP_ADDEQ  // +=
	ASIGN_OP_SUBEQ  // -=
	ASIGN_OP_MULEQ  // *=
	ASIGN_OP_DIVEQ  // /=
	ASIGN_OP_MODEQ  // %=
	ASIGN_OP_ANDEQ  // &=
	ASIGN_OP_XOREQ  // ^=
	ASIGN_OP_OREQ   // |=
)

const (
	BINOP_START = TOKEN_BINOP_START + iota
	BINOP_ADD   // +
	BINOP_SUB   // -
	BINOP_MUL   // *
	BINOP_DIV   // /
	BINOP_MOD   // %
	BINOP_AND   // &
	BINOP_XOR   // ^
	BINOP_OR    // |
	BINOP_IDIV  // //
	BINOP_SHR   // >>
	BINOP_SHL   // <<
	BINOP_LE    // <=
	BINOP_GE    // >=
	BINOP_LT    // <
	BINOP_GT    // >
	BINOP_EQ    // ==
	BINOP_NE    // !=
	BINOP_LAND  // &&
	BINOP_LOR   // ||
	BINOP_ATTR  // []
)

const (
	UNOP_NOT  = iota // ~
	UNOP_LNOT        // !
	UNOP_NEG         // -
	UNOP_DEC         // --i
	UNOP_INC         // ++i
	UNOP_DEC_        // i--
	UNOP_INC_        // i++
)
