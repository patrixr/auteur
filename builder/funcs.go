package builder

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/patrixr/q"
)

var templateFuncs = template.FuncMap{
	"Join": func(base string, elem ...string) string {
		res, err := url.JoinPath(base, elem...)
		q.AssertNoError(err)
		return res
	},
	"Vendor": func(webroot string, distfolder string, url string) string {
		resp, err := http.Get(url)
		q.AssertNoError(err)
		defer resp.Body.Close()

		out, err := os.Create(filepath.Join(distfolder, filepath.Base(url)))
		q.AssertNoError(err)
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		q.AssertNoError(err)

		return filepath.Join(webroot, filepath.Base(url))
	},
	"CleanTitle": func(txt string) string {
		return regexp.MustCompile("[-_]+").ReplaceAllString(txt, " ")
	},
}
