---
order: 3
---

# Supported Languages

Auteur supports documentation extraction from multiple programming languages.
Each language has its own comment syntax, but all follow the same principle: use the `@auteur` marker inside a comment to have auteur ingest it as markdown.
Here's how to write documentation in each supported language.

## C-Style Languages (Go, JavaScript, Java, C++, etc.)

Use either line comments with `//` or block comments with `/* */`:

```c
// @auteur
// # User Authentication
//
// This module handles user authentication and sessions.
//
// ## Usage
//
// Call `authenticate()` with user credentials.

/* @auteur
# User Authentication

This module handles user authentication and sessions.

## Usage

Call `authenticate()` with user credentials.
*/
```

## Python

Use either line comments with `#` or triple quotes:

```python
# @auteur
# # User Authentication
#
# This module handles user authentication and sessions.
#
# ## Usage
#
# Call `authenticate()` with user credentials.

"""
@auteur
# User Authentication

This module handles user authentication and sessions.

## Usage

Call `authenticate()` with user credentials.
"""
```

## Ruby

Use either line comments with `#` or block comments with `=begin/=end`:

```ruby
# @auteur
# # User Authentication
#
# This module handles user authentication and sessions.

=begin @auteur
# User Authentication

This module handles user authentication and sessions.

## Usage

Call `authenticate()` with user credentials.
=end
```

## Lua

Use block comments with `--[[` and `]]`:

```lua
--[[
  @auteur
  # User Authentication

  This module handles user authentication and sessions.

  ## Usage

  Call `authenticate()` with user credentials.
]]
```

## HTML/XML

Use HTML-style comments:

```html
<!-- @auteur
# Component Documentation

This component handles user input.

## Props

- `value`: Current input value
- `onChange`: Change handler
-->
```

## SQL

Use line comments with `--`:

```sql
-- @auteur
-- # User Table Schema
--
-- Defines the structure for storing user data.
--
-- ## Columns
--
-- - `id`: Primary key
-- - `username`: Unique identifier
```

## Common Features

For all languages:

1. Start with `@auteur` to mark documentation blocks
2. Use standard Markdown syntax after the marker
3. Maintain consistent indentation
4. Add blank lines between sections using comment markers
5. Headers, lists, code blocks, and other Markdown features are supported

## Markdown Support

You can use:

- Headers (`#`, `##`, etc.)
- Lists (ordered and unordered)
- Code blocks (fenced or indented)
- Links and images
- Tables
- Blockquotes
- All other standard Markdown syntax

Remember to maintain proper comment markers for your chosen language while writing Markdown content.
