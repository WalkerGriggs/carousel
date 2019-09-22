package cmd

import (
	"log"
	"os"

	tui "github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// menuCmd represents the menu command
var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Launch an interactive menu (recommended)",
	Long:  "Launch an interactive menu (recommended)",
	Run: func(cmd *cobra.Command, args []string) {
		result := menu()

		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		switch result {
		case "Serve":
			serve()
		case "Generate Configuration":
			generateConfig(home)
		default:
			os.Exit(0)
		}

	},
}

func init() {
	rootCmd.AddCommand(menuCmd)
}

// menu lets the user select menu items from an interactive menu
func menu() string {
	prompt := tui.Select{
		Label: "How can we help you?",
		Items: []string{
			"Serve",
			"Generate Configuration",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// generateConfig leverages a number of pipelines to fill in an entire, basic
// configuration. The configuration is then written to the user's home directory
// by default.
// TODO: Allow the user to specify directory.
func generateConfig(home string) {
	server, err := runConfigPipeline()
	if err != nil {
		log.Fatal(err)
	}

	js, err := json.MarshalIndent(server, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(home+"/.carousel/config.json", js, 0644)
}

// runConfigPipeline does the heavy lifting for generateConfig. It runs each
// pipeline, and, if the user wants, sets up a Network for the newly created
// user.
func runConfigPipeline() (*carousel.Server, error) {
	server, err := RunServerPipeline()
	if err != nil {
		return nil, err
	}

	cont, err := continuePrompt("Create a new network for that user?")
	if err != nil {
		return nil, err
	}

	if cont {
		network, err := RunNetworkPipeline()
		if err != nil {
			log.Fatal()
		}

		server.Users[0].Network = network
	}

	return server, nil
}

func continuePrompt(msg string) (bool, error) {
	prompt := tui.Prompt{
		Label:   msg,
		Default: "true",
	}

	res, err := prompt.Run()
	if err != nil {
		return false, err
	}

	b, err := strconv.ParseBool(res)
	if err != nil {
		return false, err
	}

	return b, nil
}
