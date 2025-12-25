// Package navigation provides navigation bars, title views, and tab bars
package navigation

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// TitleViewStyle defines the style of the title view
type TitleViewStyle int

const (
	// TitleViewStyleDefault shows title and subtitle vertically
	TitleViewStyleDefault TitleViewStyle = iota
	// TitleViewStyleSubtitleIsDetail shows subtitle as detail text
	TitleViewStyleSubtitleIsDetail
)

// TitleView is a customizable navigation title view
type TitleView struct {
	widget.BaseWidget

	// Content
	Title           string
	Subtitle        string
	Style           TitleViewStyle

	// Styling
	TitleColor           color.Color
	SubtitleColor        color.Color
	TitleFontSize        float32
	SubtitleFontSize     float32
	TitleTextStyle       fyne.TextStyle
	SubtitleTextStyle    fyne.TextStyle
	HorizontalAlignment  fyne.TextAlign
	SpacingBetweenTitleAndSubtitle float32

	// Accessory
	AccessoryView        fyne.CanvasObject
	AccessoryType        AccessoryType
	AccessorySpacing     float32

	// Behavior
	NeedsLoadingView     bool
	UserInteractionEnabled bool

	// Callbacks
	OnTapped func()

	mu      sync.RWMutex
	loading bool
}

// AccessoryType defines the type of accessory view
type AccessoryType int

const (
	// AccessoryTypeNone shows no accessory
	AccessoryTypeNone AccessoryType = iota
	// AccessoryTypeDisclosureIndicator shows a disclosure indicator
	AccessoryTypeDisclosureIndicator
)

// NewTitleView creates a new navigation title view
func NewTitleView() *TitleView {
	config := core.SharedConfiguration()
	ntv := &TitleView{
		Style:                TitleViewStyleDefault,
		TitleColor:           config.NavBarTitleColor,
		SubtitleColor:        config.GrayColor,
		TitleFontSize:        config.NavBarTitleFontSize,
		SubtitleFontSize:     12,
		TitleTextStyle:       fyne.TextStyle{Bold: true},
		SubtitleTextStyle:    fyne.TextStyle{},
		HorizontalAlignment:  fyne.TextAlignCenter,
		SpacingBetweenTitleAndSubtitle: 0,
		AccessoryType:        AccessoryTypeNone,
		AccessorySpacing:     4,
		NeedsLoadingView:     false,
		UserInteractionEnabled: true,
	}
	ntv.ExtendBaseWidget(ntv)
	return ntv
}

// NewTitleViewWithTitle creates a title view with title
func NewTitleViewWithTitle(title string) *TitleView {
	ntv := NewTitleView()
	ntv.Title = title
	return ntv
}

// NewTitleViewWithTitleAndSubtitle creates a title view with title and subtitle
func NewTitleViewWithTitleAndSubtitle(title, subtitle string) *TitleView {
	ntv := NewTitleView()
	ntv.Title = title
	ntv.Subtitle = subtitle
	return ntv
}

// SetTitle sets the title
func (ntv *TitleView) SetTitle(title string) {
	ntv.mu.Lock()
	ntv.Title = title
	ntv.mu.Unlock()
	ntv.Refresh()
}

// SetSubtitle sets the subtitle
func (ntv *TitleView) SetSubtitle(subtitle string) {
	ntv.mu.Lock()
	ntv.Subtitle = subtitle
	ntv.mu.Unlock()
	ntv.Refresh()
}

// SetLoading sets the loading state
func (ntv *TitleView) SetLoading(loading bool) {
	ntv.mu.Lock()
	ntv.loading = loading
	ntv.mu.Unlock()
	ntv.Refresh()
}

// CreateRenderer implements fyne.Widget
func (ntv *TitleView) CreateRenderer() fyne.WidgetRenderer {
	ntv.ExtendBaseWidget(ntv)

	titleLabel := canvas.NewText(ntv.Title, ntv.TitleColor)
	titleLabel.TextStyle = ntv.TitleTextStyle
	titleLabel.TextSize = ntv.TitleFontSize
	titleLabel.Alignment = ntv.HorizontalAlignment

	subtitleLabel := canvas.NewText(ntv.Subtitle, ntv.SubtitleColor)
	subtitleLabel.TextStyle = ntv.SubtitleTextStyle
	subtitleLabel.TextSize = ntv.SubtitleFontSize
	subtitleLabel.Alignment = ntv.HorizontalAlignment

	loading := widget.NewProgressBarInfinite()
	loading.Hide()

	return &navigationTitleRenderer{
		titleView:     ntv,
		titleLabel:    titleLabel,
		subtitleLabel: subtitleLabel,
		loading:       loading,
	}
}

// Tapped handles tap events
func (ntv *TitleView) Tapped(_ *fyne.PointEvent) {
	if !ntv.UserInteractionEnabled {
		return
	}
	if ntv.OnTapped != nil {
		ntv.OnTapped()
	}
}

// TappedSecondary handles secondary tap
func (ntv *TitleView) TappedSecondary(_ *fyne.PointEvent) {}

func (ntv *TitleView) Cursor() desktop.Cursor {
	if ntv.UserInteractionEnabled && ntv.OnTapped != nil {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

type navigationTitleRenderer struct {
	titleView     *TitleView
	titleLabel    *canvas.Text
	subtitleLabel *canvas.Text
	loading       *widget.ProgressBarInfinite
}

func (r *navigationTitleRenderer) Destroy() {}

func (r *navigationTitleRenderer) Layout(size fyne.Size) {
	titleSize := r.titleLabel.MinSize()
	subtitleSize := r.subtitleLabel.MinSize()

	r.titleView.mu.RLock()
	hasSubtitle := r.titleView.Subtitle != ""
	spacing := r.titleView.SpacingBetweenTitleAndSubtitle
	r.titleView.mu.RUnlock()

	if hasSubtitle {
		totalHeight := titleSize.Height + subtitleSize.Height + spacing
		startY := (size.Height - totalHeight) / 2
		r.titleLabel.Move(fyne.NewPos(0, startY))
		r.titleLabel.Resize(fyne.NewSize(size.Width, titleSize.Height))
		r.subtitleLabel.Move(fyne.NewPos(0, startY+titleSize.Height+spacing))
		r.subtitleLabel.Resize(fyne.NewSize(size.Width, subtitleSize.Height))
		r.subtitleLabel.Show()
	} else {
		r.titleLabel.Move(fyne.NewPos(0, (size.Height-titleSize.Height)/2))
		r.titleLabel.Resize(fyne.NewSize(size.Width, titleSize.Height))
		r.subtitleLabel.Hide()
	}
}

func (r *navigationTitleRenderer) MinSize() fyne.Size {
	titleSize := r.titleLabel.MinSize()

	r.titleView.mu.RLock()
	hasSubtitle := r.titleView.Subtitle != ""
	spacing := r.titleView.SpacingBetweenTitleAndSubtitle
	r.titleView.mu.RUnlock()

	width := titleSize.Width
	height := titleSize.Height

	if hasSubtitle {
		subtitleSize := r.subtitleLabel.MinSize()
		if subtitleSize.Width > width {
			width = subtitleSize.Width
		}
		height += subtitleSize.Height + spacing
	}

	return fyne.NewSize(width, height)
}

func (r *navigationTitleRenderer) Refresh() {
	r.titleView.mu.RLock()
	title := r.titleView.Title
	subtitle := r.titleView.Subtitle
	loading := r.titleView.loading
	r.titleView.mu.RUnlock()

	r.titleLabel.Text = title
	r.titleLabel.Color = r.titleView.TitleColor
	r.titleLabel.TextStyle = r.titleView.TitleTextStyle
	r.titleLabel.TextSize = r.titleView.TitleFontSize
	r.titleLabel.Alignment = r.titleView.HorizontalAlignment

	r.subtitleLabel.Text = subtitle
	r.subtitleLabel.Color = r.titleView.SubtitleColor
	r.subtitleLabel.TextStyle = r.titleView.SubtitleTextStyle
	r.subtitleLabel.TextSize = r.titleView.SubtitleFontSize
	r.subtitleLabel.Alignment = r.titleView.HorizontalAlignment

	if loading && r.titleView.NeedsLoadingView {
		r.loading.Show()
	} else {
		r.loading.Hide()
	}

	r.titleLabel.Refresh()
	r.subtitleLabel.Refresh()
}

func (r *navigationTitleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.titleLabel, r.subtitleLabel, r.loading}
}

// NavigationBar is a navigation bar component
type NavigationBar struct {
	widget.BaseWidget

	// Content
	TitleView        fyne.CanvasObject
	LeftBarItems     []fyne.CanvasObject
	RightBarItems    []fyne.CanvasObject

	// Styling
	BackgroundColor  color.Color
	TintColor        color.Color
	ShadowColor      color.Color
	ShadowEnabled    bool
	Height           float32

	mu sync.RWMutex
}

// NewNavigationBar creates a new navigation bar
func NewNavigationBar() *NavigationBar {
	config := core.SharedConfiguration()
	nb := &NavigationBar{
		BackgroundColor: config.NavBarBackgroundColor,
		TintColor:       config.NavBarTintColor,
		ShadowColor:     config.NavBarShadowColor,
		ShadowEnabled:   true,
		Height:          44,
	}
	nb.ExtendBaseWidget(nb)
	return nb
}

// SetTitleView sets the title view
func (nb *NavigationBar) SetTitleView(view fyne.CanvasObject) {
	nb.mu.Lock()
	nb.TitleView = view
	nb.mu.Unlock()
	nb.Refresh()
}

// SetLeftBarItems sets the left bar items
func (nb *NavigationBar) SetLeftBarItems(items []fyne.CanvasObject) {
	nb.mu.Lock()
	nb.LeftBarItems = items
	nb.mu.Unlock()
	nb.Refresh()
}

// SetRightBarItems sets the right bar items
func (nb *NavigationBar) SetRightBarItems(items []fyne.CanvasObject) {
	nb.mu.Lock()
	nb.RightBarItems = items
	nb.mu.Unlock()
	nb.Refresh()
}

// CreateRenderer implements fyne.Widget
func (nb *NavigationBar) CreateRenderer() fyne.WidgetRenderer {
	nb.ExtendBaseWidget(nb)

	background := canvas.NewRectangle(nb.BackgroundColor)
	shadow := canvas.NewRectangle(nb.ShadowColor)

	return &navigationBarRenderer{
		bar:        nb,
		background: background,
		shadow:     shadow,
	}
}

type navigationBarRenderer struct {
	bar        *NavigationBar
	background *canvas.Rectangle
	shadow     *canvas.Rectangle
}

func (r *navigationBarRenderer) Destroy() {}

func (r *navigationBarRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	// Shadow at bottom
	r.shadow.Resize(fyne.NewSize(size.Width, 0.5))
	r.shadow.Move(fyne.NewPos(0, size.Height-0.5))

	r.bar.mu.RLock()
	leftItems := r.bar.LeftBarItems
	rightItems := r.bar.RightBarItems
	titleView := r.bar.TitleView
	r.bar.mu.RUnlock()

	// Layout left items
	leftX := float32(8)
	for _, item := range leftItems {
		itemSize := item.MinSize()
		item.Resize(itemSize)
		item.Move(fyne.NewPos(leftX, (size.Height-itemSize.Height)/2))
		leftX += itemSize.Width + 8
	}

	// Layout right items
	rightX := size.Width - 8
	for i := len(rightItems) - 1; i >= 0; i-- {
		item := rightItems[i]
		itemSize := item.MinSize()
		rightX -= itemSize.Width
		item.Resize(itemSize)
		item.Move(fyne.NewPos(rightX, (size.Height-itemSize.Height)/2))
		rightX -= 8
	}

	// Layout title view
	if titleView != nil {
		titleWidth := rightX - leftX - 16
		if titleWidth < 100 {
			titleWidth = 100
		}
		titleSize := titleView.MinSize()
		titleView.Resize(fyne.NewSize(titleWidth, titleSize.Height))
		titleView.Move(fyne.NewPos(leftX+8, (size.Height-titleSize.Height)/2))
	}
}

func (r *navigationBarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, r.bar.Height)
}

func (r *navigationBarRenderer) Refresh() {
	r.background.FillColor = r.bar.BackgroundColor

	if r.bar.ShadowEnabled {
		r.shadow.FillColor = r.bar.ShadowColor
		r.shadow.Show()
	} else {
		r.shadow.Hide()
	}

	r.background.Refresh()
	r.shadow.Refresh()
}

func (r *navigationBarRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background, r.shadow}

	r.bar.mu.RLock()
	defer r.bar.mu.RUnlock()

	for _, item := range r.bar.LeftBarItems {
		objects = append(objects, item)
	}
	for _, item := range r.bar.RightBarItems {
		objects = append(objects, item)
	}
	if r.bar.TitleView != nil {
		objects = append(objects, r.bar.TitleView)
	}

	return objects
}

// TabBar is a tab bar component
type TabBar struct {
	widget.BaseWidget

	// Items
	Items         []*TabBarItem
	SelectedIndex int

	// Styling
	BackgroundColor         color.Color
	TintColor               color.Color
	UnselectedItemColor     color.Color
	SelectedItemColor       color.Color
	ShadowColor             color.Color
	ShadowEnabled           bool
	Height                  float32
	ItemTitleFontSize       float32
	ItemTitleFontSizeSelected float32

	// Callbacks
	OnItemSelected func(index int)

	mu sync.RWMutex
}

// TabBarItem represents an item in the tab bar
type TabBarItem struct {
	Title       string
	Icon        fyne.Resource
	SelectedIcon fyne.Resource
	BadgeValue  string
}

// NewTabBarItem creates a new tab bar item
func NewTabBarItem(title string, icon fyne.Resource) *TabBarItem {
	return &TabBarItem{
		Title: title,
		Icon:  icon,
	}
}

// NewTabBar creates a new tab bar
func NewTabBar(items []*TabBarItem) *TabBar {
	config := core.SharedConfiguration()
	tb := &TabBar{
		Items:                   items,
		SelectedIndex:           0,
		BackgroundColor:         config.TabBarBackgroundColor,
		TintColor:               config.BlueColor,
		UnselectedItemColor:     config.TabBarItemTitleColor,
		SelectedItemColor:       config.TabBarItemTitleColorSelected,
		ShadowColor:             config.TabBarShadowColor,
		ShadowEnabled:           true,
		Height:                  49,
		ItemTitleFontSize:       config.TabBarItemTitleFontSize,
		ItemTitleFontSizeSelected: config.TabBarItemTitleFontSizeSelected,
	}
	tb.ExtendBaseWidget(tb)
	return tb
}

// SetSelectedIndex sets the selected tab index
func (tb *TabBar) SetSelectedIndex(index int) {
	if index < 0 || index >= len(tb.Items) {
		return
	}
	tb.mu.Lock()
	tb.SelectedIndex = index
	tb.mu.Unlock()
	tb.Refresh()
	if tb.OnItemSelected != nil {
		tb.OnItemSelected(index)
	}
}

// CreateRenderer implements fyne.Widget
func (tb *TabBar) CreateRenderer() fyne.WidgetRenderer {
	tb.ExtendBaseWidget(tb)

	background := canvas.NewRectangle(tb.BackgroundColor)
	shadow := canvas.NewRectangle(tb.ShadowColor)

	return &tabBarRenderer{
		tabBar:     tb,
		background: background,
		shadow:     shadow,
		items:      make([]*tabBarItemWidget, 0),
	}
}

type tabBarRenderer struct {
	tabBar     *TabBar
	background *canvas.Rectangle
	shadow     *canvas.Rectangle
	items      []*tabBarItemWidget
}

func (r *tabBarRenderer) Destroy() {}

func (r *tabBarRenderer) updateItems() {
	if len(r.items) != len(r.tabBar.Items) {
		r.items = make([]*tabBarItemWidget, len(r.tabBar.Items))
		for i, item := range r.tabBar.Items {
			w := &tabBarItemWidget{
				tabBar: r.tabBar,
				index:  i,
				item:   item,
			}
			w.ExtendBaseWidget(w)
			r.items[i] = w
		}
	}
}

func (r *tabBarRenderer) Layout(size fyne.Size) {
	r.updateItems()

	r.background.Resize(size)

	// Shadow at top
	r.shadow.Resize(fyne.NewSize(size.Width, 0.5))
	r.shadow.Move(fyne.NewPos(0, 0))

	if len(r.items) == 0 {
		return
	}

	itemWidth := size.Width / float32(len(r.items))
	for i, item := range r.items {
		item.Resize(fyne.NewSize(itemWidth, size.Height))
		item.Move(fyne.NewPos(float32(i)*itemWidth, 0))
	}
}

func (r *tabBarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, r.tabBar.Height)
}

func (r *tabBarRenderer) Refresh() {
	r.updateItems()
	r.background.FillColor = r.tabBar.BackgroundColor

	if r.tabBar.ShadowEnabled {
		r.shadow.FillColor = r.tabBar.ShadowColor
		r.shadow.Show()
	} else {
		r.shadow.Hide()
	}

	r.background.Refresh()
	r.shadow.Refresh()

	for _, item := range r.items {
		item.Refresh()
	}
}

func (r *tabBarRenderer) Objects() []fyne.CanvasObject {
	r.updateItems()
	objects := []fyne.CanvasObject{r.background, r.shadow}
	for _, item := range r.items {
		objects = append(objects, item)
	}
	return objects
}

type tabBarItemWidget struct {
	widget.BaseWidget

	tabBar *TabBar
	index  int
	item   *TabBarItem
}

func (w *tabBarItemWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	var icon *canvas.Image
	if w.item.Icon != nil {
		icon = canvas.NewImageFromResource(w.item.Icon)
		icon.FillMode = canvas.ImageFillContain
		icon.SetMinSize(fyne.NewSize(24, 24))
	}

	title := canvas.NewText(w.item.Title, w.tabBar.UnselectedItemColor)
	title.TextSize = w.tabBar.ItemTitleFontSize
	title.Alignment = fyne.TextAlignCenter

	return &tabBarItemRenderer{
		widget: w,
		icon:   icon,
		title:  title,
	}
}

func (w *tabBarItemWidget) Tapped(_ *fyne.PointEvent) {
	w.tabBar.SetSelectedIndex(w.index)
}

func (w *tabBarItemWidget) TappedSecondary(_ *fyne.PointEvent) {}

func (w *tabBarItemWidget) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type tabBarItemRenderer struct {
	widget *tabBarItemWidget
	icon   *canvas.Image
	title  *canvas.Text
}

func (r *tabBarItemRenderer) Destroy() {}

func (r *tabBarItemRenderer) Layout(size fyne.Size) {
	titleSize := r.title.MinSize()

	if r.icon != nil {
		iconSize := fyne.NewSize(24, 24)
		totalHeight := iconSize.Height + 4 + titleSize.Height
		startY := (size.Height - totalHeight) / 2

		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos((size.Width-iconSize.Width)/2, startY))

		r.title.Move(fyne.NewPos(0, startY+iconSize.Height+4))
		r.title.Resize(fyne.NewSize(size.Width, titleSize.Height))
	} else {
		r.title.Move(fyne.NewPos(0, (size.Height-titleSize.Height)/2))
		r.title.Resize(fyne.NewSize(size.Width, titleSize.Height))
	}
}

func (r *tabBarItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(60, r.widget.tabBar.Height)
}

func (r *tabBarItemRenderer) Refresh() {
	r.widget.tabBar.mu.RLock()
	selected := r.widget.index == r.widget.tabBar.SelectedIndex
	r.widget.tabBar.mu.RUnlock()

	if selected {
		r.title.Color = r.widget.tabBar.SelectedItemColor
		r.title.TextSize = r.widget.tabBar.ItemTitleFontSizeSelected
		if r.icon != nil && r.widget.item.SelectedIcon != nil {
			r.icon.Resource = r.widget.item.SelectedIcon
		}
	} else {
		r.title.Color = r.widget.tabBar.UnselectedItemColor
		r.title.TextSize = r.widget.tabBar.ItemTitleFontSize
		if r.icon != nil {
			r.icon.Resource = r.widget.item.Icon
		}
	}

	r.title.Refresh()
	if r.icon != nil {
		r.icon.Refresh()
	}
}

func (r *tabBarItemRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.title}
	if r.icon != nil {
		objects = append(objects, r.icon)
	}
	return objects
}
