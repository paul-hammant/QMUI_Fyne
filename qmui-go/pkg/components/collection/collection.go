// Package collection provides QMUICollectionViewPagingLayout - paging collection view
// Ported from Tencent's QMUI_iOS framework
package collection

import (
	"image/color"
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/animation"
	"github.com/user/qmui-go/pkg/core"
)

// PagingStyle defines the paging behavior
type PagingStyle int

const (
	PagingStyleDefault PagingStyle = iota
	PagingStyleScale
	PagingStyleCoverFlow
)

// PagingLayout provides a horizontally scrolling collection with paging
type PagingLayout struct {
	widget.BaseWidget

	// Content
	Items           []fyne.CanvasObject
	CurrentPage     int

	// Layout
	ItemSize        fyne.Size
	ItemSpacing     float32
	PageInsets      core.EdgeInsets
	PagingStyle     PagingStyle

	// Scale style options (when PagingStyle == PagingStyleScale)
	MinimumScale       float32
	MaximumScale       float32
	ScaleInterpolation float32

	// Paging
	AllowsMultiplePages bool
	PageIndicatorEnabled bool
	PageIndicatorPosition float32 // Distance from bottom

	// Styling
	PageIndicatorColor       color.Color
	PageIndicatorActiveColor color.Color
	PageIndicatorSize        float32
	PageIndicatorSpacing     float32
	BackgroundColor          color.Color

	// Animation
	AnimationDuration time.Duration
	AnimationEasing   animation.EasingFunction

	// Callbacks
	OnPageChanged func(page int)
	OnItemTapped  func(index int)

	// State
	mu           sync.RWMutex
	offsetX      float32
	isDragging   bool
	lastDragPos  fyne.Position
	dragVelocity float32
}

// NewPagingLayout creates a new paging layout
func NewPagingLayout() *PagingLayout {
	config := core.SharedConfiguration()
	pl := &PagingLayout{
		Items:                 make([]fyne.CanvasObject, 0),
		CurrentPage:           0,
		ItemSize:              fyne.NewSize(280, 400),
		ItemSpacing:           16,
		PageInsets:            core.NewEdgeInsets(0, 20, 0, 20),
		PagingStyle:           PagingStyleDefault,
		MinimumScale:          0.8,
		MaximumScale:          1.0,
		ScaleInterpolation:    0.5,
		AllowsMultiplePages:   false,
		PageIndicatorEnabled:  true,
		PageIndicatorPosition: 20,
		PageIndicatorColor:    config.GrayLightenColor,
		PageIndicatorActiveColor: config.BlueColor,
		PageIndicatorSize:     8,
		PageIndicatorSpacing:  8,
		BackgroundColor:       color.Transparent,
		AnimationDuration:     time.Millisecond * 300,
		AnimationEasing:       animation.EaseOutCubic,
	}
	pl.ExtendBaseWidget(pl)
	return pl
}

// NewPagingLayoutWithItems creates a paging layout with items
func NewPagingLayoutWithItems(items []fyne.CanvasObject) *PagingLayout {
	pl := NewPagingLayout()
	pl.Items = items
	return pl
}

// SetItems sets the items
func (pl *PagingLayout) SetItems(items []fyne.CanvasObject) {
	pl.mu.Lock()
	pl.Items = items
	if pl.CurrentPage >= len(items) {
		pl.CurrentPage = len(items) - 1
	}
	if pl.CurrentPage < 0 {
		pl.CurrentPage = 0
	}
	pl.mu.Unlock()
	pl.Refresh()
}

// AddItem adds an item
func (pl *PagingLayout) AddItem(item fyne.CanvasObject) {
	pl.mu.Lock()
	pl.Items = append(pl.Items, item)
	pl.mu.Unlock()
	pl.Refresh()
}

// AddPage adds a page (alias for AddItem)
func (pl *PagingLayout) AddPage(page fyne.CanvasObject) {
	pl.AddItem(page)
}

// SetCurrentPage sets the current page immediately without animation
func (pl *PagingLayout) SetCurrentPage(page int) {
	pl.mu.Lock()
	if page < 0 || page >= len(pl.Items) {
		pl.mu.Unlock()
		return
	}
	oldPage := pl.CurrentPage
	pl.CurrentPage = page
	pl.offsetX = pl.calculateOffsetForPage(page)
	pl.mu.Unlock()

	pl.Refresh()

	if page != oldPage && pl.OnPageChanged != nil {
		pl.OnPageChanged(page)
	}
}

// RemoveItem removes an item at index
func (pl *PagingLayout) RemoveItem(index int) {
	pl.mu.Lock()
	if index < 0 || index >= len(pl.Items) {
		pl.mu.Unlock()
		return
	}
	pl.Items = append(pl.Items[:index], pl.Items[index+1:]...)
	if pl.CurrentPage >= len(pl.Items) {
		pl.CurrentPage = len(pl.Items) - 1
	}
	pl.mu.Unlock()
	pl.Refresh()
}

// GoToPage animates to a specific page
func (pl *PagingLayout) GoToPage(page int) {
	pl.mu.Lock()
	if page < 0 || page >= len(pl.Items) {
		pl.mu.Unlock()
		return
	}
	oldPage := pl.CurrentPage
	pl.CurrentPage = page
	pl.mu.Unlock()

	// Animate to new page
	pl.animateToPage(page)

	if page != oldPage && pl.OnPageChanged != nil {
		pl.OnPageChanged(page)
	}
}

// NextPage goes to the next page
func (pl *PagingLayout) NextPage() {
	pl.mu.RLock()
	current := pl.CurrentPage
	count := len(pl.Items)
	pl.mu.RUnlock()

	if current < count-1 {
		pl.GoToPage(current + 1)
	}
}

// PreviousPage goes to the previous page
func (pl *PagingLayout) PreviousPage() {
	pl.mu.RLock()
	current := pl.CurrentPage
	pl.mu.RUnlock()

	if current > 0 {
		pl.GoToPage(current - 1)
	}
}

// GetPageCount returns the total number of pages
func (pl *PagingLayout) GetPageCount() int {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	return len(pl.Items)
}

func (pl *PagingLayout) animateToPage(page int) {
	targetOffset := pl.calculateOffsetForPage(page)
	currentOffset := pl.offsetX

	animation.NewPropertyAnimation(
		float64(currentOffset),
		float64(targetOffset),
		pl.AnimationDuration,
		pl.AnimationEasing,
		func(value float64) {
			pl.mu.Lock()
			pl.offsetX = float32(value)
			pl.mu.Unlock()
			pl.Refresh()
		},
	).Start()
}

func (pl *PagingLayout) calculateOffsetForPage(page int) float32 {
	// Each page is at position: page * (itemSize.Width + spacing)
	pageWidth := pl.ItemSize.Width + pl.ItemSpacing
	return float32(page) * pageWidth
}

// Dragged implements fyne.Draggable
func (pl *PagingLayout) Dragged(e *fyne.DragEvent) {
	pl.mu.Lock()
	pl.isDragging = true
	pl.offsetX -= e.Dragged.DX
	pl.dragVelocity = e.Dragged.DX
	pl.mu.Unlock()
	pl.Refresh()
}

// DragEnd implements fyne.Draggable
func (pl *PagingLayout) DragEnd() {
	pl.mu.Lock()
	pl.isDragging = false
	velocity := pl.dragVelocity
	currentOffset := pl.offsetX
	pl.mu.Unlock()

	// Determine target page based on position and velocity
	pageWidth := pl.ItemSize.Width + pl.ItemSpacing
	currentPageFloat := currentOffset / pageWidth

	var targetPage int
	if math.Abs(float64(velocity)) > 5 {
		// Use velocity to determine direction
		if velocity < 0 {
			targetPage = int(math.Ceil(float64(currentPageFloat)))
		} else {
			targetPage = int(math.Floor(float64(currentPageFloat)))
		}
	} else {
		// Snap to nearest page
		targetPage = int(math.Round(float64(currentPageFloat)))
	}

	// Clamp to valid range
	pl.mu.RLock()
	itemCount := len(pl.Items)
	pl.mu.RUnlock()

	if targetPage < 0 {
		targetPage = 0
	}
	if targetPage >= itemCount {
		targetPage = itemCount - 1
	}

	pl.GoToPage(targetPage)
}

// Scrolled implements fyne.Scrollable for mouse wheel support
func (pl *PagingLayout) Scrolled(e *fyne.ScrollEvent) {
	if e.Scrolled.DX > 0 || e.Scrolled.DY > 0 {
		pl.PreviousPage()
	} else if e.Scrolled.DX < 0 || e.Scrolled.DY < 0 {
		pl.NextPage()
	}
}

// CreateRenderer implements fyne.Widget
func (pl *PagingLayout) CreateRenderer() fyne.WidgetRenderer {
	pl.ExtendBaseWidget(pl)

	background := canvas.NewRectangle(pl.BackgroundColor)

	return &pagingLayoutRenderer{
		layout:     pl,
		background: background,
		itemViews:  make([]*pagingItemView, 0),
	}
}

type pagingLayoutRenderer struct {
	layout        *PagingLayout
	background    *canvas.Rectangle
	itemViews     []*pagingItemView
	pageIndicators []*canvas.Circle
}

func (r *pagingLayoutRenderer) Destroy() {}

func (r *pagingLayoutRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.rebuildItems(size)
	r.buildPageIndicators(size)
}

func (r *pagingLayoutRenderer) rebuildItems(size fyne.Size) {
	r.layout.mu.RLock()
	items := r.layout.Items
	itemSize := r.layout.ItemSize
	spacing := r.layout.ItemSpacing
	offsetX := r.layout.offsetX
	style := r.layout.PagingStyle
	minScale := r.layout.MinimumScale
	maxScale := r.layout.MaximumScale
	currentPage := r.layout.CurrentPage
	r.layout.mu.RUnlock()

	// Ensure we have the right number of item views
	if len(r.itemViews) != len(items) {
		r.itemViews = make([]*pagingItemView, len(items))
		for i, item := range items {
			iv := &pagingItemView{
				layout:  r.layout,
				index:   i,
				content: item,
			}
			iv.ExtendBaseWidget(iv)
			r.itemViews[i] = iv
		}
	}

	// Calculate center position
	centerX := size.Width / 2

	// Position each item
	for i, iv := range r.itemViews {
		// Calculate item position
		itemX := float32(i)*(itemSize.Width+spacing) - offsetX + (size.Width-itemSize.Width)/2
		itemY := (size.Height - itemSize.Height) / 2

		// Apply scaling based on style
		scale := float32(1.0)
		if style == PagingStyleScale || style == PagingStyleCoverFlow {
			// Calculate distance from center
			itemCenterX := itemX + itemSize.Width/2
			distanceFromCenter := math.Abs(float64(itemCenterX - centerX))
			maxDistance := float64(itemSize.Width + spacing)

			// Scale based on distance
			normalizedDistance := distanceFromCenter / maxDistance
			if normalizedDistance > 1 {
				normalizedDistance = 1
			}
			scale = maxScale - (maxScale-minScale)*float32(normalizedDistance)
		}

		// Apply scale
		scaledWidth := itemSize.Width * scale
		scaledHeight := itemSize.Height * scale

		// Adjust position for scale
		scaledX := itemX + (itemSize.Width-scaledWidth)/2
		scaledY := itemY + (itemSize.Height-scaledHeight)/2

		iv.Move(fyne.NewPos(scaledX, scaledY))
		iv.Resize(fyne.NewSize(scaledWidth, scaledHeight))

		// Update content size
		if iv.content != nil {
			iv.content.Resize(fyne.NewSize(scaledWidth, scaledHeight))
		}
	}

	_ = currentPage // avoid unused warning
}

func (r *pagingLayoutRenderer) buildPageIndicators(size fyne.Size) {
	r.layout.mu.RLock()
	enabled := r.layout.PageIndicatorEnabled
	itemCount := len(r.layout.Items)
	currentPage := r.layout.CurrentPage
	indicatorSize := r.layout.PageIndicatorSize
	indicatorSpacing := r.layout.PageIndicatorSpacing
	indicatorPos := r.layout.PageIndicatorPosition
	activeColor := r.layout.PageIndicatorActiveColor
	inactiveColor := r.layout.PageIndicatorColor
	r.layout.mu.RUnlock()

	if !enabled || itemCount == 0 {
		r.pageIndicators = nil
		return
	}

	// Rebuild indicators if count changed
	if len(r.pageIndicators) != itemCount {
		r.pageIndicators = make([]*canvas.Circle, itemCount)
		for i := 0; i < itemCount; i++ {
			indicator := canvas.NewCircle(inactiveColor)
			r.pageIndicators[i] = indicator
		}
	}

	// Calculate total width of indicators
	totalWidth := float32(itemCount)*indicatorSize + float32(itemCount-1)*indicatorSpacing
	startX := (size.Width - totalWidth) / 2
	y := size.Height - indicatorPos - indicatorSize

	for i, indicator := range r.pageIndicators {
		x := startX + float32(i)*(indicatorSize+indicatorSpacing)
		indicator.Move(fyne.NewPos(x, y))
		indicator.Resize(fyne.NewSize(indicatorSize, indicatorSize))

		if i == currentPage {
			indicator.FillColor = activeColor
		} else {
			indicator.FillColor = inactiveColor
		}
	}
}

func (r *pagingLayoutRenderer) MinSize() fyne.Size {
	r.layout.mu.RLock()
	itemSize := r.layout.ItemSize
	insets := r.layout.PageInsets
	r.layout.mu.RUnlock()

	return fyne.NewSize(
		itemSize.Width+insets.Left+insets.Right,
		itemSize.Height+insets.Top+insets.Bottom+30, // Extra space for page indicator
	)
}

func (r *pagingLayoutRenderer) Refresh() {
	r.background.FillColor = r.layout.BackgroundColor
	r.background.Refresh()
	r.rebuildItems(r.layout.Size())
	r.buildPageIndicators(r.layout.Size())

	for _, iv := range r.itemViews {
		iv.Refresh()
	}
	for _, indicator := range r.pageIndicators {
		indicator.Refresh()
	}
}

func (r *pagingLayoutRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}
	for _, iv := range r.itemViews {
		objects = append(objects, iv)
	}
	for _, indicator := range r.pageIndicators {
		objects = append(objects, indicator)
	}
	return objects
}

// pagingItemView wraps an item in the paging layout
type pagingItemView struct {
	widget.BaseWidget
	layout  *PagingLayout
	index   int
	content fyne.CanvasObject
}

func (v *pagingItemView) CreateRenderer() fyne.WidgetRenderer {
	v.ExtendBaseWidget(v)
	return &pagingItemRenderer{view: v}
}

func (v *pagingItemView) Tapped(*fyne.PointEvent) {
	if v.layout.OnItemTapped != nil {
		v.layout.OnItemTapped(v.index)
	}
}

func (v *pagingItemView) TappedSecondary(*fyne.PointEvent) {}

func (v *pagingItemView) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type pagingItemRenderer struct {
	view *pagingItemView
}

func (r *pagingItemRenderer) Destroy() {}

func (r *pagingItemRenderer) Layout(size fyne.Size) {
	if r.view.content != nil {
		r.view.content.Resize(size)
	}
}

func (r *pagingItemRenderer) MinSize() fyne.Size {
	if r.view.content != nil {
		return r.view.content.MinSize()
	}
	return fyne.NewSize(0, 0)
}

func (r *pagingItemRenderer) Refresh() {
	if r.view.content != nil {
		r.view.content.Refresh()
	}
}

func (r *pagingItemRenderer) Objects() []fyne.CanvasObject {
	if r.view.content != nil {
		return []fyne.CanvasObject{r.view.content}
	}
	return nil
}

// Helper functions

// NewCardPagingLayout creates a paging layout styled for cards
func NewCardPagingLayout() *PagingLayout {
	pl := NewPagingLayout()
	pl.PagingStyle = PagingStyleScale
	pl.ItemSize = fyne.NewSize(280, 400)
	pl.ItemSpacing = 16
	pl.MinimumScale = 0.85
	pl.MaximumScale = 1.0
	return pl
}

// NewCoverFlowLayout creates a coverflow-style paging layout
func NewCoverFlowLayout() *PagingLayout {
	pl := NewPagingLayout()
	pl.PagingStyle = PagingStyleCoverFlow
	pl.ItemSize = fyne.NewSize(200, 280)
	pl.ItemSpacing = -40 // Overlapping items
	pl.MinimumScale = 0.7
	pl.MaximumScale = 1.0
	return pl
}

// PagedContainer is a simple container that pages through child views
type PagedContainer struct {
	widget.BaseWidget

	Pages       []fyne.CanvasObject
	CurrentPage int

	OnPageChanged func(page int)

	mu sync.RWMutex
}

// NewPagedContainer creates a new paged container
func NewPagedContainer(pages ...fyne.CanvasObject) *PagedContainer {
	pc := &PagedContainer{
		Pages:       pages,
		CurrentPage: 0,
	}
	pc.ExtendBaseWidget(pc)
	return pc
}

// SetPage sets the current page
func (pc *PagedContainer) SetPage(page int) {
	pc.mu.Lock()
	if page < 0 || page >= len(pc.Pages) {
		pc.mu.Unlock()
		return
	}
	oldPage := pc.CurrentPage
	pc.CurrentPage = page
	pc.mu.Unlock()

	pc.Refresh()

	if page != oldPage && pc.OnPageChanged != nil {
		pc.OnPageChanged(page)
	}
}

// NextPage goes to the next page
func (pc *PagedContainer) NextPage() {
	pc.mu.RLock()
	current := pc.CurrentPage
	count := len(pc.Pages)
	pc.mu.RUnlock()

	if current < count-1 {
		pc.SetPage(current + 1)
	}
}

// PreviousPage goes to the previous page
func (pc *PagedContainer) PreviousPage() {
	pc.mu.RLock()
	current := pc.CurrentPage
	pc.mu.RUnlock()

	if current > 0 {
		pc.SetPage(current - 1)
	}
}

// CreateRenderer implements fyne.Widget
func (pc *PagedContainer) CreateRenderer() fyne.WidgetRenderer {
	pc.ExtendBaseWidget(pc)
	return &pagedContainerRenderer{container: pc}
}

type pagedContainerRenderer struct {
	container *PagedContainer
}

func (r *pagedContainerRenderer) Destroy() {}

func (r *pagedContainerRenderer) Layout(size fyne.Size) {
	r.container.mu.RLock()
	pages := r.container.Pages
	current := r.container.CurrentPage
	r.container.mu.RUnlock()

	for i, page := range pages {
		if i == current {
			page.Resize(size)
			page.Move(fyne.NewPos(0, 0))
			page.Show()
		} else {
			page.Hide()
		}
	}
}

func (r *pagedContainerRenderer) MinSize() fyne.Size {
	r.container.mu.RLock()
	pages := r.container.Pages
	current := r.container.CurrentPage
	r.container.mu.RUnlock()

	if current >= 0 && current < len(pages) {
		return pages[current].MinSize()
	}
	return fyne.NewSize(0, 0)
}

func (r *pagedContainerRenderer) Refresh() {
	r.container.mu.RLock()
	pages := r.container.Pages
	current := r.container.CurrentPage
	r.container.mu.RUnlock()

	for i, page := range pages {
		if i == current {
			page.Show()
			page.Refresh()
		} else {
			page.Hide()
		}
	}
}

func (r *pagedContainerRenderer) Objects() []fyne.CanvasObject {
	r.container.mu.RLock()
	defer r.container.mu.RUnlock()
	return r.container.Pages
}
