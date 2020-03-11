package fvm

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

const ramSize = 2048

const (
	// registers
	ACC uint16 = 0x100

	// inputs
	P0  uint16 = 0x110
	P1         = 0x112
	P2         = 0x114
	P3         = 0x116
	EAX        = 0x118
	EBX        = 0x11A
	ECX        = 0x11C
	EDX        = 0x11E
)

var builtins = map[string]uint16{
	"ACC": ACC,
	"P0":  P0,
	"P1":  P1,
	"P2":  P2,
	"P3":  P3,
	"EAX": EAX,
	"EBX": EBX,
	"ECX": ECX,
	"EDX": EDX,
}

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
	table  [256]func(i opInfo)
}

func (f *Fvm) newTable() {
	f.table = [256]func(opInfo){
		f.nop,
		f.mov, f.mov,
		f.add, f.add,
		f.sub, f.sub,
		f.mul, f.mul,
		f.not,
		f.dgt,
		f.dst,
		f.brk,
		f.lda,
		f.adc,
		f.sta,
		f.ldx,
		f.inx,
		f.cmy,
		f.bne,
		f.sta_x,
		f.dey,
		f.ldy,
		f.jsr,
		f.rts,
		f.jmp,
	}
}

func (f *Fvm) Read(addr uint16) byte {
	switch {
	case addr < ramSize:
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

// Stack begins at last byte of second page of ram
func (f *Fvm) stackPush(v byte) {
	f.Write(uint16(f.sp)&0x1FF, v)
	f.sp--
}

func (f *Fvm) stackPush16(v uint16) {
	hi := byte(v >> 8)
	lo := byte(v & 0xFF)
	f.stackPush(hi)
	f.stackPush(lo)
}

// Stack begins at last byte of second page of ram
func (f *Fvm) stackPull() byte {
	f.sp++
	v := f.Read(uint16(f.sp) & 0x1FF)
	return v
}

func (f *Fvm) stackPull16() uint16 {
	lo := uint16(f.stackPull())
	hi := uint16(f.stackPull())
	return hi<<8 | lo
}

func (f *Fvm) Write(addr uint16, value byte) {
	switch {
	case addr < ramSize:
		f.ram[addr] = value

	default:
		log.Fatalf("invalid memory write at address 0x%02X", addr)
	}
}

func (f *Fvm) Write16(addr uint16, value uint16) {
	lo := byte(value & 0xFF)
	hi := byte(value >> 8)
	f.Write(addr, lo)
	f.Write(addr+1, hi)
}

func (f *Fvm) dumpState() {
	op := ops[f.Read(f.pc)]
	f.log.Debug("cpu state",
		zap.Uint16("pc", f.pc),
		zap.Uint8("pc_value", f.Read(f.pc)),
		zap.Uint8("op", f.Read(f.pc)),
		zap.String("op", op.name),
		zap.Uint8("a", f.a),
		zap.Uint8("x", f.x),
		zap.Uint8("y", f.y),
		zap.Uint8("z", f.z),
		zap.Uint8("sp", f.sp),
		zap.Uint8("sp_value", f.Read(uint16(f.sp+1))),
		zap.Uint16("acc", f.Read16(ACC)))
}

func (f *Fvm) dumpRam() {
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			fmt.Printf(" %02X ", f.ram[(i*16)+j])
		}
		fmt.Println()
	}
}

func (f *Fvm) Step() int {
	f.dumpState()

	cycles := f.cycles

	op := f.Read(f.pc)

	var addr uint16
	mode := ops[op].readMode
	switch mode {
	case modeAbsolute:
		addr = f.Read16(f.pc + 1)
	case modeImmediate:
		addr = f.pc + 1
	case modeZeroPage:
		addr = uint16(f.Read(f.pc + 1))
	case modeAccumulator:
		addr = 0
	case modeImplicit:
		addr = 0
	}

	f.pc += uint16(ops[op].size)
	f.cycles += 1

	f.table[op](opInfo{
		addr: addr,
		pc:   f.pc,
		mode: mode,
	})

	return int(f.cycles - cycles)
}

func (f *Fvm) Run() int {
	cycles := 0
	for f.Read(f.pc) != BRK {
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
	// Only copy the first 256 bytes of the instructions into ram
	copy(ram, instructions[:256])

	f := &Fvm{
		log: log.Named("fvm"),
		sp:  0xFF,
		ram: ram,
	}
	f.newTable()

	return f, nil
}
