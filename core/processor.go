package core

import (
	. "github.com/patrixr/auteur/common"
)

type ContentType int

const (
	HTML ContentType = iota // HTML ContentType
	Markdown
)

type Processor interface {
	Supports(extension string) bool
	Load(site *Auteur, file string) ([]Content, error)
}

type Content interface {
	Type() ContentType
	Data() string
	Path() []string
	Title() string
	Meta() Metadata
	Priority() int
	Len() int
}
