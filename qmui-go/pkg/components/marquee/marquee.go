// Package marquee provides QMUIMarqueeLabel - a scrolling/marquee text label
// Ported from Tencent's QMUI_iOS framework
package marquee

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// MarqueeDirection defines the scroll direction
type MarqueeDirection int

const (
	// MarqueeDirectionLeft scrolls left
	MarqueeDirectionLeft MarqueeDirection = iota
	// MarqueeDirectionRight scrolls right
	MarqueeDirectionRight
)

// MarqueeLabel is a scrolling text label
type MarqueeLabel struct {
	widget.BaseWidget

	// Content
	Text      string
	TextStyle fyne.TextStyle
	TextSize  float32
	TextColor color.Color

	// Scrolling behavior
	Direction     MarqueeDirection
	Speed         float32 // Pixels per second
	PauseDuration time.Duration
	FadeWidth     float32 // Width of fade effect at edges

	// State
	AutoScrollWhenFits bool // Only scroll if text doesn't fit
	IsAnimating        bool

	mu           sync.RWMutex
	offset       float32
	textWidth    float32
	animating    bool
	stopChan     chan struct{}
}

// NewMarqueeLabel creates a new marquee label
func NewMarqueeLabel(text string) *MarqueeLabel {
	ml := &MarqueeLabel{
		Text:               text,
		TextStyle:          fyne.TextStyle{},
		TextSize:           14,
		TextColor:          color.Black,
		Direction:          MarqueeDirectionLeft,
		Speed:              30,
		PauseDuration:      time.Second * 2,
		FadeWidth:          10,
		AutoScrollWhenFits: true,
	}
	ml.ExtendBaseWidget(ml)
	return ml
}

// SetText sets the label text
func (ml *MarqueeLabel) SetText(text string) {
	ml.mu.Lock()
	ml.Text = text
	ml.offset = 0
	ml.mu.Unlock()
	ml.Refresh()
}

// StartAnimation starts the marquee animation
func (ml *MarqueeLabel) StartAnimation() {
	ml.mu.Lock()
	if ml.animating {
		ml.mu.Unlock()
		return
	}
	ml.animating = true
	ml.IsAnimating = true
	ml.stopChan = make(chan struct{})
	ml.mu.Unlock()

	go ml.animate()
}

// StopAnimation stops the marquee animation
func (ml *MarqueeLabel) StopAnimation() {
	ml.mu.Lock()
	if !ml.animating {
		ml.mu.Unlock()
		return
	}
	ml.animating = false
	ml.IsAnimating = false
	if ml.stopChan != nil {
		close(ml.stopChan)
	}
	ml.offset = 0
	ml.mu.Unlock()
	ml.Refresh()
}

func (ml *MarqueeLabel) animate() {
	ticker := time.NewTicker(time.Millisecond * 16) // ~60fps
	defer ticker.Stop()

	ml.mu.RLock()
	pauseDuration := ml.PauseDuration
	ml.mu.RUnlock()

	// Initial pause
	select {
	case <-time.After(pauseDuration):
	case <-ml.stopChan:
		return
	}

	for {
		select {
		case <-ml.stopChan:
			return
		case <-ticker.C:
			ml.mu.Lock()
			containerWidth := ml.Size().Width
			needsScroll := ml.textWidth > containerWidth

			if !needsScroll && ml.AutoScrollWhenFits {
				ml.mu.Unlock()
				continue
			}

			speed := ml.Speed
			direction := ml.Direction
			textWidth := ml.textWidth

			// Calculate movement
			delta := speed * 0.016 // 16ms tick

			if direction == MarqueeDirectionLeft {
				ml.offset -= delta
				// Reset when fully scrolled
				if ml.offset < -textWidth {
					ml.offset = containerWidth
				}
			} else {
				ml.offset += delta
				// Reset when fully scrolled
				if ml.offset > containerWidth {
					ml.offset = -textWidth
				}
			}

			ml.mu.Unlock()
			ml.Refresh()
		}
	}
}

// CreateRenderer implements fyne.Widget
func (ml *MarqueeLabel) CreateRenderer() fyne.WidgetRenderer {
	ml.ExtendBaseWidget(ml)

	text := canvas.NewText(ml.Text, ml.TextColor)
	text.TextStyle = ml.TextStyle
	text.TextSize = ml.TextSize

	// Clone for seamless scrolling
	textClone := canvas.NewText(ml.Text, ml.TextColor)
	textClone.TextStyle = ml.TextStyle
	textClone.TextSize = ml.TextSize

	return &marqueeLabelRenderer{
		label:     ml,
		text:      text,
		textClone: textClone,
	}
}

type marqueeLabelRenderer struct {
	label     *MarqueeLabel
	text      *canvas.Text
	textClone *canvas.Text
}

func (r *marqueeLabelRenderer) Destroy() {
	r.label.StopAnimation()
}

func (r *marqueeLabelRenderer) Layout(size fyne.Size) {
	r.label.mu.Lock()
	r.label.textWidth = r.text.MinSize().Width
	offset := r.label.offset
	textWidth := r.label.textWidth
	r.label.mu.Unlock()

	textHeight := r.text.MinSize().Height
	y := (size.Height - textHeight) / 2

	r.text.Move(fyne.NewPos(offset, y))
	r.textClone.Move(fyne.NewPos(offset+textWidth+50, y)) // Gap between copies
}

func (r *marqueeLabelRenderer) MinSize() fyne.Size {
	textSize := r.text.MinSize()
	return fyne.NewSize(textSize.Width, textSize.Height)
}

func (r *marqueeLabelRenderer) Refresh() {
	r.label.mu.RLock()
	text := r.label.Text
	offset := r.label.offset
	textWidth := r.label.textWidth
	r.label.mu.RUnlock()

	r.text.Text = text
	r.text.Color = r.label.TextColor
	r.text.TextStyle = r.label.TextStyle
	r.text.TextSize = r.label.TextSize

	r.textClone.Text = text
	r.textClone.Color = r.label.TextColor
	r.textClone.TextStyle = r.label.TextStyle
	r.textClone.TextSize = r.label.TextSize

	// Update text positions based on current offset (for animation)
	textHeight := r.text.MinSize().Height
	size := r.label.Size()
	y := (size.Height - textHeight) / 2
	r.text.Move(fyne.NewPos(offset, y))
	r.textClone.Move(fyne.NewPos(offset+textWidth+50, y))

	r.text.Refresh()
	r.textClone.Refresh()
}

func (r *marqueeLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text, r.textClone}
}
