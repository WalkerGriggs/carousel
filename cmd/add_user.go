package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/cmd/config"
	"github.com/walkergriggs/carousel/pkg/crypto/phash"
)

type CmdUserOptions struct {
	Username string
	Password string
}

func NewCmdUser(configAccess config.ConfigAccess) *cobra.Command {
	o := &CmdUserOptions{}

	cmd := &cobra.Command{
		Use:                   "user username (--pass)",
		Short:                 "Adds user to server",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(configAccess)
		},
	}

	cmd.Flags().StringVarP(&o.Username, "username", "u", o.Username, "User name")
	cmd.Flags().StringVarP(&o.Password, "password", "p", o.Password, "User password")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

	return cmd
}

func (o *CmdUserOptions) Run(configAccess config.ConfigAccess) {
	startingConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(startingConfig.Users)

	pass, err := phash.Hash(o.Password)
	if err != nil {
		log.Fatal(err)
	}

	u := &config.UserConfig{
		Username: o.Username,
		Password: pass,
	}

	startingConfig.Users = append(startingConfig.Users, u)
	config.ModifyFile(configAccess, *startingConfig)
}
