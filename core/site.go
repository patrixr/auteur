package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/patrixr/auteur/common"
)

type Site struct {
	Content    []Content
	parent     *Site
	root       *Site
	title      string
	desc       string
	children   []*Site
	processors []Processor
}

func NewSite() *Site {
	return &Site{
		parent:     nil,
		root:       nil,
		title:      "Auteur",
		Content:    []Content{},
		processors: []Processor{},
	}
}

func (site *Site) SetTitle(title string) {
	site.title = title
}

func (site *Site) SetDesc(desc string) {
	site.desc = desc
}

func (site *Site) Title() string {
	return site.title
}

func (site *Site) Slug() string {
	return common.ToSlug(site.title)
}

func (site *Site) Desc() string {
	return site.desc
}

func (site *Site) Root() *Site {
	ref := site

	for !ref.IsRoot() {
		ref = ref.parent
	}

	return ref
}

func (site *Site) Href() string {
	if site.IsRoot() {
		return "/"
	}

	var parts []string
	current := site
	for !current.IsRoot() {
		parts = append([]string{current.Slug()}, parts...)
		current = current.parent
	}

	return "/" + strings.Join(parts, "/")
}

func (site *Site) IsRoot() bool {
	return site.root == nil
}

func (site *Site) HasChildren() bool {
	return len(site.children) > 0
}

func (site *Site) Children() []*Site {
	return site.children
}

// Registers a processor to be used when ingesting files
// Each processor is responsible for transforming specific file types
func (site *Site) RegisterProcessor(processor Processor) {
	site.processors = append(site.processors, processor)
}

// Given a folder, this function ingests all files and directories within it
// using the registered processors to transform files into site content
func (site *Site) Ingest(infolder string) error {
	files, err := os.ReadDir(infolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		abspath := filepath.Join(infolder, file.Name())

		// Recurse into directories
		if file.IsDir() {
			if err := site.Ingest(abspath); err != nil {
				return err
			}
			continue
		}

		ext := filepath.Ext(file.Name())

		for _, processor := range site.processors {
			if !processor.Supports(ext) {
				continue
			}

			contents, err := processor.Load(abspath)

			if err != nil {
				return err
			}

			for _, content := range contents {
				site.AddContent(content)
			}
		}
	}
	return nil
}

// Adds HTML/Markdown content to the site
func (site *Site) AddContent(content Content) {
	if content == nil || content.Len() == 0 {
		return
	}

	ref := site

	for _, part := range content.Path() {
		if len(strings.Trim(part, " \t\n")) == 0 {
			continue
		}
		ref = ref.GetSubpage(part)
	}

	ref.Content = append(ref.Content, content)
}

// GetSubpage retrieves a subpage with the given title. If the subpage does not exist, it creates a new one.
// The title comparison is case-insensitive and ignores leading and trailing whitespace.
func (site *Site) GetSubpage(title string) *Site {
	slug := common.ToSlug(title)

	for _, subpage := range site.children {
		if strings.EqualFold(subpage.Slug(), slug) {
			return subpage
		}
	}

	root := site

	if !site.IsRoot() {
		root = site.root
	}

	newPage := &Site{
		title:  title,
		root:   root,
		parent: site,
	}
	site.children = append(site.children, newPage)
	return newPage
}

func (site *Site) PrettyPrint() {
	var traverse func(child *Site, level int) string

	traverse = func(child *Site, level int) string {
		var sb strings.Builder
		indent := strings.Repeat("  ", level)
		sb.WriteString(fmt.Sprintf("%s- %s\n", indent, child.title))
		for _, child := range child.children {
			sb.WriteString(traverse(child, level+1))
		}
		return sb.String()
	}

	fmt.Println(traverse(site, 0))
}

func (site *Site) HasContent() bool {
	if len(site.Content) > 0 {
		return true
	}

	for _, child := range site.children {
		if child.HasContent() {
			return true
		}
	}

	return false
}
