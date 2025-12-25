// Package button provides QMUIButton - an enhanced button with advanced styling
// Ported from Tencent's QMUI_iOS framework
package button

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/paul-hammant/qmui_fyne/core"
	qmuitheme "github.com/paul-hammant/qmui_fyne/theme"
)

// ImagePosition controls the image position relative to the title
type ImagePosition int

const (
	// ImagePositionTop places image above title
	ImagePositionTop ImagePosition = iota
	// ImagePositionLeft places image left of title
	ImagePositionLeft
	// ImagePositionBottom places image below title
	ImagePositionBottom
	// ImagePositionRight places image right of title
	ImagePositionRight
)

// CornerRadiusAdjustsBounds is a special value that auto-adjusts corner radius to half the height
const CornerRadiusAdjustsBounds float32 = -1

// Button is an enhanced button widget with QMUI styling features
type Button struct {
	widget.BaseWidget

	// Text content
	Text     string
	Subtitle string

	// Icon
	Icon         fyne.Resource
	IconPosition ImagePosition

	// Styling
	TintColor                       color.Color
	HighlightedBackgroundColor      color.Color
	HighlightedBorderColor          color.Color
	DisabledColor                   color.Color
	BackgroundColor                 color.Color
	BorderColor                     color.Color
	BorderWidth                     float32
	CornerRadius                    float32
	SpacingBetweenImageAndTitle     float32
	SubtitleEdgeInsets              core.EdgeInsets
	SubtitleColor                   color.Color
	ContentEdgeInsets               core.EdgeInsets

	// Behavior
	AdjustsTitleTintColorAutomatically bool
	AdjustsImageTintColorAutomatically bool
	AdjustsButtonWhenHighlighted       bool
	AdjustsButtonWhenDisabled          bool
	Enabled                            bool

	// Callbacks
	OnTapped func()

	// State
	mu          sync.RWMutex
	hovered     bool
	pressed     bool
	highlighted bool
}

// NewButton creates a new QMUI-styled button with text
func NewButton(text string, tapped func()) *Button {
	config := core.SharedConfiguration()
	btn := &Button{
		Text:                           text,
		OnTapped:                       tapped,
		IconPosition:                   ImagePositionLeft,
		TintColor:                      config.ButtonTintColor,
		AdjustsTitleTintColorAutomatically: false,
		AdjustsImageTintColorAutomatically: false,
		AdjustsButtonWhenHighlighted:   true,
		AdjustsButtonWhenDisabled:      true,
		Enabled:                        true,
		SpacingBetweenImageAndTitle:    4,
		CornerRadius:                   0,
		ContentEdgeInsets:              core.NewEdgeInsets(8, 16, 8, 16),
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

// NewButtonWithIcon creates a new QMUI-styled button with icon and text
func NewButtonWithIcon(text string, icon fyne.Resource, tapped func()) *Button {
	btn := NewButton(text, tapped)
	btn.Icon = icon
	return btn
}

// NewIconButton creates a new QMUI-styled button with only an icon
func NewIconButton(icon fyne.Resource, tapped func()) *Button {
	btn := NewButton("", tapped)
	btn.Icon = icon
	btn.ContentEdgeInsets = core.NewEdgeInsets(8, 8, 8, 8)
	return btn
}

// SetTintColorAdjustsTitleAndImage sets both title and image to follow tint color
func (b *Button) SetTintColorAdjustsTitleAndImage(tintColor color.Color) {
	b.mu.Lock()
	b.TintColor = tintColor
	b.AdjustsTitleTintColorAutomatically = true
	b.AdjustsImageTintColorAutomatically = true
	b.mu.Unlock()
	b.Refresh()
}

// SetEnabled sets the button enabled state
func (b *Button) SetEnabled(enabled bool) {
	b.mu.Lock()
	b.Enabled = enabled
	b.mu.Unlock()
	b.Refresh()
}

// IsEnabled returns whether the button is enabled
func (b *Button) IsEnabled() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Enabled
}

// SetHighlighted sets the highlighted state
func (b *Button) SetHighlighted(highlighted bool) {
	b.mu.Lock()
	b.highlighted = highlighted
	b.mu.Unlock()
	b.Refresh()
}

// CreateRenderer implements fyne.Widget
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	// Initialize background with actual color (not transparent)
	bgColor := color.Color(color.Transparent)
	if b.BackgroundColor != nil {
		bgColor = b.BackgroundColor
	}
	background := canvas.NewRectangle(bgColor)
	background.CornerRadius = b.CornerRadius

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 0
	if b.BorderWidth > 0 && b.BorderColor != nil {
		border.StrokeColor = b.BorderColor
		border.StrokeWidth = b.BorderWidth
	}
	border.CornerRadius = b.CornerRadius

	// Initialize label with button's TintColor if set
	labelColor := theme.ForegroundColor()
	if b.TintColor != nil {
		labelColor = b.TintColor
	}
	label := canvas.NewText(b.Text, labelColor)
	label.TextStyle = fyne.TextStyle{}
	label.Alignment = fyne.TextAlignCenter

	subtitleLabel := canvas.NewText(b.Subtitle, theme.ForegroundColor())
	subtitleLabel.TextStyle = fyne.TextStyle{}
	subtitleLabel.Alignment = fyne.TextAlignCenter
	subtitleLabel.TextSize = theme.TextSize() - 2

	var iconImg *canvas.Image
	if b.Icon != nil {
		iconImg = canvas.NewImageFromResource(b.Icon)
		iconImg.FillMode = canvas.ImageFillContain
	}

	return &buttonRenderer{
		button:        b,
		background:    background,
		border:        border,
		label:         label,
		subtitleLabel: subtitleLabel,
		icon:          iconImg,
	}
}

// Tapped handles tap events
func (b *Button) Tapped(_ *fyne.PointEvent) {
	if !b.Enabled {
		return
	}
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// TappedSecondary handles secondary tap events
func (b *Button) TappedSecondary(_ *fyne.PointEvent) {}

// MouseIn handles mouse enter
func (b *Button) MouseIn(_ *desktop.MouseEvent) {
	b.mu.Lock()
	b.hovered = true
	b.mu.Unlock()
	b.Refresh()
}

// MouseMoved handles mouse movement
func (b *Button) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles mouse leave
func (b *Button) MouseOut() {
	b.mu.Lock()
	b.hovered = false
	b.mu.Unlock()
	b.Refresh()
}

// Cursor returns the cursor for this widget
func (b *Button) Cursor() desktop.Cursor {
	if b.Enabled {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

// MouseDown handles mouse button press - shows highlighted state
func (b *Button) MouseDown(_ *desktop.MouseEvent) {
	if !b.Enabled {
		return
	}
	b.mu.Lock()
	b.highlighted = true
	b.mu.Unlock()
	b.Refresh()
}

// MouseUp handles mouse button release - removes highlighted state
func (b *Button) MouseUp(_ *desktop.MouseEvent) {
	b.mu.Lock()
	b.highlighted = false
	b.mu.Unlock()
	b.Refresh()
}

type buttonRenderer struct {
	button        *Button
	background    *canvas.Rectangle
	border        *canvas.Rectangle
	label         *canvas.Text
	subtitleLabel *canvas.Text
	icon          *canvas.Image
}

func (r *buttonRenderer) Destroy() {}

func (r *buttonRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.background.Move(fyne.NewPos(0, 0))

	r.border.Resize(size)
	r.border.Move(fyne.NewPos(0, 0))

	insets := r.button.ContentEdgeInsets
	contentArea := fyne.NewSize(
		size.Width-insets.Left-insets.Right,
		size.Height-insets.Top-insets.Bottom,
	)
	contentStart := fyne.NewPos(insets.Left, insets.Top)

	// Calculate sizes
	var iconSize fyne.Size
	if r.icon != nil && r.button.Icon != nil {
		iconSize = fyne.NewSize(24, 24) // Default icon size
	}

	labelSize := fyne.NewSize(0, 0)
	if r.button.Text != "" {
		labelSize = r.label.MinSize()
	}

	subtitleSize := fyne.NewSize(0, 0)
	if r.button.Subtitle != "" {
		subtitleSize = r.subtitleLabel.MinSize()
	}

	spacing := r.button.SpacingBetweenImageAndTitle

	switch r.button.IconPosition {
	case ImagePositionLeft:
		totalWidth := iconSize.Width + labelSize.Width
		if iconSize.Width > 0 && labelSize.Width > 0 {
			totalWidth += spacing
		}
		startX := contentStart.X + (contentArea.Width-totalWidth)/2
		if r.icon != nil {
			r.icon.Move(fyne.NewPos(startX, contentStart.Y+(contentArea.Height-iconSize.Height)/2))
			r.icon.Resize(iconSize)
			startX += iconSize.Width + spacing
		}
		if r.button.Text != "" {
			r.label.Move(fyne.NewPos(startX, contentStart.Y+(contentArea.Height-labelSize.Height)/2))
		}

	case ImagePositionRight:
		totalWidth := iconSize.Width + labelSize.Width
		if iconSize.Width > 0 && labelSize.Width > 0 {
			totalWidth += spacing
		}
		startX := contentStart.X + (contentArea.Width-totalWidth)/2
		if r.button.Text != "" {
			r.label.Move(fyne.NewPos(startX, contentStart.Y+(contentArea.Height-labelSize.Height)/2))
			startX += labelSize.Width + spacing
		}
		if r.icon != nil {
			r.icon.Move(fyne.NewPos(startX, contentStart.Y+(contentArea.Height-iconSize.Height)/2))
			r.icon.Resize(iconSize)
		}

	case ImagePositionTop:
		totalHeight := iconSize.Height + labelSize.Height + subtitleSize.Height
		if iconSize.Height > 0 && labelSize.Height > 0 {
			totalHeight += spacing
		}
		startY := contentStart.Y + (contentArea.Height-totalHeight)/2
		if r.icon != nil {
			r.icon.Move(fyne.NewPos(contentStart.X+(contentArea.Width-iconSize.Width)/2, startY))
			r.icon.Resize(iconSize)
			startY += iconSize.Height + spacing
		}
		if r.button.Text != "" {
			r.label.Move(fyne.NewPos(contentStart.X+(contentArea.Width-labelSize.Width)/2, startY))
			startY += labelSize.Height
		}
		if r.button.Subtitle != "" {
			r.subtitleLabel.Move(fyne.NewPos(contentStart.X+(contentArea.Width-subtitleSize.Width)/2, startY))
		}

	case ImagePositionBottom:
		totalHeight := iconSize.Height + labelSize.Height + subtitleSize.Height
		if iconSize.Height > 0 && labelSize.Height > 0 {
			totalHeight += spacing
		}
		startY := contentStart.Y + (contentArea.Height-totalHeight)/2
		if r.button.Text != "" {
			r.label.Move(fyne.NewPos(contentStart.X+(contentArea.Width-labelSize.Width)/2, startY))
			startY += labelSize.Height
		}
		if r.button.Subtitle != "" {
			r.subtitleLabel.Move(fyne.NewPos(contentStart.X+(contentArea.Width-subtitleSize.Width)/2, startY))
			startY += subtitleSize.Height + spacing
		}
		if r.icon != nil {
			r.icon.Move(fyne.NewPos(contentStart.X+(contentArea.Width-iconSize.Width)/2, startY))
			r.icon.Resize(iconSize)
		}
	}
}

func (r *buttonRenderer) MinSize() fyne.Size {
	insets := r.button.ContentEdgeInsets

	var iconSize fyne.Size
	if r.icon != nil && r.button.Icon != nil {
		iconSize = fyne.NewSize(24, 24)
	}

	labelSize := fyne.NewSize(0, 0)
	if r.button.Text != "" {
		labelSize = r.label.MinSize()
	}

	subtitleSize := fyne.NewSize(0, 0)
	if r.button.Subtitle != "" {
		subtitleSize = r.subtitleLabel.MinSize()
	}

	spacing := r.button.SpacingBetweenImageAndTitle

	var contentWidth, contentHeight float32

	switch r.button.IconPosition {
	case ImagePositionLeft, ImagePositionRight:
		contentWidth = iconSize.Width + labelSize.Width
		if iconSize.Width > 0 && labelSize.Width > 0 {
			contentWidth += spacing
		}
		contentHeight = max(iconSize.Height, labelSize.Height+subtitleSize.Height)

	case ImagePositionTop, ImagePositionBottom:
		contentWidth = max(iconSize.Width, max(labelSize.Width, subtitleSize.Width))
		contentHeight = iconSize.Height + labelSize.Height + subtitleSize.Height
		if iconSize.Height > 0 && labelSize.Height > 0 {
			contentHeight += spacing
		}
	}

	return fyne.NewSize(
		contentWidth+insets.Left+insets.Right,
		contentHeight+insets.Top+insets.Bottom,
	)
}

func (r *buttonRenderer) Refresh() {
	config := core.SharedConfiguration()

	// Update background
	if r.button.BackgroundColor != nil {
		r.background.FillColor = r.button.BackgroundColor
	} else {
		r.background.FillColor = color.Transparent
	}

	// Apply corner radius
	cornerRadius := r.button.CornerRadius
	if cornerRadius == CornerRadiusAdjustsBounds {
		cornerRadius = r.button.Size().Height / 2
	}
	r.background.CornerRadius = cornerRadius
	r.border.CornerRadius = cornerRadius

	// Update border
	if r.button.BorderWidth > 0 && r.button.BorderColor != nil {
		r.border.StrokeWidth = r.button.BorderWidth
		r.border.StrokeColor = r.button.BorderColor
		r.border.FillColor = color.Transparent
	}

	// Handle states
	r.button.mu.RLock()
	hovered := r.button.hovered
	enabled := r.button.Enabled
	highlighted := r.button.highlighted
	r.button.mu.RUnlock()

	alpha := 1.0
	if !enabled && r.button.AdjustsButtonWhenDisabled {
		alpha = config.ButtonDisabledAlpha
	} else if (hovered || highlighted) && r.button.AdjustsButtonWhenHighlighted {
		alpha = config.ButtonHighlightedAlpha
		if r.button.HighlightedBackgroundColor != nil {
			r.background.FillColor = r.button.HighlightedBackgroundColor
		}
		if r.button.HighlightedBorderColor != nil {
			r.border.StrokeColor = r.button.HighlightedBorderColor
		}
	}

	// Update label color
	textColor := theme.ForegroundColor()
	if r.button.AdjustsTitleTintColorAutomatically && r.button.TintColor != nil {
		textColor = r.button.TintColor
	}
	if !enabled && r.button.DisabledColor != nil {
		textColor = r.button.DisabledColor
	}
	textColor = core.ColorWithAlpha(textColor, alpha)
	r.label.Color = textColor
	r.label.Text = r.button.Text

	// Update subtitle
	if r.button.Subtitle != "" {
		subtitleColor := textColor
		if r.button.SubtitleColor != nil {
			subtitleColor = r.button.SubtitleColor
		}
		r.subtitleLabel.Color = core.ColorWithAlpha(subtitleColor, alpha)
		r.subtitleLabel.Text = r.button.Subtitle
	}

	// Update icon
	if r.icon != nil && r.button.Icon != nil {
		r.icon.Resource = r.button.Icon
		r.icon.Refresh()
	}

	r.background.Refresh()
	r.border.Refresh()
	r.label.Refresh()
	r.subtitleLabel.Refresh()
}

func (r *buttonRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background, r.border, r.label}
	if r.button.Subtitle != "" {
		objects = append(objects, r.subtitleLabel)
	}
	if r.icon != nil {
		objects = append(objects, r.icon)
	}
	return objects
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// NavigationButton is a button styled for navigation bars
type NavigationButton struct {
	*Button
	IsBackButton bool
}

// NewNavigationButton creates a navigation bar button
func NewNavigationButton(text string, tapped func()) *NavigationButton {
	btn := NewButton(text, tapped)
	btn.ContentEdgeInsets = core.NewEdgeInsets(4, 8, 4, 8)
	config := core.SharedConfiguration()
	btn.TintColor = config.NavBarTintColor
	btn.AdjustsTitleTintColorAutomatically = true
	return &NavigationButton{Button: btn}
}

// NewNavigationBackButton creates a back button for navigation
func NewNavigationBackButton(tapped func()) *NavigationButton {
	btn := NewNavigationButton("Back", tapped)
	btn.IsBackButton = true
	return btn
}

// ToolbarButton is a button styled for toolbars
type ToolbarButton struct {
	*Button
}

// NewToolbarButton creates a toolbar button
func NewToolbarButton(text string, tapped func()) *ToolbarButton {
	btn := NewButton(text, tapped)
	btn.ContentEdgeInsets = core.NewEdgeInsets(4, 12, 4, 12)
	config := core.SharedConfiguration()
	btn.TintColor = config.ToolBarTintColor
	btn.AdjustsTitleTintColorAutomatically = true
	return &ToolbarButton{Button: btn}
}

// NewToolbarIconButton creates an icon-only toolbar button
func NewToolbarIconButton(icon fyne.Resource, tapped func()) *ToolbarButton {
	btn := NewIconButton(icon, tapped)
	config := core.SharedConfiguration()
	btn.TintColor = config.ToolBarTintColor
	btn.AdjustsImageTintColorAutomatically = true
	return &ToolbarButton{Button: btn}
}

// FillButton is a solid-filled button
type FillButton struct {
	*Button
}

// NewFillButton creates a filled button with background color
func NewFillButton(text string, fillColor color.Color, tapped func()) *FillButton {
	btn := NewButton(text, tapped)
	btn.BackgroundColor = fillColor
	btn.CornerRadius = 8
	btn.TintColor = color.White
	btn.AdjustsTitleTintColorAutomatically = true
	btn.HighlightedBackgroundColor = core.ColorWithAlpha(fillColor, 0.7)
	return &FillButton{Button: btn}
}

// GhostButton is an outlined button without fill
type GhostButton struct {
	*Button
}

// NewGhostButton creates an outlined ghost button
func NewGhostButton(text string, borderColor color.Color, tapped func()) *GhostButton {
	btn := NewButton(text, tapped)
	btn.BorderColor = borderColor
	btn.BorderWidth = 1
	btn.CornerRadius = 8
	btn.TintColor = borderColor
	btn.AdjustsTitleTintColorAutomatically = true
	btn.HighlightedBorderColor = core.ColorWithAlpha(borderColor, 0.5)
	return &GhostButton{Button: btn}
}

// ApplyTheme implements the Themeable interface for Button
func (b *Button) ApplyTheme(t *qmuitheme.Theme) {
	b.TintColor = t.PrimaryColor
	b.Refresh()
}

// ApplyTheme implements the Themeable interface for FillButton
func (fb *FillButton) ApplyTheme(t *qmuitheme.Theme) {
	fb.BackgroundColor = t.ButtonBackgroundColor
	fb.HighlightedBackgroundColor = core.ColorWithAlpha(t.ButtonBackgroundColor, 0.7)
	fb.TintColor = t.ButtonTextColor
	fb.Refresh()
}

// ApplyTheme implements the Themeable interface for GhostButton
func (gb *GhostButton) ApplyTheme(t *qmuitheme.Theme) {
	gb.BorderColor = t.PrimaryColor
	gb.TintColor = t.PrimaryColor
	gb.HighlightedBorderColor = core.ColorWithAlpha(t.PrimaryColor, 0.5)
	gb.Refresh()
}
