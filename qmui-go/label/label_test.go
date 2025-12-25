package label

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/paul-hammant/qmui_fyne/core"
)

func TestLabel_VisualRendering(t *testing.T) {
	lbl := NewLabel("Test Label")

	w := test.NewWindow(lbl)
	w.Resize(fyne.NewSize(200, 50))

	renderer := test.WidgetRenderer(lbl)
	renderer.Layout(fyne.NewSize(200, 50))

	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("Label renderer has no objects")
	}

	minSize := renderer.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("Label has invalid min size: %v", minSize)
	}

	t.Logf("Label min size: %v", minSize)
	w.Close()
}

func TestLabel_ContentInsetsAffectSize(t *testing.T) {
	lbl := NewLabel("Inset Test")

	w := test.NewWindow(lbl)
	renderer := test.WidgetRenderer(lbl)
	sizeNoInsets := renderer.MinSize()

	// Add insets via ContentEdgeInsets field
	lbl.ContentEdgeInsets = core.EdgeInsets{Top: 10, Left: 20, Bottom: 10, Right: 20}
	renderer.Refresh()
	sizeWithInsets := renderer.MinSize()

	// Size should increase by inset amounts
	expectedWidthIncrease := float32(40)
	expectedHeightIncrease := float32(20)

	widthDiff := sizeWithInsets.Width - sizeNoInsets.Width
	heightDiff := sizeWithInsets.Height - sizeNoInsets.Height

	if widthDiff < expectedWidthIncrease-1 || widthDiff > expectedWidthIncrease+1 {
		t.Errorf("Width increase expected ~%f, got %f", expectedWidthIncrease, widthDiff)
	}
	if heightDiff < expectedHeightIncrease-1 || heightDiff > expectedHeightIncrease+1 {
		t.Errorf("Height increase expected ~%f, got %f", expectedHeightIncrease, heightDiff)
	}

	w.Close()
}

func TestLabel_TextChangesVisually(t *testing.T) {
	lbl := NewLabel("Short")

	w := test.NewWindow(lbl)
	renderer := test.WidgetRenderer(lbl)

	shortSize := renderer.MinSize()

	lbl.SetText("This is a much longer text string")
	renderer.Refresh()

	longSize := renderer.MinSize()

	if longSize.Width <= shortSize.Width {
		t.Errorf("Longer text should have larger width. Short: %f, Long: %f",
			shortSize.Width, longSize.Width)
	}

	w.Close()
}
