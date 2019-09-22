package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	tui "github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

var config_file string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "carousel",
	Short: "A modern IRC bouncer written in Go",
	Long:  "A modern IRC bouncer written in Go",
	Run: func(cmd *cobra.Command, args []string) {
		result := menu()

		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		switch result {
		case "Generate Configuration":
			generateConfig(home)
		default:
			os.Exit(0)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&config_file, "config", "", "config file (default is $HOME/.carousel/config.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if config_file != "" {
		viper.SetConfigFile(config_file)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath(home + "/.carousel/")
		viper.AddConfigPath("../")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func menu() string {
	prompt := tui.Select{
		Label: "How can we help you?",
		Items: []string{
			"Generate Configuration",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// generateConfig leverages a number of pipelines to fill in an entire, basic
// configuration. The configuration is then written to the user's home directory
// by default.
// TODO: Allow the user to specify directory.
func generateConfig(home string) {
	server, err := runConfigPipeline()
	if err != nil {
		log.Fatal(err)
	}

	js, err := json.MarshalIndent(server, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(home+"/.carousel/config.json", js, 0644)
}

// runConfigPipeline does the heavy lifting for generateConfig. It runs each
// pipeline, and, if the user wants, sets up a Network for the newly created
// user.
func runConfigPipeline() (*carousel.Server, error) {
	server, err := RunServerPipeline()
	if err != nil {
		return nil, err
	}

	cont, err := continuePrompt("Create a new network for that user?")
	if err != nil {
		return nil, err
	}

	if cont {
		network, err := RunNetworkPipeline()
		if err != nil {
			log.Fatal()
		}

		server.Users[0].Network = network
	}

	return server, nil
}

func continuePrompt(msg string) (bool, error) {
	prompt := tui.Prompt{
		Label:   msg,
		Default: "true",
	}

	res, err := prompt.Run()
	if err != nil {
		return false, err
	}

	b, err := strconv.ParseBool(res)
	if err != nil {
		return false, err
	}

	return b, nil
}
