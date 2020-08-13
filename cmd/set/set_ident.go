package set

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/cmd/config"
	"github.com/walkergriggs/carousel/pkg/identity"
)

type CmdIdentOptions struct {
	User     string
	Username string
	Nickname string
	Realname string
	Password string
}

func NewCmdIdent(configAccess config.ConfigAccess) *cobra.Command {
	o := &CmdIdentOptions{}

	cmd := &cobra.Command{
		Use:                   "ident (--user) (--nickname) [--username] [--realname] [--password]",
		DisableFlagsInUseLine: true,
		Short:                 "Sets the identity of a user's network",
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete()
			o.Run(configAccess)
		},
	}

	cmd.Flags().StringVarP(&o.User, "user", "u", o.User, "User of network")
	cmd.Flags().StringVarP(&o.Username, "username", "s", o.Username, "Identity username")
	cmd.Flags().StringVarP(&o.Nickname, "nickname", "n", o.Nickname, "Identity nickname")
	cmd.Flags().StringVarP(&o.Realname, "realname", "r", o.Realname, "Identity realname")
	cmd.Flags().StringVarP(&o.Password, "password", "p", o.Password, "Identity password")

	cmd.MarkFlagRequired("user")
	cmd.MarkFlagRequired("nickname")

	return cmd
}

func (o *CmdIdentOptions) Complete() {
	if o.Username == "" {
		o.Username = o.Nickname
	}

	if o.Realname == "" {
		o.Realname = o.Nickname
	}
}

func (o *CmdIdentOptions) Run(configAccess config.ConfigAccess) {
	startingConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		log.Fatal(err)
	}

	ident := &identity.Identity{
		Username: o.Username,
		Nickname: o.Nickname,
		Realname: o.Realname,
		Password: o.Password,
	}

	for _, user := range startingConfig.Users {
		if user.Username == o.User {
			user.Network.Ident = ident
			config.ModifyFile(configAccess, *startingConfig)
			return
		}
	}

	log.Fatal(fmt.Errorf("User %s not found.\n", o.Username))
}
