package fvm

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

const ramSize = 256

const (
	// registers
	ACC = 0x30
	EAX = 0x32
	EBX = 0x34
	ECX = 0x36
	EDX = 0x38

	// inputs
	P0 = 0x40
	P1 = 0x42
	P2 = 0x44
	p3 = 0x46
)

type Fvm struct {
	Memory
	log    *zap.Logger
	cycles uint64
	pc     uint16
	sp     byte
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
		f.sta_x, f.dey, f.ldy, f.jsr, f.rts,
	}
}

func (f *Fvm) Read(addr uint16) byte {
	switch {
	case addr <= 0xFF:
		return f.ram[addr]

	default:
		log.Fatalf("invalid memory read at address 0x%02X", addr)
	}
	return 0
}

func (f *Fvm) Read16(addr uint16) uint16 {
	lo := uint16(f.Read(addr))
	hi := uint16(f.Read(addr + 1))
	return hi<<8 | lo
}

func (f *Fvm) stackPush(v byte) {
	f.Write(uint16(f.sp), v)
	f.sp--
}

func (f *Fvm) stackPush16(v uint16) {
	hi := byte(v >> 8)
	lo := byte(v & 0xFF)
	f.stackPush(hi)
	f.stackPush(lo)
}

func (f *Fvm) stackPull() byte {
	f.sp++
	v := f.Read(uint16(f.sp))
	return v
}

func (f *Fvm) stackPull16() uint16 {
	lo := uint16(f.stackPull())
	hi := uint16(f.stackPull())
	return hi<<8 | lo
}

func (f *Fvm) Write(addr uint16, value byte) {
	switch {
	case addr <= 0xFF:
		f.ram[addr&0xFF] = value

	default:
		log.Fatalf("invalid memory write at address 0x%02X", addr)
	}
}

func (f *Fvm) dumpState() {
	f.log.Debug("cpu state",
		zap.Uint16("pc", f.current()),
		zap.Uint8("pc_value", f.Read(f.current())),
		zap.Uint8("op", f.Read(uint16(f.current()))),
		zap.Uint8("a", f.a),
		zap.Uint8("x", f.x),
		zap.Uint8("y", f.y),
		zap.Uint8("z", f.z),
		zap.Uint8("sp", f.sp),
		zap.Uint8("sp_value", f.Read(uint16(f.sp+1))))
}

func (f *Fvm) dumpRam() {
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			fmt.Printf(" %02X ", f.ram[(i*16)+j])
		}
		fmt.Println()
	}
}

func (f *Fvm) current() uint16 {
	return f.pc
}

func (f *Fvm) Step() int {
	cycles := f.cycles

	op := f.Read(f.current())

	f.table[op]()
	f.cycles += 1
	f.pc += uint16(ops[op].size)
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
		sp:  0xff,
		ram: ram,
	}
	f.newTable()

	return f, nil
}
