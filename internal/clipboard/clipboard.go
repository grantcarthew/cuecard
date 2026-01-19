package clipboard

import (
	"fyne.io/fyne/v2"
)

// Clipboard provides clipboard operations
type Clipboard struct {
	window fyne.Window
}

// New creates a new Clipboard instance
func New(window fyne.Window) *Clipboard {
	return &Clipboard{window: window}
}

// Copy copies text to the system clipboard
func (c *Clipboard) Copy(text string) {
	c.window.Clipboard().SetContent(text)
}

// Read returns the current clipboard content
func (c *Clipboard) Read() string {
	return c.window.Clipboard().Content()
}

// ClipboardInterface abstracts clipboard operations for testing
type ClipboardInterface interface {
	Copy(text string)
	Read() string
}

// Ensure Clipboard implements ClipboardInterface
var _ ClipboardInterface = (*Clipboard)(nil)

// MockClipboard is a test implementation
type MockClipboard struct {
	Content string
}

// Copy stores text in the mock clipboard
func (m *MockClipboard) Copy(text string) {
	m.Content = text
}

// Read returns the mock clipboard content
func (m *MockClipboard) Read() string {
	return m.Content
}

// Ensure MockClipboard implements ClipboardInterface
var _ ClipboardInterface = (*MockClipboard)(nil)
