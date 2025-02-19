package processors

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/patrixr/auteur/common"
	. "github.com/patrixr/auteur/core"
	"github.com/patrixr/q"
)

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
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.ReplaceAll(title, "-", " ")
	relPath, err := site.GetRelativePath(file)
	if err != nil {
		return []Content{}, err
	}

	path := strings.Split(relPath, "/")

	filename, _ := q.Last(path)
	filename = strings.Split(filename, ".")[0]
	if strings.EqualFold(filename, "readme") || strings.EqualFold(filename, "index") {
		path = path[:len(path)-1]
	} else {
		filename = strings.ToLower(filename)
		filename = strings.ReplaceAll(filename, "_", " ")
		filename = strings.ReplaceAll(filename, "-", " ")
		path[len(path)-1] = filename
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
