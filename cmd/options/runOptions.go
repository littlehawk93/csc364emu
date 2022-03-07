package options

import "github.com/spf13/cobra"

// RunOptions program execution parameters for the run command
type RunOptions struct {
	InputFile      string
	ExecutionSpeed int
}

// GetClampedSpeed return execution speed but between 1 - 10
func (me RunOptions) GetClampedSpeed() int {

	if me.ExecutionSpeed < 1 {
		return 1
	} else if me.ExecutionSpeed > 10 {
		return 10
	}
	return me.ExecutionSpeed
}

// AddFlags populate a command's flags with this option's properties
func (me *RunOptions) AddFlags(cmd *cobra.Command) error {

	cmd.Flags().StringVarP(&(me.InputFile), "input", "i", "", "The HEX file to read the compiled program that will be excuted on the emulator")
	cmd.Flags().IntVarP(&(me.ExecutionSpeed), "speed", "s", 5, "The emulator run speed. (0 - 10)")

	if err := cmd.MarkFlagFilename("input", "hex"); err != nil {
		return err
	}
	return cmd.MarkFlagRequired("input")
}
