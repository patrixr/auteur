/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/patrixr/auteur/builder/basic"
	. "github.com/patrixr/auteur/common"
	. "github.com/patrixr/auteur/core"
	"github.com/spf13/cobra"
)

/*
---
auteur: /
---

# Auteur

Auteur is a static site generator, originally designed to generate documentation for software projects.
It traverses the file structure of a given folder and processes files to find content that can be rendered into a static site.

For code documentation, Auteur looks for comments in the source code files. Here's an example:

```go
func hello() {
world()
}
```
*/

var rootCmd = &cobra.Command{
	Use:   "auteur",
	Short: "A static site generator",
	Long:  `A static site generator`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		folder := "."
		outfolder := "out"

		if len(args) > 0 {
			folder = args[0]
		}

		Log("Booting Auteur", "folder", folder)

		auteur := NewSite()
		auteur.RegisterProcessor(NewCommentReader(folder))

		folder, err := filepath.Abs(folder)

		if err != nil {
			LogError(err)
			os.Exit(1)
		}

		if err := auteur.Ingest(folder); err != nil {
			LogError(err)
			os.Exit(1)
		}

		if auteur.HasContent() == false {
			LogErrorf("No Auteur-compatible content found in folder %s", folder)
			os.Exit(1)
		}

		builder := basic.NewBasicBuilder()

		if err := builder.Render(auteur, outfolder); err != nil {
			LogError(err)
			os.Exit(1)
		}

		Log("Auteur completed successfully", "out", outfolder)
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
