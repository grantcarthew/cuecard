# p-002: Initial Implementation

- Status: Pending
- Started: -

## Overview

Implement Cuecard v1.0 as a cross-platform GUI prompt manager using Go and Fyne. This project covers the core functionality needed for a usable first release.

## Goals

1. Implement core card-based UI with prompt display
2. Support frontmatter parsing and variable substitution
3. Add grouping, favorites, and search functionality
4. Implement system tray and window management
5. Create first-run wizard and CUE configuration
6. Package for distribution (macOS, Linux)

## Scope

In Scope:

- Card and list view layouts
- Group sections with colored headers
- Favorites with star toggle
- Search bar with filtering
- Collapsible groups
- Compact mode toggle
- Copy to clipboard with variable substitution
- Input fields for prompts requiring user input
- Card context menu (Edit, Duplicate, Delete, Reveal)
- File menu (New, Import File, Import Directory, Open Folder, Refresh, Quit)
- Help menu (Validate Prompts, About)
- System tray with minimize to tray
- Global hotkey to summon window
- Always on top toggle
- Startup launch option
- Window position/size memory
- Keyboard navigation
- First-run wizard
- CUE configuration file
- File watching for live reload
- Markdown preview for prompts (if Fyne supports it easily)

Out of Scope:

- CLI interface (future project)
- Cloud sync
- Plugin system
- Multiple prompt directories

## Success Criteria

- [ ] Application launches and displays prompt cards from configured directory
- [ ] Cards show title, copy button, and input field when needed
- [ ] Groups display with colored headers, collapsible
- [ ] Favorites pin to top of window
- [ ] Search filters cards by title, description, tags, group
- [ ] Copy button substitutes variables and copies to clipboard
- [ ] Card context menu works (Edit, Duplicate, Delete, Reveal)
- [ ] File menu operations work (New, Import, Open Folder, Refresh)
- [ ] Validate Prompts shows validation results dialog
- [ ] System tray icon present, window minimizes to tray
- [ ] Global hotkey summons window
- [ ] First-run wizard creates config on initial launch
- [ ] File changes detected and UI refreshed
- [ ] Application builds and runs on macOS and Linux

## Deliverables

Code:

- cmd/cuecard/main.go - application entry point
- internal/config/ - CUE configuration loading
- internal/prompt/ - prompt file parsing and management
- internal/ui/ - Fyne UI components
- internal/clipboard/ - clipboard operations
- internal/watcher/ - file system watching

Documentation:

- README.md - user documentation
- docs/installation.md - installation guide

Distribution:

- Homebrew formula (or tap)
- Binary releases for macOS and Linux

## Technical Approach

Project structure:

```
cuecard/
  cmd/cuecard/main.go
  internal/
    config/
    prompt/
    ui/
    clipboard/
    watcher/
  go.mod
  go.sum
```

Key dependencies:

- fyne.io/fyne/v2 - GUI toolkit
- cuelang.org/go/cue - CUE configuration parsing
- gopkg.in/yaml.v3 - YAML frontmatter parsing
- fsnotify - file system watching

## Decision Points

1. Fyne flow layout implementation

- A: Use container.NewGridWrap for responsive cards
- B: Custom flow layout widget
- C: Horizontal wrap container

2. Markdown preview rendering

- A: Use Fyne's built-in rich text (limited markdown)
- B: Integrate third-party markdown renderer
- C: Show raw text, skip markdown rendering for v1

3. Global hotkey implementation

- A: Use OS-specific APIs directly
- B: Use go-hotkey library
- C: Defer to future release if complex
