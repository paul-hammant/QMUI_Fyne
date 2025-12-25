// Package tips provides HUD overlays with loading/success/error/info icons
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

// HUDStyle defines the style of tip to show
type HUDStyle int

const (
	HUDStyleText HUDStyle = iota
	HUDStyleLoading
	HUDStyleSuccess
	HUDStyleError
	HUDStyleInfo
)

// HUD provides a convenient API for showing toast-like notifications with icons
type HUD struct {
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

// NewHUD creates a new HUD instance for a window
func NewHUD(window fyne.Window) *HUD {
	return &HUD{window: window}
}

// showTip displays a tip with the given style and text
func (t *HUD) showTip(style HUDStyle, text string, duration float64) {
	t.HideCurrent()

	config := core.SharedConfiguration()

	// Build content based on style
	var objects []fyne.CanvasObject

	switch style {
	case HUDStyleLoading:
		spinner := t.createLoadingSpinner()
		objects = append(objects, spinner)
	case HUDStyleSuccess:
		icon := t.createSuccessIcon()
		objects = append(objects, icon)
	case HUDStyleError:
		icon := t.createErrorIcon()
		objects = append(objects, icon)
	case HUDStyleInfo:
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
	if duration > 0 && style != HUDStyleLoading {
		t.mu.Lock()
		t.timer = time.AfterFunc(time.Duration(duration*float64(time.Second)), func() {
			t.HideCurrent()
		})
		t.mu.Unlock()
	}
}

// createLoadingSpinner creates an animated loading spinner
func (t *HUD) createLoadingSpinner() fyne.CanvasObject {
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

func (t *HUD) animateSpinner(spinner *loadingSpinner) {
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
func (t *HUD) createSuccessIcon() fyne.CanvasObject {
	icon := &successIcon{}
	icon.ExtendBaseWidget(icon)
	return icon
}

// createErrorIcon creates an X icon
func (t *HUD) createErrorIcon() fyne.CanvasObject {
	icon := &errorIcon{}
	icon.ExtendBaseWidget(icon)
	return icon
}

// createInfoIcon creates an info icon
func (t *HUD) createInfoIcon() fyne.CanvasObject {
	icon := &infoIcon{}
	icon.ExtendBaseWidget(icon)
	return icon
}

// ShowText shows a simple text tip
func (t *HUD) ShowText(text string) {
	t.showTip(HUDStyleText, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowTextWithDuration shows a text tip for a specific duration
func (t *HUD) ShowTextWithDuration(text string, duration float64) {
	t.showTip(HUDStyleText, text, duration)
}

// ShowLoading shows a loading tip (manual dismiss required)
func (t *HUD) ShowLoading(text string) {
	t.showTip(HUDStyleLoading, text, 0)
}

// ShowLoadingWithDuration shows a loading tip that auto-hides
func (t *HUD) ShowLoadingWithDuration(text string, duration float64) {
	t.showTip(HUDStyleLoading, text, duration)
}

// ShowSuccess shows a success tip with checkmark
func (t *HUD) ShowSuccess(text string) {
	t.showTip(HUDStyleSuccess, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowSuccessWithDuration shows a success tip for a specific duration
func (t *HUD) ShowSuccessWithDuration(text string, duration float64) {
	t.showTip(HUDStyleSuccess, text, duration)
}

// ShowError shows an error tip with X icon
func (t *HUD) ShowError(text string) {
	t.showTip(HUDStyleError, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowErrorWithDuration shows an error tip for a specific duration
func (t *HUD) ShowErrorWithDuration(text string, duration float64) {
	t.showTip(HUDStyleError, text, duration)
}

// ShowInfo shows an info tip with info icon
func (t *HUD) ShowInfo(text string) {
	t.showTip(HUDStyleInfo, text, core.SharedConfiguration().ToastDefaultDuration)
}

// ShowInfoWithDuration shows an info tip for a specific duration
func (t *HUD) ShowInfoWithDuration(text string, duration float64) {
	t.showTip(HUDStyleInfo, text, duration)
}

// HideLoading hides the loading tip
func (t *HUD) HideLoading() {
	t.HideCurrent()
}

// HideCurrent hides the currently showing tip
func (t *HUD) HideCurrent() {
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
func (t *HUD) IsVisible() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isVisible
}

// loadingSpinner widget
type loadingSpinner struct {
	widget.BaseWidget
	tips *HUD
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

// successIcon widget - circle with checkmark inside (iOS QMUI style)
type successIcon struct {
	widget.BaseWidget
}

func (s *successIcon) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	// Circle outline
	circle := canvas.NewCircle(color.Transparent)
	circle.StrokeWidth = 2
	circle.StrokeColor = color.White

	// Checkmark inside circle
	line1 := canvas.NewLine(color.White)
	line1.StrokeWidth = 2.5

	line2 := canvas.NewLine(color.White)
	line2.StrokeWidth = 2.5

	return &successIconRenderer{
		icon:   s,
		circle: circle,
		line1:  line1,
		line2:  line2,
	}
}

type successIconRenderer struct {
	icon   *successIcon
	circle *canvas.Circle
	line1  *canvas.Line
	line2  *canvas.Line
}

func (r *successIconRenderer) Destroy() {}

func (r *successIconRenderer) Layout(size fyne.Size) {
	// Circle fills the icon area
	r.circle.Resize(size)
	r.circle.Move(fyne.NewPos(0, 0))

	// Checkmark shape inside circle: from left to bottom-middle, then to top-right
	// Position relative to circle center for proper iOS look
	cx := size.Width / 2
	cy := size.Height / 2
	scale := size.Width * 0.3 // checkmark size relative to icon

	// Start point (left side)
	startX := cx - scale*0.8
	startY := cy

	// Middle point (bottom of tick)
	midX := cx - scale*0.2
	midY := cy + scale*0.5

	// End point (top right)
	endX := cx + scale*0.8
	endY := cy - scale*0.5

	r.line1.Position1 = fyne.NewPos(startX, startY)
	r.line1.Position2 = fyne.NewPos(midX, midY)

	r.line2.Position1 = fyne.NewPos(midX, midY)
	r.line2.Position2 = fyne.NewPos(endX, endY)
}

func (r *successIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (r *successIconRenderer) Refresh() {
	r.circle.Refresh()
	r.line1.Refresh()
	r.line2.Refresh()
}

func (r *successIconRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.circle, r.line1, r.line2}
}

// errorIcon widget - circle with X inside (iOS QMUI style)
type errorIcon struct {
	widget.BaseWidget
}

func (s *errorIcon) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	// Circle outline
	circle := canvas.NewCircle(color.Transparent)
	circle.StrokeWidth = 2
	circle.StrokeColor = color.White

	line1 := canvas.NewLine(color.White)
	line1.StrokeWidth = 2.5

	line2 := canvas.NewLine(color.White)
	line2.StrokeWidth = 2.5

	return &errorIconRenderer{
		icon:   s,
		circle: circle,
		line1:  line1,
		line2:  line2,
	}
}

type errorIconRenderer struct {
	icon   *errorIcon
	circle *canvas.Circle
	line1  *canvas.Line
	line2  *canvas.Line
}

func (r *errorIconRenderer) Destroy() {}

func (r *errorIconRenderer) Layout(size fyne.Size) {
	// Circle fills the icon area
	r.circle.Resize(size)
	r.circle.Move(fyne.NewPos(0, 0))

	// X shape inside circle
	cx := size.Width / 2
	cy := size.Height / 2
	offset := size.Width * 0.22 // X size relative to circle

	r.line1.Position1 = fyne.NewPos(cx-offset, cy-offset)
	r.line1.Position2 = fyne.NewPos(cx+offset, cy+offset)

	r.line2.Position1 = fyne.NewPos(cx+offset, cy-offset)
	r.line2.Position2 = fyne.NewPos(cx-offset, cy+offset)
}

func (r *errorIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (r *errorIconRenderer) Refresh() {
	r.circle.Refresh()
	r.line1.Refresh()
	r.line2.Refresh()
}

func (r *errorIconRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.circle, r.line1, r.line2}
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

var globalHUDMap sync.Map

func getHUDForWindow(window fyne.Window) *HUD {
	if tips, ok := globalHUDMap.Load(window); ok {
		return tips.(*HUD)
	}
	tips := NewHUD(window)
	globalHUDMap.Store(window, tips)
	return tips
}

// ShowText shows a text tip
func ShowText(window fyne.Window, text string) {
	getHUDForWindow(window).ShowText(text)
}

// ShowTextWithDuration shows a text tip for a duration
func ShowTextWithDuration(window fyne.Window, text string, duration float64) {
	getHUDForWindow(window).ShowTextWithDuration(text, duration)
}

// ShowLoading shows a loading tip
func ShowLoading(window fyne.Window, text string) {
	getHUDForWindow(window).ShowLoading(text)
}

// ShowSuccess shows a success tip
func ShowSuccess(window fyne.Window, text string) {
	getHUDForWindow(window).ShowSuccess(text)
}

// ShowError shows an error tip
func ShowError(window fyne.Window, text string) {
	getHUDForWindow(window).ShowError(text)
}

// ShowInfo shows an info tip
func ShowInfo(window fyne.Window, text string) {
	getHUDForWindow(window).ShowInfo(text)
}

// HideLoading hides the loading tip
func HideLoading(window fyne.Window) {
	getHUDForWindow(window).HideLoading()
}

// Hide hides any tip
func Hide(window fyne.Window) {
	getHUDForWindow(window).HideCurrent()
}
