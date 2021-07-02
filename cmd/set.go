package cmd

import (
	"github.com/spf13/cobra"
	"github.com/walkergriggs/carousel/cmd/config"
)

func NewCmdSet(configAccess config.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set SUBCOMMAND",
		Short:                 "Set specific object in config",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewCmdIdent(configAccess))

	return cmd
}
