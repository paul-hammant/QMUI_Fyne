// Package empty provides empty state views with loading, image, text, and button
package empty

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/core"
)

// EmptyState displays empty state with optional image, text, detail, and action button
type EmptyState struct {
	widget.BaseWidget

	// Content
	Image      fyne.Resource
	Text       string
	DetailText string
	ActionText string

	// Styling
	ImageTintColor       color.Color
	TextColor            color.Color
	DetailTextColor      color.Color
	ActionButtonColor    color.Color
	TextFontSize         float32
	DetailTextFontSize   float32
	ActionButtonFontSize float32
	ImageSize            fyne.Size
	ContentInsets        core.EdgeInsets
	VerticalSpacing      float32
	ImageToTextSpacing   float32
	TextToDetailSpacing  float32
	DetailToButtonSpacing float32

	// Loading state
	IsLoading     bool
	LoadingColor  color.Color

	// Callbacks
	OnActionTapped func()

	mu sync.RWMutex
}

// NewEmptyState creates a new empty view
func NewEmptyState() *EmptyState {
	config := core.SharedConfiguration()
	ev := &EmptyState{
		ImageTintColor:        config.EmptyViewImageTintColor,
		TextColor:             config.EmptyViewTextLabelColor,
		DetailTextColor:       config.EmptyViewDetailTextLabelColor,
		ActionButtonColor:     config.EmptyViewActionButtonColor,
		TextFontSize:          config.EmptyViewTextFontSize,
		DetailTextFontSize:    config.EmptyViewDetailTextFontSize,
		ActionButtonFontSize:  config.EmptyViewActionButtonFontSize,
		ImageSize:             fyne.NewSize(64, 64),
		ContentInsets:         core.NewEdgeInsets(20, 20, 20, 20),
		VerticalSpacing:       16,
		ImageToTextSpacing:    16,
		TextToDetailSpacing:   8,
		DetailToButtonSpacing: 16,
		IsLoading:             false,
		LoadingColor:          config.EmptyViewLoadingTintColor,
	}
	ev.ExtendBaseWidget(ev)
	return ev
}

// NewEmptyStateWithText creates an empty view with text
func NewEmptyStateWithText(text string) *EmptyState {
	ev := NewEmptyState()
	ev.Text = text
	return ev
}

// NewEmptyStateWithTextAndDetail creates an empty view with text and detail
func NewEmptyStateWithTextAndDetail(text, detail string) *EmptyState {
	ev := NewEmptyState()
	ev.Text = text
	ev.DetailText = detail
	return ev
}

// NewEmptyStateWithImageAndText creates an empty view with image and text
func NewEmptyStateWithImageAndText(image fyne.Resource, text string) *EmptyState {
	ev := NewEmptyState()
	ev.Image = image
	ev.Text = text
	return ev
}

// SetLoading sets the loading state
func (ev *EmptyState) SetLoading(loading bool) {
	ev.mu.Lock()
	ev.IsLoading = loading
	ev.mu.Unlock()
	ev.Refresh()
}

// SetText sets the main text
func (ev *EmptyState) SetText(text string) {
	ev.mu.Lock()
	ev.Text = text
	ev.mu.Unlock()
	ev.Refresh()
}

// SetDetailText sets the detail text
func (ev *EmptyState) SetDetailText(detail string) {
	ev.mu.Lock()
	ev.DetailText = detail
	ev.mu.Unlock()
	ev.Refresh()
}

// SetActionText sets the action button text
func (ev *EmptyState) SetActionText(text string) {
	ev.mu.Lock()
	ev.ActionText = text
	ev.mu.Unlock()
	ev.Refresh()
}

// SetImage sets the image
func (ev *EmptyState) SetImage(image fyne.Resource) {
	ev.mu.Lock()
	ev.Image = image
	ev.mu.Unlock()
	ev.Refresh()
}

// CreateRenderer implements fyne.Widget
func (ev *EmptyState) CreateRenderer() fyne.WidgetRenderer {
	ev.ExtendBaseWidget(ev)
	return &emptyViewRenderer{
		emptyView: ev,
	}
}

type emptyViewRenderer struct {
	emptyView *EmptyState
	objects   []fyne.CanvasObject
}

func (r *emptyViewRenderer) Destroy() {}

func (r *emptyViewRenderer) buildObjects() {
	r.objects = nil

	r.emptyView.mu.RLock()
	image := r.emptyView.Image
	text := r.emptyView.Text
	detailText := r.emptyView.DetailText
	actionText := r.emptyView.ActionText
	isLoading := r.emptyView.IsLoading
	r.emptyView.mu.RUnlock()

	var content []fyne.CanvasObject

	// Loading indicator
	if isLoading {
		loading := widget.NewProgressBarInfinite()
		content = append(content, loading)
	}

	// Image
	if image != nil && !isLoading {
		img := canvas.NewImageFromResource(image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(r.emptyView.ImageSize)
		content = append(content, img)
	}

	// Main text
	if text != "" {
		textLabel := canvas.NewText(text, r.emptyView.TextColor)
		textLabel.TextSize = r.emptyView.TextFontSize
		textLabel.Alignment = fyne.TextAlignCenter
		content = append(content, textLabel)
	}

	// Detail text
	if detailText != "" {
		detailLabel := canvas.NewText(detailText, r.emptyView.DetailTextColor)
		detailLabel.TextSize = r.emptyView.DetailTextFontSize
		detailLabel.Alignment = fyne.TextAlignCenter
		content = append(content, detailLabel)
	}

	// Action button
	if actionText != "" {
		actionBtn := widget.NewButton(actionText, func() {
			if r.emptyView.OnActionTapped != nil {
				r.emptyView.OnActionTapped()
			}
		})
		content = append(content, actionBtn)
	}

	if len(content) > 0 {
		vbox := container.NewVBox(content...)
		r.objects = []fyne.CanvasObject{container.NewCenter(vbox)}
	}
}

func (r *emptyViewRenderer) Layout(size fyne.Size) {
	for _, obj := range r.objects {
		obj.Resize(size)
		obj.Move(fyne.NewPos(0, 0))
	}
}

func (r *emptyViewRenderer) MinSize() fyne.Size {
	r.buildObjects()
	if len(r.objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	var width, height float32
	insets := r.emptyView.ContentInsets

	for _, obj := range r.objects {
		size := obj.MinSize()
		if size.Width > width {
			width = size.Width
		}
		height += size.Height
	}

	return fyne.NewSize(
		width+insets.Left+insets.Right,
		height+insets.Top+insets.Bottom,
	)
}

func (r *emptyViewRenderer) Refresh() {
	r.buildObjects()
	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *emptyViewRenderer) Objects() []fyne.CanvasObject {
	r.buildObjects()
	return r.objects
}

// EmptyDataSet provides a data-driven empty view
type EmptyDataSet struct {
	*EmptyState

	// Data binding
	DataSource EmptyDataSetDataSource
}

// EmptyDataSetDataSource provides data for the empty view
type EmptyDataSetDataSource interface {
	ImageForEmptyState() fyne.Resource
	TitleForEmptyState() string
	DescriptionForEmptyState() string
	ButtonTitleForEmptyState() string
}

// NewEmptyDataSet creates a data-driven empty view
func NewEmptyDataSet(dataSource EmptyDataSetDataSource) *EmptyDataSet {
	ev := NewEmptyState()
	edsv := &EmptyDataSet{
		EmptyState:  ev,
		DataSource: dataSource,
	}
	edsv.reloadData()
	return edsv
}

func (edsv *EmptyDataSet) reloadData() {
	if edsv.DataSource == nil {
		return
	}

	edsv.Image = edsv.DataSource.ImageForEmptyState()
	edsv.Text = edsv.DataSource.TitleForEmptyState()
	edsv.DetailText = edsv.DataSource.DescriptionForEmptyState()
	edsv.ActionText = edsv.DataSource.ButtonTitleForEmptyState()
	edsv.Refresh()
}

// ReloadData refreshes the empty view from the data source
func (edsv *EmptyDataSet) ReloadData() {
	edsv.reloadData()
}

// Preset empty views

// NoDataEmptyState creates an empty view for no data state
func NoDataEmptyState() *EmptyState {
	return NewEmptyStateWithTextAndDetail(
		"No Data",
		"There's nothing here yet",
	)
}

// NoNetworkEmptyState creates an empty view for no network state
func NoNetworkEmptyState(onRetry func()) *EmptyState {
	ev := NewEmptyStateWithTextAndDetail(
		"No Network Connection",
		"Please check your internet connection and try again",
	)
	ev.ActionText = "Retry"
	ev.OnActionTapped = onRetry
	return ev
}

// LoadingEmptyState creates an empty view with loading state
func LoadingEmptyState(text string) *EmptyState {
	ev := NewEmptyStateWithText(text)
	ev.SetLoading(true)
	return ev
}

// ErrorEmptyState creates an empty view for error state
func ErrorEmptyState(errorMessage string, onRetry func()) *EmptyState {
	ev := NewEmptyStateWithTextAndDetail(
		"Error",
		errorMessage,
	)
	ev.ActionText = "Retry"
	ev.OnActionTapped = onRetry
	return ev
}

// SearchEmptyState creates an empty view for no search results
func SearchEmptyState(query string) *EmptyState {
	ev := NewEmptyStateWithTextAndDetail(
		"No Results",
		"No results found for \""+query+"\"",
	)
	return ev
}
