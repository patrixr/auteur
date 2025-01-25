package core

import (
	"testing"

	. "github.com/patrixr/auteur/common"
	"github.com/stretchr/testify/assert"
)

type MockContent struct {
	path []string
	len  int
}

func (m MockContent) Path() []string    { return m.path }
func (m MockContent) Len() int          { return m.len }
func (m MockContent) Data() string      { return "" }
func (m MockContent) Meta() Metadata    { return Metadata{} }
func (m MockContent) Title() string     { return "" }
func (m MockContent) Type() ContentType { return Markdown }
func (m MockContent) Order() int        { return 0 }

type MockProcessor struct {
	supportedExt string
	contents     []Content
	loadErr      error
}

func (m MockProcessor) Supports(ext string) bool {
	return ext == m.supportedExt
}

func (m MockProcessor) Load(_ *Auteur, path string) ([]Content, error) {
	return m.contents, m.loadErr
}

func TestSite(t *testing.T) {
	t.Run("NewSite", func(t *testing.T) {
		site, err := NewAuteur()
		assert.NoError(t, err)
		assert.Nil(t, site.parent)
		assert.Nil(t, site.root)
		assert.Equal(t, "Auteur", site.Title)
		assert.Empty(t, site.Content)
		assert.Empty(t, site.processors)
	})

	t.Run("Basic Setters and Getters", func(t *testing.T) {
		site, err := NewAuteur()
		assert.NoError(t, err)

		site.SetTitle("Test Title")
		assert.Equal(t, "Test Title", site.Title)
		assert.Equal(t, "test-title", site.Slug())

		site.SetDesc("Test Description")
		assert.Equal(t, "Test Description", site.Desc)
	})

	t.Run("Root and Href detection", func(t *testing.T) {
		root, err := NewAuteur()
		assert.NoError(t, err)

		child := root.GetSubpage("Child", 0)
		grandchild := child.GetSubpage("Grandchild", 0)

		assert.True(t, root.IsRoot())
		assert.False(t, child.IsRoot())
		assert.False(t, grandchild.IsRoot())

		assert.Equal(t, "/", root.Href())
		assert.Equal(t, "/child", child.Href())
		assert.Equal(t, "/child/grandchild", grandchild.Href())

		assert.Equal(t, root, child.Root())
		assert.Equal(t, root, grandchild.Root())
	})

	t.Run("Children Management", func(t *testing.T) {
		site, err := NewAuteur()
		assert.NoError(t, err)

		assert.False(t, site.HasChildren())

		child1 := site.GetSubpage("Child1", 0)
		assert.True(t, site.HasChildren())
		assert.Len(t, site.Children(), 1)

		// Getting same subpage should return existing one
		child1Again := site.GetSubpage("Child1", 0)
		assert.Equal(t, child1, child1Again)
		assert.Len(t, site.Children(), 1, 0)

		// Case insensitive comparison
		child1Upper := site.GetSubpage("CHILD1", 0)
		assert.Equal(t, child1, child1Upper)
		assert.Len(t, site.Children(), 1, 0)
	})

	t.Run("Content Management", func(t *testing.T) {
		site, err := NewAuteur()
		assert.NoError(t, err)

		assert.False(t, site.HasContent())

		// Empty content should be ignored
		site.AddContent(MockContent{path: []string{}, len: 0})
		assert.False(t, site.HasContent())

		// Add valid content
		site.AddContent(MockContent{path: []string{"path"}, len: 1})
		assert.True(t, site.HasContent())

		// Content in child should affect parent's HasContent
		child := site.GetSubpage("Child", 0)
		assert.False(t, child.HasContent())
		child.AddContent(MockContent{path: []string{"child", "path"}, len: 1})
		assert.True(t, child.HasContent())
		assert.True(t, site.HasContent())
	})

	t.Run("Path Creation", func(t *testing.T) {
		site, err := NewAuteur()
		assert.NoError(t, err)

		content := MockContent{
			path: []string{"level1", "level2", "level3"},
			len:  1,
		}
		site.AddContent(content)

		level1 := site.GetSubpage("level1", 0)
		assert.NotNil(t, level1)

		level2 := level1.GetSubpage("level2", 0)
		assert.NotNil(t, level2)

		level3 := level2.GetSubpage("level3", 0)
		assert.NotNil(t, level3)

		assert.Equal(t, "/level1/level2/level3", level3.Href())
	})

	t.Run("Empty Path Components", func(t *testing.T) {
		site, err := NewAuteur()
		assert.NoError(t, err)

		content := MockContent{
			path: []string{"", "level1", "", "level2", " "},
			len:  1,
		}
		site.AddContent(content)

		level1 := site.GetSubpage("level1", 0)
		assert.NotNil(t, level1)

		level2 := level1.GetSubpage("level2", 0)
		assert.NotNil(t, level2)
	})

	t.Run("Processor Registration", func(t *testing.T) {
		site, err := NewAuteur()
		assert.NoError(t, err)

		processor := MockProcessor{supportedExt: ".txt"}

		site.RegisterProcessor(processor)
		assert.Len(t, site.processors, 1)
	})
}
