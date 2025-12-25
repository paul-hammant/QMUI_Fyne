// Package dialog provides customizable content dialogs
package dialog

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/animation"
	"github.com/user/qmui-go/pkg/core"
)

// DialogStyle defines the style of dialog
type DialogStyle int

const (
	DialogStyleDefault DialogStyle = iota
	DialogStyleSheet
)

// DialogAction represents a dialog button action
type DialogAction struct {
	Title       string
	Style       DialogActionStyle
	Handler     func()
	IsEnabled   bool
	textColor   color.Color
}

// DialogActionStyle defines the style of an action button
type DialogActionStyle int

const (
	DialogActionStyleDefault DialogActionStyle = iota
	DialogActionStyleCancel
	DialogActionStyleDestructive
)

// Short aliases for action styles
const (
	ActionStyleDefault     = DialogActionStyleDefault
	ActionStyleCancel      = DialogActionStyleCancel
	ActionStyleDestructive = DialogActionStyleDestructive
)

// NewDialogAction creates a new dialog action
func NewDialogAction(title string, style DialogActionStyle) *DialogAction {
	return &DialogAction{
		Title:     title,
		Style:     style,
		IsEnabled: true,
	}
}

// NewDialogActionWithHandler creates a dialog action with a handler
func NewDialogActionWithHandler(title string, style DialogActionStyle, handler func(*DialogAction)) *DialogAction {
	action := &DialogAction{
		Title:     title,
		Style:     style,
		IsEnabled: true,
	}
	if handler != nil {
		action.Handler = func() { handler(action) }
	}
	return action
}

// NewCancelAction creates a cancel action
func NewCancelAction(title string, handler func()) *DialogAction {
	return &DialogAction{
		Title:     title,
		Style:     DialogActionStyleCancel,
		Handler:   handler,
		IsEnabled: true,
	}
}

// NewDestructiveAction creates a destructive (red) action
func NewDestructiveAction(title string, handler func()) *DialogAction {
	return &DialogAction{
		Title:     title,
		Style:     DialogActionStyleDestructive,
		Handler:   handler,
		IsEnabled: true,
	}
}

// Dialog manages a custom content dialog
type Dialog struct {
	// Content
	Title           string
	Message         string
	ContentView     fyne.CanvasObject
	HeaderView      fyne.CanvasObject
	FooterView      fyne.CanvasObject
	Actions         []*DialogAction

	// Styling
	Style                  DialogStyle
	BackgroundColor        color.Color
	CornerRadius           float32
	ContentInsets          core.EdgeInsets
	TitleColor             color.Color
	MessageColor           color.Color
	TitleFontSize          float32
	MessageFontSize        float32
	TitleMessageSpacing    float32
	SeparatorColor         color.Color
	ButtonHeight           float32
	ButtonBackgroundColor  color.Color
	ButtonHighlightColor   color.Color
	MaxWidth               float32
	DimmedBackgroundColor  color.Color

	// Behavior
	DismissOnTapOutside    bool
	AnimationDuration      time.Duration
	ShowsCloseButton       bool

	// Callbacks
	OnShow    func()
	OnDismiss func()
	OnClose   func()

	// State
	mu        sync.RWMutex
	window    fyne.Window
	popup     *widget.PopUp
	visible   bool
}

// NewDialog creates a new dialog controller
func NewDialog() *Dialog {
	config := core.SharedConfiguration()
	dvc := &Dialog{
		Style:                 DialogStyleDefault,
		BackgroundColor:       config.AlertHeaderBackgroundColor,
		CornerRadius:          config.AlertContentCornerRadius,
		ContentInsets:         config.AlertHeaderInsets,
		TitleColor:            config.NavBarTitleColor,
		MessageColor:          config.GrayColor,
		TitleFontSize:         17,
		MessageFontSize:       13,
		TitleMessageSpacing:   config.AlertTitleMessageSpacing,
		SeparatorColor:        config.AlertSeparatorColor,
		ButtonHeight:          config.AlertButtonHeight,
		ButtonBackgroundColor: config.AlertButtonBackgroundColor,
		ButtonHighlightColor:  config.AlertButtonHighlightBackgroundColor,
		MaxWidth:              config.AlertContentMaximumWidth,
		DimmedBackgroundColor: config.MaskDarkColor,
		DismissOnTapOutside:   true,
		AnimationDuration:     time.Millisecond * 250,
		ShowsCloseButton:      false,
		Actions:               make([]*DialogAction, 0),
	}
	return dvc
}

// NewDialogWithTitle creates a dialog with a title
func NewDialogWithTitle(title string) *Dialog {
	dvc := NewDialog()
	dvc.Title = title
	return dvc
}

// NewDialogWithContent creates a dialog with custom content
func NewDialogWithContent(content fyne.CanvasObject) *Dialog {
	dvc := NewDialog()
	dvc.ContentView = content
	return dvc
}

// AddAction adds an action button
func (dvc *Dialog) AddAction(action *DialogAction) {
	dvc.mu.Lock()
	dvc.Actions = append(dvc.Actions, action)
	dvc.mu.Unlock()
}

// AddCancelAction adds a standard cancel action
func (dvc *Dialog) AddCancelAction(title string, handler func()) {
	dvc.AddAction(NewCancelAction(title, handler))
}

// AddSubmitAction adds a standard submit action
func (dvc *Dialog) AddSubmitAction(title string, handler func()) {
	action := &DialogAction{
		Title:     title,
		Style:     DialogActionStyleDefault,
		Handler:   handler,
		IsEnabled: true,
	}
	dvc.AddAction(action)
}

// SetContentView sets the custom content view
func (dvc *Dialog) SetContentView(view fyne.CanvasObject) {
	dvc.mu.Lock()
	dvc.ContentView = view
	dvc.mu.Unlock()
}

// Show displays the dialog in the given window
func (dvc *Dialog) Show(window fyne.Window) {
	dvc.mu.Lock()
	if dvc.visible {
		dvc.mu.Unlock()
		return
	}
	dvc.visible = true
	dvc.window = window
	dvc.mu.Unlock()

	content := dvc.buildDialogContent()

	// Create dimmed background
	dimmer := canvas.NewRectangle(dvc.DimmedBackgroundColor)

	// Create popup content with dimmer
	contentWithDimmer := container.NewStack(dimmer, content)

	dvc.popup = widget.NewModalPopUp(contentWithDimmer, window.Canvas())
	dvc.popup.Resize(window.Canvas().Size())

	// Animate in
	dvc.animateShow()

	if dvc.OnShow != nil {
		dvc.OnShow()
	}
}

// Dismiss hides the dialog
func (dvc *Dialog) Dismiss() {
	dvc.mu.Lock()
	if !dvc.visible {
		dvc.mu.Unlock()
		return
	}
	dvc.mu.Unlock()

	dvc.animateHide(func() {
		fyne.Do(func() {
			dvc.mu.Lock()
			if dvc.popup != nil {
				dvc.popup.Hide()
				dvc.popup = nil
			}
			dvc.visible = false
			dvc.mu.Unlock()

			if dvc.OnDismiss != nil {
				dvc.OnDismiss()
			}
		})
	})
}

func (dvc *Dialog) animateShow() {
	// Simple fade-in animation
	if dvc.popup != nil {
		dvc.popup.Show()
	}
}

func (dvc *Dialog) animateHide(onComplete func()) {
	// Simple fade-out animation
	animation.AnimateFloat(1, 0, dvc.AnimationDuration, animation.EaseOutQuad, func(v float64) {
		// Would animate opacity here if Fyne supported it
	}, onComplete)
}

func (dvc *Dialog) buildDialogContent() fyne.CanvasObject {
	config := core.SharedConfiguration()
	var contentObjects []fyne.CanvasObject

	// Header (title + message)
	if dvc.Title != "" || dvc.Message != "" {
		header := dvc.buildHeader()
		contentObjects = append(contentObjects, header)
	}

	// Custom header view
	if dvc.HeaderView != nil {
		contentObjects = append(contentObjects, dvc.HeaderView)
	}

	// Custom content view
	if dvc.ContentView != nil {
		contentObjects = append(contentObjects, container.NewPadded(dvc.ContentView))
	}

	// Custom footer view
	if dvc.FooterView != nil {
		contentObjects = append(contentObjects, dvc.FooterView)
	}

	// Action buttons
	if len(dvc.Actions) > 0 {
		separator := canvas.NewRectangle(dvc.SeparatorColor)
		separator.Resize(fyne.NewSize(dvc.MaxWidth, 0.5))
		contentObjects = append(contentObjects, separator)

		buttons := dvc.buildButtons()
		contentObjects = append(contentObjects, buttons)
	}

	// Close button
	if dvc.ShowsCloseButton {
		closeBtn := widget.NewButton("", func() {
			dvc.Dismiss()
			if dvc.OnClose != nil {
				dvc.OnClose()
			}
		})
		closeBtn.Importance = widget.LowImportance
	}

	// Background
	background := canvas.NewRectangle(dvc.BackgroundColor)
	background.CornerRadius = dvc.CornerRadius

	content := container.NewVBox(contentObjects...)

	// Constrain width
	dialogContent := container.NewStack(
		background,
		content,
	)

	// Wrap to constrain size
	wrapper := &dialogWrapper{
		content:  dialogContent,
		maxWidth: dvc.MaxWidth,
	}
	wrapper.ExtendBaseWidget(wrapper)

	// Center in screen
	centered := container.NewCenter(wrapper)

	_ = config // avoid unused warning
	return centered
}

func (dvc *Dialog) buildHeader() fyne.CanvasObject {
	var headerObjects []fyne.CanvasObject

	if dvc.Title != "" {
		title := canvas.NewText(dvc.Title, dvc.TitleColor)
		title.TextSize = dvc.TitleFontSize
		title.TextStyle = fyne.TextStyle{Bold: true}
		title.Alignment = fyne.TextAlignCenter
		headerObjects = append(headerObjects, title)
	}

	if dvc.Message != "" {
		message := canvas.NewText(dvc.Message, dvc.MessageColor)
		message.TextSize = dvc.MessageFontSize
		message.Alignment = fyne.TextAlignCenter
		headerObjects = append(headerObjects, message)
	}

	header := container.NewVBox(headerObjects...)
	return container.NewPadded(header)
}

func (dvc *Dialog) buildButtons() fyne.CanvasObject {
	config := core.SharedConfiguration()

	// Single button: full width
	// Two buttons: side by side
	// More buttons: stacked vertically

	if len(dvc.Actions) == 2 {
		// Side by side
		btn1 := dvc.createButton(dvc.Actions[0])
		btn2 := dvc.createButton(dvc.Actions[1])

		sep := canvas.NewRectangle(dvc.SeparatorColor)
		sep.Resize(fyne.NewSize(0.5, dvc.ButtonHeight))

		return container.NewGridWithColumns(2, btn1, btn2)
	}

	// Vertical stack
	var buttons []fyne.CanvasObject
	for i, action := range dvc.Actions {
		btn := dvc.createButton(action)
		buttons = append(buttons, btn)

		if i < len(dvc.Actions)-1 {
			sep := canvas.NewRectangle(dvc.SeparatorColor)
			sep.Resize(fyne.NewSize(dvc.MaxWidth, 0.5))
			buttons = append(buttons, sep)
		}
	}

	_ = config
	return container.NewVBox(buttons...)
}

func (dvc *Dialog) createButton(action *DialogAction) fyne.CanvasObject {
	config := core.SharedConfiguration()

	// Determine text color based on style
	textColor := config.BlueColor
	switch action.Style {
	case DialogActionStyleCancel:
		textColor = config.GrayDarkenColor
	case DialogActionStyleDestructive:
		textColor = config.RedColor
	}

	btn := &dialogButton{
		action:           action,
		dialog:           dvc,
		textColor:        textColor,
		backgroundColor:  dvc.ButtonBackgroundColor,
		highlightColor:   dvc.ButtonHighlightColor,
		height:           dvc.ButtonHeight,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}


// dialogWrapper constrains dialog width
type dialogWrapper struct {
	widget.BaseWidget
	content  fyne.CanvasObject
	maxWidth float32
}

func (w *dialogWrapper) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	return &dialogWrapperRenderer{wrapper: w}
}

type dialogWrapperRenderer struct {
	wrapper *dialogWrapper
}

func (r *dialogWrapperRenderer) Destroy() {}

func (r *dialogWrapperRenderer) Layout(size fyne.Size) {
	contentWidth := size.Width
	if contentWidth > r.wrapper.maxWidth {
		contentWidth = r.wrapper.maxWidth
	}
	r.wrapper.content.Resize(fyne.NewSize(contentWidth, size.Height))
	r.wrapper.content.Move(fyne.NewPos((size.Width-contentWidth)/2, 0))
}

func (r *dialogWrapperRenderer) MinSize() fyne.Size {
	min := r.wrapper.content.MinSize()
	if min.Width > r.wrapper.maxWidth {
		min.Width = r.wrapper.maxWidth
	}
	return min
}

func (r *dialogWrapperRenderer) Refresh() {
	r.wrapper.content.Refresh()
}

func (r *dialogWrapperRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.wrapper.content}
}

// dialogButton is a styled button for dialogs
type dialogButton struct {
	widget.BaseWidget
	action          *DialogAction
	dialog          *Dialog
	textColor       color.Color
	backgroundColor color.Color
	highlightColor  color.Color
	height          float32
	hovered         bool
	mu              sync.RWMutex
}

func (b *dialogButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	bg := canvas.NewRectangle(b.backgroundColor)
	text := canvas.NewText(b.action.Title, b.textColor)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 17

	return &dialogButtonRenderer{
		button: b,
		bg:     bg,
		text:   text,
	}
}

func (b *dialogButton) Tapped(*fyne.PointEvent) {
	if b.action.Handler != nil {
		b.action.Handler()
	}
	b.dialog.Dismiss()
}

func (b *dialogButton) TappedSecondary(*fyne.PointEvent) {}

func (b *dialogButton) MouseIn(*desktop.MouseEvent) {
	b.mu.Lock()
	b.hovered = true
	b.mu.Unlock()
	b.Refresh()
}

func (b *dialogButton) MouseMoved(*desktop.MouseEvent) {}

func (b *dialogButton) MouseOut() {
	b.mu.Lock()
	b.hovered = false
	b.mu.Unlock()
	b.Refresh()
}

func (b *dialogButton) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type dialogButtonRenderer struct {
	button *dialogButton
	bg     *canvas.Rectangle
	text   *canvas.Text
}

func (r *dialogButtonRenderer) Destroy() {}

func (r *dialogButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	textSize := r.text.MinSize()
	r.text.Move(fyne.NewPos(
		(size.Width-textSize.Width)/2,
		(size.Height-textSize.Height)/2,
	))
}

func (r *dialogButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, r.button.height)
}

func (r *dialogButtonRenderer) Refresh() {
	r.button.mu.RLock()
	hovered := r.button.hovered
	r.button.mu.RUnlock()

	if hovered {
		r.bg.FillColor = r.button.highlightColor
	} else {
		r.bg.FillColor = r.button.backgroundColor
	}
	r.bg.Refresh()
	r.text.Refresh()
}

func (r *dialogButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.text}
}

// Helper functions

// ShowDialog shows a simple dialog with title, message, and actions
func ShowDialog(window fyne.Window, title, message string, actions ...*DialogAction) *Dialog {
	dvc := NewDialog()
	dvc.Title = title
	dvc.Message = message
	for _, action := range actions {
		dvc.AddAction(action)
	}
	dvc.Show(window)
	return dvc
}

// ShowContentDialog shows a dialog with custom content
func ShowContentDialog(window fyne.Window, title string, content fyne.CanvasObject, actions ...*DialogAction) *Dialog {
	dvc := NewDialog()
	dvc.Title = title
	dvc.ContentView = content
	for _, action := range actions {
		dvc.AddAction(action)
	}
	dvc.Show(window)
	return dvc
}

// ShowConfirmDialog shows a confirmation dialog
func ShowConfirmDialog(window fyne.Window, title, message string, onConfirm, onCancel func()) *Dialog {
	confirmAction := &DialogAction{
		Title:     "Confirm",
		Style:     DialogActionStyleDefault,
		Handler:   onConfirm,
		IsEnabled: true,
	}
	return ShowDialog(window, title, message,
		NewCancelAction("Cancel", onCancel),
		confirmAction,
	)
}

// ShowInputDialog shows a dialog with a text input
func ShowInputDialog(window fyne.Window, title, placeholder string, onSubmit func(text string)) *Dialog {
	entry := widget.NewEntry()
	entry.PlaceHolder = placeholder

	dvc := NewDialog()
	dvc.Title = title
	dvc.ContentView = entry
	dvc.AddCancelAction("Cancel", nil)
	dvc.AddSubmitAction("Submit", func() {
		if onSubmit != nil {
			onSubmit(entry.Text)
		}
	})
	dvc.Show(window)
	return dvc
}
