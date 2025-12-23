// Package badge provides QMUIBadge - a badge/notification indicator system
// Ported from Tencent's QMUI_iOS framework
package badge

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// BadgeLabel is a badge label component
type BadgeLabel struct {
	widget.BaseWidget

	// Content
	Text string

	// Styling
	BackgroundColor   color.Color
	TextColor         color.Color
	FontSize          float32
	ContentEdgeInsets core.EdgeInsets
	CornerRadius      float32
	MinimumSize       fyne.Size

	mu sync.RWMutex
}

// NewBadgeLabel creates a new badge label
func NewBadgeLabel(text string) *BadgeLabel {
	config := core.SharedConfiguration()
	bl := &BadgeLabel{
		Text:              text,
		BackgroundColor:   config.BadgeBackgroundColor,
		TextColor:         config.BadgeTextColor,
		FontSize:          config.BadgeFontSize,
		ContentEdgeInsets: config.BadgeContentEdgeInsets,
		CornerRadius:      0, // Will be auto-calculated to half height
		MinimumSize:       fyne.NewSize(18, 18),
	}
	bl.ExtendBaseWidget(bl)
	return bl
}

// SetText sets the badge text
func (bl *BadgeLabel) SetText(text string) {
	bl.mu.Lock()
	bl.Text = text
	bl.mu.Unlock()
	bl.Refresh()
}

// CreateRenderer implements fyne.Widget
func (bl *BadgeLabel) CreateRenderer() fyne.WidgetRenderer {
	bl.ExtendBaseWidget(bl)

	background := canvas.NewRectangle(bl.BackgroundColor)
	text := canvas.NewText(bl.Text, bl.TextColor)
	text.TextSize = bl.FontSize
	text.Alignment = fyne.TextAlignCenter

	return &badgeLabelRenderer{
		badge:      bl,
		background: background,
		text:       text,
	}
}

type badgeLabelRenderer struct {
	badge      *BadgeLabel
	background *canvas.Rectangle
	text       *canvas.Text
}

func (r *badgeLabelRenderer) Destroy() {}

func (r *badgeLabelRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.background.CornerRadius = size.Height / 2

	// Resize text to fill available space for proper centering
	r.text.Resize(size)
	r.text.Move(fyne.NewPos(0, 0))
}

func (r *badgeLabelRenderer) MinSize() fyne.Size {
	r.badge.mu.RLock()
	text := r.badge.Text
	r.badge.mu.RUnlock()

	r.text.Text = text
	textSize := r.text.MinSize()
	insets := r.badge.ContentEdgeInsets

	width := textSize.Width + insets.Left + insets.Right
	height := textSize.Height + insets.Top + insets.Bottom

	if width < r.badge.MinimumSize.Width {
		width = r.badge.MinimumSize.Width
	}
	if height < r.badge.MinimumSize.Height {
		height = r.badge.MinimumSize.Height
	}

	return fyne.NewSize(width, height)
}

func (r *badgeLabelRenderer) Refresh() {
	r.badge.mu.RLock()
	text := r.badge.Text
	r.badge.mu.RUnlock()

	r.background.FillColor = r.badge.BackgroundColor
	r.text.Text = text
	r.text.Color = r.badge.TextColor
	r.text.TextSize = r.badge.FontSize

	r.background.Refresh()
	r.text.Refresh()
}

func (r *badgeLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.text}
}

// UpdatesIndicator is a small dot indicator for updates
type UpdatesIndicator struct {
	widget.BaseWidget

	Color         color.Color
	IndicatorSize fyne.Size
	HasUpdates    bool
}

// NewUpdatesIndicator creates a new updates indicator
func NewUpdatesIndicator() *UpdatesIndicator {
	config := core.SharedConfiguration()
	ui := &UpdatesIndicator{
		Color:         config.UpdatesIndicatorColor,
		IndicatorSize: config.UpdatesIndicatorSize,
		HasUpdates:    false,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

// CreateRenderer implements fyne.Widget
func (ui *UpdatesIndicator) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	circle := canvas.NewCircle(ui.Color)
	return &updatesIndicatorRenderer{
		indicator: ui,
		circle:    circle,
	}
}

type updatesIndicatorRenderer struct {
	indicator *UpdatesIndicator
	circle    *canvas.Circle
}

func (r *updatesIndicatorRenderer) Destroy() {}

func (r *updatesIndicatorRenderer) Layout(size fyne.Size) {
	r.circle.Resize(size)
}

func (r *updatesIndicatorRenderer) MinSize() fyne.Size {
	return r.indicator.IndicatorSize
}

func (r *updatesIndicatorRenderer) Refresh() {
	r.circle.FillColor = r.indicator.Color
	r.circle.Refresh()
}

func (r *updatesIndicatorRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.circle}
}

// BadgeView wraps any view with a badge
type BadgeView struct {
	widget.BaseWidget

	Content fyne.CanvasObject

	// Badge properties
	BadgeValue string
	BadgeColor color.Color
	BadgeTextColor color.Color
	BadgeOffset core.Offset

	// Updates indicator
	ShowUpdatesIndicator bool
	UpdatesIndicatorColor color.Color
	UpdatesIndicatorOffset core.Offset

	mu sync.RWMutex
}

// NewBadgeView creates a new badge view wrapping content
func NewBadgeView(content fyne.CanvasObject) *BadgeView {
	config := core.SharedConfiguration()
	bv := &BadgeView{
		Content:              content,
		BadgeColor:          config.BadgeBackgroundColor,
		BadgeTextColor:      config.BadgeTextColor,
		BadgeOffset:         config.BadgeOffset,
		UpdatesIndicatorColor: config.UpdatesIndicatorColor,
		UpdatesIndicatorOffset: config.UpdatesIndicatorOffset,
	}
	bv.ExtendBaseWidget(bv)
	return bv
}

// SetBadgeValue sets the badge text
func (bv *BadgeView) SetBadgeValue(value string) {
	bv.mu.Lock()
	bv.BadgeValue = value
	bv.mu.Unlock()
	bv.Refresh()
}

// ClearBadge removes the badge
func (bv *BadgeView) ClearBadge() {
	bv.SetBadgeValue("")
}

// SetShowUpdatesIndicator shows/hides the updates indicator
func (bv *BadgeView) SetShowUpdatesIndicator(show bool) {
	bv.mu.Lock()
	bv.ShowUpdatesIndicator = show
	bv.mu.Unlock()
	bv.Refresh()
}

// CreateRenderer implements fyne.Widget
func (bv *BadgeView) CreateRenderer() fyne.WidgetRenderer {
	bv.ExtendBaseWidget(bv)

	badge := NewBadgeLabel("")
	badge.BackgroundColor = bv.BadgeColor
	badge.TextColor = bv.BadgeTextColor
	badge.Hide()

	indicator := NewUpdatesIndicator()
	indicator.Color = bv.UpdatesIndicatorColor
	indicator.Hide()

	return &badgeViewRenderer{
		view:      bv,
		badge:     badge,
		indicator: indicator,
	}
}

type badgeViewRenderer struct {
	view      *BadgeView
	badge     *BadgeLabel
	indicator *UpdatesIndicator
}

func (r *badgeViewRenderer) Destroy() {}

func (r *badgeViewRenderer) Layout(size fyne.Size) {
	if r.view.Content != nil {
		r.view.Content.Resize(size)
	}

	r.view.mu.RLock()
	badgeValue := r.view.BadgeValue
	showIndicator := r.view.ShowUpdatesIndicator
	badgeOffset := r.view.BadgeOffset
	indicatorOffset := r.view.UpdatesIndicatorOffset
	r.view.mu.RUnlock()

	// Position badge at top-right
	if badgeValue != "" {
		badgeSize := r.badge.MinSize()
		badgeX := size.Width + badgeOffset.X - badgeSize.Width/2
		badgeY := badgeOffset.Y - badgeSize.Height/2
		r.badge.Resize(badgeSize)
		r.badge.Move(fyne.NewPos(badgeX, badgeY))
		r.badge.Show()
	} else {
		r.badge.Hide()
	}

	// Position indicator
	if showIndicator && badgeValue == "" {
		indicatorSize := r.indicator.MinSize()
		indicatorX := size.Width + indicatorOffset.X - indicatorSize.Width/2
		indicatorY := indicatorOffset.Y - indicatorSize.Height/2
		r.indicator.Resize(indicatorSize)
		r.indicator.Move(fyne.NewPos(indicatorX, indicatorY))
		r.indicator.Show()
	} else {
		r.indicator.Hide()
	}
}

func (r *badgeViewRenderer) MinSize() fyne.Size {
	if r.view.Content != nil {
		return r.view.Content.MinSize()
	}
	return fyne.NewSize(0, 0)
}

func (r *badgeViewRenderer) Refresh() {
	r.view.mu.RLock()
	badgeValue := r.view.BadgeValue
	r.view.mu.RUnlock()

	r.badge.SetText(badgeValue)
	r.badge.BackgroundColor = r.view.BadgeColor
	r.badge.TextColor = r.view.BadgeTextColor
	r.indicator.Color = r.view.UpdatesIndicatorColor

	if r.view.Content != nil {
		r.view.Content.Refresh()
	}
	r.badge.Refresh()
	r.indicator.Refresh()
}

func (r *badgeViewRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{}
	if r.view.Content != nil {
		objects = append(objects, r.view.Content)
	}
	objects = append(objects, r.badge, r.indicator)
	return objects
}

// TabBarBadge adds a badge to tab bar items
type TabBarBadge struct {
	value     string
	showsDot  bool
}

// NewTabBarBadge creates a tab bar badge
func NewTabBarBadge(value string) *TabBarBadge {
	return &TabBarBadge{value: value}
}

// NewTabBarDotBadge creates a tab bar dot badge
func NewTabBarDotBadge() *TabBarBadge {
	return &TabBarBadge{showsDot: true}
}
