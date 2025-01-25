package core

import (
	"testing"

	"github.com/patrixr/q"
	"github.com/stretchr/testify/assert"
)

func TestCommentExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		style    CommentStyle
		expected string
	}{
		{
			name: "Double slash comments",
			input: q.Paragraph(`
				// @auteur
				// Hello
				// World
				// @end
			`),
			style:    C_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Slash star comment",
			input: q.Paragraph(`
				/* @auteur
				// Hello
				// World
				@end */
			`),
			style:    C_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Slash star comment (with star prefix)",
			input: q.Paragraph(`
				/* @auteur
				 * Hello
				 * World
				 * @end
				 */
			`),
			style:    C_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Python style comment",
			input: q.Paragraph(`
				# @auteur
				# Hello
				# World
				# @end
			`),
			style:    PYTHON_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Ruby style comment",
			input: q.Paragraph(`
				# @auteur
				# Hello
				# World
				# @end
			`),
			style:    RUBY_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Ruby style comment with =begin",
			input: q.Paragraph(`
				=begin @auteur
				Hello
				World
				@end
				=end
			`),
			style:    RUBY_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "Lua style comment",
			input: q.Paragraph(`
				-- @auteur
				-- Hello
				-- World
				-- @end
			`),
			style:    LUA_STYLE,
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
				# @end
			`),
			style:    PYTHON_STYLE,
			expected: "<h1>Heading 1</h1>\n<p>Hello</p>\n<h2>Heading 2</h2>\n<p>World</p>\n",
		},
		{
			name: "Python style comment with triple quotes",
			input: q.Paragraph(`
				"""
				@auteur
				Hello
				World
				@end
				"""
			`),
			style:    PYTHON_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "HTML comment",
			input: q.Paragraph(`
				<!-- @auteur
				Hello
				World
				@end -->
			`),
			style:    HTML_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
		{
			name: "SQL style comment",
			input: q.Paragraph(`
				-- @auteur
				-- Hello
				-- World
				-- @end
			`),
			style:    SQL_STYLE,
			expected: "<p>Hello\nWorld</p>\n",
		},
	}

	reader := CommentProcessor{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reader.LoadFromString(tt.input, tt.style)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			assert.Equal(t, tt.expected, got[0].Data(), "Test case: %s\nInput:\n%s", tt.name, tt.input)
		})
	}
}
