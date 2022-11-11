package ui

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

type panel struct {
	// Values are initialized if panel is added as sub-panel
	isChild bool
	added   bool

	options []Option

	prompt func() string
}

func CreatePanel() *panel {
	p := &panel{}
	p.prompt = func() string {
		return "Select an option"
	}

	p.isChild = false
	p.added = false

	return p
}

func (p *panel) AddOption(o Option) *panel {
	p.options = append(p.options, o)
	return p
}

func (p *panel) SetPrompt(g func() string) *panel {
	p.prompt = g
	return p
}

func (p *panel) Start() bool {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Handle panel being sub-panel
	if p.isChild && !p.added {
		p.AddOption(Option{
			Name: "Back",
			Handler: func(o *Option) bool {
				return false
			},
		})
		p.added = true
	}

	var optionStrings []string
	for _, o := range p.options {
		optionStrings = append(optionStrings, o.Name)
	}

	for {
		// If we return to a parent panel, clear the screen
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		var index int
		prompt := &survey.Select{
			Message: p.prompt(),
			Options: optionStrings,
		}
		err := survey.AskOne(prompt, &index, survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = ""
			icons.Question.Format = "yellow+hb"
		}))

		if err != nil {
			panic(err)
		}
		selectedOpt := p.options[index]

		if loop := selectedOpt.Handler(&selectedOpt); !loop {
			// If command is terminating, propagate all the way up
			return !selectedOpt.Terminate
		}
	}
}
