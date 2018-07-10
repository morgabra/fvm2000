package fvm2000

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

var eof = rune(0)

type TokenType int

func (t TokenType) isConst() bool {
	return t >= 20 && t <= 29
}

func (t TokenType) isRegister() bool {
	return t >= 30 && t <= 39
}

func (t TokenType) isGPIO() bool {
	return t >= 40 && t <= 49
}

func (t TokenType) isInstruction() bool {
	return t >= 50 && t <= 59
}

func (t TokenType) isSrc() bool {
	return t.isConst() || t.isRegister() || t.isGPIO()
}

func (t TokenType) isDst() bool {
	return t.isRegister() || t.isGPIO()
}

func (t TokenType) isEnd() bool {
	return t == EOF || t == EOL
}

const (
	ILLEGAL TokenType = TokenType(-1)

	NONE TokenType = TokenType(0)
	EOF  TokenType = TokenType(1)
	EOL  TokenType = TokenType(2)
	WS   TokenType = TokenType(3)

	COMMA TokenType = TokenType(10)
	LABEL TokenType = TokenType(11)

	INT   TokenType = TokenType(20)
	FLOAT TokenType = TokenType(21)

	PC  TokenType = TokenType(30)
	EAX TokenType = TokenType(31)
	EBX TokenType = TokenType(32)
	ECX TokenType = TokenType(33)
	EDX TokenType = TokenType(34)

	P0 TokenType = TokenType(40)
	P1 TokenType = TokenType(41)
	P2 TokenType = TokenType(42)
	P3 TokenType = TokenType(43)

	NOP TokenType = TokenType(50)
	MOV TokenType = TokenType(51)
	ADD TokenType = TokenType(52)
	SUB TokenType = TokenType(53)
)

type Scanner struct {
	r       *bufio.Reader
	linePos int
	lineNum int
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isNewline(ch rune) bool {
	return ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == ':'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isNumber(ch rune) bool {
	return isDigit(ch) || ch == '.' || ch == '-' || ch == '+'
}

func (s *Scanner) read() rune {
	s.linePos = s.linePos + 1
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) scanWhitespace() (TokenType, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (TokenType, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch strings.ToUpper(buf.String()) {
	default:
	case "PC":
		return PC, buf.String()

	case "EAX":
		return EAX, buf.String()
	case "EBX":
		return EBX, buf.String()
	case "ECX":
		return ECX, buf.String()
	case "EDX":
		return EDX, buf.String()

	case "P0":
		return P0, buf.String()
	case "P1":
		return P1, buf.String()
	case "P2":
		return P2, buf.String()
	case "P3":
		return P3, buf.String()

	case "NOP":
		return NOP, buf.String()
	case "MOV":
		return MOV, buf.String()
	case "ADD":
		return ADD, buf.String()
	case "SUB":
		return SUB, buf.String()
	}

	if strings.HasSuffix(buf.String(), ":") {
		return LABEL, buf.String()
	}

	return ILLEGAL, buf.String()
}

func (s *Scanner) scanNumber() (TokenType, string) {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNumber(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	_, err := strconv.ParseInt(buf.String(), 10, 64)
	if err == nil {
		return INT, buf.String()
	}
	_, err = strconv.ParseFloat(buf.String(), 64)
	if err == nil {
		return FLOAT, buf.String()
	}

	return ILLEGAL, buf.String()
}

func (s *Scanner) unread() {
	s.linePos = s.linePos - 1
	s.r.UnreadRune()
}

func (s *Scanner) Scan() (int, int, TokenType, string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		pos := s.linePos
		tok, lit := s.scanWhitespace()
		return s.lineNum, pos, tok, lit
	} else if isLetter(ch) {
		s.unread()
		pos := s.linePos
		tok, lit := s.scanIdent()
		return s.lineNum, pos, tok, lit
	} else if isNumber(ch) {
		s.unread()
		pos := s.linePos
		tok, lit := s.scanNumber()
		return s.lineNum, pos, tok, lit
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return 0, s.linePos - 1, EOF, "EOF"
	case '\n':
		start := s.linePos - 1
		s.linePos = 0
		s.lineNum = s.lineNum + 1
		return s.lineNum - 1, start, EOL, string(ch)
	case ',':
		return s.lineNum, s.linePos - 1, COMMA, string(ch)
	}

	return s.lineNum, s.linePos - 1, ILLEGAL, string(ch)
}
