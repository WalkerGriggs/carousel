package cctl

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/api"
)

type CmdNetworkOptions struct {
	Meta
	Port     int
	Address  string
	User     string
	Name     string
	Username string
	Nickname string
	Realname string
	Password string
}

func NewCmdNetwork() *cobra.Command {
	o := &CmdNetworkOptions{}

	cmd := &cobra.Command{
		Use:                   "network name (--user) (--addr) [--port]",
		DisableFlagsInUseLine: true,
		Short:                 "Adds network to user",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	cmd.Flags().IntVar(&o.Port, "port", 6667, "Network port")
	cmd.Flags().StringVar(&o.Address, "address", o.Address, "Network address")
	cmd.Flags().StringVar(&o.User, "user", o.User, "User to add the network to")
	cmd.Flags().StringVar(&o.Name, "name", o.Name, "Network name")
	cmd.Flags().StringVar(&o.Username, "username", o.Username, "Identity username")
	cmd.Flags().StringVar(&o.Nickname, "nickname", o.Nickname, "Identity nickname")
	cmd.Flags().StringVar(&o.Realname, "realname", o.Realname, "Identity realname")
	cmd.Flags().StringVar(&o.Password, "password", o.Password, "Identity password")

	sharedFlags := o.Meta.FlagSet("Add network")
	cmd.Flags().AddFlagSet(sharedFlags)

	cmd.MarkFlagRequired("user")
	cmd.MarkFlagRequired("address")
	cmd.MarkFlagRequired("port")
	cmd.MarkFlagRequired("nickname")

	return cmd
}

func (o *CmdNetworkOptions) Run() {
	client, err := o.Meta.Client()
	if err != nil {
		fmt.Println(err)
		return
	}

	network := &api.Network{
		Name: o.Name,
		URI:  fmt.Sprintf("%s:%d", o.Address, o.Port),
		Ident: &api.Identity{
			Username: o.Username,
			Password: o.Password,
			Realname: o.Realname,
			Nickname: o.Nickname,
		},
		Channels: make([]string, 0),
	}

	res, err := client.Networks().Create(o.User, network)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
