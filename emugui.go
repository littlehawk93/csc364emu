package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var gui *gocui.Gui

// CloseGui - Close the terminal GUI
func CloseGui() {

	gui.Close()
}

// InitGui - Initialize the terminal GUI
func InitGui() error {

	gui, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		return err
	}

	gui.Mouse = false

	gui.SetManagerFunc(LayoutGui)

	if err = gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {

		CloseGui()
		return err
	}

	return nil
}

// BeginMainLoop - Begins the terminal GUI main event listening loop
func BeginMainLoop() {
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

// LayoutGui - Initialize and generate main layout for terminal GUI
func LayoutGui(g *gocui.Gui) error {

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

	return nil
}

// RenderEmulator - Render an Emulator's properties on the terminal GUI
func RenderEmulator(g *gocui.Gui, emu *Emulator) error {

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
