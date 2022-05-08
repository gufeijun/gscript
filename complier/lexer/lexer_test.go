package lexer

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(int64(time.Now().Unix()))
}

func Test_IDENTIFIER_Token(t *testing.T) {
	const times = 50
	for i := 0; i < times; i++ {
		tl := &TokenList{line: 1}
		for j := 0; j < 50; j++ {
			for content := range IDENTIFIER_TOKEN_MAP {
				for roll, m := rand.Int()%3+1, 0; m < roll; m++ {
					tl.newToken(TOKEN_IDENTIFIER, content, content)
				}
			}
			for kw, kind := range keywords {
				for roll, m := rand.Int()%3+1, 0; m < roll; m++ {
					tl.newToken(kind, kw, kw)
				}
			}
		}
		match(tl.getTokens(), NewLexer("Test_IDENTIFIER_Token", []byte(tl.srcCode.String())), t)
	}
}

func Test_STRING_Token(t *testing.T) {
	const times = 50
	for i := 0; i < times; i++ {
		tl := &TokenList{line: 1}
		for j := 0; j < 50; j++ {
			for content, value := range STRING_TOKEN_MAP {
				for roll, m := rand.Int()%3+1, 0; m < roll; m++ {
					tl.newToken(TOKEN_STRING, content, value)
				}
			}
		}
		match(tl.getTokens(), NewLexer("Test_STRING_Token", []byte(tl.srcCode.String())), t)
	}
}

func Test_NUMBER_Token(t *testing.T) {
	const times = 50
	for i := 0; i < times; i++ {
		tl := &TokenList{line: 1}
		for j := 0; j < 50; j++ {
			for content, value := range NUMBER_TOKEN_MAP {
				for roll, m := rand.Int()%3+1, 0; m < roll; m++ {
					tl.newToken(TOKEN_NUMBER, content, value)
				}
			}
		}
		match(tl.getTokens(), NewLexer("Test_NUMBER_Token", []byte(tl.srcCode.String())), t)
	}
}

func Test_OP_SEP_Token(t *testing.T) {
	const times = 50
	for i := 0; i < times; i++ {
		tl := &TokenList{line: 1}
		for j := 0; j < 50; j++ {
			for kind, content := range OP_SEP_TOKEN_MAP {
				for roll, m := rand.Int()%3+1, 0; m < roll; m++ {
					tl.newToken(kind, content, nil)
				}
			}
		}
		match(tl.getTokens(), NewLexer("Test_OP_Token", []byte(tl.srcCode.String())), t)
	}
}

type TokenList struct {
	tokens    []*Token
	srcCode   strings.Builder
	line, kth int
}

func (tl *TokenList) getTokens() []*Token {
	if tl.tokens == nil {
		return nil
	}
	lastToken := tl.tokens[len(tl.tokens)-1]
	if lastToken.Line == tl.line && needAddSemi(lastToken.Kind) {
		tl.tokens = append(tl.tokens, &Token{
			Kind:    TOKEN_SEP_SEMI,
			Content: ";",
			Line:    lastToken.Line,
			Kth:     lastToken.Kth + len(lastToken.Content),
		})
	}
	return tl.tokens
}

func (tl *TokenList) newToken(kind int, content string, value interface{}) {
	tl.tokens = append(tl.tokens, &Token{
		Kind:    kind,
		Line:    tl.line,
		Kth:     tl.kth,
		Content: content,
		Value:   value,
	})
	tl.srcCode.WriteString(content)
	tl.kth += len(content)
	// try to write 1+ whiteSpace into source code
	tl.writeWhiteSpace()
}

func (tl *TokenList) writeWhiteSpace() {
	random := rand.Int() % 8
	if random == 0 { // 1/8 probability of writing break line
		tl.writeCRLF()
	} else { // 7/8 probability of writing '\t' or ' '
		tl.writeTabOrSpace(random%3 + 1)
	}
}

// write k '\t' or ' ' into srcCode
func (tl *TokenList) writeTabOrSpace(k int) {
	for i := 0; i < k; i++ {
		if rand.Int()%2 == 0 {
			tl.srcCode.WriteByte('\t')
		} else {
			tl.srcCode.WriteByte(' ')
		}
		tl.kth++
	}
}

// write random break lines into srcCode
func (tl *TokenList) writeCRLF() {
	base := 1
	for {
		roll := rand.Int() % base
		if roll != 0 {
			break
		}
		if tl.tokens != nil {
			latest := tl.tokens[len(tl.tokens)-1]
			if needAddSemi(latest.Kind) {
				tl.tokens = append(tl.tokens, &Token{
					Kind:    TOKEN_SEP_SEMI,
					Line:    latest.Line,
					Kth:     latest.Kth + len(latest.Content),
					Content: ";",
				})
			}
		}
		tl.line++
		tl.kth = 0
		// 1/base probability of writing break line
		if rand.Int()%2 == 0 {
			tl.srcCode.Write([]byte("\r\n"))
		} else {
			tl.srcCode.WriteByte('\n')
		}
		base *= 2
	}
}

func match(tokens []*Token, l *Lexer, t *testing.T) {
	for _, want := range tokens {
		got := l.NextToken()
		gotV, wantV := tokenValue(got), tokenValue(want)
		if gotV != wantV {
			t.Errorf("source code:\n %s\n", l.src)
			t.Fatalf("want token %s, but got %s\n", wantV, gotV)
		}
		if got.Content != want.Content {
			// t.Errorf("source code:\n %s\n", l.src)
			t.Fatalf("want content %s, but got %s\n", want.Content, got.Content)
		}
	}
	lastToken := l.NextToken()
	if lastToken.Kind != TOKEN_EOF {
		t.Fatalf("should reach EOF, but got one more token %s", tokenValue(lastToken))
	}
}

func tokenValue(token *Token) string {
	var tokenVal string
	if token.Kind == TOKEN_EOF {
		tokenVal = fmt.Sprintf("<EOF,->")
	} else if token.Kind >= TOKEN_SEP_DOT && token.Kind <= TOKEN_OP_SHR {
		tokenVal = fmt.Sprintf("<%s,->", token.Content)
	} else if token.Kind == TOKEN_IDENTIFIER {
		tokenVal = fmt.Sprintf("<identifier,%s>", token.Value)
	} else if token.Kind == TOKEN_NUMBER {
		tokenVal = fmt.Sprintf("<number,%v>", token.Value)
	} else if token.Kind == TOKEN_STRING {
		tokenVal = fmt.Sprintf("<string,%s>", token.Value)
	} else {
		tokenVal = fmt.Sprintf("<%s,->", token.Value)
	}
	return fmt.Sprintf("%s %d:%d", tokenVal, token.Line, token.Kth)
}

var OP_SEP_TOKEN_MAP = map[int]string{
	TOKEN_SEP_DOT:    ".",
	TOKEN_SEP_VARARG: "...",
	TOKEN_SEP_COLON:  ":",
	TOKEN_SEP_SEMI:   ";",
	TOKEN_SEP_COMMA:  ",",
	TOKEN_SEP_QMARK:  "?",
	TOKEN_SEP_LBRACK: "[",
	TOKEN_SEP_RBRACK: "]",
	TOKEN_SEP_LPAREN: "(",
	TOKEN_SEP_RPAREN: ")",
	TOKEN_SEP_LCURLY: "{",
	TOKEN_SEP_RCURLY: "}",
	TOKEN_OP_ADD:     "+",
	TOKEN_OP_SUB:     "-",
	TOKEN_OP_MUL:     "*",
	TOKEN_OP_DIV:     "/",
	TOKEN_OP_IDIV:    "//",
	TOKEN_OP_MOD:     "%",
	TOKEN_OP_LAND:    "&&",
	TOKEN_OP_LOR:     "||",
	TOKEN_OP_AND:     "&",
	TOKEN_OP_OR:      "|",
	TOKEN_OP_XOR:     "^",
	TOKEN_OP_NOT:     "~",
	TOKEN_OP_LNOT:    "!",
	TOKEN_OP_EQ:      "==",
	TOKEN_OP_NE:      "!=",
	TOKEN_OP_LT:      "<",
	TOKEN_OP_GT:      ">",
	TOKEN_OP_LE:      "<=",
	TOKEN_OP_GE:      ">=",
	TOKEN_OP_ASSIGN:  "=",
	TOKEN_OP_INC:     "++",
	TOKEN_OP_DEC:     "--",
	TOKEN_OP_ADDEQ:   "+=",
	TOKEN_OP_SUBEQ:   "-=",
	TOKEN_OP_MULEQ:   "*=",
	TOKEN_OP_DIVEQ:   "/=",
	TOKEN_OP_MODEQ:   "%=",
	TOKEN_OP_ANDEQ:   "&=",
	TOKEN_OP_OREQ:    "|=",
	TOKEN_OP_XOREQ:   "^=",
	TOKEN_OP_SHL:     "<<",
	TOKEN_OP_SHR:     ">>",
}

var NUMBER_TOKEN_MAP = map[string]interface{}{
	"0x1111":   0x1111,
	"0":        0,
	"0xFFFF":   0xffff,
	"0xffff":   0xffff,
	"0xabc8":   0xabc8,
	"0765":     0765,
	"01110":    01110,
	"985":      985,
	"211":      211,
	"07654321": 07654321,
	"996":      996,
	"0.786":    0.786,
	"10.212":   10.212,
	"1.212":    1.212,
	"19.212":   19.212,
}

var STRING_TOKEN_MAP = map[string]string{
	`""`:               "",
	`"abc"`:            "abc",
	`"hello world!"`:   "hello world!",
	`"hello\tworld!"`:  "hello\tworld!",
	`"\thello world!"`: "\thello world!",
	`"hello world!\t"`: "hello world!\t",
	`"\""`:             `"`,
	`"\'"`:             `'`,
	`"\n"`:             "\n",
	`"\t"`:             "\t",
	`"\\"`:             `\`,
	`"\\\\"`:           `\\`,
	`"\\\""`:           `\"`,
	`''`:               "",
	`'abc'`:            "abc",
	`'hello world!'`:   "hello world!",
	`'hello\tworld!'`:  "hello\tworld!",
	`'\"'`:             `"`,
	`'\''`:             `'`,
	`'\n'`:             "\n",
	`'\t'`:             "\t",
	`'\\'`:             `\`,
	`'\\\\'`:           `\\`,
	`'\\\"'`:           `\"`,
}

var IDENTIFIER_TOKEN_MAP = map[string]bool{
	"_abc":      true,
	"abc":       true,
	"a_b_c":     true,
	"__abc__":   true,
	"a123":      true,
	"a1_23":     true,
	"a1_23_bc8": true,
}
