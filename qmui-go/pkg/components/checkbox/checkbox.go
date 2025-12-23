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

	// Create the checkbox visual
	box := canvas.NewRectangle(color.Transparent)
	box.StrokeWidth = 2

	checkmark := canvas.NewRectangle(color.Transparent)
	indeterminateMark := canvas.NewRectangle(color.Transparent)

	label := canvas.NewText(c.Text, c.TextColor)
	label.TextSize = c.TextSize

	return &checkboxRenderer{
		checkbox:          c,
		box:               box,
		checkmark:         checkmark,
		indeterminateMark: indeterminateMark,
		label:             label,
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
	checkbox          *Checkbox
	box               *canvas.Rectangle
	checkmark         *canvas.Rectangle
	indeterminateMark *canvas.Rectangle
	label             *canvas.Text
}

func (r *checkboxRenderer) Destroy() {}

func (r *checkboxRenderer) Layout(size fyne.Size) {
	checkSize := r.checkbox.CheckboxSize

	// Position checkbox
	r.box.Resize(checkSize)
	r.box.Move(fyne.NewPos(0, (size.Height-checkSize.Height)/2))

	// Position checkmark (smaller rectangle inside)
	markSize := fyne.NewSize(checkSize.Width*0.5, checkSize.Height*0.5)
	markPos := fyne.NewPos(
		(checkSize.Width-markSize.Width)/2,
		(size.Height-markSize.Height)/2,
	)
	r.checkmark.Resize(markSize)
	r.checkmark.Move(markPos)

	// Position indeterminate mark (horizontal bar)
	indSize := fyne.NewSize(checkSize.Width*0.6, checkSize.Height*0.2)
	indPos := fyne.NewPos(
		(checkSize.Width-indSize.Width)/2,
		(size.Height-indSize.Height)/2,
	)
	r.indeterminateMark.Resize(indSize)
	r.indeterminateMark.Move(indPos)

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

	// Update box appearance
	if !enabled {
		r.box.StrokeColor = core.ColorWithAlpha(tintColor, config.ControlDisabledAlpha)
		r.box.FillColor = color.Transparent
	} else if selected || indeterminate {
		r.box.StrokeColor = tintColor
		r.box.FillColor = tintColor
	} else {
		r.box.StrokeColor = tintColor
		r.box.FillColor = color.Transparent
	}

	if hovered && enabled {
		r.box.StrokeColor = core.ColorWithAlpha(tintColor, 0.7)
	}

	// Make it circular
	r.box.CornerRadius = r.checkbox.CheckboxSize.Width / 2

	// Show/hide checkmark
	if selected && !indeterminate {
		r.checkmark.FillColor = color.White
		r.checkmark.Show()
	} else {
		r.checkmark.Hide()
	}

	// Show/hide indeterminate mark
	if indeterminate {
		r.indeterminateMark.FillColor = color.White
		r.indeterminateMark.Show()
	} else {
		r.indeterminateMark.Hide()
	}

	// Update label
	r.label.Text = r.checkbox.Text
	r.label.Color = r.checkbox.TextColor
	r.label.TextSize = r.checkbox.TextSize
	if !enabled {
		r.label.Color = core.ColorWithAlpha(r.checkbox.TextColor, config.ControlDisabledAlpha)
	}

	r.box.Refresh()
	r.checkmark.Refresh()
	r.indeterminateMark.Refresh()
	r.label.Refresh()
}

func (r *checkboxRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.box, r.checkmark, r.indeterminateMark}
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
