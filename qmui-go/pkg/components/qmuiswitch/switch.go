// Package switch provides QOUISwitch - a custom switch control
// Ported from Tencent's QMUI_iOS framework
package qmuiswitch

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// Switch is a custom switch control that allows for different on/off tint colors
type Switch struct {
	widget.BaseWidget

	// State
	Checked bool
	Enabled bool

	// Styling
	OnTintColor  color.Color
	OffTintColor color.Color
	ThumbColor   color.Color

	// Callbacks
	OnChanged func(bool)

	mu      sync.RWMutex
	hovered bool
}

// NewSwitch creates a new custom switch
func NewSwitch(onChanged func(bool)) *Switch {
	s := &Switch{
		Checked:      false,
		Enabled:      true,
		OnTintColor:  core.SharedConfiguration().BlueColor,
		OffTintColor: color.RGBA{R: 224, G: 224, B: 224, A: 255},
		ThumbColor:   color.White,
		OnChanged:    onChanged,
	}
	s.ExtendBaseWidget(s)
	return s
}

// SetChecked sets the checked state of the switch
func (s *Switch) SetChecked(checked bool) {
	if s.Checked == checked {
		return
	}
	s.Checked = checked
	s.Refresh()
	if s.OnChanged != nil {
		s.OnChanged(s.Checked)
	}
}

// Toggle toggles the checked state of the switch
func (s *Switch) Toggle() {
	s.SetChecked(!s.Checked)
}

// Tapped is called when a regular tap event is received
func (s *Switch) Tapped(*fyne.PointEvent) {
	if !s.Enabled {
		return
	}
	s.Toggle()
}

// TappedSecondary is called when a secondary tap event is received
func (s *Switch) TappedSecondary(*fyne.PointEvent) {}

// MouseIn is called when a desktop pointer enters the widget
func (s *Switch) MouseIn(*desktop.MouseEvent) {
	s.hovered = true
	s.Refresh()
}

// MouseOut is called when a desktop pointer leaves the widget
func (s *Switch) MouseOut() {
	s.hovered = false
	s.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (s *Switch) MouseMoved(*desktop.MouseEvent) {}

// Cursor returns the cursor type of this widget
func (s *Switch) Cursor() desktop.Cursor {
	if s.Enabled {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *Switch) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	track := canvas.NewRectangle(s.OffTintColor)
	thumb := canvas.NewCircle(s.ThumbColor)

	return &switchRenderer{
		sw:      s,
		track:   track,
		thumb:   thumb,
		objects: []fyne.CanvasObject{track, thumb},
	}
}

type switchRenderer struct {
	sw      *Switch
	track   *canvas.Rectangle
	thumb   *canvas.Circle
	objects []fyne.CanvasObject
}

func (r *switchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(51, 31)
}

func (r *switchRenderer) Layout(size fyne.Size) {
	r.track.Resize(size)
	r.track.CornerRadius = size.Height / 2

	thumbSize := fyne.NewSize(size.Height-2, size.Height-2)
	r.thumb.Resize(thumbSize)

	var thumbPos fyne.Position
	if r.sw.Checked {
		thumbPos = fyne.NewPos(size.Width-thumbSize.Width-1, 1)
	} else {
		thumbPos = fyne.NewPos(1, 1)
	}
	r.thumb.Move(thumbPos)
}

func (r *switchRenderer) Refresh() {
	r.sw.mu.RLock()
	defer r.sw.mu.RUnlock()

	if r.sw.Checked {
		r.track.FillColor = r.sw.OnTintColor
	} else {
		r.track.FillColor = r.sw.OffTintColor
	}
	r.thumb.FillColor = r.sw.ThumbColor

	if !r.sw.Enabled {
		config := core.SharedConfiguration()
		r.track.FillColor = core.ColorWithAlpha(r.track.FillColor, config.ControlDisabledAlpha)
		r.thumb.FillColor = core.ColorWithAlpha(r.thumb.FillColor, config.ControlDisabledAlpha)
	} else if r.sw.hovered {
		r.track.FillColor = core.ColorWithAlpha(r.track.FillColor, 0.8)
	}

	r.Layout(r.sw.Size())
	canvas.Refresh(r.sw)
}

func (r *switchRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *switchRenderer) Destroy() {}
