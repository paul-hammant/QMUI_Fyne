// Package search provides QMUISearchBar - an enhanced search bar
// Ported from Tencent's QMUI_iOS framework
package search

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// SearchBarDelegate provides callbacks for search bar events
type SearchBarDelegate interface {
	SearchBarTextDidChange(searchBar *SearchBar, text string)
	SearchBarShouldBeginEditing(searchBar *SearchBar) bool
	SearchBarDidBeginEditing(searchBar *SearchBar)
	SearchBarShouldEndEditing(searchBar *SearchBar) bool
	SearchBarDidEndEditing(searchBar *SearchBar)
	SearchBarSearchButtonClicked(searchBar *SearchBar)
	SearchBarCancelButtonClicked(searchBar *SearchBar)
}

// SearchBar is an enhanced search input with styling options
type SearchBar struct {
	widget.BaseWidget

	// Text
	Text        string
	Placeholder string

	// Styling
	TextFieldBackgroundColor color.Color
	TextFieldBorderColor     color.Color
	BackgroundColor          color.Color
	TintColor                color.Color
	TextColor                color.Color
	PlaceholderColor         color.Color
	FontSize                 float32
	TextFieldCornerRadius    float32
	ContentInsets            core.EdgeInsets
	TextFieldInsets          core.EdgeInsets

	// Icons
	SearchIcon fyne.Resource
	ClearIcon  fyne.Resource

	// Buttons
	ShowsCancelButton  bool
	CancelButtonTitle  string
	CancelButtonColor  color.Color

	// Delegate
	Delegate SearchBarDelegate

	// Callbacks
	OnTextChanged       func(text string)
	OnSearchClicked     func()
	OnCancelClicked     func()
	OnFocusChanged      func(focused bool)

	// State
	mu       sync.RWMutex
	focused  bool
	entry    *widget.Entry
}

// NewSearchBar creates a new search bar
func NewSearchBar() *SearchBar {
	config := core.SharedConfiguration()
	sb := &SearchBar{
		Placeholder:              "Search",
		TextFieldBackgroundColor: config.SearchBarTextFieldBackgroundColor,
		TextFieldBorderColor:     config.SearchBarTextFieldBorderColor,
		BackgroundColor:          config.SearchBarBackgroundColor,
		TintColor:                config.SearchBarTintColor,
		TextColor:                config.SearchBarTextColor,
		PlaceholderColor:         config.SearchBarPlaceholderColor,
		FontSize:                 config.SearchBarFontSize,
		TextFieldCornerRadius:    config.SearchBarTextFieldCornerRadius,
		ContentInsets:            core.NewEdgeInsets(8, 8, 8, 8),
		TextFieldInsets:          core.NewEdgeInsets(8, 36, 8, 8),
		ShowsCancelButton:        false,
		CancelButtonTitle:        "Cancel",
		CancelButtonColor:        config.BlueColor,
	}
	sb.ExtendBaseWidget(sb)
	return sb
}

// NewSearchBarWithPlaceholder creates a search bar with placeholder text
func NewSearchBarWithPlaceholder(placeholder string) *SearchBar {
	sb := NewSearchBar()
	sb.Placeholder = placeholder
	return sb
}

// SetText sets the search text
func (sb *SearchBar) SetText(text string) {
	sb.mu.Lock()
	sb.Text = text
	sb.mu.Unlock()
	if sb.entry != nil {
		sb.entry.SetText(text)
	}
	sb.Refresh()
}

// GetText returns the current search text
func (sb *SearchBar) GetText() string {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.Text
}

// SetShowsCancelButton sets whether the cancel button is visible
func (sb *SearchBar) SetShowsCancelButton(show bool) {
	sb.mu.Lock()
	sb.ShowsCancelButton = show
	sb.mu.Unlock()
	sb.Refresh()
}

// Focus focuses the search bar
func (sb *SearchBar) Focus() {
	if sb.entry != nil {
		fyne.CurrentApp().Driver().CanvasForObject(sb).Focus(sb.entry)
	}
}

// Unfocus removes focus from the search bar
func (sb *SearchBar) Unfocus() {
	if sb.entry != nil {
		fyne.CurrentApp().Driver().CanvasForObject(sb).Unfocus()
	}
}

// CreateRenderer implements fyne.Widget
func (sb *SearchBar) CreateRenderer() fyne.WidgetRenderer {
	sb.ExtendBaseWidget(sb)

	background := canvas.NewRectangle(sb.BackgroundColor)

	textFieldBg := canvas.NewRectangle(sb.TextFieldBackgroundColor)
	textFieldBg.CornerRadius = sb.TextFieldCornerRadius
	textFieldBg.StrokeWidth = 0.5
	textFieldBg.StrokeColor = sb.TextFieldBorderColor

	// Search icon
	searchIcon := canvas.NewCircle(color.RGBA{R: 128, G: 128, B: 128, A: 255})
	searchIcon.StrokeWidth = 1.5
	searchIcon.StrokeColor = sb.PlaceholderColor
	searchIcon.FillColor = color.Transparent

	// Text entry
	entry := widget.NewEntry()
	entry.PlaceHolder = sb.Placeholder
	entry.OnChanged = func(text string) {
		sb.mu.Lock()
		sb.Text = text
		sb.mu.Unlock()
		if sb.OnTextChanged != nil {
			sb.OnTextChanged(text)
		}
		if sb.Delegate != nil {
			sb.Delegate.SearchBarTextDidChange(sb, text)
		}
	}
	entry.OnSubmitted = func(text string) {
		if sb.OnSearchClicked != nil {
			sb.OnSearchClicked()
		}
		if sb.Delegate != nil {
			sb.Delegate.SearchBarSearchButtonClicked(sb)
		}
	}
	sb.entry = entry

	// Cancel button
	cancelBtn := widget.NewButton(sb.CancelButtonTitle, func() {
		if sb.OnCancelClicked != nil {
			sb.OnCancelClicked()
		}
		if sb.Delegate != nil {
			sb.Delegate.SearchBarCancelButtonClicked(sb)
		}
	})

	return &searchBarRenderer{
		searchBar:   sb,
		background:  background,
		textFieldBg: textFieldBg,
		searchIcon:  searchIcon,
		entry:       entry,
		cancelBtn:   cancelBtn,
	}
}

type searchBarRenderer struct {
	searchBar   *SearchBar
	background  *canvas.Rectangle
	textFieldBg *canvas.Rectangle
	searchIcon  *canvas.Circle
	entry       *widget.Entry
	cancelBtn   *widget.Button
}

func (r *searchBarRenderer) Destroy() {}

func (r *searchBarRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	insets := r.searchBar.ContentInsets
	availableWidth := size.Width - insets.Left - insets.Right
	textFieldHeight := size.Height - insets.Top - insets.Bottom

	r.searchBar.mu.RLock()
	showCancel := r.searchBar.ShowsCancelButton
	r.searchBar.mu.RUnlock()

	var cancelWidth float32
	if showCancel {
		cancelWidth = r.cancelBtn.MinSize().Width + 8
		r.cancelBtn.Resize(fyne.NewSize(cancelWidth-8, textFieldHeight))
		r.cancelBtn.Move(fyne.NewPos(size.Width-insets.Right-cancelWidth+8, insets.Top))
		r.cancelBtn.Show()
	} else {
		r.cancelBtn.Hide()
	}

	textFieldWidth := availableWidth - cancelWidth
	r.textFieldBg.Resize(fyne.NewSize(textFieldWidth, textFieldHeight))
	r.textFieldBg.Move(fyne.NewPos(insets.Left, insets.Top))

	// Search icon
	iconSize := float32(14)
	r.searchIcon.Resize(fyne.NewSize(iconSize, iconSize))
	r.searchIcon.Move(fyne.NewPos(insets.Left+12, insets.Top+(textFieldHeight-iconSize)/2))

	// Entry
	entryInsets := r.searchBar.TextFieldInsets
	r.entry.Resize(fyne.NewSize(
		textFieldWidth-entryInsets.Left-entryInsets.Right,
		textFieldHeight-entryInsets.Top-entryInsets.Bottom,
	))
	r.entry.Move(fyne.NewPos(insets.Left+entryInsets.Left, insets.Top+entryInsets.Top))
}

func (r *searchBarRenderer) MinSize() fyne.Size {
	insets := r.searchBar.ContentInsets
	entryMin := r.entry.MinSize()
	return fyne.NewSize(
		200+insets.Left+insets.Right,
		entryMin.Height+insets.Top+insets.Bottom+8,
	)
}

func (r *searchBarRenderer) Refresh() {
	r.background.FillColor = r.searchBar.BackgroundColor
	r.textFieldBg.FillColor = r.searchBar.TextFieldBackgroundColor
	r.textFieldBg.StrokeColor = r.searchBar.TextFieldBorderColor
	r.textFieldBg.CornerRadius = r.searchBar.TextFieldCornerRadius

	r.entry.PlaceHolder = r.searchBar.Placeholder

	r.searchBar.mu.RLock()
	text := r.searchBar.Text
	r.searchBar.mu.RUnlock()

	if r.entry.Text != text {
		r.entry.SetText(text)
	}

	r.cancelBtn.SetText(r.searchBar.CancelButtonTitle)

	r.background.Refresh()
	r.textFieldBg.Refresh()
	r.searchIcon.Refresh()
	r.entry.Refresh()
	r.cancelBtn.Refresh()
}

func (r *searchBarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.background,
		r.textFieldBg,
		r.searchIcon,
		r.entry,
		r.cancelBtn,
	}
}

// SearchController manages search presentation
type SearchController struct {
	SearchBar         *SearchBar
	SearchResultsView fyne.CanvasObject

	// Behavior
	HidesNavigationBarDuringPresentation bool
	ObscuresBackgroundDuringPresentation bool

	// State
	mu     sync.RWMutex
	active bool
}

// NewSearchController creates a new search controller
func NewSearchController() *SearchController {
	return &SearchController{
		SearchBar:                            NewSearchBar(),
		HidesNavigationBarDuringPresentation: true,
		ObscuresBackgroundDuringPresentation: true,
	}
}

// SetActive sets whether the search is active
func (sc *SearchController) SetActive(active bool) {
	sc.mu.Lock()
	sc.active = active
	sc.mu.Unlock()

	sc.SearchBar.SetShowsCancelButton(active)
	if active {
		sc.SearchBar.Focus()
	} else {
		sc.SearchBar.Unfocus()
		sc.SearchBar.SetText("")
	}
}

// IsActive returns whether search is active
func (sc *SearchController) IsActive() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.active
}

// SearchSuggestion represents a search suggestion
type SearchSuggestion struct {
	Title    string
	Subtitle string
	Icon     fyne.Resource
}

// SearchSuggestionsView displays search suggestions
type SearchSuggestionsView struct {
	widget.BaseWidget

	Suggestions []SearchSuggestion

	// Styling
	BackgroundColor      color.Color
	SuggestionHeight     float32
	TitleColor           color.Color
	SubtitleColor        color.Color
	HighlightedColor     color.Color
	SeparatorColor       color.Color

	// Callbacks
	OnSuggestionSelected func(index int, suggestion SearchSuggestion)
}

// NewSearchSuggestionsView creates a search suggestions view
func NewSearchSuggestionsView() *SearchSuggestionsView {
	config := core.SharedConfiguration()
	return &SearchSuggestionsView{
		Suggestions:      make([]SearchSuggestion, 0),
		BackgroundColor:  color.White,
		SuggestionHeight: 44,
		TitleColor:       color.Black,
		SubtitleColor:    config.GrayColor,
		HighlightedColor: config.TableViewCellSelectedBackgroundColor,
		SeparatorColor:   config.SeparatorColor,
	}
}

// SetSuggestions updates the suggestions
func (sv *SearchSuggestionsView) SetSuggestions(suggestions []SearchSuggestion) {
	sv.Suggestions = suggestions
	sv.Refresh()
}

// CreateRenderer implements fyne.Widget
func (sv *SearchSuggestionsView) CreateRenderer() fyne.WidgetRenderer {
	sv.ExtendBaseWidget(sv)
	return &suggestionsRenderer{view: sv}
}

type suggestionsRenderer struct {
	view    *SearchSuggestionsView
	objects []fyne.CanvasObject
}

func (r *suggestionsRenderer) Destroy() {}

func (r *suggestionsRenderer) buildObjects() {
	r.objects = nil

	background := canvas.NewRectangle(r.view.BackgroundColor)
	r.objects = append(r.objects, background)

	for i, sug := range r.view.Suggestions {
		item := &suggestionItem{
			view:       r.view,
			index:      i,
			suggestion: sug,
		}
		item.ExtendBaseWidget(item)
		r.objects = append(r.objects, item)
	}
}

func (r *suggestionsRenderer) Layout(size fyne.Size) {
	if len(r.objects) == 0 {
		return
	}

	r.objects[0].Resize(size)

	y := float32(0)
	for i := 1; i < len(r.objects); i++ {
		r.objects[i].Resize(fyne.NewSize(size.Width, r.view.SuggestionHeight))
		r.objects[i].Move(fyne.NewPos(0, y))
		y += r.view.SuggestionHeight
	}
}

func (r *suggestionsRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, r.view.SuggestionHeight*float32(len(r.view.Suggestions)))
}

func (r *suggestionsRenderer) Refresh() {
	r.buildObjects()
}

func (r *suggestionsRenderer) Objects() []fyne.CanvasObject {
	r.buildObjects()
	return r.objects
}

type suggestionItem struct {
	widget.BaseWidget

	view       *SearchSuggestionsView
	index      int
	suggestion SearchSuggestion
	hovered    bool
	mu         sync.RWMutex
}

func (si *suggestionItem) CreateRenderer() fyne.WidgetRenderer {
	si.ExtendBaseWidget(si)

	bg := canvas.NewRectangle(color.Transparent)
	title := canvas.NewText(si.suggestion.Title, si.view.TitleColor)
	title.TextSize = theme.TextSize()

	var subtitle *canvas.Text
	if si.suggestion.Subtitle != "" {
		subtitle = canvas.NewText(si.suggestion.Subtitle, si.view.SubtitleColor)
		subtitle.TextSize = theme.TextSize() - 2
	}

	return &suggestionItemRenderer{
		item:     si,
		bg:       bg,
		title:    title,
		subtitle: subtitle,
	}
}

func (si *suggestionItem) Tapped(_ *fyne.PointEvent) {
	if si.view.OnSuggestionSelected != nil {
		si.view.OnSuggestionSelected(si.index, si.suggestion)
	}
}

func (si *suggestionItem) TappedSecondary(_ *fyne.PointEvent) {}

func (si *suggestionItem) MouseIn(_ *desktop.MouseEvent) {
	si.mu.Lock()
	si.hovered = true
	si.mu.Unlock()
	si.Refresh()
}

func (si *suggestionItem) MouseMoved(_ *desktop.MouseEvent) {}

func (si *suggestionItem) MouseOut() {
	si.mu.Lock()
	si.hovered = false
	si.mu.Unlock()
	si.Refresh()
}

func (si *suggestionItem) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type suggestionItemRenderer struct {
	item     *suggestionItem
	bg       *canvas.Rectangle
	title    *canvas.Text
	subtitle *canvas.Text
}

func (r *suggestionItemRenderer) Destroy() {}

func (r *suggestionItemRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	titleSize := r.title.MinSize()
	if r.subtitle != nil {
		subtitleSize := r.subtitle.MinSize()
		totalHeight := titleSize.Height + subtitleSize.Height
		startY := (size.Height - totalHeight) / 2
		r.title.Move(fyne.NewPos(16, startY))
		r.subtitle.Move(fyne.NewPos(16, startY+titleSize.Height))
	} else {
		r.title.Move(fyne.NewPos(16, (size.Height-titleSize.Height)/2))
	}
}

func (r *suggestionItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, r.item.view.SuggestionHeight)
}

func (r *suggestionItemRenderer) Refresh() {
	r.item.mu.RLock()
	hovered := r.item.hovered
	r.item.mu.RUnlock()

	if hovered {
		r.bg.FillColor = r.item.view.HighlightedColor
	} else {
		r.bg.FillColor = color.Transparent
	}
	r.bg.Refresh()
	r.title.Refresh()
	if r.subtitle != nil {
		r.subtitle.Refresh()
	}
}

func (r *suggestionItemRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bg, r.title}
	if r.subtitle != nil {
		objects = append(objects, r.subtitle)
	}
	return objects
}
