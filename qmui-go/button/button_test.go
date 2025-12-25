package button

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestButton_VisualRendering(t *testing.T) {
	btn := NewButton("Test Button", func() {})

	w := test.NewWindow(btn)
	w.Resize(fyne.NewSize(200, 50))

	renderer := test.WidgetRenderer(btn)
	objects := renderer.Objects()

	if len(objects) == 0 {
		t.Fatal("Button renderer has no objects")
	}

	minSize := renderer.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("Button has invalid min size: %v", minSize)
	}

	renderer.Layout(fyne.NewSize(200, 50))
	t.Logf("Button min size: %v", minSize)

	w.Close()
}

func TestFillButton_VisualStateChange(t *testing.T) {
	tapped := false
	btn := NewFillButton("Fill Button", color.RGBA{R: 0, G: 122, B: 255, A: 255}, func() {
		tapped = true
	})

	w := test.NewWindow(btn)
	w.Resize(fyne.NewSize(200, 50))

	renderer := test.WidgetRenderer(btn)
	renderer.Layout(fyne.NewSize(200, 50))

	initialObjects := renderer.Objects()
	if len(initialObjects) == 0 {
		t.Fatal("FillButton has no rendered objects")
	}

	test.Tap(btn)

	if !tapped {
		t.Error("Button tap callback was not invoked")
	}

	w.Close()
}

func TestGhostButton_VisualBorder(t *testing.T) {
	btn := NewGhostButton("Ghost", color.RGBA{R: 255, G: 0, B: 0, A: 255}, func() {})

	w := test.NewWindow(btn)
	w.Resize(fyne.NewSize(150, 40))

	renderer := test.WidgetRenderer(btn)
	renderer.Layout(fyne.NewSize(150, 40))

	minSize := renderer.MinSize()
	if minSize.Width <= 0 {
		t.Errorf("GhostButton min width should be positive: %f", minSize.Width)
	}

	w.Close()
}

func TestNavigationButton_IconAndTitle(t *testing.T) {
	btn := NewNavigationButton("Back", func() {})

	w := test.NewWindow(btn)
	w.Resize(fyne.NewSize(100, 44))

	renderer := test.WidgetRenderer(btn)
	renderer.Layout(fyne.NewSize(100, 44))

	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("NavigationButton has no rendered objects")
	}

	w.Close()
}
