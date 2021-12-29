package cadm

import (
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cadm",
	Short: "A modern IRC bouncer written in Go",
	Long:  "A modern IRC bouncer written in Go",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(NewCmdServe())
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(home + "/.carousel/")
	viper.SetConfigName("config")
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file: ", viper.ConfigFileUsed())
	}
}

func UnmarshalConfig() (*carousel.Server, error) {
	var s carousel.Server

	err := viper.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
