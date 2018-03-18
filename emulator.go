package main

import (
	"fmt"
	"ihex"
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

type emulatorOperation func(regDest, regA, regB byte, emu *Emulator) bool

// UpdateCallback - Used as a callback function for every time the emulator executes a clock cycle
type UpdateCallback func(emu *Emulator, err error)

// Emulator - A emulator of the LA Tech CSC 364 16 bit microcontroller
type Emulator struct {
	Registers []uint16
	Screen    []byte
	ROM       []uint16
	RAM       []byte
	Clock     uint64
	TickTime  uint64
}

var operationsLookup = map[byte]emulatorOperation{
	opMove:  move,
	opNot:   not,
	opAnd:   and,
	opOr:    or,
	opAdd:   add,
	opSub:   sub,
	opAddi:  addi,
	opSubi:  subi,
	opSet:   set,
	opSeth:  seth,
	opInciz: inciz,
	opDecin: decin,
	opMovez: movez,
	opMovex: movex,
	opMovep: movep,
	opMoven: moven,
}

// NewEmulator - Create an instantiate a new Emulator
func NewEmulator() *Emulator {

	var tmp Emulator

	tmp.Clock = 0
	tmp.TickTime = 1000
	tmp.Screen = make([]byte, hwScreenWidth)
	tmp.Registers = make([]uint16, hwRegCount)
	tmp.RAM = make([]byte, hwRAMSize)
	tmp.ROM = make([]uint16, hwROMSize)

	return &tmp
}

// LoadROM - Populate the ROM of this emulator using a I8HEX ROM File
func (me *Emulator) LoadROM(file *ihex.I8HEX) error {

	return nil
}

// Emulate - Begin running this emulator
func (me *Emulator) Emulate(callback UpdateCallback) {

	for me.Registers[regPC] < hwROMSize {

		if me.Registers[regOut1]&0x8000 == 0 {
			me.setInputRegister()
		}
		pcModified, err := me.executeInstruction(me.ROM[me.Registers[regPC]])

		if err == nil {

			if me.Registers[regOut1]&0x8000 != 0 {
				me.outputRegisterValues()
			}

			if !pcModified {
				me.Registers[regPC]++
			}
		}

		me.Clock++

		callback(me, err)

		if err != nil {
			return
		}

		if me.TickTime > 0 {
			time.Sleep(time.Duration(me.TickTime) * time.Millisecond)
		}
	}
}

func (me *Emulator) setInputRegister() {

	me.Registers[regIn] &= 0xFF00

	if me.Registers[regOut1]&0x4000 == 0 {
		me.Registers[regIn] |= uint16(me.RAM[me.Registers[regOut2]]) & 0x00FF
	} else {
		me.Registers[regIn] |= uint16(me.Screen[0x000F&me.Registers[regOut2]]) & 0x00FF
	}
}

func (me *Emulator) outputRegisterValues() {

	if me.Registers[regOut1]&0x4000 == 0 {
		me.RAM[me.Registers[regOut2]] = byte(0x00FF & me.Registers[regOut1])
	} else {
		me.Screen[0x00FF&me.Registers[regOut2]] = byte(0x00FF & me.Registers[regOut1])
	}
}

func (me *Emulator) executeInstruction(in uint16) (pcModified bool, err error) {

	opCode, regDest, regA, regB := me.processInstruction(in)

	if regDest >= hwRegCount || regDest < 0 {
		return false, fmt.Errorf("DESTINATION REGISTER INVALID VALUE: %d", regDest)
	}

	if regA >= hwRegCount || regA < 0 {
		return false, fmt.Errorf("REGISTER A INVALID VALUE: %d", regA)
	}

	if regB >= hwRegCount || regB < 0 {
		return false, fmt.Errorf("REGISTER B INVALID VALUE: %d", regB)
	}

	pcModified = regDest == regPC

	op, ok := operationsLookup[opCode]

	if !ok {
		return false, fmt.Errorf("Unrecognized Op Code: %d", opCode)
	}

	pcModified = pcModified && op(regDest, regA, regB, me)

	return pcModified, nil
}

func (me *Emulator) processInstruction(in uint16) (opCode, regDest, regA, regB byte) {

	opCode = byte(((in & 0xF000) >> 12) & 0xFF)

	regDest = byte(((in & 0x0F00) >> 8) & 0xFF)

	regA = byte(((in & 0x00F0) >> 4) & 0xFF)

	regB = byte(in & 0x000F)

	return
}

// move - Set Destination Register to Register A
func move(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = emu.Registers[regA]
	return true
}

// not - Set Destination Register to Binary NOT of Register A
func not(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = ^emu.Registers[regA]
	return true
}

// and - Set Destination Register to Binary AND of Register A and B
func and(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = emu.Registers[regA] & emu.Registers[regB]
	return true
}

// or - Set Destination Register to Binary OR of Register A and B
func or(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = emu.Registers[regA] | emu.Registers[regB]
	return true
}

// add - Set Destination Register to Register A plus Register B
func add(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = emu.Registers[regA] + emu.Registers[regB]
	return true
}

// sub - Set Destination Register to Register A minus Register B
func sub(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = emu.Registers[regA] - emu.Registers[regB]
	return true
}

// Set Destination Register to Register A plus Register B address value
func addi(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = emu.Registers[regA] + uint16(regB)
	return true
}

// Set Destination Register to Register A minus Register B address value
func subi(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = emu.Registers[regA] - uint16(regB)
	return true
}

// Set the lower 8 bits of Destination Register to the address values of Register A (upper 4 bits) and Register B (lower 4 bits). Clears upper 8 bits of Destination Register
func set(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] = ((0x00F0 & (uint16(regA) << 4)) | (0x000F & uint16(regB))) & 0x00FF
	return true
}

// Set the upper 8 bits of Destination Register to the address values of Register A (upper 4 bits) and Register B (lower 4 bits). Does not clear lower 8 bits of Destination Register
func seth(regDest, regA, regB byte, emu *Emulator) bool {

	emu.Registers[regDest] &= 0x00FF
	emu.Registers[regDest] |= (((0x00F0 & (uint16(regA) << 4)) | (0x000F & uint16(regB))) << 8) & 0xFF00
	return true
}

// Increment the Destination Register by the address value of Register A if Register B equals 0
func inciz(regDest, regA, regB byte, emu *Emulator) bool {

	if emu.Registers[regB] == 0 {
		emu.Registers[regDest] += uint16(regA)
		return true
	}

	return false
}

// Decrement the Destination Register by the address value of Register A if Register B is negative (most significant bit is one)
func decin(regDest, regA, regB byte, emu *Emulator) bool {

	if emu.Registers[regB]&0x8000 != 0 {
		emu.Registers[regDest] -= uint16(regA)
		return true
	}

	return false
}

// Set Destination Register to Register A if Register B is zero
func movez(regDest, regA, regB byte, emu *Emulator) bool {

	if emu.Registers[regB] == 0 {
		emu.Registers[regDest] = emu.Registers[regA]
		return true
	}

	return false
}

// Set Destination Register to Register A if Register B is not zero
func movex(regDest, regA, regB byte, emu *Emulator) bool {

	if emu.Registers[regB] != 0 {
		emu.Registers[regDest] = emu.Registers[regA]
		return true
	}

	return false
}

// Set Destination Register to Register A if Register B is positive (most significant bit is zero)
func movep(regDest, regA, regB byte, emu *Emulator) bool {

	if emu.Registers[regB]&0x8000 == 0 {
		emu.Registers[regDest] = emu.Registers[regA]
		return true
	}

	return false
}

// Set Destination Register to Register A if Register B is positive (most significant bit is one)
func moven(regDest, regA, regB byte, emu *Emulator) bool {

	if emu.Registers[regB]&0x8000 != 0 {
		emu.Registers[regDest] = emu.Registers[regA]
		return true
	}

	return false
}
