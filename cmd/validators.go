package cmd

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func validate_password(val interface{}) error {
	return validation.Validate(val,
		validation.Required,
		validation.Length(8, 32),
		validation.Match(regexp.MustCompile(`^[\w\d@%+/'!#$^?;,()-_.]+$`)),
	)
}

func validate_username(val interface{}) error {
	return validation.Validate(val,
		validation.Required,
		validation.Length(3, 16),
		validation.Match(regexp.MustCompile(`^[\w\d]+$`)),
	)
}

func validate_host(val interface{}) error {
	return validate_string_rule(val, is.Host)
}

func validate_port(val interface{}) error {
	return validate_string_rule(val, is.Port)
}

func validate_alphanumeric(val interface{}) error {
	return validate_string_rule(val, is.Alphanumeric)
}

func validate_string_rule(val interface{}, rule *validation.StringRule) error {
	return validation.Validate(val,
		validation.Required,
		rule,
	)
}
