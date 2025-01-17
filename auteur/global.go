package auteur

var globalInstance Auteur = NewAuteurEngine()

func Write(chunks ...string) {
	globalInstance.Write(chunks...)
}

func Generate(outfolder string) error {
	return globalInstance.Generate(outfolder)
}
