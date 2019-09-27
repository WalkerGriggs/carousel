package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/AlecAivazis/survey"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/walkergriggs/carousel/carousel"
)

// menuCmd represents the menu command
var configCmd = &cobra.Command{
	Use:   "configure",
	Short: "Launch an interactive menu (recommended)",
	Long:  "Launch an interactive menu (recommended)",
	Run: func(cmd *cobra.Command, args []string) {
		if survey_confirm("This command will overwrite any existing config. Continue?") {
			server := survey_server()

			home, err := homedir.Dir()
			if err != nil {
				log.Fatal(err)
			}

			js, err := json.MarshalIndent(server, "", "    ")
			if err != nil {
				log.Fatal(err)
			}

			ioutil.WriteFile(home+"/.carousel/config.json", js, 0644)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func survey_server() carousel.Server {
	fmt.Println("Lets start with Carousel's address.")
	uri := survey_uri()

	fmt.Println("Now, we need to set up an admin user.")
	admin := survey_user()

	server := carousel.Server{
		Users: []*carousel.User{&admin},
		URI:   uri,
	}

	return server
}

func survey_user() carousel.User {
	var user carousel.User
	if err := survey.Ask(user_questions, &user); err != nil {
		log.Fatal(err)
	}

	hashed_pass, err := carousel.Hash(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	user.Password = hashed_pass

	if survey_confirm("Do you want to setup a Network for this user?") {
		net := survey_network()
		user.Network = &net
	}

	return user
}

func survey_network() carousel.Network {
	var net carousel.Network
	if err := survey.Ask(network_questions, &net); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Where can we find this network?")
	net.URI = survey_uri()

	fmt.Println("Almost done! We just need to get your network identity.")
	net.Ident = survey_identity()

	return net
}

func survey_uri() carousel.URI {
	var uri carousel.URI
	if err := survey.Ask(uri_questions, &uri); err != nil {
		log.Fatal(err)
	}

	return uri
}

func survey_identity() carousel.Identity {
	var ident carousel.Identity
	err := survey.Ask(ident_questions, &ident)
	if err != nil {
		log.Fatal(err)
	}

	hashed_pass, err := carousel.Hash(ident.Password)
	if err != nil {
		log.Fatal(err)
	}

	ident.Password = hashed_pass
	return ident
}

func survey_confirm(prompt string) bool {
	b := false
	confirm := &survey.Confirm{
		Message: prompt,
	}
	survey.AskOne(confirm, &b)
	return b
}
