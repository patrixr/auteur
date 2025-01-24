package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	. "github.com/patrixr/auteur/common"
)

type CommentStyle struct {
	Start string
	End   string
}

var commentStyles = map[string]CommentStyle{
	".go":   {Start: `\/\*`, End: `\*\/`},
	".js":   {Start: `\/\*`, End: `\*\/`},
	".css":  {Start: `\/\*`, End: `\*\/`},
	".java": {Start: `\/\*`, End: `\*\/`},
	".php":  {Start: `\/\*`, End: `\*\/`},
	".py":   {Start: `"""`, End: `"""`},
	".rb":   {Start: `=begin`, End: `=end`},
	".rs":   {Start: `\/\*`, End: `\*\/`},
	".ts":   {Start: `\/\*`, End: `\*\/`},
	".zig":  {Start: `\/\*`, End: `\*\/`},
}

type CommentProcessor struct {
	folder string
}

type CommentFrontmatter struct {
	Auteur string `yaml:"auteur"`
	Title  string `yaml:"title"`
}

func NewCommentReader(folder string) Processor {
	return &CommentProcessor{
		folder: folder,
	}
}

func (r *CommentProcessor) Supports(extension string) bool {
	_, ok := commentStyles[extension]
	return ok
}

func (r *CommentProcessor) Load(file string) ([]Content, error) {
	Logf("Reading %s", file)

	out := []Content{}

	content, err := os.ReadFile(file)
	if err != nil {
		return out, err
	}

	style := GetCommentStyle(file)

	// Look for comments that start with ---
	commentRegex := regexp.MustCompile(`(?s)` + style.Start + `\s*(---\n[^\n]*\n---\n.*?)` + style.End)
	matches := commentRegex.FindAllStringSubmatch(string(content), -1)

	if len(matches) > 0 {
		fmt.Println("MATCHES::")
		fmt.Println(matches)
	}

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}

		data := match[1]

		meta, html, err := MarkdownToHTMLWithMeta([]byte(data))
		if err != nil {
			return out, err
		}

		fm, err := MetaToStruct[CommentFrontmatter](meta)
		if err != nil {
			return out, err
		}

		if len(fm.Auteur) == 0 {
			continue
		}

		path := strings.Split(fm.Auteur, "/")

		out = append(out, &ContentData{
			metadata: meta,
			data:     html,
			path:     path,
			kind:     HTML,
			title:    fm.Auteur,
		})
	}

	return out, nil
}

func GetCommentStyle(filename string) CommentStyle {
	for ext, style := range commentStyles {
		if strings.HasSuffix(filename, ext) {
			return style
		}
	}
	return CommentStyle{Start: `\/\*`, End: `\*\/`}
}
