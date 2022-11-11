package ui

type option struct {
	name      string
	terminate bool
	handler   func(o *option) bool
}

func newOption(n string) *option {
	o := &option{}
	o.terminate = false
	o.name = n
	return o
}

func (o *option) Terminator() *option {
	o.terminate = true
	return o
}

func (o *option) PanelHandler(p *panel) *option {
	o.handler = func(o *option) bool {
		p.isChild = true
		return p.Start()
	}

	return o
}

func (o *option) ConfirmationHandler(prompt string, yes, no func(), error func(error)) *option {
	o.handler = func(o *option) bool {
		confirm, err := ConfirmPrompt(prompt)

		if err != nil {
			error(err)
		}

		if confirm {
			yes()
		} else {
			no()
		}

		return !o.terminate
	}

	return o
}

func (o *option) InputHandler(prompt string, input func(string), error func(error)) *option {
	o.handler = func(o *option) bool {
		res, err := InputPrompt(prompt)

		if err != nil {
			error(err)
		}

		input(res)

		return !o.terminate
	}

	return o
}

func (o *option) CheckboxHandler(prompt string, genOpts func() []string, s func([]string), e func(error)) *option {
	o.handler = func(o *option) bool {
		checkboxOptions := genOpts()
		sel, err := Checkboxes(prompt, checkboxOptions)

		if err != nil {
			e(err)
		}

		s(sel)

		return !o.terminate
	}

	return o
}

func (o *option) CustomHandler(h func(o *option) bool) *option {
	o.handler = h
	return o
}
