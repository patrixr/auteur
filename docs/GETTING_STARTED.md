---
order: 1
---

# Getting Started

Running Auteur can be done via the go command line tool, it is as simple as running:

```sh
go run github.com/patrixr/auteur
```

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

To register a comment as a page, a comment should contain an `@auteur` tag to indicate that it should be included in the generated website, and it should end with an `@end` tag to indicate the end of the content.

Example:

```go
package main

import "fmt"

// @auteur
// # Welcome to Auteur
// This is a sample page
// @end
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
// @end

func main() {
	fmt.Println("Hello, World!")
}
```
