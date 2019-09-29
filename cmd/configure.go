package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// menuCmd represents the menu command
var configCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Carousel",
	Long:  "Configure Carousel",
	Run: func(cmd *cobra.Command, args []string) {
		if survey_confirm("This command will overwrite any existing config. Continue?") {
			server := survey_server()

			js, err := json.MarshalIndent(server, "", "    ")
			if err != nil {
				log.Fatal(err)
			}

			path, err := config_path()
			if err != nil {
				log.Fatal(err)
			}

			if err := ioutil.WriteFile(path, js, 0644); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// config_path gets the absolute path of the configuration file. If the config
// path is not specified, the default '$HOME/carousel/config.json' is used
// instead.
func config_path() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	// If config file is not specified, use default instead
	if config_file == "" {
		config_file = home + "/.carousel/config.json"
	}

	// Create .carousel directory if it does not exist
	if stat, err := os.Stat(config_file); err == nil && !stat.IsDir() {
		os.MkdirAll(home+"/.carousel", os.ModePerm)
	}

	return config_file, nil
}
