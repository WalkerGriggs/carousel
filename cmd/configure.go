package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/ssl"
)

// menuCmd represents the menu command
var configCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Carousel",
	Long:  "Configure Carousel",
	Run: func(cmd *cobra.Command, args []string) {
		if survey_confirm("This command will overwrite any existing config. Continue?") {
			config_dir, err := config_dir()
			if err != nil {
				log.Fatal(err)
			}

			// Survey the user for their configuration settings
			server := survey_server()
			server.CertificatePath = config_dir + "carousel.pem"

			// Generate a new SSL PEM file in the background
			go ssl.NewPem(server.CertificatePath)

			js, err := json.MarshalIndent(server, "", "    ")
			if err != nil {
				log.Fatal(err)
			}

			if err := ioutil.WriteFile(config_dir+"/config.json", js, 0644); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// config_path gets the absolute path of the configuration file. Defaults to
// '$HOME/carousel/config.json'.
func config_dir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	config_dir := home + "/.carousel/"
	config_file := config_dir + "/config.json"

	// Create .carousel directory if it does not exist
	if stat, err := os.Stat(config_file); err == nil && !stat.IsDir() {
		os.MkdirAll(config_dir, os.ModePerm)
	}

	return config_dir, nil
}
