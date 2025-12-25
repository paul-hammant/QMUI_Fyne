// Package textview provides QMUITextView - an enhanced multi-line text input
// Ported from Tencent's QMUI_iOS framework
package textview

import (
	"image/color"
	"sync"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// TextViewDelegate provides callbacks for text view events
type TextViewDelegate interface {
	// NewHeightAfterTextChanged is called when content height changes
	NewHeightAfterTextChanged(textView *TextView, height float32)
	// ShouldReturn is called when return key is pressed
	ShouldReturn(textView *TextView) bool
	// ShouldChangeTextInRange is called before text changes
	ShouldChangeTextInRange(textView *TextView, start, length int, replacement string, originalValue bool) bool
	// DidPreventTextChange is called when text change was prevented
	DidPreventTextChange(textView *TextView, start, length int, replacement string)
}

// TextView is an enhanced multi-line text input widget
type TextView struct {
	widget.Entry

	// Styling
	Placeholder        string
	PlaceholderColor   color.Color
	PlaceholderMargins core.EdgeInsets

	// Behavior
	ShouldResponseToProgrammaticallyTextChanges bool
	MaximumTextLength                           int
	ShouldCountingNonASCIICharacterAsTwo        bool
	MaximumHeight                               float32
	IsDeletingDuringTextChange                  bool

	// Delegate
	Delegate TextViewDelegate

	// Callbacks
	OnTextChanged      func(text string)
	OnHeightChanged    func(newHeight float32)
	OnPaste            func(sender interface{}) bool

	mu            sync.RWMutex
	lastHeight    float32
}

// NewTextView creates a new QMUI-styled text view
func NewTextView() *TextView {
	config := core.SharedConfiguration()
	tv := &TextView{
		PlaceholderColor: config.PlaceholderColor,
		ShouldResponseToProgrammaticallyTextChanges: true,
		MaximumTextLength:    -1, // No limit
		MaximumHeight:        0,  // No limit
		ShouldCountingNonASCIICharacterAsTwo: false,
	}
	tv.MultiLine = true
	tv.Wrapping = fyne.TextWrapWord
	tv.ExtendBaseWidget(tv)
	tv.Entry.OnChanged = tv.handleTextChanged
	return tv
}

// NewTextViewWithPlaceholder creates a text view with placeholder
func NewTextViewWithPlaceholder(placeholder string) *TextView {
	tv := NewTextView()
	tv.Placeholder = placeholder
	tv.PlaceHolder = placeholder
	return tv
}

// SetText sets the text and optionally triggers change events
func (tv *TextView) SetText(text string) {
	tv.mu.Lock()
	shouldNotify := tv.ShouldResponseToProgrammaticallyTextChanges
	tv.mu.Unlock()

	tv.Entry.SetText(text)
	if shouldNotify && tv.OnTextChanged != nil {
		tv.OnTextChanged(text)
	}
}

// handleTextChanged processes text changes with length limiting
func (tv *TextView) handleTextChanged(text string) {
	tv.mu.RLock()
	maxLen := tv.MaximumTextLength
	countNonASCII := tv.ShouldCountingNonASCIICharacterAsTwo
	tv.mu.RUnlock()

	if maxLen > 0 {
		length := tv.calculateTextLength(text, countNonASCII)
		if length > maxLen {
			// Trim text to max length
			trimmed := tv.trimToLength(text, maxLen, countNonASCII)
			tv.Entry.SetText(trimmed)
			if tv.Delegate != nil {
				tv.Delegate.DidPreventTextChange(tv, len(trimmed), len(text)-len(trimmed), "")
			}
			text = trimmed
		}
	}

	if tv.OnTextChanged != nil {
		tv.OnTextChanged(text)
	}

	// Check height change
	tv.checkHeightChange()
}

func (tv *TextView) calculateTextLength(s string, countNonASCIIAsTwo bool) int {
	if !countNonASCIIAsTwo {
		return utf8.RuneCountInString(s)
	}
	count := 0
	for _, r := range s {
		if r > 127 {
			count += 2
		} else {
			count++
		}
	}
	return count
}

func (tv *TextView) trimToLength(s string, maxLen int, countNonASCIIAsTwo bool) string {
	if !countNonASCIIAsTwo {
		runes := []rune(s)
		if len(runes) > maxLen {
			return string(runes[:maxLen])
		}
		return s
	}

	var result []rune
	count := 0
	for _, r := range s {
		charLen := 1
		if r > 127 {
			charLen = 2
		}
		if count+charLen > maxLen {
			break
		}
		result = append(result, r)
		count += charLen
	}
	return string(result)
}

func (tv *TextView) checkHeightChange() {
	currentHeight := tv.MinSize().Height

	tv.mu.Lock()
	maxHeight := tv.MaximumHeight
	lastHeight := tv.lastHeight
	tv.mu.Unlock()

	if maxHeight > 0 && currentHeight > maxHeight {
		currentHeight = maxHeight
	}

	if currentHeight != lastHeight {
		tv.mu.Lock()
		tv.lastHeight = currentHeight
		tv.mu.Unlock()

		if tv.OnHeightChanged != nil {
			tv.OnHeightChanged(currentHeight)
		}
		if tv.Delegate != nil {
			tv.Delegate.NewHeightAfterTextChanged(tv, currentHeight)
		}
	}
}

// CreateRenderer implements fyne.Widget
func (tv *TextView) CreateRenderer() fyne.WidgetRenderer {
	tv.ExtendBaseWidget(tv)

	background := canvas.NewRectangle(color.Transparent)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 1
	border.StrokeColor = core.SharedConfiguration().SeparatorColor

	placeholder := canvas.NewText(tv.Placeholder, tv.PlaceholderColor)

	entryRenderer := tv.Entry.CreateRenderer()

	return &textViewRenderer{
		textView:      tv,
		background:    background,
		border:        border,
		placeholder:   placeholder,
		entryRenderer: entryRenderer,
	}
}

type textViewRenderer struct {
	textView      *TextView
	background    *canvas.Rectangle
	border        *canvas.Rectangle
	placeholder   *canvas.Text
	entryRenderer fyne.WidgetRenderer
}

func (r *textViewRenderer) Destroy() {
	r.entryRenderer.Destroy()
}

func (r *textViewRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.border.Resize(size)

	margins := r.textView.PlaceholderMargins
	r.placeholder.Move(fyne.NewPos(margins.Left+4, margins.Top+4))

	r.textView.Entry.Resize(size)
	r.entryRenderer.Layout(size)
}

func (r *textViewRenderer) MinSize() fyne.Size {
	minSize := r.entryRenderer.MinSize()

	r.textView.mu.RLock()
	maxHeight := r.textView.MaximumHeight
	r.textView.mu.RUnlock()

	if maxHeight > 0 && minSize.Height > maxHeight {
		minSize.Height = maxHeight
	}

	return minSize
}

func (r *textViewRenderer) Refresh() {
	// Show/hide placeholder based on text content
	if r.textView.Text == "" && r.textView.Placeholder != "" {
		r.placeholder.Text = r.textView.Placeholder
		r.placeholder.Color = r.textView.PlaceholderColor
		r.placeholder.Show()
	} else {
		r.placeholder.Hide()
	}

	r.entryRenderer.Refresh()
	r.background.Refresh()
	r.border.Refresh()
	r.placeholder.Refresh()
}

func (r *textViewRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background, r.border}
	objects = append(objects, r.entryRenderer.Objects()...)
	objects = append(objects, r.placeholder)
	return objects
}

// AutoGrowingTextView automatically adjusts height based on content
type AutoGrowingTextView struct {
	*TextView
	MinHeight float32
}

// NewAutoGrowingTextView creates a text view that grows with content
func NewAutoGrowingTextView() *AutoGrowingTextView {
	tv := NewTextView()
	return &AutoGrowingTextView{
		TextView:  tv,
		MinHeight: 80,
	}
}

// MinSize returns the minimum size, respecting MinHeight
func (atv *AutoGrowingTextView) MinSize() fyne.Size {
	size := atv.TextView.MinSize()
	if size.Height < atv.MinHeight {
		size.Height = atv.MinHeight
	}
	return size
}
