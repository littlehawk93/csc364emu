package asm

const (
	// InstructionSize defines how many bytes a single instruction is
	InstructionSize int = 2
)

const (
	instructionMove  byte = 0x00
	instructionNot   byte = 0x10
	instructionAnd   byte = 0x20
	instructionOr    byte = 0x30
	instructionAdd   byte = 0x40
	instructionSub   byte = 0x50
	instructionAddi  byte = 0x60
	instructionSubi  byte = 0x70
	instructionSet   byte = 0x80
	instructionSeth  byte = 0x90
	instructionInciz byte = 0xA0
	instructionDecin byte = 0xB0
	instructionMovez byte = 0xC0
	instructionMovex byte = 0xD0
	instructionMovep byte = 0xE0
	instructionMoven byte = 0xF0
)

var registersMap = map[string]byte{
	"r0":   0x00,
	"r1":   0x01,
	"r2":   0x02,
	"r3":   0x03,
	"r4":   0x04,
	"r5":   0x05,
	"r6":   0x06,
	"r7":   0x07,
	"r8":   0x08,
	"r9":   0x09,
	"ra":   0x0A,
	"rb":   0x0B,
	"rc":   0x0C,
	"rd":   0x0D,
	"re":   0x0E,
	"rf":   0x0F,
	"in":   0x06,
	"out0": 0x0D,
	"out1": 0x0E,
	"pc":   0x0F,
	"r10":  0x0A,
	"r11":  0x0B,
	"r12":  0x0C,
	"r13":  0x0D,
	"r14":  0x0E,
	"r15":  0x0F,
}

var instructionsMap = map[string]byte{
	"mov":   instructionMove,
	"not":   instructionNot,
	"and":   instructionAnd,
	"or":    instructionOr,
	"add":   instructionAdd,
	"sub":   instructionSub,
	"addi":  instructionAddi,
	"subi":  instructionSubi,
	"set":   instructionSet,
	"seth":  instructionSeth,
	"incz":  instructionInciz,
	"decn":  instructionDecin,
	"movz":  instructionMovez,
	"movx":  instructionMovex,
	"movp":  instructionMovep,
	"movn":  instructionMoven,
	"move":  instructionMove,
	"inciz": instructionInciz,
	"decin": instructionDecin,
	"movez": instructionMovez,
	"movex": instructionMovex,
	"movep": instructionMovep,
	"moven": instructionMoven,
}
