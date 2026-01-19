package ui

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/grantcarthew/cuecard/internal/prompt"
)

var errInputRequired = errors.New("input is required")

// Pastel colors for groups
var groupColors = []color.Color{
	color.RGBA{255, 230, 230, 255}, // Light red
	color.RGBA{255, 243, 224, 255}, // Light orange
	color.RGBA{255, 255, 224, 255}, // Light yellow
	color.RGBA{224, 255, 224, 255}, // Light green
	color.RGBA{224, 255, 255, 255}, // Light cyan
	color.RGBA{224, 240, 255, 255}, // Light blue
	color.RGBA{240, 224, 255, 255}, // Light purple
	color.RGBA{255, 224, 240, 255}, // Light pink
}

// PromptCard is a card widget for a prompt
type PromptCard struct {
	widget.BaseWidget
	prompt     *prompt.Prompt
	app        *App
	compact    bool
	container  *fyne.Container
	inputEntry *widget.Entry
}

// NewPromptCard creates a new prompt card
func NewPromptCard(p *prompt.Prompt, app *App, compact bool) *PromptCard {
	card := &PromptCard{
		prompt:  p,
		app:     app,
		compact: compact,
	}
	card.ExtendBaseWidget(card)
	card.build()
	return card
}

func (c *PromptCard) build() {
	// Title with star
	starIcon := "☆"
	if c.prompt.Favorite {
		starIcon = "★"
	}
	starBtn := widget.NewButton(starIcon, func() {
		c.toggleFavorite()
	})
	starBtn.Importance = widget.LowImportance

	title := widget.NewLabel(c.prompt.Title)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Truncation = fyne.TextTruncateEllipsis

	titleRow := container.NewBorder(nil, nil, starBtn, nil, title)

	// Input field if needed
	var inputContainer fyne.CanvasObject
	if c.prompt.HasInput() {
		c.inputEntry = widget.NewEntry()
		c.inputEntry.SetPlaceHolder(c.prompt.InputHint)
		if c.prompt.RequiresInput() {
			c.inputEntry.Validator = func(s string) error {
				if s == "" {
					return errInputRequired
				}
				return nil
			}
		}
		inputContainer = c.inputEntry
	}

	// Copy button
	copyBtn := widget.NewButton("Copy", func() {
		c.copyToClipboard()
	})
	copyBtn.Importance = widget.HighImportance

	// Build layout
	var content *fyne.Container
	if inputContainer != nil {
		content = container.NewVBox(
			titleRow,
			inputContainer,
			container.NewCenter(copyBtn),
		)
	} else {
		content = container.NewVBox(
			titleRow,
			container.NewCenter(copyBtn),
		)
	}

	// Card background
	bg := canvas.NewRectangle(color.RGBA{245, 245, 245, 255})
	bg.CornerRadius = 8
	bg.StrokeColor = color.RGBA{220, 220, 220, 255}
	bg.StrokeWidth = 1

	c.container = container.NewStack(bg, container.NewPadded(content))
}

func (c *PromptCard) toggleFavorite() {
	c.prompt.Favorite = !c.prompt.Favorite
	if err := c.prompt.UpdateFavorite(c.prompt.Favorite); err != nil {
		return
	}
	// Rebuild to update star icon
	c.build()
	c.Refresh()
}

func (c *PromptCard) copyToClipboard() {
	if c.prompt.RequiresInput() && c.inputEntry != nil && c.inputEntry.Text == "" {
		// Show error or highlight field
		c.inputEntry.Validate()
		return
	}

	inputValue := ""
	if c.inputEntry != nil {
		inputValue = c.inputEntry.Text
	}

	c.app.CopyPrompt(c.prompt, inputValue)

	// Visual feedback - flash the card
	c.flashCopied()
}

func (c *PromptCard) flashCopied() {
	// Simple visual feedback by temporarily changing background
	// In a full implementation, this would animate
}

// CreateRenderer creates the widget renderer
func (c *PromptCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.container)
}

// Tapped handles tap events - shows context menu on long press or secondary click
func (c *PromptCard) Tapped(e *fyne.PointEvent) {
	// Normal tap does nothing special
}

// TappedSecondary handles right-click - show context menu
func (c *PromptCard) TappedSecondary(e *fyne.PointEvent) {
	c.showContextMenu(e)
}

func (c *PromptCard) showContextMenu(e *fyne.PointEvent) {
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Preview", func() {
			c.showPreview()
		}),
		fyne.NewMenuItem("Edit", func() {
			c.app.OpenInEditor(c.prompt.FilePath)
		}),
		fyne.NewMenuItem("Duplicate", func() {
			c.duplicate()
		}),
		fyne.NewMenuItem("Delete", func() {
			c.confirmDelete()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Reveal in Finder", func() {
			c.reveal()
		}),
	)

	popup := widget.NewPopUpMenu(menu, fyne.CurrentApp().Driver().CanvasForObject(c))
	popup.ShowAtPosition(e.AbsolutePosition)
}

func (c *PromptCard) showPreview() {
	content := widget.NewRichTextFromMarkdown(c.prompt.Content)
	scroll := container.NewVScroll(content)
	scroll.SetMinSize(fyne.NewSize(500, 400))

	dialog.ShowCustom(c.prompt.Title, "Close", scroll, c.app.window)
}

func (c *PromptCard) duplicate() {
	_, err := prompt.DuplicatePrompt(c.prompt)
	if err != nil {
		dialog.ShowError(err, c.app.window)
		return
	}
	c.app.refresh()
}

func (c *PromptCard) confirmDelete() {
	dialog.ShowConfirm("Delete Prompt",
		"Are you sure you want to delete \""+c.prompt.Title+"\"?",
		func(confirmed bool) {
			if confirmed {
				if err := prompt.DeletePrompt(c.prompt); err != nil {
					dialog.ShowError(err, c.app.window)
					return
				}
				c.app.refresh()
			}
		},
		c.app.window,
	)
}

func (c *PromptCard) reveal() {
	c.app.openPromptsFolder()
}

// NewGroupHeader creates a group header
func NewGroupHeader(name string, colorIndex int, isFavorites bool) fyne.CanvasObject {
	label := widget.NewLabel(name)
	label.TextStyle = fyne.TextStyle{Bold: true}

	var bgColor color.Color
	if isFavorites {
		bgColor = color.RGBA{255, 248, 220, 255} // Light gold for favorites
	} else if colorIndex >= 0 && colorIndex < len(groupColors) {
		bgColor = groupColors[colorIndex]
	} else {
		bgColor = color.RGBA{240, 240, 240, 255}
	}

	bg := canvas.NewRectangle(bgColor)
	bg.CornerRadius = 4

	return container.NewStack(bg, container.NewPadded(label))
}

// PromptListItem is a list view item for a prompt
type PromptListItem struct {
	widget.BaseWidget
	prompt    *prompt.Prompt
	app       *App
	container *fyne.Container
}

// NewPromptListItem creates a new list item
func NewPromptListItem(p *prompt.Prompt, app *App) *PromptListItem {
	item := &PromptListItem{
		prompt: p,
		app:    app,
	}
	item.ExtendBaseWidget(item)
	item.build()
	return item
}

func (li *PromptListItem) build() {
	starIcon := "☆"
	if li.prompt.Favorite {
		starIcon = "★"
	}
	starBtn := widget.NewButton(starIcon, func() {
		li.prompt.Favorite = !li.prompt.Favorite
		li.prompt.UpdateFavorite(li.prompt.Favorite)
		li.build()
		li.Refresh()
	})
	starBtn.Importance = widget.LowImportance

	title := widget.NewLabel(li.prompt.Title)
	title.TextStyle = fyne.TextStyle{Bold: true}

	group := widget.NewLabel(li.prompt.Group)
	group.Importance = widget.LowImportance

	copyBtn := widget.NewButton("Copy", func() {
		li.app.CopyPrompt(li.prompt, "")
	})
	copyBtn.Importance = widget.HighImportance

	li.container = container.NewBorder(
		nil, nil,
		container.NewHBox(starBtn, title),
		container.NewHBox(group, copyBtn),
		nil,
	)
}

// CreateRenderer creates the widget renderer
func (li *PromptListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(li.container)
}
