package options

import (
	"os"

	"github.com/spf13/cobra"
)

// AssembleOptions program execution parameters for the assemble command
type AssembleOptions struct {
	InputFile  string
	OutputFile string
}

// GetInput get the input stream to read incoming program assembly data from
func (me AssembleOptions) GetInput() (*os.File, error) {

	if me.InputFile == "" {
		return os.Stdin, nil
	}

	return os.Open(me.InputFile)
}

// GetOutput get the output stream to write outgoing binary machine code to
func (me AssembleOptions) GetOutput() (*os.File, error) {

	if me.OutputFile == "" {
		return os.Stdout, nil
	}

	return os.Open(me.OutputFile)
}

// AddFlags populate a command's flags with this option's properties
func (me *AssembleOptions) AddFlags(cmd *cobra.Command) error {

	cmd.Flags().StringVarP(&(me.InputFile), "input", "i", "", "The file to read assembly instructions that will be assembled to machine code")
	cmd.Flags().StringVarP(&(me.OutputFile), "output", "o", "", "The file to write assembled machine code to")

	if err := cmd.MarkFlagFilename("input"); err != nil {
		return err
	}

	return cmd.MarkFlagFilename("output", "hex")
}
