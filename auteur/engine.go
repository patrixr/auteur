package auteur

type AuteurEngine struct{}

func NewAuteurEngine() *AuteurEngine {
	return &AuteurEngine{}
}

func (ae *AuteurEngine) Page(name string, markdown string) {
	// Implementation here
}

func (ae *AuteurEngine) Section(page string, title string, markdown string) {
	// Implementation here
}

func (ae *AuteurEngine) Generate(folder string) error {
	// Implementation here
	return nil
}
