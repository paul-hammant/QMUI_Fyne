// Package grid provides QMUIGridView - a grid layout container
// Ported from Tencent's QMUI_iOS framework
package grid

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// GridView displays items in a grid layout
type GridView struct {
	widget.BaseWidget

	// Layout
	ColumnCount     int
	RowHeight       float32
	ColumnSpacing   float32
	RowSpacing      float32
	ContentInsets   core.EdgeInsets

	// Styling
	BackgroundColor color.Color
	SeparatorColor  color.Color
	SeparatorWidth  float32
	ShowSeparators  bool

	// Items
	items []fyne.CanvasObject

	mu sync.RWMutex
}

// NewGridView creates a new grid view
func NewGridView(columnCount int) *GridView {
	gv := &GridView{
		ColumnCount:     columnCount,
		RowHeight:       0, // Auto height
		ColumnSpacing:   0,
		RowSpacing:      0,
		ContentInsets:   core.EdgeInsets{},
		BackgroundColor: color.Transparent,
		SeparatorColor:  core.SharedConfiguration().SeparatorColor,
		SeparatorWidth:  0.5,
		ShowSeparators:  false,
		items:           make([]fyne.CanvasObject, 0),
	}
	gv.ExtendBaseWidget(gv)
	return gv
}

// NewGridViewWithSpacing creates a grid view with spacing
func NewGridViewWithSpacing(columnCount int, columnSpacing, rowSpacing float32) *GridView {
	gv := NewGridView(columnCount)
	gv.ColumnSpacing = columnSpacing
	gv.RowSpacing = rowSpacing
	return gv
}

// AddItem adds an item to the grid
func (gv *GridView) AddItem(item fyne.CanvasObject) {
	gv.mu.Lock()
	gv.items = append(gv.items, item)
	gv.mu.Unlock()
	gv.Refresh()
}

// RemoveItem removes an item from the grid
func (gv *GridView) RemoveItem(item fyne.CanvasObject) {
	gv.mu.Lock()
	for i, it := range gv.items {
		if it == item {
			gv.items = append(gv.items[:i], gv.items[i+1:]...)
			break
		}
	}
	gv.mu.Unlock()
	gv.Refresh()
}

// ClearItems removes all items
func (gv *GridView) ClearItems() {
	gv.mu.Lock()
	gv.items = make([]fyne.CanvasObject, 0)
	gv.mu.Unlock()
	gv.Refresh()
}

// SetItems sets all items
func (gv *GridView) SetItems(items []fyne.CanvasObject) {
	gv.mu.Lock()
	gv.items = items
	gv.mu.Unlock()
	gv.Refresh()
}

// ItemCount returns the number of items
func (gv *GridView) ItemCount() int {
	gv.mu.RLock()
	defer gv.mu.RUnlock()
	return len(gv.items)
}

// CreateRenderer implements fyne.Widget
func (gv *GridView) CreateRenderer() fyne.WidgetRenderer {
	gv.ExtendBaseWidget(gv)
	background := canvas.NewRectangle(gv.BackgroundColor)
	return &gridViewRenderer{
		grid:       gv,
		background: background,
		separators: make([]*canvas.Rectangle, 0),
	}
}

type gridViewRenderer struct {
	grid       *GridView
	background *canvas.Rectangle
	separators []*canvas.Rectangle
}

func (r *gridViewRenderer) Destroy() {}

func (r *gridViewRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	r.grid.mu.RLock()
	items := r.grid.items
	columnCount := r.grid.ColumnCount
	columnSpacing := r.grid.ColumnSpacing
	rowSpacing := r.grid.RowSpacing
	rowHeight := r.grid.RowHeight
	insets := r.grid.ContentInsets
	r.grid.mu.RUnlock()

	if len(items) == 0 || columnCount <= 0 {
		return
	}

	availableWidth := size.Width - insets.Left - insets.Right - float32(columnCount-1)*columnSpacing
	columnWidth := availableWidth / float32(columnCount)

	// Calculate row height if auto
	if rowHeight <= 0 {
		var maxHeight float32
		for _, item := range items {
			h := item.MinSize().Height
			if h > maxHeight {
				maxHeight = h
			}
		}
		rowHeight = maxHeight
	}

	// Layout items
	for i, item := range items {
		col := i % columnCount
		row := i / columnCount

		x := insets.Left + float32(col)*(columnWidth+columnSpacing)
		y := insets.Top + float32(row)*(rowHeight+rowSpacing)

		item.Resize(fyne.NewSize(columnWidth, rowHeight))
		item.Move(fyne.NewPos(x, y))
	}

	// Layout separators if needed
	if r.grid.ShowSeparators {
		r.layoutSeparators(size, columnWidth, rowHeight, len(items), columnCount, columnSpacing, rowSpacing, insets)
	}
}

func (r *gridViewRenderer) layoutSeparators(size fyne.Size, columnWidth, rowHeight float32, itemCount, columnCount int, columnSpacing, rowSpacing float32, insets core.EdgeInsets) {
	rowCount := (itemCount + columnCount - 1) / columnCount

	// Horizontal separators
	for row := 0; row < rowCount-1; row++ {
		y := insets.Top + float32(row+1)*(rowHeight+rowSpacing) - rowSpacing/2
		sep := canvas.NewRectangle(r.grid.SeparatorColor)
		sep.Resize(fyne.NewSize(size.Width-insets.Left-insets.Right, r.grid.SeparatorWidth))
		sep.Move(fyne.NewPos(insets.Left, y))
		r.separators = append(r.separators, sep)
	}

	// Vertical separators
	for col := 0; col < columnCount-1; col++ {
		x := insets.Left + float32(col+1)*(columnWidth+columnSpacing) - columnSpacing/2
		sep := canvas.NewRectangle(r.grid.SeparatorColor)
		sep.Resize(fyne.NewSize(r.grid.SeparatorWidth, size.Height-insets.Top-insets.Bottom))
		sep.Move(fyne.NewPos(x, insets.Top))
		r.separators = append(r.separators, sep)
	}
}

func (r *gridViewRenderer) MinSize() fyne.Size {
	r.grid.mu.RLock()
	items := r.grid.items
	columnCount := r.grid.ColumnCount
	columnSpacing := r.grid.ColumnSpacing
	rowSpacing := r.grid.RowSpacing
	rowHeight := r.grid.RowHeight
	insets := r.grid.ContentInsets
	r.grid.mu.RUnlock()

	if len(items) == 0 || columnCount <= 0 {
		return fyne.NewSize(insets.Left+insets.Right, insets.Top+insets.Bottom)
	}

	// Calculate max item width and height
	var maxWidth, maxHeight float32
	for _, item := range items {
		s := item.MinSize()
		if s.Width > maxWidth {
			maxWidth = s.Width
		}
		if s.Height > maxHeight {
			maxHeight = s.Height
		}
	}

	if rowHeight <= 0 {
		rowHeight = maxHeight
	}

	rowCount := (len(items) + columnCount - 1) / columnCount

	width := float32(columnCount)*maxWidth + float32(columnCount-1)*columnSpacing + insets.Left + insets.Right
	height := float32(rowCount)*rowHeight + float32(rowCount-1)*rowSpacing + insets.Top + insets.Bottom

	return fyne.NewSize(width, height)
}

func (r *gridViewRenderer) Refresh() {
	r.background.FillColor = r.grid.BackgroundColor
	r.background.Refresh()

	r.grid.mu.RLock()
	items := r.grid.items
	r.grid.mu.RUnlock()

	for _, item := range items {
		item.Refresh()
	}
}

func (r *gridViewRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}

	r.grid.mu.RLock()
	items := r.grid.items
	r.grid.mu.RUnlock()

	objects = append(objects, items...)

	if r.grid.ShowSeparators {
		for _, sep := range r.separators {
			objects = append(objects, sep)
		}
	}

	return objects
}

// GridViewItem is a wrapper for items in the grid
type GridViewItem struct {
	widget.BaseWidget

	Content         fyne.CanvasObject
	BackgroundColor color.Color
	SelectedColor   color.Color
	CornerRadius    float32
	ContentInsets   core.EdgeInsets

	OnTapped func()

	mu       sync.RWMutex
	selected bool
}

// NewGridViewItem creates a new grid item
func NewGridViewItem(content fyne.CanvasObject) *GridViewItem {
	item := &GridViewItem{
		Content:         content,
		BackgroundColor: color.Transparent,
		SelectedColor:   color.RGBA{R: 0, G: 0, B: 0, A: 20},
		CornerRadius:    0,
		ContentInsets:   core.EdgeInsets{},
	}
	item.ExtendBaseWidget(item)
	return item
}

func (i *GridViewItem) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	background := canvas.NewRectangle(i.BackgroundColor)
	background.CornerRadius = i.CornerRadius
	return &gridItemRenderer{
		item:       i,
		background: background,
	}
}

func (i *GridViewItem) Tapped(_ *fyne.PointEvent) {
	if i.OnTapped != nil {
		i.OnTapped()
	}
}

func (i *GridViewItem) TappedSecondary(_ *fyne.PointEvent) {}

type gridItemRenderer struct {
	item       *GridViewItem
	background *canvas.Rectangle
}

func (r *gridItemRenderer) Destroy() {}

func (r *gridItemRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	if r.item.Content != nil {
		insets := r.item.ContentInsets
		r.item.Content.Resize(fyne.NewSize(
			size.Width-insets.Left-insets.Right,
			size.Height-insets.Top-insets.Bottom,
		))
		r.item.Content.Move(fyne.NewPos(insets.Left, insets.Top))
	}
}

func (r *gridItemRenderer) MinSize() fyne.Size {
	if r.item.Content == nil {
		return fyne.NewSize(0, 0)
	}
	insets := r.item.ContentInsets
	contentSize := r.item.Content.MinSize()
	return fyne.NewSize(
		contentSize.Width+insets.Left+insets.Right,
		contentSize.Height+insets.Top+insets.Bottom,
	)
}

func (r *gridItemRenderer) Refresh() {
	r.item.mu.RLock()
	selected := r.item.selected
	r.item.mu.RUnlock()

	if selected {
		r.background.FillColor = r.item.SelectedColor
	} else {
		r.background.FillColor = r.item.BackgroundColor
	}
	r.background.CornerRadius = r.item.CornerRadius
	r.background.Refresh()

	if r.item.Content != nil {
		r.item.Content.Refresh()
	}
}

func (r *gridItemRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}
	if r.item.Content != nil {
		objects = append(objects, r.item.Content)
	}
	return objects
}
