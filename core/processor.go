package core

import (
	. "github.com/patrixr/auteur/common"
)

/*
 * @auteur
 * ---
 * path: processors
 * ---
 *
 * # Processors
 *
 * Auteur processors are responsible for loading and parsing files into a structured format that can be used by the Auteur engine
 *
 * ## Comments
 *
 * The "Comment" processor is responsible for parsing comments from various code languages and convert them into a structured format.
 * It supports multiple programming languages and can be extended to support more.
 *
 * ## Supported Languages
 *
 * - Go
 *   - Python
 * - JavaScript
 * - Java
 * - C++
 * @end
 */

type ContentType int

const (
	HTML ContentType = iota // HTML ContentType
	Markdown
)

type Processor interface {
	Supports(extension string) bool
	Load(site *Site, file string) ([]Content, error)
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
