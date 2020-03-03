package fvm

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

const ramSize = 256

type Fvm struct {
	Memory
	log    *zap.Logger
	cycles uint64
	pc     byte
	ram    Ram
	a      byte
	x      byte
	y      byte
	z      byte
	table  [256]func()
}

func (f *Fvm) newTable() {
	f.table = [256]func(){
		f.nop, f.mov, f.add, f.sub, f.mul, f.not, f.dgt, f.dst,
		f.brk, f.lda, f.adc, f.sta, f.ldx, f.inx, f.cmy, f.bne,
		f.sta_x, f.dey, f.ldy,
	}
}

func (f *Fvm) Read(addr byte) byte {
	switch {
	case addr <= 0xFF:
		return f.ram[addr]

	default:
		log.Fatalf("invalid memory read at address 0x%02X", addr)
	}
	return 0
}

func (f *Fvm) Read16(addr byte) uint16 {
	lo := uint16(f.Read(addr))
	hi := uint16(f.Read(addr + 1))
	return hi<<8 | lo
}

func (f *Fvm) Write(addr byte, value byte) {
	switch {
	case addr <= 0xFF:
		f.ram[addr] = value

	default:
		log.Fatalf("invalid memory write at address 0x%02X", addr)
	}
}

func (f *Fvm) dumpState() {
	f.log.Debug("cpu state",
		zap.Uint8("pc", f.current()),
		zap.Uint8("pc_value", f.Read(f.current())),
		zap.Uint8("a", f.a),
		zap.Uint8("x", f.x),
		zap.Uint8("y", f.y))
}

func (f *Fvm) dumpRam() {
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			fmt.Printf(" %02X ", f.ram[(i*16)+j])
		}
		fmt.Println()
	}
}

func (f *Fvm) advance() {
	f.pc++
}

func (f *Fvm) current() byte {
	return f.pc
}

func (f *Fvm) Step() int {
	cycles := f.cycles

	op := f.Read(f.current())

	f.table[op]()
	f.cycles += 1 //uint64(instructionCycles[op])

	f.dumpState()
	return int(f.cycles - cycles)
}

func (f *Fvm) Run() int {
	cycles := 0
	f.dumpState()
	for f.Read(f.current()) != BRK {
		cycles += f.Step()
	}
	f.dumpState()

	return cycles
}

func NewCPU(instructions []byte) (*Fvm, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	ram := make([]byte, ramSize)
	copy(ram, instructions)

	f := &Fvm{
		log: log.Named("fvm"),
		ram: ram,
	}
	f.newTable()

	return f, nil
}
