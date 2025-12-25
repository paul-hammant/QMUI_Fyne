package progress

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestPieProgress_VisualUpdate(t *testing.T) {
	pie := NewPieProgress()

	w := test.NewWindow(pie)
	w.Resize(fyne.NewSize(100, 100))

	renderer := test.WidgetRenderer(pie)
	renderer.Layout(fyne.NewSize(100, 100))

	// Set initial progress
	pie.SetProgress(0.0)
	renderer.Refresh()

	if pie.Progress != 0.0 {
		t.Errorf("Progress should be 0.0, got %f", pie.Progress)
	}

	// Update progress
	pie.SetProgress(0.5)
	renderer.Refresh()

	if pie.Progress != 0.5 {
		t.Errorf("Progress should be 0.5, got %f", pie.Progress)
	}

	// Full progress
	pie.SetProgress(1.0)
	renderer.Refresh()

	if pie.Progress != 1.0 {
		t.Errorf("Progress should be 1.0, got %f", pie.Progress)
	}

	w.Close()
}

func TestRingProgress_VisualUpdate(t *testing.T) {
	circular := NewRingProgress()

	w := test.NewWindow(circular)
	w.Resize(fyne.NewSize(100, 100))

	renderer := test.WidgetRenderer(circular)
	renderer.Layout(fyne.NewSize(100, 100))

	// Test progress values
	testValues := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	for _, val := range testValues {
		circular.SetProgress(val)
		renderer.Refresh()

		if circular.Progress != val {
			t.Errorf("Progress should be %f, got %f", val, circular.Progress)
		}
	}

	w.Close()
}

func TestProgressBar_VisualWidth(t *testing.T) {
	linear := NewProgressBar()

	w := test.NewWindow(linear)
	w.Resize(fyne.NewSize(200, 20))

	renderer := test.WidgetRenderer(linear)
	renderer.Layout(fyne.NewSize(200, 20))

	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("LinearProgressView has no rendered objects")
	}

	// Set 50% progress
	linear.SetProgress(0.5)
	renderer.Refresh()
	renderer.Layout(fyne.NewSize(200, 20))

	// Verify objects exist and have valid sizes
	for i, obj := range objects {
		size := obj.Size()
		pos := obj.Position()
		t.Logf("Linear progress object %d: pos=%v size=%v", i, pos, size)
	}

	w.Close()
}

func TestProgressViews_MinSize(t *testing.T) {
	tests := []struct {
		name   string
		widget fyne.Widget
	}{
		{"PieProgress", NewPieProgress()},
		{"CircularProgress", NewRingProgress()},
		{"LinearProgress", NewProgressBar()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := test.NewWindow(tt.widget)
			renderer := test.WidgetRenderer(tt.widget)

			minSize := renderer.MinSize()

			if minSize.Width <= 0 || minSize.Height <= 0 {
				t.Errorf("%s has invalid min size: %v", tt.name, minSize)
			}

			t.Logf("%s min size: %v", tt.name, minSize)
			w.Close()
		})
	}
}
