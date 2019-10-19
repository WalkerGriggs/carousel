package cmd

import (
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(home + "/.carousel/")
	viper.SetConfigName("config")

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file: ", viper.ConfigFileUsed())
	}
}
