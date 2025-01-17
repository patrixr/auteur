package auteur

type Auteur interface {
	Write(chunks ...string)
	Generate(foloutfolderder string) error
}
