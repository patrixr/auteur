package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/patrixr/auteur/common"
)

type Auteur struct {
	AuteurConfig

	Content    []Content
	parent     *Auteur
	root       *Auteur
	children   []*Auteur
	processors []Processor
}

func NewAuteur() *Auteur {
	return NewSiteWithConfig(AuteurConfig{
		Title:   "Auteur",
		Desc:    "Static site generated with Auteur",
		Exclude: []string{"node_modules", ".git", ".gitignore", ".DS_Store"},
	})
}

func NewSiteWithConfig(config AuteurConfig) *Auteur {
	if len(config.Title) == 0 {
		config.Title = "Auteur"
	}

	return &Auteur{
		AuteurConfig: config,
		parent:       nil,
		root:         nil,
		Content:      []Content{},
		processors:   []Processor{},
	}
}

func (site *Auteur) SetTitle(title string) {
	site.Title = title
}

func (site *Auteur) SetDesc(desc string) {
	site.Desc = desc
}

func (site *Auteur) Slug() string {
	return common.ToSlug(site.Title)
}

func (site *Auteur) Root() *Auteur {
	ref := site

	for !ref.IsRoot() {
		ref = ref.parent
	}

	return ref
}

func (site *Auteur) Href() string {
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

func (site *Auteur) IsRoot() bool {
	return site.root == nil
}

func (site *Auteur) HasChildren() bool {
	return len(site.children) > 0
}

func (site *Auteur) Children() []*Auteur {
	return site.children
}

// Registers a processor to be used when ingesting files
// Each processor is responsible for transforming specific file types
func (site *Auteur) RegisterProcessor(processor Processor) {
	site.processors = append(site.processors, processor)
}

// Given a folder, this function ingests all files and directories within it
// using the registered processors to transform files into site content
func (site *Auteur) Ingest(infolder string) error {
	files, err := os.ReadDir(infolder)

	if err != nil {
		return err
	}

	for _, file := range files {
		abspath, err := filepath.Abs(filepath.Join(infolder, file.Name()))

		if err != nil {
			return err
		}

		fmt.Println(abspath)

		if IsExcluded(file.Name(), site.Exclude) {
			common.Log("Excluding " + abspath)
			continue
		}

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

			contents, err := processor.Load(site, abspath)

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
func (site *Auteur) AddContent(content Content) {
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

	ordered := make([]Content, len(ref.Content)+1)

	for i := 0; i < len(ordered); i++ {
		if i == len(ordered)-1 {
			ordered[i] = content
			break
		}

		existing := ref.Content[i]

		if existing.Order() <= content.Order() {
			ordered[i] = existing
			continue
		}

		ordered[i] = content
		ordered[i+1] = existing
		i += 1
	}

	ref.Content = ordered
}

// GetSubpage retrieves a subpage with the given title. If the subpage does not exist, it creates a new one.
// The title comparison is case-insensitive and ignores leading and trailing whitespace.
func (site *Auteur) GetSubpage(title string) *Auteur {
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

	newPage := &Auteur{
		AuteurConfig: site.ExtendConfig(&AuteurConfig{
			Title: title,
		}),
		root:   root,
		parent: site,
	}
	site.children = append(site.children, newPage)
	return newPage
}

func (site *Auteur) PrettyPrint() {
	var traverse func(child *Auteur, level int) string

	traverse = func(child *Auteur, level int) string {
		var sb strings.Builder
		indent := strings.Repeat("  ", level)
		sb.WriteString(fmt.Sprintf("%s- %s\n", indent, child.Title))
		sb.WriteString(fmt.Sprintf("%s  href=%s\n", indent, child.Href()))
		for _, child := range child.children {
			sb.WriteString(traverse(child, level+1))
		}
		return sb.String()
	}

	fmt.Println(traverse(site, 0))
}

func (site *Auteur) HasContent() bool {
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

func IsExcluded(filename string, patterns []string) bool {
	for _, pattern := range patterns {

		if pattern == filename {
			return true
		}

		pattern = filepath.FromSlash(pattern)

		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
	}
	return false
}
