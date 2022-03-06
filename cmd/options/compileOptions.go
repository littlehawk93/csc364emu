package options

import (
	"os"

	"github.com/spf13/cobra"
)

type CompileOptions struct {
	InputFile  string
	OutputFile string
}

func (me CompileOptions) GetInput() (*os.File, error) {

	if me.InputFile == "" {
		return os.Stdin, nil
	}

	return os.Open(me.InputFile)
}

func (me CompileOptions) GetOuput() (*os.File, error) {

	if me.OutputFile == "" {
		return os.Stdout, nil
	}

	return os.Open(me.OutputFile)
}

func (me *CompileOptions) AddFlags(cmd *cobra.Command) error {

	cmd.Flags().StringVarP(&(me.InputFile), "input", "i", "", "The file to read assembly instructions in to compile to machine code")
	cmd.Flags().StringVarP(&(me.OutputFile), "output", "o", "", "The file to write compiled machine code to")

	if err := cmd.MarkFlagFilename("input"); err != nil {
		return err
	}

	return cmd.MarkFlagFilename("output", "hex")
}
