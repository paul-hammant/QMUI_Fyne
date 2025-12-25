// Package floatlayout provides a flow/wrap layout container
package floatlayout

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// FlowLayout arranges items in a flow layout that wraps to new lines
type FlowLayout struct {
	widget.BaseWidget

	// Layout
	ItemSpacing   float32
	LineSpacing   float32
	ContentInsets core.EdgeInsets
	MaximumWidth  float32

	// Styling
	BackgroundColor color.Color

	// Items
	items []fyne.CanvasObject

	mu sync.RWMutex
}

// NewFlowLayout creates a new float layout view
func NewFlowLayout() *FlowLayout {
	fv := &FlowLayout{
		ItemSpacing:     8,
		LineSpacing:     8,
		ContentInsets:   core.EdgeInsets{},
		MaximumWidth:    0,
		BackgroundColor: color.Transparent,
		items:           make([]fyne.CanvasObject, 0),
	}
	fv.ExtendBaseWidget(fv)
	return fv
}

// NewFlowLayoutWithSpacing creates a float layout with custom spacing
func NewFlowLayoutWithSpacing(itemSpacing, lineSpacing float32) *FlowLayout {
	fv := NewFlowLayout()
	fv.ItemSpacing = itemSpacing
	fv.LineSpacing = lineSpacing
	return fv
}

// AddItem adds an item to the layout
func (fv *FlowLayout) AddItem(item fyne.CanvasObject) {
	fv.mu.Lock()
	fv.items = append(fv.items, item)
	fv.mu.Unlock()
	fv.Refresh()
}

// AddItems adds multiple items to the layout
func (fv *FlowLayout) AddItems(items []fyne.CanvasObject) {
	fv.mu.Lock()
	fv.items = append(fv.items, items...)
	fv.mu.Unlock()
	fv.Refresh()
}

// RemoveItem removes an item from the layout
func (fv *FlowLayout) RemoveItem(item fyne.CanvasObject) {
	fv.mu.Lock()
	for i, it := range fv.items {
		if it == item {
			fv.items = append(fv.items[:i], fv.items[i+1:]...)
			break
		}
	}
	fv.mu.Unlock()
	fv.Refresh()
}

// RemoveAllItems removes all items
func (fv *FlowLayout) RemoveAllItems() {
	fv.mu.Lock()
	fv.items = make([]fyne.CanvasObject, 0)
	fv.mu.Unlock()
	fv.Refresh()
}

// SetItems sets all items
func (fv *FlowLayout) SetItems(items []fyne.CanvasObject) {
	fv.mu.Lock()
	fv.items = items
	fv.mu.Unlock()
	fv.Refresh()
}

// ItemCount returns the number of items
func (fv *FlowLayout) ItemCount() int {
	fv.mu.RLock()
	defer fv.mu.RUnlock()
	return len(fv.items)
}

// CreateRenderer implements fyne.Widget
func (fv *FlowLayout) CreateRenderer() fyne.WidgetRenderer {
	fv.ExtendBaseWidget(fv)
	background := canvas.NewRectangle(fv.BackgroundColor)
	return &floatLayoutRenderer{
		layout:     fv,
		background: background,
	}
}

type floatLayoutRenderer struct {
	layout     *FlowLayout
	background *canvas.Rectangle
}

func (r *floatLayoutRenderer) Destroy() {}

func (r *floatLayoutRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	r.layout.mu.RLock()
	items := r.layout.items
	itemSpacing := r.layout.ItemSpacing
	lineSpacing := r.layout.LineSpacing
	insets := r.layout.ContentInsets
	r.layout.mu.RUnlock()

	if len(items) == 0 {
		return
	}

	availableWidth := size.Width - insets.Left - insets.Right
	x := insets.Left
	y := insets.Top
	var lineHeight float32

	for _, item := range items {
		itemSize := item.MinSize()

		// Check if we need to wrap to next line
		if x+itemSize.Width > insets.Left+availableWidth && x > insets.Left {
			x = insets.Left
			y += lineHeight + lineSpacing
			lineHeight = 0
		}

		item.Resize(itemSize)
		item.Move(fyne.NewPos(x, y))

		x += itemSize.Width + itemSpacing
		if itemSize.Height > lineHeight {
			lineHeight = itemSize.Height
		}
	}
}

func (r *floatLayoutRenderer) MinSize() fyne.Size {
	r.layout.mu.RLock()
	items := r.layout.items
	itemSpacing := r.layout.ItemSpacing
	lineSpacing := r.layout.LineSpacing
	insets := r.layout.ContentInsets
	maxWidth := r.layout.MaximumWidth
	r.layout.mu.RUnlock()

	if len(items) == 0 {
		return fyne.NewSize(insets.Left+insets.Right, insets.Top+insets.Bottom)
	}

	if maxWidth <= 0 {
		// Calculate based on all items in one line
		var totalWidth, maxHeight float32
		for i, item := range items {
			s := item.MinSize()
			totalWidth += s.Width
			if i < len(items)-1 {
				totalWidth += itemSpacing
			}
			if s.Height > maxHeight {
				maxHeight = s.Height
			}
		}
		return fyne.NewSize(
			totalWidth+insets.Left+insets.Right,
			maxHeight+insets.Top+insets.Bottom,
		)
	}

	// Calculate with wrapping
	availableWidth := maxWidth - insets.Left - insets.Right
	x := float32(0)
	y := float32(0)
	var lineHeight float32
	var maxX float32

	for _, item := range items {
		itemSize := item.MinSize()

		if x+itemSize.Width > availableWidth && x > 0 {
			if x > maxX {
				maxX = x
			}
			x = 0
			y += lineHeight + lineSpacing
			lineHeight = 0
		}

		x += itemSize.Width + itemSpacing
		if itemSize.Height > lineHeight {
			lineHeight = itemSize.Height
		}
	}

	if x > maxX {
		maxX = x
	}

	return fyne.NewSize(
		maxX+insets.Left+insets.Right,
		y+lineHeight+insets.Top+insets.Bottom,
	)
}

func (r *floatLayoutRenderer) Refresh() {
	r.background.FillColor = r.layout.BackgroundColor
	r.background.Refresh()

	r.layout.mu.RLock()
	items := r.layout.items
	r.layout.mu.RUnlock()

	for _, item := range items {
		item.Refresh()
	}
}

func (r *floatLayoutRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}

	r.layout.mu.RLock()
	items := r.layout.items
	r.layout.mu.RUnlock()

	return append(objects, items...)
}

// TagCloud is a common use case for FlowLayout - displaying multiple tags
type TagCloud struct {
	*FlowLayout

	// Styling
	TagBackgroundColor color.Color
	TagTextColor       color.Color
	TagFontSize        float32
	TagCornerRadius    float32
	TagPadding         core.EdgeInsets

	// Callbacks
	OnTagTapped func(index int, text string)
}

// NewTagCloud creates a new tag cloud
func NewTagCloud() *TagCloud {
	tc := &TagCloud{
		FlowLayout:         NewFlowLayout(),
		TagBackgroundColor: color.RGBA{R: 230, G: 230, B: 230, A: 255},
		TagTextColor:       color.Black,
		TagFontSize:        14,
		TagCornerRadius:    4,
		TagPadding:         core.NewEdgeInsets(4, 8, 4, 8),
	}
	tc.ItemSpacing = 8
	tc.LineSpacing = 8
	return tc
}

// SetTags sets the tag strings
func (tc *TagCloud) SetTags(tags []string) {
	items := make([]fyne.CanvasObject, len(tags))
	for i, tag := range tags {
		idx := i
		text := tag
		tagWidget := NewTag(tag, tc.TagBackgroundColor, tc.TagTextColor, tc.TagFontSize, tc.TagCornerRadius, tc.TagPadding)
		tagWidget.OnTapped = func() {
			if tc.OnTagTapped != nil {
				tc.OnTagTapped(idx, text)
			}
		}
		items[i] = tagWidget
	}
	tc.SetItems(items)
}

// Tag is a single tag widget
type Tag struct {
	widget.BaseWidget

	Text            string
	BackgroundColor color.Color
	TextColor       color.Color
	FontSize        float32
	CornerRadius    float32
	Padding         core.EdgeInsets
	OnTapped        func()

	mu      sync.RWMutex
	hovered bool
}

// NewSimpleTag creates a new tag with default styling
func NewSimpleTag(text string) *Tag {
	return NewTag(
		text,
		color.RGBA{R: 230, G: 230, B: 230, A: 255},
		color.Black,
		14,
		4,
		core.NewEdgeInsets(4, 8, 4, 8),
	)
}

// NewTag creates a new tag widget
func NewTag(text string, bg, textColor color.Color, fontSize, cornerRadius float32, padding core.EdgeInsets) *Tag {
	t := &Tag{
		Text:            text,
		BackgroundColor: bg,
		TextColor:       textColor,
		FontSize:        fontSize,
		CornerRadius:    cornerRadius,
		Padding:         padding,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *Tag) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	background := canvas.NewRectangle(t.BackgroundColor)
	background.CornerRadius = t.CornerRadius

	text := canvas.NewText(t.Text, t.TextColor)
	text.TextSize = t.FontSize

	return &tagRenderer{
		tag:        t,
		background: background,
		text:       text,
	}
}

func (t *Tag) Tapped(_ *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped()
	}
}

func (t *Tag) TappedSecondary(_ *fyne.PointEvent) {}

type tagRenderer struct {
	tag        *Tag
	background *canvas.Rectangle
	text       *canvas.Text
}

func (r *tagRenderer) Destroy() {}

func (r *tagRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	padding := r.tag.Padding
	r.text.Move(fyne.NewPos(padding.Left, padding.Top))
}

func (r *tagRenderer) MinSize() fyne.Size {
	textSize := r.text.MinSize()
	padding := r.tag.Padding
	return fyne.NewSize(
		textSize.Width+padding.Left+padding.Right,
		textSize.Height+padding.Top+padding.Bottom,
	)
}

func (r *tagRenderer) Refresh() {
	r.background.FillColor = r.tag.BackgroundColor
	r.background.CornerRadius = r.tag.CornerRadius
	r.text.Text = r.tag.Text
	r.text.Color = r.tag.TextColor
	r.text.TextSize = r.tag.FontSize
	r.background.Refresh()
	r.text.Refresh()
}

func (r *tagRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.text}
}
