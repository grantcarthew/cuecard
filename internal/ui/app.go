package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/grantcarthew/cuecard/internal/clipboard"
	"github.com/grantcarthew/cuecard/internal/config"
	"github.com/grantcarthew/cuecard/internal/prompt"
	"github.com/grantcarthew/cuecard/internal/watcher"
)

const appID = "com.grantcarthew.cuecard"

// App represents the main application
type App struct {
	fyneApp   fyne.App
	window    fyne.Window
	config    *config.Config
	prompts   []*prompt.Prompt
	clipboard *clipboard.Clipboard
	watcher   *watcher.Watcher
	mainView  *MainView
}

// New creates a new application instance
func New() *App {
	fyneApp := app.NewWithID(appID)
	return &App{
		fyneApp: fyneApp,
	}
}

// Run starts the application
func (a *App) Run() error {
	// Create main window
	a.window = a.fyneApp.NewWindow("Cuecard")
	a.window.Resize(fyne.NewSize(1024, 768))
	a.window.CenterOnScreen()

	// Check if config exists
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("failed to check config: %w", err)
	}

	if !exists {
		// Show first-run wizard in main window
		wizard := NewWizard(a.window, func() {
			// Wizard complete - load main UI
			a.loadMainUI()
		})
		a.window.SetContent(wizard.Content())
	} else {
		// Load main UI directly
		if err := a.loadMainUI(); err != nil {
			return err
		}
	}

	// Handle window close
	a.window.SetCloseIntercept(func() {
		if desk, ok := a.fyneApp.(desktop.App); ok {
			// Minimize to tray instead of quitting
			a.window.Hide()
			_ = desk
		} else {
			a.fyneApp.Quit()
		}
	})

	// Show and run
	a.window.ShowAndRun()

	// Cleanup
	if a.watcher != nil {
		a.watcher.Stop()
	}

	return nil
}

func (a *App) loadMainUI() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.config = cfg

	// Apply theme
	a.applyTheme()

	// Resize window based on config
	a.window.Resize(fyne.NewSize(float32(cfg.Window.Width), float32(cfg.Window.Height)))

	// Initialize clipboard
	a.clipboard = clipboard.New(a.window)

	// Load prompts
	if err := a.loadPrompts(); err != nil {
		return fmt.Errorf("failed to load prompts: %w", err)
	}

	// Create main view
	a.mainView = NewMainView(a)
	a.window.SetContent(a.mainView.Container())

	// Set up menus
	a.setupMenus()

	// Set up system tray
	a.setupSystemTray()

	// Set up file watcher
	if err := a.setupWatcher(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: file watcher failed: %v\n", err)
	}

	return nil
}

func (a *App) applyTheme() {
	switch a.config.Theme {
	case "light":
		a.fyneApp.Settings().SetTheme(&lightTheme{})
	case "dark":
		a.fyneApp.Settings().SetTheme(&darkTheme{})
	default:
		// Use system theme (default Fyne behavior)
	}
}

func (a *App) loadPrompts() error {
	prompts, err := prompt.LoadDirectory(a.config.PromptsDir)
	if err != nil {
		return err
	}
	a.prompts = prompts
	return nil
}

func (a *App) setupWatcher() error {
	w, err := watcher.New(a.config.PromptsDir, func() {
		// Reload prompts on change
		if err := a.loadPrompts(); err != nil {
			return
		}
		// Refresh UI on main thread
		a.mainView.Refresh(a.prompts)
	})
	if err != nil {
		return err
	}

	a.watcher = w
	return w.Start()
}

func (a *App) setupSystemTray() {
	if desk, ok := a.fyneApp.(desktop.App); ok {
		menu := fyne.NewMenu("Cuecard",
			fyne.NewMenuItem("Show", func() {
				a.window.Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				a.fyneApp.Quit()
			}),
		)
		desk.SetSystemTrayMenu(menu)
	}
}

func (a *App) setupMenus() {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("New Prompt", a.showNewPromptDialog),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Import File...", a.importFile),
		fyne.NewMenuItem("Import Directory...", a.importDirectory),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Open Folder", a.openPromptsFolder),
		fyne.NewMenuItem("Refresh", a.refresh),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Validate Prompts", a.showValidation),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("About", a.showAbout),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	a.window.SetMainMenu(mainMenu)
}

func (a *App) showNewPromptDialog() {
	ShowNewPromptDialog(a.window, a.config.PromptsDir, func() {
		a.refresh()
	})
}

func (a *App) importFile() {
	ShowImportFileDialog(a.window, a.config.PromptsDir, func() {
		a.refresh()
	})
}

func (a *App) importDirectory() {
	ShowImportDirectoryDialog(a.window, a.config.PromptsDir, func() {
		a.refresh()
	})
}

func (a *App) openPromptsFolder() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", a.config.PromptsDir)
	case "linux":
		cmd = exec.Command("xdg-open", a.config.PromptsDir)
	default:
		return
	}
	cmd.Start()
}

func (a *App) refresh() {
	if err := a.loadPrompts(); err != nil {
		return
	}
	a.mainView.Refresh(a.prompts)
}

func (a *App) showValidation() {
	ShowValidationDialog(a.window, a.config.PromptsDir)
}

func (a *App) showAbout() {
	ShowAboutDialog(a.window)
}

// OpenInEditor opens a file in the configured editor
func (a *App) OpenInEditor(path string) {
	cmd := exec.Command(a.config.Editor, path)
	cmd.Start()
}

// CopyPrompt copies a prompt to clipboard with variable substitution
func (a *App) CopyPrompt(p *prompt.Prompt, inputValue string) {
	resolver := &prompt.VariableResolver{
		Input:     inputValue,
		Clipboard: a.clipboard.Read(),
		FileSelector: func() string {
			return ""
		},
	}

	content := prompt.Substitute(p.Content, resolver)
	a.clipboard.Copy(content)
}

// GetPrompts returns the current prompts
func (a *App) GetPrompts() []*prompt.Prompt {
	return a.prompts
}

// GetConfig returns the configuration
func (a *App) GetConfig() *config.Config {
	return a.config
}
