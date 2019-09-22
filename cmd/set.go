package cmd

import "github.com/spf13/cobra"

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use: "set",
}

func init() {
	rootCmd.AddCommand(setCmd)
}
