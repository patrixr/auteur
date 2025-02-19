package processors

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	. "github.com/patrixr/auteur/common"
	. "github.com/patrixr/auteur/core"
	"github.com/patrixr/q"
)

const AUTEUR_TAG = "@auteur"

type CommentStyle struct {
	BlockStart  string
	BlockEnd    string
	LineComment []string
	LineBegin   string
}

var (
	C_STYLE = CommentStyle{
		BlockStart:  `\/\*`,
		BlockEnd:    `\*\/`,
		LineComment: []string{`[/]{2,}`},
		LineBegin:   `\*`,
	}

	PYTHON_STYLE = CommentStyle{
		BlockStart:  `"""`,
		BlockEnd:    `"""`,
		LineComment: []string{`#+`},
		LineBegin:   ``,
	}

	RUBY_STYLE = CommentStyle{
		BlockStart:  `=begin`,
		BlockEnd:    `=end`,
		LineComment: []string{`#+`},
		LineBegin:   ``,
	}

	HTML_STYLE = CommentStyle{
		BlockStart:  `<!--`,
		BlockEnd:    `-->`,
		LineComment: []string{},
		LineBegin:   ``,
	}

	CSS_STYLE = CommentStyle{
		BlockStart:  `\/\*`,
		BlockEnd:    `\*\/`,
		LineComment: []string{},
		LineBegin:   ``,
	}

	SCSS_STYLE = CommentStyle{
		BlockStart:  `\/\*`,
		BlockEnd:    `\*\/`,
		LineComment: []string{`/{2,}`},
		LineBegin:   ``,
	}

	HASH_STYLE = CommentStyle{
		BlockStart:  ``,
		BlockEnd:    ``,
		LineComment: []string{`#{1,}`},
		LineBegin:   ``,
	}

	LUA_STYLE = CommentStyle{
		BlockStart:  `--\[\[`,
		BlockEnd:    `\]\]`,
		LineComment: []string{},
		LineBegin:   ``,
	}

	HASKELL_STYLE = CommentStyle{
		BlockStart:  `{-`,
		BlockEnd:    `-}`,
		LineComment: []string{`-{2,}`},
		LineBegin:   ``,
	}

	SQL_STYLE = CommentStyle{
		BlockStart:  `\/\*`,
		BlockEnd:    `\*\/`,
		LineComment: []string{`-{2,}`},
		LineBegin:   ``,
	}

	PERL_STYLE = CommentStyle{
		BlockStart:  `=pod`,
		BlockEnd:    `=cut`,
		LineComment: []string{`#+`},
		LineBegin:   ``,
	}

	MATLAB_STYLE = CommentStyle{
		BlockStart:  `%{`,
		BlockEnd:    `%}`,
		LineComment: []string{`%+`},
		LineBegin:   ``,
	}

	VB_STYLE = CommentStyle{
		BlockStart:  ``,
		BlockEnd:    ``,
		LineComment: []string{`'`},
		LineBegin:   ``,
	}

	PHP_STYLE = CommentStyle{
		BlockStart:  `\/\*`,
		BlockEnd:    `\*\/`,
		LineComment: []string{`/{2,}`},
		LineBegin:   `\*`,
	}
)

var commentStyles = map[string]CommentStyle{
	".go":    C_STYLE,
	".js":    C_STYLE,
	".jsx":   C_STYLE,
	".ts":    C_STYLE,
	".tsx":   C_STYLE,
	".java":  C_STYLE,
	".c":     C_STYLE,
	".cpp":   C_STYLE,
	".cs":    C_STYLE,
	".kt":    C_STYLE,
	".swift": C_STYLE,
	".scala": C_STYLE,
	".rs":    C_STYLE,
	".py":    PYTHON_STYLE,
	".rb":    RUBY_STYLE,
	".php":   PHP_STYLE,
	".r":     HASH_STYLE,
	".yaml":  HASH_STYLE,
	".yml":   HASH_STYLE,
	".toml":  HASH_STYLE,
	".sh":    HASH_STYLE,
	".bash":  HASH_STYLE,
	".zsh":   HASH_STYLE,
	".html":  HTML_STYLE,
	".xml":   HTML_STYLE,
	".svg":   HTML_STYLE,
	".css":   CSS_STYLE,
	".scss":  SCSS_STYLE,
	".sass":  SCSS_STYLE,
	".less":  SCSS_STYLE,
	".lua":   LUA_STYLE,
	".hs":    HASKELL_STYLE,
	".sql":   SQL_STYLE,
	".pl":    PERL_STYLE,
	".m":     MATLAB_STYLE,
	".vb":    VB_STYLE,
}

type CommentProcessor struct{}

func NewCommentReader() Processor {
	return &CommentProcessor{}
}

func (r *CommentProcessor) Supports(extension string) bool {
	_, ok := commentStyles[extension]
	return ok
}

func (r *CommentProcessor) Load(auteur *Auteur, file string) ([]Content, error) {
	Logf("Reading %s", file)

	content, err := os.ReadFile(file)
	if err != nil {
		return []Content{}, err
	}

	style := getCommentStyle(file)

	out := []Content{}

	comments := findCommentsInText(string(content), style)

	relPath, err := auteur.GetRelativePath(file)
	folderPath := filepath.Dir(relPath)
	// Default to the path of the file the content is contained in
	path := strings.Split(folderPath, "/")

	for _, comment := range comments {
		include, args, trimmed := extractAuteurMetaFromComment(comment)

		if !include {
			continue
		}

		meta, html, err := MarkdownToHTMLWithMeta([]byte(trimmed))
		if err != nil {
			return out, err
		}

		fm, err := MetaToStruct[AuteurFrontmatter](meta)
		if err != nil {
			return out, err
		}

		if fm.Ignore {
			continue
		}

		// If a path is provided as arg or in the frontmatter, use that instead
		if len(args) > 0 {
			path = strings.Split(args[0], "/")
		} else if fm.Path != "" {
			path = strings.Split(fm.Path, "/")
		}

		fmt.Println("folderPath", folderPath)
		fmt.Println("path", path)
		fmt.Println("args", args)
		out = append(out, &ContentData{
			metadata: meta,
			data:     html,
			path:     path,
			kind:     HTML,
			title:    fm.Title,
			order:    fm.Order,
		})
	}

	return out, nil
}

// -----------------------------------
// Helpers
// -----------------------------------

func getCommentStyle(filename string) CommentStyle {
	for ext, style := range commentStyles {
		if strings.HasSuffix(filename, ext) {
			return style
		}
	}
	return CommentStyle{BlockStart: `\/\*`, BlockEnd: `\*\/`, LineComment: []string{`//`, `\*`}, LineBegin: `\*`}
}

// Given a text, tries to find if there is an @auteur(...)
// tag inside of it. If yes, returns true, alongside the arguments of the tag
// and the text trimmed of said tag
func extractAuteurMetaFromComment(text string) (present bool, args []string, trimmed string) {
	argsPattern := `\s*((?:("[^"\n]*")\s*,\s*)*(?:("[^"\n]*")))?\s*`
	auteurPattern := fmt.Sprintf(`%s(?:[\s\n]|$)|%s\(%s\)`, AUTEUR_TAG, AUTEUR_TAG, argsPattern)
	auteurRexp := regexp.MustCompile(auteurPattern)
	auteurMatches := auteurRexp.FindAllStringSubmatch(text, -1)

	if len(auteurMatches) == 0 {
		trimmed = text
		return
	}

	present = true
	trimmed = strings.Trim(q.TrimIndent(auteurRexp.ReplaceAllString(text, "")), "\n")

	// Check to see if we've found some arguments
	fmt.Println("matches", auteurMatches)
	if len(auteurMatches[0]) > 2 {
		args = q.Filter(auteurMatches[0][2:], func(arg string) bool {
			return arg != ""
		})

		for i, arg := range args {
			args[i] = strings.Trim(arg, `"`)
		}
	}

	return
}

// Given a code file, finds all the comment blocks present inside ot it
// Comment symbols (//) are trimmed before returning the text content
func findCommentsInText(text string, style CommentStyle) []string {
	comments := []string{}

	// Block comments
	pattern := fmt.Sprintf(`(?s)%s(.*?)%s`, style.BlockStart, style.BlockEnd)
	matches := regexp.MustCompile(pattern).FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		comments = append(comments, match[1])
	}

	// Line comments
	for _, comment := range style.LineComment {
		pattern = fmt.Sprintf(`(?m)(^\s*%s[^\n]*\n)+`, comment)
		blocks := regexp.MustCompile(pattern).FindAllString(text, -1)
		for _, block := range blocks {
			trimPattern := fmt.Sprintf(`(?m)^\s*%s`, comment)
			trimmed := regexp.MustCompile(trimPattern).ReplaceAllString(block, "")
			comments = append(comments, trimmed)
		}
	}

	return comments
}
