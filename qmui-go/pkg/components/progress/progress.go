// Package progress provides QMUIPieProgressView - a pie chart progress view
// Ported from Tencent's QMUI_iOS framework
package progress

import (
	"fmt"
	"image/color"
	"math"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
	"github.com/user/qmui-go/pkg/theme"
)

// PieProgressView displays progress as a pie chart
type PieProgressView struct {
	widget.BaseWidget

	// Progress value (0.0 - 1.0)
	Progress float64

	// Styling
	TintColor       color.Color
	TrackColor      color.Color
	BackgroundColor color.Color
	LineWidth       float32
	ViewSize        fyne.Size

	// Shape
	ShowsTrack bool
	Clockwise  bool

	mu sync.RWMutex
}

// NewPieProgressView creates a new pie progress view
func NewPieProgressView() *PieProgressView {
	config := core.SharedConfiguration()
	ppv := &PieProgressView{
		Progress:        0,
		TintColor:       config.BlueColor,
		TrackColor:      color.RGBA{R: 200, G: 200, B: 200, A: 100},
		BackgroundColor: color.Transparent,
		LineWidth:       3,
		ViewSize:        fyne.NewSize(37, 37),
		ShowsTrack:      true,
		Clockwise:       true,
	}
	ppv.ExtendBaseWidget(ppv)
	return ppv
}

// SetProgress sets the progress value (0.0 - 1.0)
func (ppv *PieProgressView) SetProgress(progress float64) {
	ppv.mu.Lock()
	ppv.Progress = core.ClampFloat64(progress, 0, 1)
	ppv.mu.Unlock()
	fyne.Do(func() {
		ppv.Refresh()
	})
}

// GetProgress returns the current progress value
func (ppv *PieProgressView) GetProgress() float64 {
	ppv.mu.RLock()
	defer ppv.mu.RUnlock()
	return ppv.Progress
}

// CreateRenderer implements fyne.Widget
func (ppv *PieProgressView) CreateRenderer() fyne.WidgetRenderer {
	ppv.ExtendBaseWidget(ppv)
	return &pieProgressRenderer{
		view: ppv,
	}
}

type pieProgressRenderer struct {
	view    *PieProgressView
	objects []fyne.CanvasObject
}

func (r *pieProgressRenderer) Destroy() {}

func (r *pieProgressRenderer) buildObjects(size fyne.Size) {
	r.objects = nil

	// Background
	bg := canvas.NewRectangle(r.view.BackgroundColor)
	bg.Resize(size)
	r.objects = append(r.objects, bg)

	r.view.mu.RLock()
	progress := r.view.Progress
	showTrack := r.view.ShowsTrack
	r.view.mu.RUnlock()

	centerX := size.Width / 2
	centerY := size.Height / 2
	radius := (min(size.Width, size.Height) - r.view.LineWidth) / 2

	// Track circle
	if showTrack {
		track := canvas.NewCircle(color.Transparent)
		track.StrokeColor = r.view.TrackColor
		track.StrokeWidth = r.view.LineWidth
		track.Resize(fyne.NewSize(radius*2, radius*2))
		track.Move(fyne.NewPos(centerX-radius, centerY-radius))
		r.objects = append(r.objects, track)
	}

	// Progress pie wedge (filled from center)
	if progress > 0 {
		wedge := r.createPieWedge(centerX, centerY, radius, progress)
		r.objects = append(r.objects, wedge...)
	}
}

func (r *pieProgressRenderer) createPieWedge(cx, cy, radius float32, progress float64) []fyne.CanvasObject {
	var objects []fyne.CanvasObject

	// Draw filled pie wedge using triangles from center
	segments := int(progress * 72) // More segments for smoother pie
	if segments < 1 {
		segments = 1
	}

	startAngle := -math.Pi / 2 // Start from top (12 o'clock)
	endAngle := startAngle + 2*math.Pi*progress

	for i := 0; i < segments; i++ {
		t1 := float64(i) / float64(segments)
		t2 := float64(i+1) / float64(segments)

		angle1 := startAngle + (endAngle-startAngle)*t1
		angle2 := startAngle + (endAngle-startAngle)*t2

		// Triangle from center to two points on circumference
		x1 := cx + radius*float32(math.Cos(angle1))
		y1 := cy + radius*float32(math.Sin(angle1))
		x2 := cx + radius*float32(math.Cos(angle2))
		y2 := cy + radius*float32(math.Sin(angle2))

		// Draw two lines from center to create filled wedge effect
		line1 := canvas.NewLine(r.view.TintColor)
		line1.StrokeWidth = r.view.LineWidth + 2
		line1.Position1 = fyne.NewPos(cx, cy)
		line1.Position2 = fyne.NewPos(x1, y1)
		objects = append(objects, line1)

		line2 := canvas.NewLine(r.view.TintColor)
		line2.StrokeWidth = r.view.LineWidth + 2
		line2.Position1 = fyne.NewPos(cx, cy)
		line2.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line2)

		// Arc segment
		arcLine := canvas.NewLine(r.view.TintColor)
		arcLine.StrokeWidth = r.view.LineWidth
		arcLine.Position1 = fyne.NewPos(x1, y1)
		arcLine.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, arcLine)
	}

	return objects
}

func (r *pieProgressRenderer) Layout(size fyne.Size) {
	r.buildObjects(size)
}

func (r *pieProgressRenderer) MinSize() fyne.Size {
	return r.view.ViewSize
}

func (r *pieProgressRenderer) Refresh() {
	r.buildObjects(r.view.ViewSize)
}

func (r *pieProgressRenderer) Objects() []fyne.CanvasObject {
	if len(r.objects) == 0 {
		r.buildObjects(r.view.ViewSize)
	}
	return r.objects
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// CircularProgressView displays progress as a circular bar
type CircularProgressView struct {
	widget.BaseWidget

	// Progress value (0.0 - 1.0)
	Progress float64

	// Styling
	TintColor       color.Color
	TrackColor      color.Color
	LineWidth       float32
	ViewSize        fyne.Size
	LineCap         string // "round" or "square"

	// Label
	ShowsText     bool
	LabelFormat   string // e.g., "%.0f%%"
	LabelColor    color.Color
	LabelFontSize float32

	mu sync.RWMutex
}

// NewCircularProgressView creates a new circular progress view
func NewCircularProgressView() *CircularProgressView {
	config := core.SharedConfiguration()
	cpv := &CircularProgressView{
		Progress:      0,
		TintColor:     config.BlueColor,
		TrackColor:    color.RGBA{R: 200, G: 200, B: 200, A: 100},
		LineWidth:     4,
		ViewSize:      fyne.NewSize(50, 50),
		LineCap:       "round",
		ShowsText:     false,
		LabelFormat:   "%.0f%%",
		LabelColor:    color.Black,
		LabelFontSize: 12,
	}
	cpv.ExtendBaseWidget(cpv)
	return cpv
}

// SetProgress sets the progress value
func (cpv *CircularProgressView) SetProgress(progress float64) {
	cpv.mu.Lock()
	cpv.Progress = core.ClampFloat64(progress, 0, 1)
	cpv.mu.Unlock()
	fyne.Do(func() {
		cpv.Refresh()
	})
}

// CreateRenderer implements fyne.Widget
func (cpv *CircularProgressView) CreateRenderer() fyne.WidgetRenderer {
	cpv.ExtendBaseWidget(cpv)
	return &circularProgressRenderer{view: cpv}
}

type circularProgressRenderer struct {
	view    *CircularProgressView
	objects []fyne.CanvasObject
}

func (r *circularProgressRenderer) Destroy() {}

func (r *circularProgressRenderer) buildObjects(size fyne.Size) {
	r.objects = nil

	centerX := size.Width / 2
	centerY := size.Height / 2
	radius := (min(size.Width, size.Height) - r.view.LineWidth) / 2

	// Track circle
	track := canvas.NewCircle(color.Transparent)
	track.StrokeColor = r.view.TrackColor
	track.StrokeWidth = r.view.LineWidth
	track.Resize(fyne.NewSize(radius*2, radius*2))
	track.Move(fyne.NewPos(centerX-radius, centerY-radius))
	r.objects = append(r.objects, track)

	r.view.mu.RLock()
	progress := r.view.Progress
	showLabel := r.view.ShowsText
	r.view.mu.RUnlock()

	// Progress arc
	if progress > 0 {
		arc := r.createArc(centerX, centerY, radius, progress)
		r.objects = append(r.objects, arc...)
	}

	// Label
	if showLabel {
		labelText := r.view.LabelFormat
		if labelText == "" {
			labelText = "%.0f%%"
		}
		label := canvas.NewText(
			fmt.Sprintf(labelText, progress*100),
			r.view.LabelColor,
		)
		label.TextSize = r.view.LabelFontSize
		label.Alignment = fyne.TextAlignCenter
		labelSize := label.MinSize()
		label.Move(fyne.NewPos(centerX-labelSize.Width/2, centerY-labelSize.Height/2))
		r.objects = append(r.objects, label)
	}
}

func (r *circularProgressRenderer) createArc(cx, cy, radius float32, progress float64) []fyne.CanvasObject {
	var objects []fyne.CanvasObject

	segments := int(progress * 36)
	if segments < 1 {
		segments = 1
	}

	startAngle := -math.Pi / 2
	endAngle := startAngle + 2*math.Pi*progress

	for i := 0; i < segments; i++ {
		t1 := float64(i) / float64(segments)
		t2 := float64(i+1) / float64(segments)

		angle1 := startAngle + (endAngle-startAngle)*t1
		angle2 := startAngle + (endAngle-startAngle)*t2

		x1 := cx + radius*float32(math.Cos(angle1))
		y1 := cy + radius*float32(math.Sin(angle1))
		x2 := cx + radius*float32(math.Cos(angle2))
		y2 := cy + radius*float32(math.Sin(angle2))

		line := canvas.NewLine(r.view.TintColor)
		line.StrokeWidth = r.view.LineWidth
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line)
	}

	return objects
}

func (r *circularProgressRenderer) Layout(size fyne.Size) {
	r.buildObjects(size)
}

func (r *circularProgressRenderer) MinSize() fyne.Size {
	return r.view.ViewSize
}

func (r *circularProgressRenderer) Refresh() {
	r.buildObjects(r.view.ViewSize)
}

func (r *circularProgressRenderer) Objects() []fyne.CanvasObject {
	if len(r.objects) == 0 {
		r.buildObjects(r.view.ViewSize)
	}
	return r.objects
}

// LinearProgressView displays progress as a linear bar
type LinearProgressView struct {
	widget.BaseWidget

	// Progress value (0.0 - 1.0)
	Progress float64

	// Styling
	TintColor       color.Color
	TrackColor      color.Color
	Height          float32
	CornerRadius    float32

	mu sync.RWMutex
}

// NewLinearProgressView creates a new linear progress view
func NewLinearProgressView() *LinearProgressView {
	config := core.SharedConfiguration()
	lpv := &LinearProgressView{
		Progress:     0,
		TintColor:    config.BlueColor,
		TrackColor:   color.RGBA{R: 200, G: 200, B: 200, A: 100},
		Height:       4,
		CornerRadius: 2,
	}
	lpv.ExtendBaseWidget(lpv)
	return lpv
}

// SetProgress sets the progress value
func (lpv *LinearProgressView) SetProgress(progress float64) {
	lpv.mu.Lock()
	lpv.Progress = core.ClampFloat64(progress, 0, 1)
	lpv.mu.Unlock()
	fyne.Do(func() {
		lpv.Refresh()
	})
}

// CreateRenderer implements fyne.Widget
func (lpv *LinearProgressView) CreateRenderer() fyne.WidgetRenderer {
	lpv.ExtendBaseWidget(lpv)

	track := canvas.NewRectangle(lpv.TrackColor)
	track.CornerRadius = lpv.CornerRadius

	progress := canvas.NewRectangle(lpv.TintColor)
	progress.CornerRadius = lpv.CornerRadius

	return &linearProgressRenderer{
		view:     lpv,
		track:    track,
		progress: progress,
	}
}

type linearProgressRenderer struct {
	view     *LinearProgressView
	track    *canvas.Rectangle
	progress *canvas.Rectangle
	size     fyne.Size
}

func (r *linearProgressRenderer) Destroy() {}

func (r *linearProgressRenderer) Layout(size fyne.Size) {
	r.size = size
	r.track.Resize(size)

	r.view.mu.RLock()
	progress := r.view.Progress
	r.view.mu.RUnlock()

	progressWidth := size.Width * float32(progress)
	r.progress.Resize(fyne.NewSize(progressWidth, size.Height))
}

func (r *linearProgressRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, r.view.Height)
}

func (r *linearProgressRenderer) Refresh() {
	r.track.FillColor = r.view.TrackColor
	r.track.CornerRadius = r.view.CornerRadius
	r.progress.FillColor = r.view.TintColor
	r.progress.CornerRadius = r.view.CornerRadius

	// Update progress bar width based on current progress
	if r.size.Width > 0 {
		r.view.mu.RLock()
		progress := r.view.Progress
		r.view.mu.RUnlock()
		progressWidth := r.size.Width * float32(progress)
		r.progress.Resize(fyne.NewSize(progressWidth, r.size.Height))
	}

	r.track.Refresh()
	r.progress.Refresh()
}

func (r *linearProgressRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.track, r.progress}
}

// ApplyTheme implements the Themeable interface for PieProgressView
func (ppv *PieProgressView) ApplyTheme(t *theme.Theme) {
	ppv.TintColor = t.PrimaryColor
	ppv.Refresh()
}

// ApplyTheme implements the Themeable interface for CircularProgressView
func (cpv *CircularProgressView) ApplyTheme(t *theme.Theme) {
	cpv.TintColor = t.PrimaryColor
	cpv.Refresh()
}

// ApplyTheme implements the Themeable interface for LinearProgressView
func (lpv *LinearProgressView) ApplyTheme(t *theme.Theme) {
	lpv.TintColor = t.PrimaryColor
	lpv.Refresh()
}
