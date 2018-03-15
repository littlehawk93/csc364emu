package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	hwROMSize     = 65535
	hwRAMSize     = 65536
	hwRegCount    = 16
	hwScreenWidth = 16

	regPC   = 15
	regIn   = 6
	regOut1 = 13
	regOut2 = 14

	opMove  = 0
	opNot   = 1
	opAnd   = 2
	opOr    = 3
	opAdd   = 4
	opSub   = 5
	opAddi  = 6
	opSubi  = 7
	opSet   = 8
	opSeth  = 9
	opInciz = 10
	opDecin = 11
	opMovez = 12
	opMovex = 13
	opMovep = 14
	opMoven = 15
)

type emulator struct {
	Registers []uint16
	Screen    []byte
	ROM       []uint16
	RAM       []byte
	Clock     uint64
	TickTime  uint64
}

func main() {

	tickTime := flag.Uint64("t", 1000, "Sleep time between emulator clock cycles")

	romFile := flag.String("f", "", "Input file to initialize emulator ROM")

	flag.Parse()

	if romFile == nil || *romFile == "" {
		panic("No ROM file provided")
	} else if _, err := os.Stat(*romFile); err != nil {
		panic(fmt.Sprintf("Invalid ROM file provided: %s", err.Error()))
	}

	emu := newEmulator()

	if tickTime != nil && *tickTime >= 0 {
		emu.TickTime = *tickTime
	}
}

func newEmulator() *emulator {

	var emu emulator

	emu.Clock = 0
	emu.TickTime = 1000
	emu.Screen = make([]byte, hwScreenWidth)
	emu.Registers = make([]uint16, hwRegCount)
	emu.RAM = make([]byte, hwRAMSize)
	emu.ROM = make([]uint16, hwROMSize)

	return &emu
}

func (me *emulator) emulate() {

	for me.Registers[regPC] < hwROMSize {

		if me.Registers[regOut1]&0x8000 == 0 {
			me.setInputRegister()
		}
		pcModified := me.executeInstruction(me.ROM[me.Registers[regPC]])

		if me.Registers[regOut1]&0x8000 != 0 {
			me.outputRegisterValues()
		}

		if !pcModified {
			me.Registers[regPC]++
		}

		me.Clock++

		if me.TickTime > 0 {
			time.Sleep(time.Duration(me.TickTime) * time.Millisecond)
		}
	}
}

func (me *emulator) setInputRegister() {

	me.Registers[regIn] &= 0xFF00

	if me.Registers[regOut1]&0x4000 == 0 {
		me.Registers[regIn] |= uint16(me.RAM[me.Registers[regOut2]]) & 0x00FF
	} else {
		me.Registers[regIn] |= uint16(me.Screen[0x000F&me.Registers[regOut2]]) & 0x00FF
	}
}

func (me *emulator) outputRegisterValues() {

	if me.Registers[regOut1]&0x4000 == 0 {
		me.RAM[me.Registers[regOut2]] = byte(0x00FF & me.Registers[regOut1])
	} else {
		me.Screen[0x00FF&me.Registers[regOut2]] = byte(0x00FF & me.Registers[regOut1])
	}
}

func (me *emulator) executeInstruction(in uint16) (pcModified bool) {

	opCode, regDest, regA, regB := processInstruction(in)

	if regDest >= hwRegCount || regDest < 0 {
		panic(fmt.Sprintf("DESTINATION REGISTER INVALID VALUE: %d", regDest))
	}

	if regA >= hwRegCount || regA < 0 {
		panic(fmt.Sprintf("REGISTER A INVALID VALUE: %d", regA))
	}

	if regB >= hwRegCount || regB < 0 {
		panic(fmt.Sprintf("REGISTER B INVALID VALUE: %d", regB))
	}

	pcModified = regDest == regPC

	switch opCode {

	// Set Destination Register to Register A
	case opMove:
		me.Registers[regDest] = me.Registers[regA]
		break

	// Set Destination Register to Binary NOT of Register A
	case opNot:
		me.Registers[regDest] = ^me.Registers[regA]
		break

	// Set Destination Register to Binary AND of Register A and B
	case opAnd:
		me.Registers[regDest] = me.Registers[regA] & me.Registers[regB]
		break

	// Set Destination Register to Binary OR of Register A and B
	case opOr:
		me.Registers[regDest] = me.Registers[regA] | me.Registers[regB]
		break

	// Set Destination Register to Register A plus Register B
	case opAdd:
		me.Registers[regDest] = me.Registers[regA] + me.Registers[regB]
		break

	// Set Destination Register to Register A minus Register B
	case opSub:
		me.Registers[regDest] = me.Registers[regA] - me.Registers[regB]
		break

	// Set Destination Register to Register A plus Register B address value
	case opAddi:
		me.Registers[regDest] = me.Registers[regA] + uint16(regB)
		break

	// Set Destination Register to Register A minus Register B address value
	case opSubi:
		me.Registers[regDest] = me.Registers[regA] - uint16(regB)
		break

	// Set the lower 8 bits of Destination Register to the address values of Register A (upper 4 bits) and Register B (lower 4 bits). Clears upper 8 bits of Destination Register
	case opSet:
		me.Registers[regDest] = ((0x00F0 & (uint16(regA) << 4)) | (0x000F & uint16(regB))) & 0x00FF
		break

	// Set the upper 8 bits of Destination Register to the address values of Register A (upper 4 bits) and Register B (lower 4 bits). Does not clear lower 8 bits of Destination Register
	case opSeth:
		me.Registers[regDest] &= 0x00FF
		me.Registers[regDest] |= (((0x00F0 & (uint16(regA) << 4)) | (0x000F & uint16(regB))) << 8) & 0xFF00
		break

	// Increment the Destination Register by the address value of Register A if Register B equals 0
	case opInciz:
		if me.Registers[regB] == 0 {
			me.Registers[regDest] += uint16(regA)
		} else if pcModified {
			pcModified = false
		}
		break

	// Decrement the Destination Register by the address value of Register A if Register B is negative (most significant bit is one)
	case opDecin:
		if me.Registers[regB]&0x8000 != 0 {
			me.Registers[regDest] -= uint16(regA)
		} else if pcModified {
			pcModified = false
		}
		break

	// Set Destination Register to Register A if Register B is zero
	case opMovez:
		if me.Registers[regB] == 0 {
			me.Registers[regDest] = me.Registers[regA]
		} else if pcModified {
			pcModified = false
		}
		break

	// Set Destination Register to Register A if Register B is not zero
	case opMovex:
		if me.Registers[regB] != 0 {
			me.Registers[regDest] = me.Registers[regA]
		} else if pcModified {
			pcModified = false
		}
		break

	// Set Destination Register to Register A if Register B is positive (most significant bit is zero)
	case opMovep:
		if me.Registers[regB]&0x8000 == 0 {
			me.Registers[regDest] = me.Registers[regA]
		} else if pcModified {
			pcModified = false
		}
		break

	// Set Destination Register to Register A if Register B is positive (most significant bit is one)
	case opMoven:
		if me.Registers[regB]&0x8000 != 0 {
			me.Registers[regDest] = me.Registers[regA]
		} else if pcModified {
			pcModified = false
		}
		break

	default:
		panic(fmt.Sprintf("UNEXPECTED OPCODE: %d", opCode))
	}

	return
}

func processInstruction(in uint16) (opCode, regDest, regA, regB byte) {

	opCode = byte(((in & 0xF000) >> 12) & 0xFF)

	regDest = byte(((in & 0x0F00) >> 8) & 0xFF)

	regA = byte(((in & 0x00F0) >> 4) & 0xFF)

	regB = byte(in & 0x000F)

	return
}
