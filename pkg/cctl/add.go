package cctl

import (
	"github.com/spf13/cobra"
)

func NewCmdAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Add specific object to config",
	}

	cmd.AddCommand(NewCmdUser())
	cmd.AddCommand(NewCmdNetwork())

	return cmd
}
