# Cuecard

Cross-platform GUI prompt manager for AI workflows. Point, click, prompt — brain cells optional.

## Features

- Card-based UI for visual prompt discovery
- Groups with colored headers for organization
- Favorites pinned to top for quick access
- Search filtering by title, description, tags
- Variable substitution: `${INPUT}`, `${DATE}`, `${CLIPBOARD}`, `${FILE}`
- System tray for always-available access
- Global hotkey to summon window
- First-run wizard for easy setup
- File watching for live prompt updates

## Installation

TODO: Installation instructions will be added during implementation.

```bash
# Future Homebrew installation
brew install grantcarthew/tap/cuecard
```

## Usage

### Prompt Files

Prompts are markdown files with YAML frontmatter:

```markdown
---
title: Refactor Code
description: Improve code structure
group: Coding
tags: [refactor, cleanup]
input: required
input_hint: file path or component name
---

Review and refactor the following code for better readability:

${INPUT}
```

### Frontmatter Fields

| Field       | Required | Description                           |
| ----------- | -------- | ------------------------------------- |
| title       | Yes      | Display name on card                  |
| description | No       | Tooltip text on hover                 |
| group       | No       | Category for visual grouping          |
| tags        | No       | Keywords for search filtering         |
| input       | No       | "required" or "optional" for ${INPUT} |
| input_hint  | No       | Placeholder text for input field      |
| favorite    | No       | Pin to top of window (true/false)     |

### Variables

| Variable       | Description                              |
| -------------- | ---------------------------------------- |
| `${INPUT}`     | User-provided text from card input field |
| `${DATE}`      | Current date (ISO format)                |
| `${DATETIME}`  | Current date and time                    |
| `${CLIPBOARD}` | Current clipboard content                |
| `${FILE}`      | Opens file picker, inserts selected path |

## Configuration

Config location: `~/.config/cuecard/config.cue`

```cue
prompts_dir: "/path/to/prompts"
editor: "code"
theme: "system"
```

## Development

### Prerequisites

- Go 1.21+
- Fyne dependencies (see [Fyne Getting Started](https://docs.fyne.io/started/))

### Build

```bash
go build -o cuecard ./cmd/cuecard
```

### Run

```bash
./cuecard
```

## Project Structure

This project uses Documentation Driven Development (DDD).

- Design decisions: `.ai/design/design-records/`
- Project tracking: `.ai/projects/`
- Development workflow: `.ai/workflow.md`

## License

MIT License - See [LICENSE](LICENSE) for details.
