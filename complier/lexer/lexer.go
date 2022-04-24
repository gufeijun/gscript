package lexer

import (
	"fmt"
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
	line       int    // current line number
	kth        int    // index of current charactor in current line
	src        []byte // source code
	cursor     int    // number of next character needs to be parse
	srcFile    string // source file path
	curToken   *Token
	aheadToken *Token // save LookAhead token temporarily
}

func NewLexer(srcFile string, src []byte) *Lexer {
	return &Lexer{
		line:    1,
		src:     src,
		cursor:  0,
		srcFile: srcFile,
	}
}

// Look ahead 1 token
func (l *Lexer) LookAhead() (token *Token) {
	l.aheadToken = l.nextToken()
	return l.aheadToken
}

func (l *Lexer) NextToken() (token *Token) {
	// if LookAhead before
	if l.aheadToken != nil {
		token, l.aheadToken = l.aheadToken, nil
		return
	}
	return l.nextToken()
}
func (l *Lexer) nextToken() (token *Token) {
again:
	l.skipWhiteSpace()
	if l.reachEndOfFile() {
		return eofToken
	}
	curCh := l.src[l.cursor]
	switch curCh {
	case '.':
		nextCh, gapCh := l.lookAhead(1), l.lookAhead(2)
		if nextCh == '.' && gapCh == '.' {
			l.genToken(TOKEN_SEP_VARARG, 3)
			l.forward(2)
			break
		}
		l.genToken(TOKEN_SEP_DOT, 1)
	case '"', '\'':
		l.scanStringLiteral()
	case ':':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_CLONE, 2)
			l.forward(1)
			break
		}
		l.genToken(TOKEN_SEP_COLON, 1)
	case ';':
		l.genToken(TOKEN_SEP_SEMI, 1)
	case ',':
		l.genToken(TOKEN_SEP_COMMA, 1)
	case '?':
		l.genToken(TOKEN_SEP_QMARK, 1)
	case '[':
		l.genToken(TOKEN_SEP_LBRACK, 1)
	case ']':
		l.genToken(TOKEN_SEP_RBRACK, 1)
	case '(':
		l.genToken(TOKEN_SEP_LPAREN, 1)
	case ')':
		l.genToken(TOKEN_SEP_RPAREN, 1)
	case '{':
		l.genToken(TOKEN_SEP_LCURLY, 1)
	case '}':
		l.genToken(TOKEN_SEP_RCURLY, 1)
	case '+':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_ADDEQ, 2)
			l.forward(1)
		} else if nextCh == '+' {
			l.genToken(TOKEN_OP_INC, 2)
			l.forward(1)
		} else if isDigit(nextCh) {
			l.scanNumber(1)
		} else {
			l.genToken(TOKEN_OP_ADD, 1)
		}
	case '-':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_SUBEQ, 2)
			l.forward(1)
		} else if nextCh == '-' {
			l.genToken(TOKEN_OP_DEC, 2)
			l.forward(1)
		} else if isDigit(nextCh) {
			l.scanNumber(-1)
		} else {
			l.genToken(TOKEN_OP_SUB, 1)
		}
	case '*':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_MULEQ, 2)
			l.forward(1)
			break
		}
		l.genToken(TOKEN_OP_MUL, 1)
	case '/':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_DIVEQ, 2)
			l.forward(1)
		} else if nextCh == '/' {
			l.genToken(TOKEN_OP_IDIV, 2)
			l.forward(1)
		} else {
			l.genToken(TOKEN_OP_DIV, 1)
		}
	case '%':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_MODEQ, 2)
			l.forward(1)
			break
		}
		l.genToken(TOKEN_OP_MOD, 1)
	case '&':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_ANDEQ, 2)
			l.forward(1)
		} else if nextCh == '&' {
			l.genToken(TOKEN_OP_LAND, 2)
			l.forward(1)
		} else {
			l.genToken(TOKEN_OP_AND, 1)
		}
	case '|':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_OREQ, 2)
			l.forward(1)
		} else if nextCh == '|' {
			l.genToken(TOKEN_OP_LOR, 2)
			l.forward(1)
		} else {
			l.genToken(TOKEN_OP_OR, 1)
		}
	case '^':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_XOREQ, 2)
			l.forward(1)
			break
		}
		l.genToken(TOKEN_OP_XOR, 1)
	case '~':
		l.genToken(TOKEN_OP_NOT, 1)
	case '=':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_EQ, 2)
			l.forward(1)
			break
		}
		l.genToken(TOKEN_OP_ASSIGN, 1)
	case '!':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_NE, 2)
			l.forward(1)
			break
		}
		l.genToken(TOKEN_OP_LNOT, 1)
	case '<':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_LE, 2)
			l.forward(1)
		} else if nextCh == '<' {
			l.genToken(TOKEN_OP_SHL, 2)
			l.forward(1)
		} else {
			l.genToken(TOKEN_OP_LT, 1)
		}
	case '>':
		nextCh := l.lookAhead(1)
		if nextCh == '=' {
			l.genToken(TOKEN_OP_GE, 2)
			l.forward(1)
		} else if nextCh == '>' {
			l.genToken(TOKEN_OP_SHR, 2)
			l.forward(1)
		} else {
			l.genToken(TOKEN_OP_GT, 1)
		}
	case '#': // comment
		l.skipComment()
		goto again
	default:
		if isDigit(curCh) {
			l.scanNumber(0)
		} else if isLetter_(curCh) {
			l.scanIdentifier()
		} else {
			l.error("unexpected symbol near %c", curCh)
		}
	}

	l.forward(1)
	return l.curToken
}

func (l *Lexer) scanIdentifier() {
	k := l.cursor + 1
	for ; k < len(l.src) && (isLetter_(l.src[k]) || isDigit(l.src[k])); k++ {
	}

	l.genToken(TOKEN_IDENTIFIER, k-l.cursor)
	id := string(l.src[l.cursor:k])
	l.curToken.Value = id
	if kind, ok := keywords[id]; ok {
		l.curToken.Kind = kind
	}
	l.forward(k - l.cursor - 1)
}

// dec: 1234, hex: 0xFF, oct: 017, float: 0.123
func (l *Lexer) scanNumber(signed int64) {
	var num, base int64
	var skip int
	var value interface{}

	if signed != 0 { // curCh == '+' or '-'
		skip = 1
	} else {
		signed = 1
	}

	k := l.cursor + skip
	firstDigit := l.src[k]
	if firstDigit != '0' { // dec
		base = 10
	} else {
		nextCh := l.lookAhead(skip + 1)
		if nextCh == '.' { // float
			k, value = l.scanFloat(k + 2)
			goto end
		}
		if isDigit(nextCh) { // oct
			base = 8
			k += 1
		} else if nextCh == 'x' { // hex
			if gapCh := l.lookAhead(skip + 2); gapCh == CHAR_EOF || !isHexDigit(gapCh) {
				l.error("invalid hex number near %c", firstDigit)
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
	value = num * signed
end:
	contentLength := k - l.cursor
	l.genToken(TOKEN_NUMBER, contentLength)
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
		if ahead == CHAR_EOF || ahead == CHAR_CR || ahead == CHAR_LF || isQuoteAndUnmatched(ahead, curCh) {
			l.error("mismatch %c", curCh)
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
			l.error("mismatch %c", curCh)
		}
		// only support \n, \t, \', \", \\,
		if gap != 'n' && gap != 't' && gap != '\'' && gap != '"' && gap != '\\' {
			l.cursor += k
			l.error("unknow escape \\%c", gap)
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
	l.genToken(TOKEN_STRING, k+1)
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
func (l *Lexer) skipWhiteSpace() {
	for l.cursor < len(l.src) {
		switch l.src[l.cursor] {
		case CHAR_CR:
			nextCh := l.lookAhead(1)
			if nextCh == CHAR_EOF {
				break
			}
			if nextCh != CHAR_LF {
				l.kth = 0
				l.line++
				break
			}
			l.cursor++
			fallthrough
		case CHAR_LF:
			l.kth = 0
			l.line++
		case CHAR_SPACE:
			l.kth++
		case CHAR_TAB:
			l.kth++
		default:
			return
		}
		l.cursor++
	}
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
	fmt.Printf("[%s:%d:%d] %s\n", l.srcFile, l.line, l.kth, errMsg)
	os.Exit(-1)
}

// move cursor and kth @k steps
func (l *Lexer) forward(k int) {
	l.cursor += k
	l.kth += k
}

func (l *Lexer) reachEndOfFile() bool {
	return l.cursor >= len(l.src)
}

func (l *Lexer) genToken(kind, contentLength int) {
	l.curToken = &Token{
		Kind:    kind,
		Line:    l.line,
		Kth:     l.kth,
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
