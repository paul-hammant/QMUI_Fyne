// Package alert provides customizable alert and action sheet dialogs
package alert

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// ActionStyle defines the style of an action button
type ActionStyle int

const (
	// ActionStyleDefault is the default button style
	ActionStyleDefault ActionStyle = iota
	// ActionStyleCancel is the cancel button style
	ActionStyleCancel
	// ActionStyleDestructive is the destructive/dangerous action style
	ActionStyleDestructive
)

// ControllerStyle defines the presentation style
type ControllerStyle int

const (
	// ControllerStyleActionSheet presents as a bottom action sheet
	ControllerStyleActionSheet ControllerStyle = iota
	// ControllerStyleAlert presents as a centered alert dialog
	ControllerStyleAlert
)

// Action represents a button action in the alert controller
type Action struct {
	Title   string
	Style   ActionStyle
	Enabled bool
	Handler func(controller *Alert, action *Action)

	// Custom attributes
	TextColor         color.Color
	DisabledTextColor color.Color
	BackgroundColor   color.Color
}

// NewAction creates a new alert action
func NewAction(title string, style ActionStyle, handler func(controller *Alert, action *Action)) *Action {
	return &Action{
		Title:   title,
		Style:   style,
		Enabled: true,
		Handler: handler,
	}
}

// Delegate provides callbacks for alert controller events
type Delegate interface {
	WillShow(controller *Alert)
	DidShow(controller *Alert)
	WillHide(controller *Alert)
	DidHide(controller *Alert)
	ShouldHide(controller *Alert) bool
}

// Alert is a customizable alert/action sheet dialog
type Alert struct {
	widget.BaseWidget

	// Content
	Title   string
	Message string
	Style   ControllerStyle

	// Actions
	actions []*Action

	// Text fields (for alert style)
	textFields []*widget.Entry

	// Custom view
	customView fyne.CanvasObject

	// Delegate
	Delegate Delegate

	// Alert styling
	AlertContentMargin            core.EdgeInsets
	AlertContentMaximumWidth      float32
	AlertSeparatorColor           color.Color
	AlertTitleColor               color.Color
	AlertTitleFontSize            float32
	AlertMessageColor             color.Color
	AlertMessageFontSize          float32
	AlertButtonTextColor          color.Color
	AlertButtonDisabledTextColor  color.Color
	AlertCancelButtonTextColor    color.Color
	AlertDestructiveButtonTextColor color.Color
	AlertContentCornerRadius      float32
	AlertButtonHeight             float32
	AlertHeaderBackgroundColor    color.Color
	AlertButtonBackgroundColor    color.Color
	AlertButtonHighlightBackgroundColor color.Color
	AlertHeaderInsets             core.EdgeInsets
	AlertTitleMessageSpacing      float32

	// Sheet styling
	SheetContentMargin            core.EdgeInsets
	SheetContentMaximumWidth      float32
	SheetSeparatorColor           color.Color
	SheetTitleColor               color.Color
	SheetTitleFontSize            float32
	SheetMessageColor             color.Color
	SheetMessageFontSize          float32
	SheetButtonTextColor          color.Color
	SheetButtonDisabledTextColor  color.Color
	SheetCancelButtonTextColor    color.Color
	SheetDestructiveButtonTextColor color.Color
	SheetCancelButtonMarginTop    float32
	SheetContentCornerRadius      float32
	SheetButtonHeight             float32
	SheetHeaderBackgroundColor    color.Color
	SheetButtonBackgroundColor    color.Color
	SheetButtonHighlightBackgroundColor color.Color
	SheetHeaderInsets             core.EdgeInsets
	SheetTitleMessageSpacing      float32
	SheetButtonColumnCount        int

	// Behavior
	OrderActionsByAddedOrdered   bool
	ShouldRespondDimmingViewTouch bool
	IsExtendBottomLayout         bool

	// State
	mu       sync.RWMutex
	visible  bool
	window   fyne.Window
	overlay  *widget.PopUp
}

// NewAlert creates a new alert controller
func NewAlert(title, message string, style ControllerStyle) *Alert {
	config := core.SharedConfiguration()

	ac := &Alert{
		Title:   title,
		Message: message,
		Style:   style,
		actions: make([]*Action, 0),

		// Alert defaults
		AlertContentMargin:            config.AlertContentMargin,
		AlertContentMaximumWidth:      config.AlertContentMaximumWidth,
		AlertSeparatorColor:           config.AlertSeparatorColor,
		AlertTitleColor:               color.Black,
		AlertTitleFontSize:            17,
		AlertMessageColor:             color.Black,
		AlertMessageFontSize:          13,
		AlertButtonTextColor:          config.BlueColor,
		AlertButtonDisabledTextColor:  config.GrayColor,
		AlertCancelButtonTextColor:    config.BlueColor,
		AlertDestructiveButtonTextColor: config.RedColor,
		AlertContentCornerRadius:      config.AlertContentCornerRadius,
		AlertButtonHeight:             config.AlertButtonHeight,
		AlertHeaderBackgroundColor:    config.AlertHeaderBackgroundColor,
		AlertButtonBackgroundColor:    config.AlertButtonBackgroundColor,
		AlertButtonHighlightBackgroundColor: config.AlertButtonHighlightBackgroundColor,
		AlertHeaderInsets:             config.AlertHeaderInsets,
		AlertTitleMessageSpacing:      config.AlertTitleMessageSpacing,

		// Sheet defaults
		SheetContentMargin:            config.SheetContentMargin,
		SheetContentMaximumWidth:      config.SheetContentMaximumWidth,
		SheetSeparatorColor:           config.SheetSeparatorColor,
		SheetTitleColor:               config.GrayColor,
		SheetTitleFontSize:            13,
		SheetMessageColor:             config.GrayColor,
		SheetMessageFontSize:          13,
		SheetButtonTextColor:          config.BlueColor,
		SheetButtonDisabledTextColor:  config.GrayColor,
		SheetCancelButtonTextColor:    config.BlueColor,
		SheetDestructiveButtonTextColor: config.RedColor,
		SheetCancelButtonMarginTop:    config.SheetCancelButtonMarginTop,
		SheetContentCornerRadius:      config.SheetContentCornerRadius,
		SheetButtonHeight:             config.SheetButtonHeight,
		SheetHeaderBackgroundColor:    config.SheetHeaderBackgroundColor,
		SheetButtonBackgroundColor:    config.SheetButtonBackgroundColor,
		SheetButtonHighlightBackgroundColor: config.SheetButtonHighlightBackgroundColor,
		SheetHeaderInsets:             config.SheetHeaderInsets,
		SheetTitleMessageSpacing:      config.SheetTitleMessageSpacing,
		SheetButtonColumnCount:        config.SheetButtonColumnCount,

		OrderActionsByAddedOrdered:   false,
		ShouldRespondDimmingViewTouch: style == ControllerStyleActionSheet,
		IsExtendBottomLayout:         false,
	}
	ac.ExtendBaseWidget(ac)
	return ac
}

// AddAction adds an action button
func (ac *Alert) AddAction(action *Action) {
	ac.mu.Lock()
	ac.actions = append(ac.actions, action)
	ac.mu.Unlock()
}

// AddCancelAction adds a cancel action
func (ac *Alert) AddCancelAction() {
	ac.AddAction(NewAction("Cancel", ActionStyleCancel, nil))
}

// AddTextField adds a text field (alert style only)
func (ac *Alert) AddTextField(configHandler func(entry *widget.Entry)) {
	entry := widget.NewEntry()
	if configHandler != nil {
		configHandler(entry)
	}
	ac.mu.Lock()
	ac.textFields = append(ac.textFields, entry)
	ac.mu.Unlock()
}

// AddCustomView adds a custom view to the alert
func (ac *Alert) AddCustomView(view fyne.CanvasObject) {
	ac.mu.Lock()
	ac.customView = view
	ac.mu.Unlock()
}

// GetActions returns all actions
func (ac *Alert) GetActions() []*Action {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.actions
}

// GetTextFields returns all text fields
func (ac *Alert) GetTextFields() []*widget.Entry {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.textFields
}

// Show displays the alert controller
func (ac *Alert) ShowIn(window fyne.Window) {
	ac.mu.Lock()
	ac.window = window
	ac.visible = true
	ac.mu.Unlock()

	if ac.Delegate != nil {
		ac.Delegate.WillShow(ac)
	}

	content := ac.buildContent()

	// For ActionSheet style with ShouldRespondDimmingViewTouch, use PopUp so
	// tapping outside dismisses it
	if ac.Style == ControllerStyleActionSheet && ac.ShouldRespondDimmingViewTouch {
		ac.overlay = widget.NewPopUp(content, window.Canvas())
		// Position at bottom of screen
		canvasSize := window.Canvas().Size()
		contentSize := content.MinSize()
		ac.overlay.Move(fyne.NewPos(
			(canvasSize.Width-contentSize.Width)/2,
			canvasSize.Height-contentSize.Height-20,
		))
		ac.overlay.Show()
	} else {
		ac.overlay = widget.NewModalPopUp(content, window.Canvas())
		ac.overlay.Show()
	}

	if ac.Delegate != nil {
		ac.Delegate.DidShow(ac)
	}
}

// ShowWithAnimated displays the alert with optional animation
func (ac *Alert) ShowWithAnimated(window fyne.Window, animated bool) {
	ac.ShowIn(window)
}

// Hide hides the alert controller
func (ac *Alert) Hide() {
	ac.mu.Lock()
	visible := ac.visible
	ac.mu.Unlock()

	if !visible {
		return
	}

	if ac.Delegate != nil {
		if !ac.Delegate.ShouldHide(ac) {
			return
		}
		ac.Delegate.WillHide(ac)
	}

	if ac.overlay != nil {
		ac.overlay.Hide()
	}

	ac.mu.Lock()
	ac.visible = false
	ac.mu.Unlock()

	if ac.Delegate != nil {
		ac.Delegate.DidHide(ac)
	}
}

// HideWithAnimated hides the alert with optional animation
func (ac *Alert) HideWithAnimated(animated bool) {
	ac.Hide()
}

func (ac *Alert) buildContent() fyne.CanvasObject {
	if ac.Style == ControllerStyleAlert {
		return ac.buildAlertContent()
	}
	return ac.buildActionSheetContent()
}

func (ac *Alert) buildAlertContent() fyne.CanvasObject {
	// Header (title + message)
	var headerObjects []fyne.CanvasObject

	if ac.Title != "" {
		titleLabel := canvas.NewText(ac.Title, ac.AlertTitleColor)
		titleLabel.TextStyle = fyne.TextStyle{Bold: true}
		titleLabel.TextSize = ac.AlertTitleFontSize
		titleLabel.Alignment = fyne.TextAlignCenter
		headerObjects = append(headerObjects, titleLabel)
	}

	if ac.Message != "" {
		messageLabel := canvas.NewText(ac.Message, ac.AlertMessageColor)
		messageLabel.TextSize = ac.AlertMessageFontSize
		messageLabel.Alignment = fyne.TextAlignCenter
		headerObjects = append(headerObjects, messageLabel)
	}

	// Text fields
	for _, tf := range ac.textFields {
		headerObjects = append(headerObjects, tf)
	}

	// Custom view
	if ac.customView != nil {
		headerObjects = append(headerObjects, ac.customView)
	}

	header := container.NewVBox(headerObjects...)

	// Action buttons
	var buttonObjects []fyne.CanvasObject
	ac.mu.RLock()
	actions := ac.actions
	ac.mu.RUnlock()

	for _, action := range actions {
		btn := ac.createActionButton(action)
		buttonObjects = append(buttonObjects, btn)
	}

	var buttons fyne.CanvasObject
	if len(buttonObjects) == 2 {
		// Two buttons side by side
		buttons = container.NewGridWithColumns(2, buttonObjects...)
	} else {
		// Stacked vertically
		buttons = container.NewVBox(buttonObjects...)
	}

	// Main container
	background := canvas.NewRectangle(ac.AlertHeaderBackgroundColor)
	background.CornerRadius = ac.AlertContentCornerRadius

	content := container.NewVBox(
		container.NewPadded(header),
		canvas.NewLine(ac.AlertSeparatorColor),
		buttons,
	)

	return container.NewStack(
		background,
		container.NewPadded(content),
	)
}

func (ac *Alert) buildActionSheetContent() fyne.CanvasObject {
	// Header (title + message)
	var headerObjects []fyne.CanvasObject

	if ac.Title != "" {
		titleLabel := canvas.NewText(ac.Title, ac.SheetTitleColor)
		titleLabel.TextStyle = fyne.TextStyle{Bold: true}
		titleLabel.TextSize = ac.SheetTitleFontSize
		titleLabel.Alignment = fyne.TextAlignCenter
		headerObjects = append(headerObjects, titleLabel)
	}

	if ac.Message != "" {
		messageLabel := canvas.NewText(ac.Message, ac.SheetMessageColor)
		messageLabel.TextSize = ac.SheetMessageFontSize
		messageLabel.Alignment = fyne.TextAlignCenter
		headerObjects = append(headerObjects, messageLabel)
	}

	// Separate cancel action from other actions
	var regularActions []*Action
	var cancelAction *Action

	ac.mu.RLock()
	for _, action := range ac.actions {
		if action.Style == ActionStyleCancel {
			cancelAction = action
		} else {
			regularActions = append(regularActions, action)
		}
	}
	ac.mu.RUnlock()

	// Build main container with header and regular action buttons
	var mainObjects []fyne.CanvasObject

	if len(headerObjects) > 0 {
		header := container.NewVBox(headerObjects...)
		mainObjects = append(mainObjects, container.NewPadded(header))
		mainObjects = append(mainObjects, canvas.NewLine(ac.SheetSeparatorColor))
	}

	// Regular action buttons
	for i, action := range regularActions {
		btn := ac.createActionButton(action)
		mainObjects = append(mainObjects, btn)
		if i < len(regularActions)-1 {
			mainObjects = append(mainObjects, canvas.NewLine(ac.SheetSeparatorColor))
		}
	}

	mainBackground := canvas.NewRectangle(ac.SheetHeaderBackgroundColor)
	mainBackground.CornerRadius = ac.SheetContentCornerRadius

	mainContent := container.NewStack(
		mainBackground,
		container.NewVBox(mainObjects...),
	)

	// Cancel button (separate container)
	var allContent fyne.CanvasObject
	if cancelAction != nil {
		cancelBackground := canvas.NewRectangle(ac.SheetButtonBackgroundColor)
		cancelBackground.CornerRadius = ac.SheetContentCornerRadius

		cancelBtn := ac.createActionButton(cancelAction)
		cancelContent := container.NewStack(
			cancelBackground,
			cancelBtn,
		)

		allContent = container.NewVBox(
			mainContent,
			layout.NewSpacer(),
			cancelContent,
		)
	} else {
		allContent = mainContent
	}

	return container.NewPadded(allContent)
}

func (ac *Alert) createActionButton(action *Action) fyne.CanvasObject {
	var textColor color.Color
	var fontSize float32
	var buttonHeight float32

	if ac.Style == ControllerStyleAlert {
		fontSize = 17
		buttonHeight = ac.AlertButtonHeight

		switch action.Style {
		case ActionStyleCancel:
			textColor = ac.AlertCancelButtonTextColor
		case ActionStyleDestructive:
			textColor = ac.AlertDestructiveButtonTextColor
		default:
			textColor = ac.AlertButtonTextColor
		}

		if !action.Enabled {
			textColor = ac.AlertButtonDisabledTextColor
		}
	} else {
		fontSize = 20
		buttonHeight = ac.SheetButtonHeight

		switch action.Style {
		case ActionStyleCancel:
			textColor = ac.SheetCancelButtonTextColor
		case ActionStyleDestructive:
			textColor = ac.SheetDestructiveButtonTextColor
		default:
			textColor = ac.SheetButtonTextColor
		}

		if !action.Enabled {
			textColor = ac.SheetButtonDisabledTextColor
		}
	}

	// Override with custom colors if set
	if action.TextColor != nil {
		textColor = action.TextColor
	}
	if !action.Enabled && action.DisabledTextColor != nil {
		textColor = action.DisabledTextColor
	}

	label := canvas.NewText(action.Title, textColor)
	label.TextSize = fontSize
	if action.Style == ActionStyleCancel {
		label.TextStyle = fyne.TextStyle{Bold: true}
	}
	label.Alignment = fyne.TextAlignCenter

	btn := &actionButton{
		label:       label,
		action:      action,
		controller:  ac,
		height:      buttonHeight,
		normalBg:    color.Transparent,
		highlightBg: ac.AlertButtonHighlightBackgroundColor,
	}
	if ac.Style == ControllerStyleActionSheet {
		btn.highlightBg = ac.SheetButtonHighlightBackgroundColor
	}
	btn.ExtendBaseWidget(btn)

	return btn
}

type actionButton struct {
	widget.BaseWidget

	label       *canvas.Text
	action      *Action
	controller  *Alert
	height      float32
	normalBg    color.Color
	highlightBg color.Color
	hovered     bool
	pressed     bool
	mu          sync.RWMutex
}

func (b *actionButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	background := canvas.NewRectangle(b.normalBg)
	return &actionButtonRenderer{
		button:     b,
		background: background,
		label:      b.label,
	}
}

func (b *actionButton) Tapped(_ *fyne.PointEvent) {
	if !b.action.Enabled {
		return
	}
	if b.action.Handler != nil {
		b.action.Handler(b.controller, b.action)
	}
	b.controller.Hide()
}

func (b *actionButton) TappedSecondary(_ *fyne.PointEvent) {}

func (b *actionButton) MouseIn(_ *desktop.MouseEvent) {
	b.mu.Lock()
	b.hovered = true
	b.mu.Unlock()
	b.Refresh()
}

func (b *actionButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *actionButton) MouseOut() {
	b.mu.Lock()
	b.hovered = false
	b.mu.Unlock()
	b.Refresh()
}

func (b *actionButton) Cursor() desktop.Cursor {
	if b.action.Enabled {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

type actionButtonRenderer struct {
	button     *actionButton
	background *canvas.Rectangle
	label      *canvas.Text
}

func (r *actionButtonRenderer) Destroy() {}

func (r *actionButtonRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	labelSize := r.label.MinSize()
	r.label.Move(fyne.NewPos(
		(size.Width-labelSize.Width)/2,
		(size.Height-labelSize.Height)/2,
	))
}

func (r *actionButtonRenderer) MinSize() fyne.Size {
	labelSize := r.label.MinSize()
	return fyne.NewSize(labelSize.Width+32, r.button.height)
}

func (r *actionButtonRenderer) Refresh() {
	r.button.mu.RLock()
	hovered := r.button.hovered
	r.button.mu.RUnlock()

	if hovered {
		r.background.FillColor = r.button.highlightBg
	} else {
		r.background.FillColor = r.button.normalBg
	}
	r.background.Refresh()
	r.label.Refresh()
}

func (r *actionButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.label}
}

// CreateRenderer implements fyne.Widget
func (ac *Alert) CreateRenderer() fyne.WidgetRenderer {
	ac.ExtendBaseWidget(ac)
	return &alertRenderer{controller: ac}
}

type alertRenderer struct {
	controller *Alert
}

func (r *alertRenderer) Destroy()                      {}
func (r *alertRenderer) Layout(size fyne.Size)         {}
func (r *alertRenderer) MinSize() fyne.Size            { return fyne.NewSize(0, 0) }
func (r *alertRenderer) Refresh()                      {}
func (r *alertRenderer) Objects() []fyne.CanvasObject { return nil }

// IsAnyAlertVisible returns whether any alert is visible
var activeAlerts sync.Map

func IsAnyAlertVisible() bool {
	visible := false
	activeAlerts.Range(func(_, _ interface{}) bool {
		visible = true
		return false
	})
	return visible
}

// ShowAlert is a convenience function to show a simple alert
func ShowAlert(window fyne.Window, title, message string, actions ...*Action) *Alert {
	ac := NewAlert(title, message, ControllerStyleAlert)
	for _, action := range actions {
		ac.AddAction(action)
	}
	if len(actions) == 0 {
		ac.AddAction(NewAction("OK", ActionStyleDefault, nil))
	}
	ac.ShowIn(window)
	return ac
}

// ShowConfirm is a convenience function to show a confirmation dialog
func ShowConfirm(window fyne.Window, title, message string, onConfirm func(), onCancel func()) *Alert {
	ac := NewAlert(title, message, ControllerStyleAlert)
	ac.AddAction(NewAction("Cancel", ActionStyleCancel, func(_ *Alert, _ *Action) {
		if onCancel != nil {
			onCancel()
		}
	}))
	ac.AddAction(NewAction("OK", ActionStyleDefault, func(_ *Alert, _ *Action) {
		if onConfirm != nil {
			onConfirm()
		}
	}))
	ac.ShowIn(window)
	return ac
}

// ShowActionSheet is a convenience function to show an action sheet
func ShowActionSheet(window fyne.Window, title, message string, actions ...*Action) *Alert {
	ac := NewAlert(title, message, ControllerStyleActionSheet)
	for _, action := range actions {
		ac.AddAction(action)
	}
	ac.AddCancelAction()
	ac.ShowIn(window)
	return ac
}

// ShowDestructiveConfirm shows a confirmation with destructive action
func ShowDestructiveConfirm(window fyne.Window, title, message, destructiveTitle string, onDestructive func(), onCancel func()) *Alert {
	ac := NewAlert(title, message, ControllerStyleAlert)
	ac.AddAction(NewAction("Cancel", ActionStyleCancel, func(_ *Alert, _ *Action) {
		if onCancel != nil {
			onCancel()
		}
	}))
	ac.AddAction(NewAction(destructiveTitle, ActionStyleDestructive, func(_ *Alert, _ *Action) {
		if onDestructive != nil {
			onDestructive()
		}
	}))
	ac.ShowIn(window)
	return ac
}

// ShowTextInput shows an alert with a text input field
func ShowTextInput(window fyne.Window, title, message, placeholder string, onSubmit func(text string), onCancel func()) *Alert {
	ac := NewAlert(title, message, ControllerStyleAlert)
	var textField *widget.Entry
	ac.AddTextField(func(entry *widget.Entry) {
		entry.SetPlaceHolder(placeholder)
		textField = entry
	})
	ac.AddAction(NewAction("Cancel", ActionStyleCancel, func(_ *Alert, _ *Action) {
		if onCancel != nil {
			onCancel()
		}
	}))
	ac.AddAction(NewAction("OK", ActionStyleDefault, func(_ *Alert, _ *Action) {
		if onSubmit != nil && textField != nil {
			onSubmit(textField.Text)
		}
	}))
	ac.ShowIn(window)
	return ac
}
