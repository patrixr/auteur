package core

import (
	"os"
	"path/filepath"

	. "github.com/patrixr/auteur/common"
	"github.com/patrixr/q"
	"gopkg.in/yaml.v3"
)

type Link struct {
	Title string `yaml:"title"`
	Href  string `yaml:"url"`
	Icon  string `yaml:"icon"`
}

type AuteurConfig struct {
	Exclude   []string `yaml:"exclude"`
	Title     string   `yaml:"title"`
	Desc      string   `yaml:"desc"`
	Version   string   `yaml:"version"`
	Outfolder string   `yaml:"outfolder"`
	Rootdir   string   `yaml:"root"`
	Webroot   string   `yaml:"webroot"`
	Links     []Link   `yaml:"links"`
	Order     int      `yaml:"order"`
	Theme     string   `yaml:"theme"`
}

// ExtendConfig returns a new AuteurConfig with the values of the other config
// merged into the current one. The other config takes precedence over the current one.
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

	if other.Theme != "" {
		ac.Theme = other.Theme
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

	if other.Order != ac.Order {
		ac.Order = other.Order
	}

	return ac
}

// DetectConfig reads the configuration file from the current directory and returns
// an AuteurConfig struct with the values from the configuration file.
// Environment variables can be used to override the values in the configuration file.
func DetectConfig() (AuteurConfig, error) {
	config := AuteurConfig{
		Title:     "Auteur",
		Desc:      "Static site generated with Auteur",
		Rootdir:   ".",
		Outfolder: "out",
		Webroot:   "/",
		Version:   "0.0.1",
		Theme:     "default",
		Exclude: []string{
			"node_modules",
			".git",
			".gitignore",
			".DS_Store",
			"*_test.go",
		},
	}

	candidates := []string{
		"auteur.yml",
		"auteur.yaml",
		"auteur.json",
	}

	for _, configFile := range candidates {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			continue
		}

		Log("Configuration detected", "file", configFile)
		fileContent, err := os.ReadFile(configFile)

		if err != nil && os.IsNotExist(err) {
			continue
		}

		if err != nil {
			return config, err
		}

		if err := yaml.Unmarshal(fileContent, &config); err != nil {
			return config, err
		}
	}

	config.Webroot = q.ReadEnv("AUTEUR_WEBROOT", config.Webroot)
	config.Outfolder = q.ReadEnv("AUTEUR_OUTFOLDER", config.Outfolder)
	config.Rootdir = q.ReadEnv("AUTEUR_ROOTDIR", config.Rootdir)
	config.Title = q.ReadEnv("AUTEUR_TITLE", config.Title)
	config.Version = q.ReadEnv("AUTEUR_VERSION", config.Version)
	config.Desc = q.ReadEnv("AUTEUR_DESC", config.Desc)

	absRootdir, err := filepath.Abs(config.Rootdir)
	if err != nil {
		return config, err
	}
	config.Rootdir = absRootdir

	absOutfolder, err := filepath.Abs(config.Outfolder)
	if err != nil {
		return config, err
	}
	config.Outfolder = absOutfolder

	return config, nil
}
