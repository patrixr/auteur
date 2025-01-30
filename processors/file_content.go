package processors

import (
	. "github.com/patrixr/auteur/common"
	. "github.com/patrixr/auteur/core"
)

// A generic implementation of the Content interface that
// can be used by processors to returns chunks of file content.
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
