package main

import (
	"flag"
	"fmt"
	"ihex"
	"os"
)

func main() {

	tickTime := flag.Uint64("t", 1000, "Sleep time between emulator clock cycles")

	romFile := flag.String("f", "", "Input file to initialize emulator ROM")

	flag.Parse()

	if romFile == nil || *romFile == "" {
		panic("No ROM file provided")
	} else if _, err := os.Stat(*romFile); err != nil {
		panic(fmt.Sprintf("Invalid ROM file provided: %s", err.Error()))
	}

	file, err := os.Open(*romFile)

	if err != nil {
		panic(fmt.Sprintf("Unable to read ROM file: %s", err.Error()))
	}

	hexFile, err := ihex.NewI8HEX(file)

	if err != nil {
		panic(fmt.Sprintf("Unable to read I8HEX ROM file: %s", err.Error()))
	}

	emu := NewEmulator()

	err = emu.LoadROM(hexFile)

	if err != nil {
		panic(fmt.Sprintf("Unable to load I8HEX ROM file: %s", err.Error()))
	}

	gui, err := NewGui(emu)

	if err != nil {
		panic(err)
	}

	defer gui.Close()

	if tickTime != nil && *tickTime >= 0 {
		emu.TickTime = *tickTime
	}

	gui.MainLoop()
}
