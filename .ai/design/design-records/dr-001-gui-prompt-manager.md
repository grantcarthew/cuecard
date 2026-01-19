# dr-001: GUI Prompt Manager

- Date: 2026-01-19
- Status: Accepted
- Category: architecture/gui

## Problem

Managing AI prompts via CLI requires remembering aliases or navigating fzf each time. For frequent prompt usage, a visual interface would be faster and more discoverable. The existing bash-based ai script works but lacks:

- Visual prompt discovery without terminal
- Quick access via system tray
- Rich editing experience with markdown preview
- Cross-platform desktop integration

## Decision

Build Cuecard as a cross-platform GUI application using Go and the Fyne toolkit. The application displays prompts as interactive cards, supports grouping and search, and copies selected prompts to the system clipboard.

Key design elements:

Technology stack:

- Language: Go
- GUI: Fyne toolkit
- Config: CUE format
- Prompts: Markdown files with YAML frontmatter

Core functionality:

- Display prompts as cards in a responsive grid/flow layout
- Group cards by optional group field with colored section headers
- Copy prompt content to clipboard on button click
- Support variable substitution in prompts
- Watch filesystem for prompt file changes

## Why

Go with Fyne:

- Single binary distribution, easy to package with Homebrew
- Cross-platform (Linux, macOS) from single codebase
- No runtime dependencies
- Modern look without heavy frameworks like Electron

Card-based UI:

- Visual scanning faster than reading text lists
- Familiar metaphor (like cue cards for presenters)
- Natural grouping with colored backgrounds
- Input fields can be embedded directly in cards

CUE for configuration:

- Superset of JSON with validation capabilities
- Supports comments and constraints
- Aligns with user preference and existing tooling

Frontmatter in markdown:

- Human-readable and editable
- Version control friendly
- No database required
- Same format as the original ai script (migration path)

## Trade-offs

Accept:

- Fyne has limited markdown rendering (may need custom widget for preview)
- Desktop app requires more resources than CLI
- GUI development more complex than scripts
- File watching adds background processing

Gain:

- Visual discovery without remembering aliases
- One-click prompt copying
- System tray for always-available access
- Rich editing with markdown preview
- Cross-platform from single codebase

## Alternatives

Electron:

- Pro: Full web technologies, rich ecosystem
- Pro: Excellent markdown support
- Con: Large binary size (100MB+)
- Con: High memory usage
- Rejected: Too heavy for a utility app

Tauri:

- Pro: Small binaries, Rust backend
- Pro: Web UI flexibility
- Con: Requires Rust toolchain
- Con: More complex build setup
- Rejected: Go is simpler for this use case

Python + Qt:

- Pro: Mature GUI framework
- Pro: Good cross-platform support
- Con: Requires Python runtime
- Con: Packaging complexity
- Rejected: Go single binary is cleaner

Web app (localhost):

- Pro: Simple development
- Con: Requires browser, less integrated
- Con: No system tray
- Rejected: Desktop integration is important

## Structure

### Frontmatter Schema

Prompt files use YAML frontmatter:

```yaml
---
title: Prompt Title
description: Brief explanation of what this prompt does
group: Category Name
tags: [tag1, tag2, tag3]
alias: shortname
input: optional
input_hint: description of expected input
favorite: false
---

Prompt content here.

${INPUT}
```

Field definitions:

title (string, required):

- Display name shown on card
- Used for search matching
- Used to auto-generate filename

description (string, optional):

- Shown as tooltip on hover
- Used for search matching

group (string, optional):

- Visual grouping with colored header
- Groups sorted alphabetically
- Prompts without group appear last, no header, no background color

tags (array of strings, optional):

- Keywords for search filtering
- Not displayed on cards

alias (string, optional):

- Reserved for future CLI interface
- Must be unique if provided

input (enum, optional):

- Values: "required", "optional"
- Omit if prompt has no variables
- "required": card shows text field, must have value before copy
- "optional": card shows text field, can be empty

input_hint (string, optional):

- Placeholder text for input field
- Shown when input field is empty

favorite (boolean, optional):

- Marks prompt as favorite
- Favorited prompts appear at top of window
- Default: false
- Cleared/ignored during import

### Variable Substitution

Supported variables in prompt content:

| Variable | Description |
| --- | --- |
| ${INPUT} | User-provided text from card input field |
| ${DATE} | Current date (ISO format) |
| ${DATETIME} | Current date and time |
| ${CLIPBOARD} | Current clipboard content |
| ${FILE} | Opens file picker, inserts selected path |

Variables are substituted when copying to clipboard.

### Configuration

Location: ~/.config/cuecard/config.cue

```cue
prompts_dir: "/home/user/ai-prompts"
editor: "code"
theme: "system"
window: {
    width: 1024
    height: 768
    position: "remember"
}
```

Fields:

prompts_dir (string, required):

- Path to directory containing prompt markdown files
- Set during first-run wizard

editor (string, required):

- Command to open external editor
- Examples: "code", "nvim", "vim", "subl"
- Set during first-run wizard

theme (string, optional):

- Values: "light", "dark", "system"
- Default: "system"

window (object, optional):

- width: initial window width in pixels
- height: initial window height in pixels
- position: "remember" to restore last position, "center" to center on screen

### File Naming

Filename auto-generated from title:

- Convert to lowercase
- Replace spaces with hyphens
- Remove special characters
- Add index suffix for duplicates

Examples:

- "Prompt Start" creates prompt-start.md
- Second "Prompt Start" creates prompt-start-2.md

README.md files are ignored when scanning prompts directory.

## GUI Layout

### Main Window

```
+---------------------------------------------------------------+
| Cuecard                                    [search] [toggles] |
+---------------------------------------------------------------+
| -- Favorites ------------------------------------------------ |
| [Card] [Card] [Card]                                          |
|                                                               |
| -- Coding (colored background) ------------------------------ |
| [Card] [Card with input field]                                |
|                                                               |
| -- Project (colored background) ----------------------------- |
| [Card]                                                        |
|                                                               |
| [Card] [Card]  (ungrouped, no header, no color)               |
+---------------------------------------------------------------+
```

### Card Types

Simple card (no input):

```
+-------------------+
| [star] Title      |
|                   |
|     [Copy]        |
+-------------------+
```

Card with input:

```
+-------------------+
| [star] Title      |
| [input field    ] |
|     [Copy]        |
+-------------------+
```

### UI Elements

Search bar:

- Hidden when all cards fit in window
- Appears when cards overflow
- Filters by title, description, tags, group
- Instant filtering as user types
- Keyboard shortcut: / or Cmd+F

View toggles:

- Card view / List view toggle
- Compact mode toggle
- Always on top toggle
- Theme toggle (light/dark/system, minimal UI footprint)

Group headers:

- Collapsible sections
- Background color deterministically generated from group name hash
- Light pastel colors for readability
- Click header to filter to that group only (click again to show all)

Favorites section:

- Star icon on each card to toggle favorite
- Favorited cards appear at top before groups
- Favorites sorted alphabetically

### Menus

File menu:

- New: create new prompt with full form
- Import File...: import single markdown file
- Import Directory...: bulk import files with valid frontmatter
- Open Folder: open prompts directory in file manager
- Refresh: reload prompts from disk
- Quit: exit application

Help menu:

- Validate Prompts: check all prompts for schema issues
- About: version and credits

Card context menu (right-click):

- Preview: view full prompt content (markdown rendered if supported)
- Edit: open in configured editor
- Duplicate: copy file with new name
- Delete: remove with confirmation dialog
- Reveal: show in file manager

### System Integration

System tray:

- App runs in background when closed
- Tray icon for quick access
- Minimize to tray on close

Startup launch:

- Option to launch on system boot

Window memory:

- Remember size and position per monitor

### Feedback

Copy action:

- Brief card highlight animation
- Visual confirmation of success

Delete action:

- Confirmation dialog "Are you sure?"
- Card removed from view on confirm

### First-Run Wizard

Shown when no config exists:

1. Welcome screen with app description
2. Select prompts directory (browse or create default)
3. Select editor (dropdown with common options + custom)
4. Create sample prompt (optional)
5. Save config, launch main window

### Keyboard Navigation

- Tab: move between cards
- Enter: copy focused card
- /: focus search field
- Escape: clear search, unfocus

## Validation

Doctor/Validate function checks:

- Valid YAML frontmatter syntax
- title field present
- input field valid enum if present
- tags is array if present
- group is string if present
- No duplicate titles (warning)
- Files parseable (error if not)

Results shown in dialog:

```
Validation Results
------------------
12 prompts valid
2 warnings
  - prompt-start.md: missing description
  - old-thing.md: duplicate title
1 error
  - broken.md: invalid YAML frontmatter
```

## Import Behavior

Import File:

- Select single markdown file
- If valid frontmatter exists: copy directly to prompts dir
- If no frontmatter: show wizard form to add metadata
- Filename preserved or regenerated from title
- favorite field cleared/ignored during import

Import Directory:

- Select directory
- Scan all .md files
- Import files with valid frontmatter schema
- Skip files without frontmatter or invalid schema
- Skip README.md
- favorite field cleared/ignored during import
- Show summary dialog with counts and skipped files

## Updates

- 2026-01-19: Initial design
