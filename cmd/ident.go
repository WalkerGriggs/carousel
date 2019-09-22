package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

// identCmd represents the ident command
var identCmd = &cobra.Command{
	Use:   "ident",
	Short: "Set a network's ident",
	Long:  "Set a network's ident",
	Run: func(cmd *cobra.Command, args []string) {
		ident, err := parseIdent(cmd, args)
		if err != nil {
			log.Fatal(err)
		}

		var users []*carousel.User
		user, err := unmarshalUser(cmd, &users)
		if err != nil {
			log.Fatal(err)
		}

		user.Network.Ident = *ident
		viper.Set("users", users)
		viper.WriteConfig()
	},
}

func init() {
	setCmd.AddCommand(identCmd)

	identCmd.Flags().StringP("user", "", "", "The username to attach the network to")
	identCmd.Flags().StringP("nickname", "n", "", "The nickname used to authenticate with the network")
	identCmd.Flags().StringP("username", "u", "", "The username used to...")
	identCmd.Flags().StringP("realname", "r", "", "The realname used to... (Defaults to the username)")
	identCmd.Flags().StringP("password", "p", "", "The username's password used to...")

	identCmd.MarkFlagRequired("user")
	identCmd.MarkFlagRequired("nickname")
	identCmd.MarkFlagRequired("username")
}

func parseIdent(cmd *cobra.Command, args []string) (*carousel.Identity, error) {
	nickname, _ := cmd.Flags().GetString("nickname")
	username, _ := cmd.Flags().GetString("username")
	realname, _ := cmd.Flags().GetString("realname")
	password, _ := cmd.Flags().GetString("password")

	fmt.Println(password)
	if password != "" {
		password, _ = carousel.Hash(password)
	}

	if realname == "" {
		realname = username
	}

	return &carousel.Identity{
		Nickname: nickname,
		Username: username,
		Realname: realname,
		Password: password,
	}, nil
}

func unmarshalUser(cmd *cobra.Command, users *[]*carousel.User) (*carousel.User, error) {
	if err := viper.UnmarshalKey("users", users); err != nil {
		return nil, err
	}

	username, _ := cmd.Flags().GetString("user")
	user := carousel.GetUser(*users, username)
	if user == nil {
		return nil, errors.New("User not found")
	}

	return user, nil
}
