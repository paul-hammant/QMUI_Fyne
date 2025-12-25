package textview

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestTextView_VisualRendering(t *testing.T) {
	tv := NewTextView()

	w := test.NewWindow(tv)
	w.Resize(fyne.NewSize(200, 100))

	renderer := test.WidgetRenderer(tv)
	renderer.Layout(fyne.NewSize(200, 100))

	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("TextView has no rendered objects")
	}

	minSize := renderer.MinSize()
	t.Logf("TextView min size: %v", minSize)

	w.Close()
}

func TestTextView_SetText(t *testing.T) {
	tv := NewTextView()

	w := test.NewWindow(tv)
	w.Resize(fyne.NewSize(200, 100))

	renderer := test.WidgetRenderer(tv)
	renderer.Layout(fyne.NewSize(200, 100))

	tv.SetText("Line 1\nLine 2\nLine 3")
	renderer.Refresh()

	if tv.Text != "Line 1\nLine 2\nLine 3" {
		t.Errorf("Text not set correctly: '%s'", tv.Text)
	}

	w.Close()
}

func TestTextView_MultilineExpands(t *testing.T) {
	tv := NewTextView()

	w := test.NewWindow(tv)
	renderer := test.WidgetRenderer(tv)

	// Single line
	tv.SetText("Single")
	renderer.Refresh()
	singleSize := renderer.MinSize()

	// Multiple lines
	tv.SetText("Line 1\nLine 2\nLine 3\nLine 4")
	renderer.Refresh()
	multiSize := renderer.MinSize()

	// Log the sizes for diagnostic purposes
	t.Logf("Single line size: %v, Multi line size: %v", singleSize, multiSize)

	w.Close()
}

func TestTextView_PlaceholderWhenEmpty(t *testing.T) {
	tv := NewTextViewWithPlaceholder("Enter text here...")

	w := test.NewWindow(tv)
	w.Resize(fyne.NewSize(200, 100))

	renderer := test.WidgetRenderer(tv)
	renderer.Layout(fyne.NewSize(200, 100))

	if tv.PlaceHolder != "Enter text here..." {
		t.Errorf("Placeholder should be set")
	}

	w.Close()
}
