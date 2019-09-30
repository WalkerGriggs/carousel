package cmd

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey"

	"github.com/walkergriggs/carousel/crypto"
	"github.com/walkergriggs/carousel/network"
	"github.com/walkergriggs/carousel/server"
	"github.com/walkergriggs/carousel/uri"
	"github.com/walkergriggs/carousel/user"
)

func survey_server() server.Server {
	var server server.Server

	fmt.Println("Lets start with Carousel's address.")
	server.URI = survey_uri()

	if err := survey.Ask(server_questions, &server); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Now, we need to set up an admin user.")
	server.Users = []*user.User{survey_user()}

	return server
}

func survey_user() *user.User {
	var user user.User
	if err := survey.Ask(user_questions, &user); err != nil {
		log.Fatal(err)
	}

	hashed_pass, err := crypto.Hash(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	user.Password = hashed_pass

	if survey_confirm("Do you want to setup a Network for this user?") {
		net := survey_network()
		user.Network = &net
	}

	return &user
}

func survey_network() network.Network {
	var net network.Network
	if err := survey.Ask(network_questions, &net); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Where can we find this network?")
	net.URI = survey_uri()

	fmt.Println("Almost done! We just need to get your network identity.")
	net.Ident = survey_identity()

	return net
}

func survey_uri() uri.URI {
	var uri uri.URI
	if err := survey.Ask(uri_questions, &uri); err != nil {
		log.Fatal(err)
	}

	return uri
}

func survey_identity() network.Identity {
	var ident network.Identity
	err := survey.Ask(ident_questions, &ident)
	if err != nil {
		log.Fatal(err)
	}

	if ident.Password != "" {
		hashed_pass, err := crypto.Hash(ident.Password)
		if err != nil {
			log.Fatal(err)
		}

		ident.Password = hashed_pass
	}

	if ident.Realname == "" {
		ident.Realname = ident.Username
	}

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
