package add

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/cmd/config"
	"github.com/walkergriggs/carousel/pkg/server"
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

func unmarshalConfig() (*server.Server, error) {
	var s server.Server

	err := viper.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
