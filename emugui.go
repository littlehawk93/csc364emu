package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// EmulatorGUI - Struct for managing displaying emulator properties on a Terminal GUI
type EmulatorGUI struct {
	emulatorStarted bool
	emulator        *Emulator
	gui             *gocui.Gui
	guiConstructed  bool
}

// Close - Close the terminal GUI
func (me *EmulatorGUI) Close() {

	me.gui.Close()
}

// NewGui - Initialize the terminal GUI
func NewGui(emu *Emulator) (*EmulatorGUI, error) {

	var emuGui EmulatorGUI

	emuGui.guiConstructed = false

	emuGui.emulator = emu

	gui, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		return nil, err
	}

	gui.Mouse = false

	gui.SetManagerFunc(emuGui.renderEmulator)

	if err = gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {

		emuGui.Close()
		return nil, err
	}

	if err = gui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, emuGui.startEmulator); err != nil {

		emuGui.Close()
		return nil, err
	}

	emuGui.gui = gui

	return &emuGui, nil
}

// MainLoop - Begins the terminal GUI main event listening loop
func (me *EmulatorGUI) MainLoop() {

	if err := me.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func (me *EmulatorGUI) layoutGui(g *gocui.Gui) error {

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

	clockView, err := g.SetView("clock", 49, 0, 69, 2)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	clockView.BgColor = gocui.ColorBlack
	clockView.FgColor = gocui.ColorWhite
	clockView.Editable = false
	clockView.Frame = false

	clockView.Clear()

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

	instructionView, err := g.SetView("instruction", 2, 14, 22, 16)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	instructionView.BgColor = gocui.ColorBlack
	instructionView.FgColor = gocui.ColorWhite
	instructionView.Editable = false
	instructionView.Frame = true
	instructionView.Title = "Instruction"

	instructionView.Clear()

	romView, err := g.SetView("rom", 54, 3, 69, 16)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	romView.BgColor = gocui.ColorBlack
	romView.FgColor = gocui.ColorWhite
	romView.Editable = false
	romView.Frame = true
	romView.Title = "ROM"

	romView.Clear()

	return nil
}

func (me *EmulatorGUI) renderEmulator(g *gocui.Gui) error {

	if !me.guiConstructed {

		err := me.layoutGui(g)

		if err != nil {
			return err
		}

		me.guiConstructed = true
	}

	clockView, err := g.View("clock")

	if err != nil {
		return err
	}

	clockView.Clear()

	fmt.Fprintf(clockView, "Clock: %12d", me.emulator.Clock)

	registersView, err := g.View("registers")

	if err != nil {
		return err
	}

	registersView.Clear()

	for i := 0; i < len(me.emulator.Registers)/2; i++ {

		fmt.Fprintf(registersView, " %s  | %04X | %s | %04X \n", getRegisterTitle(i), me.emulator.Registers[i], getRegisterTitle(i+8), me.emulator.Registers[i+8])
	}

	screenView, err := g.View("screen")

	screenView.Clear()

	for x := 0x80; x != 0; x >>= 1 {

		for i := 0; i < len(me.emulator.Screen); i++ {

			if me.emulator.Screen[i]&byte(x) == 0 {
				fmt.Fprintf(screenView, "%c", ' ')
			} else {
				fmt.Fprintf(screenView, "%c", '\u2588')
			}
		}

		fmt.Fprintf(screenView, "%c", '\n')
	}

	instructionView, err := g.View("instruction")

	if err != nil {
		return nil
	}

	instructionView.Clear()

	if me.emulator.Registers[15] < 0xFFFF {
		fmt.Fprintf(instructionView, "   %1X   %1X   %1X   %1X   ", byte(me.emulator.ROM[me.emulator.Registers[15]]>>12&0x000F), byte(me.emulator.ROM[me.emulator.Registers[15]]>>8&0x000F), byte(me.emulator.ROM[me.emulator.Registers[15]]>>4&0x000F), byte(me.emulator.ROM[me.emulator.Registers[15]]&0x000F))
	} else {
		fmt.Fprint(instructionView, "   0   0   0   0   ")
	}

	romView, err := g.View("rom")

	if err != nil {
		return nil
	}

	romView.Clear()

	romAddress := int(me.emulator.Registers[15])

	var startAddress int
	var endAddress int

	if romAddress < 6 {
		startAddress = 0
	} else if romAddress > len(me.emulator.ROM)-6 {
		startAddress = len(me.emulator.ROM) - 12
	} else {
		startAddress = romAddress - 6
	}

	endAddress = startAddress + 12

	for i := startAddress; i < endAddress; i++ {

		var colorCode string

		if i == romAddress {
			colorCode = "\u001b[42m"
		} else {
			colorCode = "\u001b[40m"
		}

		fmt.Fprintf(romView, "%s %05d | %04X \n", colorCode, i, me.emulator.ROM[i])
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

func (me *EmulatorGUI) startEmulator(g *gocui.Gui, v *gocui.View) error {

	if !me.emulatorStarted {

		go me.emulator.Emulate(me.updateLayout)
		me.emulatorStarted = true
	}

	return nil
}

func (me *EmulatorGUI) updateLayout(emu *Emulator, err error) {

	if err != nil {

		me.Close()
		panic(err)
	}

	me.gui.Update(me.renderEmulator)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
