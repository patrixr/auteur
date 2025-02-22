---
priority: 100
---

# Getting Started

1. Install Auteur using Homebrew or Golang. Instructions are available in the [Installation](/installation) section.
2. Naviage to a directory containing markdown files or code comments
3. Run Auteur via a simple command line:

**Important**: Comments should follow the [Auteur syntax](/supported-languages) to be included in the generated website.

```sh
# Homebrew
auteur
# or for Go
go run github.com/patrixr/auteur
```

Auteur is also available as a [Homebrew formule](/installation)

## Markown Pages

Auteur supports markdown pages, which can be used to generate static content as you would a traditional static site generator.

Simply add markdown files to the directory you are running Auteur from, and they will be included in the generated website.

Example file structure:

```
docs/
  index.md
  about.md
  contact.md
```

## Writing Comments

To register a comment as a page, a comment should contain an `@auteur` tag to indicate that it should be included in the generated website.

Example:

```go
package main

import "fmt"

// @auteur
// # Welcome to Auteur
// This is a sample page
func main() {
	fmt.Println("Hello, World!")
}
```

## Specifying the path

To specify the path of the generated page directly within the `@auteur` annotation, you can add a string argument to the annotation.
This string argument should represent the desired path for the page.

Example:

```go
package main

import "fmt"

// @auteur("/my-subpage/my-subsubpage")
// # Custom Path Page
func main() {
 fmt.Println("Hello, World!")
}
```

## Frontmatter

Auteur also supports frontmatter in the form of a YAML object at the beginning of a comment block.
This allows to control:

- The title of the page
- The order in which the content appears
- The path to the page
- Whether the page should be ignored or not

Example:

```go
package main

import "fmt"

// @auteur
// ---
// title: Welcome to Auteur
// order: 1
// path: /welcome
// ---
//
// # Hello
// This is a sample page

func main() {
	fmt.Println("Hello, World!")
}
```
