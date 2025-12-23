package badge

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestBadgeLabel_VisualRendering(t *testing.T) {
	badge := NewBadgeLabel("99+")

	w := test.NewWindow(badge)
	w.Resize(fyne.NewSize(100, 50))

	renderer := test.WidgetRenderer(badge)
	renderer.Layout(fyne.NewSize(100, 50))

	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("BadgeLabel has no rendered objects")
	}

	minSize := renderer.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("BadgeLabel has invalid min size: %v", minSize)
	}

	// Badge should be relatively small
	if minSize.Height > 30 {
		t.Errorf("Badge height seems too large: %f", minSize.Height)
	}

	t.Logf("Badge min size: %v", minSize)
	w.Close()
}

func TestBadgeLabel_TextChangesSize(t *testing.T) {
	badge := NewBadgeLabel("1")

	w := test.NewWindow(badge)
	renderer := test.WidgetRenderer(badge)

	smallSize := renderer.MinSize()

	badge.SetText("99999")
	renderer.Refresh()

	largeSize := renderer.MinSize()

	if largeSize.Width <= smallSize.Width {
		t.Errorf("Larger badge text should increase width. Small: %f, Large: %f",
			smallSize.Width, largeSize.Width)
	}

	w.Close()
}

func TestBadgeLabel_VisualPosition(t *testing.T) {
	badge := NewBadgeLabel("NEW")

	w := test.NewWindow(badge)
	w.Resize(fyne.NewSize(100, 50))

	renderer := test.WidgetRenderer(badge)
	renderer.Layout(fyne.NewSize(100, 50))

	objects := renderer.Objects()
	for i, obj := range objects {
		pos := obj.Position()
		size := obj.Size()
		t.Logf("Badge object %d: pos=%v size=%v", i, pos, size)

		// Objects should have valid sizes
		if size.Width < 0 || size.Height < 0 {
			t.Errorf("Object %d has negative size: %v", i, size)
		}
	}

	w.Close()
}
