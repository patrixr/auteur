/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/patrixr/auteur/builder"
	. "github.com/patrixr/auteur/common"
	. "github.com/patrixr/auteur/core"
	. "github.com/patrixr/auteur/processors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auteur",
	Short: "A static site generator",
	Long:  `A static site generator`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		auteur, err := NewAuteur()
		if err != nil {
			LogError(err)
			os.Exit(1)
		}

		Log("Booting Auteur", "root", auteur.Rootdir)

		auteur.RegisterProcessor(NewCommentReader())
		auteur.RegisterProcessor(NewMarkdownProcessor())

		if err := auteur.Ingest(auteur.Rootdir); err != nil {
			LogError(err)
			os.Exit(1)
		}

		if auteur.HasContent() == false {
			LogErrorf("No Auteur-compatible content found in folder %s", auteur.Rootdir)
			os.Exit(1)
		}

		builder := builder.NewDefaultBuilder()

		if err := builder.Render(auteur, auteur.Outfolder); err != nil {
			LogError(err)
			os.Exit(1)
		}

		Log("Auteur completed successfully", "out", auteur.Outfolder)
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
