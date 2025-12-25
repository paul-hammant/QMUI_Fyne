// Package layouter provides QMUILayouter - a linear layout system
// Ported from Tencent's QMUI_iOS framework
package layouter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// LayouterItem wraps a canvas object with layout properties
type LayouterItem struct {
	View            fyne.CanvasObject
	MinimumSize     fyne.Size
	MaximumSize     fyne.Size
	SizeThatFits    fyne.Size
	SpacingBefore   float32
	SpacingAfter    float32
	Visible         bool
	FlexibleSpacing bool // If true, this item's spacing can grow to fill space
}

// NewLayouterItem creates a new layouter item
func NewLayouterItem(view fyne.CanvasObject) *LayouterItem {
	return &LayouterItem{
		View:    view,
		Visible: true,
	}
}

// NewLayouterItemWithSpacing creates a layouter item with spacing
func NewLayouterItemWithSpacing(view fyne.CanvasObject, before, after float32) *LayouterItem {
	item := NewLayouterItem(view)
	item.SpacingBefore = before
	item.SpacingAfter = after
	return item
}

// Layouter is the base interface for layout managers
type Layouter interface {
	AddItem(item *LayouterItem)
	RemoveItem(item *LayouterItem)
	LayoutItems(containerSize fyne.Size)
	CalculateSize() fyne.Size
}

// LinearHorizontalLayouter arranges items horizontally
type LinearHorizontalLayouter struct {
	widget.BaseWidget

	Items         []*LayouterItem
	ContentInsets core.EdgeInsets
	ItemSpacing   float32
	Alignment     VerticalAlignment
}

// VerticalAlignment defines vertical alignment for horizontal layouts
type VerticalAlignment int

const (
	// VerticalAlignmentTop aligns items to top
	VerticalAlignmentTop VerticalAlignment = iota
	// VerticalAlignmentCenter centers items vertically
	VerticalAlignmentCenter
	// VerticalAlignmentBottom aligns items to bottom
	VerticalAlignmentBottom
	// VerticalAlignmentFill stretches items to fill height
	VerticalAlignmentFill
)

// NewLinearHorizontalLayouter creates a horizontal layouter
func NewLinearHorizontalLayouter() *LinearHorizontalLayouter {
	l := &LinearHorizontalLayouter{
		Items:       make([]*LayouterItem, 0),
		ItemSpacing: 0,
		Alignment:   VerticalAlignmentCenter,
	}
	l.ExtendBaseWidget(l)
	return l
}

// AddItem adds an item to the layout
func (l *LinearHorizontalLayouter) AddItem(item *LayouterItem) {
	l.Items = append(l.Items, item)
	l.Refresh()
}

// AddView adds a view with default settings
func (l *LinearHorizontalLayouter) AddView(view fyne.CanvasObject) {
	l.AddItem(NewLayouterItem(view))
}

// AddViewWithSpacing adds a view with spacing
func (l *LinearHorizontalLayouter) AddViewWithSpacing(view fyne.CanvasObject, before, after float32) {
	l.AddItem(NewLayouterItemWithSpacing(view, before, after))
}

// AddFlexibleSpace adds a flexible space that expands
func (l *LinearHorizontalLayouter) AddFlexibleSpace() {
	item := &LayouterItem{
		FlexibleSpacing: true,
		Visible:         true,
	}
	l.Items = append(l.Items, item)
	l.Refresh()
}

// RemoveItem removes an item
func (l *LinearHorizontalLayouter) RemoveItem(item *LayouterItem) {
	for i, it := range l.Items {
		if it == item {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			break
		}
	}
	l.Refresh()
}

// CreateRenderer implements fyne.Widget
func (l *LinearHorizontalLayouter) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	return &horizontalLayouterRenderer{layouter: l}
}

type horizontalLayouterRenderer struct {
	layouter *LinearHorizontalLayouter
}

func (r *horizontalLayouterRenderer) Destroy() {}

func (r *horizontalLayouterRenderer) Layout(size fyne.Size) {
	insets := r.layouter.ContentInsets
	availableWidth := size.Width - insets.Left - insets.Right
	availableHeight := size.Height - insets.Top - insets.Bottom

	// Calculate total fixed width and count flexible items
	var totalFixedWidth float32
	var flexibleCount int
	for _, item := range r.layouter.Items {
		if !item.Visible {
			continue
		}
		if item.FlexibleSpacing {
			flexibleCount++
		} else if item.View != nil {
			itemSize := item.View.MinSize()
			totalFixedWidth += itemSize.Width + item.SpacingBefore + item.SpacingAfter
		}
	}

	// Add item spacing
	visibleCount := 0
	for _, item := range r.layouter.Items {
		if item.Visible {
			visibleCount++
		}
	}
	if visibleCount > 1 {
		totalFixedWidth += r.layouter.ItemSpacing * float32(visibleCount-1)
	}

	// Calculate flexible space width
	var flexWidth float32
	if flexibleCount > 0 {
		flexWidth = (availableWidth - totalFixedWidth) / float32(flexibleCount)
		if flexWidth < 0 {
			flexWidth = 0
		}
	}

	// Layout items
	x := insets.Left
	for _, item := range r.layouter.Items {
		if !item.Visible {
			continue
		}

		x += item.SpacingBefore

		if item.FlexibleSpacing {
			x += flexWidth
		} else if item.View != nil {
			itemSize := item.View.MinSize()

			var y float32
			var itemHeight float32

			switch r.layouter.Alignment {
			case VerticalAlignmentTop:
				y = insets.Top
				itemHeight = itemSize.Height
			case VerticalAlignmentCenter:
				y = insets.Top + (availableHeight-itemSize.Height)/2
				itemHeight = itemSize.Height
			case VerticalAlignmentBottom:
				y = size.Height - insets.Bottom - itemSize.Height
				itemHeight = itemSize.Height
			case VerticalAlignmentFill:
				y = insets.Top
				itemHeight = availableHeight
			}

			item.View.Resize(fyne.NewSize(itemSize.Width, itemHeight))
			item.View.Move(fyne.NewPos(x, y))
			x += itemSize.Width
		}

		x += item.SpacingAfter + r.layouter.ItemSpacing
	}
}

func (r *horizontalLayouterRenderer) MinSize() fyne.Size {
	insets := r.layouter.ContentInsets
	var totalWidth, maxHeight float32

	for i, item := range r.layouter.Items {
		if !item.Visible {
			continue
		}

		totalWidth += item.SpacingBefore + item.SpacingAfter

		if item.View != nil {
			itemSize := item.View.MinSize()
			totalWidth += itemSize.Width
			if itemSize.Height > maxHeight {
				maxHeight = itemSize.Height
			}
		}

		if i < len(r.layouter.Items)-1 {
			totalWidth += r.layouter.ItemSpacing
		}
	}

	return fyne.NewSize(
		totalWidth+insets.Left+insets.Right,
		maxHeight+insets.Top+insets.Bottom,
	)
}

func (r *horizontalLayouterRenderer) Refresh() {
	for _, item := range r.layouter.Items {
		if item.View != nil {
			item.View.Refresh()
		}
	}
}

func (r *horizontalLayouterRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0)
	for _, item := range r.layouter.Items {
		if item.View != nil && item.Visible {
			objects = append(objects, item.View)
		}
	}
	return objects
}

// LinearVerticalLayouter arranges items vertically
type LinearVerticalLayouter struct {
	widget.BaseWidget

	Items         []*LayouterItem
	ContentInsets core.EdgeInsets
	ItemSpacing   float32
	Alignment     HorizontalAlignment
}

// HorizontalAlignment defines horizontal alignment for vertical layouts
type HorizontalAlignment int

const (
	// HorizontalAlignmentLeft aligns items to left
	HorizontalAlignmentLeft HorizontalAlignment = iota
	// HorizontalAlignmentCenter centers items horizontally
	HorizontalAlignmentCenter
	// HorizontalAlignmentRight aligns items to right
	HorizontalAlignmentRight
	// HorizontalAlignmentFill stretches items to fill width
	HorizontalAlignmentFill
)

// NewLinearVerticalLayouter creates a vertical layouter
func NewLinearVerticalLayouter() *LinearVerticalLayouter {
	l := &LinearVerticalLayouter{
		Items:       make([]*LayouterItem, 0),
		ItemSpacing: 0,
		Alignment:   HorizontalAlignmentCenter,
	}
	l.ExtendBaseWidget(l)
	return l
}

// AddItem adds an item to the layout
func (l *LinearVerticalLayouter) AddItem(item *LayouterItem) {
	l.Items = append(l.Items, item)
	l.Refresh()
}

// AddView adds a view with default settings
func (l *LinearVerticalLayouter) AddView(view fyne.CanvasObject) {
	l.AddItem(NewLayouterItem(view))
}

// AddViewWithSpacing adds a view with spacing
func (l *LinearVerticalLayouter) AddViewWithSpacing(view fyne.CanvasObject, before, after float32) {
	l.AddItem(NewLayouterItemWithSpacing(view, before, after))
}

// AddFlexibleSpace adds a flexible space that expands
func (l *LinearVerticalLayouter) AddFlexibleSpace() {
	item := &LayouterItem{
		FlexibleSpacing: true,
		Visible:         true,
	}
	l.Items = append(l.Items, item)
	l.Refresh()
}

// RemoveItem removes an item
func (l *LinearVerticalLayouter) RemoveItem(item *LayouterItem) {
	for i, it := range l.Items {
		if it == item {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			break
		}
	}
	l.Refresh()
}

// CreateRenderer implements fyne.Widget
func (l *LinearVerticalLayouter) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	return &verticalLayouterRenderer{layouter: l}
}

type verticalLayouterRenderer struct {
	layouter *LinearVerticalLayouter
}

func (r *verticalLayouterRenderer) Destroy() {}

func (r *verticalLayouterRenderer) Layout(size fyne.Size) {
	insets := r.layouter.ContentInsets
	availableWidth := size.Width - insets.Left - insets.Right
	availableHeight := size.Height - insets.Top - insets.Bottom

	// Calculate total fixed height and count flexible items
	var totalFixedHeight float32
	var flexibleCount int
	for _, item := range r.layouter.Items {
		if !item.Visible {
			continue
		}
		if item.FlexibleSpacing {
			flexibleCount++
		} else if item.View != nil {
			itemSize := item.View.MinSize()
			totalFixedHeight += itemSize.Height + item.SpacingBefore + item.SpacingAfter
		}
	}

	// Add item spacing
	visibleCount := 0
	for _, item := range r.layouter.Items {
		if item.Visible {
			visibleCount++
		}
	}
	if visibleCount > 1 {
		totalFixedHeight += r.layouter.ItemSpacing * float32(visibleCount-1)
	}

	// Calculate flexible space height
	var flexHeight float32
	if flexibleCount > 0 {
		flexHeight = (availableHeight - totalFixedHeight) / float32(flexibleCount)
		if flexHeight < 0 {
			flexHeight = 0
		}
	}

	// Layout items
	y := insets.Top
	for _, item := range r.layouter.Items {
		if !item.Visible {
			continue
		}

		y += item.SpacingBefore

		if item.FlexibleSpacing {
			y += flexHeight
		} else if item.View != nil {
			itemSize := item.View.MinSize()

			var x float32
			var itemWidth float32

			switch r.layouter.Alignment {
			case HorizontalAlignmentLeft:
				x = insets.Left
				itemWidth = itemSize.Width
			case HorizontalAlignmentCenter:
				x = insets.Left + (availableWidth-itemSize.Width)/2
				itemWidth = itemSize.Width
			case HorizontalAlignmentRight:
				x = size.Width - insets.Right - itemSize.Width
				itemWidth = itemSize.Width
			case HorizontalAlignmentFill:
				x = insets.Left
				itemWidth = availableWidth
			}

			item.View.Resize(fyne.NewSize(itemWidth, itemSize.Height))
			item.View.Move(fyne.NewPos(x, y))
			y += itemSize.Height
		}

		y += item.SpacingAfter + r.layouter.ItemSpacing
	}
}

func (r *verticalLayouterRenderer) MinSize() fyne.Size {
	insets := r.layouter.ContentInsets
	var totalHeight, maxWidth float32

	for i, item := range r.layouter.Items {
		if !item.Visible {
			continue
		}

		totalHeight += item.SpacingBefore + item.SpacingAfter

		if item.View != nil {
			itemSize := item.View.MinSize()
			totalHeight += itemSize.Height
			if itemSize.Width > maxWidth {
				maxWidth = itemSize.Width
			}
		}

		if i < len(r.layouter.Items)-1 {
			totalHeight += r.layouter.ItemSpacing
		}
	}

	return fyne.NewSize(
		maxWidth+insets.Left+insets.Right,
		totalHeight+insets.Top+insets.Bottom,
	)
}

func (r *verticalLayouterRenderer) Refresh() {
	for _, item := range r.layouter.Items {
		if item.View != nil {
			item.View.Refresh()
		}
	}
}

func (r *verticalLayouterRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0)
	for _, item := range r.layouter.Items {
		if item.View != nil && item.Visible {
			objects = append(objects, item.View)
		}
	}
	return objects
}
