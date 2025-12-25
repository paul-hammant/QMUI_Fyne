package textfield

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestTextField_VisualRendering(t *testing.T) {
	tf := NewTextField()

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(200, 40))

	renderer := test.WidgetRenderer(tf)
	renderer.Layout(fyne.NewSize(200, 40))

	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("TextField has no rendered objects")
	}

	minSize := renderer.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("TextField has invalid min size: %v", minSize)
	}

	t.Logf("TextField min size: %v", minSize)
	w.Close()
}

func TestTextField_TextEntryVisual(t *testing.T) {
	tf := NewTextField()

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(200, 40))

	renderer := test.WidgetRenderer(tf)
	renderer.Layout(fyne.NewSize(200, 40))

	// Set text
	tf.SetText("Hello World")
	renderer.Refresh()

	if tf.Text != "Hello World" {
		t.Errorf("Text should be 'Hello World', got '%s'", tf.Text)
	}

	// Clear text
	tf.SetText("")
	renderer.Refresh()

	if tf.Text != "" {
		t.Errorf("Text should be empty, got '%s'", tf.Text)
	}

	w.Close()
}

func TestTextField_WithPlaceholder(t *testing.T) {
	tf := NewTextFieldWithPlaceholder("Enter name...")

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(200, 40))

	renderer := test.WidgetRenderer(tf)
	renderer.Layout(fyne.NewSize(200, 40))

	if tf.PlaceHolder != "Enter name..." {
		t.Errorf("Placeholder should be set, got '%s'", tf.PlaceHolder)
	}

	w.Close()
}

func TestTextField_MaxLength(t *testing.T) {
	tf := NewTextField()
	tf.MaximumTextLength = 5

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(200, 40))

	renderer := test.WidgetRenderer(tf)
	renderer.Layout(fyne.NewSize(200, 40))

	// Try to set text longer than max
	tf.Entry.SetText("1234567890")
	renderer.Refresh()

	// Text should be trimmed
	if len(tf.Text) > 5 {
		t.Errorf("Text should be limited to 5 chars, got %d: '%s'", len(tf.Text), tf.Text)
	}

	w.Close()
}

func TestTextField_OnChangedCallback(t *testing.T) {
	tf := NewTextField()
	var lastText string
	tf.OnTextChanged = func(text string) {
		lastText = text
	}

	w := test.NewWindow(tf)
	renderer := test.WidgetRenderer(tf)

	tf.SetText("test")
	renderer.Refresh()

	if lastText != "test" {
		t.Errorf("OnTextChanged callback should have received 'test', got '%s'", lastText)
	}

	w.Close()
}
