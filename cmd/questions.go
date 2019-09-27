package cmd

import (
	"github.com/AlecAivazis/survey"
)

var ident_questions = []*survey.Question{
	{
		Name:     "nickname",
		Prompt:   &survey.Input{Message: "Nickname?"},
		Validate: survey.Required,
	},
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "User?"},
		Validate: survey.Required,
	},
	{
		Name:   "realname",
		Prompt: &survey.Input{Message: "Read name (Optional)?"},
	},
	{
		Name:   "password",
		Prompt: &survey.Password{Message: "Password (Probably optional)?"},
	},
}

var network_questions = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "What are we calling this network?"},
		Validate: survey.Required,
	},
}

var uri_questions = []*survey.Question{
	{
		Name: "address",
		Prompt: &survey.Input{
			Message: "Hostname or IP?",
		},
		Validate: survey.Required,
	},
	{
		Name: "port",
		Prompt: &survey.Input{
			Message: "Port?",
		},
		Validate: survey.Required,
	},
}

var user_questions = []*survey.Question{
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "Username?"},
		Validate: survey.Required,
	},
	{
		Name:     "password",
		Prompt:   &survey.Password{Message: "Password?"},
		Validate: survey.Required,
	},
}
