// Package console provides QMUIConsole - an in-app debug console
// Ported from Tencent's QMUI_iOS framework
package console

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/log"
)

// Console is an in-app debug console
type Console struct {
	widget.BaseWidget

	// Styling
	BackgroundColor color.Color
	TextColor       color.Color
	FontSize        float32
	MaxLines        int

	// State
	logs     []string
	visible  bool
	window   fyne.Window
	popup    *widget.PopUp

	mu sync.RWMutex
}

var (
	sharedConsole *Console
	consoleOnce   sync.Once
)

// SharedConsole returns the shared console instance
func SharedConsole() *Console {
	consoleOnce.Do(func() {
		sharedConsole = NewConsole()
	})
	return sharedConsole
}

// NewConsole creates a new console
func NewConsole() *Console {
	c := &Console{
		BackgroundColor: color.RGBA{R: 0, G: 0, B: 0, A: 220},
		TextColor:       color.RGBA{R: 0, G: 255, B: 0, A: 255},
		FontSize:        12,
		MaxLines:        100,
		logs:            make([]string, 0),
	}
	c.ExtendBaseWidget(c)

	// Hook into QMUI logging
	log.SharedLogManager().GetLogger("QMUI").AddHandler(func(item *log.LogItem) {
		c.Log(item.String())
	})

	return c
}

// Log adds a log message
func (c *Console) Log(message string) {
	c.mu.Lock()
	c.logs = append(c.logs, message)
	if len(c.logs) > c.MaxLines {
		c.logs = c.logs[len(c.logs)-c.MaxLines:]
	}
	c.mu.Unlock()
	c.Refresh()
}

// Clear clears all logs
func (c *Console) Clear() {
	c.mu.Lock()
	c.logs = make([]string, 0)
	c.mu.Unlock()
	c.Refresh()
}

// Show shows the console
func (c *Console) ShowIn(window fyne.Window) {
	c.mu.Lock()
	c.window = window
	c.visible = true
	c.mu.Unlock()

	content := c.buildContent()
	c.popup = widget.NewPopUp(content, window.Canvas())

	// Position at bottom of screen
	canvasSize := window.Canvas().Size()
	consoleHeight := float32(200)
	c.popup.Resize(fyne.NewSize(canvasSize.Width, consoleHeight))
	c.popup.Move(fyne.NewPos(0, canvasSize.Height-consoleHeight))
	c.popup.Show()
}

// Hide hides the console
func (c *Console) Hide() {
	c.mu.Lock()
	c.visible = false
	if c.popup != nil {
		c.popup.Hide()
		c.popup = nil
	}
	c.mu.Unlock()
}

// Toggle toggles the console visibility
func (c *Console) Toggle(window fyne.Window) {
	c.mu.RLock()
	visible := c.visible
	c.mu.RUnlock()

	if visible {
		c.Hide()
	} else {
		c.ShowIn(window)
	}
}

// IsVisible returns whether the console is visible
func (c *Console) IsVisible() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.visible
}

func (c *Console) buildContent() fyne.CanvasObject {
	// Background
	bg := canvas.NewRectangle(c.BackgroundColor)

	// Toolbar
	toolbar := container.NewHBox(
		widget.NewButton("Clear", func() {
			c.Clear()
		}),
		widget.NewButton("Close", func() {
			c.Hide()
		}),
	)

	// Log text
	c.mu.RLock()
	logText := ""
	for _, line := range c.logs {
		logText += line + "\n"
	}
	c.mu.RUnlock()

	text := widget.NewLabel(logText)
	text.Wrapping = fyne.TextWrapWord

	scroll := container.NewVScroll(text)

	content := container.NewBorder(toolbar, nil, nil, nil, scroll)

	return container.NewStack(bg, content)
}

// CreateRenderer implements fyne.Widget
func (c *Console) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	return &consoleRenderer{console: c}
}

type consoleRenderer struct {
	console *Console
}

func (r *consoleRenderer) Destroy()                      {}
func (r *consoleRenderer) Layout(size fyne.Size)         {}
func (r *consoleRenderer) MinSize() fyne.Size            { return fyne.NewSize(0, 0) }
func (r *consoleRenderer) Refresh()                      {}
func (r *consoleRenderer) Objects() []fyne.CanvasObject { return nil }

// ConsoleToolbar provides quick actions for the console
type ConsoleToolbar struct {
	widget.BaseWidget

	BackgroundColor color.Color
	ButtonColor     color.Color

	OnClear  func()
	OnFilter func()
	OnClose  func()
}

// NewConsoleToolbar creates a console toolbar
func NewConsoleToolbar() *ConsoleToolbar {
	return &ConsoleToolbar{
		BackgroundColor: color.RGBA{R: 40, G: 40, B: 40, A: 255},
		ButtonColor:     color.White,
	}
}

// CreateRenderer implements fyne.Widget
func (ct *ConsoleToolbar) CreateRenderer() fyne.WidgetRenderer {
	ct.ExtendBaseWidget(ct)

	bg := canvas.NewRectangle(ct.BackgroundColor)

	clearBtn := widget.NewButton("Clear", func() {
		if ct.OnClear != nil {
			ct.OnClear()
		}
	})

	filterBtn := widget.NewButton("Filter", func() {
		if ct.OnFilter != nil {
			ct.OnFilter()
		}
	})

	closeBtn := widget.NewButton("X", func() {
		if ct.OnClose != nil {
			ct.OnClose()
		}
	})

	buttons := container.NewHBox(clearBtn, filterBtn, closeBtn)

	return &consoleToolbarRenderer{
		toolbar: ct,
		bg:      bg,
		buttons: buttons,
	}
}

type consoleToolbarRenderer struct {
	toolbar *ConsoleToolbar
	bg      *canvas.Rectangle
	buttons *fyne.Container
}

func (r *consoleToolbarRenderer) Destroy() {}

func (r *consoleToolbarRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.buttons.Resize(size)
}

func (r *consoleToolbarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 44)
}

func (r *consoleToolbarRenderer) Refresh() {
	r.bg.FillColor = r.toolbar.BackgroundColor
	r.bg.Refresh()
	r.buttons.Refresh()
}

func (r *consoleToolbarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.buttons}
}

// DebugPanel provides a debug panel with various tools
type DebugPanel struct {
	widget.BaseWidget

	ShowsConsole     bool
	ShowsFPS         bool
	ShowsMemory      bool

	window fyne.Window
}

// NewDebugPanel creates a new debug panel
func NewDebugPanel() *DebugPanel {
	return &DebugPanel{
		ShowsConsole: true,
		ShowsFPS:     false,
		ShowsMemory:  false,
	}
}

// ShowIn shows the debug panel
func (dp *DebugPanel) ShowIn(window fyne.Window) {
	dp.window = window
	if dp.ShowsConsole {
		SharedConsole().ShowIn(window)
	}
}

// Hide hides the debug panel
func (dp *DebugPanel) Hide() {
	SharedConsole().Hide()
}

// CreateRenderer implements fyne.Widget
func (dp *DebugPanel) CreateRenderer() fyne.WidgetRenderer {
	dp.ExtendBaseWidget(dp)
	return &debugPanelRenderer{panel: dp}
}

type debugPanelRenderer struct {
	panel *DebugPanel
}

func (r *debugPanelRenderer) Destroy()                      {}
func (r *debugPanelRenderer) Layout(size fyne.Size)         {}
func (r *debugPanelRenderer) MinSize() fyne.Size            { return fyne.NewSize(0, 0) }
func (r *debugPanelRenderer) Refresh()                      {}
func (r *debugPanelRenderer) Objects() []fyne.CanvasObject { return nil }
