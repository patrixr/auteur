package auteur

type Auteur interface {
	Page(name string, markdown string)
	Section(page string, title string, markdown string)
	Generate(foloutfolderder string) error
}
