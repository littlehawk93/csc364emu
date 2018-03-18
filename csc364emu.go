package main

import (
	"flag"
	"fmt"
	"os"
)

var emu *Emulator

func main() {

	tickTime := flag.Uint64("t", 1000, "Sleep time between emulator clock cycles")

	romFile := flag.String("f", "", "Input file to initialize emulator ROM")

	flag.Parse()

	if romFile == nil || *romFile == "" {
		panic("No ROM file provided")
	} else if _, err := os.Stat(*romFile); err != nil {
		panic(fmt.Sprintf("Invalid ROM file provided: %s", err.Error()))
	}

	emu = NewEmulator()

	err := InitGui()

	if err != nil {
		panic(err)
	}

	defer CloseGui()

	if tickTime != nil && *tickTime >= 0 {
		emu.TickTime = *tickTime
	}
}
