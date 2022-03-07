package options

import "github.com/spf13/cobra"

// OptionAdder interface for an options that need to add flags to a command
type OptionAdder interface {
	AddFlags(cmd *cobra.Command) error
}
