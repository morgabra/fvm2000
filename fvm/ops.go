package fvm

const (
	NOP = iota
	MOV
	ADD
	SUB
	MUL
	NOT
	DGT
	DST
	BRK
	LDA
	ADC
	STA
	LDX
	INX
	CMY
	BNE
	STA_X
	DEY
	LDY
	JSR
	RTS
)

var opNames = map[string]byte{
	"NOP":   NOP,
	"MOV":   MOV,
	"ADD":   ADD,
	"SUB":   SUB,
	"MUL":   MUL,
	"NOT":   NOT,
	"DGT":   DGT,
	"DST":   DST,
	"BRK":   BRK,
	"LDA":   LDA,
	"ADC":   ADC,
	"STA":   STA,
	"LDX":   LDX,
	"INX":   INX,
	"CMY":   CMY,
	"BNE":   BNE,
	"STA_X": STA_X,
	"DEY":   DEY,
	"LDY":   LDY,
	"JSR":   JSR,
	"RTS":   RTS,
}

type Op struct {
	size   byte
	cycles byte
}

var ops = map[byte]Op{
	NOP: {
		size:   1,
		cycles: 1,
	},
	MOV: {
		size:   3,
		cycles: 1,
	},
	ADD: {
		size:   2,
		cycles: 1,
	},
	SUB: {
		size:   2,
		cycles: 1,
	},
	MUL: {
		size:   2,
		cycles: 1,
	},
	NOT: {
		size:   1,
		cycles: 1,
	},
	DGT: {
		size:   2,
		cycles: 1,
	},
	DST: {
		size:   3,
		cycles: 1,
	},
	BRK: {
		size:   1,
		cycles: 1,
	},
	LDA: {
		size:   2,
		cycles: 1,
	},
	ADC: {
		size:   2,
		cycles: 1,
	},
	STA: {
		size:   2,
		cycles: 1,
	},
	LDX: {
		size:   2,
		cycles: 1,
	},
	INX: {
		size:   1,
		cycles: 1,
	},
	CMY: {
		size:   2,
		cycles: 1,
	},
	BNE: {
		size:   2,
		cycles: 1,
	},
	STA_X: {
		size:   1,
		cycles: 1,
	},
	DEY: {
		size:   1,
		cycles: 1,
	},
	LDY: {
		size:   2,
		cycles: 1,
	},
	JSR: {
		size:   0,
		cycles: 1,
	},
	RTS: {
		size:   1,
		cycles: 1,
	},
}

func (f *Fvm) brk() {
}

func (f *Fvm) lda() {
	f.a = f.Read(f.current() + 1)
}

func (f *Fvm) adc() {
	f.a += f.Read(f.current() + 1)
}

func (f *Fvm) sta() {
	f.Write(uint16(f.Read(f.current()+1)), f.a)
}

func (f *Fvm) ldx() {
	f.x = f.Read(f.current() + 1)
}

func (f *Fvm) inx() {
	f.x++
}

func (f *Fvm) cmy() {
	if f.y == f.Read(f.current()+1) {
		f.z = 1
	} else {
		f.z = 0
	}
}

func (f *Fvm) bne() {
	if f.z == 0 {
		f.pc = uint16(f.Read(f.current()+1) - 2)
		return
	}
}

func (f *Fvm) sta_x() {
	f.Write(uint16(f.x), f.a)
}

func (f *Fvm) dey() {
	f.y--
}

func (f *Fvm) ldy() {
	f.y = f.Read(f.current() + 1)
}

func (f *Fvm) jsr() {
	f.stackPush16(f.current() + 1)
	f.pc = uint16(f.Read(f.current() + 1))
}

func (f *Fvm) rts() {
	f.pc = f.stackPull16()
}

func (f *Fvm) jmp() {
	f.pc = uint16(f.Read(f.current() + 1))
}

func (f *Fvm) nop() {}

func (f *Fvm) mov() {
}

func (f *Fvm) add() {
	f.Write(ACC, f.Read(ACC)+f.Read(f.current()+1))
}

func (f *Fvm) sub() {
}

func (f *Fvm) mul() {
}

func (f *Fvm) not() {
}

func (f *Fvm) dgt() {
}

func (f *Fvm) dst() {
}
