package cmd

import (
	"github.com/spf13/cobra"
	"github.com/walkergriggs/carousel/cmd/config"
)

func NewCmdGenerate(configAccess config.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "generate SUBCOMMAND",
		Short:                 "Generate new objects",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewCmdKey(configAccess))
	cmd.AddCommand(NewCmdGenerateConfig(configAccess))

	return cmd
}
