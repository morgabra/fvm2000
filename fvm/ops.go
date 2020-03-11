package fvm

const (
	NOP = iota
	MOVI
	MOVR
	ADDI
	ADDR
	SUBI
	SUBR
	MULI
	MULR
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
	JMP
)

const (
	modeAbsolute = iota + 1
	modeImmediate
	modeAccumulator
	modeZeroPage
	modeImplicit
)

var opNames = map[string]byte{
	"NOP":   NOP,
	"MOV":   MOVR,
	"ADD":   ADDR,
	"SUB":   SUBR,
	"MUL":   MULR,
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
	"JMP":   JMP,
}

type Op struct {
	name     string
	size     byte
	cycles   byte
	readMode byte
}

type opInfo struct {
	addr uint16
	mode byte
	pc   uint16
}

var ops = map[byte]Op{
	NOP: {
		name:     "NOP",
		size:     1,
		cycles:   1,
		readMode: modeImplicit,
	},
	MOVR: {
		name:     "MOV",
		size:     5,
		cycles:   1,
		readMode: modeAbsolute,
	},
	MOVI: {
		name:     "MOV",
		size:     5,
		cycles:   1,
		readMode: modeImmediate,
	},
	ADDR: {
		name:     "ADD",
		size:     3,
		cycles:   1,
		readMode: modeAbsolute,
	},
	ADDI: {
		name:     "ADD",
		size:     3,
		cycles:   1,
		readMode: modeImmediate,
	},
	SUBR: {
		name:     "SUB",
		size:     3,
		cycles:   1,
		readMode: modeAbsolute,
	},
	SUBI: {
		name:     "SUB",
		size:     3,
		cycles:   1,
		readMode: modeImmediate,
	},
	MULR: {
		name:     "MUL",
		size:     3,
		cycles:   1,
		readMode: modeAbsolute,
	},
	MULI: {
		name:     "MUL",
		size:     3,
		cycles:   1,
		readMode: modeImmediate,
	},
	NOT: {
		name:     "NOT",
		size:     1,
		cycles:   1,
		readMode: modeZeroPage,
	},
	DGT: {
		name:     "DGT",
		size:     2,
		cycles:   1,
		readMode: modeImmediate,
	},
	DST: {
		name:     "DST",
		size:     3,
		cycles:   1,
		readMode: modeImmediate,
	},
	BRK: {
		name:     "BRK",
		size:     1,
		cycles:   1,
		readMode: modeZeroPage,
	},
	LDA: {
		name:     "LDA",
		size:     2,
		cycles:   1,
		readMode: modeImmediate,
	},
	ADC: {
		name:     "ADC",
		size:     2,
		cycles:   1,
		readMode: modeImmediate,
	},
	STA: {
		name:     "STA",
		size:     2,
		cycles:   1,
		readMode: modeImmediate,
	},
	LDX: {
		name:     "LDX",
		size:     2,
		cycles:   1,
		readMode: modeImmediate,
	},
	INX: {
		name:     "INX",
		size:     1,
		cycles:   1,
		readMode: modeImplicit,
	},
	CMY: {
		name:     "CMY",
		size:     2,
		cycles:   1,
		readMode: modeImmediate,
	},
	BNE: {
		name:     "BNE",
		size:     2,
		cycles:   1,
		readMode: modeAbsolute,
	},
	STA_X: {
		name:     "STA_X",
		size:     1,
		cycles:   1,
		readMode: modeImplicit,
	},
	DEY: {
		name:     "DEY",
		size:     1,
		cycles:   1,
		readMode: modeImplicit,
	},
	LDY: {
		name:     "LDY",
		size:     2,
		cycles:   1,
		readMode: modeImmediate,
	},
	JSR: {
		name:     "JSR",
		size:     3,
		cycles:   1,
		readMode: modeImmediate,
	},
	RTS: {
		name:     "RTS",
		size:     1,
		cycles:   1,
		readMode: modeImplicit,
	},
	JMP: {
		name:     "JMP",
		size:     3,
		cycles:   1,
		readMode: modeImmediate,
	},
}

func (f *Fvm) brk(i opInfo) {
}

func (f *Fvm) lda(i opInfo) {
	f.a = f.Read(i.addr)
}

func (f *Fvm) adc(i opInfo) {
	f.a += f.Read(i.addr)
}

func (f *Fvm) sta(i opInfo) {
	f.Write(uint16(f.Read(i.addr)), f.a)
}

func (f *Fvm) ldx(i opInfo) {
	f.x = f.Read(i.addr)
}

func (f *Fvm) inx(i opInfo) {
	f.x++
}

func (f *Fvm) cmy(i opInfo) {
	if f.y == f.Read(i.addr) {
		f.z = 1
	} else {
		f.z = 0
	}
}

func (f *Fvm) bne(i opInfo) {
	if f.z == 0 {
		f.pc = i.addr
		return
	}
}

func (f *Fvm) sta_x(i opInfo) {
	f.Write(uint16(f.x), f.a)
}

func (f *Fvm) dey(i opInfo) {
	f.y--
}

func (f *Fvm) ldy(i opInfo) {
	f.y = f.Read(i.addr)
}

func (f *Fvm) jsr(i opInfo) {
	f.stackPush16(f.pc)
	f.pc = f.Read16(i.addr)
}

func (f *Fvm) rts(i opInfo) {
	f.pc = f.stackPull16()
}

func (f *Fvm) jmp(i opInfo) {
	f.pc = f.Read16(i.addr)
}

func (f *Fvm) nop(i opInfo) {}

func (f *Fvm) mov(i opInfo) {
	if i.mode == modeImmediate {
		f.Write16(f.Read16(i.addr+2), f.Read16(i.addr))
	} else if i.mode == modeAbsolute {
		f.Write16(f.Read16(f.pc-2), f.Read16(i.addr))
	}
}

func (f *Fvm) add(i opInfo) {
	f.Write16(ACC, f.Read16(ACC)+f.Read16(i.addr))
}

func (f *Fvm) sub(i opInfo) {
	f.Write16(ACC, f.Read16(ACC)-f.Read16(i.addr))
}

func (f *Fvm) mul(i opInfo) {
	f.Write16(ACC, f.Read16(ACC)*f.Read16(i.addr))
}

func (f *Fvm) not(i opInfo) {
	acc := f.Read16(ACC)
	if acc == 0 {
		f.Write16(ACC, 100)
	} else {
		f.Write16(ACC, 0)
	}
}

func (f *Fvm) dgt(i opInfo) {
}

func (f *Fvm) dst(i opInfo) {
}
