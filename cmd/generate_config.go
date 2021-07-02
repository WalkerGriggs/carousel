package cmd

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/cmd/config"
)

type GenerateConfigOptions struct {
	Dir string
}

func NewGenerateConfigOptions() *GenerateConfigOptions {
	return &GenerateConfigOptions{}
}

func (o *GenerateConfigOptions) Complete(configAccess config.ConfigAccess) {
	if o.Dir == "" {
		o.Dir = filepath.Dir(configAccess.GetDefaultFilename())
	}

	abs, err := filepath.Abs(o.Dir)
	if err != nil {
		log.Fatal(err)
	}

	o.Dir = abs
}

func (o *GenerateConfigOptions) Run(configAccess config.ConfigAccess) {
	config.ModifyFile(configAccess, *config.EmptyConfig())
}

func NewCmdGenerateConfig(configAccess config.ConfigAccess) *cobra.Command {
	o := NewGenerateConfigOptions()

	cmd := &cobra.Command{
		Use:                   "config [--dir]",
		Short:                 "Generates a new EMPTY config file",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			o.Complete(configAccess)
			o.Run(configAccess)
		},
	}

	cmd.Flags().StringVarP(&o.Dir, "dir", "d", "$HOME/.carousel/", "Certificate directory")

	return cmd
}
