package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var config_file string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "carousel",
	Short: "A modern IRC bouncer written in Go",
	Long:  "A modern IRC bouncer written in Go",
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
