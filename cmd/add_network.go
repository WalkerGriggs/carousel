package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/cmd/config"
)

type CmdNetworkOptions struct {
	User    string
	Address string
	Port    int
	Name    string
}

func NewCmdNetwork(configAccess config.ConfigAccess) *cobra.Command {
	o := &CmdNetworkOptions{}

	cmd := &cobra.Command{
		Use:                   "network name (--user) (--addr) [--port]",
		DisableFlagsInUseLine: true,
		Short:                 "Adds network to user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Validate(cmd, args); err != nil {
				panic(err)
			}

			o.Complete(cmd, args)
			o.Run(configAccess)
		},
	}

	cmd.Flags().StringVarP(&o.User, "user", "u", o.User, "User to add the network to")
	cmd.Flags().StringVarP(&o.Address, "address", "a", o.Address, "Network address")
	cmd.Flags().IntVarP(&o.Port, "port", "p", 6667, "Network port")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("address")

	return cmd
}

func (o *CmdNetworkOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("First argument must be the network name")
	}
}

func (o *CmdNetworkOptions) Complete(cmd *cobra.Command, args []string) {
	o.Name = args[0]
}

func (o *CmdNetworkOptions) Run(configAccess config.ConfigAccess) {
	startingConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		log.Fatal(err)
	}

	networkConfig := &config.NetworkConfig{
		Name: o.Name,
		URI:  fmt.Sprintf("%s:%s", o.Address, o.Port),
	}

	for _, user := range startingConfig.Users {
		if user.Username == o.User {
			user.Networks = append(user.Networks, networkConfig)
			config.ModifyFile(configAccess, *startingConfig)
			return
		}
	}

	log.Fatal(fmt.Errorf("User %s not found.\n", o.User))
}
