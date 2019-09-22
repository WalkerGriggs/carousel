package cmd

import (
	"strconv"

	tui "github.com/manifoldco/promptui"

	"github.com/walkergriggs/carousel/carousel"
)

// RunServerPipeline combines the User and Uri pipelines into a single function.
// Both User and Uri results are parsed into a Server object.
func RunServerPipeline() (*carousel.Server, error) {
	uri, err := RunUriPipeline()
	if err != nil {
		return nil, err
	}

	user, err := RunUserPipeline()
	if err != nil {
		return nil, err
	}

	return &carousel.Server{
		URI:   *uri,
		Users: []*carousel.User{user},
	}, nil
}

// RunUriPipeline reads in required URI fields and parses the results.
func RunUriPipeline() (*carousel.URI, error) {
	res, err := UriPipeline.Run()
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(res["Port"])
	if err != nil {
		return nil, err
	}

	return &carousel.URI{
		Address: res["Address"],
		Port:    port,
	}, nil
}

var UriPipeline = Pipeline{
	Prompts: []Prompt{
		AddressPrompt,
		PortPrompt,
	},
}

var PortPrompt = Prompt{
	Prompt: &tui.Prompt{
		Label:    "Carousel port",
		Default:  "6667",
		Validate: validatePort,
	},
	Field: "Port",
}

var AddressPrompt = Prompt{
	Prompt: &tui.Prompt{
		Label:   "Carousel address",
		Default: "0.0.0.0",
	},
	Field: "Address",
}
