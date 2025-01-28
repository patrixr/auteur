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
	Order() int
	Len() int
}

type ContentData struct {
	kind     ContentType
	data     string
	path     []string
	metadata Metadata
	title    string
	order    int
}

func (c *ContentData) Type() ContentType {
	return c.kind
}

func (c *ContentData) Order() int {
	return c.order
}

func (c *ContentData) Data() string {
	return c.data
}

func (c *ContentData) Path() []string {
	return c.path
}

func (c *ContentData) Meta() Metadata {
	return c.metadata
}

func (c *ContentData) Title() string {
	return c.title
}

func (c *ContentData) Len() int {
	return len(c.data)
}
