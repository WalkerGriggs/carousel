package set

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/cmd/config"
	"github.com/walkergriggs/carousel/pkg/server"
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

func unmarshalConfig() (*server.Server, error) {
	var s server.Server

	err := viper.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
