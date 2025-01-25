package core

type AuteurConfig struct {
	Exclude   []string `yaml:"exclude"`
	Title     string   `yaml:"title"`
	Desc      string   `yaml:"desc"`
	Version   string   `yaml:"version"`
	Outfolder string   `yaml:"outfolder"`
	Rootdir   string   `yaml:"root"`
	Webroot   string   `yaml:"webroot"`
}

func (ac AuteurConfig) ExtendConfig(other *AuteurConfig) AuteurConfig {
	if other == nil {
		return ac
	}

	if other.Exclude != nil {
		ac.Exclude = append(ac.Exclude, other.Exclude...)
	}

	if other.Title != "" {
		ac.Title = other.Title
	}

	if other.Desc != "" {
		ac.Desc = other.Desc
	}

	if other.Version != "" {
		ac.Version = other.Version
	}

	if other.Outfolder != "" {
		ac.Outfolder = other.Outfolder
	}

	if other.Rootdir != "" {
		ac.Rootdir = other.Rootdir
	}

	if other.Webroot != "" {
		ac.Webroot = other.Webroot
	}

	return ac
}
