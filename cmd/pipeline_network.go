package cmd

import (
	"strconv"

	tui "github.com/manifoldco/promptui"

	"github.com/walkergriggs/carousel/carousel"
)

// RunNetworkPipleline reads in all Network, Uri, and Identity information, and
// returns a single Network object. Passwords are never stored in plaintext, and
// port ranges are verified.
func RunNetworkPipeline() (*carousel.Network, error) {
	res, err := NetworkPipeline.Run()
	if err != nil {
		return nil, err
	}

	password := res["Password"]
	if password != "" {
		password, err = carousel.Hash(res["Password"])
		if err != nil {
			return nil, err
		}
	}

	port, err := strconv.Atoi(res["Port"])
	if err != nil {
		return nil, err
	}

	return &carousel.Network{
		Name: res["Name"],
		URI: carousel.URI{
			Address: res["Address"],
			Port:    port,
		},
		Ident: carousel.Identity{
			Nickname: res["Nickname"],
			Username: res["Username"],
			Realname: res["Realname"],
			Password: password,
		},
	}, nil
}

var NetworkPipeline = Pipeline{
	Prompts: []Prompt{
		NamePrompt,
		NetworkAddressPrompt,
		NetworkPortPrompt,
		NicknamePrompt,
		UsernamePrompt,
		RealnamePrompt,
		PasswordPrompt,
	},
}

var NamePrompt = Prompt{
	Prompt: &tui.Prompt{
		Label:   "Network name",
		Default: "Freenode",
	},
	Field: "Name",
}

var NetworkAddressPrompt = Prompt{
	Prompt: &tui.Prompt{
		Label:   "Network hostname",
		Default: "chat.freenode.net",
	},
	Field: "Address",
}

var NetworkPortPrompt = Prompt{
	Prompt: &tui.Prompt{
		Label:    "Network port",
		Default:  "6667",
		Validate: validatePort,
	},
	Field: "Port",
}

var NicknamePrompt = GenericPrompt("Network nick", "Nickname")
var UsernamePrompt = GenericPrompt("Network user", "Username")
var RealnamePrompt = GenericPrompt("Real name", "Realname")
var PasswordPrompt = GenericPrompt("SASL password", "Password")
