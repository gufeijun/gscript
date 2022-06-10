package token

// token kind
const (
	TOKEN_EOF        = iota // end of file
	TOKEN_IDENTIFIER        // identifier
	TOKEN_NUMBER            // number
	TOKEN_STRING            // string

	// seperator
	TOKEN_SEP_DOT    // .
	TOKEN_SEP_VARARG // ...
	TOKEN_SEP_COLON  // :
	TOKEN_SEP_SEMI   // ;
	TOKEN_SEP_COMMA  // ,
	TOKEN_SEP_QMARK  // ?
	TOKEN_SEP_LBRACK // [
	TOKEN_SEP_RBRACK // ]
	TOKEN_SEP_LPAREN // (
	TOKEN_SEP_RPAREN // )
	TOKEN_SEP_LCURLY // {
	TOKEN_SEP_RCURLY // }

	// operator
	TOKEN_BINOP_START
	TOKEN_OP_ADD  // +
	TOKEN_OP_SUB  // -
	TOKEN_OP_MUL  // *
	TOKEN_OP_DIV  // /
	TOKEN_OP_MOD  // %
	TOKEN_OP_AND  // &
	TOKEN_OP_XOR  // ^
	TOKEN_OP_OR   // |
	TOKEN_OP_IDIV // //
	TOKEN_OP_SHR  // >>
	TOKEN_OP_SHL  // <<
	TOKEN_OP_LE   // <=
	TOKEN_OP_GE   // >=
	TOKEN_OP_LT   // <
	TOKEN_OP_GT   // >
	TOKEN_OP_EQ   // ==
	TOKEN_OP_NE   // !=
	TOKEN_OP_LAND // &&
	TOKEN_OP_LOR  // ||
	TOKEN_BINOP_END

	TOKEN_OP_NOT  // ~
	TOKEN_OP_LNOT // !
	TOKEN_OP_INC  // ++
	TOKEN_OP_DEC  // --

	// don't change the order of constants between TOKEN_ASIGN_START and TOKEN_ASIGN_END
	TOKEN_ASIGN_START
	TOKEN_OP_ASSIGN // =
	TOKEN_OP_ADDEQ  // +=
	TOKEN_OP_SUBEQ  // -=
	TOKEN_OP_MULEQ  // *=
	TOKEN_OP_DIVEQ  // /=
	TOKEN_OP_MODEQ  // %=
	TOKEN_OP_ANDEQ  // &=
	TOKEN_OP_XOREQ  // ^=
	TOKEN_OP_OREQ   // |=
	TOKEN_ASIGN_END

	// keywords
	TOKEN_KW_BREAK       // break
	TOKEN_KW_CONTINUE    // continue
	TOKEN_KW_FOR         // for
	TOKEN_KW_WHILE       // while
	TOKEN_KW_IF          // if
	TOKEN_KW_ELIF        // elif
	TOKEN_KW_ELSE        // else
	TOKEN_KW_SWITCH      // switch
	TOKEN_KW_CASE        // case
	TOKEN_KW_FALLTHROUGH // fallthrough
	TOKEN_KW_DEFAULT     // default
	TOKEN_KW_RETURN      // return
	TOKEN_KW_FUNC        // func
	TOKEN_KW_LET         // let
	TOKEN_KW_TRUE        // true
	TOKEN_KW_FALSE       // false
	TOKEN_KW_NEW         // new
	TOKEN_KW_NIL         // nil
	TOKEN_KW_CLASS       // class
	TOKEN_KW_CONST       // const
	TOKEN_KW_ENUM        // enum
	TOKEN_KW_LOOP        // loop
	TOKEN_KW_IMPORT      // import
	TOKEN_KW_EXPORT      // export
	TOKEN_KW_AS          // as
	TOKEN_KW_GOTO        // goto
	TOKEN_KW_TRY         // try
	TOKEN_KW_CATCH       // catch
	TOKEN_KW_FINALLY     // finally
)

var _eofToken = Token{Kind: TOKEN_EOF}
var EOFToken = &_eofToken

type Token struct {
	Kind    int         // token kind
	Line    int         // line number
	Kth     int         // index of first charactor of Content in current line
	Content string      // original string
	Value   interface{} // token value, string literal, number literal or identifier
}

var Keywords = map[string]int{
	"break":       TOKEN_KW_BREAK,
	"continue":    TOKEN_KW_CONTINUE,
	"for":         TOKEN_KW_FOR,
	"while":       TOKEN_KW_WHILE,
	"if":          TOKEN_KW_IF,
	"elif":        TOKEN_KW_ELIF,
	"else":        TOKEN_KW_ELSE,
	"switch":      TOKEN_KW_SWITCH,
	"case":        TOKEN_KW_CASE,
	"fallthrough": TOKEN_KW_FALLTHROUGH,
	"default":     TOKEN_KW_DEFAULT,
	"return":      TOKEN_KW_RETURN,
	"func":        TOKEN_KW_FUNC,
	"let":         TOKEN_KW_LET,
	"true":        TOKEN_KW_TRUE,
	"false":       TOKEN_KW_FALSE,
	"new":         TOKEN_KW_NEW,
	"nil":         TOKEN_KW_NIL,
	"class":       TOKEN_KW_CLASS,
	"const":       TOKEN_KW_CONST,
	"enum":        TOKEN_KW_ENUM,
	"loop":        TOKEN_KW_LOOP,
	"import":      TOKEN_KW_IMPORT,
	"export":      TOKEN_KW_EXPORT,
	"as":          TOKEN_KW_AS,
	"goto":        TOKEN_KW_GOTO,
	"try":         TOKEN_KW_TRY,
	"catch":       TOKEN_KW_CATCH,
}
