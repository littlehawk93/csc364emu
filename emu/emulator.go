package emu

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/littlehawk93/ihex"
)

const (
	// HWROMSize maximum number of instructions that can be stored in this emulator's ROM
	HWROMSize = 65535

	// HWRAMSize maximum number of bytes that can be stored in this emulator's RAM
	HWRAMSize = 65535

	// HWRegistersCount the number of registers in this emulator
	HWRegistersCount = 16

	// HWScreenWidth number of columns wide the emulator screen is
	HWScreenWidth = 16

	// HWInstructionSize number of bytes per instruction
	HWInstructionSize = 2

	// RegisterInput address value for the input register
	RegisterInput = 6

	// RegisterOutput1 address value for the output1 register
	RegisterOutput1 = 13

	// RegisterOutput2 address value for the output2 register
	RegisterOutput2 = 14

	// RegisterProgramCounter address value for the program counter
	RegisterProgramCounter = 15
)

// EmulatorUpdate is a callback function that executes with each clock cycle the emulator executes or if any errors occur
type EmulatorUpdate func(emu *Emulator, err error)

// Emulator is the CSC 364 emulator
type Emulator struct {
	registers []uint16
	Screen    []byte
	ROM       []uint16
	RAM       []byte
	clock     uint64

	pcModified bool
}

// GetRegister return the value at a particular register
func (me Emulator) GetRegister(register byte) uint16 {
	return me.registers[register]
}

// SetRegister modify the value within a particular register
func (me *Emulator) SetRegister(register byte, value uint16) {

	if register == RegisterProgramCounter {
		me.pcModified = true
	}
	me.registers[register] = value
}

// GetClock get total number of clock cycles that have been executed by the emulator
func (me Emulator) GetClock() uint64 {
	return me.clock
}

// Begin begin executing this emulator's program
func (me *Emulator) Begin(stepSpeed int, update EmulatorUpdate) error {

	me.pcModified = true

	for me.registers[RegisterProgramCounter] < HWROMSize {
		if stepSpeed < 10 {
			time.Sleep(time.Duration((10-stepSpeed)*100) * time.Millisecond)
		}

		if err := me.step(); err != nil {
			update(me, err)
			return err
		}

		if update != nil {
			update(me, nil)
		}
	}
	return nil
}

func (me *Emulator) step() error {

	if !me.pcModified {
		me.registers[RegisterProgramCounter]++
	} else {
		me.pcModified = false
	}

	instruction := me.ROM[me.registers[RegisterProgramCounter]]
	opCode, dest, optionA, optionB := convertToBytes(instruction)

	if err := executeInstruction(me, opCode, dest, optionA, optionB); err != nil {
		return err
	}

	me.clock++
	return nil
}

func (me *Emulator) setInputRegister() {

	me.registers[RegisterInput] &= 0xFF00

	if me.registers[RegisterOutput1]&0x4000 == 0 {
		me.registers[RegisterInput] |= uint16(me.RAM[me.registers[RegisterOutput2]]) & 0x00FF
	} else {
		me.registers[RegisterInput] |= uint16(me.Screen[me.registers[RegisterOutput2]]) & 0x00FF
	}
}

func (me *Emulator) outputRegisterValues() {

	if me.registers[RegisterOutput1]&0x4000 == 0 {
		me.RAM[me.registers[RegisterOutput2]] = byte(0x00FF & me.registers[RegisterOutput1])
	} else {
		me.Screen[me.registers[RegisterOutput2]] = byte(0x00FF & me.registers[RegisterOutput1])
	}
}

// LoadProgram reads a HEX file and loads it into ROM. Returns any errors
func (me *Emulator) LoadProgram(r ihex.File) error {

	recordNumber := 0

	for record, ok := r.ReadNext(); ok; record, ok = r.ReadNext() {
		recordNumber++
		if record.Type == ihex.RecordData {
			if len(record.Data) != HWInstructionSize {
				return fmt.Errorf("[Error at record %d]: Unexpected record data size: %d (expected %d)", recordNumber, len(record.Data), HWInstructionSize)
			} else if record.AddressOffset >= HWROMSize {
				return fmt.Errorf("[Error at record %d]: Invalid Address: %d", recordNumber, record.AddressOffset)
			}

			instruction := (uint16(record.Data[0]) << 8) | (uint16(record.Data[0]) & 0x00FF)

			me.ROM[record.AddressOffset] = instruction
		} else if record.Type == ihex.RecordEOF {
			break
		}
	}
	return nil
}

// LoadProgramReader reads the contents of a HEX data stream and load it into ROM. Returns any errors
func (me *Emulator) LoadProgramReader(r io.Reader) error {

	h, err := ihex.NewFile(r)

	if err != nil {
		return err
	}

	return me.LoadProgram(h)
}

// LoadProgramFile reads the contents of a HEX file and load it into ROM. Returns any errors
func (me *Emulator) LoadProgramFile(file string) error {

	f, err := os.Open(file)

	if err != nil {
		return err
	}

	defer f.Close()

	return me.LoadProgramReader(f)
}

// New create and instantiate a new Emulator
func New() *Emulator {
	return &Emulator{
		registers: make([]uint16, HWRegistersCount),
		ROM:       make([]uint16, HWROMSize),
		RAM:       make([]byte, HWRAMSize),
		Screen:    make([]byte, HWScreenWidth),
		clock:     0,
	}
}
