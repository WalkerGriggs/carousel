package cctl

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/api"
)

type CmdUserOptions struct {
	Meta
	Username string
	Password string
}

func NewCmdUser() *cobra.Command {
	o := &CmdUserOptions{}

	cmd := &cobra.Command{
		Use:                   "user username (--pass)",
		Short:                 "Adds user to server",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.Username, "username", "u", o.Username, "User name")
	cmd.Flags().StringVarP(&o.Password, "password", "p", o.Password, "User password")

	sharedFlags := o.Meta.FlagSet("Add user")
	cmd.Flags().AddFlagSet(sharedFlags)

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

	return cmd
}

func (o *CmdUserOptions) Run() {
	client, err := o.Meta.Client()
	if err != nil {
		fmt.Println(err)
		return
	}

	user := &api.User{
		Username: o.Username,
		Password: o.Password,
	}

	res, err := client.Users().Create(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
