package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Create a user",
	Long:  "Create a user",
	Run: func(cmd *cobra.Command, args []string) {
		var users []*carousel.User
		if err := viper.UnmarshalKey("users", &users); err != nil {
			log.Fatal(err)
		}

		user, err := parseUser(cmd, args)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
		viper.Set("users", users)
		viper.WriteConfig()
	},
}

func init() {
	createCmd.AddCommand(userCmd)

	userCmd.Flags().StringP("username", "u", "", "The user's name for Carousel (not a network)")
	userCmd.Flags().StringP("password", "p", "", "The user's password to be hashed and stored")
	userCmd.MarkFlagRequired("username")
	userCmd.MarkFlagRequired("password")
}

func parseUser(cmd *cobra.Command, args []string) (*carousel.User, error) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	hash, err := carousel.Hash(password)
	if err != nil {
		return nil, err
	}

	return &carousel.User{
		Username: username,
		Password: hash,
	}, nil
}
