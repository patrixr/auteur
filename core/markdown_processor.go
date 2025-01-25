package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/patrixr/auteur/common"
	"github.com/patrixr/q"
)

/**
@auteur
---
path: Processors
---

## Markdown Processor

@end
*/

type AuteurFrontmatter struct {
	Path   string `yaml:"path"`
	Title  string `yaml:"title"`
	Order  int    `yaml:"order"`
	Ignore bool   `yaml:"ignore"`
}

type MarkdownProcessor struct {
	folder string
}

func NewMarkdownProcessor() Processor {
	return &MarkdownProcessor{}
}

func (r *MarkdownProcessor) Supports(extension string) bool {
	return extension == ".md" || extension == ".markdown"
}

func (r *MarkdownProcessor) Load(site *Auteur, file string) ([]Content, error) {
	Logf("Reading %s", file)

	content, err := os.ReadFile(file)
	if err != nil {
		return []Content{}, err
	}

	title := strings.Split(filepath.Base(file), ".")[0]
	relPath, err := r.getRelativePath(site, file)
	if err != nil {
		return []Content{}, err
	}

	path := strings.Split(relPath, "/")

	filename, _ := q.Last(path)
	filename = strings.Split(filename, ".")[0]
	if strings.EqualFold(filename, "readme") || strings.EqualFold(filename, "index") {
		path = path[:len(path)-1]
	} else {
		path[len(path)-1] = strings.ToLower(filename)
	}

	meta, html, err := MarkdownToHTMLWithMeta([]byte(content))
	if err != nil {
		return []Content{}, err
	}

	fm, err := MetaToStruct[AuteurFrontmatter](meta)
	if err != nil {
		return []Content{}, err
	}

	if fm.Ignore {
		return []Content{}, nil
	}

	if fm.Title != "" {
		title = fm.Title
	}

	if len(fm.Path) != 0 {
		path = strings.Split(fm.Path, "/")
	}

	return []Content{
		&ContentData{
			metadata: meta,
			data:     string(html),
			kind:     HTML,
			title:    title,
			path:     path,
			order:    fm.Order,
		},
	}, nil
}

func (r *MarkdownProcessor) getRelativePath(site *Auteur, path string) (string, error) {
	cwd, err := filepath.Abs(site.Rootdir)

	if err != nil {
		return "", fmt.Errorf("failed to get an absolute path for working directory: %w", err)
	}

	relPath, err := filepath.Rel(cwd, path)

	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	if relPath == "." {
		return "", nil
	}

	return relPath, nil
}
