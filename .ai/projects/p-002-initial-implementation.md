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
- Global hotkey to summon window (deferred to future release)
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
- [ ] ~~Global hotkey summons window~~ (deferred)
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

## Current State

Codebase:

- No Go source files exist yet (greenfield implementation)
- No go.mod/go.sum initialized
- Project structure not created (cmd/, internal/ directories don't exist)
- Design documentation complete (DR-001 accepted)
- README.md with user documentation drafted

Development environment:

- Go 1.25.6 installed (darwin/arm64)
- macOS development target confirmed
- Linux build target planned

Verified dependency versions (latest stable):

| Package | Version | Purpose |
| --- | --- | --- |
| fyne.io/fyne/v2 | v2.7.2 | GUI toolkit |
| cuelang.org/go | v0.15.3 | CUE configuration parsing |
| gopkg.in/yaml.v3 | v3.0.1 | YAML frontmatter parsing |
| github.com/fsnotify/fsnotify | v1.9.0 | File system watching |
| github.com/golang-design/hotkey | v0.4.1 | Global hotkey support (deferred) |

Implementation notes:

- `.gitignore` needs update to include `cuecard` binary (build output)
- Go module path follows convention: `github.com/grantcarthew/cuecard`

Research findings:

Fyne v2.7.2:

- System tray: Built-in support, works well on macOS, Linux requires AppIndicator
- GridWrap container: Available for responsive card layouts
- RichText: Limited markdown (bold, italic, headers, links only - no tables or code blocks)
- Requires CGO and C compiler for builds
- Linux needs: gcc, libgl1-mesa-dev, xorg-dev packages

CUE v0.15.3:

- Uses cuecontext.New() API (Runtime deprecated)
- Context instances grow in memory; recreate periodically
- Decode() method maps CUE to Go structs via json tags

golang-design/hotkey v0.4.1:

- Cross-platform global hotkey support (macOS, Linux, Windows)
- Actively maintained
- Lightweight compared to alternatives like robotn/gohook

## Decision Points

1. Fyne flow layout implementation

Decision: **A - Use container.NewGridWrap for responsive cards**

2. Markdown preview rendering

Decision: **A - Use Fyne's built-in rich text**

3. Global hotkey implementation

Decision: **C - Defer to future release** (not included in v1)
