package serve

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/cmd/config"
	"github.com/walkergriggs/carousel/pkg/server"
	"github.com/walkergriggs/carousel/pkg/uri"
)

type CmdServeOptions struct {
	Host    string
	Port    int
	SSL     bool
	Verbose bool
}

func NewCmdServe(configAccess config.ConfigAccess) *cobra.Command {
	o := &CmdServeOptions{}

	cmd := &cobra.Command{
		Use:   "serve [--host] [--port] [--ssl]",
		Short: "Serves carousel",
		Long:  "Serves carousel",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(configAccess)
		},
	}

	cmd.Flags().StringVarP(&o.Host, "host", "s", "0.0.0.0", "Server host")
	cmd.Flags().IntVarP(&o.Port, "port", "p", 6667, "Server port")
	cmd.Flags().BoolVarP(&o.SSL, "ssl", "", false, "Enables SSL")
	cmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false, "Logs debug messages")

	return cmd
}

func (o *CmdServeOptions) Run(configAccess config.ConfigAccess) {
	if o.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	startingConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		log.Fatal(err)
	}

	opts := server.Options{
		Users:           startingConfig.Users,
		CertificatePath: startingConfig.CertificatePath,
		SSLEnabled:      o.SSL,
		URI: uri.URI{
			Host: o.Host,
			Port: o.Port,
		},
	}

	s, err := server.New(opts)
	if err != nil {
		log.Fatal(err)
	}

	s.Serve()
}
