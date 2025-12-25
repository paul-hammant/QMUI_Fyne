// Package toast provides QMUIToastView - a toast notification system
// Ported from Tencent's QMUI_iOS framework
package toast

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/core"
)

// ToastPosition defines where the toast appears
type ToastPosition int

const (
	// ToastPositionCenter centers the toast
	ToastPositionCenter ToastPosition = iota
	// ToastPositionTop shows toast at top
	ToastPositionTop
	// ToastPositionBottom shows toast at bottom
	ToastPositionBottom
)

// ToastContentView defines custom content for toast
type ToastContentView interface {
	fyne.CanvasObject
	SetText(text string)
	SetDetailText(detail string)
}

// ToastAnimator handles toast animations
type ToastAnimator interface {
	ShowAnimation(toast *ToastView, completion func())
	HideAnimation(toast *ToastView, completion func())
}

// DefaultToastAnimator provides fade in/out animation
type DefaultToastAnimator struct{}

func (a *DefaultToastAnimator) ShowAnimation(toast *ToastView, completion func()) {
	// In Fyne, we don't have direct opacity animation, so we just complete
	if completion != nil {
		completion()
	}
}

func (a *DefaultToastAnimator) HideAnimation(toast *ToastView, completion func()) {
	if completion != nil {
		completion()
	}
}

// ToastBackgroundView is the background of the toast
type ToastBackgroundView struct {
	widget.BaseWidget

	CornerRadius    float32
	BackgroundColor color.Color
	BlurEnabled     bool
}

// NewToastBackgroundView creates a new toast background
func NewToastBackgroundView() *ToastBackgroundView {
	config := core.SharedConfiguration()
	bg := &ToastBackgroundView{
		CornerRadius:    config.ToastCornerRadius,
		BackgroundColor: config.ToastBackgroundColor,
		BlurEnabled:     false,
	}
	bg.ExtendBaseWidget(bg)
	return bg
}

func (bg *ToastBackgroundView) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(bg.BackgroundColor)
	rect.CornerRadius = bg.CornerRadius
	return &toastBackgroundRenderer{
		bg:   bg,
		rect: rect,
	}
}

type toastBackgroundRenderer struct {
	bg   *ToastBackgroundView
	rect *canvas.Rectangle
}

func (r *toastBackgroundRenderer) Destroy()              {}
func (r *toastBackgroundRenderer) Layout(size fyne.Size) { r.rect.Resize(size) }
func (r *toastBackgroundRenderer) MinSize() fyne.Size    { return fyne.NewSize(0, 0) }
func (r *toastBackgroundRenderer) Refresh() {
	r.rect.FillColor = r.bg.BackgroundColor
	r.rect.CornerRadius = r.bg.CornerRadius
	r.rect.Refresh()
}
func (r *toastBackgroundRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.rect} }

// ToastView is the main toast component
type ToastView struct {
	widget.BaseWidget

	// Content
	Text       string
	DetailText string
	Icon       fyne.Resource

	// Styling
	DisplayPosition   ToastPosition
	ContentInsets     core.EdgeInsets
	CornerRadius      float32
	BackgroundColor   color.Color
	TextColor         color.Color
	DetailTextColor   color.Color
	TextSize          float32
	DetailTextSize    float32
	MarginFromScreen  float32
	IconSize          fyne.Size
	SpacingBetweenIconAndText float32
	SpacingBetweenTextAndDetail float32

	// Behavior
	Duration             float64 // Duration in seconds, 0 means manual dismiss
	ShouldTintIcon       bool
	HidesWhenTapped      bool
	MaskUserInteraction  bool

	// Animation
	Animator ToastAnimator

	// Callbacks
	OnShow func()
	OnHide func()

	// State
	mu       sync.RWMutex
	visible  bool
	timer    *time.Timer
	popup    *widget.PopUp
	window   fyne.Window
}

// NewToastView creates a new toast view
func NewToastView() *ToastView {
	config := core.SharedConfiguration()
	tv := &ToastView{
		DisplayPosition:   ToastPositionCenter,
		ContentInsets:     config.ToastContentInsets,
		CornerRadius:      config.ToastCornerRadius,
		BackgroundColor:   config.ToastBackgroundColor,
		TextColor:         config.ToastTextColor,
		DetailTextColor:   config.ToastTextColor,
		TextSize:          config.ToastFontSize,
		DetailTextSize:    config.ToastFontSize - 2,
		MarginFromScreen:  config.ToastMarginFromScreen,
		IconSize:          fyne.NewSize(32, 32),
		SpacingBetweenIconAndText: 8,
		SpacingBetweenTextAndDetail: 4,
		Duration:          config.ToastDefaultDuration,
		ShouldTintIcon:    true,
		HidesWhenTapped:   true,
		MaskUserInteraction: false,
		Animator:          &DefaultToastAnimator{},
	}
	tv.ExtendBaseWidget(tv)
	return tv
}

// NewToastViewWithText creates a toast with text
func NewToastViewWithText(text string) *ToastView {
	tv := NewToastView()
	tv.Text = text
	return tv
}

// Show displays the toast in the given window
func (tv *ToastView) ShowIn(window fyne.Window) {
	tv.mu.Lock()
	if tv.visible {
		tv.mu.Unlock()
		return
	}
	tv.visible = true
	tv.window = window
	tv.mu.Unlock()

	content := tv.buildContent()
	tv.popup = widget.NewPopUp(content, window.Canvas())

	// Position the popup
	canvasSize := window.Canvas().Size()
	contentSize := content.MinSize()

	var pos fyne.Position
	switch tv.DisplayPosition {
	case ToastPositionTop:
		pos = fyne.NewPos(
			(canvasSize.Width-contentSize.Width)/2,
			tv.MarginFromScreen,
		)
	case ToastPositionBottom:
		pos = fyne.NewPos(
			(canvasSize.Width-contentSize.Width)/2,
			canvasSize.Height-contentSize.Height-tv.MarginFromScreen,
		)
	default: // Center
		pos = fyne.NewPos(
			(canvasSize.Width-contentSize.Width)/2,
			(canvasSize.Height-contentSize.Height)/2,
		)
	}

	tv.popup.Move(pos)
	tv.popup.Show()

	if tv.Animator != nil {
		tv.Animator.ShowAnimation(tv, nil)
	}

	if tv.OnShow != nil {
		tv.OnShow()
	}

	// Set up auto-hide timer
	if tv.Duration > 0 {
		tv.mu.Lock()
		tv.timer = time.AfterFunc(time.Duration(tv.Duration*float64(time.Second)), func() {
			fyne.Do(func() {
				tv.Hide()
			})
		})
		tv.mu.Unlock()
	}
}

// Hide hides the toast
func (tv *ToastView) Hide() {
	tv.mu.Lock()
	if !tv.visible {
		tv.mu.Unlock()
		return
	}
	if tv.timer != nil {
		tv.timer.Stop()
		tv.timer = nil
	}
	tv.mu.Unlock()

	hideFunc := func() {
		tv.mu.Lock()
		if tv.popup != nil {
			tv.popup.Hide()
			tv.popup = nil
		}
		tv.visible = false
		tv.mu.Unlock()

		if tv.OnHide != nil {
			tv.OnHide()
		}
	}

	if tv.Animator != nil {
		tv.Animator.HideAnimation(tv, hideFunc)
	} else {
		hideFunc()
	}
}

// IsVisible returns whether the toast is currently visible
func (tv *ToastView) IsVisible() bool {
	tv.mu.RLock()
	defer tv.mu.RUnlock()
	return tv.visible
}

func (tv *ToastView) buildContent() fyne.CanvasObject {
	var objects []fyne.CanvasObject

	// Icon
	if tv.Icon != nil {
		icon := canvas.NewImageFromResource(tv.Icon)
		icon.FillMode = canvas.ImageFillContain
		icon.SetMinSize(tv.IconSize)
		objects = append(objects, icon)
	}

	// Text
	if tv.Text != "" {
		text := canvas.NewText(tv.Text, tv.TextColor)
		text.TextSize = tv.TextSize
		text.Alignment = fyne.TextAlignCenter
		objects = append(objects, text)
	}

	// Detail text
	if tv.DetailText != "" {
		detail := canvas.NewText(tv.DetailText, tv.DetailTextColor)
		detail.TextSize = tv.DetailTextSize
		detail.Alignment = fyne.TextAlignCenter
		objects = append(objects, detail)
	}

	// Background
	background := canvas.NewRectangle(tv.BackgroundColor)
	background.CornerRadius = tv.CornerRadius

	content := container.NewVBox(objects...)
	padded := container.NewPadded(content)

	return container.NewStack(background, padded)
}

func (tv *ToastView) CreateRenderer() fyne.WidgetRenderer {
	tv.ExtendBaseWidget(tv)
	return &toastRenderer{toast: tv}
}

type toastRenderer struct {
	toast *ToastView
}

func (r *toastRenderer) Destroy()                      {}
func (r *toastRenderer) Layout(size fyne.Size)         {}
func (r *toastRenderer) MinSize() fyne.Size            { return fyne.NewSize(0, 0) }
func (r *toastRenderer) Refresh()                      {}
func (r *toastRenderer) Objects() []fyne.CanvasObject { return nil }

// Tips is a higher-level toast API similar to QMUITips
type Tips struct {
	window fyne.Window
	toast  *ToastView
}

// NewTips creates a new Tips instance for a window
func NewTips(window fyne.Window) *Tips {
	return &Tips{window: window}
}

// ShowText shows a simple text toast
func (t *Tips) ShowText(text string) {
	t.HideCurrent()
	t.toast = NewToastViewWithText(text)
	t.toast.ShowIn(t.window)
}

// ShowTextWithDuration shows a toast for a specific duration
func (t *Tips) ShowTextWithDuration(text string, duration float64) {
	t.HideCurrent()
	t.toast = NewToastViewWithText(text)
	t.toast.Duration = duration
	t.toast.ShowIn(t.window)
}

// ShowSuccess shows a success toast
func (t *Tips) ShowSuccess(text string) {
	t.HideCurrent()
	t.toast = NewToastView()
	t.toast.Text = text
	// You would set a success icon here
	t.toast.ShowIn(t.window)
}

// ShowError shows an error toast
func (t *Tips) ShowError(text string) {
	t.HideCurrent()
	t.toast = NewToastView()
	t.toast.Text = text
	t.toast.BackgroundColor = color.RGBA{R: 180, G: 50, B: 50, A: 230}
	t.toast.ShowIn(t.window)
}

// ShowInfo shows an info toast
func (t *Tips) ShowInfo(text string) {
	t.HideCurrent()
	t.toast = NewToastView()
	t.toast.Text = text
	t.toast.ShowIn(t.window)
}

// ShowLoading shows a loading toast (manual dismiss required)
func (t *Tips) ShowLoading(text string) {
	t.HideCurrent()
	t.toast = NewToastView()
	t.toast.Text = text
	t.toast.Duration = 0 // Manual dismiss
	// You would add a loading indicator here
	t.toast.ShowIn(t.window)
}

// HideLoading hides the loading toast
func (t *Tips) HideLoading() {
	t.HideCurrent()
}

// HideCurrent hides any currently showing toast
func (t *Tips) HideCurrent() {
	if t.toast != nil && t.toast.IsVisible() {
		t.toast.Hide()
		t.toast = nil
	}
}

// Global Tips functions

var (
	globalTipsMap sync.Map
)

func getTipsForWindow(window fyne.Window) *Tips {
	if tips, ok := globalTipsMap.Load(window); ok {
		return tips.(*Tips)
	}
	tips := NewTips(window)
	globalTipsMap.Store(window, tips)
	return tips
}

// ShowText shows a text toast in the given window
func ShowText(window fyne.Window, text string) {
	getTipsForWindow(window).ShowText(text)
}

// ShowTextWithDuration shows a toast for a specific duration
func ShowTextWithDuration(window fyne.Window, text string, duration float64) {
	getTipsForWindow(window).ShowTextWithDuration(text, duration)
}

// ShowSuccess shows a success toast
func ShowSuccess(window fyne.Window, text string) {
	getTipsForWindow(window).ShowSuccess(text)
}

// ShowError shows an error toast
func ShowError(window fyne.Window, text string) {
	getTipsForWindow(window).ShowError(text)
}

// ShowInfo shows an info toast
func ShowInfo(window fyne.Window, text string) {
	getTipsForWindow(window).ShowInfo(text)
}

// ShowLoading shows a loading toast
func ShowLoading(window fyne.Window, text string) {
	getTipsForWindow(window).ShowLoading(text)
}

// HideLoading hides the loading toast
func HideLoading(window fyne.Window) {
	getTipsForWindow(window).HideLoading()
}

// Hide hides any toast in the window
func Hide(window fyne.Window) {
	getTipsForWindow(window).HideCurrent()
}

// ShowMessage is an alias for ShowText
func ShowMessage(window fyne.Window, text string) {
	ShowText(window, text)
}
