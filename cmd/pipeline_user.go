package cmd

import (
	tui "github.com/manifoldco/promptui"

	"github.com/walkergriggs/carousel/carousel"
)

// RunUserPipeline runs the user pipeline and parses the results into a user
// object. Passwords passed in through a prompt are immediately hashedl; their
// plain text is never stored.
func RunUserPipeline() (*carousel.User, error) {
	res, err := UserPipeline.Run()
	if err != nil {
		return nil, err
	}

	hashed_passwd, err := carousel.Hash(res["Password"])
	if err != nil {
		return nil, err
	}

	return &carousel.User{
		Username: res["Username"],
		Password: hashed_passwd,
	}, nil
}

var UserPipeline = Pipeline{
	Prompts: []Prompt{
		AdminUsernamePrompt,
		AdminPasswordPrompt,
	},
}

var AdminUsernamePrompt = Prompt{
	Prompt: &tui.Prompt{
		Label: "Admin username",
	},
	Field: "Username",
}

var AdminPasswordPrompt = Prompt{
	Prompt: &tui.Prompt{
		Label:    "Admin password",
		Validate: validatePassword,
	},
	Field: "Password",
}
