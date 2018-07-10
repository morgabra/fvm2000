package fvm2000

import (
	"fmt"
	"io"
)

/*
 * [LABEL:] INSTRUCTION SRC [, DST]
 */
type Token struct {
	tok  TokenType
	lit  string
	line int
	pos  int
}

type Parser struct {
	s *Scanner

	unscanned bool
	buf       Token
}

type ParseError struct {
	msg string
	tok Token
}

func NewParseError(msg string, tok Token) *ParseError {
	return &ParseError{
		msg: msg,
		tok: tok,
	}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line:%d pos:%d msg:%s", e.tok.line, e.tok.pos, e.msg)
}

type Instruction struct {
	lbl Token
	ins Token
	src Token
	dst Token
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r), buf: Token{}}
}

func (p *Parser) scan() Token {
	if p.unscanned {
		p.unscanned = false
		return p.buf
	}

	line, pos, tok, lit := p.s.Scan()

	// Skip whitespace, scanner is greedy with WS so the next scan will be some other token.
	if tok == WS {
		line, pos, tok, lit = p.s.Scan()
	}

	p.buf.line, p.buf.pos, p.buf.tok, p.buf.lit = line, pos, tok, lit
	return p.buf
}

func (p *Parser) unscan() { p.unscanned = true }

func (p *Parser) parse() (*Instruction, error) {
	st := &Instruction{}

	tok := p.scan()

	// Optional label
	if tok.tok == LABEL {
		st.lbl = tok
		tok = p.scan()
	}

	if tok.tok.isEnd() {
		return st, nil
	}

	// First token should be an instruction
	st.ins = tok
	if !tok.tok.isInstruction() {
		return nil, NewParseError("expected instruction", tok)
	}

	tok = p.scan()
	if tok.tok.isEnd() {
		return st, nil
	}

	// Second token should be a src
	st.src = tok
	if !st.src.tok.isSrc() {
		return nil, NewParseError("expected src address or const", st.src)
	}

	tok = p.scan()
	if tok.tok.isEnd() {
		return st, nil
	}

	if tok.tok != COMMA {
		return nil, NewParseError("expected delimiter", tok)
	}

	// Third token should be a dst
	tok = p.scan()
	st.dst = tok
	if !st.dst.tok.isDst() {
		return nil, NewParseError("expected dst address", st.dst)
	}

	tok = p.scan()
	if !tok.tok.isEnd() {
		return nil, NewParseError("unexpected token", tok)
	}

	return st, nil
}

func (p *Parser) Parse() ([]*Instruction, error) {
	st := []*Instruction{}

	for {
		tok := p.scan()

		// if EOF, we're done parsing
		if tok.tok == EOF {
			return st, nil
		}

		// skip empty lines
		if tok.tok == EOL {
			continue
		}

		p.unscan()
		s, err := p.parse()
		if err != nil {
			return nil, err
		}
		st = append(st, s)
	}
}
