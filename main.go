package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// TokenType represents different types of JSON tokens
type TokenType int

const (
	TokenEOF          TokenType = iota
	TokenLeftBrace              // {
	TokenRightBrace             // }
	TokenLeftBracket            // [
	TokenRightBracket           // ]
	TokenColon                  // :
	TokenComma                  // ,
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNull
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

// Lexer performs lexical analysis
type Lexer struct {
	input string
	pos   int
	ch    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
	}
	l.pos++
}

func (l *Lexer) peekChar() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() (string, error) {
	var result []rune
	l.readChar() // skip opening "

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '"', '\\', '/':
				result = append(result, rune(l.ch))
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'u':
				// Unicode escape sequence
				l.readChar()
				hex := ""
				for i := 0; i < 4; i++ {
					hex += string(l.ch)
					l.readChar()
				}
				val, err := strconv.ParseInt(hex, 16, 32)
				if err != nil {
					return "", fmt.Errorf("invalid unicode escape: %s", hex)
				}
				result = append(result, rune(val))
				continue
			default:
				return "", fmt.Errorf("invalid escape sequence: \\%c", l.ch)
			}
			l.readChar()
		} else {
			result = append(result, rune(l.ch))
			l.readChar()
		}
	}

	if l.ch != '"' {
		return "", fmt.Errorf("unterminated string")
	}
	l.readChar() // skip closing "
	return string(result), nil
}

func (l *Lexer) readNumber() string {
	start := l.pos - 1

	if l.ch == '-' {
		l.readChar()
	}

	if l.ch == '0' {
		l.readChar()
	} else {
		for unicode.IsDigit(rune(l.ch)) {
			l.readChar()
		}
	}

	if l.ch == '.' {
		l.readChar()
		for unicode.IsDigit(rune(l.ch)) {
			l.readChar()
		}
	}

	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for unicode.IsDigit(rune(l.ch)) {
			l.readChar()
		}
	}

	return l.input[start : l.pos-1]
}

func (l *Lexer) readIdentifier() string {
	start := l.pos - 1
	for unicode.IsLetter(rune(l.ch)) {
		l.readChar()
	}
	return l.input[start : l.pos-1]
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	tok := Token{Pos: l.pos - 1}

	switch l.ch {
	case 0:
		tok.Type = TokenEOF
	case '{':
		tok.Type = TokenLeftBrace
		l.readChar()
	case '}':
		tok.Type = TokenRightBrace
		l.readChar()
	case '[':
		tok.Type = TokenLeftBracket
		l.readChar()
	case ']':
		tok.Type = TokenRightBracket
		l.readChar()
	case ':':
		tok.Type = TokenColon
		l.readChar()
	case ',':
		tok.Type = TokenComma
		l.readChar()
	case '"':
		str, err := l.readString()
		if err != nil {
			return tok, err
		}
		tok.Type = TokenString
		tok.Value = str
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok.Type = TokenNumber
		tok.Value = l.readNumber()
	default:
		if unicode.IsLetter(rune(l.ch)) {
			ident := l.readIdentifier()
			switch ident {
			case "true":
				tok.Type = TokenTrue
			case "false":
				tok.Type = TokenFalse
			case "null":
				tok.Type = TokenNull
			default:
				return tok, fmt.Errorf("unexpected identifier: %s", ident)
			}
		} else {
			return tok, fmt.Errorf("unexpected character: %c", l.ch)
		}
	}

	return tok, nil
}

// Parser builds data structures from tokens
type Parser struct {
	lexer    *Lexer
	curToken Token
}

func NewParser(input string) (*Parser, error) {
	p := &Parser{lexer: NewLexer(input)}
	tok, err := p.lexer.NextToken()
	if err != nil {
		return nil, err
	}
	p.curToken = tok
	return p, nil
}

func (p *Parser) advance() error {
	tok, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.curToken = tok
	return nil
}

func (p *Parser) Parse() (interface{}, error) {
	return p.parseValue()
}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.curToken.Type {
	case TokenLeftBrace:
		return p.parseObject()
	case TokenLeftBracket:
		return p.parseArray()
	case TokenString:
		val := p.curToken.Value
		p.advance()
		return val, nil
	case TokenNumber:
		val, err := strconv.ParseFloat(p.curToken.Value, 64)
		if err != nil {
			return nil, err
		}
		p.advance()
		return val, nil
	case TokenTrue:
		p.advance()
		return true, nil
	case TokenFalse:
		p.advance()
		return false, nil
	case TokenNull:
		p.advance()
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.curToken.Type)
	}
}

func (p *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	p.advance() // skip {

	if p.curToken.Type == TokenRightBrace {
		p.advance()
		return obj, nil
	}

	for {
		if p.curToken.Type != TokenString {
			return nil, fmt.Errorf("expected string key, got %v", p.curToken.Type)
		}

		key := p.curToken.Value
		p.advance()

		if p.curToken.Type != TokenColon {
			return nil, fmt.Errorf("expected colon after key")
		}
		p.advance()

		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = val

		if p.curToken.Type == TokenRightBrace {
			p.advance()
			return obj, nil
		}

		if p.curToken.Type != TokenComma {
			return nil, fmt.Errorf("expected comma or closing brace")
		}
		p.advance()
	}
}

func (p *Parser) parseArray() ([]interface{}, error) {
	arr := []interface{}{}

	p.advance() // skip [

	if p.curToken.Type == TokenRightBracket {
		p.advance()
		return arr, nil
	}

	for {
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arr = append(arr, val)

		if p.curToken.Type == TokenRightBracket {
			p.advance()
			return arr, nil
		}

		if p.curToken.Type != TokenComma {
			return nil, fmt.Errorf("expected comma or closing bracket")
		}
		p.advance()
	}
}

// Public API
func Parse(input string) (interface{}, error) {
	parser, err := NewParser(input)
	if err != nil {
		return nil, err
	}
	return parser.Parse()
}

// Example usage
func main() {
	testCases := []string{
		`{"name": "John", "age": 30, "active": true}`,
		`[1, 2, 3, "hello", null, false]`,
		`{"user": {"name": "Alice", "scores": [95, 87, 92]}}`,
		`{"unicode": "Hello \u0057orld"}`,
		`{"number": -123.45e-6}`,
	}

	for i, tc := range testCases {
		fmt.Printf("\nTest case %d: %s\n", i+1, tc)
		result, err := Parse(tc)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Result: %+v\n", result)
		}
	}
}
