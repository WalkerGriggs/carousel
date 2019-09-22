package cmd

import (
	tui "github.com/manifoldco/promptui"
)

// Pipeline represents a series of Prompts, each of which can be run in order
// to present a single, cohesive dialog.
type Pipeline struct {
	Prompts []Prompt
}

// Prompt represnts a single tui prompt, but also stores the name of a struct's
// field. This field name is used in correlate the prompt's results with a
// struct's specific field without having to rely on the prompt's label.
type Prompt struct {
	Prompt *tui.Prompt
	Field  string
}

// Run run's each of the Pipeline's prompts, and packs the results into a single
// map. That map can be used later to construct the Pipeline's matching struct.
func (p Pipeline) Run() (map[string]string, error) {
	results := make(map[string]string)

	for _, prompt := range p.Prompts {
		result, err := prompt.Prompt.Run()
		if err != nil {
			return results, err
		}

		results[prompt.Field] = result
	}

	return results, nil
}

func GenericPrompt(label, field string) Prompt {
	return Prompt{
		Prompt: &tui.Prompt{
			Label: label,
		},
		Field: field,
	}
}
