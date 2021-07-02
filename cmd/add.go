package cmd

import (
	"github.com/spf13/cobra"
	"github.com/walkergriggs/carousel/cmd/config"
)

func NewCmdAdd(configAccess config.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Add specific object to config",
	}

	cmd.AddCommand(NewCmdUser(configAccess))
	cmd.AddCommand(NewCmdNetwork(configAccess))

	return cmd
}
