// Package imagepicker provides QMUIImagePickerViewController - image picker functionality
// Ported from Tencent's QMUI_iOS framework
package imagepicker

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/core"
)

// ImageAsset represents a selected image
type ImageAsset struct {
	URI         fyne.URI
	Path        string
	Image       image.Image
	Thumbnail   image.Image
	Selected    bool
	Index       int
}

// ImagePickerDelegate provides callbacks for image picker events
type ImagePickerDelegate interface {
	ImagePickerDidSelectImages(picker *ImagePickerViewController, images []*ImageAsset)
	ImagePickerDidCancel(picker *ImagePickerViewController)
	ImagePickerShouldSelectImage(picker *ImagePickerViewController, asset *ImageAsset) bool
}

// ImagePickerViewController manages image selection
type ImagePickerViewController struct {
	widget.BaseWidget

	// Configuration
	MaximumSelectCount int
	MinimumSelectCount int
	AllowsMultipleSelection bool
	ShowsPreview       bool
	PreviewSize        fyne.Size

	// Styling
	BackgroundColor           color.Color
	CellBackgroundColor       color.Color
	CellSelectedBorderColor   color.Color
	CellSelectedBorderWidth   float32
	ThumbnailSize             fyne.Size
	CellSpacing               float32
	ColumnsCount              int
	CheckboxSize              float32
	CheckboxBackgroundColor   color.Color
	CheckboxSelectedColor     color.Color
	ToolbarBackgroundColor    color.Color
	ToolbarHeight             float32

	// Content
	Images         []*ImageAsset
	SelectedImages []*ImageAsset

	// Delegate
	Delegate ImagePickerDelegate

	// Callbacks
	OnImagesSelected func(images []*ImageAsset)
	OnCancel         func()

	// State
	mu      sync.RWMutex
	window  fyne.Window
	popup   *widget.PopUp
}

// NewImagePickerViewController creates a new image picker
func NewImagePickerViewController() *ImagePickerViewController {
	config := core.SharedConfiguration()
	ipvc := &ImagePickerViewController{
		MaximumSelectCount:       9,
		MinimumSelectCount:       1,
		AllowsMultipleSelection:  true,
		ShowsPreview:             true,
		PreviewSize:              fyne.NewSize(300, 300),
		BackgroundColor:          config.BackgroundColor,
		CellBackgroundColor:      color.RGBA{R: 240, G: 240, B: 240, A: 255},
		CellSelectedBorderColor:  config.BlueColor,
		CellSelectedBorderWidth:  3,
		ThumbnailSize:            fyne.NewSize(80, 80),
		CellSpacing:              2,
		ColumnsCount:             4,
		CheckboxSize:             24,
		CheckboxBackgroundColor:  color.RGBA{R: 0, G: 0, B: 0, A: 100},
		CheckboxSelectedColor:    config.BlueColor,
		ToolbarBackgroundColor:   config.ToolBarBackgroundColor,
		ToolbarHeight:            44,
		Images:                   make([]*ImageAsset, 0),
		SelectedImages:           make([]*ImageAsset, 0),
	}
	ipvc.ExtendBaseWidget(ipvc)
	return ipvc
}

// SetImages sets the available images
func (ipvc *ImagePickerViewController) SetImages(images []*ImageAsset) {
	ipvc.mu.Lock()
	ipvc.Images = images
	ipvc.mu.Unlock()
	ipvc.Refresh()
}

// SelectImage selects an image
func (ipvc *ImagePickerViewController) SelectImage(asset *ImageAsset) bool {
	ipvc.mu.Lock()
	defer ipvc.mu.Unlock()

	// Check if already at max
	if len(ipvc.SelectedImages) >= ipvc.MaximumSelectCount {
		return false
	}

	// Check delegate
	if ipvc.Delegate != nil && !ipvc.Delegate.ImagePickerShouldSelectImage(ipvc, asset) {
		return false
	}

	// Toggle selection
	if asset.Selected {
		// Deselect
		asset.Selected = false
		for i, img := range ipvc.SelectedImages {
			if img == asset {
				ipvc.SelectedImages = append(ipvc.SelectedImages[:i], ipvc.SelectedImages[i+1:]...)
				break
			}
		}
	} else {
		// Select
		asset.Selected = true
		ipvc.SelectedImages = append(ipvc.SelectedImages, asset)
	}

	return true
}

// GetSelectedCount returns the number of selected images
func (ipvc *ImagePickerViewController) GetSelectedCount() int {
	ipvc.mu.RLock()
	defer ipvc.mu.RUnlock()
	return len(ipvc.SelectedImages)
}

// ClearSelection clears all selections
func (ipvc *ImagePickerViewController) ClearSelection() {
	ipvc.mu.Lock()
	for _, img := range ipvc.SelectedImages {
		img.Selected = false
	}
	ipvc.SelectedImages = make([]*ImageAsset, 0)
	ipvc.mu.Unlock()
	ipvc.Refresh()
}

// Confirm confirms the selection
func (ipvc *ImagePickerViewController) Confirm() {
	ipvc.mu.RLock()
	selected := make([]*ImageAsset, len(ipvc.SelectedImages))
	copy(selected, ipvc.SelectedImages)
	ipvc.mu.RUnlock()

	if ipvc.OnImagesSelected != nil {
		ipvc.OnImagesSelected(selected)
	}
	if ipvc.Delegate != nil {
		ipvc.Delegate.ImagePickerDidSelectImages(ipvc, selected)
	}
}

// Cancel cancels the picker
func (ipvc *ImagePickerViewController) Cancel() {
	if ipvc.OnCancel != nil {
		ipvc.OnCancel()
	}
	if ipvc.Delegate != nil {
		ipvc.Delegate.ImagePickerDidCancel(ipvc)
	}
}

// ShowIn displays the picker in a window
func (ipvc *ImagePickerViewController) ShowIn(window fyne.Window) {
	ipvc.mu.Lock()
	ipvc.window = window
	ipvc.mu.Unlock()

	content := ipvc.buildContent()

	ipvc.popup = widget.NewModalPopUp(content, window.Canvas())
	ipvc.popup.Resize(window.Canvas().Size())
	ipvc.popup.Show()
}

// Dismiss hides the picker
func (ipvc *ImagePickerViewController) Dismiss() {
	ipvc.mu.Lock()
	if ipvc.popup != nil {
		ipvc.popup.Hide()
		ipvc.popup = nil
	}
	ipvc.mu.Unlock()
}

func (ipvc *ImagePickerViewController) buildContent() fyne.CanvasObject {
	// Toolbar
	toolbar := ipvc.buildToolbar()

	// Grid of images
	grid := ipvc.buildImageGrid()

	// Layout
	content := container.NewBorder(toolbar, nil, nil, nil, container.NewScroll(grid))

	return content
}

func (ipvc *ImagePickerViewController) buildToolbar() fyne.CanvasObject {
	bg := canvas.NewRectangle(ipvc.ToolbarBackgroundColor)

	cancelBtn := widget.NewButton("Cancel", func() {
		ipvc.Cancel()
		ipvc.Dismiss()
	})

	title := widget.NewLabel("Select Images")
	title.Alignment = fyne.TextAlignCenter

	confirmBtn := widget.NewButton("Done", func() {
		ipvc.Confirm()
		ipvc.Dismiss()
	})

	toolbar := container.NewBorder(nil, nil, cancelBtn, confirmBtn, title)

	return container.NewStack(bg, toolbar)
}

func (ipvc *ImagePickerViewController) buildImageGrid() fyne.CanvasObject {
	ipvc.mu.RLock()
	images := ipvc.Images
	cols := ipvc.ColumnsCount
	ipvc.mu.RUnlock()

	var cells []fyne.CanvasObject
	for _, img := range images {
		cell := ipvc.buildImageCell(img)
		cells = append(cells, cell)
	}

	grid := container.NewGridWithColumns(cols, cells...)
	return grid
}

func (ipvc *ImagePickerViewController) buildImageCell(asset *ImageAsset) fyne.CanvasObject {
	cell := &imagePickerCell{
		picker: ipvc,
		asset:  asset,
	}
	cell.ExtendBaseWidget(cell)
	return cell
}

// CreateRenderer implements fyne.Widget
func (ipvc *ImagePickerViewController) CreateRenderer() fyne.WidgetRenderer {
	ipvc.ExtendBaseWidget(ipvc)
	return &imagePickerRenderer{picker: ipvc}
}

type imagePickerRenderer struct {
	picker *ImagePickerViewController
}

func (r *imagePickerRenderer) Destroy()              {}
func (r *imagePickerRenderer) Layout(size fyne.Size) {}
func (r *imagePickerRenderer) MinSize() fyne.Size    { return fyne.NewSize(0, 0) }
func (r *imagePickerRenderer) Refresh()              {}
func (r *imagePickerRenderer) Objects() []fyne.CanvasObject { return nil }

// imagePickerCell represents a single image cell
type imagePickerCell struct {
	widget.BaseWidget
	picker  *ImagePickerViewController
	asset   *ImageAsset
	hovered bool
	mu      sync.RWMutex
}

func (c *imagePickerCell) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	bg := canvas.NewRectangle(c.picker.CellBackgroundColor)

	// Placeholder image
	placeholder := canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255})

	// Selection checkbox
	checkbox := canvas.NewCircle(c.picker.CheckboxBackgroundColor)

	// Selection border
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = c.picker.CellSelectedBorderWidth
	border.StrokeColor = c.picker.CellSelectedBorderColor

	return &imagePickerCellRenderer{
		cell:        c,
		bg:          bg,
		placeholder: placeholder,
		checkbox:    checkbox,
		border:      border,
	}
}

func (c *imagePickerCell) Tapped(*fyne.PointEvent) {
	c.picker.SelectImage(c.asset)
	c.picker.Refresh()
}

func (c *imagePickerCell) TappedSecondary(*fyne.PointEvent) {}

func (c *imagePickerCell) MouseIn(*desktop.MouseEvent) {
	c.mu.Lock()
	c.hovered = true
	c.mu.Unlock()
	c.Refresh()
}

func (c *imagePickerCell) MouseMoved(*desktop.MouseEvent) {}

func (c *imagePickerCell) MouseOut() {
	c.mu.Lock()
	c.hovered = false
	c.mu.Unlock()
	c.Refresh()
}

func (c *imagePickerCell) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type imagePickerCellRenderer struct {
	cell        *imagePickerCell
	bg          *canvas.Rectangle
	placeholder *canvas.Rectangle
	checkbox    *canvas.Circle
	border      *canvas.Rectangle
	image       *canvas.Image
}

func (r *imagePickerCellRenderer) Destroy() {}

func (r *imagePickerCellRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.placeholder.Resize(size)
	r.border.Resize(size)

	// Checkbox in top-right
	checkSize := r.cell.picker.CheckboxSize
	r.checkbox.Resize(fyne.NewSize(checkSize, checkSize))
	r.checkbox.Move(fyne.NewPos(size.Width-checkSize-4, 4))

	// Image if available
	if r.image != nil {
		r.image.Resize(size)
	}
}

func (r *imagePickerCellRenderer) MinSize() fyne.Size {
	return r.cell.picker.ThumbnailSize
}

func (r *imagePickerCellRenderer) Refresh() {
	r.cell.mu.RLock()
	hovered := r.cell.hovered
	r.cell.mu.RUnlock()

	if r.cell.asset.Selected {
		r.border.StrokeColor = r.cell.picker.CellSelectedBorderColor
		r.checkbox.FillColor = r.cell.picker.CheckboxSelectedColor
	} else {
		r.border.StrokeColor = color.Transparent
		r.checkbox.FillColor = r.cell.picker.CheckboxBackgroundColor
	}

	if hovered {
		r.bg.FillColor = core.ColorWithAlpha(r.cell.picker.CellBackgroundColor, 0.8)
	} else {
		r.bg.FillColor = r.cell.picker.CellBackgroundColor
	}

	r.bg.Refresh()
	r.placeholder.Refresh()
	r.checkbox.Refresh()
	r.border.Refresh()
}

func (r *imagePickerCellRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bg, r.placeholder, r.border, r.checkbox}
	if r.image != nil {
		objects = append([]fyne.CanvasObject{r.bg, r.image, r.border, r.checkbox}, nil)
		objects = objects[:4]
	}
	return objects
}

// Helper functions

// ShowImagePicker shows a native file picker for images
func ShowImagePicker(window fyne.Window, callback func([]fyne.URI)) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			callback(nil)
			return
		}
		callback([]fyne.URI{reader.URI()})
		reader.Close()
	}, window)
}

// ShowMultiImagePicker shows a picker for multiple images
func ShowMultiImagePicker(window fyne.Window, callback func([]fyne.URI)) {
	// Fyne doesn't have native multi-file picker, so we use the image picker view
	picker := NewImagePickerViewController()
	picker.OnImagesSelected = func(images []*ImageAsset) {
		uris := make([]fyne.URI, len(images))
		for i, img := range images {
			uris[i] = img.URI
		}
		callback(uris)
	}
	picker.OnCancel = func() {
		callback(nil)
	}
	picker.ShowIn(window)
}

// LoadImagesFromDirectory loads images from a directory
func LoadImagesFromDirectory(dirPath string) ([]*ImageAsset, error) {
	var assets []*ImageAsset
	index := 0

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Check if it's an image
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" {
			uri := storage.NewFileURI(path)
			asset := &ImageAsset{
				URI:   uri,
				Path:  path,
				Index: index,
			}
			assets = append(assets, asset)
			index++
		}

		return nil
	})

	return assets, err
}

// CreateImageAsset creates an image asset from a file path
func CreateImageAsset(path string) *ImageAsset {
	uri := storage.NewFileURI(path)
	return &ImageAsset{
		URI:  uri,
		Path: path,
	}
}

// CreateImageAssetsFromPaths creates image assets from file paths
func CreateImageAssetsFromPaths(paths []string) []*ImageAsset {
	assets := make([]*ImageAsset, len(paths))
	for i, path := range paths {
		assets[i] = &ImageAsset{
			URI:   storage.NewFileURI(path),
			Path:  path,
			Index: i,
		}
	}
	return assets
}
