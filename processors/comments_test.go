package processors

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/patrixr/auteur/core"
	"github.com/patrixr/q"
	"github.com/stretchr/testify/assert"
)

func TestExtractAuteurMetadata(t *testing.T) {
	t.Run("test detect auteur tag with multiple args", func(t *testing.T) {
		text := q.TrimIndent(`
			@auteur("arg1", "arg2")

			# Hello
		`)

		include, args, trimmed := extractAuteurMetaFromComment(text)
		assert.True(t, include)
		assert.Equal(t, []string{"arg1", "arg2"}, args)
		assert.Equal(t, "# Hello", trimmed)
	})

	t.Run("test detect auteur tag with single arg", func(t *testing.T) {
		text := q.TrimIndent(`
			@auteur("arg1")

			# Hello
		`)

		include, args, trimmed := extractAuteurMetaFromComment(text)
		assert.True(t, include)
		assert.Equal(t, []string{"arg1"}, args)
		assert.Equal(t, "# Hello", trimmed)
	})

	t.Run("test detect auteur tag with no args", func(t *testing.T) {
		text := q.TrimIndent(`
			@auteur()

			# Hello
		`)

		include, args, trimmed := extractAuteurMetaFromComment(text)
		assert.True(t, include)
		assert.Equal(t, []string{}, args)
		assert.Equal(t, "# Hello", trimmed)
	})

	t.Run("test detect auteur tag with no args and no parenthesis", func(t *testing.T) {
		text := q.TrimIndent(`
			@auteur

			# Hello
		`)

		include, args, trimmed := extractAuteurMetaFromComment(text)
		assert.True(t, include)
		assert.Equal(t, []string{}, args)
		assert.Equal(t, "# Hello", trimmed)
	})
}

func TestFindCommentsInText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		style    CommentStyle
		expected []string
	}{
		{
			name: "C-style block comments",
			input: `
                /* First comment */
                code here
                /* Multi
                   line
                   comment */
            `,
			style: C_STYLE,
			expected: []string{
				"First comment",
				"Multi\n                   line\n                   comment",
			},
		},
		{
			name: "C-style line comments",
			input: `
// First comment
code here
//   Second comment
    // Indented comment
//Multiple
//Consecutive
//Comments
		          `,
			style: C_STYLE,
			expected: []string{
				"First comment",
				"    Second comment\n" +
					" Indented comment\n" +
					"Multiple\n" +
					"Consecutive\n" +
					"Comments\n",
			},
		},
		{
			name: "Lua-style comments",
			input: `
		              --[[ Block comment ]]
		          `,
			style: LUA_STYLE,
			expected: []string{
				"Block comment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findCommentsInText(tt.input, tt.style)

			if len(got) != len(tt.expected) {
				t.Errorf("findCommentsInText() got %d comments, expected %d comments",
					len(got), len(tt.expected))
				return
			}

			for i := range got {
				got[i] = strings.TrimSpace(got[i])
				if got[i] != strings.TrimSpace(tt.expected[i]) {
					t.Errorf("findCommentsInText() got[%d] = %q, expected %q",
						i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestAuteurCommentParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		ext      string
		expected string
	}{
		{
			name: "Double slash comments",
			input: q.Paragraph(`
				// @auteur
				// Hello
				// World
			`),
			ext:      ".c",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Slash star comment",
			input: q.Paragraph(`
				/* @auteur
				Hello
				World  */
			`),
			ext:      ".c",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Python style comment",
			input: q.Paragraph(`
				# @auteur
				# Hello
				# World
			`),
			ext:      ".py",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Ruby style comment",
			input: q.Paragraph(`
				# @auteur
				# Hello
				# World
			`),
			ext:      ".rb",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Ruby style comment with =begin",
			input: q.Paragraph(`
				=begin @auteur
				Hello
				World
				=end
			`),
			ext:      ".rb",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Lua style comment",
			input: q.Paragraph(`
				--[[
					@auteur
					Hello
					World
				]]
			`),
			ext:      ".lua",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Python style comment with headings",
			input: q.Paragraph(`
				# @auteur
				# # Heading 1
				# Hello
				#
				# Heading 2
				# ---------
				# World
			`),
			ext:      ".py",
			expected: "<h1>Heading 1</h1>\n<p>Hello</p>\n<h2>Heading 2</h2>\n<p>World</p>\n",
		},
		{
			name: "Python style comment with triple quotes",
			input: q.Paragraph(`
				"""
				@auteur
				Hello
				World
				"""
			`),
			ext:      ".py",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "HTML comment",
			input: q.Paragraph(`
				<!-- @auteur
				Hello
				World
			  -->
			`),
			ext:      ".html",
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "SQL style comment",
			input: q.Paragraph(`
				-- @auteur
				-- Hello
				-- World
			`),
			ext:      ".sql",
			expected: "<p>Hello\nWorld</p>\n",
		},
	}

	reader := CommentProcessor{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			tempfile := filepath.Join(tmpdir, tt.name+tt.ext)
			os.WriteFile(tempfile, []byte(tt.input), 0644)
			auteur, err := NewAuteur()

			assert.NoError(t, err)
			got, err := reader.Load(auteur, tempfile)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			assert.Equal(t, tt.expected, got[0].Data(), "Test case: %s\nInput:\n%s", tt.name, tt.input)
		})
	}
}
