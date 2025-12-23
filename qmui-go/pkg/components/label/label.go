// Package label provides QMUILabel - an enhanced label with padding and copy support
// Ported from Tencent's QMUI_iOS framework
package label

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

// Label is an enhanced label widget with QMUI styling features
type Label struct {
	widget.BaseWidget

	// Text content
	Text      string
	TextStyle fyne.TextStyle
	TextSize  float32
	Alignment fyne.TextAlign
	Wrapping  fyne.TextWrap
	Truncation fyne.TextTruncation

	// Styling
	Color                    color.Color
	ContentEdgeInsets        core.EdgeInsets
	HighlightedBackgroundColor color.Color

	// Copy feature
	CanPerformCopyAction   bool
	MenuItemTitleForCopyAction string
	OnCopy                 func(label *Label, copiedText string)

	// Truncating tail view (custom view shown when text truncates)
	TruncatingTailView fyne.CanvasObject

	// State
	mu          sync.RWMutex
	hovered     bool
	highlighted bool
	longPressed bool
}

// NewLabel creates a new QMUI-styled label
func NewLabel(text string) *Label {
	l := &Label{
		Text:              text,
		TextStyle:         fyne.TextStyle{},
		TextSize:          theme.TextSize(),
		Alignment:         fyne.TextAlignLeading,
		Wrapping:          fyne.TextWrapOff,
		Color:             theme.ForegroundColor(),
		ContentEdgeInsets: core.EdgeInsets{},
		CanPerformCopyAction: false,
	}
	l.ExtendBaseWidget(l)
	return l
}

// NewLabelWithStyle creates a new QMUI-styled label with custom text style
func NewLabelWithStyle(text string, alignment fyne.TextAlign, style fyne.TextStyle) *Label {
	l := NewLabel(text)
	l.Alignment = alignment
	l.TextStyle = style
	return l
}

// SetText updates the label text
func (l *Label) SetText(text string) {
	l.mu.Lock()
	l.Text = text
	l.mu.Unlock()
	l.Refresh()
}

// SetHighlighted sets the highlighted state
func (l *Label) SetHighlighted(highlighted bool) {
	l.mu.Lock()
	l.highlighted = highlighted
	l.mu.Unlock()
	l.Refresh()
}

// CreateRenderer implements fyne.Widget
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)

	background := canvas.NewRectangle(color.Transparent)
	text := canvas.NewText(l.Text, l.Color)
	text.TextStyle = l.TextStyle
	text.TextSize = l.TextSize
	text.Alignment = l.Alignment

	return &labelRenderer{
		label:      l,
		background: background,
		text:       text,
	}
}

// Tapped handles tap events
func (l *Label) Tapped(_ *fyne.PointEvent) {}

// TappedSecondary handles secondary tap (right-click for copy menu)
func (l *Label) TappedSecondary(e *fyne.PointEvent) {
	if l.CanPerformCopyAction {
		l.copyText()
	}
}

// MouseIn handles mouse enter
func (l *Label) MouseIn(_ *desktop.MouseEvent) {
	l.mu.Lock()
	l.hovered = true
	l.mu.Unlock()
}

// MouseMoved handles mouse movement
func (l *Label) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles mouse leave
func (l *Label) MouseOut() {
	l.mu.Lock()
	l.hovered = false
	l.mu.Unlock()
}

func (l *Label) copyText() {
	if clipboard := fyne.CurrentApp().Driver().AllWindows(); len(clipboard) > 0 {
		clipboard[0].Clipboard().SetContent(l.Text)
		if l.OnCopy != nil {
			l.OnCopy(l, l.Text)
		}
	}
}

// Cursor returns the cursor for this widget
func (l *Label) Cursor() desktop.Cursor {
	if l.CanPerformCopyAction {
		return desktop.TextCursor
	}
	return desktop.DefaultCursor
}

type labelRenderer struct {
	label      *Label
	background *canvas.Rectangle
	text       *canvas.Text
}

func (r *labelRenderer) Destroy() {}

func (r *labelRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.background.Move(fyne.NewPos(0, 0))

	insets := r.label.ContentEdgeInsets
	textPos := fyne.NewPos(insets.Left, insets.Top)
	textSize := fyne.NewSize(
		size.Width-insets.Left-insets.Right,
		size.Height-insets.Top-insets.Bottom,
	)

	r.text.Move(textPos)
	r.text.Resize(textSize)
}

func (r *labelRenderer) MinSize() fyne.Size {
	textSize := r.text.MinSize()
	insets := r.label.ContentEdgeInsets
	return fyne.NewSize(
		textSize.Width+insets.Left+insets.Right,
		textSize.Height+insets.Top+insets.Bottom,
	)
}

func (r *labelRenderer) Refresh() {
	r.label.mu.RLock()
	highlighted := r.label.highlighted
	r.label.mu.RUnlock()

	// Update background
	if highlighted && r.label.HighlightedBackgroundColor != nil {
		r.background.FillColor = r.label.HighlightedBackgroundColor
	} else {
		r.background.FillColor = color.Transparent
	}

	// Update text
	r.text.Text = r.label.Text
	r.text.Color = r.label.Color
	r.text.TextStyle = r.label.TextStyle
	r.text.TextSize = r.label.TextSize
	r.text.Alignment = r.label.Alignment

	r.background.Refresh()
	r.text.Refresh()
}

func (r *labelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.text}
}

// RichLabel is a label with attributed text support
type RichLabel struct {
	widget.BaseWidget

	Segments []RichTextSegment
	ContentEdgeInsets core.EdgeInsets
}

// RichTextSegment represents a segment of styled text
type RichTextSegment struct {
	Text      string
	Style     fyne.TextStyle
	Color     color.Color
	Size      float32
	Inline    bool
}

// NewRichLabel creates a new rich text label
func NewRichLabel() *RichLabel {
	rl := &RichLabel{
		Segments: []RichTextSegment{},
	}
	rl.ExtendBaseWidget(rl)
	return rl
}

// AppendText adds a text segment
func (rl *RichLabel) AppendText(text string, style fyne.TextStyle, textColor color.Color, size float32) {
	rl.Segments = append(rl.Segments, RichTextSegment{
		Text:   text,
		Style:  style,
		Color:  textColor,
		Size:   size,
		Inline: true,
	})
	rl.Refresh()
}

// Clear removes all segments
func (rl *RichLabel) Clear() {
	rl.Segments = []RichTextSegment{}
	rl.Refresh()
}

// CreateRenderer implements fyne.Widget
func (rl *RichLabel) CreateRenderer() fyne.WidgetRenderer {
	rl.ExtendBaseWidget(rl)
	return &richLabelRenderer{
		label: rl,
		texts: make([]*canvas.Text, 0),
	}
}

type richLabelRenderer struct {
	label *RichLabel
	texts []*canvas.Text
}

func (r *richLabelRenderer) Destroy() {}

func (r *richLabelRenderer) Layout(size fyne.Size) {
	insets := r.label.ContentEdgeInsets
	x := insets.Left
	y := insets.Top

	for _, text := range r.texts {
		textSize := text.MinSize()
		text.Move(fyne.NewPos(x, y))
		x += textSize.Width
		if x > size.Width-insets.Right {
			x = insets.Left
			y += textSize.Height
		}
	}
}

func (r *richLabelRenderer) MinSize() fyne.Size {
	var width, height float32
	var lineHeight float32

	for _, text := range r.texts {
		textSize := text.MinSize()
		width += textSize.Width
		if textSize.Height > lineHeight {
			lineHeight = textSize.Height
		}
	}

	height = lineHeight
	insets := r.label.ContentEdgeInsets
	return fyne.NewSize(width+insets.Left+insets.Right, height+insets.Top+insets.Bottom)
}

func (r *richLabelRenderer) Refresh() {
	// Rebuild text objects
	r.texts = make([]*canvas.Text, 0, len(r.label.Segments))
	for _, seg := range r.label.Segments {
		text := canvas.NewText(seg.Text, seg.Color)
		text.TextStyle = seg.Style
		if seg.Size > 0 {
			text.TextSize = seg.Size
		}
		r.texts = append(r.texts, text)
	}
}

func (r *richLabelRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(r.texts))
	for i, t := range r.texts {
		objects[i] = t
	}
	return objects
}
