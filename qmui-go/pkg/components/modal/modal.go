// Package modal provides animated modal presentations
package modal

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/animation"
	"github.com/user/qmui-go/pkg/core"
)

// ModalAnimationStyle defines the animation style
type ModalAnimationStyle int

const (
	ModalAnimationStyleFade ModalAnimationStyle = iota
	ModalAnimationStyleSlideUp
	ModalAnimationStyleSlideDown
	ModalAnimationStyleSlideLeft
	ModalAnimationStyleSlideRight
	ModalAnimationStyleZoom
	ModalAnimationStyleBounce
)

// ModalContentPosition defines where content appears
type ModalContentPosition int

const (
	ModalContentPositionCenter ModalContentPosition = iota
	ModalContentPositionTop
	ModalContentPositionBottom
	ModalContentPositionLeft
	ModalContentPositionRight
)

// Modal manages animated modal presentations
type Modal struct {
	widget.BaseWidget

	// Content
	ContentView fyne.CanvasObject
	ContentSize fyne.Size

	// Animation
	AnimationStyle    ModalAnimationStyle
	AnimationDuration time.Duration
	AnimationEasing   animation.EasingFunction

	// Styling
	DimmingColor        color.Color
	CornerRadius        float32
	BackgroundColor     color.Color
	ContentPosition     ModalContentPosition
	ContentMargin       core.EdgeInsets
	ShadowEnabled       bool
	ShadowColor         color.Color
	ShadowOffset        fyne.Position
	ShadowRadius        float32

	// Behavior
	DismissOnTapOutside     bool
	KeyboardFollowsContent  bool

	// Callbacks
	OnWillPresent func()
	OnDidPresent  func()
	OnWillDismiss func()
	OnDidDismiss  func()
	OnTapOutside  func() bool // Return true to allow dismiss

	// State
	mu          sync.RWMutex
	window      fyne.Window
	popup       *widget.PopUp
	visible     bool
	animating   bool
	dimmer      *canvas.Rectangle
	contentWrapper fyne.CanvasObject
}

// NewModal creates a new modal presentation controller
func NewModal() *Modal {
	config := core.SharedConfiguration()
	mpvc := &Modal{
		AnimationStyle:    ModalAnimationStyleFade,
		AnimationDuration: time.Millisecond * 300,
		AnimationEasing:   animation.EaseOutCubic,
		DimmingColor:      config.MaskDarkColor,
		CornerRadius:      12,
		BackgroundColor:   config.BackgroundColor,
		ContentPosition:   ModalContentPositionCenter,
		ContentMargin:     core.NewEdgeInsets(20, 20, 20, 20),
		ShadowEnabled:     true,
		ShadowColor:       color.RGBA{R: 0, G: 0, B: 0, A: 50},
		ShadowOffset:      fyne.NewPos(0, 4),
		ShadowRadius:      8,
		DismissOnTapOutside: true,
	}
	mpvc.ExtendBaseWidget(mpvc)
	return mpvc
}

// NewModalWithContent creates a modal with content
func NewModalWithContent(content fyne.CanvasObject) *Modal {
	mpvc := NewModal()
	mpvc.ContentView = content
	return mpvc
}

// SetContentView sets the content view
func (mpvc *Modal) SetContentView(content fyne.CanvasObject) {
	mpvc.mu.Lock()
	mpvc.ContentView = content
	mpvc.mu.Unlock()
}

// Present shows the modal with animation
func (mpvc *Modal) Present(window fyne.Window) {
	mpvc.mu.Lock()
	if mpvc.visible || mpvc.animating {
		mpvc.mu.Unlock()
		return
	}
	mpvc.animating = true
	mpvc.window = window
	mpvc.mu.Unlock()

	if mpvc.OnWillPresent != nil {
		mpvc.OnWillPresent()
	}

	// Build content
	content := mpvc.buildContent()
	mpvc.popup = widget.NewModalPopUp(content, window.Canvas())
	mpvc.popup.Resize(window.Canvas().Size())
	mpvc.popup.Show()

	// Animate in
	mpvc.animatePresent(func() {
		mpvc.mu.Lock()
		mpvc.visible = true
		mpvc.animating = false
		mpvc.mu.Unlock()

		if mpvc.OnDidPresent != nil {
			mpvc.OnDidPresent()
		}
	})
}

// Dismiss hides the modal with animation
func (mpvc *Modal) Dismiss() {
	mpvc.mu.Lock()
	if !mpvc.visible || mpvc.animating {
		mpvc.mu.Unlock()
		return
	}
	mpvc.animating = true
	mpvc.mu.Unlock()

	if mpvc.OnWillDismiss != nil {
		mpvc.OnWillDismiss()
	}

	// Animate out
	mpvc.animateDismiss(func() {
		fyne.Do(func() {
			mpvc.mu.Lock()
			if mpvc.popup != nil {
				mpvc.popup.Hide()
				mpvc.popup = nil
			}
			mpvc.visible = false
			mpvc.animating = false
			mpvc.mu.Unlock()

			if mpvc.OnDidDismiss != nil {
				mpvc.OnDidDismiss()
			}
		})
	})
}

// IsVisible returns whether the modal is visible
func (mpvc *Modal) IsVisible() bool {
	mpvc.mu.RLock()
	defer mpvc.mu.RUnlock()
	return mpvc.visible
}

func (mpvc *Modal) buildContent() fyne.CanvasObject {
	// Dimmer background
	mpvc.dimmer = canvas.NewRectangle(mpvc.DimmingColor)

	// Content background
	contentBg := canvas.NewRectangle(mpvc.BackgroundColor)
	contentBg.CornerRadius = mpvc.CornerRadius

	// Shadow (if enabled)
	var shadow *canvas.Rectangle
	if mpvc.ShadowEnabled {
		shadow = canvas.NewRectangle(mpvc.ShadowColor)
		shadow.CornerRadius = mpvc.CornerRadius
	}

	// Wrap content
	var wrappedContent fyne.CanvasObject
	if mpvc.ContentView != nil {
		wrappedContent = container.NewStack(contentBg, container.NewPadded(mpvc.ContentView))
	} else {
		wrappedContent = contentBg
	}

	// Position content
	positioned := mpvc.positionContent(wrappedContent)
	mpvc.contentWrapper = positioned

	// Stack everything
	if shadow != nil {
		return container.NewStack(mpvc.dimmer, positioned)
	}
	return container.NewStack(mpvc.dimmer, positioned)
}

func (mpvc *Modal) positionContent(content fyne.CanvasObject) fyne.CanvasObject {
	switch mpvc.ContentPosition {
	case ModalContentPositionTop:
		return container.NewBorder(content, nil, nil, nil)
	case ModalContentPositionBottom:
		return container.NewBorder(nil, content, nil, nil)
	case ModalContentPositionLeft:
		return container.NewBorder(nil, nil, content, nil)
	case ModalContentPositionRight:
		return container.NewBorder(nil, nil, nil, content)
	default:
		return container.NewCenter(content)
	}
}

func (mpvc *Modal) animatePresent(onComplete func()) {
	switch mpvc.AnimationStyle {
	case ModalAnimationStyleFade:
		mpvc.animateFadeIn(onComplete)
	case ModalAnimationStyleSlideUp:
		mpvc.animateSlideIn(0, 1, onComplete)
	case ModalAnimationStyleSlideDown:
		mpvc.animateSlideIn(0, -1, onComplete)
	case ModalAnimationStyleSlideLeft:
		mpvc.animateSlideIn(1, 0, onComplete)
	case ModalAnimationStyleSlideRight:
		mpvc.animateSlideIn(-1, 0, onComplete)
	case ModalAnimationStyleZoom:
		mpvc.animateZoomIn(onComplete)
	case ModalAnimationStyleBounce:
		mpvc.animateBounceIn(onComplete)
	default:
		if onComplete != nil {
			onComplete()
		}
	}
}

func (mpvc *Modal) animateDismiss(onComplete func()) {
	switch mpvc.AnimationStyle {
	case ModalAnimationStyleFade:
		mpvc.animateFadeOut(onComplete)
	case ModalAnimationStyleSlideUp:
		mpvc.animateSlideOut(0, -1, onComplete)
	case ModalAnimationStyleSlideDown:
		mpvc.animateSlideOut(0, 1, onComplete)
	case ModalAnimationStyleSlideLeft:
		mpvc.animateSlideOut(-1, 0, onComplete)
	case ModalAnimationStyleSlideRight:
		mpvc.animateSlideOut(1, 0, onComplete)
	case ModalAnimationStyleZoom:
		mpvc.animateZoomOut(onComplete)
	default:
		if onComplete != nil {
			onComplete()
		}
	}
}

func (mpvc *Modal) animateFadeIn(onComplete func()) {
	// Fyne doesn't directly support opacity animation, so we just complete
	if onComplete != nil {
		go func() {
			time.Sleep(mpvc.AnimationDuration)
			onComplete()
		}()
	}
}

func (mpvc *Modal) animateFadeOut(onComplete func()) {
	if onComplete != nil {
		go func() {
			time.Sleep(mpvc.AnimationDuration)
			onComplete()
		}()
	}
}

func (mpvc *Modal) animateSlideIn(dirX, dirY float64, onComplete func()) {
	if mpvc.contentWrapper == nil || mpvc.window == nil {
		if onComplete != nil {
			onComplete()
		}
		return
	}

	canvasSize := mpvc.window.Canvas().Size()
	contentSize := mpvc.contentWrapper.MinSize()

	// Calculate start position
	startX := float64(0)
	startY := float64(0)
	if dirX > 0 {
		startX = float64(canvasSize.Width)
	} else if dirX < 0 {
		startX = float64(-contentSize.Width)
	}
	if dirY > 0 {
		startY = float64(canvasSize.Height)
	} else if dirY < 0 {
		startY = float64(-contentSize.Height)
	}

	// Calculate end position (center)
	endX := float64((canvasSize.Width - contentSize.Width) / 2)
	endY := float64((canvasSize.Height - contentSize.Height) / 2)

	animation.NewPositionAnimation(
		startX, startY, endX, endY,
		mpvc.AnimationDuration,
		mpvc.AnimationEasing,
		func(x, y float64) {
			if mpvc.contentWrapper != nil {
				mpvc.contentWrapper.Move(fyne.NewPos(float32(x), float32(y)))
			}
		},
	).Start()

	if onComplete != nil {
		go func() {
			time.Sleep(mpvc.AnimationDuration)
			onComplete()
		}()
	}
}

func (mpvc *Modal) animateSlideOut(dirX, dirY float64, onComplete func()) {
	if mpvc.contentWrapper == nil || mpvc.window == nil {
		if onComplete != nil {
			onComplete()
		}
		return
	}

	canvasSize := mpvc.window.Canvas().Size()
	contentSize := mpvc.contentWrapper.MinSize()
	currentPos := mpvc.contentWrapper.Position()

	// Calculate end position
	endX := float64(currentPos.X)
	endY := float64(currentPos.Y)
	if dirX > 0 {
		endX = float64(canvasSize.Width)
	} else if dirX < 0 {
		endX = float64(-contentSize.Width)
	}
	if dirY > 0 {
		endY = float64(canvasSize.Height)
	} else if dirY < 0 {
		endY = float64(-contentSize.Height)
	}

	animation.NewPositionAnimation(
		float64(currentPos.X), float64(currentPos.Y), endX, endY,
		mpvc.AnimationDuration,
		mpvc.AnimationEasing,
		func(x, y float64) {
			if mpvc.contentWrapper != nil {
				mpvc.contentWrapper.Move(fyne.NewPos(float32(x), float32(y)))
			}
		},
	).Start()

	if onComplete != nil {
		go func() {
			time.Sleep(mpvc.AnimationDuration)
			onComplete()
		}()
	}
}

func (mpvc *Modal) animateZoomIn(onComplete func()) {
	// Zoom animation is tricky without scale support in Fyne
	// For now, just fade in
	mpvc.animateFadeIn(onComplete)
}

func (mpvc *Modal) animateZoomOut(onComplete func()) {
	mpvc.animateFadeOut(onComplete)
}

func (mpvc *Modal) animateBounceIn(onComplete func()) {
	// Use spring easing for bounce effect
	if mpvc.contentWrapper == nil || mpvc.window == nil {
		if onComplete != nil {
			onComplete()
		}
		return
	}

	canvasSize := mpvc.window.Canvas().Size()
	contentSize := mpvc.contentWrapper.MinSize()

	// Start from bottom
	startY := float64(canvasSize.Height)
	endY := float64((canvasSize.Height - contentSize.Height) / 2)
	x := float64((canvasSize.Width - contentSize.Width) / 2)

	springEasing := animation.Spring(8, 12)

	animation.NewPositionAnimation(
		x, startY, x, endY,
		mpvc.AnimationDuration*2,
		springEasing,
		func(px, py float64) {
			if mpvc.contentWrapper != nil {
				mpvc.contentWrapper.Move(fyne.NewPos(float32(px), float32(py)))
			}
		},
	).Start()

	if onComplete != nil {
		go func() {
			time.Sleep(mpvc.AnimationDuration * 2)
			onComplete()
		}()
	}
}

// CreateRenderer implements fyne.Widget
func (mpvc *Modal) CreateRenderer() fyne.WidgetRenderer {
	mpvc.ExtendBaseWidget(mpvc)
	return &modalRenderer{modal: mpvc}
}

type modalRenderer struct {
	modal *Modal
}

func (r *modalRenderer) Destroy()              {}
func (r *modalRenderer) Layout(size fyne.Size) {}
func (r *modalRenderer) MinSize() fyne.Size    { return fyne.NewSize(0, 0) }
func (r *modalRenderer) Refresh()              {}
func (r *modalRenderer) Objects() []fyne.CanvasObject { return nil }

// Helper functions

// PresentModal shows content as a modal with animation
func PresentModal(window fyne.Window, content fyne.CanvasObject, style ModalAnimationStyle) *Modal {
	mpvc := NewModalWithContent(content)
	mpvc.AnimationStyle = style
	mpvc.Present(window)
	return mpvc
}

// PresentModalFromBottom shows content sliding up from bottom
func PresentModalFromBottom(window fyne.Window, content fyne.CanvasObject) *Modal {
	return PresentModal(window, content, ModalAnimationStyleSlideUp)
}

// PresentCenteredModal shows content centered with fade animation
func PresentCenteredModal(window fyne.Window, content fyne.CanvasObject) *Modal {
	return PresentModal(window, content, ModalAnimationStyleFade)
}

// PresentBounceModal shows content with bounce animation
func PresentBounceModal(window fyne.Window, content fyne.CanvasObject) *Modal {
	return PresentModal(window, content, ModalAnimationStyleBounce)
}
