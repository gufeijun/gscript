package lexer

import (
	"fmt"
	"gscript/complier/token"
	"os"
	"strconv"
	"strings"
)

const (
	CHAR_EOF   byte = 0
	CHAR_CR    byte = '\r'
	CHAR_LF    byte = '\n'
	CHAR_SPACE byte = ' '
	CHAR_TAB   byte = '\t'
)

type Lexer struct {
	line       int          // current line number
	column     int          // column of current charactor in current line
	cursor     int          // number of next character needs to be parse
	src        []byte       // source code
	srcFile    string       // source file path
	curToken   *token.Token // current token
	aheadToken *token.Token // save LookAhead token temporarily
}

func NewLexer(srcFile string, src []byte) *Lexer {
	return &Lexer{
		line:    1,
		src:     src,
		cursor:  0,
		srcFile: srcFile,
	}
}

func (l *Lexer) Line() int {
	return l.line
}

func (l *Lexer) Column() int {
	return l.column
}

func (l *Lexer) SrcFile() string {
	return l.srcFile
}

// Look ahead 1 token
func (l *Lexer) LookAhead() (token *token.Token) {
	if l.aheadToken != nil {
		return l.aheadToken
	}
	l.aheadToken = l.nextToken()
	return l.aheadToken
}

func (l *Lexer) NextToken() (token *token.Token) {
	// if LookAhead before
	if l.aheadToken != nil {
		token, l.aheadToken = l.aheadToken, nil
		return
	}
	return l.nextToken()
}

func (l *Lexer) nextToken() *token.Token {
again:
	if breakLine := l.skipWhiteSpace(); breakLine && l.needAddSemi() {
		l.genSemiToken()
		return l.curToken
	}
	if l.reachEndOfFile() {
		return token.EOFToken
	}
	curCh := l.src[l.cursor]
	switch curCh {
	case '.':
		nextCh, gapCh := l.lookAhead(1), l.lookAhead(2)
		if nextCh == '.' && gapCh == '.' {
			l.genToken(token.TOKEN_SEP_VARARG, 3)
			l.forward(2)
			break
		}
		l.genToken(token.TOKEN_SEP_DOT, 1)
	case '"', '\'':
		l.scanStringLiteral()
	case ':':
		l.genToken(token.TOKEN_SEP_COLON, 1)
	case ';':
		l.genToken(token.TOKEN_SEP_SEMI, 1)
	case ',':
		l.genToken(token.TOKEN_SEP_COMMA, 1)
	case '?':
		l.genToken(token.TOKEN_SEP_QMARK, 1)
	case '[':
		l.genToken(token.TOKEN_SEP_LBRACK, 1)
	case ']':
		l.genToken(token.TOKEN_SEP_RBRACK, 1)
	case '(':
		l.genToken(token.TOKEN_SEP_LPAREN, 1)
	case ')':
		l.genToken(token.TOKEN_SEP_RPAREN, 1)
	case '{':
		l.genToken(token.TOKEN_SEP_LCURLY, 1)
	case '}':
		l.genToken(token.TOKEN_SEP_RCURLY, 1)
	case '+':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_ADDEQ, 2)
			l.forward(1)
		} else if nextCh == '+' {
			l.genToken(token.TOKEN_OP_INC, 2)
			l.forward(1)
		} else {
			l.genToken(token.TOKEN_OP_ADD, 1)
		}
	case '-':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_SUBEQ, 2)
			l.forward(1)
		} else if nextCh == '-' {
			l.genToken(token.TOKEN_OP_DEC, 2)
			l.forward(1)
		} else {
			l.genToken(token.TOKEN_OP_SUB, 1)
		}
	case '*':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_MULEQ, 2)
			l.forward(1)
			break
		}
		l.genToken(token.TOKEN_OP_MUL, 1)
	case '/':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_DIVEQ, 2)
			l.forward(1)
		} else if nextCh == '/' {
			l.genToken(token.TOKEN_OP_IDIV, 2)
			l.forward(1)
		} else {
			l.genToken(token.TOKEN_OP_DIV, 1)
		}
	case '%':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_MODEQ, 2)
			l.forward(1)
			break
		}
		l.genToken(token.TOKEN_OP_MOD, 1)
	case '&':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_ANDEQ, 2)
			l.forward(1)
		} else if nextCh == '&' {
			l.genToken(token.TOKEN_OP_LAND, 2)
			l.forward(1)
		} else {
			l.genToken(token.TOKEN_OP_AND, 1)
		}
	case '|':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_OREQ, 2)
			l.forward(1)
		} else if nextCh == '|' {
			l.genToken(token.TOKEN_OP_LOR, 2)
			l.forward(1)
		} else {
			l.genToken(token.TOKEN_OP_OR, 1)
		}
	case '^':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_XOREQ, 2)
			l.forward(1)
			break
		}
		l.genToken(token.TOKEN_OP_XOR, 1)
	case '~':
		l.genToken(token.TOKEN_OP_NOT, 1)
	case '=':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_EQ, 2)
			l.forward(1)
			break
		}
		l.genToken(token.TOKEN_OP_ASSIGN, 1)
	case '!':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_NE, 2)
			l.forward(1)
			break
		}
		l.genToken(token.TOKEN_OP_LNOT, 1)
	case '<':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_LE, 2)
			l.forward(1)
		} else if nextCh == '<' {
			l.genToken(token.TOKEN_OP_SHL, 2)
			l.forward(1)
		} else {
			l.genToken(token.TOKEN_OP_LT, 1)
		}
	case '>':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(token.TOKEN_OP_GE, 2)
			l.forward(1)
		} else if nextCh == '>' {
			l.genToken(token.TOKEN_OP_SHR, 2)
			l.forward(1)
		} else {
			l.genToken(token.TOKEN_OP_GT, 1)
		}
	case '#': // comment
		l.skipComment()
		goto again
	default:
		if isDigit(curCh) {
			l.scanNumber()
		} else if isLetter_(curCh) {
			l.scanIdentifier()
		} else {
			l.error("unexpected symbol near '%c'", curCh)
		}
	}

	l.forward(1)
	return l.curToken
}

func (l *Lexer) scanIdentifier() {
	k := l.cursor + 1
	for ; k < len(l.src) && (isLetter_(l.src[k]) || isDigit(l.src[k])); k++ {
	}

	l.genToken(token.TOKEN_IDENTIFIER, k-l.cursor)
	id := string(l.src[l.cursor:k])
	l.curToken.Value = id
	if kind, ok := token.Keywords[id]; ok {
		l.curToken.Kind = kind
	}
	l.forward(k - l.cursor - 1)
}

// dec: 1234, hex: 0xFF, oct: 017, float: 0.123
func (l *Lexer) scanNumber() {
	var num, base int64
	var value interface{}

	k := l.cursor
	firstDigit := l.src[k]
	if firstDigit != '0' { // dec
		base = 10
	} else {
		nextCh := l.lookAhead(1)
		if nextCh == '.' { // float
			k, value = l.scanFloat(k + 2)
			goto end
		}
		if isDigit(nextCh) { // oct
			base = 8
			k += 1
		} else if nextCh == 'x' { // hex
			if gapCh := l.lookAhead(2); gapCh == CHAR_EOF || !isHexDigit(gapCh) {
				l.error("invalid hex number near '%c'", firstDigit)
			}
			base = 16
			k += 2
		} else { // 0
			base = 10
		}
	}

	for ; k < len(l.src) &&
		(l.src[k] == '.' ||
			base == 8 && l.src[k] < '8' && l.src[k] >= '0' ||
			base == 10 && isDigit(l.src[k]) ||
			base == 16 && isHexDigit(l.src[k])); k++ {
		if l.src[k] == '.' {
			k, value = l.scanFloat(k + 1)
			goto end
		}
		num = num*base + toNumber(l.src[k])
	}
	value = num
end:
	contentLength := k - l.cursor
	l.genToken(token.TOKEN_NUMBER, contentLength)
	l.curToken.Value = value
	l.forward(contentLength - 1)
}

func (l *Lexer) scanFloat(start int) (end int, value float64) {
	for end = start; end < len(l.src) && isDigit(l.src[end]); end++ {
	}
	value, err := strconv.ParseFloat(string(l.src[l.cursor:end]), 64)
	if err != nil {
		l.error(err.Error())
	}
	return end, value
}

func (l *Lexer) scanStringLiteral() {
	curCh := l.src[l.cursor]
	k := 1
	var b strings.Builder
	escape := false
	start := l.cursor + 1 // skip first " or '
	for {
		ahead := l.lookAhead(k)
		if ahead == CHAR_EOF || ahead == CHAR_CR || ahead == CHAR_LF {
			l.error("expect another quotation mark before end of file or newline")
		}
		if isQuoteAndUnmatched(ahead, curCh) {
			l.error("expect another %c, but got %c", curCh, ahead)
		}
		if ahead == l.src[l.cursor] { //matched
			break
		}
		if ahead != '\\' {
			k++
			continue
		}
		gap := l.lookAhead(k + 1)
		if gap == CHAR_EOF {
			l.error("expect another quotation mark before end of file or newline")
		}
		// only support \n, \t, \', \", \\,
		if gap != 'n' && gap != 't' && gap != '\'' && gap != '"' && gap != '\\' {
			l.cursor += k
			l.error("invalid escape character \\%c", gap)
		}
		escape = true
		b.Write(l.src[start : l.cursor+k])
		if gap == 'n' {
			b.WriteByte('\n')
		} else if gap == 't' {
			b.WriteByte('\t')
		} else {
			b.WriteByte(gap)
		}
		k += 2
		start = l.cursor + k
	}
	l.genToken(token.TOKEN_STRING, k+1)
	if escape {
		b.Write(l.src[start : l.cursor+k])
		l.curToken.Value = b.String()
	} else {
		l.curToken.Value = string(l.src[l.cursor+1 : l.cursor+k])
	}
	l.forward(k)
}

func (l *Lexer) skipComment() {
	for k := 1; ; k++ {
		ahead := l.lookAhead(k)
		if ahead == CHAR_EOF || ahead == CHAR_CR || ahead == CHAR_LF {
			l.forward(k)
			break
		}
	}
}

// skip \r, \n, \r\n, \t, space
func (l *Lexer) skipWhiteSpace() (breakLine bool) {
	for l.cursor < len(l.src) {
		switch l.src[l.cursor] {
		case CHAR_CR:
			breakLine = true
			nextCh := l.lookAhead(1)
			if nextCh == CHAR_EOF {
				break
			}
			if nextCh != CHAR_LF {
				l.column = 0
				l.line++
				break
			}
			l.cursor++
			fallthrough
		case CHAR_LF:
			breakLine = true
			l.column = 0
			l.line++
		case CHAR_SPACE:
			l.column++
		case CHAR_TAB:
			l.column++
		default:
			return
		}
		l.cursor++
	}
	// add breakLine at last of code so that we can tell
	// if need add semicolon for the last statement
	if l.cursor == len(l.src) {
		return true
	}
	return
}

// lookAhead @k characters
func (l *Lexer) lookAhead(k int) byte {
	idx := l.cursor + k
	if idx >= len(l.src) {
		return CHAR_EOF
	}
	return l.src[idx]
}

func (l *Lexer) error(format string, args ...interface{}) {
	errMsg := fmt.Sprintf(format, args...)
	fmt.Printf("Lexer error: [%s:%d:%d]\n", l.srcFile, l.line, l.column+1)
	fmt.Printf("\t%s\n", errMsg)
	os.Exit(0)
}

// move cursor and kth @k steps
func (l *Lexer) forward(k int) {
	l.cursor += k
	l.column += k
}

func (l *Lexer) reachEndOfFile() bool {
	return l.cursor >= len(l.src)
}

func (l *Lexer) genToken(kind, contentLength int) {
	l.curToken = &token.Token{
		Kind:    kind,
		Line:    l.line,
		Kth:     l.column,
		Content: string(l.src[l.cursor : l.cursor+contentLength]),
	}
}

func isQuoteAndUnmatched(a, b byte) bool {
	return (a == '"' || a == '\'') && a != b
}

func isDigit(ch byte) bool {
	return ch <= '9' && ch >= '0'
}

func isHexDigit(ch byte) bool {
	return ch <= '9' && ch >= '0' || ch <= 'f' && ch >= 'a' || ch <= 'F' && ch >= 'A'
}

func isLetter_(ch byte) bool {
	return ch == '_' || ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func toNumber(digit byte) (result int64) {
	if !isDigit(digit) {
		result += 9
	}
	return result + int64(digit&15)
}

var addSemiTokens = map[int]struct{}{
	token.TOKEN_IDENTIFIER:     {},
	token.TOKEN_NUMBER:         {},
	token.TOKEN_STRING:         {},
	token.TOKEN_KW_BREAK:       {},
	token.TOKEN_KW_FALLTHROUGH: {},
	token.TOKEN_KW_CONTINUE:    {},
	token.TOKEN_KW_RETURN:      {},
	token.TOKEN_OP_INC:         {},
	token.TOKEN_OP_DEC:         {},
	token.TOKEN_SEP_RBRACK:     {},
	token.TOKEN_SEP_RCURLY:     {},
	token.TOKEN_SEP_RPAREN:     {},
	token.TOKEN_KW_NIL:         {},
	token.TOKEN_KW_TRUE:        {},
	token.TOKEN_KW_FALSE:       {},
}

func needAddSemi(kind int) bool {
	_, ok := addSemiTokens[kind]
	return ok
}

func (l *Lexer) needAddSemi() bool {
	return l.curToken != nil && needAddSemi(l.curToken.Kind)
}

func (l *Lexer) genSemiToken() {
	l.curToken = &token.Token{
		Kind:    token.TOKEN_SEP_SEMI,
		Line:    l.curToken.Line,
		Kth:     l.curToken.Kth + len(l.curToken.Content),
		Content: ";",
	}
}
