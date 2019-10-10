package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create one or many resources",
	Long:  "Create one or many resources",
}

func init() {
	rootCmd.AddCommand(createCmd)
}
