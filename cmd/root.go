package cmd

import (
	"path"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/cmd/add"
	"github.com/walkergriggs/carousel/cmd/config"
	"github.com/walkergriggs/carousel/cmd/generate"
	"github.com/walkergriggs/carousel/cmd/serve"
	"github.com/walkergriggs/carousel/cmd/set"
	"github.com/walkergriggs/carousel/pkg/server"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "carousel",
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
	globalPath, err := globalPath()
	if err != nil {
		log.Fatal(err)
	}

	pathOptions := &config.PathOptions{
		GlobalPath: globalPath,
	}

	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(add.NewCmdAdd(pathOptions))
	rootCmd.AddCommand(set.NewCmdSet(pathOptions))
	rootCmd.AddCommand(serve.NewCmdServe(pathOptions))
	rootCmd.AddCommand(generate.NewCmdGenerate(pathOptions))
}

func globalPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	path := path.Join(home, ".carousel/config.json")
	return path, nil
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

func UnmarshalConfig() (*server.Server, error) {
	var s server.Server

	err := viper.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
