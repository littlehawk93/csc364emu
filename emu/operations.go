package emu

import "fmt"

const (
	opCodeMove  = 0
	opCodeNot   = 1
	opCodeAnd   = 2
	opCodeOr    = 3
	opCodeAdd   = 4
	opCodeSub   = 5
	opCodeAddi  = 6
	opCodeSubi  = 7
	opCodeSet   = 8
	opCodeSeth  = 9
	opCodeInciz = 10
	opCodeDecin = 11
	opCodeMovez = 12
	opCodeMovex = 13
	opCodeMovep = 14
	opCodeMoven = 15
)

// emulatorOperation executes a single instruction for the Emulator. Returns true if the program counter should increment, false otherwise
type emulatorOperation func(e *Emulator, dest, optionA, optionB byte)

var operationsMap = map[byte]emulatorOperation{
	opCodeMove:  operationMove,
	opCodeNot:   operationNot,
	opCodeAnd:   operationAnd,
	opCodeOr:    operationOr,
	opCodeAdd:   operationAdd,
	opCodeSub:   operationSub,
	opCodeAddi:  operationAddi,
	opCodeSubi:  operationSubi,
	opCodeSet:   operationSet,
	opCodeSeth:  operationSeth,
	opCodeInciz: operationInciz,
	opCodeDecin: operationDecin,
	opCodeMovez: operationMovez,
	opCodeMovex: operationMovex,
	opCodeMovep: operationMovep,
	opCodeMoven: operationMoven,
}

func convertToBytes(instruction uint16) (byte, byte, byte, byte) {
	return byte((instruction >> 12) & 0x000F), byte((instruction >> 8) & 0x000F), byte((instruction >> 4) & 0x000F), byte(instruction & 0x000F)
}

func executeInstruction(e *Emulator, instruction, dest, optionA, optionB byte) error {

	operation, ok := operationsMap[instruction]

	if !ok {
		return fmt.Errorf("Unrecognized op code: 0x'%X'", instruction)
	}
	operation(e, dest, optionA, optionB)
	return nil
}

// operationMove performs the MOVE instruction. Copies the value of the register in optionA into the destination register
func operationMove(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, e.GetRegister(optionA))
}

// operationNot performs the NOT instruction. Gets the binary NOT value of the register in optionA and saves the result into the destination register
func operationNot(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, ^e.GetRegister(optionA))
}

// operationAnd performs the AND instruction. Gets the binary AND value of the two registers in options A and B and saves the result into the destination register
func operationAnd(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, e.GetRegister(optionA)&e.GetRegister(optionB))
}

// operationOr performs the OR instruction. Gets the binary OR value of the two registers in options A and B and saves the result into the destination register
func operationOr(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, e.GetRegister(optionA)|e.GetRegister(optionB))
}

// operationAdd performs the ADD instruction. Adds the value of the two registers in options A and B together and saves the result into the destination register
func operationAdd(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, e.GetRegister(optionA)+e.GetRegister(optionB))
}

// operationSub performs the SUB instruction. Subtracts the value of register in optionA by the value of the register in optionB and saves the result into the destination register
func operationSub(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, e.GetRegister(optionA)-e.GetRegister(optionB))
}

// operationAddi performs the ADDI instruction. Adds the 4 bit binary value in optionB to the value of the register in optionA and saves the result into the destination register
func operationAddi(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, e.GetRegister(optionA)+uint16(optionB))
}

// operationSubi performs the SUBI instruction. Subtracts the 4 bit binary value in optionB from the value of the register in optionA and saves the result into the destination register
func operationSubi(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, e.GetRegister(optionA)-uint16(optionB))
}

// operationSet performs the SET instruction. Clears the destination register and sets it's lower 8 bits to the values of optionA and optionB
func operationSet(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, uint16(0x00FF&(optionA<<4|(optionB&0x0F))))
}

// operationSeth performs the SETH instruction. Sets the upper 8 bits of the destination register to the values of optionA and optionB without clearing the lower bits of the register
func operationSeth(e *Emulator, dest, optionA, optionB byte) {
	e.SetRegister(dest, uint16(0xFF00&uint16(optionA<<4|(optionB&0x0F))<<8)|e.GetRegister(dest))
}

// operationInciz performs the INCIZ instruction. Adds the 4 bit binary value in optionA to the value of the destination register if the register in optionB is equal to 0
func operationInciz(e *Emulator, dest, optionA, optionB byte) {
	if e.GetRegister(optionB) == 0 {
		e.SetRegister(dest, e.GetRegister(dest)+uint16(optionA))
	}
}

// operationDecin performs the DECIN instruction. Subtracts the 4 bit binary value in optionA from the value of the destination register if the register in optionB is negative (most significant bit is 1)
func operationDecin(e *Emulator, dest, optionA, optionB byte) {
	if e.GetRegister(optionB)&0x8000 != 0 {
		e.SetRegister(dest, e.GetRegister(dest)-uint16(optionA))
	}
}

// operationMovez performs the MOVEZ instruction. Copies the value from the register in optionA to the destination register if the register in optionB is zero
func operationMovez(e *Emulator, dest, optionA, optionB byte) {
	if e.GetRegister(optionB) == 0 {
		e.SetRegister(dest, e.GetRegister(optionA))
	}
}

// operationMovex performs the MOVEX instruction. Copies the value from the register in optionA to the destination register if the register in optionB is not zero
func operationMovex(e *Emulator, dest, optionA, optionB byte) {
	if e.GetRegister(optionB) != 0 {
		e.SetRegister(dest, e.GetRegister(optionA))
	}
}

// operationMovep performs the MOVEP instruction. Copies the value from the register in optionA to the destination register if the register in optionB is positive (most significant bit is 0)
func operationMovep(e *Emulator, dest, optionA, optionB byte) {
	if e.GetRegister(optionB)&0x8000 == 0 {
		e.SetRegister(dest, e.GetRegister(optionA))
	}
}

// operationMoven performs the MOVEN instruction. Copies the value from the register in optionA to the destination register if the register in optionB is negative (most significant bit is 1)
func operationMoven(e *Emulator, dest, optionA, optionB byte) {
	if e.GetRegister(optionB)&0x8000 != 0 {
		e.SetRegister(dest, e.GetRegister(optionA))
	}
}
