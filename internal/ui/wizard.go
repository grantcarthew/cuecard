package ui

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/grantcarthew/cuecard/internal/config"
	"github.com/grantcarthew/cuecard/internal/prompt"
)

// Wizard handles first-run setup
type Wizard struct {
	window     fyne.Window
	onComplete func()
	container  *fyne.Container
	promptsDir string
	editor     string
}

// NewWizard creates a new setup wizard
func NewWizard(window fyne.Window, onComplete func()) *Wizard {
	w := &Wizard{
		window:     window,
		onComplete: onComplete,
	}
	w.showWelcome()
	return w
}

// Content returns the wizard's content
func (w *Wizard) Content() fyne.CanvasObject {
	return w.container
}

func (w *Wizard) showWelcome() {
	title := widget.NewLabel("Welcome to Cuecard")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	desc := widget.NewLabel("Cuecard is a prompt manager for AI workflows.\n\nLet's set up your prompts directory and preferences.")
	desc.Wrapping = fyne.TextWrapWord
	desc.Alignment = fyne.TextAlignCenter

	nextBtn := widget.NewButton("Get Started", func() {
		w.showPromptsDir()
	})
	nextBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		widget.NewLabel(""),
		title,
		widget.NewLabel(""),
		desc,
		widget.NewLabel(""),
		container.NewCenter(nextBtn),
	)

	w.container = container.NewCenter(content)
	w.window.SetContent(w.container)
}

func (w *Wizard) showPromptsDir() {
	title := widget.NewLabel("Prompts Directory")
	title.TextStyle = fyne.TextStyle{Bold: true}

	desc := widget.NewLabel("Choose where to store your prompt files:")
	desc.Wrapping = fyne.TextWrapWord

	// Default path
	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, "prompts")

	pathEntry := widget.NewEntry()
	pathEntry.SetText(defaultPath)

	browseBtn := widget.NewButton("Browse...", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			pathEntry.SetText(uri.Path())
		}, w.window)
		fd.Show()
	})

	createDefaultBtn := widget.NewButton("Use Default", func() {
		pathEntry.SetText(defaultPath)
	})

	pathRow := container.NewBorder(nil, nil, nil, browseBtn, pathEntry)

	backBtn := widget.NewButton("Back", func() {
		w.showWelcome()
	})

	nextBtn := widget.NewButton("Next", func() {
		w.promptsDir = pathEntry.Text

		// Create directory if it doesn't exist
		if err := os.MkdirAll(w.promptsDir, 0755); err != nil {
			dialog.ShowError(err, w.window)
			return
		}

		w.showEditor()
	})
	nextBtn.Importance = widget.HighImportance

	buttons := container.NewHBox(backBtn, createDefaultBtn, widget.NewLabel(""), nextBtn)

	content := container.NewVBox(
		title,
		widget.NewLabel(""),
		desc,
		pathRow,
		widget.NewLabel(""),
		buttons,
	)

	w.container = container.NewPadded(content)
	w.window.SetContent(w.container)
}

func (w *Wizard) showEditor() {
	title := widget.NewLabel("Default Editor")
	title.TextStyle = fyne.TextStyle{Bold: true}

	desc := widget.NewLabel("Choose your preferred text editor for editing prompts:")
	desc.Wrapping = fyne.TextWrapWord

	editors := []string{
		"code",
		"cursor",
		"vim",
		"nvim",
		"nano",
		"subl",
		"Custom...",
	}

	editorSelect := widget.NewSelect(editors, nil)
	editorSelect.SetSelected("code")

	customEntry := widget.NewEntry()
	customEntry.SetPlaceHolder("Enter custom editor command")
	customEntry.Hide()

	editorSelect.OnChanged = func(s string) {
		if s == "Custom..." {
			customEntry.Show()
		} else {
			customEntry.Hide()
		}
	}

	backBtn := widget.NewButton("Back", func() {
		w.showPromptsDir()
	})

	nextBtn := widget.NewButton("Next", func() {
		if editorSelect.Selected == "Custom..." {
			w.editor = customEntry.Text
		} else {
			w.editor = editorSelect.Selected
		}

		if w.editor == "" {
			w.editor = "code"
		}

		w.showSamplePrompt()
	})
	nextBtn.Importance = widget.HighImportance

	buttons := container.NewHBox(backBtn, widget.NewLabel(""), nextBtn)

	content := container.NewVBox(
		title,
		widget.NewLabel(""),
		desc,
		editorSelect,
		customEntry,
		widget.NewLabel(""),
		buttons,
	)

	w.container = container.NewPadded(content)
	w.window.SetContent(w.container)
}

func (w *Wizard) showSamplePrompt() {
	title := widget.NewLabel("Create Sample Prompt")
	title.TextStyle = fyne.TextStyle{Bold: true}

	desc := widget.NewLabel("Would you like to create a sample prompt to get started?")
	desc.Wrapping = fyne.TextWrapWord

	backBtn := widget.NewButton("Back", func() {
		w.showEditor()
	})

	skipBtn := widget.NewButton("Skip", func() {
		w.finishSetup(false)
	})

	createBtn := widget.NewButton("Create Sample", func() {
		w.finishSetup(true)
	})
	createBtn.Importance = widget.HighImportance

	buttons := container.NewHBox(backBtn, skipBtn, widget.NewLabel(""), createBtn)

	content := container.NewVBox(
		title,
		widget.NewLabel(""),
		desc,
		widget.NewLabel(""),
		buttons,
	)

	w.container = container.NewPadded(content)
	w.window.SetContent(w.container)
}

func (w *Wizard) finishSetup(createSample bool) {
	// Create config
	cfg := &config.Config{
		PromptsDir: w.promptsDir,
		Editor:     w.editor,
		Theme:      "system",
		Window: config.WindowConfig{
			Width:    1024,
			Height:   768,
			Position: "center",
		},
	}

	if err := cfg.Save(); err != nil {
		dialog.ShowError(err, w.window)
		return
	}

	// Create sample prompt if requested
	if createSample {
		sample := &prompt.Prompt{
			Title:       "Hello World",
			Description: "A simple greeting prompt",
			Group:       "Examples",
			Tags:        []string{"sample", "greeting"},
			Content:     "Hello! This is a sample prompt.\n\nYou can use variables like ${DATE} to insert the current date.\n\nEdit or delete this prompt to get started!",
		}

		if _, err := prompt.CreatePromptFile(w.promptsDir, sample); err != nil {
			dialog.ShowError(err, w.window)
			return
		}
	}

	// Call completion callback
	if w.onComplete != nil {
		w.onComplete()
	}
}
