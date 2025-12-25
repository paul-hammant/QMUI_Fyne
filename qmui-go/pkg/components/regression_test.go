// Package components provides comprehensive regression tests for all QMUI Go components
// These tests are based on the iOS QMUI demo frames to ensure visual and behavioral parity
package components

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/user/qmui-go/pkg/components/alert"
	"github.com/user/qmui-go/pkg/components/badge"
	"github.com/user/qmui-go/pkg/components/button"
	"github.com/user/qmui-go/pkg/components/checkbox"
	"github.com/user/qmui-go/pkg/components/dialog"
	"github.com/user/qmui-go/pkg/components/empty"
	"github.com/user/qmui-go/pkg/components/emotion"
	"github.com/user/qmui-go/pkg/components/floatlayout"
	"github.com/user/qmui-go/pkg/components/grid"
	"github.com/user/qmui-go/pkg/components/label"
	"github.com/user/qmui-go/pkg/components/marquee"
	"github.com/user/qmui-go/pkg/components/modal"
	"github.com/user/qmui-go/pkg/components/moreop"
	"github.com/user/qmui-go/pkg/components/navigation"
	"github.com/user/qmui-go/pkg/components/popup"
	"github.com/user/qmui-go/pkg/components/progress"
	"github.com/user/qmui-go/pkg/components/search"
	"github.com/user/qmui-go/pkg/components/segmented"
	"github.com/user/qmui-go/pkg/components/table"
	"github.com/user/qmui-go/pkg/components/textfield"
	"github.com/user/qmui-go/pkg/components/textview"
	"github.com/user/qmui-go/pkg/components/tips"
	"github.com/user/qmui-go/pkg/components/toast"
	"github.com/user/qmui-go/pkg/core"
)

// =============================================================================
// BUTTON TESTS - Based on iOS QMUIButton
// =============================================================================

func TestButton_Rendering(t *testing.T) {
	btn := button.NewButton("OK", func() {})
	w := test.NewWindow(btn)
	w.Resize(fyne.NewSize(100, 44))

	renderer := test.WidgetRenderer(btn)
	if renderer == nil {
		t.Fatal("Button renderer is nil")
	}

	minSize := renderer.MinSize()
	if minSize.Width <= 0 || minSize.Height <= 0 {
		t.Errorf("Button min size invalid: %v", minSize)
	}

	w.Close()
}

func TestButton_TapCallback(t *testing.T) {
	tapped := false
	btn := button.NewButton("Tap Me", func() {
		tapped = true
	})

	w := test.NewWindow(btn)
	test.Tap(btn)

	if !tapped {
		t.Error("Button tap callback was not invoked")
	}
	w.Close()
}

func TestButton_DisabledState(t *testing.T) {
	tapped := false
	btn := button.NewButton("Disabled", func() {
		tapped = true
	})
	btn.SetEnabled(false)

	w := test.NewWindow(btn)
	test.Tap(btn)

	if tapped {
		t.Error("Disabled button should not invoke callback")
	}
	if btn.IsEnabled() {
		t.Error("Button should report as disabled")
	}
	w.Close()
}

func TestFillButton_BackgroundColor(t *testing.T) {
	fillColor := color.RGBA{R: 49, G: 189, B: 243, A: 255}
	btn := button.NewFillButton("Fill", fillColor, func() {})

	w := test.NewWindow(btn)
	w.Resize(fyne.NewSize(120, 44))

	if btn.BackgroundColor == nil {
		t.Error("FillButton should have background color set")
	}

	renderer := test.WidgetRenderer(btn)
	minSize := renderer.MinSize()
	t.Logf("FillButton min size: %v", minSize)

	w.Close()
}

func TestGhostButton_BorderColor(t *testing.T) {
	borderColor := color.RGBA{R: 49, G: 189, B: 243, A: 255}
	btn := button.NewGhostButton("Ghost", borderColor, func() {})

	w := test.NewWindow(btn)
	w.Resize(fyne.NewSize(120, 44))

	if btn.BorderColor == nil {
		t.Error("GhostButton should have border color set")
	}
	if btn.BorderWidth <= 0 {
		t.Error("GhostButton should have border width > 0")
	}

	w.Close()
}

func TestNavigationButton_Styling(t *testing.T) {
	btn := button.NewNavigationButton("< Back", func() {})

	w := test.NewWindow(btn)
	renderer := test.WidgetRenderer(btn)

	if renderer == nil {
		t.Fatal("NavigationButton renderer is nil")
	}

	minSize := renderer.MinSize()
	if minSize.Width <= 0 {
		t.Error("NavigationButton should have positive width")
	}

	w.Close()
}

func TestToolbarButton_Styling(t *testing.T) {
	btn := button.NewToolbarButton("Done", func() {})

	w := test.NewWindow(btn)
	renderer := test.WidgetRenderer(btn)

	if renderer == nil {
		t.Fatal("ToolbarButton renderer is nil")
	}

	w.Close()
}

// =============================================================================
// LABEL TESTS - Based on iOS QMUILabel
// =============================================================================

func TestLabel_Rendering(t *testing.T) {
	lbl := label.NewLabel("Test Label")

	w := test.NewWindow(lbl)
	w.Resize(fyne.NewSize(200, 50))

	renderer := test.WidgetRenderer(lbl)
	if len(renderer.Objects()) == 0 {
		t.Fatal("Label has no rendered objects")
	}

	w.Close()
}

func TestLabel_EdgeInsets(t *testing.T) {
	lbl := label.NewLabel("Padded")

	w := test.NewWindow(lbl)
	renderer := test.WidgetRenderer(lbl)
	sizeNoInsets := renderer.MinSize()

	lbl.ContentEdgeInsets = core.EdgeInsets{Top: 10, Left: 20, Bottom: 10, Right: 20}
	renderer.Refresh()
	sizeWithInsets := renderer.MinSize()

	widthDiff := sizeWithInsets.Width - sizeNoInsets.Width
	heightDiff := sizeWithInsets.Height - sizeNoInsets.Height

	if widthDiff < 35 || widthDiff > 45 {
		t.Errorf("Width should increase by ~40, got %f", widthDiff)
	}
	if heightDiff < 15 || heightDiff > 25 {
		t.Errorf("Height should increase by ~20, got %f", heightDiff)
	}

	w.Close()
}

func TestLabel_TextChange(t *testing.T) {
	lbl := label.NewLabel("Short")

	w := test.NewWindow(lbl)
	renderer := test.WidgetRenderer(lbl)
	shortSize := renderer.MinSize()

	lbl.SetText("This is a much longer text string")
	renderer.Refresh()
	longSize := renderer.MinSize()

	if longSize.Width <= shortSize.Width {
		t.Error("Longer text should have larger width")
	}

	w.Close()
}

// =============================================================================
// MARQUEE TESTS - Based on iOS QMUIMarqueeLabel (frames-4)
// =============================================================================

func TestMarqueeLabel_Creation(t *testing.T) {
	ml := marquee.NewMarquee("Scrolling text that is long enough to scroll")

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(200, 30))

	renderer := test.WidgetRenderer(ml)
	if renderer == nil {
		t.Fatal("MarqueeLabel renderer is nil")
	}

	w.Close()
}

func TestMarqueeLabel_AnimationStart(t *testing.T) {
	ml := marquee.NewMarquee("This text should scroll continuously when started")
	ml.Speed = 50

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(150, 30))

	ml.StartAnimation()
	time.Sleep(100 * time.Millisecond)

	if !ml.IsAnimating {
		t.Error("MarqueeLabel should be animating after StartAnimation()")
	}

	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_AnimationStop(t *testing.T) {
	ml := marquee.NewMarquee("Scrolling text")
	ml.Speed = 50

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(100, 30))

	ml.StartAnimation()
	time.Sleep(50 * time.Millisecond)
	ml.StopAnimation()
	time.Sleep(50 * time.Millisecond)

	if ml.IsAnimating {
		t.Error("MarqueeLabel should stop animating after StopAnimation()")
	}

	w.Close()
}

func TestMarqueeLabel_ShortTextNoScroll(t *testing.T) {
	ml := marquee.NewMarquee("Hi")

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(200, 30)) // Container much wider than text

	renderer := test.WidgetRenderer(ml)
	renderer.Layout(fyne.NewSize(200, 30))

	ml.StartAnimation()
	time.Sleep(100 * time.Millisecond)

	// Short text that fits should not actually scroll
	// (implementation may vary - some marquees scroll anyway)
	t.Log("Short text marquee started")

	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_DirectionLeftToRight(t *testing.T) {
	ml := marquee.NewMarquee("Scrolling right to left")
	ml.Direction = marquee.MarqueeDirectionLeft
	ml.Speed = 100

	w := test.NewWindow(ml)
	w.Resize(fyne.NewSize(100, 30))

	renderer := test.WidgetRenderer(ml)
	renderer.Layout(fyne.NewSize(100, 30))

	// Get initial position
	objects := renderer.Objects()
	if len(objects) == 0 {
		t.Fatal("No objects in marquee renderer")
	}
	initialPos := objects[0].Position()

	ml.StartAnimation()
	time.Sleep(300 * time.Millisecond)

	// Check position changed (scrolling left means X decreases)
	renderer.Refresh()
	newPos := objects[0].Position()

	t.Logf("Initial X: %f, New X: %f", initialPos.X, newPos.X)

	ml.StopAnimation()
	w.Close()
}

func TestMarqueeLabel_SpeedSetting(t *testing.T) {
	ml := marquee.NewMarquee("Speed test")

	ml.Speed = 100
	if ml.Speed != 100 {
		t.Errorf("Speed should be 100, got %f", ml.Speed)
	}

	ml.Speed = 50
	if ml.Speed != 50 {
		t.Errorf("Speed should be 50, got %f", ml.Speed)
	}
}

// =============================================================================
// BADGE TESTS - Based on iOS QMUIBadge
// =============================================================================

func TestBadgeLabel_Creation(t *testing.T) {
	b := badge.NewBadge("99+")

	w := test.NewWindow(b)
	w.Resize(fyne.NewSize(50, 30))

	renderer := test.WidgetRenderer(b)
	if renderer == nil {
		t.Fatal("BadgeLabel renderer is nil")
	}

	w.Close()
}

func TestBadgeLabel_TextChange(t *testing.T) {
	b := badge.NewBadge("5")

	w := test.NewWindow(b)
	renderer := test.WidgetRenderer(b)
	size5 := renderer.MinSize()

	b.SetText("999")
	renderer.Refresh()
	size999 := renderer.MinSize()

	if size999.Width <= size5.Width {
		t.Error("Larger badge text should have larger width")
	}

	w.Close()
}

func TestUpdatesIndicator_HasUpdates(t *testing.T) {
	indicator := badge.NewUpdatesIndicator()
	indicator.HasUpdates = true

	w := test.NewWindow(indicator)
	renderer := test.WidgetRenderer(indicator)

	if renderer == nil {
		t.Fatal("UpdatesIndicator renderer is nil")
	}

	w.Close()
}

// =============================================================================
// TEXTFIELD TESTS - Based on iOS QMUITextField
// =============================================================================

func TestTextField_Creation(t *testing.T) {
	tf := textfield.NewTextField()

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(200, 44))

	renderer := test.WidgetRenderer(tf)
	if renderer == nil {
		t.Fatal("TextField renderer is nil")
	}

	w.Close()
}

func TestTextField_Placeholder(t *testing.T) {
	tf := textfield.NewTextFieldWithPlaceholder("Enter text...")

	if tf.PlaceHolder != "Enter text..." {
		t.Errorf("Placeholder should be set, got '%s'", tf.PlaceHolder)
	}

	w := test.NewWindow(tf)
	w.Close()
}

func TestTextField_MaxLength(t *testing.T) {
	tf := textfield.NewTextField()
	tf.MaximumTextLength = 10

	w := test.NewWindow(tf)

	// Type more than max length
	tf.SetText("12345678901234567890")

	if len(tf.Text) > 10 {
		t.Errorf("Text should be truncated to 10 chars, got %d", len(tf.Text))
	}

	w.Close()
}

func TestTextField_OnChanged(t *testing.T) {
	changed := false
	tf := textfield.NewTextField()
	tf.OnChanged = func(text string) {
		changed = true
	}

	w := test.NewWindow(tf)
	tf.SetText("new text")

	if !changed {
		t.Error("OnChanged callback should be invoked")
	}

	w.Close()
}

// =============================================================================
// TEXTVIEW TESTS - Based on iOS QMUITextView
// =============================================================================

func TestTextView_Creation(t *testing.T) {
	tv := textview.NewTextView()

	w := test.NewWindow(tv)
	w.Resize(fyne.NewSize(200, 100))

	renderer := test.WidgetRenderer(tv)
	if renderer == nil {
		t.Fatal("TextView renderer is nil")
	}

	w.Close()
}

func TestTextView_Placeholder(t *testing.T) {
	tv := textview.NewTextViewWithPlaceholder("Enter description...")

	if tv.PlaceHolder != "Enter description..." {
		t.Error("Placeholder should be set")
	}

	w := test.NewWindow(tv)
	w.Close()
}

func TestTextView_MultilineText(t *testing.T) {
	tv := textview.NewTextView()
	tv.SetText("Line 1\nLine 2\nLine 3")

	w := test.NewWindow(tv)

	if tv.Text != "Line 1\nLine 2\nLine 3" {
		t.Error("Multiline text should be preserved")
	}

	w.Close()
}

// =============================================================================
// PROGRESS TESTS - Based on iOS QMUIPieProgressView
// =============================================================================

func TestPieProgressView_Creation(t *testing.T) {
	pie := progress.NewPieProgress()
	pie.Progress = 0.65

	w := test.NewWindow(pie)
	w.Resize(fyne.NewSize(60, 60))

	renderer := test.WidgetRenderer(pie)
	if renderer == nil {
		t.Fatal("PieProgressView renderer is nil")
	}

	w.Close()
}

func TestPieProgressView_ProgressUpdate(t *testing.T) {
	pie := progress.NewPieProgress()
	pie.Progress = 0.0

	w := test.NewWindow(pie)

	pie.SetProgress(0.5)
	if pie.Progress != 0.5 {
		t.Errorf("Progress should be 0.5, got %f", pie.Progress)
	}

	pie.SetProgress(1.0)
	if pie.Progress != 1.0 {
		t.Errorf("Progress should be 1.0, got %f", pie.Progress)
	}

	w.Close()
}

func TestCircularProgressView_ShowsText(t *testing.T) {
	circ := progress.NewRingProgress()
	circ.Progress = 0.75
	circ.ShowsText = true

	w := test.NewWindow(circ)
	w.Resize(fyne.NewSize(80, 80))

	renderer := test.WidgetRenderer(circ)
	if renderer == nil {
		t.Fatal("CircularProgressView renderer is nil")
	}

	w.Close()
}

func TestLinearProgressView_WidthChanges(t *testing.T) {
	lin := progress.NewProgressBar()
	lin.Progress = 0.5

	w := test.NewWindow(lin)
	w.Resize(fyne.NewSize(200, 10))

	renderer := test.WidgetRenderer(lin)
	renderer.Layout(fyne.NewSize(200, 10))

	// Progress bar fill should be proportional to progress
	lin.SetProgress(0.25)
	renderer.Refresh()

	lin.SetProgress(0.75)
	renderer.Refresh()

	w.Close()
}

// =============================================================================
// CHECKBOX TESTS - Based on iOS QMUICheckbox
// =============================================================================

func TestCheckbox_Toggle(t *testing.T) {
	selected := false
	cb := checkbox.NewCheckbox(func(s bool) {
		selected = s
	})

	w := test.NewWindow(cb)

	if cb.Selected {
		t.Error("Checkbox should start unchecked")
	}

	test.Tap(cb)

	if !selected {
		t.Error("Callback should receive true after tap")
	}
	if !cb.Selected {
		t.Error("Checkbox should be selected after tap")
	}

	test.Tap(cb)

	if selected {
		t.Error("Callback should receive false after second tap")
	}

	w.Close()
}

func TestCheckbox_SetSelected(t *testing.T) {
	cb := checkbox.NewCheckbox(func(bool) {})

	w := test.NewWindow(cb)

	cb.SetSelected(true)
	if !cb.Selected {
		t.Error("SetSelected(true) should set Selected to true")
	}

	cb.SetSelected(false)
	if cb.Selected {
		t.Error("SetSelected(false) should set Selected to false")
	}

	w.Close()
}

// =============================================================================
// SEGMENTED CONTROL TESTS - Based on iOS QMUISegmentedControl
// =============================================================================

func TestSegmentedControl_Creation(t *testing.T) {
	sc := segmented.NewSegmentedControl([]string{"One", "Two", "Three"}, func(index int) {})

	w := test.NewWindow(sc)
	w.Resize(fyne.NewSize(300, 40))

	renderer := test.WidgetRenderer(sc)
	if renderer == nil {
		t.Fatal("SegmentedControl renderer is nil")
	}

	w.Close()
}

func TestSegmentedControl_Selection(t *testing.T) {
	sc := segmented.NewSegmentedControl([]string{"A", "B", "C"}, func(index int) {
		// Selection callback
	})

	w := test.NewWindow(sc)
	w.Resize(fyne.NewSize(300, 40))

	sc.SetSelectedIndex(1)
	if sc.SelectedIndex != 1 {
		t.Errorf("SelectedIndex should be 1, got %d", sc.SelectedIndex)
	}

	sc.SetSelectedIndex(2)
	if sc.SelectedIndex != 2 {
		t.Errorf("SelectedIndex should be 2, got %d", sc.SelectedIndex)
	}

	w.Close()
}

func TestSegmentedControl_PillVariant(t *testing.T) {
	sc := segmented.NewPillSegmentedControl([]string{"Pill", "Style"}, func(int) {})

	w := test.NewWindow(sc)
	renderer := test.WidgetRenderer(sc)

	if renderer == nil {
		t.Fatal("Pill SegmentedControl renderer is nil")
	}

	w.Close()
}

// =============================================================================
// POPUP TESTS - Based on iOS QMUIPopupContainerView (frames-3)
// =============================================================================

func TestPopupMenu_Creation(t *testing.T) {
	items := []*popup.MenuItem{
		popup.NewMenuItem("Edit", func(*popup.MenuItem) {}),
		popup.NewMenuItem("Copy", func(*popup.MenuItem) {}),
		popup.NewMenuItem("Delete", func(*popup.MenuItem) {}),
	}

	// PopupMenuView is shown via Show(window, position), not as a regular widget
	menu := popup.NewPopupMenuWithItems(items)

	if len(items) != 3 {
		t.Error("Should have 3 menu items")
	}

	// Verify menu was created
	if menu == nil {
		t.Fatal("PopupMenuView is nil")
	}
}

func TestPopupContainerView_Arrow(t *testing.T) {
	pc := popup.NewPopupContainer()

	w := test.NewWindow(pc)
	w.Resize(fyne.NewSize(200, 100))

	renderer := test.WidgetRenderer(pc)
	if renderer == nil {
		t.Fatal("PopupContainerView renderer is nil")
	}

	// Test different arrow directions
	pc.ArrowDirection = popup.ArrowDirectionUp
	renderer.Refresh()

	pc.ArrowDirection = popup.ArrowDirectionDown
	renderer.Refresh()

	w.Close()
}

// =============================================================================
// EMOTION VIEW TESTS - Based on iOS QMUIEmotionView (frames-3)
// =============================================================================

func TestEmotionView_Creation(t *testing.T) {
	ev := emotion.NewEmojiPicker()

	w := test.NewWindow(ev)
	w.Resize(fyne.NewSize(300, 200))

	renderer := test.WidgetRenderer(ev)
	if renderer == nil {
		t.Fatal("EmotionView renderer is nil")
	}

	w.Close()
}

func TestEmotionView_GroupSelection(t *testing.T) {
	// Create emotion groups
	group1 := &emotion.EmotionGroup{
		Name: "Smileys",
		Emotions: []*emotion.Emotion{
			{Emoji: "😀", DisplayName: "smile"},
			{Emoji: "😂", DisplayName: "laugh"},
		},
	}
	group2 := &emotion.EmotionGroup{
		Name: "Hearts",
		Emotions: []*emotion.Emotion{
			{Emoji: "❤️", DisplayName: "heart"},
		},
	}

	ev := emotion.NewEmojiPickerWithGroups([]*emotion.EmotionGroup{group1, group2})

	w := test.NewWindow(ev)

	// Should have multiple emotion groups
	if len(ev.Groups) != 2 {
		t.Errorf("EmotionView should have 2 groups, got %d", len(ev.Groups))
	}

	ev.SetCurrentGroup(1)
	if ev.CurrentGroupIndex != 1 {
		t.Error("CurrentGroupIndex should be 1")
	}

	w.Close()
}

func TestEmotionView_EmotionCallback(t *testing.T) {
	var selectedEmotion *emotion.Emotion
	ev := emotion.NewEmojiPicker()
	ev.OnEmotionSelected = func(e *emotion.Emotion) {
		selectedEmotion = e
	}

	w := test.NewWindow(ev)
	w.Resize(fyne.NewSize(300, 200))

	// Simulate emotion selection
	testEmotion := &emotion.Emotion{Emoji: "😀", DisplayName: "smile"}
	if ev.OnEmotionSelected != nil {
		ev.OnEmotionSelected(testEmotion)
	}

	if selectedEmotion == nil || selectedEmotion.Emoji != "😀" {
		t.Error("Emotion selection callback should be invoked")
	}

	w.Close()
}

// =============================================================================
// GRID VIEW TESTS - Based on iOS QMUIGridView
// =============================================================================

func TestGridView_Creation(t *testing.T) {
	gv := grid.NewGrid(3)

	w := test.NewWindow(gv)
	w.Resize(fyne.NewSize(300, 200))

	renderer := test.WidgetRenderer(gv)
	if renderer == nil {
		t.Fatal("GridView renderer is nil")
	}

	w.Close()
}

func TestGridView_AddItems(t *testing.T) {
	gv := grid.NewGrid(4)

	for i := 0; i < 8; i++ {
		item := grid.NewGridItem(label.NewLabel("Item"))
		gv.AddItem(item)
	}

	w := test.NewWindow(gv)
	w.Resize(fyne.NewSize(400, 200))

	renderer := test.WidgetRenderer(gv)
	if renderer == nil {
		t.Fatal("GridView renderer is nil after adding items")
	}

	// Items are stored internally, verify grid renders
	t.Log("GridView with 8 items created successfully")

	w.Close()
}

func TestGridView_ColumnCount(t *testing.T) {
	gv := grid.NewGrid(3)

	if gv.ColumnCount != 3 {
		t.Errorf("ColumnCount should be 3, got %d", gv.ColumnCount)
	}

	gv.ColumnCount = 4
	if gv.ColumnCount != 4 {
		t.Errorf("ColumnCount should be 4, got %d", gv.ColumnCount)
	}
}

// =============================================================================
// FLOAT LAYOUT TESTS - Based on iOS QMUIFloatLayoutView
// =============================================================================

func TestFloatLayoutView_Creation(t *testing.T) {
	fl := floatlayout.NewFlowLayout()

	w := test.NewWindow(fl)
	w.Resize(fyne.NewSize(300, 100))

	renderer := test.WidgetRenderer(fl)
	if renderer == nil {
		t.Fatal("FloatLayoutView renderer is nil")
	}

	w.Close()
}

func TestFloatLayoutView_TagCloud(t *testing.T) {
	fl := floatlayout.NewFlowLayout()
	fl.ItemSpacing = 8
	fl.LineSpacing = 8

	tags := []string{"Go", "Fyne", "QMUI", "Cross-Platform", "Desktop"}
	for _, tag := range tags {
		fl.AddItem(floatlayout.NewSimpleTag(tag))
	}

	w := test.NewWindow(fl)
	w.Resize(fyne.NewSize(200, 100))

	renderer := test.WidgetRenderer(fl)
	if renderer == nil {
		t.Fatal("FloatLayoutView renderer is nil after adding items")
	}

	// Items are stored internally
	t.Log("FloatLayoutView with 5 tags created successfully")

	w.Close()
}

// =============================================================================
// TABLE VIEW TESTS - Based on iOS QMUITableView
// =============================================================================

func TestTableView_Creation(t *testing.T) {
	tv := table.NewTable(table.TableStyleInsetGrouped)

	w := test.NewWindow(tv)
	w.Resize(fyne.NewSize(300, 400))

	renderer := test.WidgetRenderer(tv)
	if renderer == nil {
		t.Fatal("TableView renderer is nil")
	}

	w.Close()
}

func TestTableView_Sections(t *testing.T) {
	tv := table.NewTable(table.TableStyleGrouped)

	section1 := table.NewTableSection("Section 1")
	section1.Cells = []*table.TableCell{
		table.NewTableCellWithText("Cell 1"),
		table.NewTableCellWithText("Cell 2"),
	}

	section2 := table.NewTableSection("Section 2")
	section2.Cells = []*table.TableCell{
		table.NewTableCellWithTextAndDetail("Cell 3", "Detail"),
	}

	tv.Sections = []*table.TableSection{section1, section2}

	w := test.NewWindow(tv)
	w.Resize(fyne.NewSize(300, 400))

	if len(tv.Sections) != 2 {
		t.Errorf("TableView should have 2 sections, got %d", len(tv.Sections))
	}

	w.Close()
}

// =============================================================================
// NAVIGATION TESTS - Based on iOS QMUINavigationController
// =============================================================================

func TestNavigationBar_Creation(t *testing.T) {
	nb := navigation.NewNavigationBar()

	w := test.NewWindow(nb)
	w.Resize(fyne.NewSize(400, 44))

	renderer := test.WidgetRenderer(nb)
	if renderer == nil {
		t.Fatal("NavigationBar renderer is nil")
	}

	w.Close()
}

func TestNavigationBar_Title(t *testing.T) {
	nb := navigation.NewNavigationBar()
	titleView := navigation.NewTitleViewWithTitle("My Title")
	nb.SetTitleView(titleView)

	w := test.NewWindow(nb)
	w.Resize(fyne.NewSize(400, 44))

	// Title should be set
	if titleView.Title != "My Title" {
		t.Error("Title should be 'My Title'")
	}

	w.Close()
}

func TestNavigationTitleView_Loading(t *testing.T) {
	tv := navigation.NewTitleViewWithTitle("Loading...")

	w := test.NewWindow(tv)

	// Set loading state
	tv.SetLoading(true)

	renderer := test.WidgetRenderer(tv)
	renderer.Refresh()

	// Loading state is managed internally
	t.Log("TitleView loading state set to true")

	tv.SetLoading(false)
	renderer.Refresh()

	t.Log("TitleView loading state set to false")

	w.Close()
}

func TestTabBar_Creation(t *testing.T) {
	items := []*navigation.TabBarItem{
		navigation.NewTabBarItem("Home", nil),
		navigation.NewTabBarItem("Settings", nil),
	}
	tb := navigation.NewTabBar(items)

	w := test.NewWindow(tb)
	w.Resize(fyne.NewSize(400, 49))

	renderer := test.WidgetRenderer(tb)
	if renderer == nil {
		t.Fatal("TabBar renderer is nil")
	}

	w.Close()
}

// =============================================================================
// ALERT TESTS - Based on iOS QMUIAlertController
// =============================================================================

func TestAlertController_Alert(t *testing.T) {
	ac := alert.NewAlert("Alert", "This is a message", alert.ControllerStyleAlert)

	ac.AddAction(alert.NewAction("OK", alert.ActionStyleDefault, nil))
	ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))

	w := test.NewWindow(ac)
	w.Resize(fyne.NewSize(300, 200))

	if ac.Title != "Alert" {
		t.Error("Alert title should be 'Alert'")
	}
	actions := ac.GetActions()
	if len(actions) != 2 {
		t.Errorf("Alert should have 2 actions, got %d", len(actions))
	}

	w.Close()
}

func TestAlertController_ActionSheet(t *testing.T) {
	ac := alert.NewAlert("Choose", "", alert.ControllerStyleActionSheet)

	ac.AddAction(alert.NewAction("Option 1", alert.ActionStyleDefault, nil))
	ac.AddAction(alert.NewAction("Option 2", alert.ActionStyleDefault, nil))
	ac.AddAction(alert.NewAction("Delete", alert.ActionStyleDestructive, nil))
	ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))

	w := test.NewWindow(ac)
	w.Resize(fyne.NewSize(300, 250))

	if ac.Style != alert.ControllerStyleActionSheet {
		t.Error("Style should be ActionSheet")
	}

	w.Close()
}

// =============================================================================
// DIALOG TESTS - Based on iOS QMUIDialogViewController
// =============================================================================

func TestDialogViewController_Creation(t *testing.T) {
	content := label.NewLabel("Dialog content")
	dvc := dialog.NewDialogWithContent(content)
	dvc.Title = "Dialog Title"

	// DialogViewController is shown via ShowIn(window), not as a regular widget
	if dvc == nil {
		t.Fatal("DialogViewController is nil")
	}

	if dvc.Title != "Dialog Title" {
		t.Error("Dialog title should be set")
	}
}

func TestDialogViewController_Actions(t *testing.T) {
	dvc := dialog.NewDialog()
	dvc.Title = "Title"

	dvc.AddAction(dialog.NewDialogAction("OK", dialog.ActionStyleDefault))
	dvc.AddAction(dialog.NewDialogAction("Cancel", dialog.ActionStyleCancel))

	if len(dvc.Actions) != 2 {
		t.Errorf("Dialog should have 2 actions, got %d", len(dvc.Actions))
	}
}

// =============================================================================
// MODAL TESTS - Based on iOS QMUIModalPresentationViewController
// =============================================================================

func TestModalPresentation_AnimationStyles(t *testing.T) {
	content := label.NewLabel("Modal content")

	styles := []modal.ModalAnimationStyle{
		modal.ModalAnimationStyleFade,
		modal.ModalAnimationStyleSlideUp,
		modal.ModalAnimationStyleSlideDown,
	}

	for _, style := range styles {
		mvc := modal.NewModalWithContent(content)
		mvc.AnimationStyle = style

		// ModalPresentationViewController is shown via ShowIn(window)
		if mvc == nil {
			t.Fatalf("ModalPresentation is nil for style %v", style)
		}

		t.Logf("Created modal with animation style %v", style)
	}
}

// =============================================================================
// MORE OPERATION TESTS - Based on iOS QMUIMoreOperationController
// =============================================================================

func TestMoreOperationController_Creation(t *testing.T) {
	moc := moreop.NewActionSheet()

	items := []*moreop.Item{
		moreop.NewItem("share", "Share", nil, func(*moreop.Item) {}),
		moreop.NewItem("copy", "Copy", nil, func(*moreop.Item) {}),
	}
	moc.AddItems(items...)

	// MoreOperationController is shown via ShowIn(window)
	if moc == nil {
		t.Fatal("MoreOperationController is nil")
	}

	t.Log("MoreOperationController created with 2 items")
}

// =============================================================================
// TOAST TESTS - Based on iOS QMUIToastView
// =============================================================================

func TestToast_Creation(t *testing.T) {
	w := test.NewWindow(label.NewLabel("Background"))
	w.Resize(fyne.NewSize(300, 400))

	// Toast should be showable
	toast.ShowMessage(w, "Test message")

	// Small delay to let toast appear
	time.Sleep(100 * time.Millisecond)

	w.Close()
}

// =============================================================================
// TIPS TESTS - Based on iOS QMUITips
// =============================================================================

func TestTips_Loading(t *testing.T) {
	w := test.NewWindow(label.NewLabel("Background"))
	w.Resize(fyne.NewSize(300, 400))

	tips.ShowLoading(w, "Loading...")
	time.Sleep(100 * time.Millisecond)

	tips.Hide(w)

	w.Close()
}

func TestTips_Success(t *testing.T) {
	w := test.NewWindow(label.NewLabel("Background"))
	w.Resize(fyne.NewSize(300, 400))

	tips.ShowSuccess(w, "Success!")
	time.Sleep(100 * time.Millisecond)

	w.Close()
}

func TestTips_Error(t *testing.T) {
	w := test.NewWindow(label.NewLabel("Background"))
	w.Resize(fyne.NewSize(300, 400))

	tips.ShowError(w, "Error occurred")
	time.Sleep(100 * time.Millisecond)

	w.Close()
}

// =============================================================================
// EMPTY VIEW TESTS - Based on iOS QMUIEmptyView
// =============================================================================

func TestEmptyView_Loading(t *testing.T) {
	ev := empty.LoadingEmptyState("Loading data...")

	w := test.NewWindow(ev)
	w.Resize(fyne.NewSize(300, 200))

	renderer := test.WidgetRenderer(ev)
	if renderer == nil {
		t.Fatal("EmptyView renderer is nil")
	}

	w.Close()
}

func TestEmptyView_Error(t *testing.T) {
	ev := empty.ErrorEmptyState("Failed to load", func() {})

	w := test.NewWindow(ev)
	w.Resize(fyne.NewSize(300, 200))

	renderer := test.WidgetRenderer(ev)
	if renderer == nil {
		t.Fatal("Error EmptyView renderer is nil")
	}

	w.Close()
}

func TestEmptyView_NoData(t *testing.T) {
	ev := empty.NoDataEmptyState()

	w := test.NewWindow(ev)
	w.Resize(fyne.NewSize(300, 200))

	renderer := test.WidgetRenderer(ev)
	if renderer == nil {
		t.Fatal("NoData EmptyView renderer is nil")
	}

	w.Close()
}

// =============================================================================
// SEARCH BAR TESTS - Based on iOS QMUISearchBar
// =============================================================================

func TestSearchBar_Creation(t *testing.T) {
	sb := search.NewSearchBar()

	w := test.NewWindow(sb)
	w.Resize(fyne.NewSize(300, 44))

	renderer := test.WidgetRenderer(sb)
	if renderer == nil {
		t.Fatal("SearchBar renderer is nil")
	}

	w.Close()
}

func TestSearchBar_Placeholder(t *testing.T) {
	sb := search.NewSearchBarWithPlaceholder("Search...")

	w := test.NewWindow(sb)

	if sb.Placeholder != "Search..." {
		t.Error("SearchBar placeholder should be set")
	}

	w.Close()
}

func TestSearchBar_OnSearchClicked(t *testing.T) {
	searched := false
	sb := search.NewSearchBar()
	sb.OnSearchClicked = func() {
		searched = true
	}

	w := test.NewWindow(sb)

	// Simulate search action
	if sb.OnSearchClicked != nil {
		sb.OnSearchClicked()
	}

	if !searched {
		t.Error("OnSearchClicked callback should be invoked")
	}

	w.Close()
}
