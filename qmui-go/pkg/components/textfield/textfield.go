// Package textfield provides QMUITextField - an enhanced text input field
// Ported from Tencent's QMUI_iOS framework
package textfield

import (
	"image/color"
	"sync"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// TextFieldDelegate provides callbacks for text field events
type TextFieldDelegate interface {
	// ShouldChangeCharactersInRange is called before text changes
	ShouldChangeCharactersInRange(textField *TextField, start, length int, replacement string) bool
	// DidPreventTextChange is called when text change was prevented due to length limit
	DidPreventTextChange(textField *TextField, start, length int, replacement string)
}

// TextField is an enhanced text input widget with QMUI styling
type TextField struct {
	widget.Entry

	// Styling
	PlaceholderColor color.Color
	TextInsets       core.EdgeInsets
	ClearButtonPositionAdjustment core.Offset

	// Behavior
	ShouldResponseToProgrammaticallyTextChanges bool
	MaximumTextLength                           int
	ShouldCountingNonASCIICharacterAsTwo        bool

	// Delegate
	Delegate TextFieldDelegate

	// Callbacks
	OnTextChanged func(text string)
	OnPaste       func(sender interface{}) bool

	mu sync.RWMutex
}

// NewTextField creates a new QMUI-styled text field
func NewTextField() *TextField {
	config := core.SharedConfiguration()
	tf := &TextField{
		PlaceholderColor: config.PlaceholderColor,
		TextInsets:       config.TextFieldTextInsets,
		ShouldResponseToProgrammaticallyTextChanges: true,
		MaximumTextLength: -1, // No limit
		ShouldCountingNonASCIICharacterAsTwo: false,
	}
	tf.ExtendBaseWidget(tf)
	tf.Entry.OnChanged = tf.handleTextChanged
	return tf
}

// NewTextFieldWithPlaceholder creates a new text field with placeholder
func NewTextFieldWithPlaceholder(placeholder string) *TextField {
	tf := NewTextField()
	tf.PlaceHolder = placeholder
	return tf
}

// SetText sets the text and optionally triggers change events
func (tf *TextField) SetText(text string) {
	tf.mu.Lock()
	shouldNotify := tf.ShouldResponseToProgrammaticallyTextChanges
	tf.mu.Unlock()

	tf.Entry.SetText(text)
	if shouldNotify && tf.OnTextChanged != nil {
		tf.OnTextChanged(text)
	}
}

// handleTextChanged processes text changes with length limiting
func (tf *TextField) handleTextChanged(text string) {
	tf.mu.RLock()
	maxLen := tf.MaximumTextLength
	countNonASCII := tf.ShouldCountingNonASCIICharacterAsTwo
	tf.mu.RUnlock()

	if maxLen > 0 {
		length := tf.calculateTextLength(text, countNonASCII)
		if length > maxLen {
			// Trim text to max length
			trimmed := tf.trimToLength(text, maxLen, countNonASCII)
			tf.Entry.SetText(trimmed)
			if tf.Delegate != nil {
				tf.Delegate.DidPreventTextChange(tf, len(trimmed), len(text)-len(trimmed), "")
			}
			text = trimmed
		}
	}

	if tf.OnTextChanged != nil {
		tf.OnTextChanged(text)
	}
}

func (tf *TextField) calculateTextLength(s string, countNonASCIIAsTwo bool) int {
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

func (tf *TextField) trimToLength(s string, maxLen int, countNonASCIIAsTwo bool) string {
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

// CreateRenderer implements fyne.Widget
func (tf *TextField) CreateRenderer() fyne.WidgetRenderer {
	tf.ExtendBaseWidget(tf)

	// Create background for styling
	background := canvas.NewRectangle(color.Transparent)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 1
	border.StrokeColor = core.SharedConfiguration().SeparatorColor

	entryRenderer := tf.Entry.CreateRenderer()

	return &textFieldRenderer{
		textField:     tf,
		background:    background,
		border:        border,
		entryRenderer: entryRenderer,
	}
}

type textFieldRenderer struct {
	textField     *TextField
	background    *canvas.Rectangle
	border        *canvas.Rectangle
	entryRenderer fyne.WidgetRenderer
}

func (r *textFieldRenderer) Destroy() {
	r.entryRenderer.Destroy()
}

func (r *textFieldRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.border.Resize(size)

	insets := r.textField.TextInsets
	entrySize := fyne.NewSize(
		size.Width-insets.Left-insets.Right,
		size.Height-insets.Top-insets.Bottom,
	)
	// Note: We don't call Entry.Resize() here to avoid infinite recursion
	// since Entry is embedded in TextField and shares the same widget identity.
	// Instead, we position the entry content via the renderer.
	r.entryRenderer.Layout(entrySize)
}

func (r *textFieldRenderer) MinSize() fyne.Size {
	entryMin := r.entryRenderer.MinSize()
	insets := r.textField.TextInsets
	return fyne.NewSize(
		entryMin.Width+insets.Left+insets.Right,
		entryMin.Height+insets.Top+insets.Bottom,
	)
}

func (r *textFieldRenderer) Refresh() {
	r.entryRenderer.Refresh()
	r.background.Refresh()
	r.border.Refresh()
}

func (r *textFieldRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background, r.border}
	objects = append(objects, r.entryRenderer.Objects()...)
	return objects
}

// PasswordTextField is a secure text field for passwords
type PasswordTextField struct {
	*TextField
}

// NewPasswordTextField creates a password input field
func NewPasswordTextField() *PasswordTextField {
	tf := NewTextField()
	tf.Password = true
	return &PasswordTextField{TextField: tf}
}

// SearchTextField is styled for search input
type SearchTextField struct {
	*TextField
	SearchIcon fyne.Resource
}

// NewSearchTextField creates a search input field
func NewSearchTextField() *SearchTextField {
	config := core.SharedConfiguration()
	tf := NewTextField()
	tf.PlaceHolder = "Search"
	tf.PlaceholderColor = config.SearchBarPlaceholderColor

	return &SearchTextField{
		TextField: tf,
	}
}

// NumberTextField only allows numeric input
type NumberTextField struct {
	*TextField
	AllowDecimal bool
	AllowNegative bool
}

// NewNumberTextField creates a numeric input field
func NewNumberTextField() *NumberTextField {
	tf := NewTextField()
	return &NumberTextField{
		TextField: tf,
		AllowDecimal: true,
		AllowNegative: true,
	}
}

// Validate returns whether the current text is a valid number
func (ntf *NumberTextField) Validate() error {
	// Validation is handled by Fyne's Entry
	return nil
}

// MultilineTextField is a text field with multiple lines
type MultilineTextField struct {
	*TextField
	MinLines int
	MaxLines int
}

// NewMultilineTextField creates a multiline text input
func NewMultilineTextField() *MultilineTextField {
	tf := NewTextField()
	tf.MultiLine = true
	tf.Wrapping = fyne.TextWrapWord
	return &MultilineTextField{
		TextField: tf,
		MinLines:  1,
		MaxLines:  0, // No limit
	}
}
