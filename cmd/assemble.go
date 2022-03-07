package cmd

import (
	"log"

	"github.com/littlehawk93/csc364emu/asm"
	"github.com/littlehawk93/csc364emu/cmd/options"
	"github.com/littlehawk93/ihex"
	"github.com/spf13/cobra"
)

var assembleOptions *options.AssembleOptions = &options.AssembleOptions{}

// assembleCmd represents the assemble command
var assembleCmd = &cobra.Command{
	Use:   "assemble",
	Short: "Convert CSC 364 plaintext assembly code into a binary HEX file that can be executed by the CSC 364 emulator",
	Run:   executeAssembleCommand,
}

func init() {
	rootCmd.AddCommand(assembleCmd)

	if err := assembleOptions.AddFlags(assembleCmd); err != nil {
		log.Fatal(err)
	}
}

func executeAssembleCommand(cmd *cobra.Command, args []string) {

	in, err := assembleOptions.GetInput()

	if err != nil {
		log.Fatal(err)
	}

	defer in.Close()

	out, err := assembleOptions.GetOutput()

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	hexWriter, err := ihex.NewFileWriterType(out, asm.InstructionSize, ihex.I8HEX)

	if err != nil {
		log.Fatal(err)
	}

	defer hexWriter.Close()

	if err = asm.Assemble(in, hexWriter); err != nil {
		log.Fatal(err)
	}
}
