package cmd

import (
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/cmd/config"
	"github.com/walkergriggs/carousel/pkg/crypto/ssl"
)

type CmdKeyOptions struct {
	Dir  string
	Name string
}

func NewCmdKey(configAccess config.ConfigAccess) *cobra.Command {
	o := &CmdKeyOptions{}

	cmd := &cobra.Command{
		Use:                   "key [--name] [--path]",
		Short:                 "Generates a new encryption KEY",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(configAccess)
			o.Run(configAccess)
		},
	}

	cmd.Flags().StringVarP(&o.Name, "name", "n", "carousel", "Certificate name")
	cmd.Flags().StringVarP(&o.Dir, "dir", "d", "$HOME/.carousel/", "Certificate directory")

	return cmd
}

func (o *CmdKeyOptions) Complete(configAccess config.ConfigAccess) {
	if o.Name == "" {
		o.Name = "carousel"
	}

	if o.Dir == "" {
		o.Dir = filepath.Dir(configAccess.GetDefaultFilename())
	}

	abs, err := filepath.Abs(o.Dir)
	if err != nil {
		log.Fatal(err)
	}

	o.Dir = abs
}

func (o *CmdKeyOptions) Run(configAccess config.ConfigAccess) {
	startingConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		log.Fatal(err)
	}

	path := path.Join(o.Dir, o.Name+".pem")

	err = ssl.NewPem(path)
	if err != nil {
		log.Fatal(err)
	}

	startingConfig.Server.CertificatePath = path
	config.ModifyFile(configAccess, *startingConfig)
}
