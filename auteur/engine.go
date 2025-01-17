package auteur

type AuteurEngine struct{}

func NewAuteurEngine() *AuteurEngine {
	return &AuteurEngine{}
}

func (ae *AuteurEngine) Write(chunks ...string) {
	// Implementation here
}

func (ae *AuteurEngine) Generate(folder string) error {
	// Implementation here
	return nil
}
