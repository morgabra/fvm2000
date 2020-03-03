package fvm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFvm_Run(t *testing.T) {
	program := `
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

	asm := Assemble(program)

	c, err := NewCPU(asm)
	require.NoError(t, err)

	c.Run()
}
