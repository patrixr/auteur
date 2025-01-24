package core

import (
	. "github.com/patrixr/auteur/common"
)

/*
---
Auteur: Processors
---

# Processors

Auteur processors are responsible for loading and parsing files into a structured format that can be used by the Auteur engine.
The `Processor` interface defines the methods that a processor must implement in order to be used by the engine.
*/

/*
---
auteur: "Processors/Comments"
---

# Comments

*/

/*
---
auteur: "Processors/Comments/Nested"
---

# Comments Nested

*/

/*
---
auteur: "Processors/Comments/Nested/More"
---

# Comments

*/

type ContentType int

const (
	HTML ContentType = iota // HTML ContentType
	Markdown
)

type Processor interface {
	Supports(extension string) bool
	Load(file string) ([]Content, error)
}

type Content interface {
	Type() ContentType
	Data() string
	Path() []string
	Title() string
	Meta() Metadata
	Len() int
}

type ContentData struct {
	kind     ContentType
	data     string
	path     []string
	metadata Metadata
	title    string
}

func (c *ContentData) Type() ContentType {
	return c.kind
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
