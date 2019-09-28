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
		Prompt: &survey.Input{Message: "Real name (Optional)?"},
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
		Validate: validate_alphanumeric,
	},
}

var uri_questions = []*survey.Question{
	{
		Name: "host",
		Prompt: &survey.Input{
			Message: "Hostname or IP?",
			Default: "0.0.0.0",
		},
		Validate: validate_host,
	},
	{
		Name: "port",
		Prompt: &survey.Input{
			Message: "Port?",
			Default: "6667",
		},
		Validate: validate_port,
	},
}

var user_questions = []*survey.Question{
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "Username?"},
		Validate: validate_username,
	},
	{
		Name:     "password",
		Prompt:   &survey.Password{Message: "Password?"},
		Validate: validate_password,
	},
}
