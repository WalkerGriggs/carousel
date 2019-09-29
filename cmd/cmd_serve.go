package cmd

import (
	_ "fmt"
	"log"

	"encoding/json"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve carousel",
	Long:  "Serve carousel",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve() {
	var c server.Server

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}

	js, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	c.Serve()
}
