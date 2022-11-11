package ui

type Option struct {
	Name      string
	Terminate bool
	Handler   func(o *Option) bool
}

func PanelHandler(p *panel) func(o *Option) bool {
	return func(o *Option) bool {
		p.isChild = true
		return p.Start()
	}
}

func ConfirmationHandler(prompt string, yes, no func(), error func(error)) func(o *Option) bool {
	return func(o *Option) bool {
		confirm, err := ConfirmPrompt(prompt)

		if err != nil {
			error(err)
		}

		if confirm {
			yes()
		} else {
			no()
		}

		return !o.Terminate
	}
}

func InputHandler(prompt string, input func(string), error func(error)) func(o *Option) bool {
	return func(o *Option) bool {
		res, err := InputPrompt(prompt)

		if err != nil {
			error(err)
		}

		input(res)

		return !o.Terminate
	}
}

func CheckboxHandler(prompt string, genOpts func() []string, s func([]string), e func(error)) func(o *Option) bool {
	return func(o *Option) bool {
		checkboxOptions := genOpts()
		sel, err := Checkboxes(prompt, checkboxOptions)

		if err != nil {
			e(err)
		}

		s(sel)

		return !o.Terminate
	}
}
