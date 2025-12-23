// Package tips provides QMUITips - toast variants with icons for Loading/Success/Error/Info
// Ported from Tencent's QMUI_iOS framework
package tips

import (
	"image/color"
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// TipsStyle defines the style of tip to show
type TipsStyle int

const (
	TipsStyleText TipsStyle = iota
	TipsStyleLoading
	TipsStyleSuccess
	TipsStyleError
	TipsStyleInfo
)

// Tips provides a convenient API for showing toast-like notifications with icons
type Tips struct {
	window    fyne.Window
	popup     *widget.PopUp
	mu        sync.RWMutex
	isVisible bool
	timer     *time.Timer

	// Animation state for loading spinner
	spinnerAngle float64
	spinnerAnim  bool
	stopSpinner  chan struct{}
}

// NewTips creates a new Tips instance for a window
func NewTips(window fyne.Window) *Tips {
	return &Tips{window: window}
}

// showTip displays a tip with the given style and text
func (t *Tips) showTip(style TipsStyle, text string, duration float64) {
	t.HideCurrent()

	config := core.SharedConfiguration()

	// Build content based on style
	var objects []fyne.CanvasObject

	switch style {
	case TipsStyleLoading:
		spinner := t.createLoadingSpinner()
		objects = append(objects, spinner)
	case TipsStyleSuccess:
		icon := t.createSuccessIcon()
		objects = append(objects, icon)
	case TipsStyleError:
		icon := t.createErrorIcon()
		objects = append(objects, icon)
	case TipsStyleInfo:
		icon := t.createInfoIcon()
		objects = append(objects, icon)
	}

	// Add text label
	if text != "" {
		label := canvas.NewText(text, config.ToastTextColor)
		label.TextSize = config.ToastFontSize
		label.Alignment = fyne.TextAlignCenter
		objects = append(objects, label)
	}

	// Background
	background := canvas.NewRectangle(config.ToastBackgroundColor)
	background.CornerRadius = config.ToastCornerRadius

	content := container.NewVBox(objects...)
	padded := container.NewPadded(content)

	popupContent := container.NewStack(background, padded)

	t.mu.Lock()
	t.popup = widget.NewPopUp(popupContent, t.window.Canvas())
	t.isVisible = true
	t.mu.Unlock()

	// Position the popup at center
	canvasSize := t.window.Canvas().Size()
	contentSize := popupContent.MinSize()

	pos := fyne.NewPos(
		(canvasSize.Width-contentSize.Width)/2,
		(canvasSize.Height-contentSize.Height)/2,
	)

	t.popup.Move(pos)
	t.popup.Show()

	// Set up auto-hide timer (except for loading which requires manual dismiss)
	if duration > 0 && style != TipsStyleLoading {
		t.mu.Lock()
		t.timer = time.AfterFunc(time.Duration(duration*float64(time.Second)), func() {
			t.HideCurrent()
		})
		t.mu.Unlock()
	}
}

// createLoadingSpinner creates an animated loading spinner
func (t *Tips) createLoadingSpinner() fyne.CanvasObject {
	spinner := &loadingSpinner{tips: t}
	spinner.ExtendBaseWidget(spinner)

	// Start animation
	t.mu.Lock()
	t.spinnerAnim = true
	t.stopSpinner = make(chan struct{})
	t.mu.Unlock()

	go t.animateSpinner(spinner)

	return spinner
}

func (t *Tips) animateSpinner(spinner *loadingSpinner) {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopSpinner:
			return
		case <-ticker.C:
			t.mu.Lock()
			t.spinnerAngle += 30
			if t.spinnerAngle >= 360 {
				t.spinnerAngle = 0
			}
			t.mu.Unlock()
			spinner.Refresh()
		}
	}
}

// createSuccessIcon creates a checkmark icon
func (t *Tips) createSuccessIcon() fyne.CanvasObject {
	icon := &successIcon{}
	icon.ExtendBaseWidget(icon)
	return icon
}

// createErrorIcon creates an X icon
func (t *Tips) createErrorIcon() fyne.CanvasObject {
	icon := &errorIcon{}
	icon.ExtendBaseWidget(icon)
	return icon
}

// createInfoIcon creates an info icon
func (t *Tips) createInfoIcon() fyne.CanvasObject {
	icon := &infoIcon{}
	icon.ExtendBaseWidget(icon)
	return icon
}

// ShowText shows a simple text tip
func (t *Tips) ShowText(text string) {
	t.showTip(TipsStyleText, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowTextWithDuration shows a text tip for a specific duration
func (t *Tips) ShowTextWithDuration(text string, duration float64) {
	t.showTip(TipsStyleText, text, duration)
}

// ShowLoading shows a loading tip (manual dismiss required)
func (t *Tips) ShowLoading(text string) {
	t.showTip(TipsStyleLoading, text, 0)
}

// ShowLoadingWithDuration shows a loading tip that auto-hides
func (t *Tips) ShowLoadingWithDuration(text string, duration float64) {
	t.showTip(TipsStyleLoading, text, duration)
}

// ShowSuccess shows a success tip with checkmark
func (t *Tips) ShowSuccess(text string) {
	t.showTip(TipsStyleSuccess, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowSuccessWithDuration shows a success tip for a specific duration
func (t *Tips) ShowSuccessWithDuration(text string, duration float64) {
	t.showTip(TipsStyleSuccess, text, duration)
}

// ShowError shows an error tip with X icon
func (t *Tips) ShowError(text string) {
	t.showTip(TipsStyleError, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowErrorWithDuration shows an error tip for a specific duration
func (t *Tips) ShowErrorWithDuration(text string, duration float64) {
	t.showTip(TipsStyleError, text, duration)
}

// ShowInfo shows an info tip with info icon
func (t *Tips) ShowInfo(text string) {
	t.showTip(TipsStyleInfo, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowInfoWithDuration shows an info tip for a specific duration
func (t *Tips) ShowInfoWithDuration(text string, duration float64) {
	t.showTip(TipsStyleInfo, text, duration)
}

// HideLoading hides the loading tip
func (t *Tips) HideLoading() {
	t.HideCurrent()
}

// HideCurrent hides the currently showing tip
func (t *Tips) HideCurrent() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Stop spinner animation
	if t.spinnerAnim && t.stopSpinner != nil {
		close(t.stopSpinner)
		t.spinnerAnim = false
	}

	// Stop timer
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	// Hide popup
	if t.popup != nil && t.isVisible {
		t.popup.Hide()
		t.popup = nil
		t.isVisible = false
	}
}

// IsVisible returns whether a tip is currently showing
func (t *Tips) IsVisible() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isVisible
}

// loadingSpinner widget
type loadingSpinner struct {
	widget.BaseWidget
	tips *Tips
}

func (s *loadingSpinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	return &loadingSpinnerRenderer{spinner: s}
}

type loadingSpinnerRenderer struct {
	spinner *loadingSpinner
	objects []fyne.CanvasObject
}

func (r *loadingSpinnerRenderer) Destroy() {}

func (r *loadingSpinnerRenderer) Layout(size fyne.Size) {
	r.buildObjects(size)
}

func (r *loadingSpinnerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (r *loadingSpinnerRenderer) buildObjects(size fyne.Size) {
	r.objects = nil
	config := core.SharedConfiguration()

	centerX := size.Width / 2
	centerY := size.Height / 2
	radius := float32(15)

	r.spinner.tips.mu.RLock()
	angle := r.spinner.tips.spinnerAngle
	r.spinner.tips.mu.RUnlock()

	// Draw 12 lines in a circle
	numLines := 12
	for i := 0; i < numLines; i++ {
		lineAngle := (float64(i)*30 + angle) * math.Pi / 180

		// Calculate opacity based on position
		opacity := uint8(255 - uint8(i*20))
		if opacity < 50 {
			opacity = 50
		}

		lineColor := color.RGBA{R: 255, G: 255, B: 255, A: opacity}

		x1 := centerX + float32(math.Cos(lineAngle)*float64(radius-6))
		y1 := centerY + float32(math.Sin(lineAngle)*float64(radius-6))
		x2 := centerX + float32(math.Cos(lineAngle)*float64(radius))
		y2 := centerY + float32(math.Sin(lineAngle)*float64(radius))

		line := canvas.NewLine(lineColor)
		line.StrokeWidth = 2.5
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)

		r.objects = append(r.objects, line)
	}

	_ = config // avoid unused warning
}

func (r *loadingSpinnerRenderer) Refresh() {
	r.buildObjects(r.spinner.Size())
}

func (r *loadingSpinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// successIcon widget - checkmark
type successIcon struct {
	widget.BaseWidget
}

func (s *successIcon) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	// Draw a checkmark
	line1 := canvas.NewLine(color.White)
	line1.StrokeWidth = 3

	line2 := canvas.NewLine(color.White)
	line2.StrokeWidth = 3

	return &successIconRenderer{
		icon:  s,
		line1: line1,
		line2: line2,
	}
}

type successIconRenderer struct {
	icon  *successIcon
	line1 *canvas.Line
	line2 *canvas.Line
}

func (r *successIconRenderer) Destroy() {}

func (r *successIconRenderer) Layout(size fyne.Size) {
	// Checkmark shape: from bottom-left to middle-bottom, then to top-right
	r.line1.Position1 = fyne.NewPos(size.Width*0.2, size.Height*0.5)
	r.line1.Position2 = fyne.NewPos(size.Width*0.4, size.Height*0.7)

	r.line2.Position1 = fyne.NewPos(size.Width*0.4, size.Height*0.7)
	r.line2.Position2 = fyne.NewPos(size.Width*0.8, size.Height*0.3)
}

func (r *successIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (r *successIconRenderer) Refresh() {
	r.line1.Refresh()
	r.line2.Refresh()
}

func (r *successIconRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.line1, r.line2}
}

// errorIcon widget - X
type errorIcon struct {
	widget.BaseWidget
}

func (s *errorIcon) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	line1 := canvas.NewLine(color.White)
	line1.StrokeWidth = 3

	line2 := canvas.NewLine(color.White)
	line2.StrokeWidth = 3

	return &errorIconRenderer{
		icon:  s,
		line1: line1,
		line2: line2,
	}
}

type errorIconRenderer struct {
	icon  *errorIcon
	line1 *canvas.Line
	line2 *canvas.Line
}

func (r *errorIconRenderer) Destroy() {}

func (r *errorIconRenderer) Layout(size fyne.Size) {
	// X shape
	padding := float32(8)
	r.line1.Position1 = fyne.NewPos(padding, padding)
	r.line1.Position2 = fyne.NewPos(size.Width-padding, size.Height-padding)

	r.line2.Position1 = fyne.NewPos(size.Width-padding, padding)
	r.line2.Position2 = fyne.NewPos(padding, size.Height-padding)
}

func (r *errorIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (r *errorIconRenderer) Refresh() {
	r.line1.Refresh()
	r.line2.Refresh()
}

func (r *errorIconRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.line1, r.line2}
}

// infoIcon widget - i in circle
type infoIcon struct {
	widget.BaseWidget
}

func (s *infoIcon) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	circle := canvas.NewCircle(color.Transparent)
	circle.StrokeWidth = 2
	circle.StrokeColor = color.White

	dot := canvas.NewCircle(color.White)

	line := canvas.NewLine(color.White)
	line.StrokeWidth = 2

	return &infoIconRenderer{
		icon:   s,
		circle: circle,
		dot:    dot,
		line:   line,
	}
}

type infoIconRenderer struct {
	icon   *infoIcon
	circle *canvas.Circle
	dot    *canvas.Circle
	line   *canvas.Line
}

func (r *infoIconRenderer) Destroy() {}

func (r *infoIconRenderer) Layout(size fyne.Size) {
	// Circle
	r.circle.Resize(size)
	r.circle.Move(fyne.NewPos(0, 0))

	// Dot (top of i)
	dotSize := float32(4)
	r.dot.Resize(fyne.NewSize(dotSize, dotSize))
	r.dot.Move(fyne.NewPos(size.Width/2-dotSize/2, size.Height*0.25))

	// Line (stem of i)
	r.line.Position1 = fyne.NewPos(size.Width/2, size.Height*0.4)
	r.line.Position2 = fyne.NewPos(size.Width/2, size.Height*0.75)
}

func (r *infoIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (r *infoIconRenderer) Refresh() {
	r.circle.Refresh()
	r.dot.Refresh()
	r.line.Refresh()
}

func (r *infoIconRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.circle, r.dot, r.line}
}

// Global convenience functions

var globalTipsMap sync.Map

func getTipsForWindow(window fyne.Window) *Tips {
	if tips, ok := globalTipsMap.Load(window); ok {
		return tips.(*Tips)
	}
	tips := NewTips(window)
	globalTipsMap.Store(window, tips)
	return tips
}

// ShowText shows a text tip
func ShowText(window fyne.Window, text string) {
	getTipsForWindow(window).ShowText(text)
}

// ShowTextWithDuration shows a text tip for a duration
func ShowTextWithDuration(window fyne.Window, text string, duration float64) {
	getTipsForWindow(window).ShowTextWithDuration(text, duration)
}

// ShowLoading shows a loading tip
func ShowLoading(window fyne.Window, text string) {
	getTipsForWindow(window).ShowLoading(text)
}

// ShowSuccess shows a success tip
func ShowSuccess(window fyne.Window, text string) {
	getTipsForWindow(window).ShowSuccess(text)
}

// ShowError shows an error tip
func ShowError(window fyne.Window, text string) {
	getTipsForWindow(window).ShowError(text)
}

// ShowInfo shows an info tip
func ShowInfo(window fyne.Window, text string) {
	getTipsForWindow(window).ShowInfo(text)
}

// HideLoading hides the loading tip
func HideLoading(window fyne.Window) {
	getTipsForWindow(window).HideLoading()
}

// Hide hides any tip
func Hide(window fyne.Window) {
	getTipsForWindow(window).HideCurrent()
}
