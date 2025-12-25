// Package tile provides QMUIComponentTile - styled grid tiles for component showcases
// Matches the beautiful QMUI iOS styling with rounded borders and accent colors
package tile

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/core"
)

// ComponentTile is a styled tile for showcasing components in a grid
// Matches the QMUI iOS demo styling with rounded borders and icons
type ComponentTile struct {
	widget.BaseWidget

	// Content
	Title string
	Icon  fyne.CanvasObject // Custom icon/content to display

	// Styling
	TintColor       color.Color
	BackgroundColor color.Color
	BorderWidth     float32
	CornerRadius    float32
	TitleColor      color.Color
	TitleFontSize   float32
	IconSize        fyne.Size
	Padding         float32

	// Callbacks
	OnTapped func()

	mu      sync.RWMutex
	hovered bool
}

// NewComponentTile creates a new styled component tile
func NewComponentTile(title string, icon fyne.CanvasObject) *ComponentTile {
	cfg := core.SharedConfiguration()
	t := &ComponentTile{
		Title:           title,
		Icon:            icon,
		TintColor:       cfg.BlueColor,
		BackgroundColor: color.White,
		BorderWidth:     1.5,
		CornerRadius:    8,
		TitleColor:      cfg.GrayDarkenColor,
		TitleFontSize:   11,
		IconSize:        fyne.NewSize(48, 48),
		Padding:         12,
	}
	t.ExtendBaseWidget(t)
	return t
}

// NewComponentTileWithDrawing creates a tile with a custom drawing function
func NewComponentTileWithDrawing(title string, drawFunc func(size fyne.Size, tint color.Color) fyne.CanvasObject) *ComponentTile {
	cfg := core.SharedConfiguration()
	icon := drawFunc(fyne.NewSize(48, 48), cfg.BlueColor)
	return NewComponentTile(title, icon)
}

// CreateRenderer implements fyne.Widget
func (t *ComponentTile) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	background := canvas.NewRectangle(t.BackgroundColor)
	background.CornerRadius = t.CornerRadius
	background.StrokeWidth = t.BorderWidth
	background.StrokeColor = t.TintColor

	title := canvas.NewText(t.Title, t.TitleColor)
	title.TextSize = t.TitleFontSize
	title.Alignment = fyne.TextAlignCenter

	return &tileRenderer{
		tile:       t,
		background: background,
		title:      title,
	}
}

func (t *ComponentTile) Tapped(_ *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped()
	}
}

func (t *ComponentTile) TappedSecondary(_ *fyne.PointEvent) {}

func (t *ComponentTile) MouseIn(_ *desktop.MouseEvent) {
	t.mu.Lock()
	t.hovered = true
	t.mu.Unlock()
	t.Refresh()
}

func (t *ComponentTile) MouseMoved(_ *desktop.MouseEvent) {}

func (t *ComponentTile) MouseOut() {
	t.mu.Lock()
	t.hovered = false
	t.mu.Unlock()
	t.Refresh()
}

func (t *ComponentTile) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type tileRenderer struct {
	tile       *ComponentTile
	background *canvas.Rectangle
	title      *canvas.Text
}

func (r *tileRenderer) Destroy() {}

func (r *tileRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	padding := r.tile.Padding
	iconSize := r.tile.IconSize
	titleSize := r.title.MinSize()

	// Center icon horizontally, position in upper area
	if r.tile.Icon != nil {
		iconX := (size.Width - iconSize.Width) / 2
		iconY := padding + 8
		r.tile.Icon.Resize(iconSize)
		r.tile.Icon.Move(fyne.NewPos(iconX, iconY))
	}

	// Title at bottom, centered
	titleY := size.Height - titleSize.Height - padding
	r.title.Resize(fyne.NewSize(size.Width, titleSize.Height))
	r.title.Move(fyne.NewPos(0, titleY))
}

func (r *tileRenderer) MinSize() fyne.Size {
	padding := r.tile.Padding
	iconSize := r.tile.IconSize
	titleSize := r.title.MinSize()

	width := iconSize.Width + padding*2
	if titleSize.Width+padding*2 > width {
		width = titleSize.Width + padding*2
	}

	height := padding + iconSize.Height + 8 + titleSize.Height + padding

	return fyne.NewSize(width, height)
}

func (r *tileRenderer) Refresh() {
	r.tile.mu.RLock()
	hovered := r.tile.hovered
	r.tile.mu.RUnlock()

	r.background.FillColor = r.tile.BackgroundColor
	r.background.StrokeColor = r.tile.TintColor
	r.background.StrokeWidth = r.tile.BorderWidth
	r.background.CornerRadius = r.tile.CornerRadius

	if hovered {
		// Subtle highlight on hover
		r.background.FillColor = color.RGBA{R: 245, G: 250, B: 255, A: 255}
	}

	r.title.Text = r.tile.Title
	r.title.Color = r.tile.TitleColor
	r.title.TextSize = r.tile.TitleFontSize

	r.background.Refresh()
	r.title.Refresh()
	if r.tile.Icon != nil {
		r.tile.Icon.Refresh()
	}
}

func (r *tileRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}
	if r.tile.Icon != nil {
		objects = append(objects, r.tile.Icon)
	}
	objects = append(objects, r.title)
	return objects
}

// Icon drawing helpers for common QMUI component icons

// ButtonIcon creates a button icon (rounded rect with "OK" text)
func ButtonIcon(size fyne.Size, tint color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeColor = tint
	rect.StrokeWidth = 1.5
	rect.CornerRadius = size.Height / 4
	rect.Resize(fyne.NewSize(size.Width*0.7, size.Height*0.4))

	text := canvas.NewText("OK", tint)
	text.TextSize = 10
	text.Alignment = fyne.TextAlignCenter

	return &iconContainer{
		rect: rect,
		text: text,
		size: size,
	}
}

// LabelIcon creates a label icon ("A" in a box)
func LabelIcon(size fyne.Size, tint color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeColor = tint
	rect.StrokeWidth = 1.5
	rect.CornerRadius = 4

	text := canvas.NewText("A", tint)
	text.TextSize = 18
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Alignment = fyne.TextAlignCenter

	return &iconContainer{
		rect: rect,
		text: text,
		size: size,
	}
}

// TextViewIcon creates a text view icon (lines)
func TextViewIcon(size fyne.Size, tint color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeColor = tint
	rect.StrokeWidth = 1.5
	rect.CornerRadius = 4

	// Create lines inside
	line1 := canvas.NewRectangle(tint)
	line2 := canvas.NewRectangle(tint)
	line3 := canvas.NewRectangle(tint)

	return &linesIconContainer{
		rect:  rect,
		lines: []*canvas.Rectangle{line1, line2, line3},
		size:  size,
		tint:  tint,
	}
}

// SliderIcon creates a slider icon
func SliderIcon(size fyne.Size, tint color.Color) fyne.CanvasObject {
	track := canvas.NewRectangle(tint)
	track.CornerRadius = 2

	thumb := canvas.NewCircle(tint)

	return &sliderIconContainer{
		track: track,
		thumb: thumb,
		size:  size,
	}
}

// TableIcon creates a table view icon
func TableIcon(size fyne.Size, tint color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeColor = tint
	rect.StrokeWidth = 1.5
	rect.CornerRadius = 4

	row1 := canvas.NewRectangle(tint)
	row2 := canvas.NewRectangle(tint)
	row3 := canvas.NewRectangle(tint)

	return &tableIconContainer{
		rect: rect,
		rows: []*canvas.Rectangle{row1, row2, row3},
		size: size,
	}
}

// SearchIcon creates a search/magnifier icon
func SearchIcon(size fyne.Size, tint color.Color) fyne.CanvasObject {
	circle := canvas.NewCircle(color.Transparent)
	circle.StrokeColor = tint
	circle.StrokeWidth = 2

	handle := canvas.NewRectangle(tint)

	return &searchIconContainer{
		circle: circle,
		handle: handle,
		size:   size,
	}
}

// Helper containers for complex icons

type iconContainer struct {
	rect *canvas.Rectangle
	text *canvas.Text
	size fyne.Size
}

func (c *iconContainer) MinSize() fyne.Size { return c.size }
func (c *iconContainer) Move(pos fyne.Position) {
	c.rect.Move(fyne.NewPos(pos.X+c.size.Width*0.15, pos.Y+c.size.Height*0.3))
	c.text.Move(fyne.NewPos(pos.X, pos.Y+c.size.Height*0.3))
}
func (c *iconContainer) Position() fyne.Position { return c.rect.Position() }
func (c *iconContainer) Resize(size fyne.Size) {
	c.size = size
	c.rect.Resize(fyne.NewSize(size.Width*0.7, size.Height*0.4))
	c.text.Resize(fyne.NewSize(size.Width*0.7, size.Height*0.4))
}
func (c *iconContainer) Size() fyne.Size    { return c.size }
func (c *iconContainer) Show()              { c.rect.Show(); c.text.Show() }
func (c *iconContainer) Hide()              { c.rect.Hide(); c.text.Hide() }
func (c *iconContainer) Visible() bool      { return c.rect.Visible() }
func (c *iconContainer) Refresh()           { c.rect.Refresh(); c.text.Refresh() }

type linesIconContainer struct {
	rect  *canvas.Rectangle
	lines []*canvas.Rectangle
	size  fyne.Size
	tint  color.Color
}

func (c *linesIconContainer) MinSize() fyne.Size { return c.size }
func (c *linesIconContainer) Move(pos fyne.Position) {
	c.rect.Move(pos)
	lineH := float32(2)
	spacing := c.size.Height * 0.2
	startY := pos.Y + c.size.Height*0.25
	for i, line := range c.lines {
		w := c.size.Width * 0.6
		if i == 2 {
			w = c.size.Width * 0.4
		}
		line.Move(fyne.NewPos(pos.X+c.size.Width*0.2, startY+float32(i)*spacing))
		line.Resize(fyne.NewSize(w, lineH))
	}
}
func (c *linesIconContainer) Position() fyne.Position { return c.rect.Position() }
func (c *linesIconContainer) Resize(size fyne.Size)   { c.size = size; c.rect.Resize(size) }
func (c *linesIconContainer) Size() fyne.Size         { return c.size }
func (c *linesIconContainer) Show()                   { c.rect.Show() }
func (c *linesIconContainer) Hide()                   { c.rect.Hide() }
func (c *linesIconContainer) Visible() bool           { return c.rect.Visible() }
func (c *linesIconContainer) Refresh() {
	c.rect.Refresh()
	for _, l := range c.lines {
		l.Refresh()
	}
}

type sliderIconContainer struct {
	track *canvas.Rectangle
	thumb *canvas.Circle
	size  fyne.Size
}

func (c *sliderIconContainer) MinSize() fyne.Size { return c.size }
func (c *sliderIconContainer) Move(pos fyne.Position) {
	trackY := pos.Y + c.size.Height/2 - 2
	c.track.Move(fyne.NewPos(pos.X+c.size.Width*0.1, trackY))
	c.track.Resize(fyne.NewSize(c.size.Width*0.8, 4))
	thumbX := pos.X + c.size.Width*0.5
	thumbY := pos.Y + c.size.Height/2 - 8
	c.thumb.Move(fyne.NewPos(thumbX, thumbY))
	c.thumb.Resize(fyne.NewSize(16, 16))
}
func (c *sliderIconContainer) Position() fyne.Position { return c.track.Position() }
func (c *sliderIconContainer) Resize(size fyne.Size)   { c.size = size }
func (c *sliderIconContainer) Size() fyne.Size         { return c.size }
func (c *sliderIconContainer) Show()                   { c.track.Show(); c.thumb.Show() }
func (c *sliderIconContainer) Hide()                   { c.track.Hide(); c.thumb.Hide() }
func (c *sliderIconContainer) Visible() bool           { return c.track.Visible() }
func (c *sliderIconContainer) Refresh()                { c.track.Refresh(); c.thumb.Refresh() }

type tableIconContainer struct {
	rect *canvas.Rectangle
	rows []*canvas.Rectangle
	size fyne.Size
}

func (c *tableIconContainer) MinSize() fyne.Size { return c.size }
func (c *tableIconContainer) Move(pos fyne.Position) {
	c.rect.Move(pos)
	c.rect.Resize(c.size)
	rowH := float32(3)
	spacing := c.size.Height * 0.22
	startY := pos.Y + c.size.Height*0.2
	for i, row := range c.rows {
		row.Move(fyne.NewPos(pos.X+c.size.Width*0.15, startY+float32(i)*spacing))
		row.Resize(fyne.NewSize(c.size.Width*0.7, rowH))
	}
}
func (c *tableIconContainer) Position() fyne.Position { return c.rect.Position() }
func (c *tableIconContainer) Resize(size fyne.Size)   { c.size = size }
func (c *tableIconContainer) Size() fyne.Size         { return c.size }
func (c *tableIconContainer) Show()                   { c.rect.Show() }
func (c *tableIconContainer) Hide()                   { c.rect.Hide() }
func (c *tableIconContainer) Visible() bool           { return c.rect.Visible() }
func (c *tableIconContainer) Refresh() {
	c.rect.Refresh()
	for _, r := range c.rows {
		r.Refresh()
	}
}

type searchIconContainer struct {
	circle *canvas.Circle
	handle *canvas.Rectangle
	size   fyne.Size
}

func (c *searchIconContainer) MinSize() fyne.Size { return c.size }
func (c *searchIconContainer) Move(pos fyne.Position) {
	circleSize := c.size.Width * 0.5
	c.circle.Move(fyne.NewPos(pos.X+c.size.Width*0.15, pos.Y+c.size.Height*0.15))
	c.circle.Resize(fyne.NewSize(circleSize, circleSize))
	// Handle from bottom-right of circle
	handleX := pos.X + c.size.Width*0.55
	handleY := pos.Y + c.size.Height*0.55
	c.handle.Move(fyne.NewPos(handleX, handleY))
	c.handle.Resize(fyne.NewSize(c.size.Width*0.3, 3))
}
func (c *searchIconContainer) Position() fyne.Position { return c.circle.Position() }
func (c *searchIconContainer) Resize(size fyne.Size)   { c.size = size }
func (c *searchIconContainer) Size() fyne.Size         { return c.size }
func (c *searchIconContainer) Show()                   { c.circle.Show(); c.handle.Show() }
func (c *searchIconContainer) Hide()                   { c.circle.Hide(); c.handle.Hide() }
func (c *searchIconContainer) Visible() bool           { return c.circle.Visible() }
func (c *searchIconContainer) Refresh()                { c.circle.Refresh(); c.handle.Refresh() }
