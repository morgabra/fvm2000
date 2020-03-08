package fvm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAsm(t *testing.T) {
	p := `
LDA 10
add100:
ADC 100
BNE add100
`
	expectedBytes := []byte{
		LDA, 10,
		ADC, 100,
		BNE, 2, 0,
	}
	out := parseAsm(p)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}

func TestParseAsm2(t *testing.T) {
	p := `
		LDA 0
		STA 15
		LDY 35
		loop:
		  CMY 30
		  DEY
		  ADC 1
		  BNE loop
		BRK
	`
	expectedBytes := []byte{
		LDA, 0,
		STA, 15,
		LDY, 35,
		CMY, 30,
		DEY,
		ADC, 1,
		BNE, 6, 0,
		BRK,
	}
	out := parseAsm(p)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}

func TestParseAsm3(t *testing.T) {
	p := `
		LDA 34
		STA 15
		LDY 35
		loop:
          JSR decrementY
		  CMY 30
          BNE loop
		  BRK
decrementY:
DEY
RTS
	`
	expectedBytes := []byte{
		LDA, 34,
		STA, 15,
		LDY, 35,
		JSR, 15, 0,
		CMY, 30,
		BNE, 6, 0,
		BRK,
		DEY,
		RTS,
	}
	out := parseAsm(p)

	fmt.Println(out)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}

func TestParseAsm4(t *testing.T) {
	p := `
	LDA 34
    JMP 0x06
	ADC 1
	BRK
	`
	expectedBytes := []byte{
		LDA, 34,
		JMP, 6, 0,
		ADC, 1,
		BRK,
	}
	out := parseAsm(p)

	fmt.Println(out)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}

func TestParseAsmDefineLabels(t *testing.T) {
	p := `
	#define jump 0x06
	LDA 34
    JMP jump
	ADC 1
	BRK
	`
	expectedBytes := []byte{
		LDA, 34,
		JMP, 6, 0,
		ADC, 1,
		BRK,
	}
	out := parseAsm(p)

	fmt.Println(out)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}

func TestParseAsmDefineLabelsStore(t *testing.T) {
	p := `
	#define foo 0x21e4
	LDA 34
	STA foo
	`
	expectedBytes := []byte{
		LDA, 34,
		STA, 0xe4, 0x21,
	}
	out := parseAsm(p)

	fmt.Println(out)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}

func TestParseAsmLiterals(t *testing.T) {
	p := `
	#define foo 0x200
	MOV #500 foo
	MOV 0x200 0x744
	BRK
	`
	expectedBytes := []byte{
		MOVI, 0xf4, 0x01, 0x00, 0x02,
		MOVR, 0x00, 0x02, 0x44, 0x07,
		BRK,
	}
	out := parseAsm(p)

	fmt.Println(out)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}
func TestParseAsmADD(t *testing.T) {
	p := `
	#define foo 0x200
	ADD #500
	ADD 0x200
	BRK
	`
	expectedBytes := []byte{
		ADDI, 0xf4, 0x01,
		ADDR, 0x00, 0x02,
		BRK,
	}
	out := parseAsm(p)

	fmt.Println(out)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}

func TestParseAsmBuiltins(t *testing.T) {
	p := `
	#define foo 0x200
	ADD #500
	ADD 0x200
	ADD foo
	SUB ACC
	ADD #1
	MUL foo
	MUL #0
	BRK
	`
	expectedBytes := []byte{
		ADDI, 0xf4, 0x01,
		ADDR, 0x00, 0x02,
		ADDR, 0x00, 0x02,
		SUBR, byte(ACC & 0xFF), byte(ACC >> 8),
		ADDI, 0x01, 0x00,
		MULR, 0x00, 0x02,
		MULI, 0x00, 0x00,
		BRK,
	}
	out := parseAsm(p)

	fmt.Println(out)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}
