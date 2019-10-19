package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/walkergriggs/carousel/cmd"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

func main() {
	cmd.Execute()
}
