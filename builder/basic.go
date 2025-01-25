package builder

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	. "github.com/patrixr/auteur/common"
	. "github.com/patrixr/auteur/core"
)

//go:embed assets/*
var tmplFS embed.FS
var templates = template.Must(
	template.New("").Funcs(template.FuncMap{}).ParseFS(tmplFS, "assets/**/*.tmpl"),
)

type DefaultBuilder struct{}

func NewDefaultBuilder() Builder {
	return DefaultBuilder{}
}

// Render generates the site files and folders inside the output folder
// This is the main entry point for the builder
func (builder DefaultBuilder) Render(site *Auteur, outfolder string) error {
	pageKey := site.Slug()

	if site.IsRoot() {
		site.PrettyPrint()
		pageKey = "index"
		if err := Rmdir(outfolder); err != nil {
			return err
		}
	} else if site.HasChildren() {
		// We create a folder and use an index.html file
		// to accomodate having children files
		outfolder = filepath.Join(outfolder, site.Slug())
		pageKey = "index"
	}

	if err := Mkdirp(outfolder); err != nil {
		return err
	}

	if len(site.Content) > 0 || site.IsRoot() {
		fileName := fmt.Sprintf("%s.html", pageKey)
		fragFileName := fmt.Sprintf("%s.frag.html", pageKey)

		// Write page file
		file, err := os.Create(filepath.Join(outfolder, fileName))

		if err != nil {
			return err
		}

		html, err := builder.GetHTML(site)

		if err != nil {
			return err
		}

		err = templates.ExecuteTemplate(file, "page.html.tmpl", struct {
			Fragment string
			Site     *Auteur
			Title    string
			Webroot  string
		}{
			Fragment: html.String(),
			Site:     site.Root(),
			Title:    site.Title,
			Webroot:  strings.TrimRight(site.Webroot, "/"),
		})

		// Close manually (instead of defer) to avoid stacking up open files
		file.Close()

		if err != nil {
			return err
		}

		// Create frag file
		if err := os.WriteFile(filepath.Join(outfolder, fragFileName), html.Bytes(), 0644); err != nil {
			return err
		}
	}

	for _, child := range site.Children() {
		if err := builder.Render(child, outfolder); err != nil {
			return err
		}
	}

	if site.IsRoot() {
		if err := builder.CopyAssets(outfolder); err != nil {
			return err
		}
	}

	return nil
}

func (t DefaultBuilder) GetHTML(site *Auteur) (bytes.Buffer, error) {
	buffer := bytes.Buffer{}

	for _, content := range site.Content {
		switch content.Type() {
		case Markdown:
			if err := ConvertMarkdown([]byte(content.Data()), &buffer); err != nil {
				return buffer, err
			}
		case HTML:
			if _, err := buffer.WriteString(content.Data()); err != nil {
				return buffer, err
			}
		default:
			return buffer, fmt.Errorf("Unknown content type: %d", content.Type())
		}
	}

	return buffer, nil
}

func (t DefaultBuilder) CopyAssets(outfolder string) error {
	filesToCopy := []string{"assets/default/script.js", "assets/default/style.css"}

	for _, file := range filesToCopy {
		src, err := tmplFS.Open(file)
		if err != nil {
			return err
		}
		defer src.Close()

		destPath := filepath.Join(outfolder, filepath.Base(file))
		dest, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		if _, err := io.Copy(dest, src); err != nil {
			return err
		}
	}
	return nil
}
