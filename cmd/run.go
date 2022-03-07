/*
Copyright © 2022 github.com/littlehawk93
*/
package cmd

import (
	"log"

	"github.com/littlehawk93/csc364emu/cmd/options"
	"github.com/littlehawk93/csc364emu/emu"
	"github.com/littlehawk93/csc364emu/emu/gui"
	"github.com/spf13/cobra"
)

var runOptions *options.RunOptions = &options.RunOptions{}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the CSC 364 Emulator using a pre-made compiled program",
	Run:   executeRunCommand,
}

func init() {
	rootCmd.AddCommand(runCmd)

	if err := runOptions.AddFlags(runCmd); err != nil {
		log.Fatal(err)
	}
}

func executeRunCommand(cmd *cobra.Command, args []string) {

	emulator := emu.New()

	if err := emulator.LoadProgramFile(runOptions.InputFile); err != nil {
		log.Fatal(err)
	}

	emuGui, err := gui.NewGui(emulator, runOptions.GetClampedSpeed())

	if err != nil {
		log.Fatal(err)
	}

	defer emuGui.Close()

	emuGui.MainLoop()
}
