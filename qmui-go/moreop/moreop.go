// Package moreop provides grid-style action sheet bottom sheets
package moreop

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/animation"
	"github.com/paul-hammant/qmui_fyne/core"
)

// Item represents an item in the operation sheet
type Item struct {
	Identifier  string
	Title       string
	Icon        fyne.Resource
	Handler     func(item *Item)
	Tag         int
	IsEnabled   bool
	ShowsBadge  bool
	BadgeValue  string
}

// NewItem creates a new operation item
func NewItem(identifier, title string, icon fyne.Resource, handler func(*Item)) *Item {
	return &Item{
		Identifier: identifier,
		Title:      title,
		Icon:       icon,
		Handler:    handler,
		IsEnabled:  true,
	}
}

// ItemGroup represents a group/row of items
type ItemGroup struct {
	Items []*Item
}

// NewItemGroup creates a new item group
func NewItemGroup(items ...*Item) *ItemGroup {
	return &ItemGroup{Items: items}
}

// ActionSheet manages a grid-style bottom sheet
type ActionSheet struct {
	// Content
	ItemGroups     []*ItemGroup
	CancelButton   *Item

	// Layout
	ItemsPerRow            int
	ItemHeight             float32
	ItemIconSize           fyne.Size
	ItemTitleFontSize      float32
	ItemSpacing            float32
	GroupSpacing           float32
	ContentInsets          core.EdgeInsets

	// Styling
	BackgroundColor        color.Color
	ItemBackgroundColor    color.Color
	ItemHighlightColor     color.Color
	ItemTitleColor         color.Color
	ItemTitleHighlightColor color.Color
	SeparatorColor         color.Color
	CancelButtonColor      color.Color
	DimmingColor           color.Color
	CornerRadius           float32

	// Animation
	AnimationDuration time.Duration

	// Behavior
	DismissOnTapOutside bool
	DismissOnItemSelected bool

	// Callbacks
	OnShow     func()
	OnDismiss  func()
	OnItemSelected func(item *Item)

	// State
	mu      sync.RWMutex
	window  fyne.Window
	popup   *widget.PopUp
	visible bool
}

// NewActionSheet creates a new operation controller
func NewActionSheet() *ActionSheet {
	config := core.SharedConfiguration()
	moc := &ActionSheet{
		ItemGroups:            make([]*ItemGroup, 0),
		ItemsPerRow:           4,
		ItemHeight:            80,
		ItemIconSize:          fyne.NewSize(44, 44),
		ItemTitleFontSize:     12,
		ItemSpacing:           0,
		GroupSpacing:          8,
		ContentInsets:         core.NewEdgeInsets(16, 0, 0, 0),
		BackgroundColor:       config.SheetHeaderBackgroundColor,
		ItemBackgroundColor:   color.Transparent,
		ItemHighlightColor:    config.SheetButtonHighlightBackgroundColor,
		ItemTitleColor:        config.TableViewCellTitleLabelColor,
		ItemTitleHighlightColor: config.BlueColor,
		SeparatorColor:        config.SeparatorColor,
		CancelButtonColor:     config.BlueColor,
		DimmingColor:          config.MaskDarkColor,
		CornerRadius:          config.SheetContentCornerRadius,
		AnimationDuration:     time.Millisecond * 300,
		DismissOnTapOutside:   true,
		DismissOnItemSelected: true,
	}
	return moc
}

// AddItemGroup adds a group of items
func (moc *ActionSheet) AddItemGroup(group *ItemGroup) {
	moc.mu.Lock()
	moc.ItemGroups = append(moc.ItemGroups, group)
	moc.mu.Unlock()
}

// AddItems adds items as a new group
func (moc *ActionSheet) AddItems(items ...*Item) {
	moc.AddItemGroup(NewItemGroup(items...))
}

// SetCancelButton sets the cancel button
func (moc *ActionSheet) SetCancelButton(title string, handler func(*Item)) {
	moc.mu.Lock()
	moc.CancelButton = &Item{
		Identifier: "cancel",
		Title:      title,
		Handler:    handler,
		IsEnabled:  true,
	}
	moc.mu.Unlock()
}

// Show displays the operation sheet
func (moc *ActionSheet) Show(window fyne.Window) {
	moc.mu.Lock()
	if moc.visible {
		moc.mu.Unlock()
		return
	}
	moc.visible = true
	moc.window = window
	moc.mu.Unlock()

	content := moc.buildContent()

	// Create dimmer that dismisses on tap
	dimmer := &tappableDimmer{
		color:   moc.DimmingColor,
		onTap:   func() {
			if moc.DismissOnTapOutside {
				moc.Dismiss()
			}
		},
	}
	dimmer.ExtendBaseWidget(dimmer)

	// Position content at bottom
	bottomContent := container.NewBorder(nil, content, nil, nil)
	fullContent := container.NewStack(dimmer, bottomContent)

	moc.popup = widget.NewModalPopUp(fullContent, window.Canvas())
	moc.popup.Resize(window.Canvas().Size())

	// Animate in
	moc.animateShow(content)

	if moc.OnShow != nil {
		moc.OnShow()
	}
}

// Dismiss hides the operation sheet
func (moc *ActionSheet) Dismiss() {
	moc.mu.Lock()
	if !moc.visible {
		moc.mu.Unlock()
		return
	}
	moc.mu.Unlock()

	moc.animateHide(func() {
		moc.mu.Lock()
		if moc.popup != nil {
			moc.popup.Hide()
			moc.popup = nil
		}
		moc.visible = false
		moc.mu.Unlock()

		if moc.OnDismiss != nil {
			moc.OnDismiss()
		}
	})
}

// IsVisible returns whether the sheet is visible
func (moc *ActionSheet) IsVisible() bool {
	moc.mu.RLock()
	defer moc.mu.RUnlock()
	return moc.visible
}

func (moc *ActionSheet) animateShow(content fyne.CanvasObject) {
	// Slide up animation
	if moc.window == nil {
		return
	}

	canvasSize := moc.window.Canvas().Size()
	contentSize := content.MinSize()

	startY := float64(canvasSize.Height)
	endY := float64(canvasSize.Height - contentSize.Height)

	animation.NewPropertyAnimation(startY, endY, moc.AnimationDuration, animation.EaseOutCubic, func(y float64) {
		content.Move(fyne.NewPos(0, float32(y)))
	}).Start()
}

func (moc *ActionSheet) animateHide(onComplete func()) {
	if onComplete != nil {
		go func() {
			time.Sleep(moc.AnimationDuration)
			onComplete()
		}()
	}
}

func (moc *ActionSheet) buildContent() fyne.CanvasObject {
	moc.mu.RLock()
	groups := moc.ItemGroups
	cancelBtn := moc.CancelButton
	moc.mu.RUnlock()

	var contentObjects []fyne.CanvasObject

	// Background with rounded top corners
	bg := canvas.NewRectangle(moc.BackgroundColor)
	bg.CornerRadius = moc.CornerRadius

	// Build item groups
	for i, group := range groups {
		groupContent := moc.buildItemGroup(group)
		contentObjects = append(contentObjects, groupContent)

		// Add separator between groups
		if i < len(groups)-1 {
			sep := canvas.NewRectangle(moc.SeparatorColor)
			sep.Resize(fyne.NewSize(0, 8))
			contentObjects = append(contentObjects, sep)
		}
	}

	// Cancel button
	if cancelBtn != nil {
		sep := canvas.NewRectangle(moc.SeparatorColor)
		sep.Resize(fyne.NewSize(0, 8))
		contentObjects = append(contentObjects, sep)

		cancelBtnWidget := moc.buildCancelButton(cancelBtn)
		contentObjects = append(contentObjects, cancelBtnWidget)
	}

	content := container.NewVBox(contentObjects...)
	padded := container.NewPadded(content)

	return container.NewStack(bg, padded)
}

func (moc *ActionSheet) buildItemGroup(group *ItemGroup) fyne.CanvasObject {
	var rows []fyne.CanvasObject

	// Split items into rows based on ItemsPerRow
	for i := 0; i < len(group.Items); i += moc.ItemsPerRow {
		end := i + moc.ItemsPerRow
		if end > len(group.Items) {
			end = len(group.Items)
		}

		rowItems := group.Items[i:end]
		row := moc.buildItemRow(rowItems)
		rows = append(rows, row)
	}

	return container.NewVBox(rows...)
}

func (moc *ActionSheet) buildItemRow(items []*Item) fyne.CanvasObject {
	var widgets []fyne.CanvasObject

	for _, item := range items {
		itemWidget := moc.buildItem(item)
		widgets = append(widgets, itemWidget)
	}

	// Pad with empty containers if needed
	for len(widgets) < moc.ItemsPerRow {
		widgets = append(widgets, widget.NewLabel(""))
	}

	return container.NewGridWithColumns(moc.ItemsPerRow, widgets...)
}

func (moc *ActionSheet) buildItem(item *Item) fyne.CanvasObject {
	itemWidget := &operationItemWidget{
		controller: moc,
		item:       item,
	}
	itemWidget.ExtendBaseWidget(itemWidget)
	return itemWidget
}

func (moc *ActionSheet) buildCancelButton(item *Item) fyne.CanvasObject {
	btn := &cancelButtonWidget{
		controller: moc,
		item:       item,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (moc *ActionSheet) selectItem(item *Item) {
	if item.Handler != nil {
		item.Handler(item)
	}
	if moc.OnItemSelected != nil {
		moc.OnItemSelected(item)
	}
	if moc.DismissOnItemSelected {
		moc.Dismiss()
	}
}


// operationItemWidget represents a single operation item
type operationItemWidget struct {
	widget.BaseWidget
	controller *ActionSheet
	item       *Item
	hovered    bool
	mu         sync.RWMutex
}

func (w *operationItemWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	bg := canvas.NewRectangle(color.Transparent)

	var icon *canvas.Image
	if w.item.Icon != nil {
		icon = canvas.NewImageFromResource(w.item.Icon)
		icon.FillMode = canvas.ImageFillContain
	}

	title := canvas.NewText(w.item.Title, w.controller.ItemTitleColor)
	title.TextSize = w.controller.ItemTitleFontSize
	title.Alignment = fyne.TextAlignCenter

	return &operationItemRenderer{
		widget: w,
		bg:     bg,
		icon:   icon,
		title:  title,
	}
}

func (w *operationItemWidget) Tapped(*fyne.PointEvent) {
	if w.item.IsEnabled {
		w.controller.selectItem(w.item)
	}
}

func (w *operationItemWidget) TappedSecondary(*fyne.PointEvent) {}

func (w *operationItemWidget) MouseIn(*desktop.MouseEvent) {
	w.mu.Lock()
	w.hovered = true
	w.mu.Unlock()
	w.Refresh()
}

func (w *operationItemWidget) MouseMoved(*desktop.MouseEvent) {}

func (w *operationItemWidget) MouseOut() {
	w.mu.Lock()
	w.hovered = false
	w.mu.Unlock()
	w.Refresh()
}

func (w *operationItemWidget) Cursor() desktop.Cursor {
	if w.item.IsEnabled {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

type operationItemRenderer struct {
	widget *operationItemWidget
	bg     *canvas.Rectangle
	icon   *canvas.Image
	title  *canvas.Text
}

func (r *operationItemRenderer) Destroy() {}

func (r *operationItemRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	iconSize := r.widget.controller.ItemIconSize
	titleSize := r.title.MinSize()

	totalHeight := iconSize.Height + 4 + titleSize.Height
	startY := (size.Height - totalHeight) / 2

	if r.icon != nil {
		r.icon.Resize(iconSize)
		r.icon.Move(fyne.NewPos((size.Width-iconSize.Width)/2, startY))
	}

	r.title.Move(fyne.NewPos(0, startY+iconSize.Height+4))
	r.title.Resize(fyne.NewSize(size.Width, titleSize.Height))
}

func (r *operationItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(80, r.widget.controller.ItemHeight)
}

func (r *operationItemRenderer) Refresh() {
	r.widget.mu.RLock()
	hovered := r.widget.hovered
	r.widget.mu.RUnlock()

	if hovered && r.widget.item.IsEnabled {
		r.bg.FillColor = r.widget.controller.ItemHighlightColor
		r.title.Color = r.widget.controller.ItemTitleHighlightColor
	} else {
		r.bg.FillColor = color.Transparent
		r.title.Color = r.widget.controller.ItemTitleColor
	}

	if !r.widget.item.IsEnabled {
		r.title.Color = core.SharedConfiguration().DisabledColor
	}

	r.bg.Refresh()
	r.title.Refresh()
	if r.icon != nil {
		r.icon.Refresh()
	}
}

func (r *operationItemRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bg, r.title}
	if r.icon != nil {
		objects = append(objects, r.icon)
	}
	return objects
}

// cancelButtonWidget represents the cancel button
type cancelButtonWidget struct {
	widget.BaseWidget
	controller *ActionSheet
	item       *Item
	hovered    bool
	mu         sync.RWMutex
}

func (w *cancelButtonWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	bg := canvas.NewRectangle(w.controller.BackgroundColor)
	title := canvas.NewText(w.item.Title, w.controller.CancelButtonColor)
	title.TextSize = 17
	title.Alignment = fyne.TextAlignCenter

	return &cancelButtonRenderer{
		widget: w,
		bg:     bg,
		title:  title,
	}
}

func (w *cancelButtonWidget) Tapped(*fyne.PointEvent) {
	if w.item.Handler != nil {
		w.item.Handler(w.item)
	}
	w.controller.Dismiss()
}

func (w *cancelButtonWidget) TappedSecondary(*fyne.PointEvent) {}

func (w *cancelButtonWidget) MouseIn(*desktop.MouseEvent) {
	w.mu.Lock()
	w.hovered = true
	w.mu.Unlock()
	w.Refresh()
}

func (w *cancelButtonWidget) MouseMoved(*desktop.MouseEvent) {}

func (w *cancelButtonWidget) MouseOut() {
	w.mu.Lock()
	w.hovered = false
	w.mu.Unlock()
	w.Refresh()
}

func (w *cancelButtonWidget) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type cancelButtonRenderer struct {
	widget *cancelButtonWidget
	bg     *canvas.Rectangle
	title  *canvas.Text
}

func (r *cancelButtonRenderer) Destroy() {}

func (r *cancelButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	titleSize := r.title.MinSize()
	r.title.Move(fyne.NewPos(
		(size.Width-titleSize.Width)/2,
		(size.Height-titleSize.Height)/2,
	))
}

func (r *cancelButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 56)
}

func (r *cancelButtonRenderer) Refresh() {
	r.widget.mu.RLock()
	hovered := r.widget.hovered
	r.widget.mu.RUnlock()

	if hovered {
		r.bg.FillColor = r.widget.controller.ItemHighlightColor
	} else {
		r.bg.FillColor = r.widget.controller.BackgroundColor
	}

	r.bg.Refresh()
	r.title.Refresh()
}

func (r *cancelButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.title}
}

// tappableDimmer is a tappable dimmed background
type tappableDimmer struct {
	widget.BaseWidget
	color color.Color
	onTap func()
}

func (d *tappableDimmer) CreateRenderer() fyne.WidgetRenderer {
	d.ExtendBaseWidget(d)
	rect := canvas.NewRectangle(d.color)
	return &dimmerRenderer{dimmer: d, rect: rect}
}

func (d *tappableDimmer) Tapped(*fyne.PointEvent) {
	if d.onTap != nil {
		d.onTap()
	}
}

func (d *tappableDimmer) TappedSecondary(*fyne.PointEvent) {}

type dimmerRenderer struct {
	dimmer *tappableDimmer
	rect   *canvas.Rectangle
}

func (r *dimmerRenderer) Destroy() {}

func (r *dimmerRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *dimmerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *dimmerRenderer) Refresh() {
	r.rect.FillColor = r.dimmer.color
	r.rect.Refresh()
}

func (r *dimmerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect}
}

// Helper functions

// ShowOperationSheet creates and shows an operation sheet with items
func ShowOperationSheet(window fyne.Window, items []*Item, cancelTitle string) *ActionSheet {
	moc := NewActionSheet()
	moc.AddItems(items...)
	if cancelTitle != "" {
		moc.SetCancelButton(cancelTitle, nil)
	}
	moc.Show(window)
	return moc
}

// ShowShareSheet creates a share-style operation sheet
func ShowShareSheet(window fyne.Window, onItemSelected func(item *Item)) *ActionSheet {
	moc := NewActionSheet()

	// Example share items (in real use, you'd pass actual icons)
	shareItems := []*Item{
		NewItem("message", "Message", nil, nil),
		NewItem("mail", "Mail", nil, nil),
		NewItem("copy", "Copy Link", nil, nil),
		NewItem("more", "More", nil, nil),
	}

	moc.AddItems(shareItems...)
	moc.SetCancelButton("Cancel", nil)
	moc.OnItemSelected = onItemSelected
	moc.Show(window)
	return moc
}
