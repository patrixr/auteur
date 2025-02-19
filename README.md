# Auteur

Auteur is a simple static site generator, designed to generate a static website from an arbitrary folder structure.
Its main use case is building out documentation from a codebase.

Auteur is built using:

- [Go](https://golang.org)
- [Web Awesome](https://webawesome.com)

> Note: Both Auteur and Webawesome are in their alpha stages, changes are expected

## Getting Started

Running Auteur can be done via the go command line tool, it is as simple as running:

```sh
go run github.com/patrixr/auteur
```

Auteur will look for Markdown files and code comments in the current directory and generate a static website from them.

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

## Release process

Currently auteur is available via:

- `go get` and `go run` commands
- Homebrew

To make a release, you can use the `make tag release` command. This command automates the process of tagging a new release version in your repository and pushing it to github.

Requirements:

- `git`
- `gh` github command line tool

Before running the command, ensure that you have the necessary permissions to push tags to the remote repository.
