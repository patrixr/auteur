/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/patrixr/auteur/builder"
	. "github.com/patrixr/auteur/common"
	. "github.com/patrixr/auteur/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "auteur",
	Short: "A static site generator",
	Long:  `A static site generator`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		config := AuteurConfig{
			Title:     "Auteur",
			Desc:      "Static site generated with Auteur",
			Rootdir:   ".",
			Outfolder: "out",
			Exclude: []string{
				"node_modules",
				".git",
				".gitignore",
				".DS_Store",
				"*_test.go",
			},
		}

		configFiles := []string{
			"auteur.yml",
			"auteur.yaml",
			"auteur.json",
		}

		for _, configFile := range configFiles {
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				continue
			}

			Log("Configuration detected", "file", configFile)
			fileContent, err := os.ReadFile(configFile)

			if err != nil {
				LogError(err)
				os.Exit(1)
			}

			if err := yaml.Unmarshal(fileContent, &config); err != nil {
				LogError(err)
				os.Exit(1)
			}
		}

		Log("Configuration", "config", fmt.Sprintf("%+v", config))

		folder, err := filepath.Abs(config.Rootdir)
		if err != nil {
			LogError(err)
			os.Exit(1)
		}

		Log("Booting Auteur", "root", config.Rootdir)

		auteur := NewSiteWithConfig(config)
		auteur.RegisterProcessor(NewCommentReader())
		auteur.RegisterProcessor(NewMarkdownProcessor())

		if err := auteur.Ingest(config.Rootdir); err != nil {
			LogError(err)
			os.Exit(1)
		}

		if auteur.HasContent() == false {
			LogErrorf("No Auteur-compatible content found in folder %s", folder)
			os.Exit(1)
		}

		builder := builder.NewDefaultBuilder()

		if err := builder.Render(auteur, config.Outfolder); err != nil {
			LogError(err)
			os.Exit(1)
		}

		Log("Auteur completed successfully", "out", config.Outfolder)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.auteur.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//
	// e.g.:
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
