package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/grantcarthew/cuecard/internal/prompt"
)

// ShowNewPromptDialog shows the new prompt creation dialog
func ShowNewPromptDialog(window fyne.Window, promptsDir string, onCreated func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Prompt title")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("Brief description (optional)")

	groupEntry := widget.NewEntry()
	groupEntry.SetPlaceHolder("Group name (optional)")

	tagsEntry := widget.NewEntry()
	tagsEntry.SetPlaceHolder("tag1, tag2, tag3 (optional)")

	inputSelect := widget.NewSelect([]string{"None", "Optional", "Required"}, nil)
	inputSelect.SetSelected("None")

	inputHintEntry := widget.NewEntry()
	inputHintEntry.SetPlaceHolder("Input field hint (optional)")

	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetPlaceHolder("Prompt content...\n\nUse ${INPUT} for user input, ${DATE}, ${DATETIME}, ${CLIPBOARD}, ${FILE}")
	contentEntry.SetMinRowsVisible(8)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Description", Widget: descEntry},
			{Text: "Group", Widget: groupEntry},
			{Text: "Tags", Widget: tagsEntry},
			{Text: "Input", Widget: inputSelect},
			{Text: "Input Hint", Widget: inputHintEntry},
			{Text: "Content", Widget: contentEntry},
		},
		OnSubmit: func() {
			// Parse tags
			var tags []string
			if tagsEntry.Text != "" {
				for _, t := range strings.Split(tagsEntry.Text, ",") {
					t = strings.TrimSpace(t)
					if t != "" {
						tags = append(tags, t)
					}
				}
			}

			// Determine input value
			inputValue := ""
			switch inputSelect.Selected {
			case "Optional":
				inputValue = "optional"
			case "Required":
				inputValue = "required"
			}

			p := &prompt.Prompt{
				Title:       titleEntry.Text,
				Description: descEntry.Text,
				Group:       groupEntry.Text,
				Tags:        tags,
				Input:       inputValue,
				InputHint:   inputHintEntry.Text,
				Content:     contentEntry.Text,
			}

			if _, err := prompt.CreatePromptFile(promptsDir, p); err != nil {
				dialog.ShowError(err, window)
				return
			}

			if onCreated != nil {
				onCreated()
			}
		},
	}

	d := dialog.NewForm("New Prompt", "Create", "Cancel", form.Items, func(submitted bool) {
		if submitted {
			form.OnSubmit()
		}
	}, window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}

// ShowImportFileDialog shows the file import dialog
func ShowImportFileDialog(window fyne.Window, promptsDir string, onImported func()) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if reader == nil {
			return // Cancelled
		}
		defer reader.Close()

		// Read file content
		data := make([]byte, 0)
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				data = append(data, buf[:n]...)
			}
			if err != nil {
				break
			}
		}

		content := string(data)

		// Check if it has valid frontmatter
		if prompt.HasFrontmatter(content) {
			p, err := prompt.Parse(content)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid frontmatter: %w", err), window)
				return
			}

			// Clear favorite for imports
			p.Favorite = false
			p.Content = content[strings.Index(content, "---\n")+4:]
			if idx := strings.Index(p.Content, "---\n"); idx >= 0 {
				p.Content = strings.TrimSpace(p.Content[idx+4:])
			}

			// Create the file
			if _, err := prompt.CreatePromptFile(promptsDir, p); err != nil {
				dialog.ShowError(err, window)
				return
			}

			if onImported != nil {
				onImported()
			}
		} else {
			// Show form to add metadata
			ShowImportFormDialog(window, promptsDir, content, onImported)
		}
	}, window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".md"}))
	fd.Show()
}

// ShowImportFormDialog shows a form to add metadata to imported content
func ShowImportFormDialog(window fyne.Window, promptsDir string, content string, onImported func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Prompt title")

	groupEntry := widget.NewEntry()
	groupEntry.SetPlaceHolder("Group name (optional)")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Group", Widget: groupEntry},
		},
		OnSubmit: func() {
			p := &prompt.Prompt{
				Title:   titleEntry.Text,
				Group:   groupEntry.Text,
				Content: content,
			}

			if _, err := prompt.CreatePromptFile(promptsDir, p); err != nil {
				dialog.ShowError(err, window)
				return
			}

			if onImported != nil {
				onImported()
			}
		},
	}

	d := dialog.NewForm("Import Prompt", "Import", "Cancel", form.Items, func(submitted bool) {
		if submitted {
			form.OnSubmit()
		}
	}, window)
	d.Show()
}

// ShowImportDirectoryDialog shows the directory import dialog
func ShowImportDirectoryDialog(window fyne.Window, promptsDir string, onImported func()) {
	fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if uri == nil {
			return // Cancelled
		}

		importDir := uri.Path()
		imported := 0
		skipped := 0
		var errors []string

		entries, err := os.ReadDir(importDir)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			if !strings.HasSuffix(strings.ToLower(name), ".md") {
				continue
			}
			if strings.EqualFold(name, "readme.md") {
				skipped++
				continue
			}

			path := filepath.Join(importDir, name)
			data, err := os.ReadFile(path)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: read error", name))
				continue
			}

			content := string(data)
			if !prompt.HasFrontmatter(content) {
				skipped++
				continue
			}

			p, err := prompt.Parse(content)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: parse error", name))
				continue
			}

			// Clear favorite for imports
			p.Favorite = false

			if _, err := prompt.CreatePromptFile(promptsDir, p); err != nil {
				errors = append(errors, fmt.Sprintf("%s: create error", name))
				continue
			}

			imported++
		}

		// Show summary
		summary := fmt.Sprintf("Imported: %d\nSkipped: %d", imported, skipped)
		if len(errors) > 0 {
			summary += fmt.Sprintf("\nErrors: %d\n\n%s", len(errors), strings.Join(errors, "\n"))
		}

		dialog.ShowInformation("Import Complete", summary, window)

		if onImported != nil && imported > 0 {
			onImported()
		}
	}, window)

	fd.Show()
}

// ShowValidationDialog shows the validation results
func ShowValidationDialog(window fyne.Window, promptsDir string) {
	valid, warnings, errors := prompt.ValidateDirectory(promptsDir)

	var content strings.Builder
	content.WriteString(fmt.Sprintf("%d prompts valid\n", valid))

	if len(warnings) > 0 {
		content.WriteString(fmt.Sprintf("\n%d warnings:\n", len(warnings)))
		for _, w := range warnings {
			for _, issue := range w.Issues {
				if issue.Level == prompt.ValidationWarning {
					content.WriteString(fmt.Sprintf("  - %s: %s\n", w.FileName, issue.Message))
				}
			}
		}
	}

	if len(errors) > 0 {
		content.WriteString(fmt.Sprintf("\n%d errors:\n", len(errors)))
		for _, e := range errors {
			for _, issue := range e.Issues {
				if issue.Level == prompt.ValidationError {
					content.WriteString(fmt.Sprintf("  - %s: %s\n", e.FileName, issue.Message))
				}
			}
		}
	}

	label := widget.NewLabel(content.String())
	label.Wrapping = fyne.TextWrapWord

	scroll := container.NewVScroll(label)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	dialog.ShowCustom("Validation Results", "Close", scroll, window)
}

// ShowAboutDialog shows the about dialog
func ShowAboutDialog(window fyne.Window) {
	content := widget.NewRichTextFromMarkdown(`# Cuecard

Cross-platform GUI prompt manager for AI workflows.

**Version:** 1.0.0

Point, click, prompt — brain cells optional.

---

[GitHub](https://github.com/grantcarthew/cuecard)
`)

	dialog.ShowCustom("About Cuecard", "Close", content, window)
}
