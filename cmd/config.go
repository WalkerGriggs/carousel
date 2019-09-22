package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/carousel"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a config",
	Long:  "Create a config",
	Run: func(cmd *cobra.Command, args []string) {
		uri, err := parseURI(cmd, args)
		if err != nil {
			log.Fatal(err)
		}

		s := carousel.Server{
			URI:   *uri,
			Users: []*carousel.User{},
		}

		if err := writeConfig(cmd, s); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	createCmd.AddCommand(configCmd)

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	home += "/.carousel/"

	configCmd.Flags().StringP("address", "a", "0.0.0.0", "Carousel's host address")
	configCmd.Flags().IntP("port", "p", 6667, "Carousel's host port")
	configCmd.Flags().StringP("dir", "d", home, "The directory in which to create the config")
}

// writeConfig writes the Server struct to disk with human readable spacing.
func writeConfig(cmd *cobra.Command, server carousel.Server) error {
	js, err := json.MarshalIndent(server, "", "    ")
	if err != nil {
		return err
	}

	dir, _ := cmd.Flags().GetString("dir")
	return ioutil.WriteFile(dir+"config.json", js, 0644)
}

// parseURI gets each flag's value from the command, validates the results where
// necessary, and returns the formatted URI.
func parseURI(cmd *cobra.Command, args []string) (*carousel.URI, error) {
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")

	if err := validateURI(address, port); err != nil {
		return nil, err
	}

	return &carousel.URI{
		Address: address,
		Port:    port,
	}, nil
}

// validateURI checks that the host port is within the valid range.
func validateURI(address string, port int) error {
	if port < 1 || port > 65535 {
		return errors.New("Port invalid. Port must be between 1 and 65535")
	}

	return nil
}
