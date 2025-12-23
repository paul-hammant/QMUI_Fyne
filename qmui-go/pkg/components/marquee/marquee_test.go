package marquee

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestMarqueeLabel_VisualPositionChanges(t *testing.T) {
	// Create a marquee with long text
	longText := "This is a very long marquee text that should definitely scroll because it is much wider than the container"
	ml := NewMarqueeLabel(longText)
	ml.PauseDuration = 100 * time.Millisecond // Short pause for faster test
	ml.Speed = 100
	ml.AutoScrollWhenFits = false // Force scrolling

	// Create test window with narrow width
	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(100, 50))

	// Get renderer to check actual visual positions
	renderer := test.WidgetRenderer(ml)
	objects := renderer.Objects()

	if len(objects) < 1 {
		t.Fatal("Expected at least 1 text object in renderer")
	}

	// Force initial layout
	renderer.Layout(fyne.NewSize(100, 50))

	// Record initial position of the text
	initialPos := objects[0].Position()
	t.Logf("Initial text position: %v", initialPos)

	// Start animation
	ml.StartAnimation()

	// Wait for pause + animation time
	time.Sleep(300 * time.Millisecond)

	// Force refresh and layout
	renderer.Refresh()
	renderer.Layout(fyne.NewSize(100, 50))

	// Check new position
	newPos := objects[0].Position()
	t.Logf("New text position after animation: %v", newPos)

	// The X position should have changed
	if newPos.X == initialPos.X {
		t.Errorf("Text X position did not change. Initial: %f, Current: %f", initialPos.X, newPos.X)
	}

	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_AnimationStarts(t *testing.T) {
	// Create a marquee with long text that will definitely need to scroll
	longText := "This is a very long marquee text that should definitely scroll because it is much wider than the container width"
	ml := NewMarqueeLabel(longText)

	// Create test window with narrow width to force scrolling
	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(100, 50)) // Narrow window forces text to scroll

	// Force layout to calculate text width
	w.Canvas().Refresh(ml)

	// Get initial offset
	ml.mu.RLock()
	initialOffset := ml.offset
	ml.mu.RUnlock()

	// Start animation
	ml.StartAnimation()

	// Verify animation is marked as running
	if !ml.IsAnimating {
		t.Error("Expected IsAnimating to be true after StartAnimation()")
	}

	ml.mu.RLock()
	animating := ml.animating
	ml.mu.RUnlock()

	if !animating {
		t.Error("Expected animating flag to be true")
	}

	// Wait for animation to progress (longer than pause duration + some animation time)
	// Default pause is 2 seconds, so we wait 2.5 seconds
	time.Sleep(2500 * time.Millisecond)

	// Check that offset has changed
	ml.mu.RLock()
	newOffset := ml.offset
	textWidth := ml.textWidth
	ml.mu.RUnlock()

	if textWidth == 0 {
		t.Error("textWidth was not calculated - Layout may not have been called")
	}

	if newOffset == initialOffset {
		t.Errorf("Marquee offset did not change after animation started. Initial: %f, Current: %f, TextWidth: %f",
			initialOffset, newOffset, textWidth)
	}

	// Clean up
	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_AnimationStops(t *testing.T) {
	ml := NewMarqueeLabel("Scrolling text for stop test")

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(50, 50))
	w.Canvas().Refresh(ml)

	ml.StartAnimation()

	if !ml.IsAnimating {
		t.Error("Expected IsAnimating to be true")
	}

	ml.StopAnimation()

	if ml.IsAnimating {
		t.Error("Expected IsAnimating to be false after StopAnimation()")
	}

	ml.mu.RLock()
	offset := ml.offset
	ml.mu.RUnlock()

	if offset != 0 {
		t.Errorf("Expected offset to reset to 0 after stop, got %f", offset)
	}

	w.Close()
}

func TestMarqueeLabel_NoScrollWhenTextFits(t *testing.T) {
	// Short text that fits in container
	ml := NewMarqueeLabel("Hi")
	ml.AutoScrollWhenFits = true // Default, but explicit

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(500, 50)) // Wide window - text fits
	w.Canvas().Refresh(ml)

	ml.StartAnimation()

	// Wait past the pause duration
	time.Sleep(2500 * time.Millisecond)

	ml.mu.RLock()
	offset := ml.offset
	textWidth := ml.textWidth
	containerWidth := ml.Size().Width
	ml.mu.RUnlock()

	// Text should fit, so no scrolling should occur
	if textWidth > containerWidth {
		t.Skipf("Text unexpectedly wider than container (text: %f, container: %f)", textWidth, containerWidth)
	}

	if offset != 0 {
		t.Errorf("Expected no scrolling when text fits, but offset is %f", offset)
	}

	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_OffsetChangesOverTime(t *testing.T) {
	longText := "This marquee text is intentionally very long to ensure it exceeds container width and triggers scrolling animation"
	ml := NewMarqueeLabel(longText)
	ml.PauseDuration = 100 * time.Millisecond // Short pause for faster test
	ml.Speed = 100                             // Faster speed

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(80, 50))
	w.Canvas().Refresh(ml)

	ml.StartAnimation()

	// Wait for pause to complete
	time.Sleep(200 * time.Millisecond)

	// Collect multiple offset readings
	var offsets []float32
	for i := 0; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)
		ml.mu.RLock()
		offsets = append(offsets, ml.offset)
		ml.mu.RUnlock()
	}

	// Check that offsets are changing (animation is progressing)
	allSame := true
	for i := 1; i < len(offsets); i++ {
		if offsets[i] != offsets[0] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Errorf("Marquee offset not changing over time. All readings: %v", offsets)
	}

	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_DirectionRight(t *testing.T) {
	longText := "Right scrolling marquee text that is long enough to scroll"
	ml := NewMarqueeLabel(longText)
	ml.Direction = MarqueeDirectionRight
	ml.PauseDuration = 100 * time.Millisecond
	ml.Speed = 100

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(80, 50))
	w.Canvas().Refresh(ml)

	ml.mu.RLock()
	initialOffset := ml.offset
	ml.mu.RUnlock()

	ml.StartAnimation()

	// Wait for animation
	time.Sleep(300 * time.Millisecond)

	ml.mu.RLock()
	newOffset := ml.offset
	ml.mu.RUnlock()

	// Right direction should increase offset
	if newOffset <= initialOffset {
		t.Errorf("Right direction should increase offset. Initial: %f, Current: %f", initialOffset, newOffset)
	}

	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_DiagnoseSizing(t *testing.T) {
	// This test diagnoses sizing issues that might prevent scrolling
	longText := "This is long text for the marquee that should scroll"
	ml := NewMarqueeLabel(longText)
	ml.AutoScrollWhenFits = false
	ml.PauseDuration = 50 * time.Millisecond
	ml.Speed = 100

	w := test.NewWindow(ml)

	// Test various container sizes
	testSizes := []fyne.Size{
		{Width: 50, Height: 30},
		{Width: 100, Height: 30},
		{Width: 200, Height: 30},
		{Width: 500, Height: 30},
		{Width: 960, Height: 30}, // Wide window like the demo
	}

	for _, size := range testSizes {
		w.Resize(size)
		w.Canvas().Refresh(ml)

		// Force layout
		renderer := test.WidgetRenderer(ml)
		renderer.Layout(size)

		ml.mu.RLock()
		textWidth := ml.textWidth
		ml.mu.RUnlock()

		containerWidth := ml.Size().Width

		t.Logf("Container: %.0f, TextWidth: %.0f, NeedsScroll: %v",
			containerWidth, textWidth, textWidth > containerWidth)

		// Start animation briefly
		ml.offset = 0
		ml.StartAnimation()
		time.Sleep(150 * time.Millisecond)

		ml.mu.RLock()
		offset := ml.offset
		ml.mu.RUnlock()

		ml.StopAnimation()

		if offset == 0 {
			t.Errorf("At container width %.0f: offset did not change (textWidth=%.0f)",
				containerWidth, textWidth)
		} else {
			t.Logf("At container width %.0f: offset changed to %.2f", containerWidth, offset)
		}
	}

	w.Close()
}
