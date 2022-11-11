package ui

import (
	"fmt"
)

type panel struct {
	// Values are initialized if panel is added as sub-panel
	isChild bool
	added   bool

	options []*option

	prompt func() string

	errorHandler func(error)
}

func NewPanel() *panel {
	p := &panel{}
	p.prompt = func() string {
		return "Select an option"
	}

	p.isChild = false
	p.added = false

	p.errorHandler = func(err error) {
		panic(err)
	}

	return p
}

func (p *panel) AddOption(n string) *option {
	o := newOption(n)
	p.options = append(p.options, o)

	return o
}

func (p *panel) SetPrompt(g func() string) *panel {
	p.prompt = g
	return p
}

func (p *panel) ErrorHandler(h func(error)) *panel {
	p.errorHandler = h
	return p
}

func (p *panel) Start() bool {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Handle panel being sub-panel
	if p.isChild && !p.added {
		p.AddOption("Back").
			CustomHandler(
				func(o *option) bool {
					return false
				},
			)
		p.added = true
	}

	var optionStrings []string
	for _, o := range p.options {
		optionStrings = append(optionStrings, o.name)
	}

	for {
		// If we return to a parent panel, clear the screen
		// Clear the screen
		fmt.Print("\033[H\033[2J")

		index, err := SingleCheckboxesIndex(p.prompt(), optionStrings)
		if err != nil {
			p.errorHandler(err)
		}

		selectedOpt := p.options[index]

		if loop := selectedOpt.handler(selectedOpt); !loop {
			// If command is terminating, propagate all the way up
			return !selectedOpt.terminate
		}
	}
}
