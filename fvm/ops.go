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
}

func (f *Fvm) brk() {
}

func (f *Fvm) lda() {
	f.advance()
	f.a = f.Read(f.current())
	f.advance()
}

func (f *Fvm) adc() {
	f.advance()
	f.a += f.Read(f.current())
	f.advance()
}

func (f *Fvm) sta() {
	f.advance()
	f.Write(f.Read(f.current()), f.a)
	f.advance()
}

func (f *Fvm) ldx() {
	f.advance()
	f.x = f.Read(f.current())
	f.advance()
}

func (f *Fvm) inx() {
	f.x++
	f.advance()
}

func (f *Fvm) cmy() {
	f.advance()
	if f.y == f.Read(f.current()) {
		f.z = 1
	} else {
		f.z = 0
	}
	f.advance()
}

func (f *Fvm) bne() {
	f.advance()
	if f.z == 0 {
		f.pc = f.Read(f.current())
		return
	}
	f.advance()
}

func (f *Fvm) sta_x() {
	f.Write(f.x, f.a)
	f.advance()
}

func (f *Fvm) dey() {
	f.y--
	f.advance()
}

func (f *Fvm) ldy() {
	f.advance()
	f.y = f.Read(f.current())
	f.advance()
}

func (f *Fvm) nop() {
}

func (f *Fvm) mov() {
}

func (f *Fvm) add() {
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
