package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/grantcarthew/cuecard/internal/prompt"
)

// MainView is the main content view
type MainView struct {
	app           *App
	container     *fyne.Container
	searchEntry   *widget.Entry
	cardsScroll   *container.Scroll
	cardsContent  *fyne.Container
	compactMode   bool
	listView      bool
	alwaysOnTop   bool
	currentFilter string
}

// NewMainView creates a new main view
func NewMainView(app *App) *MainView {
	mv := &MainView{
		app: app,
	}
	mv.build()
	return mv
}

func (mv *MainView) build() {
	// Search bar
	mv.searchEntry = widget.NewEntry()
	mv.searchEntry.SetPlaceHolder("Search prompts...")
	mv.searchEntry.OnChanged = mv.onSearch

	// View toggles
	compactToggle := widget.NewCheck("Compact", func(checked bool) {
		mv.compactMode = checked
		mv.rebuildCards()
	})

	listToggle := widget.NewCheck("List", func(checked bool) {
		mv.listView = checked
		mv.rebuildCards()
	})

	topToggle := widget.NewCheck("Always on Top", func(checked bool) {
		mv.alwaysOnTop = checked
		// Note: Fyne doesn't have direct always-on-top support
		// This would require platform-specific code
	})

	toolbar := container.NewBorder(
		nil, nil,
		nil,
		container.NewHBox(compactToggle, listToggle, topToggle),
		mv.searchEntry,
	)

	// Cards container
	mv.cardsContent = container.NewVBox()
	mv.cardsScroll = container.NewVScroll(mv.cardsContent)

	// Build the cards
	mv.rebuildCards()

	// Main layout
	mv.container = container.NewBorder(
		toolbar,
		nil, nil, nil,
		mv.cardsScroll,
	)
}

func (mv *MainView) onSearch(query string) {
	mv.currentFilter = query
	mv.rebuildCards()
}

func (mv *MainView) rebuildCards() {
	mv.cardsContent.RemoveAll()

	prompts := mv.app.GetPrompts()

	// Apply filter if any
	if mv.currentFilter != "" {
		prompts = prompt.Filter(prompts, mv.currentFilter)
	}

	// Group prompts
	favorites, groups, ungrouped := prompt.GroupPrompts(prompts)

	// Add favorites section
	if len(favorites) > 0 {
		header := NewGroupHeader("Favorites", -1, true)
		mv.cardsContent.Add(header)

		cardsContainer := mv.createCardsContainer(favorites)
		mv.cardsContent.Add(cardsContainer)
	}

	// Add groups
	groupNames := prompt.SortedGroupNames(groups)
	for i, name := range groupNames {
		header := NewGroupHeader(name, i%8, false)
		mv.cardsContent.Add(header)

		cardsContainer := mv.createCardsContainer(groups[name])
		mv.cardsContent.Add(cardsContainer)
	}

	// Add ungrouped
	if len(ungrouped) > 0 {
		if len(favorites) > 0 || len(groups) > 0 {
			// Add some spacing
			mv.cardsContent.Add(widget.NewSeparator())
		}

		cardsContainer := mv.createCardsContainer(ungrouped)
		mv.cardsContent.Add(cardsContainer)
	}

	mv.cardsContent.Refresh()
}

func (mv *MainView) createCardsContainer(prompts []*prompt.Prompt) fyne.CanvasObject {
	if mv.listView {
		return mv.createListView(prompts)
	}
	return mv.createGridView(prompts)
}

func (mv *MainView) createGridView(prompts []*prompt.Prompt) fyne.CanvasObject {
	cardSize := fyne.NewSize(250, 120)
	if mv.compactMode {
		cardSize = fyne.NewSize(200, 80)
	}

	var cards []fyne.CanvasObject
	for _, p := range prompts {
		card := NewPromptCard(p, mv.app, mv.compactMode)
		cards = append(cards, card)
	}

	grid := container.NewGridWrap(cardSize, cards...)
	return grid
}

func (mv *MainView) createListView(prompts []*prompt.Prompt) fyne.CanvasObject {
	var items []fyne.CanvasObject
	for _, p := range prompts {
		item := NewPromptListItem(p, mv.app)
		items = append(items, item)
	}
	return container.NewVBox(items...)
}

// Container returns the main container
func (mv *MainView) Container() fyne.CanvasObject {
	return mv.container
}

// Refresh rebuilds the view with new prompts
func (mv *MainView) Refresh(prompts []*prompt.Prompt) {
	mv.rebuildCards()
}
