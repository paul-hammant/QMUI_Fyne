// Package checkbox provides QMUICheckbox - a circular checkbox control
// Ported from Tencent's QMUI_iOS framework
package checkbox

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

// Checkbox is a circular checkbox control with three states:
// - Unchecked (Selected = false, Indeterminate = false)
// - Checked (Selected = true, Indeterminate = false)
// - Indeterminate (Indeterminate = true, Selected = false)
type Checkbox struct {
	widget.BaseWidget

	// State
	Selected      bool
	Indeterminate bool
	Enabled       bool

	// Styling
	TintColor     color.Color
	CheckboxSize  fyne.Size
	NormalImage   fyne.Resource
	SelectedImage fyne.Resource
	IndeterminateImage fyne.Resource
	DisabledImage fyne.Resource

	// Optional label
	Text          string
	TextColor     color.Color
	TextSize      float32
	SpacingBetweenCheckboxAndText float32

	// Callbacks
	OnChanged func(selected bool)

	mu      sync.RWMutex
	hovered bool
}

// NewCheckbox creates a new checkbox
func NewCheckbox(onChanged func(selected bool)) *Checkbox {
	c := &Checkbox{
		Selected:      false,
		Indeterminate: false,
		Enabled:       true,
		TintColor:     core.SharedConfiguration().BlueColor,
		CheckboxSize:  fyne.NewSize(16, 16),
		SpacingBetweenCheckboxAndText: 8,
		TextSize:      theme.TextSize(),
		OnChanged:     onChanged,
	}
	c.ExtendBaseWidget(c)
	return c
}

// NewCheckboxWithLabel creates a checkbox with a text label
func NewCheckboxWithLabel(text string, onChanged func(selected bool)) *Checkbox {
	c := NewCheckbox(onChanged)
	c.Text = text
	c.TextColor = theme.ForegroundColor()
	return c
}

// SetSelected sets the selected state
func (c *Checkbox) SetSelected(selected bool) {
	c.mu.Lock()
	c.Selected = selected
	if selected {
		c.Indeterminate = false
	}
	c.mu.Unlock()
	c.Refresh()
	if c.OnChanged != nil {
		c.OnChanged(selected)
	}
}

// SetIndeterminate sets the indeterminate state
func (c *Checkbox) SetIndeterminate(indeterminate bool) {
	c.mu.Lock()
	c.Indeterminate = indeterminate
	if indeterminate {
		c.Selected = false
	}
	c.mu.Unlock()
	c.Refresh()
}

// Toggle toggles the checkbox state
func (c *Checkbox) Toggle() {
	c.mu.Lock()
	if c.Indeterminate {
		c.Indeterminate = false
		c.Selected = true
	} else {
		c.Selected = !c.Selected
	}
	selected := c.Selected
	c.mu.Unlock()
	c.Refresh()
	if c.OnChanged != nil {
		c.OnChanged(selected)
	}
}

// CreateRenderer implements fyne.Widget
func (c *Checkbox) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	// Create the checkbox visual - use circle for iOS-style
	circle := canvas.NewCircle(color.Transparent)
	circle.StrokeWidth = 2

	label := canvas.NewText(c.Text, c.TextColor)
	label.TextSize = c.TextSize

	return &checkboxRenderer{
		checkbox: c,
		circle:   circle,
		label:    label,
	}
}

// Tapped handles tap events
func (c *Checkbox) Tapped(_ *fyne.PointEvent) {
	if !c.Enabled {
		return
	}
	c.Toggle()
}

// TappedSecondary handles secondary tap
func (c *Checkbox) TappedSecondary(_ *fyne.PointEvent) {}

// MouseIn handles mouse enter
func (c *Checkbox) MouseIn(_ *desktop.MouseEvent) {
	c.mu.Lock()
	c.hovered = true
	c.mu.Unlock()
	c.Refresh()
}

// MouseMoved handles mouse movement
func (c *Checkbox) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles mouse leave
func (c *Checkbox) MouseOut() {
	c.mu.Lock()
	c.hovered = false
	c.mu.Unlock()
	c.Refresh()
}

// Cursor returns the cursor for this widget
func (c *Checkbox) Cursor() desktop.Cursor {
	if c.Enabled {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

type checkboxRenderer struct {
	checkbox *Checkbox
	circle   *canvas.Circle
	// Checkmark drawn with two lines (tick shape)
	checkLine1 *canvas.Line
	checkLine2 *canvas.Line
	// Indeterminate drawn with horizontal line
	indeterminateLine *canvas.Line
	label             *canvas.Text
}

func (r *checkboxRenderer) Destroy() {}

func (r *checkboxRenderer) Layout(size fyne.Size) {
	checkSize := r.checkbox.CheckboxSize
	centerY := (size.Height - checkSize.Height) / 2

	// Position circle
	r.circle.Resize(checkSize)
	r.circle.Move(fyne.NewPos(0, centerY))

	// Calculate checkmark positions (tick shape: short line down-right, long line up-right)
	cx := checkSize.Width / 2
	cy := centerY + checkSize.Height/2

	// Checkmark proportions for a nice tick
	// Start point (left), middle point (bottom of tick), end point (top right)
	startX := cx - checkSize.Width*0.25
	startY := cy
	midX := cx - checkSize.Width*0.05
	midY := cy + checkSize.Height*0.2
	endX := cx + checkSize.Width*0.3
	endY := cy - checkSize.Height*0.25

	// Initialize lines if needed
	if r.checkLine1 == nil {
		r.checkLine1 = canvas.NewLine(color.White)
		r.checkLine2 = canvas.NewLine(color.White)
		r.indeterminateLine = canvas.NewLine(color.White)
	}

	r.checkLine1.Position1 = fyne.NewPos(startX, startY)
	r.checkLine1.Position2 = fyne.NewPos(midX, midY)
	r.checkLine1.StrokeWidth = 2

	r.checkLine2.Position1 = fyne.NewPos(midX, midY)
	r.checkLine2.Position2 = fyne.NewPos(endX, endY)
	r.checkLine2.StrokeWidth = 2

	// Indeterminate line (horizontal bar in center)
	r.indeterminateLine.Position1 = fyne.NewPos(cx-checkSize.Width*0.25, cy)
	r.indeterminateLine.Position2 = fyne.NewPos(cx+checkSize.Width*0.25, cy)
	r.indeterminateLine.StrokeWidth = 2

	// Position label
	if r.checkbox.Text != "" {
		labelX := checkSize.Width + r.checkbox.SpacingBetweenCheckboxAndText
		r.label.Move(fyne.NewPos(labelX, (size.Height-r.label.MinSize().Height)/2))
	}
}

func (r *checkboxRenderer) MinSize() fyne.Size {
	checkSize := r.checkbox.CheckboxSize
	width := checkSize.Width
	height := checkSize.Height

	if r.checkbox.Text != "" {
		labelSize := r.label.MinSize()
		width += r.checkbox.SpacingBetweenCheckboxAndText + labelSize.Width
		if labelSize.Height > height {
			height = labelSize.Height
		}
	}

	return fyne.NewSize(width, height)
}

func (r *checkboxRenderer) Refresh() {
	r.checkbox.mu.RLock()
	selected := r.checkbox.Selected
	indeterminate := r.checkbox.Indeterminate
	enabled := r.checkbox.Enabled
	hovered := r.checkbox.hovered
	r.checkbox.mu.RUnlock()

	tintColor := r.checkbox.TintColor
	config := core.SharedConfiguration()

	// Initialize lines if needed
	if r.checkLine1 == nil {
		r.checkLine1 = canvas.NewLine(color.White)
		r.checkLine2 = canvas.NewLine(color.White)
		r.indeterminateLine = canvas.NewLine(color.White)
	}

	// Update circle appearance - iOS style: stroke only for normal, filled for selected
	if !enabled {
		r.circle.StrokeColor = core.ColorWithAlpha(tintColor, config.ControlDisabledAlpha)
		r.circle.FillColor = color.Transparent
	} else if selected || indeterminate {
		// Filled circle when selected/indeterminate
		r.circle.StrokeColor = tintColor
		r.circle.FillColor = tintColor
	} else {
		// Outline only when not selected (iOS style)
		r.circle.StrokeColor = tintColor
		r.circle.FillColor = color.Transparent
	}

	if hovered && enabled && !selected && !indeterminate {
		r.circle.StrokeColor = core.ColorWithAlpha(tintColor, 0.7)
	}

	r.circle.StrokeWidth = 2

	// Show/hide checkmark lines
	if selected && !indeterminate {
		r.checkLine1.StrokeColor = color.White
		r.checkLine2.StrokeColor = color.White
		r.checkLine1.Show()
		r.checkLine2.Show()
	} else {
		r.checkLine1.Hide()
		r.checkLine2.Hide()
	}

	// Show/hide indeterminate line
	if indeterminate {
		r.indeterminateLine.StrokeColor = color.White
		r.indeterminateLine.Show()
	} else {
		r.indeterminateLine.Hide()
	}

	// Update label
	r.label.Text = r.checkbox.Text
	r.label.Color = r.checkbox.TextColor
	r.label.TextSize = r.checkbox.TextSize
	if !enabled {
		r.label.Color = core.ColorWithAlpha(r.checkbox.TextColor, config.ControlDisabledAlpha)
	}

	r.circle.Refresh()
	if r.checkLine1 != nil {
		r.checkLine1.Refresh()
		r.checkLine2.Refresh()
		r.indeterminateLine.Refresh()
	}
	r.label.Refresh()
}

func (r *checkboxRenderer) Objects() []fyne.CanvasObject {
	// Initialize lines if needed
	if r.checkLine1 == nil {
		r.checkLine1 = canvas.NewLine(color.White)
		r.checkLine2 = canvas.NewLine(color.White)
		r.indeterminateLine = canvas.NewLine(color.White)
	}
	objects := []fyne.CanvasObject{r.circle, r.checkLine1, r.checkLine2, r.indeterminateLine}
	if r.checkbox.Text != "" {
		objects = append(objects, r.label)
	}
	return objects
}

// CheckboxGroup manages a group of checkboxes
type CheckboxGroup struct {
	widget.BaseWidget

	Checkboxes []*Checkbox
	Orientation string // "horizontal" or "vertical"
	Spacing     float32
}

// NewCheckboxGroup creates a group of checkboxes
func NewCheckboxGroup(options []string, onChanged func(selected []string)) *CheckboxGroup {
	group := &CheckboxGroup{
		Orientation: "vertical",
		Spacing:     8,
	}

	for _, opt := range options {
		checkbox := NewCheckboxWithLabel(opt, func(selected bool) {
			if onChanged != nil {
				var selectedOpts []string
				for _, cb := range group.Checkboxes {
					if cb.Selected {
						selectedOpts = append(selectedOpts, cb.Text)
					}
				}
				onChanged(selectedOpts)
			}
		})
		group.Checkboxes = append(group.Checkboxes, checkbox)
	}

	group.ExtendBaseWidget(group)
	return group
}

// SetSelected sets the selected options
func (g *CheckboxGroup) SetSelected(options []string) {
	optSet := make(map[string]bool)
	for _, opt := range options {
		optSet[opt] = true
	}

	for _, cb := range g.Checkboxes {
		cb.SetSelected(optSet[cb.Text])
	}
}

// GetSelected returns all selected options
func (g *CheckboxGroup) GetSelected() []string {
	var selected []string
	for _, cb := range g.Checkboxes {
		if cb.Selected {
			selected = append(selected, cb.Text)
		}
	}
	return selected
}

// CreateRenderer implements fyne.Widget
func (g *CheckboxGroup) CreateRenderer() fyne.WidgetRenderer {
	g.ExtendBaseWidget(g)

	objects := make([]fyne.CanvasObject, len(g.Checkboxes))
	for i, cb := range g.Checkboxes {
		objects[i] = cb
	}

	return &checkboxGroupRenderer{
		group:   g,
		objects: objects,
	}
}

type checkboxGroupRenderer struct {
	group   *CheckboxGroup
	objects []fyne.CanvasObject
}

func (r *checkboxGroupRenderer) Destroy() {}

func (r *checkboxGroupRenderer) Layout(size fyne.Size) {
	var x, y float32

	for _, obj := range r.objects {
		objSize := obj.MinSize()
		obj.Resize(objSize)
		obj.Move(fyne.NewPos(x, y))

		if r.group.Orientation == "horizontal" {
			x += objSize.Width + r.group.Spacing
		} else {
			y += objSize.Height + r.group.Spacing
		}
	}
}

func (r *checkboxGroupRenderer) MinSize() fyne.Size {
	var width, height float32

	for i, obj := range r.objects {
		objSize := obj.MinSize()

		if r.group.Orientation == "horizontal" {
			width += objSize.Width
			if i < len(r.objects)-1 {
				width += r.group.Spacing
			}
			if objSize.Height > height {
				height = objSize.Height
			}
		} else {
			height += objSize.Height
			if i < len(r.objects)-1 {
				height += r.group.Spacing
			}
			if objSize.Width > width {
				width = objSize.Width
			}
		}
	}

	return fyne.NewSize(width, height)
}

func (r *checkboxGroupRenderer) Refresh() {
	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *checkboxGroupRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
