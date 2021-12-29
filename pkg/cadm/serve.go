package cadm

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/carousel"
)

type CmdServeOptions struct {
	Host    string
	Port    int
	SSL     bool
	Verbose bool
}

func NewCmdServe() *cobra.Command {
	o := &CmdServeOptions{}

	cmd := &cobra.Command{
		Use:   "serve [--host] [--port] [--ssl]",
		Short: "Serves carousel",
		Long:  "Serves carousel",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.Host, "host", "s", "0.0.0.0", "Server host")
	cmd.Flags().IntVarP(&o.Port, "port", "p", 6667, "Server port")
	cmd.Flags().BoolVarP(&o.SSL, "ssl", "", false, "Enables SSL")
	cmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false, "Logs debug messages")

	return cmd
}

func (o *CmdServeOptions) Run() {
	if o.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := path.Join(dirname, "./carousel/config.json")

	config, err := unmarshalConfig(path)
	if err != nil {
		return
	}

	serverConfig, _ := ConvertServerConfig(config)

	s, err := carousel.NewServer(serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	carousel.NewHTTPServer(s, &carousel.HTTPConfig{
		Advertise: "127.0.0.1:8080",
	})

	s.Serve()
}
