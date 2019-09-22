package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Create a user's network",
	Long:  "Create a user's network",
	Run: func(cmd *cobra.Command, args []string) {
		network, err := parseNetwork(cmd, args)
		if err != nil {
			log.Fatal(err)
		}

		var users []*carousel.User
		user, err := unmarshalUser(cmd, &users)
		if err != nil {
			log.Fatal(err)
		}

		user.Network = network
		viper.Set("users", users)
		viper.WriteConfig()
	},
}

func init() {
	setCmd.AddCommand(networkCmd)

	networkCmd.Flags().StringP("user", "", "", "The username to attach the network to")
	networkCmd.Flags().StringP("name", "n", "", "The name of the network")
	networkCmd.Flags().StringP("address", "a", "", "The network host")
	networkCmd.Flags().IntP("port", "p", 6667, "The network's port")

	networkCmd.MarkFlagRequired("username")
	networkCmd.MarkFlagRequired("name")
	networkCmd.MarkFlagRequired("address")
	networkCmd.MarkFlagRequired("port")
}

func parseNetwork(cmd *cobra.Command, args []string) (*carousel.Network, error) {
	name, _ := cmd.Flags().GetString("name")

	uri, err := parseURI(cmd, args)
	if err != nil {
		return nil, err
	}

	return &carousel.Network{
		Name:  name,
		URI:   *uri,
		Ident: carousel.Identity{},
	}, nil
}
