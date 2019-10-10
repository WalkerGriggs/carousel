package cmd

import (
	"log"
	"sync"

	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/crypto/ssl"
)

var certificateCmd = &cobra.Command{
	Use:   "certificate",
	Short: "Generate a new SSL certificate.",
	Long:  "Generate a new SSL certificate.",
	Run: func(cmd *cobra.Command, args []string) {
		if survey_confirm("This will overwrite any existing SSL certificate. Continue?") {
			config_dir, err := config_dir()
			if err != nil {
				log.Fatal(err)
			}

			var wg sync.WaitGroup
			wg.Add(1)

			go func(certificatePath string) {
				ssl.NewPem(certificatePath)
				wg.Done()
			}(config_dir + "carousel.pem")

			wg.Wait()
		}
	},
}

func init() {
	createCmd.AddCommand(certificateCmd)
}
