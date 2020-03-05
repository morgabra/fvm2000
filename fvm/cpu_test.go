package fvm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFvm_Run(t *testing.T) {
	program := `
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

	asm := Assemble(program)
	fmt.Println(asm)

	c, err := NewCPU(asm)
	require.NoError(t, err)
	c.Run()

	require.EqualValues(t, 12, c.pc)
	require.EqualValues(t, 34, c.a)
	require.EqualValues(t, 0, c.x)
	require.EqualValues(t, 30, c.y)
	require.EqualValues(t, 1, c.z)
}
