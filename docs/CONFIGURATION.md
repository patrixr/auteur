---
order: 2
---

# Auteur Configuration

Auteur is designed to be flexible and configurable. You can customize the behavior of the tool by providing a configuration file in YAML or JSON format.

The following file names are recognized:

- `auteur.yaml`
- `auteur.yml`
- `auteur.json`

## Configuration File

The configuration is specified in YAML format, typically named `config.yaml` or similar.

## Basic Settings

| Setting     | Type   | Description                                    | Default |
| ----------- | ------ | ---------------------------------------------- | ------- |
| `title`     | string | Project name displayed in documentation        | Auteur  |
| `version`   | string | Version number of the project                  | 0.0.1   |
| `outfolder` | string | Output directory for generated files           | ./dist  |
| `root`      | string | Root directory containing source documentation | .       |
| `webroot`   | string | Base URL path for web serving                  | /       |

## Exclusion Rules

The `exclude` section defines patterns and directories to ignore during processing:

```yml
exclude:
  - .git
  - tmp
  - dist
  - node_modules
  - "*_test.go"
```

Files and directories matching these patterns will be skipped during document generation. The exclusion supports:

- Directory names
- File patterns using glob syntax
- Hidden files and directories

## Example Configuration

```yml
title: "Auteur"
version: 0.0.6
outfolder: ./dist
root: "./docs"
webroot: "/auteur"
exclude:
  - .git
  - tmp
  - dist
  - node_modules
  - "*_test.go"
```
