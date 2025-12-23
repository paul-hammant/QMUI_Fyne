// Package zoomimage provides QMUIZoomImageView - a zoomable image view
// Ported from Tencent's QMUI_iOS framework
package zoomimage

import (
	"image/color"
	"math"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ZoomImageView displays an image with zoom and pan support
type ZoomImageView struct {
	widget.BaseWidget

	// Image
	Image    fyne.Resource
	FillMode canvas.ImageFill

	// Zoom settings
	MinZoomScale     float32
	MaxZoomScale     float32
	CurrentZoomScale float32

	// Pan settings
	PanEnabled bool
	OffsetX    float32
	OffsetY    float32

	// Double tap to zoom
	DoubleTapEnabled  bool
	DoubleTapZoomScale float32

	// Callbacks
	OnZoomChanged func(scale float32)
	OnPanChanged  func(x, y float32)

	// State
	mu           sync.RWMutex
	isDragging   bool
	lastMousePos fyne.Position
	imageSize    fyne.Size
}

// NewZoomImageView creates a new zoom image view
func NewZoomImageView() *ZoomImageView {
	ziv := &ZoomImageView{
		MinZoomScale:       0.5,
		MaxZoomScale:       3.0,
		CurrentZoomScale:   1.0,
		PanEnabled:         true,
		DoubleTapEnabled:   true,
		DoubleTapZoomScale: 2.0,
		FillMode:           canvas.ImageFillContain,
	}
	ziv.ExtendBaseWidget(ziv)
	return ziv
}

// NewZoomImageViewWithResource creates a zoom image view with an image
func NewZoomImageViewWithResource(img fyne.Resource) *ZoomImageView {
	ziv := NewZoomImageView()
	ziv.Image = img
	return ziv
}

// SetImage sets the image to display
func (ziv *ZoomImageView) SetImage(img fyne.Resource) {
	ziv.mu.Lock()
	ziv.Image = img
	ziv.CurrentZoomScale = 1.0
	ziv.OffsetX = 0
	ziv.OffsetY = 0
	ziv.mu.Unlock()
	ziv.Refresh()
}

// SetZoomScale sets the current zoom scale
func (ziv *ZoomImageView) SetZoomScale(scale float32) {
	ziv.mu.Lock()
	if scale < ziv.MinZoomScale {
		scale = ziv.MinZoomScale
	}
	if scale > ziv.MaxZoomScale {
		scale = ziv.MaxZoomScale
	}
	ziv.CurrentZoomScale = scale
	ziv.mu.Unlock()
	ziv.Refresh()

	if ziv.OnZoomChanged != nil {
		ziv.OnZoomChanged(scale)
	}
}

// ZoomIn increases the zoom level
func (ziv *ZoomImageView) ZoomIn() {
	ziv.mu.RLock()
	scale := ziv.CurrentZoomScale * 1.25
	ziv.mu.RUnlock()
	ziv.SetZoomScale(scale)
}

// ZoomOut decreases the zoom level
func (ziv *ZoomImageView) ZoomOut() {
	ziv.mu.RLock()
	scale := ziv.CurrentZoomScale / 1.25
	ziv.mu.RUnlock()
	ziv.SetZoomScale(scale)
}

// ResetZoom resets to default zoom and position
func (ziv *ZoomImageView) ResetZoom() {
	ziv.mu.Lock()
	ziv.CurrentZoomScale = 1.0
	ziv.OffsetX = 0
	ziv.OffsetY = 0
	ziv.mu.Unlock()
	ziv.Refresh()

	if ziv.OnZoomChanged != nil {
		ziv.OnZoomChanged(1.0)
	}
	if ziv.OnPanChanged != nil {
		ziv.OnPanChanged(0, 0)
	}
}

// FitToView scales the image to fit within the view
func (ziv *ZoomImageView) FitToView() {
	ziv.mu.Lock()
	ziv.CurrentZoomScale = 1.0
	ziv.OffsetX = 0
	ziv.OffsetY = 0
	ziv.mu.Unlock()
	ziv.Refresh()
}

// FillView scales the image to fill the view
func (ziv *ZoomImageView) FillView() {
	ziv.mu.Lock()
	viewSize := ziv.Size()
	imgSize := ziv.imageSize

	if imgSize.Width == 0 || imgSize.Height == 0 {
		ziv.mu.Unlock()
		return
	}

	scaleX := viewSize.Width / imgSize.Width
	scaleY := viewSize.Height / imgSize.Height
	scale := float32(math.Max(float64(scaleX), float64(scaleY)))

	ziv.CurrentZoomScale = scale
	ziv.OffsetX = 0
	ziv.OffsetY = 0
	ziv.mu.Unlock()
	ziv.Refresh()

	if ziv.OnZoomChanged != nil {
		ziv.OnZoomChanged(scale)
	}
}

// Pan moves the image by the given offset
func (ziv *ZoomImageView) Pan(deltaX, deltaY float32) {
	if !ziv.PanEnabled {
		return
	}

	ziv.mu.Lock()
	ziv.OffsetX += deltaX
	ziv.OffsetY += deltaY
	ziv.constrainOffset()
	offsetX := ziv.OffsetX
	offsetY := ziv.OffsetY
	ziv.mu.Unlock()
	ziv.Refresh()

	if ziv.OnPanChanged != nil {
		ziv.OnPanChanged(offsetX, offsetY)
	}
}

func (ziv *ZoomImageView) constrainOffset() {
	viewSize := ziv.Size()
	scaledWidth := ziv.imageSize.Width * ziv.CurrentZoomScale
	scaledHeight := ziv.imageSize.Height * ziv.CurrentZoomScale

	// Calculate max offsets
	maxOffsetX := float32(math.Max(0, float64(scaledWidth-viewSize.Width)/2))
	maxOffsetY := float32(math.Max(0, float64(scaledHeight-viewSize.Height)/2))

	// Constrain
	if ziv.OffsetX > maxOffsetX {
		ziv.OffsetX = maxOffsetX
	}
	if ziv.OffsetX < -maxOffsetX {
		ziv.OffsetX = -maxOffsetX
	}
	if ziv.OffsetY > maxOffsetY {
		ziv.OffsetY = maxOffsetY
	}
	if ziv.OffsetY < -maxOffsetY {
		ziv.OffsetY = -maxOffsetY
	}
}

// DoubleTapped implements fyne.DoubleTappable
func (ziv *ZoomImageView) DoubleTapped(e *fyne.PointEvent) {
	if !ziv.DoubleTapEnabled {
		return
	}

	ziv.mu.RLock()
	currentScale := ziv.CurrentZoomScale
	targetScale := ziv.DoubleTapZoomScale
	ziv.mu.RUnlock()

	if currentScale > 1.0 {
		ziv.ResetZoom()
	} else {
		ziv.SetZoomScale(targetScale)
	}
}

// Scrolled implements fyne.Scrollable for mouse wheel zoom
func (ziv *ZoomImageView) Scrolled(e *fyne.ScrollEvent) {
	ziv.mu.RLock()
	scale := ziv.CurrentZoomScale
	ziv.mu.RUnlock()

	if e.Scrolled.DY > 0 {
		scale *= 1.1
	} else if e.Scrolled.DY < 0 {
		scale /= 1.1
	}

	ziv.SetZoomScale(scale)
}

// MouseDown implements desktop.Mouseable
func (ziv *ZoomImageView) MouseDown(e *desktop.MouseEvent) {
	ziv.mu.Lock()
	ziv.isDragging = true
	ziv.lastMousePos = e.Position
	ziv.mu.Unlock()
}

// MouseUp implements desktop.Mouseable
func (ziv *ZoomImageView) MouseUp(e *desktop.MouseEvent) {
	ziv.mu.Lock()
	ziv.isDragging = false
	ziv.mu.Unlock()
}

// Dragged implements fyne.Draggable
func (ziv *ZoomImageView) Dragged(e *fyne.DragEvent) {
	if !ziv.PanEnabled {
		return
	}

	ziv.Pan(e.Dragged.DX, e.Dragged.DY)
}

// DragEnd implements fyne.Draggable
func (ziv *ZoomImageView) DragEnd() {
	ziv.mu.Lock()
	ziv.isDragging = false
	ziv.mu.Unlock()
}

// CreateRenderer implements fyne.Widget
func (ziv *ZoomImageView) CreateRenderer() fyne.WidgetRenderer {
	ziv.ExtendBaseWidget(ziv)

	var img *canvas.Image
	if ziv.Image != nil {
		img = canvas.NewImageFromResource(ziv.Image)
		img.FillMode = ziv.FillMode
	}

	return &zoomImageRenderer{
		view:  ziv,
		image: img,
	}
}

type zoomImageRenderer struct {
	view  *ZoomImageView
	image *canvas.Image
}

func (r *zoomImageRenderer) Destroy() {}

func (r *zoomImageRenderer) Layout(size fyne.Size) {
	r.view.mu.Lock()

	if r.image == nil && r.view.Image != nil {
		r.image = canvas.NewImageFromResource(r.view.Image)
		r.image.FillMode = r.view.FillMode
	}

	if r.image == nil {
		r.view.mu.Unlock()
		return
	}

	// Update image resource if changed
	if r.view.Image != nil {
		r.image.Resource = r.view.Image
	}

	// Calculate scaled image size
	imgMinSize := r.image.MinSize()
	if imgMinSize.Width == 0 {
		imgMinSize = fyne.NewSize(size.Width, size.Height)
	}
	r.view.imageSize = imgMinSize

	scale := r.view.CurrentZoomScale
	offsetX := r.view.OffsetX
	offsetY := r.view.OffsetY

	scaledWidth := imgMinSize.Width * scale
	scaledHeight := imgMinSize.Height * scale

	// Center the image with offset
	x := (size.Width-scaledWidth)/2 + offsetX
	y := (size.Height-scaledHeight)/2 + offsetY

	r.view.mu.Unlock()

	r.image.Move(fyne.NewPos(x, y))
	r.image.Resize(fyne.NewSize(scaledWidth, scaledHeight))
}

func (r *zoomImageRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (r *zoomImageRenderer) Refresh() {
	if r.image != nil {
		r.view.mu.RLock()
		if r.view.Image != nil {
			r.image.Resource = r.view.Image
		}
		r.image.FillMode = r.view.FillMode
		r.view.mu.RUnlock()
		r.image.Refresh()
	}
}

func (r *zoomImageRenderer) Objects() []fyne.CanvasObject {
	if r.image != nil {
		return []fyne.CanvasObject{r.image}
	}
	return nil
}

// ZoomImageViewWrapper wraps a ZoomImageView with additional controls
type ZoomImageViewWrapper struct {
	widget.BaseWidget

	ZoomView    *ZoomImageView
	ShowsControls bool

	// Styling
	ControlsBackgroundColor color.Color
	ControlsTextColor       color.Color
}

// NewZoomImageViewWrapper creates a wrapped zoom image view
func NewZoomImageViewWrapper(img fyne.Resource) *ZoomImageViewWrapper {
	wrapper := &ZoomImageViewWrapper{
		ZoomView:                NewZoomImageViewWithResource(img),
		ShowsControls:           true,
		ControlsBackgroundColor: color.RGBA{R: 0, G: 0, B: 0, A: 150},
		ControlsTextColor:       color.White,
	}
	wrapper.ExtendBaseWidget(wrapper)
	return wrapper
}

// CreateRenderer implements fyne.Widget
func (w *ZoomImageViewWrapper) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	bg := canvas.NewRectangle(color.Transparent)

	return &zoomImageWrapperRenderer{
		wrapper: w,
		bg:      bg,
	}
}

type zoomImageWrapperRenderer struct {
	wrapper *ZoomImageViewWrapper
	bg      *canvas.Rectangle
}

func (r *zoomImageWrapperRenderer) Destroy() {}

func (r *zoomImageWrapperRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.wrapper.ZoomView.Resize(size)
}

func (r *zoomImageWrapperRenderer) MinSize() fyne.Size {
	return r.wrapper.ZoomView.MinSize()
}

func (r *zoomImageWrapperRenderer) Refresh() {
	r.wrapper.ZoomView.Refresh()
}

func (r *zoomImageWrapperRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.wrapper.ZoomView}
}

// GestureZoomHandler handles pinch-to-zoom gestures
type GestureZoomHandler struct {
	View         *ZoomImageView
	InitialScale float32

	mu sync.RWMutex
}

// NewGestureZoomHandler creates a gesture zoom handler
func NewGestureZoomHandler(view *ZoomImageView) *GestureZoomHandler {
	return &GestureZoomHandler{
		View:         view,
		InitialScale: 1.0,
	}
}

// BeginGesture starts a pinch gesture
func (h *GestureZoomHandler) BeginGesture() {
	h.mu.Lock()
	h.InitialScale = h.View.CurrentZoomScale
	h.mu.Unlock()
}

// UpdateGesture updates the zoom based on pinch scale
func (h *GestureZoomHandler) UpdateGesture(scale float32) {
	h.mu.RLock()
	initial := h.InitialScale
	h.mu.RUnlock()

	newScale := initial * scale
	h.View.SetZoomScale(newScale)
}

// EndGesture ends a pinch gesture
func (h *GestureZoomHandler) EndGesture() {
	// Optional: snap to nice zoom levels
}

// ImageViewerController provides a full image viewer experience
type ImageViewerController struct {
	Images       []fyne.Resource
	CurrentIndex int

	ZoomView *ZoomImageView

	OnIndexChanged func(index int)
	OnDismiss      func()

	window fyne.Window
	popup  *widget.PopUp

	mu sync.RWMutex
}

// NewImageViewerController creates an image viewer controller
func NewImageViewerController(images []fyne.Resource) *ImageViewerController {
	ivc := &ImageViewerController{
		Images:       images,
		CurrentIndex: 0,
		ZoomView:     NewZoomImageView(),
	}

	if len(images) > 0 {
		ivc.ZoomView.SetImage(images[0])
	}

	return ivc
}

// Show displays the image viewer
func (ivc *ImageViewerController) Show(window fyne.Window) {
	ivc.window = window

	ivc.popup = widget.NewModalPopUp(ivc.ZoomView, window.Canvas())
	ivc.popup.Resize(window.Canvas().Size())
	ivc.popup.Show()
}

// Hide hides the image viewer
func (ivc *ImageViewerController) Hide() {
	if ivc.popup != nil {
		ivc.popup.Hide()
		ivc.popup = nil
	}

	if ivc.OnDismiss != nil {
		ivc.OnDismiss()
	}
}

// Next shows the next image
func (ivc *ImageViewerController) Next() {
	ivc.mu.Lock()
	if ivc.CurrentIndex < len(ivc.Images)-1 {
		ivc.CurrentIndex++
		ivc.ZoomView.SetImage(ivc.Images[ivc.CurrentIndex])
	}
	index := ivc.CurrentIndex
	ivc.mu.Unlock()

	if ivc.OnIndexChanged != nil {
		ivc.OnIndexChanged(index)
	}
}

// Previous shows the previous image
func (ivc *ImageViewerController) Previous() {
	ivc.mu.Lock()
	if ivc.CurrentIndex > 0 {
		ivc.CurrentIndex--
		ivc.ZoomView.SetImage(ivc.Images[ivc.CurrentIndex])
	}
	index := ivc.CurrentIndex
	ivc.mu.Unlock()

	if ivc.OnIndexChanged != nil {
		ivc.OnIndexChanged(index)
	}
}

// GoToIndex jumps to a specific image
func (ivc *ImageViewerController) GoToIndex(index int) {
	ivc.mu.Lock()
	if index >= 0 && index < len(ivc.Images) {
		ivc.CurrentIndex = index
		ivc.ZoomView.SetImage(ivc.Images[index])
	}
	ivc.mu.Unlock()

	if ivc.OnIndexChanged != nil {
		ivc.OnIndexChanged(index)
	}
}
