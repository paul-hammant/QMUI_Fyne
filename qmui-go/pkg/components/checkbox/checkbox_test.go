package checkbox

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestCheckbox_VisualStateToggle(t *testing.T) {
	changed := false
	var lastState bool

	cb := NewCheckbox(func(selected bool) {
		changed = true
		lastState = selected
	})

	w := test.NewWindow(cb)
	w.Resize(fyne.NewSize(100, 50))

	renderer := test.WidgetRenderer(cb)
	renderer.Layout(fyne.NewSize(100, 50))

	// Initial state should be unchecked
	if cb.Selected {
		t.Error("Checkbox should start unchecked")
	}

	// Get initial visual state
	initialObjects := renderer.Objects()
	if len(initialObjects) == 0 {
		t.Fatal("Checkbox has no rendered objects")
	}

	// Tap to toggle
	test.Tap(cb)

	if !changed {
		t.Error("Checkbox callback was not invoked")
	}
	if !lastState {
		t.Error("Checkbox should be selected after tap")
	}
	if !cb.Selected {
		t.Error("Checkbox.Selected should be true after tap")
	}

	// Tap again to toggle back
	test.Tap(cb)

	if lastState {
		t.Error("Checkbox should be deselected after second tap")
	}

	w.Close()
}

func TestCheckbox_SetSelectedVisualUpdate(t *testing.T) {
	cb := NewCheckbox(func(bool) {})

	w := test.NewWindow(cb)
	w.Resize(fyne.NewSize(100, 50))

	renderer := test.WidgetRenderer(cb)
	renderer.Layout(fyne.NewSize(100, 50))

	cb.SetSelected(true)
	renderer.Refresh()

	if !cb.Selected {
		t.Error("SetSelected(true) should set Selected to true")
	}

	cb.SetSelected(false)
	renderer.Refresh()

	if cb.Selected {
		t.Error("SetSelected(false) should set Selected to false")
	}

	w.Close()
}

func TestCheckbox_MinSizeReasonable(t *testing.T) {
	cb := NewCheckbox(func(bool) {})

	w := test.NewWindow(cb)
	renderer := test.WidgetRenderer(cb)

	minSize := renderer.MinSize()

	// Checkbox should have some minimum size
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("Checkbox min size invalid: %v", minSize)
	}

	t.Logf("Checkbox min size: %v", minSize)
	w.Close()
}
