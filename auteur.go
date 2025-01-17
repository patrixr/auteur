package auteur

import "github.com/patrixr/auteur/core"

type Auteur interface {
	Write(chunks ...string)
	Generate(foloutfolderder string) error
}

var globalInstance Auteur = core.NewAuteurEngine()

func Write(chunks ...string) {
	globalInstance.Write(chunks...)
}

func Generate(outfolder string) error {
	return globalInstance.Generate(outfolder)
}

func New() Auteur {
	return core.NewAuteurEngine()
}
