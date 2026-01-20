package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Prompt represents a parsed prompt file
type Prompt struct {
	// Metadata from frontmatter
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Group       string   `yaml:"group"`
	Tags        []string `yaml:"tags"`
	Alias       string   `yaml:"alias"`
	Input       string   `yaml:"input"` // "required", "optional", or ""
	InputHint   string   `yaml:"input_hint"`
	Favorite    bool     `yaml:"favorite"`

	// Content and file info
	Content  string // The prompt content (after frontmatter)
	FilePath string // Full path to the file
	FileName string // Just the filename
}

// RequiresInput returns true if the prompt requires user input
func (p *Prompt) RequiresInput() bool {
	return p.Input == "required"
}

// HasOptionalInput returns true if the prompt has optional input
func (p *Prompt) HasOptionalInput() bool {
	return p.Input == "optional"
}

// HasInput returns true if the prompt has any input field
func (p *Prompt) HasInput() bool {
	return p.Input == "required" || p.Input == "optional"
}

// HasVariables returns true if the content contains any variables
func (p *Prompt) HasVariables() bool {
	return strings.Contains(p.Content, "${")
}

// LoadDirectory loads all prompts from a directory
func LoadDirectory(dir string) ([]*Prompt, error) {
	var prompts []*Prompt

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Skip non-markdown files
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			continue
		}
		// Skip README files
		if strings.EqualFold(name, "readme.md") {
			continue
		}

		path := filepath.Join(dir, name)
		prompt, err := LoadFile(path)
		if err != nil {
			// Log error but continue loading other prompts
			continue
		}
		prompts = append(prompts, prompt)
	}

	return prompts, nil
}

// LoadFile loads a single prompt from a file
func LoadFile(path string) (*Prompt, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	prompt, err := Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	prompt.FilePath = path
	prompt.FileName = filepath.Base(path)

	return prompt, nil
}

// GroupPrompts organizes prompts by group, with favorites first
func GroupPrompts(prompts []*Prompt) (favorites []*Prompt, groups map[string][]*Prompt, ungrouped []*Prompt) {
	groups = make(map[string][]*Prompt)

	for _, p := range prompts {
		if p.Favorite {
			favorites = append(favorites, p)
		}

		if p.Group != "" {
			groups[p.Group] = append(groups[p.Group], p)
		} else {
			ungrouped = append(ungrouped, p)
		}
	}

	// Sort each list by title
	sortByTitle := func(list []*Prompt) {
		sort.Slice(list, func(i, j int) bool {
			return strings.ToLower(list[i].Title) < strings.ToLower(list[j].Title)
		})
	}

	sortByTitle(favorites)
	sortByTitle(ungrouped)
	for _, group := range groups {
		sortByTitle(group)
	}

	return favorites, groups, ungrouped
}

// SortedGroupNames returns group names sorted alphabetically
func SortedGroupNames(groups map[string][]*Prompt) []string {
	names := make([]string, 0, len(groups))
	for name := range groups {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GenerateFilename creates a filename from a title
func GenerateFilename(title string, existingFiles []string) string {
	// Convert to lowercase
	name := strings.ToLower(title)

	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")

	// Remove special characters (keep alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}
	name = result.String()

	// Remove consecutive hyphens
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	// Trim leading/trailing hyphens
	name = strings.Trim(name, "-")

	if name == "" {
		name = "prompt"
	}

	// Add .md extension
	baseName := name + ".md"

	// Check for duplicates and add suffix if needed
	exists := func(n string) bool {
		for _, f := range existingFiles {
			if strings.EqualFold(f, n) {
				return true
			}
		}
		return false
	}

	if !exists(baseName) {
		return baseName
	}

	// Try with numeric suffix
	for i := 2; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d.md", name, i)
		if !exists(candidate) {
			return candidate
		}
	}

	return baseName // fallback
}

// Filter returns prompts matching the search query
func Filter(prompts []*Prompt, query string) []*Prompt {
	if query == "" {
		return prompts
	}

	query = strings.ToLower(query)
	var matches []*Prompt

	for _, p := range prompts {
		if matchesQuery(p, query) {
			matches = append(matches, p)
		}
	}

	return matches
}

func matchesQuery(p *Prompt, query string) bool {
	// Check title
	if strings.Contains(strings.ToLower(p.Title), query) {
		return true
	}

	// Check description
	if strings.Contains(strings.ToLower(p.Description), query) {
		return true
	}

	// Check group
	if strings.Contains(strings.ToLower(p.Group), query) {
		return true
	}

	// Check tags
	for _, tag := range p.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	return false
}

// Validate checks if a prompt is valid
func (p *Prompt) Validate() []ValidationResult {
	var results []ValidationResult

	// Title is required
	if p.Title == "" {
		results = append(results, ValidationResult{
			Level:   ValidationError,
			Message: "missing required field: title",
		})
	}

	// Input must be valid enum
	if p.Input != "" && p.Input != "required" && p.Input != "optional" {
		results = append(results, ValidationResult{
			Level:   ValidationError,
			Message: fmt.Sprintf("invalid input value: %q (must be 'required' or 'optional')", p.Input),
		})
	}

	// Warn if no description
	if p.Description == "" {
		results = append(results, ValidationResult{
			Level:   ValidationWarning,
			Message: "missing description",
		})
	}

	return results
}

// ValidationLevel represents the severity of a validation issue
type ValidationLevel int

const (
	ValidationWarning ValidationLevel = iota
	ValidationError
)

// ValidationResult represents a validation issue
type ValidationResult struct {
	Level   ValidationLevel
	Message string
}

// ValidateDirectory checks all prompts in a directory
func ValidateDirectory(dir string) (valid int, warnings []FileValidation, errors []FileValidation) {
	prompts, err := LoadDirectory(dir)
	if err != nil {
		errors = append(errors, FileValidation{
			FileName: dir,
			Issues:   []ValidationResult{{Level: ValidationError, Message: err.Error()}},
		})
		return
	}

	// Check for duplicate titles
	titleCounts := make(map[string][]string)
	for _, p := range prompts {
		lower := strings.ToLower(p.Title)
		titleCounts[lower] = append(titleCounts[lower], p.FileName)
	}

	for _, p := range prompts {
		issues := p.Validate()

		// Add duplicate title warning
		if files := titleCounts[strings.ToLower(p.Title)]; len(files) > 1 {
			issues = append(issues, ValidationResult{
				Level:   ValidationWarning,
				Message: "duplicate title",
			})
		}

		hasError := false
		hasWarning := false
		for _, issue := range issues {
			if issue.Level == ValidationError {
				hasError = true
			} else {
				hasWarning = true
			}
		}

		if hasError {
			errors = append(errors, FileValidation{FileName: p.FileName, Issues: issues})
		} else if hasWarning {
			warnings = append(warnings, FileValidation{FileName: p.FileName, Issues: issues})
		} else {
			valid++
		}
	}

	return valid, warnings, errors
}

// FileValidation represents validation results for a file
type FileValidation struct {
	FileName string
	Issues   []ValidationResult
}

// CreatePromptFile creates a new prompt file with the given metadata
func CreatePromptFile(dir string, p *Prompt) (string, error) {
	// Get existing files for duplicate check
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	var existingFiles []string
	for _, e := range entries {
		existingFiles = append(existingFiles, e.Name())
	}

	filename := GenerateFilename(p.Title, existingFiles)
	path := filepath.Join(dir, filename)

	content := p.ToMarkdown()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return path, nil
}

// ToMarkdown converts the prompt to markdown format with frontmatter
func (p *Prompt) ToMarkdown() string {
	var fm strings.Builder
	fm.WriteString("---\n")
	fm.WriteString(fmt.Sprintf("title: %s\n", p.Title))

	if p.Description != "" {
		fm.WriteString(fmt.Sprintf("description: %s\n", p.Description))
	}
	if p.Group != "" {
		fm.WriteString(fmt.Sprintf("group: %s\n", p.Group))
	}
	if len(p.Tags) > 0 {
		fm.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(p.Tags, ", ")))
	}
	if p.Alias != "" {
		fm.WriteString(fmt.Sprintf("alias: %s\n", p.Alias))
	}
	if p.Input != "" {
		fm.WriteString(fmt.Sprintf("input: %s\n", p.Input))
	}
	if p.InputHint != "" {
		fm.WriteString(fmt.Sprintf("input_hint: %s\n", p.InputHint))
	}
	if p.Favorite {
		fm.WriteString("favorite: true\n")
	}

	fm.WriteString("---\n\n")
	fm.WriteString(p.Content)

	return fm.String()
}

// DuplicatePrompt creates a copy of a prompt with a new filename
func DuplicatePrompt(p *Prompt) (*Prompt, error) {
	dir := filepath.Dir(p.FilePath)

	newPrompt := &Prompt{
		Title:       p.Title + " (Copy)",
		Description: p.Description,
		Group:       p.Group,
		Tags:        append([]string{}, p.Tags...),
		Alias:       "", // Don't copy alias
		Input:       p.Input,
		InputHint:   p.InputHint,
		Favorite:    false, // Don't copy favorite
		Content:     p.Content,
	}

	path, err := CreatePromptFile(dir, newPrompt)
	if err != nil {
		return nil, err
	}

	newPrompt.FilePath = path
	newPrompt.FileName = filepath.Base(path)

	return newPrompt, nil
}

// DeletePrompt removes a prompt file
func DeletePrompt(p *Prompt) error {
	return os.Remove(p.FilePath)
}

// UpdateFavorite updates the favorite status in the file
func (p *Prompt) UpdateFavorite(favorite bool) error {
	p.Favorite = favorite
	content := p.ToMarkdown()
	return os.WriteFile(p.FilePath, []byte(content), 0644)
}

// colorForGroup generates a deterministic color index for a group name
func ColorIndexForGroup(group string) int {
	if group == "" {
		return -1
	}

	// FNV-1a hash for better distribution
	const fnvPrime = 16777619
	const fnvOffset = 2166136261
	var hash uint32 = fnvOffset
	for _, c := range group {
		hash ^= uint32(c)
		hash *= fnvPrime
	}

	// Return index 0-7 for 8 pastel colors
	return int(hash % 8)
}

// GetVariables returns all variable names used in the content
func (p *Prompt) GetVariables() []string {
	re := regexp.MustCompile(`\$\{([A-Z_]+)\}`)
	matches := re.FindAllStringSubmatch(p.Content, -1)

	seen := make(map[string]bool)
	var vars []string
	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			seen[match[1]] = true
			vars = append(vars, match[1])
		}
	}
	return vars
}
