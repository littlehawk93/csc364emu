package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jroimartin/gocui"
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

var emu *emulator
var gui *gocui.Gui

func main() {

	tickTime := flag.Uint64("t", 1000, "Sleep time between emulator clock cycles")

	romFile := flag.String("f", "", "Input file to initialize emulator ROM")

	flag.Parse()

	if romFile == nil || *romFile == "" {
		panic("No ROM file provided")
	} else if _, err := os.Stat(*romFile); err != nil {
		panic(fmt.Sprintf("Invalid ROM file provided: %s", err.Error()))
	}

	emu = newEmulator()

	gui = initGui()

	defer gui.Close()

	if tickTime != nil && *tickTime >= 0 {
		emu.TickTime = *tickTime
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func initGui() *gocui.Gui {

	tmp, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		panic(err)
	}

	tmp.Mouse = false

	tmp.SetManagerFunc(layout)

	if err = tmp.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {

		tmp.Close()
		panic(err)
	}

	return tmp
}

func layout(g *gocui.Gui) error {

	g.BgColor = gocui.ColorBlack
	g.FgColor = gocui.ColorWhite

	titleView, err := g.SetView("title", 2, 0, 30, 2)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	titleView.BgColor = gocui.ColorBlack
	titleView.FgColor = gocui.ColorWhite
	titleView.Editable = false
	titleView.Frame = false

	titleView.Clear()

	fmt.Fprint(titleView, " LA Tech CSC 364 Emulator ")

	clockView, err := g.SetView("clock", 32, 0, 52, 2)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	clockView.BgColor = gocui.ColorBlack
	clockView.FgColor = gocui.ColorWhite
	clockView.Editable = false
	clockView.Frame = false

	clockView.Clear()

	fmt.Fprintf(clockView, "Clock: %12d", emu.Clock)

	screenView, err := g.SetView("screen", 2, 3, 19, 12)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	screenView.BgColor = gocui.ColorBlack
	screenView.FgColor = gocui.ColorGreen
	screenView.Editable = false
	screenView.Frame = true
	screenView.Title = "Screen"

	screenView.Clear()

	registersView, err := g.SetView("registers", 21, 3, 52, 12)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	registersView.BgColor = gocui.ColorBlack
	registersView.FgColor = gocui.ColorWhite
	registersView.Editable = false
	registersView.Frame = true
	registersView.Title = "Registers"

	registersView.Clear()

	renderEmulator(g)

	return nil
}

func renderEmulator(g *gocui.Gui) error {

	registersView, err := g.View("registers")

	if err != nil {
		return err
	}

	registersView.Clear()

	for i := 0; i < len(emu.Registers)/2; i++ {

		fmt.Fprintf(registersView, " %s  | %04X | %s | %04X \n", getRegisterTitle(i), emu.Registers[i], getRegisterTitle(i+8), emu.Registers[i+8])
	}

	screenView, err := g.View("screen")

	screenView.Clear()

	for x := 0x80; x != 0; x >>= 1 {

		for i := 0; i < len(emu.Screen); i++ {

			if emu.Screen[i]&byte(x) == 0 {
				fmt.Fprintf(screenView, "%c", ' ')
			} else {
				fmt.Fprintf(screenView, "%c", '\u2588')
			}
		}

		fmt.Fprintf(screenView, "%c", '\n')
	}

	return nil
}

func getRegisterTitle(index int) string {

	if index == regIn {
		return "INPUT"
	} else if index == regOut1 {
		return "OUT 1"
	} else if index == regOut2 {
		return "OUT 2"
	} else if index == regPC {
		return "PROGC"
	} else {
		return fmt.Sprintf("REG %1X", index)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func newEmulator() *emulator {

	var tmp emulator

	tmp.Clock = 0
	tmp.TickTime = 1000
	tmp.Screen = make([]byte, hwScreenWidth)
	tmp.Registers = make([]uint16, hwRegCount)
	tmp.RAM = make([]byte, hwRAMSize)
	tmp.ROM = make([]uint16, hwROMSize)

	return &tmp
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

		gui.Update(renderEmulator)

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
