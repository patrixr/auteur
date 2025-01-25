package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	. "github.com/patrixr/auteur/common"
	"github.com/patrixr/q"
)

var (
	// Use Sprintf to avoid being detected by auteur itself
	BEGIN_TAG = fmt.Sprintf("@%s", "auteur")
	END_TAG   = fmt.Sprintf("@%s", "end")
)

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
		LineComment: []string{`//`, `*`},
		LineBegin:   `\*`,
	}

	PYTHON_STYLE = CommentStyle{
		BlockStart:  `"""`,
		BlockEnd:    `"""`,
		LineComment: []string{`#`},
		LineBegin:   ``,
	}

	RUBY_STYLE = CommentStyle{
		BlockStart:  `=begin`,
		BlockEnd:    `=end`,
		LineComment: []string{`#`},
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
		LineComment: []string{`//`},
		LineBegin:   ``,
	}

	HASH_STYLE = CommentStyle{
		BlockStart:  ``,
		BlockEnd:    ``,
		LineComment: []string{`#`},
		LineBegin:   ``,
	}

	LUA_STYLE = CommentStyle{
		BlockStart:  `--\[\[`,
		BlockEnd:    `\]\]`,
		LineComment: []string{`--`},
		LineBegin:   ``,
	}

	HASKELL_STYLE = CommentStyle{
		BlockStart:  `{-`,
		BlockEnd:    `-}`,
		LineComment: []string{`--`},
		LineBegin:   ``,
	}

	SQL_STYLE = CommentStyle{
		BlockStart:  `\/\*`,
		BlockEnd:    `\*\/`,
		LineComment: []string{`--`},
		LineBegin:   ``,
	}

	PERL_STYLE = CommentStyle{
		BlockStart:  `=pod`,
		BlockEnd:    `=cut`,
		LineComment: []string{`#`},
		LineBegin:   ``,
	}

	MATLAB_STYLE = CommentStyle{
		BlockStart:  `%{`,
		BlockEnd:    `%}`,
		LineComment: []string{`%`},
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
		LineComment: []string{`//`, `#`, `*`},
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

func (r *CommentProcessor) Load(_ *Site, file string) ([]Content, error) {
	Logf("Reading %s", file)

	content, err := os.ReadFile(file)
	if err != nil {
		return []Content{}, err
	}

	return r.LoadFromString(string(content), getCommentStyle(file))
}

func (r *CommentProcessor) LoadFromString(content string, style CommentStyle) ([]Content, error) {
	out := []Content{}

	// find all beginning to end blocks
	pattern := fmt.Sprintf(`(?m)%s\s*(.*?)%s`, BEGIN_TAG, END_TAG)
	regex := regexp.MustCompile(`(?s)` + pattern)
	matches := regex.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}

		data := trimCommentArtifacts(match[1], style)

		meta, html, err := MarkdownToHTMLWithMeta([]byte(data))
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

		path := strings.Split(fm.Path, "/")

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

func trimCommentArtifacts(content string, style CommentStyle) string {
	lines := strings.Split(content, "\n")

	// Build regex pattern for comment prefixes
	var prefixesToIgnore []string
	for _, comment := range style.LineComment {
		prefixesToIgnore = append(prefixesToIgnore, regexp.QuoteMeta(comment))
	}
	pattern := fmt.Sprintf(`^(?:%s)`, strings.Join(prefixesToIgnore, "|"))
	commentRegex := regexp.MustCompile(pattern)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		if !commentRegex.MatchString(trimmed) {
			// Lines are not all starting with the same comment prefix
			// we can't clean the prefixes
			return content
		}
	}

	// Second pass: clean the prefixes since we know they're consistent
	cleaned := make([]string, len(lines))
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if len(trimmed) == 0 {
			cleaned[i] = trimmed
			continue
		}

		cleanLine := commentRegex.ReplaceAllString(trimmed, "")
		cleaned[i] = cleanLine
	}

	return q.TrimIndent(
		strings.Join(cleaned, "\n"),
	)
}
