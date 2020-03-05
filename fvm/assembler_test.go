package fvm

import (
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
	expectedBytes := []byte{9, 10, 10, 100, 15, 2}
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
	expectedBytes := []byte{9, 0, 11, 15, 18, 35, 14, 30, 17, 10, 1, 15, 6, 8}
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
	expectedBytes := []byte{9, 34, 11, 15, 18, 35, 19, 13, 14, 30, 15, 6, 8, 17, 20}
	out := parseAsm(p)

	require.Equal(t, expectedBytes, out[:len(expectedBytes)])
}
