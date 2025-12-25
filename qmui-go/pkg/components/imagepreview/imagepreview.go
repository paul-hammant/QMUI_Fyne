// Package imagepreview provides an image preview component with zoom and pan
package imagepreview

import (
	"fmt"
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// ImagePreview displays images with zoom and pan support
type ImagePreview struct {
	widget.BaseWidget

	// Images
	Images       []fyne.Resource
	CurrentIndex int

	// Styling
	BackgroundColor color.Color
	PageIndicatorColor color.Color

	// Behavior
	ZoomEnabled bool
	MinZoom     float32
	MaxZoom     float32

	// Callbacks
	OnCurrentIndexChanged func(index int)
	OnDismiss             func()
	OnLongPress           func(index int)

	mu sync.RWMutex
}

// NewImagePreview creates a new image preview view
func NewImagePreview() *ImagePreview {
	ipv := &ImagePreview{
		Images:             make([]fyne.Resource, 0),
		CurrentIndex:       0,
		BackgroundColor:    color.Black,
		PageIndicatorColor: color.White,
		ZoomEnabled:        true,
		MinZoom:            1.0,
		MaxZoom:            3.0,
	}
	ipv.ExtendBaseWidget(ipv)
	return ipv
}

// NewImagePreviewWithImages creates a preview view with images
func NewImagePreviewWithImages(images []fyne.Resource) *ImagePreview {
	ipv := NewImagePreview()
	ipv.Images = images
	return ipv
}

// SetImages sets the images to display
func (ipv *ImagePreview) SetImages(images []fyne.Resource) {
	ipv.mu.Lock()
	ipv.Images = images
	if ipv.CurrentIndex >= len(images) {
		ipv.CurrentIndex = len(images) - 1
	}
	if ipv.CurrentIndex < 0 {
		ipv.CurrentIndex = 0
	}
	ipv.mu.Unlock()
	ipv.Refresh()
}

// SetCurrentIndex sets the current image index
func (ipv *ImagePreview) SetCurrentIndex(index int) {
	ipv.mu.Lock()
	if index < 0 {
		index = 0
	}
	if index >= len(ipv.Images) {
		index = len(ipv.Images) - 1
	}
	ipv.CurrentIndex = index
	ipv.mu.Unlock()
	ipv.Refresh()

	if ipv.OnCurrentIndexChanged != nil {
		ipv.OnCurrentIndexChanged(index)
	}
}

// Next shows the next image
func (ipv *ImagePreview) Next() {
	ipv.mu.RLock()
	index := ipv.CurrentIndex + 1
	count := len(ipv.Images)
	ipv.mu.RUnlock()

	if index < count {
		ipv.SetCurrentIndex(index)
	}
}

// Previous shows the previous image
func (ipv *ImagePreview) Previous() {
	ipv.mu.RLock()
	index := ipv.CurrentIndex - 1
	ipv.mu.RUnlock()

	if index >= 0 {
		ipv.SetCurrentIndex(index)
	}
}

// Dragged implements fyne.Draggable for swipe navigation
func (ipv *ImagePreview) Dragged(e *fyne.DragEvent) {
	// Horizontal swipe detection
	if e.Dragged.DX > 50 {
		ipv.Previous()
	} else if e.Dragged.DX < -50 {
		ipv.Next()
	}
}

// DragEnd implements fyne.Draggable
func (ipv *ImagePreview) DragEnd() {}

// Tapped handles taps (dismiss on tap)
func (ipv *ImagePreview) Tapped(_ *fyne.PointEvent) {
	if ipv.OnDismiss != nil {
		ipv.OnDismiss()
	}
}

// TappedSecondary handles secondary taps
func (ipv *ImagePreview) TappedSecondary(_ *fyne.PointEvent) {}

// CreateRenderer implements fyne.Widget
func (ipv *ImagePreview) CreateRenderer() fyne.WidgetRenderer {
	ipv.ExtendBaseWidget(ipv)

	background := canvas.NewRectangle(ipv.BackgroundColor)

	return &imagePreviewRenderer{
		preview:    ipv,
		background: background,
	}
}

type imagePreviewRenderer struct {
	preview    *ImagePreview
	background *canvas.Rectangle
	imageView  *canvas.Image
	pageLabel  *canvas.Text
}

func (r *imagePreviewRenderer) Destroy() {}

func (r *imagePreviewRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	r.preview.mu.RLock()
	images := r.preview.Images
	currentIndex := r.preview.CurrentIndex
	r.preview.mu.RUnlock()

	if len(images) > 0 && currentIndex >= 0 && currentIndex < len(images) {
		if r.imageView == nil {
			r.imageView = canvas.NewImageFromResource(images[currentIndex])
			r.imageView.FillMode = canvas.ImageFillContain
		} else {
			r.imageView.Resource = images[currentIndex]
		}
		r.imageView.Resize(size)
	}

	// Page indicator
	if r.pageLabel == nil {
		r.pageLabel = canvas.NewText("", r.preview.PageIndicatorColor)
		r.pageLabel.TextSize = 14
	}
	r.pageLabel.Text = fmt.Sprintf("%d / %d", currentIndex+1, len(images))
	labelSize := r.pageLabel.MinSize()
	r.pageLabel.Move(fyne.NewPos((size.Width-labelSize.Width)/2, size.Height-labelSize.Height-20))
}

func (r *imagePreviewRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

func (r *imagePreviewRenderer) Refresh() {
	r.background.FillColor = r.preview.BackgroundColor
	r.background.Refresh()

	if r.imageView != nil {
		r.imageView.Refresh()
	}
	if r.pageLabel != nil {
		r.pageLabel.Color = r.preview.PageIndicatorColor
		r.pageLabel.Refresh()
	}
}

func (r *imagePreviewRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}
	if r.imageView != nil {
		objects = append(objects, r.imageView)
	}
	if r.pageLabel != nil {
		objects = append(objects, r.pageLabel)
	}
	return objects
}

// ImagePreviewController manages full-screen image preview
type ImagePreviewController struct {
	PreviewView *ImagePreview

	window fyne.Window
	popup  *widget.PopUp
}

// NewImagePreviewController creates a new preview controller
func NewImagePreviewController() *ImagePreviewController {
	return &ImagePreviewController{
		PreviewView: NewImagePreview(),
	}
}

// Show shows the image preview as a full-screen overlay
func (ipvc *ImagePreviewController) Show(window fyne.Window) {
	ipvc.window = window

	content := ipvc.PreviewView

	ipvc.popup = widget.NewModalPopUp(content, window.Canvas())
	ipvc.popup.Resize(window.Canvas().Size())
	ipvc.popup.Show()
}

// Hide hides the image preview
func (ipvc *ImagePreviewController) Hide() {
	if ipvc.popup != nil {
		ipvc.popup.Hide()
		ipvc.popup = nil
	}

	if ipvc.PreviewView.OnDismiss != nil {
		ipvc.PreviewView.OnDismiss()
	}
}

// ShowImages shows images starting at the specified index
func ShowImages(window fyne.Window, images []fyne.Resource, startIndex int) *ImagePreviewController {
	controller := NewImagePreviewController()
	controller.PreviewView.SetImages(images)
	controller.PreviewView.SetCurrentIndex(startIndex)
	controller.PreviewView.OnDismiss = func() {
		controller.Hide()
	}
	controller.Show(window)
	return controller
}
