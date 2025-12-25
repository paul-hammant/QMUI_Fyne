package segmented

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestSegmentedControl_VisualRendering(t *testing.T) {
	var selectedIndex int
	sc := NewSegmentedControl([]string{"One", "Two", "Three"}, func(idx int) {
		selectedIndex = idx
	})

	w := test.NewWindow(sc)
	w.Resize(fyne.NewSize(300, 50))

	renderer := test.WidgetRenderer(sc)
	renderer.Layout(fyne.NewSize(300, 50))

	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("SegmentedControl has no rendered objects")
	}

	// Should have segments for each item
	t.Logf("SegmentedControl has %d objects", len(objects))

	w.Close()
	_ = selectedIndex
}

func TestSegmentedControl_SelectionChangesVisually(t *testing.T) {
	var selectedIndex int
	sc := NewSegmentedControl([]string{"A", "B", "C"}, func(idx int) {
		selectedIndex = idx
	})

	w := test.NewWindow(sc)
	w.Resize(fyne.NewSize(300, 50))

	renderer := test.WidgetRenderer(sc)
	renderer.Layout(fyne.NewSize(300, 50))

	// Initial selection
	if sc.SelectedIndex != 0 {
		t.Errorf("Initial selection should be 0, got %d", sc.SelectedIndex)
	}

	// Change selection
	sc.SetSelectedIndex(1)
	renderer.Refresh()

	if sc.SelectedIndex != 1 {
		t.Errorf("After SetSelectedIndex(1), should be 1, got %d", sc.SelectedIndex)
	}

	// Verify callback was invoked
	if selectedIndex != 1 {
		t.Errorf("Callback should have received index 1, got %d", selectedIndex)
	}

	w.Close()
}

func TestSegmentedControl_WidthDistribution(t *testing.T) {
	sc := NewSegmentedControl([]string{"First", "Second", "Third"}, func(int) {})

	w := test.NewWindow(sc)
	w.Resize(fyne.NewSize(300, 50))

	renderer := test.WidgetRenderer(sc)
	renderer.Layout(fyne.NewSize(300, 50))

	minSize := renderer.MinSize()

	// Should have reasonable width for 3 segments
	if minSize.Width < 100 {
		t.Errorf("SegmentedControl min width too small: %f", minSize.Width)
	}

	t.Logf("SegmentedControl min size: %v", minSize)
	w.Close()
}
