// Package popup provides customizable popup menus with arrow indicators
package popup

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// ArrowDirection defines the popup arrow direction
type ArrowDirection int

const (
	// ArrowDirectionUp points up
	ArrowDirectionUp ArrowDirection = iota
	// ArrowDirectionDown points down
	ArrowDirectionDown
	// ArrowDirectionLeft points left
	ArrowDirectionLeft
	// ArrowDirectionRight points right
	ArrowDirectionRight
	// ArrowDirectionNone no arrow
	ArrowDirectionNone
)

// MenuItem represents an item in the popup menu
type MenuItem struct {
	Title    string
	Subtitle string
	Icon     fyne.Resource
	Enabled  bool
	Handler  func(item *MenuItem)

	// Styling
	TitleColor    color.Color
	SubtitleColor color.Color
	IconTintColor color.Color
}

// NewMenuItem creates a new menu item
func NewMenuItem(title string, handler func(item *MenuItem)) *MenuItem {
	return &MenuItem{
		Title:   title,
		Enabled: true,
		Handler: handler,
	}
}

// NewMenuItemWithIcon creates a menu item with icon
func NewMenuItemWithIcon(title string, icon fyne.Resource, handler func(item *MenuItem)) *MenuItem {
	item := NewMenuItem(title, handler)
	item.Icon = icon
	return item
}

// PopupContainer is a generic popup container with arrow
type PopupContainer struct {
	widget.BaseWidget

	// Styling
	BackgroundColor   color.Color
	BorderColor       color.Color
	BorderWidth       float32
	CornerRadius      float32
	ArrowSize         fyne.Size
	ArrowDirection    ArrowDirection
	ShadowEnabled     bool
	ShadowColor       color.Color
	ShadowOffset      core.Offset
	ShadowRadius      float32
	ContentEdgeInsets core.EdgeInsets
	MaximumWidth      float32
	MaximumHeight     float32

	// Content
	ContentView fyne.CanvasObject

	// State
	mu      sync.RWMutex
	popup   *widget.PopUp
	window  fyne.Window
	visible bool
}

// NewPopupContainer creates a new popup container
func NewPopupContainer() *PopupContainer {
	config := core.SharedConfiguration()
	pcv := &PopupContainer{
		BackgroundColor:   color.White,
		BorderColor:       config.SeparatorColor,
		BorderWidth:       0.5,
		CornerRadius:      8,
		ArrowSize:         fyne.NewSize(16, 8),
		ArrowDirection:    ArrowDirectionUp,
		ShadowEnabled:     true,
		ShadowColor:       color.RGBA{R: 0, G: 0, B: 0, A: 40},
		ShadowOffset:      core.NewOffset(0, 2),
		ShadowRadius:      8,
		ContentEdgeInsets: core.NewEdgeInsets(4, 0, 4, 0),
		MaximumWidth:      0,
		MaximumHeight:     0,
	}
	pcv.ExtendBaseWidget(pcv)
	return pcv
}

// Show displays the popup at the specified position
func (pcv *PopupContainer) ShowAt(window fyne.Window, position fyne.Position) {
	pcv.mu.Lock()
	pcv.window = window
	pcv.visible = true
	pcv.mu.Unlock()

	content := pcv.buildContent()
	pcv.popup = widget.NewPopUp(content, window.Canvas())
	pcv.popup.Move(position)
	pcv.popup.Show()
}

// ShowBelowView shows the popup below a view
func (pcv *PopupContainer) ShowBelowView(window fyne.Window, view fyne.CanvasObject) {
	viewPos := view.Position()
	viewSize := view.Size()
	position := fyne.NewPos(viewPos.X, viewPos.Y+viewSize.Height+pcv.ArrowSize.Height)
	pcv.ArrowDirection = ArrowDirectionUp
	pcv.ShowAt(window, position)
}

// ShowAboveView shows the popup above a view
func (pcv *PopupContainer) ShowAboveView(window fyne.Window, view fyne.CanvasObject) {
	viewPos := view.Position()
	contentSize := pcv.buildContent().MinSize()
	position := fyne.NewPos(viewPos.X, viewPos.Y-contentSize.Height-pcv.ArrowSize.Height)
	pcv.ArrowDirection = ArrowDirectionDown
	pcv.ShowAt(window, position)
}

// Hide hides the popup
func (pcv *PopupContainer) Hide() {
	pcv.mu.Lock()
	defer pcv.mu.Unlock()

	if pcv.popup != nil {
		pcv.popup.Hide()
		pcv.popup = nil
	}
	pcv.visible = false
}

// IsVisible returns whether the popup is visible
func (pcv *PopupContainer) IsVisible() bool {
	pcv.mu.RLock()
	defer pcv.mu.RUnlock()
	return pcv.visible
}

func (pcv *PopupContainer) buildContent() fyne.CanvasObject {
	background := canvas.NewRectangle(pcv.BackgroundColor)
	background.CornerRadius = pcv.CornerRadius
	background.StrokeWidth = pcv.BorderWidth
	background.StrokeColor = pcv.BorderColor

	var content fyne.CanvasObject = background
	if pcv.ContentView != nil {
		content = container.NewStack(background, container.NewPadded(pcv.ContentView))
	}

	return content
}

func (pcv *PopupContainer) CreateRenderer() fyne.WidgetRenderer {
	pcv.ExtendBaseWidget(pcv)
	return &popupContainerRenderer{container: pcv}
}

type popupContainerRenderer struct {
	container *PopupContainer
}

func (r *popupContainerRenderer) Destroy()                      {}
func (r *popupContainerRenderer) Layout(size fyne.Size)         {}
func (r *popupContainerRenderer) MinSize() fyne.Size            { return fyne.NewSize(0, 0) }
func (r *popupContainerRenderer) Refresh()                      {}
func (r *popupContainerRenderer) Objects() []fyne.CanvasObject { return nil }

// PopupMenu is a popup menu with multiple items
type PopupMenu struct {
	*PopupContainer

	// Items
	Items []*MenuItem

	// Styling
	ItemHeight             float32
	ItemPaddingHorizontal  float32
	SeparatorColor         color.Color
	SeparatorInsets        core.EdgeInsets
	TitleFontSize          float32
	SubtitleFontSize       float32
	TitleColor             color.Color
	SubtitleColor          color.Color
	HighlightedColor       color.Color
	IconSize               fyne.Size
	SpacingBetweenIconAndTitle float32

	// Behavior
	ShouldDismissAfterSelection bool

	// Callbacks
	OnItemSelected func(index int, item *MenuItem)
	OnDismiss      func()
}

// NewPopupMenu creates a new popup menu
func NewPopupMenu() *PopupMenu {
	config := core.SharedConfiguration()
	pmv := &PopupMenu{
		PopupContainer:          NewPopupContainer(),
		Items:                       make([]*MenuItem, 0),
		ItemHeight:                  44,
		ItemPaddingHorizontal:       16,
		SeparatorColor:              config.SeparatorColor,
		SeparatorInsets:             core.NewEdgeInsets(0, 16, 0, 16),
		TitleFontSize:               theme.TextSize(),
		SubtitleFontSize:            theme.TextSize() - 2,
		TitleColor:                  color.Black,
		SubtitleColor:               config.GrayColor,
		HighlightedColor:            color.RGBA{R: 0, G: 0, B: 0, A: 20},
		IconSize:                    fyne.NewSize(20, 20),
		SpacingBetweenIconAndTitle:  12,
		ShouldDismissAfterSelection: true,
	}
	return pmv
}

// NewPopupMenuWithItems creates a popup menu with items
func NewPopupMenuWithItems(items []*MenuItem) *PopupMenu {
	pmv := NewPopupMenu()
	pmv.Items = items
	pmv.buildMenuContent()
	return pmv
}

// AddItem adds an item to the menu
func (pmv *PopupMenu) AddItem(item *MenuItem) {
	pmv.Items = append(pmv.Items, item)
	pmv.buildMenuContent()
}

// RemoveItem removes an item from the menu
func (pmv *PopupMenu) RemoveItem(index int) {
	if index < 0 || index >= len(pmv.Items) {
		return
	}
	pmv.Items = append(pmv.Items[:index], pmv.Items[index+1:]...)
	pmv.buildMenuContent()
}

func (pmv *PopupMenu) buildMenuContent() {
	var objects []fyne.CanvasObject

	for i, item := range pmv.Items {
		itemView := pmv.createItemView(i, item)
		objects = append(objects, itemView)

		// Add separator (except after last item)
		if i < len(pmv.Items)-1 {
			sep := canvas.NewRectangle(pmv.SeparatorColor)
			sep.SetMinSize(fyne.NewSize(0, 0.5))
			objects = append(objects, sep)
		}
	}

	pmv.ContentView = container.NewVBox(objects...)
}

func (pmv *PopupMenu) createItemView(index int, item *MenuItem) fyne.CanvasObject {
	itemWidget := &menuItemWidget{
		menu:  pmv,
		index: index,
		item:  item,
	}
	itemWidget.ExtendBaseWidget(itemWidget)
	return itemWidget
}

// Show displays the popup menu
func (pmv *PopupMenu) Show(window fyne.Window, position fyne.Position) {
	pmv.buildMenuContent()
	pmv.PopupContainer.ShowAt(window, position)
}

// menuItemWidget represents a menu item in the popup
type menuItemWidget struct {
	widget.BaseWidget

	menu    *PopupMenu
	index   int
	item    *MenuItem
	hovered bool
	mu      sync.RWMutex
}

func (w *menuItemWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	background := canvas.NewRectangle(color.Transparent)

	var icon *canvas.Image
	if w.item.Icon != nil {
		icon = canvas.NewImageFromResource(w.item.Icon)
		icon.FillMode = canvas.ImageFillContain
		icon.SetMinSize(w.menu.IconSize)
	}

	titleColor := w.menu.TitleColor
	if w.item.TitleColor != nil {
		titleColor = w.item.TitleColor
	}

	title := canvas.NewText(w.item.Title, titleColor)
	title.TextSize = w.menu.TitleFontSize

	var subtitle *canvas.Text
	if w.item.Subtitle != "" {
		subtitleColor := w.menu.SubtitleColor
		if w.item.SubtitleColor != nil {
			subtitleColor = w.item.SubtitleColor
		}
		subtitle = canvas.NewText(w.item.Subtitle, subtitleColor)
		subtitle.TextSize = w.menu.SubtitleFontSize
	}

	return &menuItemRenderer{
		widget:     w,
		background: background,
		icon:       icon,
		title:      title,
		subtitle:   subtitle,
	}
}

func (w *menuItemWidget) Tapped(_ *fyne.PointEvent) {
	if !w.item.Enabled {
		return
	}

	if w.item.Handler != nil {
		w.item.Handler(w.item)
	}

	if w.menu.OnItemSelected != nil {
		w.menu.OnItemSelected(w.index, w.item)
	}

	if w.menu.ShouldDismissAfterSelection {
		w.menu.Hide()
		if w.menu.OnDismiss != nil {
			w.menu.OnDismiss()
		}
	}
}

func (w *menuItemWidget) TappedSecondary(_ *fyne.PointEvent) {}

func (w *menuItemWidget) MouseIn(_ *desktop.MouseEvent) {
	w.mu.Lock()
	w.hovered = true
	w.mu.Unlock()
	w.Refresh()
}

func (w *menuItemWidget) MouseMoved(_ *desktop.MouseEvent) {}

func (w *menuItemWidget) MouseOut() {
	w.mu.Lock()
	w.hovered = false
	w.mu.Unlock()
	w.Refresh()
}

func (w *menuItemWidget) Cursor() desktop.Cursor {
	if w.item.Enabled {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

type menuItemRenderer struct {
	widget     *menuItemWidget
	background *canvas.Rectangle
	icon       *canvas.Image
	title      *canvas.Text
	subtitle   *canvas.Text
}

func (r *menuItemRenderer) Destroy() {}

func (r *menuItemRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	padding := r.widget.menu.ItemPaddingHorizontal
	x := padding
	centerY := size.Height / 2

	if r.icon != nil {
		iconSize := r.widget.menu.IconSize
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos(x, centerY-iconSize.Height/2))
		x += iconSize.Width + r.widget.menu.SpacingBetweenIconAndTitle
	}

	titleSize := r.title.MinSize()
	if r.subtitle != nil {
		subtitleSize := r.subtitle.MinSize()
		totalHeight := titleSize.Height + subtitleSize.Height
		titleY := centerY - totalHeight/2
		r.title.Move(fyne.NewPos(x, titleY))
		r.subtitle.Move(fyne.NewPos(x, titleY+titleSize.Height))
	} else {
		r.title.Move(fyne.NewPos(x, centerY-titleSize.Height/2))
	}
}

func (r *menuItemRenderer) MinSize() fyne.Size {
	width := r.widget.menu.ItemPaddingHorizontal * 2

	if r.icon != nil {
		width += r.widget.menu.IconSize.Width + r.widget.menu.SpacingBetweenIconAndTitle
	}

	titleSize := r.title.MinSize()
	width += titleSize.Width

	if r.subtitle != nil {
		subtitleSize := r.subtitle.MinSize()
		if subtitleSize.Width > titleSize.Width {
			width = width - titleSize.Width + subtitleSize.Width
		}
	}

	return fyne.NewSize(width, r.widget.menu.ItemHeight)
}

func (r *menuItemRenderer) Refresh() {
	r.widget.mu.RLock()
	hovered := r.widget.hovered
	r.widget.mu.RUnlock()

	if hovered && r.widget.item.Enabled {
		r.background.FillColor = r.widget.menu.HighlightedColor
	} else {
		r.background.FillColor = color.Transparent
	}

	if !r.widget.item.Enabled {
		r.title.Color = core.ColorWithAlpha(r.title.Color, 0.5)
	}

	r.background.Refresh()
	r.title.Refresh()
	if r.subtitle != nil {
		r.subtitle.Refresh()
	}
	if r.icon != nil {
		r.icon.Refresh()
	}
}

func (r *menuItemRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background}
	if r.icon != nil {
		objects = append(objects, r.icon)
	}
	objects = append(objects, r.title)
	if r.subtitle != nil {
		objects = append(objects, r.subtitle)
	}
	return objects
}

// ContextMenu shows a context menu at the given position
func ContextMenu(window fyne.Window, position fyne.Position, items []*MenuItem) *PopupMenu {
	menu := NewPopupMenuWithItems(items)
	menu.ArrowDirection = ArrowDirectionNone
	menu.Show(window, position)
	return menu
}
