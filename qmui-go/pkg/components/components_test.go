// Package components provides comprehensive tests for all QMUI Go components
// Uses Fyne's test framework for headless UI testing
package components

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/components/alert"
	"github.com/user/qmui-go/pkg/components/badge"
	"github.com/user/qmui-go/pkg/components/button"
	"github.com/user/qmui-go/pkg/components/checkbox"
	"github.com/user/qmui-go/pkg/components/collection"
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
	"github.com/user/qmui-go/pkg/theme"
)

var testApp fyne.App
var testWindow fyne.Window

func setupTest() {
	if testApp == nil {
		testApp = test.NewApp()
	}
	testWindow = testApp.NewWindow("Test")
	testWindow.Resize(fyne.NewSize(400, 600))
}

// ============ Button Tests ============

func TestButton(t *testing.T) {
	setupTest()

	tapped := false
	btn := button.NewButton("Test Button", func() {
		tapped = true
	})

	testWindow.SetContent(btn)
	test.Tap(btn)

	if !tapped {
		t.Error("Button tap handler was not called")
	}
}

func TestFillButton(t *testing.T) {
	setupTest()

	tapped := false
	btn := button.NewFillButton("Fill", color.RGBA{R: 100, G: 100, B: 255, A: 255}, func() {
		tapped = true
	})

	testWindow.SetContent(btn)
	test.Tap(btn)

	if !tapped {
		t.Error("FillButton tap handler was not called")
	}
}

func TestGhostButton(t *testing.T) {
	setupTest()

	tapped := false
	btn := button.NewGhostButton("Ghost", color.RGBA{R: 100, G: 100, B: 255, A: 255}, func() {
		tapped = true
	})

	testWindow.SetContent(btn)
	test.Tap(btn)

	if !tapped {
		t.Error("GhostButton tap handler was not called")
	}
}

func TestNavigationButton(t *testing.T) {
	setupTest()

	tapped := false
	btn := button.NewNavigationButton("Back", func() {
		tapped = true
	})

	testWindow.SetContent(btn)
	test.Tap(btn)

	if !tapped {
		t.Error("NavigationButton tap handler was not called")
	}
}

// ============ Label Tests ============

func TestLabel(t *testing.T) {
	setupTest()

	lbl := label.NewLabel("Test Label")
	lbl.ContentEdgeInsets = core.EdgeInsets{Top: 10, Left: 10, Bottom: 10, Right: 10}

	testWindow.SetContent(lbl)

	if lbl.Text != "Test Label" {
		t.Errorf("Expected 'Test Label', got '%s'", lbl.Text)
	}
}

func TestMarqueeLabel(t *testing.T) {
	setupTest()

	m := marquee.NewMarqueeLabel("Scrolling text")
	m.Speed = 50

	testWindow.SetContent(m)
	m.StartAnimation()

	// Let it animate briefly
	time.Sleep(100 * time.Millisecond)
	m.StopAnimation()
}

func TestBadgeLabel(t *testing.T) {
	setupTest()

	b := badge.NewBadgeLabel("99+")
	testWindow.SetContent(b)

	if b.Text != "99+" {
		t.Errorf("Expected '99+', got '%s'", b.Text)
	}
}

// ============ Text Input Tests ============

func TestTextField(t *testing.T) {
	setupTest()

	tf := textfield.NewTextFieldWithPlaceholder("Enter text")

	testWindow.SetContent(tf)

	tf.SetText("Hello World")
	if tf.Text != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", tf.Text)
	}
}

func TestTextFieldMaxLength(t *testing.T) {
	setupTest()

	tf := textfield.NewTextField()
	tf.MaximumTextLength = 5

	testWindow.SetContent(tf)

	tf.SetText("Hello World")
	// MaximumTextLength should limit the text
	if len(tf.Text) > 5 {
		t.Errorf("Text should be limited to 5 characters, got %d", len(tf.Text))
	}
}

func TestTextView(t *testing.T) {
	setupTest()

	tv := textview.NewTextView()
	tv.Placeholder = "Enter multi-line text"

	testWindow.SetContent(tv)

	tv.SetText("Line 1\nLine 2")
	if tv.Text != "Line 1\nLine 2" {
		t.Errorf("Expected multi-line text, got '%s'", tv.Text)
	}
}

func TestSearchBar(t *testing.T) {
	setupTest()

	sb := search.NewSearchBar()

	testWindow.SetContent(sb)

	sb.SetText("query")
	if sb.Text != "query" {
		t.Errorf("Expected 'query', got '%s'", sb.Text)
	}
}

// ============ Progress Tests ============

func TestPieProgressView(t *testing.T) {
	setupTest()

	pie := progress.NewPieProgressView()
	pie.Progress = 0.5
	pie.TintColor = color.RGBA{R: 0, G: 122, B: 255, A: 255}

	testWindow.SetContent(pie)

	pie.SetProgress(0.75)
	if pie.Progress != 0.75 {
		t.Errorf("Expected 0.75, got %f", pie.Progress)
	}
}

func TestCircularProgressView(t *testing.T) {
	setupTest()

	circ := progress.NewCircularProgressView()
	circ.Progress = 0.5
	circ.ShowsText = true

	testWindow.SetContent(circ)

	circ.SetProgress(0.9)
	if circ.Progress != 0.9 {
		t.Errorf("Expected 0.9, got %f", circ.Progress)
	}
}

func TestLinearProgressView(t *testing.T) {
	setupTest()

	lin := progress.NewLinearProgressView()
	lin.Progress = 0.3

	testWindow.SetContent(lin)

	lin.SetProgress(0.6)
	if lin.Progress != 0.6 {
		t.Errorf("Expected 0.6, got %f", lin.Progress)
	}
}

// ============ Selection Tests ============

func TestCheckbox(t *testing.T) {
	setupTest()

	changed := false
	cb := checkbox.NewCheckbox(func(selected bool) {
		changed = true
	})

	testWindow.SetContent(cb)

	cb.SetSelected(true)
	if !cb.Selected {
		t.Error("Checkbox should be checked")
	}
	if !changed {
		t.Error("OnChanged should have been called")
	}
}

func TestSegmentedControl(t *testing.T) {
	setupTest()

	changed := false
	sc := segmented.NewSegmentedControl([]string{"A", "B", "C"}, func(index int) {
		changed = true
	})

	testWindow.SetContent(sc)

	sc.SetSelectedIndex(2)
	if sc.SelectedIndex != 2 {
		t.Errorf("Expected index 2, got %d", sc.SelectedIndex)
	}
	_ = changed // prevent unused warning
}

// ============ Layout Tests ============

func TestFloatLayoutView(t *testing.T) {
	setupTest()

	fl := floatlayout.NewFloatLayoutView()
	fl.ItemSpacing = 5
	fl.LineSpacing = 5

	fl.AddItem(floatlayout.NewTagView("Tag1"))
	fl.AddItem(floatlayout.NewTagView("Tag2"))
	fl.AddItem(floatlayout.NewTagView("Tag3"))

	testWindow.SetContent(fl)
	// Widget renders without error
}

func TestGridView(t *testing.T) {
	setupTest()

	gv := grid.NewGridView(3)
	gv.RowSpacing = 4
	gv.ColumnSpacing = 4

	for i := 0; i < 6; i++ {
		gv.AddItem(grid.NewGridViewItem(widget.NewLabel("Item")))
	}

	testWindow.SetContent(gv)
	// Widget renders without error
}

func TestEmptyView(t *testing.T) {
	setupTest()

	ev := empty.LoadingEmptyView("Loading...")
	testWindow.SetContent(ev)

	// Just verify it creates without error
}

func TestTableView(t *testing.T) {
	setupTest()

	tv := table.NewTableView(table.TableViewStyleInsetGrouped)
	section := table.NewTableSection("Section 1")
	section.Cells = []*table.TableViewCell{
		table.NewTableViewCellWithTextAndDetail("Cell 1", "Detail 1"),
		table.NewTableViewCellWithTextAndDetail("Cell 2", "Detail 2"),
	}
	tv.Sections = []*table.TableSection{section}

	testWindow.SetContent(tv)

	if len(tv.Sections) != 1 {
		t.Errorf("Expected 1 section, got %d", len(tv.Sections))
	}
	if len(tv.Sections[0].Cells) != 2 {
		t.Errorf("Expected 2 cells, got %d", len(tv.Sections[0].Cells))
	}
}

// ============ Navigation Tests ============

func TestNavigationBar(t *testing.T) {
	setupTest()

	nav := navigation.NewNavigationBar()
	nav.SetTitleView(navigation.NewNavigationTitleViewWithTitle("Title"))
	nav.TintColor = color.RGBA{R: 0, G: 122, B: 255, A: 255}

	testWindow.SetContent(nav)
}

func TestNavigationTitleView(t *testing.T) {
	setupTest()

	tv := navigation.NewNavigationTitleViewWithTitleAndSubtitle("Title", "Subtitle")

	testWindow.SetContent(tv)

	if tv.Title != "Title" {
		t.Errorf("Expected 'Title', got '%s'", tv.Title)
	}
	if tv.Subtitle != "Subtitle" {
		t.Errorf("Expected 'Subtitle', got '%s'", tv.Subtitle)
	}
}

func TestTabBar(t *testing.T) {
	setupTest()

	items := []*navigation.TabBarItem{
		navigation.NewTabBarItem("Home", nil),
		navigation.NewTabBarItem("Settings", nil),
	}
	tb := navigation.NewTabBar(items)

	testWindow.SetContent(tb)

	tb.SetSelectedIndex(1)
	if tb.SelectedIndex != 1 {
		t.Errorf("Expected index 1, got %d", tb.SelectedIndex)
	}
}

// ============ Dialog Tests ============

func TestAlertController(t *testing.T) {
	setupTest()

	ac := alert.NewAlertController("Title", "Message", alert.ControllerStyleAlert)
	ac.AddAction(alert.NewAction("OK", alert.ActionStyleDefault, func(_ *alert.AlertController, _ *alert.Action) {
		// Action tapped
	}))

	testWindow.SetContent(container.NewVBox())
	ac.ShowIn(testWindow)

	// Hide after brief delay
	time.Sleep(50 * time.Millisecond)
	ac.Hide()
}

func TestToast(t *testing.T) {
	setupTest()
	testWindow.SetContent(container.NewVBox())

	toast.ShowMessage(testWindow, "Test toast message")
	time.Sleep(50 * time.Millisecond)
}

func TestTips(t *testing.T) {
	setupTest()
	testWindow.SetContent(container.NewVBox())

	tp := tips.NewTips(testWindow)
	tp.ShowLoading("Loading...")
	time.Sleep(50 * time.Millisecond)
	tp.HideCurrent()

	tp.ShowSuccess("Success!")
	time.Sleep(50 * time.Millisecond)
	tp.HideCurrent()

	tp.ShowError("Error!")
	time.Sleep(50 * time.Millisecond)
	tp.HideCurrent()
}

func TestDialogViewController(t *testing.T) {
	setupTest()
	testWindow.SetContent(container.NewVBox())

	dvc := dialog.NewDialogViewController()
	dvc.Title = "Dialog Title"
	dvc.Message = "Dialog message"
	dvc.AddAction(dialog.NewDialogAction("OK", dialog.ActionStyleDefault))

	dvc.Show(testWindow)
	time.Sleep(50 * time.Millisecond)
	dvc.Dismiss()
}

func TestModalPresentation(t *testing.T) {
	setupTest()
	testWindow.SetContent(container.NewVBox())

	content := widget.NewLabel("Modal Content")
	mvc := modal.PresentModalFromBottom(testWindow, content)

	time.Sleep(50 * time.Millisecond)
	mvc.Dismiss()
}

func TestMoreOperationController(t *testing.T) {
	setupTest()
	testWindow.SetContent(container.NewVBox())

	items := []*moreop.MoreOperationItem{
		moreop.NewMoreOperationItem("share", "Share", nil, nil),
		moreop.NewMoreOperationItem("copy", "Copy", nil, nil),
	}

	ctrl := moreop.NewMoreOperationController()
	ctrl.AddItems(items...)
	ctrl.Show(testWindow)

	time.Sleep(50 * time.Millisecond)
	ctrl.Dismiss()
}

// ============ Popup Tests ============

func TestPopupMenu(t *testing.T) {
	setupTest()
	testWindow.SetContent(container.NewVBox())

	items := []*popup.MenuItem{
		popup.NewMenuItem("Edit", func(_ *popup.MenuItem) {
			// Item tapped
		}),
		popup.NewMenuItem("Delete", nil),
	}

	pm := popup.NewPopupMenuViewWithItems(items)
	pm.Show(testWindow, fyne.NewPos(100, 100))

	time.Sleep(50 * time.Millisecond)
	pm.Hide()
}

// ============ Special Components Tests ============

func TestEmotionView(t *testing.T) {
	setupTest()

	ev := emotion.NewEmotionView()
	ev.OnEmotionSelected = func(e *emotion.Emotion) {
		// Emotion selected
		_ = e
	}

	testWindow.SetContent(ev)
}

func TestPagingLayout(t *testing.T) {
	setupTest()

	pl := collection.NewPagingLayout()
	pl.AddPage(widget.NewLabel("Page 1"))
	pl.AddPage(widget.NewLabel("Page 2"))
	pl.AddPage(widget.NewLabel("Page 3"))

	testWindow.SetContent(pl)

	pl.SetCurrentPage(1)
	if pl.CurrentPage != 1 {
		t.Errorf("Expected page 1, got %d", pl.CurrentPage)
	}

	pl.SetCurrentPage(2)
	if pl.CurrentPage != 2 {
		t.Errorf("Expected page 2, got %d", pl.CurrentPage)
	}
}

// ============ Theme Tests ============

func TestThemeManager(t *testing.T) {
	tm := theme.SharedThemeManager()

	// Switch to dark
	tm.SetCurrentTheme(theme.ThemeIdentifierDark)
	if tm.CurrentTheme().Identifier != theme.ThemeIdentifierDark {
		t.Errorf("Expected dark theme, got %s", tm.CurrentTheme().Identifier)
	}

	// Switch to light
	tm.SetCurrentTheme(theme.ThemeIdentifierDefault)
	if tm.CurrentTheme().Identifier != theme.ThemeIdentifierDefault {
		t.Errorf("Expected default theme, got %s", tm.CurrentTheme().Identifier)
	}
}

func TestThemeChangeListener(t *testing.T) {
	tm := theme.SharedThemeManager()

	changed := false
	tm.AddThemeChangeListener(func(th *theme.Theme) {
		changed = true
	})

	tm.SetCurrentTheme(theme.ThemeIdentifierDark)
	if !changed {
		t.Error("Theme change listener was not called")
	}

	// Reset
	tm.SetCurrentTheme(theme.ThemeIdentifierDefault)
}

// ============ Configuration Tests ============

func TestConfiguration(t *testing.T) {
	cfg := core.SharedConfiguration()

	// Verify default colors exist
	if cfg.BlueColor == nil {
		t.Error("BlueColor should not be nil")
	}
	if cfg.RedColor == nil {
		t.Error("RedColor should not be nil")
	}
	if cfg.GreenColor == nil {
		t.Error("GreenColor should not be nil")
	}
}

// ============ Comprehensive Widget Cycle Test ============

func TestAllWidgetsCycle(t *testing.T) {
	setupTest()

	// This test creates every widget type and verifies they render without panic

	widgets := []fyne.CanvasObject{
		// Buttons
		button.NewButton("Button", nil),
		button.NewFillButton("Fill", color.RGBA{R: 100, G: 100, B: 255, A: 255}, nil),
		button.NewGhostButton("Ghost", color.RGBA{R: 100, G: 100, B: 255, A: 255}, nil),
		button.NewNavigationButton("Nav", nil),
		button.NewToolbarButton("Tool", nil),

		// Labels
		label.NewLabel("Label"),
		marquee.NewMarqueeLabel("Marquee"),
		badge.NewBadgeLabel("99"),

		// Text Input
		textfield.NewTextField(),
		textview.NewTextView(),
		search.NewSearchBar(),

		// Progress
		progress.NewPieProgressView(),
		progress.NewCircularProgressView(),
		progress.NewLinearProgressView(),

		// Selection
		checkbox.NewCheckbox(func(bool) {}),
		segmented.NewSegmentedControl([]string{"A", "B"}, func(int) {}),

		// Layout
		floatlayout.NewFloatLayoutView(),
		grid.NewGridView(2),
		empty.NewEmptyView(),
		table.NewTableView(table.TableViewStyleGrouped),

		// Navigation
		navigation.NewNavigationBar(),
		navigation.NewTabBar([]*navigation.TabBarItem{}),
		navigation.NewNavigationTitleViewWithTitle("Title"),

		// Special
		emotion.NewEmotionView(),
		collection.NewPagingLayout(),
	}

	for i, w := range widgets {
		testWindow.SetContent(w)
		testWindow.Canvas().Refresh(w)
		t.Logf("Widget %d rendered successfully", i)
	}
}

// ============ Integration Test ============

func TestFullAppFlow(t *testing.T) {
	setupTest()

	// Simulate a complete app flow

	// 1. Show loading
	testWindow.SetContent(empty.LoadingEmptyView("Loading..."))
	time.Sleep(50 * time.Millisecond)

	// 2. Show main content with tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Home", widget.NewLabel("Home")),
		container.NewTabItem("Settings", widget.NewLabel("Settings")),
	)
	testWindow.SetContent(tabs)
	time.Sleep(50 * time.Millisecond)

	// 3. Switch tabs
	tabs.SelectIndex(1)
	time.Sleep(50 * time.Millisecond)

	// 4. Show a dialog
	dvc := dialog.NewDialogViewController()
	dvc.Title = "Confirm"
	dvc.Show(testWindow)
	time.Sleep(50 * time.Millisecond)
	dvc.Dismiss()

	// 5. Show toast
	toast.ShowMessage(testWindow, "Action completed")
	time.Sleep(50 * time.Millisecond)

	t.Log("Full app flow completed successfully")
}
