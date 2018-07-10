package fvm2000

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func validateStatement(t *testing.T, st *Instruction, lbl, ins, src, dst *Token, msg string) {
	if lbl == nil {
		require.Equal(t, NONE, st.lbl.tok, msg)
		require.Equal(t, "", st.lbl.lit, msg)
		require.Equal(t, 0, st.lbl.line, msg)
		require.Equal(t, 0, st.lbl.pos, msg)
	} else {
		require.Equal(t, lbl.tok, st.lbl.tok, msg)
		require.Equal(t, lbl.lit, st.lbl.lit, msg)
		require.Equal(t, lbl.line, st.lbl.line, msg)
		require.Equal(t, lbl.pos, st.lbl.pos, msg)
	}

	if ins == nil {
		require.Equal(t, NONE, st.ins.tok, msg)
		require.Equal(t, "", st.ins.lit, msg)
		require.Equal(t, 0, st.ins.line, msg)
		require.Equal(t, 0, st.ins.pos, msg)
	} else {
		require.Equal(t, ins.tok, st.ins.tok, msg)
		require.Equal(t, ins.lit, st.ins.lit, msg)
		require.Equal(t, ins.line, st.ins.line, msg)
		require.Equal(t, ins.pos, st.ins.pos, msg)
	}

	if src == nil {
		require.Equal(t, NONE, st.src.tok, msg)
		require.Equal(t, "", st.src.lit, msg)
		require.Equal(t, 0, st.src.line, msg)
		require.Equal(t, 0, st.src.pos, msg)
	} else {
		require.Equal(t, src.tok, st.src.tok, msg)
		require.Equal(t, src.lit, st.src.lit, msg)
		require.Equal(t, src.line, st.src.line, msg)
		require.Equal(t, src.pos, st.src.pos, msg)
	}

	if dst == nil {
		require.Equal(t, NONE, st.dst.tok, msg)
		require.Equal(t, "", st.dst.lit, msg)
		require.Equal(t, 0, st.dst.line, msg)
		require.Equal(t, 0, st.dst.pos, msg)
	} else {
		require.Equal(t, dst.tok, st.dst.tok, msg)
		require.Equal(t, dst.lit, st.dst.lit, msg)
		require.Equal(t, dst.line, st.dst.line, msg)
		require.Equal(t, dst.pos, st.dst.pos, msg)
	}
}

func TestParserEmpty(t *testing.T) {
	c := ""

	p := NewParser(strings.NewReader(c))

	s, err := p.Parse()
	require.NoError(t, err)
	require.Equal(t, 0, len(s))

	// tab
	c = "\t"

	p = NewParser(strings.NewReader(c))

	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 0, len(s))

	// space
	c = ""

	p = NewParser(strings.NewReader(c))

	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 0, len(s))

	// mixed
	c = " \t"

	p = NewParser(strings.NewReader(c))

	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 0, len(s))

	// newline
	c = "\n"

	p = NewParser(strings.NewReader(c))

	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 0, len(s))

	// newlines
	c = "\n\n\n"

	p = NewParser(strings.NewReader(c))

	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 0, len(s))

	// chaos monkey
	c = "        \n\t    \n\t\t\t\t\t   \n     "

	p = NewParser(strings.NewReader(c))

	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 0, len(s))
}

func TestParserLabel(t *testing.T) {

	var expectedLabel *Token
	var resetFixture = func() {
		expectedLabel = &Token{
			tok:  LABEL,
			lit:  "label:",
			line: 0,
			pos:  0,
		}
	}
	resetFixture()

	// base case
	c := "label:"
	p := NewParser(strings.NewReader(c))
	s, err := p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	validateStatement(t, s[0], expectedLabel, nil, nil, nil, "base case")
	resetFixture()

	// leading spaces
	c = "     label:"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	expectedLabel.pos = 5
	validateStatement(t, s[0], expectedLabel, nil, nil, nil, "leading spaces")
	resetFixture()

	// leading tabs
	c = "\tlabel:"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	expectedLabel.pos = 1
	validateStatement(t, s[0], expectedLabel, nil, nil, nil, "leading tabs")
	resetFixture()

	// leading newline
	c = "\nlabel:"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	expectedLabel.line = 1
	validateStatement(t, s[0], expectedLabel, nil, nil, nil, "leading newlines")
	resetFixture()

	// labeled instruction
	c = "label: add eax"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	validateStatement(t, s[0], expectedLabel, &Token{ADD, "add", 0, 7}, &Token{EAX, "eax", 0, 11}, nil, "labeled instruction")
	resetFixture()

	// labeled instructions
	c = "label: add eax\nlabeltwo: mov eax, ebx"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 2, len(s))

	validateStatement(t, s[0], expectedLabel, &Token{ADD, "add", 0, 7}, &Token{EAX, "eax", 0, 11}, nil, "multi-line")
	expectedLabel.line = 1
	expectedLabel.lit = "labeltwo:"
	validateStatement(t, s[1], expectedLabel, &Token{MOV, "mov", 1, 10}, &Token{EAX, "eax", 1, 14}, &Token{EBX, "ebx", 1, 19}, "multi-line")
	resetFixture()
}

func TestParserNumber(t *testing.T) {

	// int
	c := "add 1"
	p := NewParser(strings.NewReader(c))
	s, err := p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{INT, "1", 0, 4}, nil, "int")

	// float
	c = "add 1.1"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{FLOAT, "1.1", 0, 4}, nil, "float")

	// neg
	c = "add -1"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{INT, "-1", 0, 4}, nil, "float")

	// neg float
	c = "add -1.1"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))

	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{FLOAT, "-1.1", 0, 4}, nil, "float")

	// invalid stuff
	c = "add 1.1-"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.Error(t, err)
	require.Equal(t, "line:0 pos:4 msg:expected src address or const", err.Error())

	c = "add 1.1.1"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.Error(t, err)
	require.Equal(t, "line:0 pos:4 msg:expected src address or const", err.Error())

}

func TestParserInstruction(t *testing.T) {

	// labeled
	c := "label: nop"
	p := NewParser(strings.NewReader(c))
	s, err := p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))
	validateStatement(t, s[0], &Token{LABEL, "label:", 0, 0}, &Token{NOP, "nop", 0, 7}, nil, nil, "labeled")

	// base case
	c = "nop"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))
	validateStatement(t, s[0], nil, &Token{NOP, "nop", 0, 0}, nil, nil, "base case")

	// src-only register
	c = "add eax"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))
	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{EAX, "eax", 0, 4}, nil, "src-only")

	// src-only gpio
	c = "add p0"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))
	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{P0, "p0", 0, 4}, nil, "src-only")

	// src-only const
	c = "add 1.1"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))
	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{FLOAT, "1.1", 0, 4}, nil, "src-only")

	// invalid src
	c = "add foobar"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.Error(t, err)
	require.Equal(t, "line:0 pos:4 msg:expected src address or const", err.(*ParseError).Error())

	// missing delimiter
	c = "add eax eax"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.Error(t, err)
	require.Equal(t, "line:0 pos:8 msg:expected delimiter", err.(*ParseError).Error())

	// src and dst
	c = "add eax, ebx"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))
	validateStatement(t, s[0], nil, &Token{ADD, "add", 0, 0}, &Token{EAX, "eax", 0, 4}, &Token{EBX, "ebx", 0, 9}, "src and dst")

	// invalid dst
	c = "add eax, 4.0"
	p = NewParser(strings.NewReader(c))
	s, err = p.Parse()
	require.Error(t, err)
	require.Equal(t, "line:0 pos:9 msg:expected dst address", err.(*ParseError).Error())
}
