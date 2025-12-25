// Package album provides QMUIAlbumViewController - album browser functionality
// Ported from Tencent's QMUI_iOS framework
package album

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/core"
)

// Album represents a photo album
type Album struct {
	Identifier    string
	Name          string
	ThumbnailPath string
	PhotoCount    int
	Photos        []*Photo
}

// Photo represents a photo in an album
type Photo struct {
	URI        fyne.URI
	Path       string
	Album      *Album
	CreateDate int64
}

// AlbumViewDelegate provides callbacks for album view events
type AlbumViewDelegate interface {
	AlbumViewDidSelectAlbum(view *AlbumView, album *Album)
	AlbumViewDidSelectPhoto(view *AlbumView, photo *Photo)
}

// AlbumView displays a list of albums
type AlbumView struct {
	widget.BaseWidget

	// Content
	Albums []*Album

	// Styling
	BackgroundColor       color.Color
	CellBackgroundColor   color.Color
	CellHighlightColor    color.Color
	TitleColor            color.Color
	CountColor            color.Color
	TitleFontSize         float32
	CountFontSize         float32
	CellHeight            float32
	ThumbnailSize         fyne.Size
	ThumbnailCornerRadius float32
	SeparatorColor        color.Color

	// Delegate
	Delegate AlbumViewDelegate

	// Callbacks
	OnAlbumSelected func(album *Album)
	OnPhotoSelected func(photo *Photo)

	// State
	mu       sync.RWMutex
	window   fyne.Window
}

// NewAlbumView creates a new album view
func NewAlbumView() *AlbumView {
	config := core.SharedConfiguration()
	av := &AlbumView{
		Albums:                make([]*Album, 0),
		BackgroundColor:       config.BackgroundColor,
		CellBackgroundColor:   config.TableViewCellBackgroundColor,
		CellHighlightColor:    config.TableViewCellSelectedBackgroundColor,
		TitleColor:            config.TableViewCellTitleLabelColor,
		CountColor:            config.TableViewCellDetailLabelColor,
		TitleFontSize:         16,
		CountFontSize:         14,
		CellHeight:            60,
		ThumbnailSize:         fyne.NewSize(50, 50),
		ThumbnailCornerRadius: 4,
		SeparatorColor:        config.SeparatorColor,
	}
	av.ExtendBaseWidget(av)
	return av
}

// NewAlbumViewWithAlbums creates an album view with albums
func NewAlbumViewWithAlbums(albums []*Album) *AlbumView {
	av := NewAlbumView()
	av.Albums = albums
	return av
}

// SetAlbums sets the albums to display
func (av *AlbumView) SetAlbums(albums []*Album) {
	av.mu.Lock()
	av.Albums = albums
	av.mu.Unlock()
	av.Refresh()
}

// AddAlbum adds an album
func (av *AlbumView) AddAlbum(album *Album) {
	av.mu.Lock()
	av.Albums = append(av.Albums, album)
	av.mu.Unlock()
	av.Refresh()
}

// SelectAlbum handles album selection
func (av *AlbumView) SelectAlbum(album *Album) {
	if av.OnAlbumSelected != nil {
		av.OnAlbumSelected(album)
	}
	if av.Delegate != nil {
		av.Delegate.AlbumViewDidSelectAlbum(av, album)
	}
}

// Show displays the album view
func (av *AlbumView) ShowIn(window fyne.Window) {
	av.mu.Lock()
	av.window = window
	av.mu.Unlock()
}

// CreateRenderer implements fyne.Widget
func (av *AlbumView) CreateRenderer() fyne.WidgetRenderer {
	av.ExtendBaseWidget(av)

	background := canvas.NewRectangle(av.BackgroundColor)

	return &albumViewRenderer{
		view:       av,
		background: background,
		cells:      make([]*albumCell, 0),
	}
}

type albumViewRenderer struct {
	view       *AlbumView
	background *canvas.Rectangle
	cells      []*albumCell
}

func (r *albumViewRenderer) Destroy() {}

func (r *albumViewRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.rebuildCells(size)
}

func (r *albumViewRenderer) rebuildCells(size fyne.Size) {
	r.view.mu.RLock()
	albums := r.view.Albums
	cellHeight := r.view.CellHeight
	r.view.mu.RUnlock()

	// Rebuild cells if count changed
	if len(r.cells) != len(albums) {
		r.cells = make([]*albumCell, len(albums))
		for i, album := range albums {
			cell := &albumCell{
				view:  r.view,
				album: album,
			}
			cell.ExtendBaseWidget(cell)
			r.cells[i] = cell
		}
	}

	// Position cells
	y := float32(0)
	for _, cell := range r.cells {
		cell.Resize(fyne.NewSize(size.Width, cellHeight))
		cell.Move(fyne.NewPos(0, y))
		y += cellHeight
	}
}

func (r *albumViewRenderer) MinSize() fyne.Size {
	r.view.mu.RLock()
	count := len(r.view.Albums)
	cellHeight := r.view.CellHeight
	r.view.mu.RUnlock()

	return fyne.NewSize(200, float32(count)*cellHeight)
}

func (r *albumViewRenderer) Refresh() {
	r.background.FillColor = r.view.BackgroundColor
	r.background.Refresh()

	for _, cell := range r.cells {
		cell.Refresh()
	}
}

func (r *albumViewRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}
	for _, cell := range r.cells {
		objects = append(objects, cell)
	}
	return objects
}

// albumCell represents a single album cell
type albumCell struct {
	widget.BaseWidget
	view    *AlbumView
	album   *Album
	hovered bool
	mu      sync.RWMutex
}

func (c *albumCell) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	bg := canvas.NewRectangle(c.view.CellBackgroundColor)

	thumbnail := canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	thumbnail.CornerRadius = c.view.ThumbnailCornerRadius

	title := canvas.NewText(c.album.Name, c.view.TitleColor)
	title.TextSize = c.view.TitleFontSize
	title.TextStyle = fyne.TextStyle{Bold: true}

	count := canvas.NewText("", c.view.CountColor)
	count.TextSize = c.view.CountFontSize

	separator := canvas.NewRectangle(c.view.SeparatorColor)

	disclosure := canvas.NewText(">", c.view.CountColor)
	disclosure.TextSize = c.view.TitleFontSize

	return &albumCellRenderer{
		cell:       c,
		bg:         bg,
		thumbnail:  thumbnail,
		title:      title,
		count:      count,
		separator:  separator,
		disclosure: disclosure,
	}
}

func (c *albumCell) Tapped(*fyne.PointEvent) {
	c.view.SelectAlbum(c.album)
}

func (c *albumCell) TappedSecondary(*fyne.PointEvent) {}

func (c *albumCell) MouseIn(*desktop.MouseEvent) {
	c.mu.Lock()
	c.hovered = true
	c.mu.Unlock()
	c.Refresh()
}

func (c *albumCell) MouseMoved(*desktop.MouseEvent) {}

func (c *albumCell) MouseOut() {
	c.mu.Lock()
	c.hovered = false
	c.mu.Unlock()
	c.Refresh()
}

func (c *albumCell) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type albumCellRenderer struct {
	cell       *albumCell
	bg         *canvas.Rectangle
	thumbnail  *canvas.Rectangle
	title      *canvas.Text
	count      *canvas.Text
	separator  *canvas.Rectangle
	disclosure *canvas.Text
	image      *canvas.Image
}

func (r *albumCellRenderer) Destroy() {}

func (r *albumCellRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	thumbSize := r.cell.view.ThumbnailSize
	padding := (size.Height - thumbSize.Height) / 2

	// Thumbnail
	r.thumbnail.Resize(thumbSize)
	r.thumbnail.Move(fyne.NewPos(padding, padding))

	// Title and count
	titleX := padding + thumbSize.Width + padding
	titleY := padding
	r.title.Move(fyne.NewPos(titleX, titleY))

	countY := size.Height - padding - r.count.MinSize().Height
	r.count.Move(fyne.NewPos(titleX, countY))

	// Disclosure
	disclosureSize := r.disclosure.MinSize()
	r.disclosure.Move(fyne.NewPos(size.Width-padding-disclosureSize.Width, (size.Height-disclosureSize.Height)/2))

	// Separator
	r.separator.Resize(fyne.NewSize(size.Width-titleX, 0.5))
	r.separator.Move(fyne.NewPos(titleX, size.Height-0.5))
}

func (r *albumCellRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, r.cell.view.CellHeight)
}

func (r *albumCellRenderer) Refresh() {
	r.cell.mu.RLock()
	hovered := r.cell.hovered
	r.cell.mu.RUnlock()

	if hovered {
		r.bg.FillColor = r.cell.view.CellHighlightColor
	} else {
		r.bg.FillColor = r.cell.view.CellBackgroundColor
	}

	r.title.Text = r.cell.album.Name
	r.title.Color = r.cell.view.TitleColor

	if r.cell.album.PhotoCount > 0 {
		r.count.Text = string(rune('0'+r.cell.album.PhotoCount%10)) + " photos"
		if r.cell.album.PhotoCount >= 10 {
			r.count.Text = string(rune('0'+r.cell.album.PhotoCount/10)) + string(rune('0'+r.cell.album.PhotoCount%10)) + " photos"
		}
	} else {
		r.count.Text = "Empty"
	}

	r.bg.Refresh()
	r.thumbnail.Refresh()
	r.title.Refresh()
	r.count.Refresh()
	r.separator.Refresh()
	r.disclosure.Refresh()
}

func (r *albumCellRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.thumbnail, r.title, r.count, r.separator, r.disclosure}
}

// PhotoGridView displays photos in a grid
type PhotoGridView struct {
	widget.BaseWidget

	// Content
	Photos []*Photo
	Album  *Album

	// Layout
	ColumnsCount int
	PhotoSpacing float32
	PhotoSize    fyne.Size

	// Styling
	BackgroundColor         color.Color
	CellBackgroundColor     color.Color
	CellSelectedBorderColor color.Color
	CellSelectedBorderWidth float32

	// Callbacks
	OnPhotoSelected func(photo *Photo)

	// State
	mu sync.RWMutex
}

// NewPhotoGridView creates a new photo grid view
func NewPhotoGridView() *PhotoGridView {
	config := core.SharedConfiguration()
	pgv := &PhotoGridView{
		Photos:                  make([]*Photo, 0),
		ColumnsCount:            4,
		PhotoSpacing:            2,
		PhotoSize:               fyne.NewSize(80, 80),
		BackgroundColor:         config.BackgroundColor,
		CellBackgroundColor:     color.RGBA{R: 240, G: 240, B: 240, A: 255},
		CellSelectedBorderColor: config.BlueColor,
		CellSelectedBorderWidth: 3,
	}
	pgv.ExtendBaseWidget(pgv)
	return pgv
}

// NewPhotoGridViewWithPhotos creates a photo grid with photos
func NewPhotoGridViewWithPhotos(photos []*Photo) *PhotoGridView {
	pgv := NewPhotoGridView()
	pgv.Photos = photos
	return pgv
}

// SetPhotos sets the photos to display
func (pgv *PhotoGridView) SetPhotos(photos []*Photo) {
	pgv.mu.Lock()
	pgv.Photos = photos
	pgv.mu.Unlock()
	pgv.Refresh()
}

// CreateRenderer implements fyne.Widget
func (pgv *PhotoGridView) CreateRenderer() fyne.WidgetRenderer {
	pgv.ExtendBaseWidget(pgv)

	bg := canvas.NewRectangle(pgv.BackgroundColor)

	return &photoGridRenderer{
		view:   pgv,
		bg:     bg,
		cells:  make([]*photoCell, 0),
	}
}

type photoGridRenderer struct {
	view  *PhotoGridView
	bg    *canvas.Rectangle
	cells []*photoCell
}

func (r *photoGridRenderer) Destroy() {}

func (r *photoGridRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.rebuildCells(size)
}

func (r *photoGridRenderer) rebuildCells(size fyne.Size) {
	r.view.mu.RLock()
	photos := r.view.Photos
	cols := r.view.ColumnsCount
	photoSize := r.view.PhotoSize
	spacing := r.view.PhotoSpacing
	r.view.mu.RUnlock()

	// Rebuild cells if count changed
	if len(r.cells) != len(photos) {
		r.cells = make([]*photoCell, len(photos))
		for i, photo := range photos {
			cell := &photoCell{
				view:  r.view,
				photo: photo,
			}
			cell.ExtendBaseWidget(cell)
			r.cells[i] = cell
		}
	}

	// Position cells in grid
	for i, cell := range r.cells {
		col := i % cols
		row := i / cols
		x := float32(col) * (photoSize.Width + spacing)
		y := float32(row) * (photoSize.Height + spacing)
		cell.Resize(photoSize)
		cell.Move(fyne.NewPos(x, y))
	}
}

func (r *photoGridRenderer) MinSize() fyne.Size {
	r.view.mu.RLock()
	count := len(r.view.Photos)
	cols := r.view.ColumnsCount
	photoSize := r.view.PhotoSize
	spacing := r.view.PhotoSpacing
	r.view.mu.RUnlock()

	rows := (count + cols - 1) / cols
	width := float32(cols)*photoSize.Width + float32(cols-1)*spacing
	height := float32(rows)*photoSize.Height + float32(rows-1)*spacing

	return fyne.NewSize(width, height)
}

func (r *photoGridRenderer) Refresh() {
	r.bg.FillColor = r.view.BackgroundColor
	r.bg.Refresh()

	for _, cell := range r.cells {
		cell.Refresh()
	}
}

func (r *photoGridRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bg}
	for _, cell := range r.cells {
		objects = append(objects, cell)
	}
	return objects
}

// photoCell represents a single photo cell
type photoCell struct {
	widget.BaseWidget
	view    *PhotoGridView
	photo   *Photo
	hovered bool
	mu      sync.RWMutex
}

func (c *photoCell) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	bg := canvas.NewRectangle(c.view.CellBackgroundColor)

	return &photoCellRenderer{
		cell: c,
		bg:   bg,
	}
}

func (c *photoCell) Tapped(*fyne.PointEvent) {
	if c.view.OnPhotoSelected != nil {
		c.view.OnPhotoSelected(c.photo)
	}
}

func (c *photoCell) TappedSecondary(*fyne.PointEvent) {}

func (c *photoCell) MouseIn(*desktop.MouseEvent) {
	c.mu.Lock()
	c.hovered = true
	c.mu.Unlock()
	c.Refresh()
}

func (c *photoCell) MouseMoved(*desktop.MouseEvent) {}

func (c *photoCell) MouseOut() {
	c.mu.Lock()
	c.hovered = false
	c.mu.Unlock()
	c.Refresh()
}

func (c *photoCell) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type photoCellRenderer struct {
	cell  *photoCell
	bg    *canvas.Rectangle
	image *canvas.Image
}

func (r *photoCellRenderer) Destroy() {}

func (r *photoCellRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	if r.image != nil {
		r.image.Resize(size)
	}
}

func (r *photoCellRenderer) MinSize() fyne.Size {
	return r.cell.view.PhotoSize
}

func (r *photoCellRenderer) Refresh() {
	r.cell.mu.RLock()
	hovered := r.cell.hovered
	r.cell.mu.RUnlock()

	if hovered {
		r.bg.FillColor = core.ColorWithAlpha(r.cell.view.CellBackgroundColor, 0.7)
	} else {
		r.bg.FillColor = r.cell.view.CellBackgroundColor
	}
	r.bg.Refresh()
}

func (r *photoCellRenderer) Objects() []fyne.CanvasObject {
	if r.image != nil {
		return []fyne.CanvasObject{r.bg, r.image}
	}
	return []fyne.CanvasObject{r.bg}
}

// AlbumViewController is a full album browser view controller
type AlbumViewController struct {
	widget.BaseWidget

	// Views
	AlbumView    *AlbumView
	PhotoGrid    *PhotoGridView

	// State
	CurrentAlbum *Album

	// Callbacks
	OnPhotoSelected func(photo *Photo)
	OnCancel        func()

	// State
	mu     sync.RWMutex
	window fyne.Window
	popup  *widget.PopUp
}

// NewAlbumViewController creates a new album view controller
func NewAlbumViewController() *AlbumViewController {
	avc := &AlbumViewController{
		AlbumView: NewAlbumView(),
		PhotoGrid: NewPhotoGridView(),
	}
	avc.ExtendBaseWidget(avc)

	// Wire up callbacks
	avc.AlbumView.OnAlbumSelected = func(album *Album) {
		avc.ShowAlbum(album)
	}

	avc.PhotoGrid.OnPhotoSelected = func(photo *Photo) {
		if avc.OnPhotoSelected != nil {
			avc.OnPhotoSelected(photo)
		}
	}

	return avc
}

// SetAlbums sets the albums
func (avc *AlbumViewController) SetAlbums(albums []*Album) {
	avc.AlbumView.SetAlbums(albums)
}

// ShowAlbum shows photos from an album
func (avc *AlbumViewController) ShowAlbum(album *Album) {
	avc.mu.Lock()
	avc.CurrentAlbum = album
	avc.mu.Unlock()

	avc.PhotoGrid.SetPhotos(album.Photos)
	avc.Refresh()
}

// BackToAlbums goes back to album list
func (avc *AlbumViewController) BackToAlbums() {
	avc.mu.Lock()
	avc.CurrentAlbum = nil
	avc.mu.Unlock()
	avc.Refresh()
}

// ShowIn displays the album browser
func (avc *AlbumViewController) ShowIn(window fyne.Window) {
	avc.mu.Lock()
	avc.window = window
	avc.mu.Unlock()

	content := avc.buildContent()

	avc.popup = widget.NewModalPopUp(content, window.Canvas())
	avc.popup.Resize(window.Canvas().Size())
	avc.popup.Show()
}

// Dismiss hides the album browser
func (avc *AlbumViewController) Dismiss() {
	avc.mu.Lock()
	if avc.popup != nil {
		avc.popup.Hide()
		avc.popup = nil
	}
	avc.mu.Unlock()
}

func (avc *AlbumViewController) buildContent() fyne.CanvasObject {
	config := core.SharedConfiguration()

	// Toolbar
	bg := canvas.NewRectangle(config.ToolBarBackgroundColor)

	cancelBtn := widget.NewButton("Cancel", func() {
		if avc.OnCancel != nil {
			avc.OnCancel()
		}
		avc.Dismiss()
	})

	avc.mu.RLock()
	currentAlbum := avc.CurrentAlbum
	avc.mu.RUnlock()

	var title *widget.Label
	var backBtn *widget.Button

	if currentAlbum != nil {
		title = widget.NewLabel(currentAlbum.Name)
		backBtn = widget.NewButton("Back", func() {
			avc.BackToAlbums()
		})
	} else {
		title = widget.NewLabel("Albums")
	}
	title.Alignment = fyne.TextAlignCenter

	var toolbar fyne.CanvasObject
	if backBtn != nil {
		toolbar = container.NewBorder(nil, nil, backBtn, cancelBtn, title)
	} else {
		toolbar = container.NewBorder(nil, nil, nil, cancelBtn, title)
	}

	toolbarWithBg := container.NewStack(bg, toolbar)

	// Content
	var content fyne.CanvasObject
	if currentAlbum != nil {
		content = container.NewScroll(avc.PhotoGrid)
	} else {
		content = container.NewScroll(avc.AlbumView)
	}

	return container.NewBorder(toolbarWithBg, nil, nil, nil, content)
}

// CreateRenderer implements fyne.Widget
func (avc *AlbumViewController) CreateRenderer() fyne.WidgetRenderer {
	avc.ExtendBaseWidget(avc)
	return &albumVCRenderer{controller: avc}
}

type albumVCRenderer struct {
	controller *AlbumViewController
}

func (r *albumVCRenderer) Destroy()              {}
func (r *albumVCRenderer) Layout(size fyne.Size) {}
func (r *albumVCRenderer) MinSize() fyne.Size    { return fyne.NewSize(0, 0) }
func (r *albumVCRenderer) Refresh()              {}
func (r *albumVCRenderer) Objects() []fyne.CanvasObject { return nil }

// Helper functions

// LoadAlbumsFromDirectory loads albums from a directory structure
func LoadAlbumsFromDirectory(basePath string) ([]*Album, error) {
	var albums []*Album

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			albumPath := filepath.Join(basePath, entry.Name())
			photos, err := loadPhotosFromDirectory(albumPath)
			if err != nil {
				continue
			}

			album := &Album{
				Identifier: entry.Name(),
				Name:       entry.Name(),
				PhotoCount: len(photos),
				Photos:     photos,
			}

			// Set thumbnail path to first photo
			if len(photos) > 0 {
				album.ThumbnailPath = photos[0].Path
			}

			albums = append(albums, album)
		}
	}

	return albums, nil
}

func loadPhotosFromDirectory(dirPath string) ([]*Photo, error) {
	var photos []*Photo

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" {
			path := filepath.Join(dirPath, entry.Name())
			photo := &Photo{
				URI:  storage.NewFileURI(path),
				Path: path,
			}
			photos = append(photos, photo)
		}
	}

	return photos, nil
}

// CreateAlbum creates an album with the given name and photos
func CreateAlbum(name string, photoPaths []string) *Album {
	photos := make([]*Photo, len(photoPaths))
	for i, path := range photoPaths {
		photos[i] = &Photo{
			URI:  storage.NewFileURI(path),
			Path: path,
		}
	}

	album := &Album{
		Identifier: name,
		Name:       name,
		PhotoCount: len(photos),
		Photos:     photos,
	}

	if len(photoPaths) > 0 {
		album.ThumbnailPath = photoPaths[0]
	}

	return album
}
