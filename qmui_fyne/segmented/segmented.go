// Package segmented provides QMUISegmentedControl - a styled segmented control
// Ported from Tencent's QMUI_iOS framework
package segmented

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/core"
)

// SegmentedControl is a horizontal control with multiple segments
type SegmentedControl struct {
	widget.BaseWidget

	// Segments
	Segments []string
	SelectedIndex int

	// Styling
	TintColor             color.Color
	BackgroundColor       color.Color
	SelectedBackgroundColor color.Color
	TextColor             color.Color
	SelectedTextColor     color.Color
	BorderColor           color.Color
	BorderWidth           float32
	CornerRadius          float32
	ContentEdgeInsets     core.EdgeInsets
	SegmentSpacing        float32
	TextSize              float32

	// Callbacks
	OnValueChanged func(selectedIndex int)

	mu          sync.RWMutex
	hoveredIndex int
}

// NewSegmentedControl creates a new segmented control
func NewSegmentedControl(segments []string, onValueChanged func(selectedIndex int)) *SegmentedControl {
	config := core.SharedConfiguration()
	sc := &SegmentedControl{
		Segments:              segments,
		SelectedIndex:         0,
		TintColor:             config.BlueColor,
		BackgroundColor:       color.Transparent,
		SelectedBackgroundColor: config.BlueColor,
		TextColor:             config.BlueColor,
		SelectedTextColor:     color.White,
		BorderColor:           config.BlueColor,
		BorderWidth:           1,
		CornerRadius:          4,
		ContentEdgeInsets:     core.NewEdgeInsets(6, 12, 6, 12),
		SegmentSpacing:        0,
		TextSize:              theme.TextSize(),
		OnValueChanged:        onValueChanged,
		hoveredIndex:          -1,
	}
	sc.ExtendBaseWidget(sc)
	return sc
}

// SetSelectedIndex sets the selected segment index
func (sc *SegmentedControl) SetSelectedIndex(index int) {
	if index < 0 || index >= len(sc.Segments) {
		return
	}
	sc.mu.Lock()
	sc.SelectedIndex = index
	sc.mu.Unlock()
	sc.Refresh()
	if sc.OnValueChanged != nil {
		sc.OnValueChanged(index)
	}
}

// GetSelectedIndex returns the currently selected index
func (sc *SegmentedControl) GetSelectedIndex() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.SelectedIndex
}

// GetSelectedSegment returns the currently selected segment text
func (sc *SegmentedControl) GetSelectedSegment() string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if sc.SelectedIndex >= 0 && sc.SelectedIndex < len(sc.Segments) {
		return sc.Segments[sc.SelectedIndex]
	}
	return ""
}

// InsertSegment adds a new segment at the specified index
func (sc *SegmentedControl) InsertSegment(title string, index int) {
	sc.mu.Lock()
	if index < 0 {
		index = 0
	}
	if index > len(sc.Segments) {
		index = len(sc.Segments)
	}
	sc.Segments = append(sc.Segments[:index], append([]string{title}, sc.Segments[index:]...)...)
	sc.mu.Unlock()
	sc.Refresh()
}

// RemoveSegment removes the segment at the specified index
func (sc *SegmentedControl) RemoveSegment(index int) {
	sc.mu.Lock()
	if index < 0 || index >= len(sc.Segments) {
		sc.mu.Unlock()
		return
	}
	sc.Segments = append(sc.Segments[:index], sc.Segments[index+1:]...)
	if sc.SelectedIndex >= len(sc.Segments) {
		sc.SelectedIndex = len(sc.Segments) - 1
	}
	sc.mu.Unlock()
	sc.Refresh()
}

// CreateRenderer implements fyne.Widget
func (sc *SegmentedControl) CreateRenderer() fyne.WidgetRenderer {
	sc.ExtendBaseWidget(sc)

	background := canvas.NewRectangle(sc.BackgroundColor)
	background.CornerRadius = sc.CornerRadius
	background.StrokeWidth = sc.BorderWidth
	background.StrokeColor = sc.BorderColor

	return &segmentedRenderer{
		control:    sc,
		background: background,
		segments:   make([]*segmentItem, 0),
	}
}

type segmentItem struct {
	background *canvas.Rectangle
	label      *canvas.Text
	separator  *canvas.Rectangle
}

type segmentedRenderer struct {
	control    *SegmentedControl
	background *canvas.Rectangle
	segments   []*segmentItem
}

func (r *segmentedRenderer) Destroy() {}

func (r *segmentedRenderer) updateSegments() {
	// Rebuild segment items if count changed
	if len(r.segments) != len(r.control.Segments) {
		r.segments = make([]*segmentItem, len(r.control.Segments))
		for i := range r.control.Segments {
			bg := canvas.NewRectangle(color.Transparent)
			label := canvas.NewText(r.control.Segments[i], r.control.TextColor)
			label.TextSize = r.control.TextSize
			label.Alignment = fyne.TextAlignCenter
			sep := canvas.NewRectangle(r.control.BorderColor)
			r.segments[i] = &segmentItem{
				background: bg,
				label:      label,
				separator:  sep,
			}
		}
	}
}

func (r *segmentedRenderer) Layout(size fyne.Size) {
	r.updateSegments()

	r.background.Resize(size)
	r.background.Move(fyne.NewPos(0, 0))

	if len(r.segments) == 0 {
		return
	}

	segmentWidth := size.Width / float32(len(r.segments))
	insets := r.control.ContentEdgeInsets

	for i, seg := range r.segments {
		x := float32(i) * segmentWidth

		// Background for segment
		seg.background.Resize(fyne.NewSize(segmentWidth, size.Height))
		seg.background.Move(fyne.NewPos(x, 0))

		// Apply corner radius to first and last segments
		if i == 0 {
			seg.background.CornerRadius = r.control.CornerRadius
		} else if i == len(r.segments)-1 {
			seg.background.CornerRadius = r.control.CornerRadius
		} else {
			seg.background.CornerRadius = 0
		}

		// Label
		labelSize := seg.label.MinSize()
		labelX := x + (segmentWidth-labelSize.Width)/2
		labelY := insets.Top + (size.Height-insets.Top-insets.Bottom-labelSize.Height)/2
		seg.label.Move(fyne.NewPos(labelX, labelY))

		// Separator (except for last segment)
		if i < len(r.segments)-1 {
			seg.separator.Resize(fyne.NewSize(r.control.BorderWidth, size.Height*0.6))
			seg.separator.Move(fyne.NewPos(x+segmentWidth-r.control.BorderWidth/2, size.Height*0.2))
		}
	}
}

func (r *segmentedRenderer) MinSize() fyne.Size {
	r.updateSegments()

	var maxWidth, maxHeight float32
	insets := r.control.ContentEdgeInsets

	for _, seg := range r.segments {
		labelSize := seg.label.MinSize()
		width := labelSize.Width + insets.Left + insets.Right
		height := labelSize.Height + insets.Top + insets.Bottom
		if width > maxWidth {
			maxWidth = width
		}
		if height > maxHeight {
			maxHeight = height
		}
	}

	totalWidth := maxWidth * float32(len(r.segments))
	return fyne.NewSize(totalWidth, maxHeight)
}

func (r *segmentedRenderer) Refresh() {
	r.updateSegments()

	r.background.FillColor = r.control.BackgroundColor
	r.background.StrokeColor = r.control.BorderColor
	r.background.StrokeWidth = r.control.BorderWidth
	r.background.CornerRadius = r.control.CornerRadius

	r.control.mu.RLock()
	selectedIndex := r.control.SelectedIndex
	hoveredIndex := r.control.hoveredIndex
	r.control.mu.RUnlock()

	for i, seg := range r.segments {
		if i == selectedIndex {
			seg.background.FillColor = r.control.SelectedBackgroundColor
			seg.label.Color = r.control.SelectedTextColor
		} else if i == hoveredIndex {
			seg.background.FillColor = core.ColorWithAlpha(r.control.SelectedBackgroundColor, 0.3)
			seg.label.Color = r.control.TextColor
		} else {
			seg.background.FillColor = color.Transparent
			seg.label.Color = r.control.TextColor
		}

		seg.label.Text = r.control.Segments[i]
		seg.label.TextSize = r.control.TextSize

		if i < len(r.segments)-1 {
			if i == selectedIndex || i+1 == selectedIndex {
				seg.separator.Hide()
			} else {
				seg.separator.FillColor = r.control.BorderColor
				seg.separator.Show()
			}
		}

		seg.background.Refresh()
		seg.label.Refresh()
		seg.separator.Refresh()
	}

	r.background.Refresh()
}

func (r *segmentedRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}
	for _, seg := range r.segments {
		objects = append(objects, seg.background, seg.label, seg.separator)
	}
	return objects
}

// Tapped handles tap events
func (sc *SegmentedControl) Tapped(e *fyne.PointEvent) {
	index := sc.indexAtPosition(e.Position)
	if index >= 0 && index < len(sc.Segments) {
		sc.SetSelectedIndex(index)
	}
}

// TappedSecondary handles secondary tap
func (sc *SegmentedControl) TappedSecondary(_ *fyne.PointEvent) {}

// MouseIn handles mouse enter
func (sc *SegmentedControl) MouseIn(e *desktop.MouseEvent) {
	index := sc.indexAtPosition(e.Position)
	sc.mu.Lock()
	sc.hoveredIndex = index
	sc.mu.Unlock()
	sc.Refresh()
}

// MouseMoved handles mouse movement
func (sc *SegmentedControl) MouseMoved(e *desktop.MouseEvent) {
	index := sc.indexAtPosition(e.Position)
	sc.mu.Lock()
	sc.hoveredIndex = index
	sc.mu.Unlock()
	sc.Refresh()
}

// MouseOut handles mouse leave
func (sc *SegmentedControl) MouseOut() {
	sc.mu.Lock()
	sc.hoveredIndex = -1
	sc.mu.Unlock()
	sc.Refresh()
}

func (sc *SegmentedControl) indexAtPosition(pos fyne.Position) int {
	if len(sc.Segments) == 0 {
		return -1
	}
	segmentWidth := sc.Size().Width / float32(len(sc.Segments))
	index := int(pos.X / segmentWidth)
	if index < 0 {
		index = 0
	}
	if index >= len(sc.Segments) {
		index = len(sc.Segments) - 1
	}
	return index
}

// Cursor returns the cursor for this widget
func (sc *SegmentedControl) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// PillSegmentedControl is a segmented control with pill-shaped selection
type PillSegmentedControl struct {
	*SegmentedControl
}

// NewPillSegmentedControl creates a pill-style segmented control
func NewPillSegmentedControl(segments []string, onValueChanged func(selectedIndex int)) *PillSegmentedControl {
	sc := NewSegmentedControl(segments, onValueChanged)
	sc.CornerRadius = 16
	sc.BorderWidth = 0
	sc.BackgroundColor = color.RGBA{R: 229, G: 229, B: 234, A: 255}
	sc.SelectedBackgroundColor = color.White
	sc.TextColor = color.Black
	sc.SelectedTextColor = color.Black
	return &PillSegmentedControl{SegmentedControl: sc}
}

// UnderlineSegmentedControl uses an underline for selection indication
type UnderlineSegmentedControl struct {
	*SegmentedControl
	UnderlineHeight float32
	UnderlineColor  color.Color
}

// NewUnderlineSegmentedControl creates an underline-style segmented control
func NewUnderlineSegmentedControl(segments []string, onValueChanged func(selectedIndex int)) *UnderlineSegmentedControl {
	config := core.SharedConfiguration()
	sc := NewSegmentedControl(segments, onValueChanged)
	sc.BorderWidth = 0
	sc.CornerRadius = 0
	sc.BackgroundColor = color.Transparent
	sc.SelectedBackgroundColor = color.Transparent
	sc.TextColor = config.GrayColor
	sc.SelectedTextColor = config.BlueColor
	return &UnderlineSegmentedControl{
		SegmentedControl: sc,
		UnderlineHeight:  2,
		UnderlineColor:   config.BlueColor,
	}
}
