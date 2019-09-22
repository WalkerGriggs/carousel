package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve carousel",
	Long:  "Serve carousel",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	var c carousel.Server

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on ", c.URI.Format())

	c.Serve()
}
