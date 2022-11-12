package manage

type manageOpts struct {
	Delete struct {
		Source string
		Query  []string
	}
}
