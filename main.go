/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/patrixr/auteur/cmd"

/*
 * ---
 * auteur: /
 * ---
 *
 * # Auteur
 *
 * Auteur is a static site generator, originally designed to generate documentation for software projects.
 * It traverses the file structure of a given folder and processes files to find content that can be rendered into a static site.
 *
 * For code documentation, Auteur looks for comments in the source code files. Here's an example:
 *
 * ```markdown
 * ---
 * auteur: /path
 * ---
 *
 * # Lorem ipsum dolor
 *
 * Consectetur adipiscing elit. Integer nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla.
 * Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero. Sed dignissim lacinia nunc. Curabitur tortor. Pellentesque nibh. Aenean quam. In scelerisque sem at dolor. Maecenas mattis. Sed conv
 * ```
 *
 */

func main() {
	cmd.Execute()
}
