package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/walkergriggs/carousel/pkg/cadm"
)

func init() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}

	log.SetOutput(os.Stdout)
	log.SetFormatter(formatter)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cadm.Execute()
}
