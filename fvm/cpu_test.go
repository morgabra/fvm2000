package fvm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func check16(t *testing.T, f *Fvm, addr, value uint16) {
	require.EqualValues(t, value&0xFF, f.Read(addr))
	require.EqualValues(t, value>>8, f.Read(addr+1))
}

func TestFvm_Run(t *testing.T) {
	program := `
		LDA 34
		STA 35
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

	require.EqualValues(t, 34, c.a)
	require.EqualValues(t, 0, c.x)
	require.EqualValues(t, 30, c.y)
	require.EqualValues(t, 1, c.z)
}

func TestFvm_JMP_ADDR(t *testing.T) {
	program := `
	LDA 34
    JMP 0x07
	ADC 1
	BRK
	`

	asm := Assemble(program)
	fmt.Println(asm)

	c, err := NewCPU(asm)
	require.NoError(t, err)
	c.Run()

	require.EqualValues(t, 34, c.a)
	require.EqualValues(t, 0, c.x)
	require.EqualValues(t, 0, c.y)
	require.EqualValues(t, 0, c.z)
}

func TestFvmLabelStore(t *testing.T) {
	p := `
	#define foo 0xe4
	LDA 34
	STA foo
	BRK
	`

	asm := Assemble(p)
	fmt.Println(asm)

	c, err := NewCPU(asm)
	require.NoError(t, err)

	c.Run()

	require.EqualValues(t, 34, c.a)
	require.EqualValues(t, 34, c.ram[0xe4])
}

func TestFvmMov(t *testing.T) {
	p := `
	#define foo 0x200
	#define bar 0x300
	MOV #500 foo
	MOV 0x200 0x744
	MOV 0x744 bar
	BRK
	`

	asm := Assemble(p)
	fmt.Println(asm)

	c, err := NewCPU(asm)
	require.NoError(t, err)

	c.Run()
	require.EqualValues(t, 0, c.a)
	check16(t, c, 0x200, 500)
	check16(t, c, 0x744, 500)
	check16(t, c, 0x300, 500)
}

func TestFvmAdd(t *testing.T) {
	p := `
	#define foo 0x200
	MOV #500 foo
	ADD foo
	ADD #100
	ADD #0x100
	`

	asm := Assemble(p)
	fmt.Println(asm)

	c, err := NewCPU(asm)
	require.NoError(t, err)

	c.Step()
	check16(t, c, 0x200, 500)
	check16(t, c, ACC, 0)

	c.Step()
	check16(t, c, 0x200, 500)
	check16(t, c, ACC, 500)

	c.Step()
	check16(t, c, 0x200, 500)
	check16(t, c, ACC, 600)

	c.Step()
	check16(t, c, 0x200, 500)
	check16(t, c, ACC, 856)
}

func TestFvmSub(t *testing.T) {
	p := `
	#define foo 0x200
	MOV #500 ACC
	MOV #100 foo
	SUB foo
	SUB #50
	SUB 0x100
	ADD #40
	MUL #2
	MUL ACC
	MOV #0 foo
	MUL foo
	NOT
	NOT
	MOV #42 ACC
	NOT
	BRK
	`

	asm := Assemble(p)
	fmt.Println(asm)

	c, err := NewCPU(asm)
	require.NoError(t, err)
	check16(t, c, ACC, 0)

	c.Step()
	check16(t, c, ACC, 500)

	c.Step()
	check16(t, c, 0x200, 100)
	check16(t, c, ACC, 500)

	c.Step()
	check16(t, c, 0x200, 100)
	check16(t, c, ACC, 400)

	c.Step()
	check16(t, c, 0x200, 100)
	check16(t, c, ACC, 350)

	c.Step()
	check16(t, c, 0x200, 100)
	check16(t, c, ACC, 0)

	c.Step()
	check16(t, c, 0x200, 100)
	check16(t, c, ACC, 40)

	c.Step()
	check16(t, c, 0x200, 100)
	check16(t, c, ACC, 80)

	c.Step()
	check16(t, c, 0x200, 100)
	check16(t, c, ACC, 6400)

	c.Step()
	check16(t, c, 0x200, 0)
	check16(t, c, ACC, 6400)

	c.Step()
	check16(t, c, 0x200, 0)
	check16(t, c, ACC, 0)

	c.Step()
	check16(t, c, 0x200, 0)
	check16(t, c, ACC, 100)

	c.Step()
	check16(t, c, 0x200, 0)
	check16(t, c, ACC, 0)

	c.Step()
	check16(t, c, 0x200, 0)
	check16(t, c, ACC, 42)

	c.Step()
	check16(t, c, 0x200, 0)
	check16(t, c, ACC, 0)
}
