package builder

import (
	. "github.com/patrixr/auteur/core"
)

type Builder interface {
	Render(site *Site, outfolder string) error
}
