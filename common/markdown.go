package common

import (
	"bytes"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/mermaid"
)

// Project-wide markdown converter
var converter = goldmark.New(
	goldmark.WithExtensions(
		&frontmatter.Extender{},
		&mermaid.Extender{
			RenderMode: mermaid.RenderModeClient,
		},
		extension.NewTable(),
	),
)

func MarkdownToHTML(md []byte) (string, error) {
	var buf bytes.Buffer
	if err := ConvertMarkdown(md, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func MarkdownToHTMLWithMeta(md []byte) (Metadata, string, error) {
	var buf bytes.Buffer
	meta, err := ConvertMarkdownWithMeta(md, &buf)
	if err != nil {
		return meta, "", err
	}
	return meta, buf.String(), nil
}

func ConvertMarkdown(md []byte, w io.Writer) error {
	_, err := ConvertMarkdownWithMeta(md, w)
	return err
}

func ConvertMarkdownWithMeta(md []byte, w io.Writer) (Metadata, error) {
	var meta Metadata

	ctx := parser.NewContext()

	if err := converter.Convert(md, w, parser.WithContext(ctx)); err != nil {
		return meta, err
	}

	d := frontmatter.Get(ctx)
	if d == nil {
		return meta, nil
	}

	if err := d.Decode(&meta); err != nil {
		return meta, err
	}

	return meta, nil
}
