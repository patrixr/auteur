package auteur

var globalInstance Auteur = NewAuteurEngine()

func Page(name string, markdown string) {
	globalInstance.Page(name, markdown)
}

func Section(page string, title string, markdown string) {
	globalInstance.Section(page, title, markdown)
}

func Generate(outfolder string) error {
	return globalInstance.Generate(folder)
}
