// Automated UI test/tour for QMUI Go demo
// This program cycles through every component and interacts with them
// Suitable for video recording to showcase all features
package main

import (
	"fmt"
	"image/color"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	fynetheme "fyne.io/fyne/v2/theme"
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
	"github.com/user/qmui-go/pkg/console"
	"github.com/user/qmui-go/pkg/core"
	"github.com/user/qmui-go/pkg/theme"
)

// Colors
var (
	qmuiCyan       = color.RGBA{R: 49, G: 189, B: 243, A: 255}
	qmuiBackground = color.RGBA{R: 246, G: 247, B: 249, A: 255}
	qmuiWhite      = color.White
	qmuiGrayText   = color.RGBA{R: 134, G: 144, B: 156, A: 255}
	qmuiDarkText   = color.RGBA{R: 51, G: 51, B: 51, A: 255}
)

var (
	mainWindow fyne.Window
	mainApp    fyne.App
	tabs       *container.AppTabs
	statusText *canvas.Text
	tourActive bool
)

// Tour step timing
const (
	shortPause  = 800 * time.Millisecond
	mediumPause = 1200 * time.Millisecond
	longPause   = 2000 * time.Millisecond
)

func main() {
	mainApp = app.New()
	mainWindow = mainApp.NewWindow("QMUI Go - Automated UI Tour")
	mainWindow.Resize(fyne.NewSize(960, 800))

	// Make this the master window - closing it quits the app
	mainWindow.SetMaster()

	// Handle OS signals (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		tourActive = false
		mainApp.Quit()
	}()

	// Ensure clean exit when window is closed
	mainWindow.SetCloseIntercept(func() {
		fmt.Println("Window closed, exiting...")
		tourActive = false
		mainApp.Quit()
	})

	// Center window if possible
	if drv, ok := mainApp.Driver().(desktop.Driver); ok {
		_ = drv // Could use for positioning
	}

	_ = core.SharedConfiguration()

	showTourUI()
	mainApp.Run()
}

func showTourUI() {
	// Status bar at top showing current action
	statusText = canvas.NewText("Ready to start automated tour", qmuiWhite)
	statusText.TextSize = 12
	statusBg := canvas.NewRectangle(color.RGBA{R: 40, G: 40, B: 40, A: 255})
	statusBg.SetMinSize(fyne.NewSize(0, 28))
	statusBar := container.NewStack(statusBg, container.NewCenter(statusText))

	// Main content tabs
	tabs = container.NewAppTabs(
		container.NewTabItem("Components", createComponentsTab()),
		container.NewTabItem("Buttons", createButtonsTab()),
		container.NewTabItem("Text Input", createTextInputTab()),
		container.NewTabItem("Progress", createProgressTab()),
		container.NewTabItem("Dialogs", createDialogsTab()),
		container.NewTabItem("Navigation", createNavigationTab()),
		container.NewTabItem("Advanced", createAdvancedTab()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	// Control buttons
	startBtn := widget.NewButtonWithIcon("Start Tour", fynetheme.MediaPlayIcon(), func() {
		if !tourActive {
			go runAutomatedTour()
		}
	})
	startBtn.Importance = widget.HighImportance

	stopBtn := widget.NewButtonWithIcon("Stop", fynetheme.MediaStopIcon(), func() {
		tourActive = false
		setStatus("Tour stopped")
	})

	controls := container.NewHBox(layout.NewSpacer(), startBtn, stopBtn, layout.NewSpacer())
	controlBar := container.NewPadded(controls)

	navBar := createNavBar("QMUI Go - Component Tour")
	bg := canvas.NewRectangle(qmuiBackground)

	content := container.NewBorder(
		container.NewVBox(navBar, statusBar, controlBar),
		nil, nil, nil,
		container.NewStack(bg, tabs),
	)

	mainWindow.SetContent(content)
	mainWindow.Show()
}

func setStatus(msg string) {
	fmt.Println(msg)
	statusText.Text = msg
	statusText.Refresh()
}

func createNavBar(title string) fyne.CanvasObject {
	bg := canvas.NewRectangle(qmuiCyan)
	bg.SetMinSize(fyne.NewSize(0, 44))

	titleText := canvas.NewText(title, qmuiWhite)
	titleText.TextSize = 17
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewStack(bg, container.NewCenter(titleText))
}

// ============ Automated Tour ============

func runAutomatedTour() {
	tourActive = true
	setStatus("Starting automated tour...")
	time.Sleep(longPause)

	// Tab 0: Components
	tourComponentsTab()
	if !tourActive {
		return
	}

	// Tab 1: Buttons
	tourButtonsTab()
	if !tourActive {
		return
	}

	// Tab 2: Text Input
	tourTextInputTab()
	if !tourActive {
		return
	}

	// Tab 3: Progress
	tourProgressTab()
	if !tourActive {
		return
	}

	// Tab 4: Dialogs
	tourDialogsTab()
	if !tourActive {
		return
	}

	// Tab 5: Navigation
	tourNavigationTab()
	if !tourActive {
		return
	}

	// Tab 6: Advanced
	tourAdvancedTab()
	if !tourActive {
		return
	}

	// Final
	setStatus("Tour complete! All components demonstrated.")
	tourActive = false
}

func selectTab(index int) {
	tabs.SelectIndex(index)
	time.Sleep(shortPause)
}

// ============ Tour: Components Tab ============

func tourComponentsTab() {
	setStatus("Tab 1/7: Components - Labels & Badges")
	selectTab(0)
	time.Sleep(longPause)

	setStatus("Showing: marquee.MarqueeLabel - animated scrolling text")
	time.Sleep(longPause)

	setStatus("Showing: badge.BadgeLabel - notification badges")
	time.Sleep(mediumPause)

	setStatus("Showing: label.Label - padded labels with edge insets")
	time.Sleep(mediumPause)

	setStatus("Showing: floatlayout.FloatLayoutView - tag cloud")
	time.Sleep(mediumPause)

	setStatus("Showing: grid.GridView - colored grid")
	time.Sleep(mediumPause)

	setStatus("Showing: empty.EmptyView - loading states")
	time.Sleep(mediumPause)

	setStatus("Showing: table.TableView - iOS-style lists")
	time.Sleep(longPause)
}

// ============ Tour: Buttons Tab ============

var (
	demoBtnStandard *button.Button
	demoBtnFill     *button.FillButton
	demoBtnGhost    *button.GhostButton
)

func tourButtonsTab() {
	setStatus("Tab 2/7: Buttons - Interactive button variants")
	selectTab(1)
	time.Sleep(longPause)

	setStatus("Tapping: button.Button (standard)")
	if demoBtnStandard != nil {
		simulateTap(demoBtnStandard)
	}
	time.Sleep(mediumPause)

	setStatus("Tapping: button.FillButton (solid)")
	if demoBtnFill != nil {
		simulateTap(demoBtnFill)
	}
	time.Sleep(mediumPause)

	setStatus("Tapping: button.GhostButton (outlined)")
	if demoBtnGhost != nil {
		simulateTap(demoBtnGhost)
	}
	time.Sleep(mediumPause)

	setStatus("Showing: Color variants (Red, Green, Blue)")
	time.Sleep(longPause)
}

// ============ Tour: Text Input Tab ============

var (
	demoTextField *textfield.TextField
	demoTextView  *textview.TextView
	demoSearchBar *search.SearchBar
)

func tourTextInputTab() {
	setStatus("Tab 3/7: Text Input - Fields and search")
	selectTab(2)
	time.Sleep(longPause)

	setStatus("Typing in: textfield.TextField")
	if demoTextField != nil {
		simulateTyping(demoTextField, "Hello QMUI!")
	}
	time.Sleep(mediumPause)

	setStatus("Typing in: textview.TextView (multiline)")
	if demoTextView != nil {
		simulateTypingTextView(demoTextView, "Multi-line\ntext input")
	}
	time.Sleep(mediumPause)

	setStatus("Searching in: search.SearchBar")
	if demoSearchBar != nil {
		demoSearchBar.SetText("QMUI Go")
		demoSearchBar.Refresh()
	}
	time.Sleep(longPause)
}

// ============ Tour: Progress Tab ============

var (
	demoPie      *progress.PieProgressView
	demoCircular *progress.CircularProgressView
	demoLinear   *progress.LinearProgressView
)

func tourProgressTab() {
	setStatus("Tab 4/7: Progress - Animated indicators")
	selectTab(3)
	time.Sleep(longPause)

	setStatus("Animating: progress.PieProgressView")
	animateProgress(demoPie, demoCircular, demoLinear, 0.33)
	time.Sleep(mediumPause)

	setStatus("Animating: progress.CircularProgressView")
	animateProgress(demoPie, demoCircular, demoLinear, 0.66)
	time.Sleep(mediumPause)

	setStatus("Animating: progress.LinearProgressView")
	animateProgress(demoPie, demoCircular, demoLinear, 1.0)
	time.Sleep(longPause)

	// Reset
	animateProgress(demoPie, demoCircular, demoLinear, 0)
}

func animateProgress(pie *progress.PieProgressView, circ *progress.CircularProgressView, lin *progress.LinearProgressView, target float64) {
	if pie != nil {
		pie.SetProgress(target)
	}
	if circ != nil {
		circ.SetProgress(target)
	}
	if lin != nil {
		lin.SetProgress(target)
	}
}

// ============ Tour: Dialogs Tab ============

func tourDialogsTab() {
	setStatus("Tab 5/7: Dialogs - Toasts, Alerts, Popups")
	selectTab(4)
	time.Sleep(longPause)

	setStatus("Showing: toast.ToastView")
	toast.ShowMessage(mainWindow, "This is a toast message!")
	time.Sleep(longPause)

	setStatus("Showing: tips.Tips (loading)")
	t := tips.NewTips(mainWindow)
	t.ShowLoading("Loading...")
	time.Sleep(mediumPause)
	t.HideCurrent()

	setStatus("Showing: tips.Tips (success)")
	t.ShowSuccess("Success!")
	time.Sleep(mediumPause)
	t.HideCurrent()

	setStatus("Showing: alert.AlertController")
	showDemoAlert()
	time.Sleep(longPause)

	setStatus("Showing: alert.AlertController (action sheet)")
	showDemoActionSheet()
	time.Sleep(longPause)

	setStatus("Showing: dialog.DialogViewController")
	showDemoDialog()
	time.Sleep(longPause)

	setStatus("Showing: popup.PopupMenu (context menu)")
	showDemoPopupMenu()
	time.Sleep(longPause)

	setStatus("Showing: moreop.MoreOperationController")
	showDemoMoreOp()
	time.Sleep(longPause)

	setStatus("Showing: modal.ModalPresentationViewController")
	showDemoModal()
	time.Sleep(longPause)
}

func showDemoAlert() {
	ac := alert.NewAlertController("Alert Demo", "This demonstrates alert.AlertController", alert.ControllerStyleAlert)
	ac.AddAction(alert.NewAction("OK", alert.ActionStyleDefault, func(_ *alert.AlertController, _ *alert.Action) {}))
	ac.ShowIn(mainWindow)

	// Auto-dismiss after delay
	go func() {
		time.Sleep(mediumPause)
		ac.Hide()
	}()
}

func showDemoActionSheet() {
	ac := alert.NewAlertController("Action Sheet", "Choose an option", alert.ControllerStyleActionSheet)
	ac.AddAction(alert.NewAction("Option 1", alert.ActionStyleDefault, nil))
	ac.AddAction(alert.NewAction("Option 2", alert.ActionStyleDefault, nil))
	ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))
	ac.ShowIn(mainWindow)

	go func() {
		time.Sleep(mediumPause)
		ac.Hide()
	}()
}

func showDemoDialog() {
	dvc := dialog.NewDialogViewController()
	dvc.Title = "Dialog Demo"
	dvc.Message = "This is dialog.DialogViewController with custom content"
	dvc.AddAction(dialog.NewDialogActionWithHandler("Close", dialog.ActionStyleCancel, func(_ *dialog.DialogAction) {
		dvc.Dismiss()
	}))
	dvc.Show(mainWindow)

	go func() {
		time.Sleep(mediumPause)
		dvc.Dismiss()
	}()
}

func showDemoPopupMenu() {
	items := []*popup.MenuItem{
		popup.NewMenuItem("Edit", nil),
		popup.NewMenuItem("Copy", nil),
		popup.NewMenuItem("Delete", nil),
	}
	pm := popup.NewPopupMenuViewWithItems(items)
	pm.Show(mainWindow, fyne.NewPos(200, 400))

	go func() {
		time.Sleep(mediumPause)
		pm.Hide()
	}()
}

func showDemoMoreOp() {
	items := []*moreop.MoreOperationItem{
		moreop.NewMoreOperationItem("share", "Share", nil, nil),
		moreop.NewMoreOperationItem("copy", "Copy", nil, nil),
		moreop.NewMoreOperationItem("save", "Save", nil, nil),
		moreop.NewMoreOperationItem("delete", "Delete", nil, nil),
	}
	ctrl := moreop.NewMoreOperationController()
	ctrl.AddItems(items...)
	ctrl.Show(mainWindow)

	go func() {
		time.Sleep(mediumPause)
		ctrl.Dismiss()
	}()
}

func showDemoModal() {
	content := container.NewVBox(
		widget.NewLabel("Modal Content"),
		widget.NewLabel("This slides up from bottom"),
	)
	mvc := modal.PresentModalFromBottom(mainWindow, container.NewPadded(content))

	go func() {
		time.Sleep(mediumPause)
		mvc.Dismiss()
	}()
}

// ============ Tour: Navigation Tab ============

var (
	demoSegmented *segmented.SegmentedControl
	demoCheckbox1 *checkbox.Checkbox
	demoCheckbox2 *checkbox.Checkbox
)

func tourNavigationTab() {
	setStatus("Tab 6/7: Navigation - Bars and controls")
	selectTab(5)
	time.Sleep(longPause)

	setStatus("Showing: navigation.NavigationBar")
	time.Sleep(mediumPause)

	setStatus("Showing: navigation.TabBar")
	time.Sleep(mediumPause)

	setStatus("Toggling: segmented.SegmentedControl")
	if demoSegmented != nil {
		demoSegmented.SetSelectedIndex(1)
		time.Sleep(shortPause)
		demoSegmented.SetSelectedIndex(2)
		time.Sleep(shortPause)
		demoSegmented.SetSelectedIndex(0)
	}
	time.Sleep(mediumPause)

	setStatus("Toggling: checkbox.Checkbox")
	if demoCheckbox1 != nil {
		demoCheckbox1.SetSelected(true)
		time.Sleep(shortPause)
	}
	if demoCheckbox2 != nil {
		demoCheckbox2.SetSelected(true)
	}
	time.Sleep(longPause)
}

// ============ Tour: Advanced Tab ============

var demoPaging *collection.PagingLayout

func tourAdvancedTab() {
	setStatus("Tab 7/7: Advanced - Theme, Console, Special")
	selectTab(6)
	time.Sleep(longPause)

	setStatus("Showing: emotion.EmotionView (emoji picker)")
	time.Sleep(mediumPause)

	setStatus("Swiping: collection.PagingLayout")
	if demoPaging != nil {
		demoPaging.SetCurrentPage(1)
		time.Sleep(shortPause)
		demoPaging.SetCurrentPage(2)
		time.Sleep(shortPause)
		demoPaging.SetCurrentPage(0)
	}
	time.Sleep(mediumPause)

	setStatus("Switching to: Dark Theme")
	theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDark)
	time.Sleep(longPause)

	setStatus("Switching to: Light Theme")
	theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDefault)
	time.Sleep(mediumPause)

	setStatus("Opening: console.Console")
	console.SharedConsole().ShowIn(mainWindow)
	time.Sleep(longPause)
	console.SharedConsole().Hide()
	time.Sleep(mediumPause)
}

// ============ Tab Content Creators ============

func createComponentsTab() fyne.CanvasObject {
	cards := []fyne.CanvasObject{
		createCard("marquee.MarqueeLabel", "Scrolling text animation", createMarquee()),
		createCard("badge.BadgeLabel", "Notification badges", createBadges()),
		createCard("label.Label", "Label with edge insets", createLabel()),
		createCard("floatlayout.FloatLayoutView", "Tag cloud layout", createFloatLayout()),
		createCard("grid.GridView", "Grid arrangement", createGrid()),
		createCard("empty.EmptyView", "Loading state", createEmpty()),
		createCard("table.TableView", "iOS-style grouped list", createTable()),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createButtonsTab() fyne.CanvasObject {
	demoBtnStandard = button.NewButton("Standard Button", func() {
		toast.ShowMessage(mainWindow, "Standard button tapped!")
	})

	demoBtnFill = button.NewFillButton("Fill Button", qmuiCyan, func() {
		toast.ShowMessage(mainWindow, "Fill button tapped!")
	})

	demoBtnGhost = button.NewGhostButton("Ghost Button", qmuiCyan, func() {
		toast.ShowMessage(mainWindow, "Ghost button tapped!")
	})

	navBtn := button.NewNavigationButton("< Back", func() {
		toast.ShowMessage(mainWindow, "Navigation button tapped!")
	})

	cfg := core.SharedConfiguration()
	colorBtns := container.NewHBox(
		button.NewFillButton("Red", cfg.RedColor, nil),
		button.NewFillButton("Green", cfg.GreenColor, nil),
		button.NewFillButton("Blue", cfg.BlueColor, nil),
	)

	cards := []fyne.CanvasObject{
		createCard("button.Button", "Standard tappable button", demoBtnStandard),
		createCard("button.FillButton", "Solid filled button", demoBtnFill),
		createCard("button.GhostButton", "Outlined button", demoBtnGhost),
		createCard("button.NavigationButton", "Navigation bar button", navBtn),
		createCard("Color Variants", "Multiple tint colors", colorBtns),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createTextInputTab() fyne.CanvasObject {
	demoTextField = textfield.NewTextFieldWithPlaceholder("Enter text here...")
	demoTextField.PlaceholderColor = qmuiGrayText

	demoTextView = textview.NewTextView()
	demoTextView.PlaceHolder = "Multi-line text input..."

	demoSearchBar = search.NewSearchBar()
	demoSearchBar.Placeholder = "Search..."

	cards := []fyne.CanvasObject{
		createCard("textfield.TextField", "Single-line input", demoTextField),
		createCard("textview.TextView", "Multi-line input", demoTextView),
		createCard("search.SearchBar", "Search with suggestions", demoSearchBar),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createProgressTab() fyne.CanvasObject {
	demoPie = progress.NewPieProgressView()
	demoPie.Progress = 0
	demoPie.TintColor = qmuiCyan

	demoCircular = progress.NewCircularProgressView()
	demoCircular.Progress = 0
	demoCircular.TintColor = color.RGBA{R: 52, G: 199, B: 89, A: 255}
	demoCircular.ShowsText = true

	demoLinear = progress.NewLinearProgressView()
	demoLinear.Progress = 0
	demoLinear.TintColor = color.RGBA{R: 255, G: 149, B: 0, A: 255}

	cards := []fyne.CanvasObject{
		createCard("progress.PieProgressView", "Pie chart progress", demoPie),
		createCard("progress.CircularProgressView", "Ring with percentage", demoCircular),
		createCard("progress.LinearProgressView", "Horizontal bar", demoLinear),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createDialogsTab() fyne.CanvasObject {
	toastBtn := widget.NewButton("Show Toast", func() {
		toast.ShowMessage(mainWindow, "Toast message!")
	})

	tipsBtn := widget.NewButton("Show Tips", func() {
		t := tips.NewTips(mainWindow)
		t.ShowSuccess("Operation completed!")
		go func() {
			time.Sleep(2 * time.Second)
			t.HideCurrent()
		}()
	})

	alertBtn := widget.NewButton("Show Alert", func() {
		ac := alert.NewAlertController("Alert", "This is an alert dialog", alert.ControllerStyleAlert)
		ac.AddAction(alert.NewAction("OK", alert.ActionStyleDefault, nil))
		ac.ShowIn(mainWindow)
	})

	dialogBtn := widget.NewButton("Show Dialog", func() {
		dialog.ShowConfirmDialog(mainWindow, "Confirm", "Are you sure?", func() {
			toast.ShowMessage(mainWindow, "Confirmed!")
		}, nil)
	})

	popupBtn := widget.NewButton("Show Popup Menu", func() {
		items := []*popup.MenuItem{
			popup.NewMenuItem("Edit", func(_ *popup.MenuItem) { toast.ShowMessage(mainWindow, "Edit") }),
			popup.NewMenuItem("Delete", func(_ *popup.MenuItem) { toast.ShowMessage(mainWindow, "Delete") }),
		}
		popup.ContextMenu(mainWindow, fyne.NewPos(200, 300), items)
	})

	moreOpBtn := widget.NewButton("Show More Operations", func() {
		items := []*moreop.MoreOperationItem{
			moreop.NewMoreOperationItem("share", "Share", nil, func(_ *moreop.MoreOperationItem) {
				toast.ShowMessage(mainWindow, "Share tapped")
			}),
			moreop.NewMoreOperationItem("copy", "Copy", nil, func(_ *moreop.MoreOperationItem) {
				toast.ShowMessage(mainWindow, "Copy tapped")
			}),
		}
		moreop.ShowOperationSheet(mainWindow, items, "Cancel")
	})

	modalBtn := widget.NewButton("Show Modal", func() {
		content := container.NewVBox(
			widget.NewLabel("Modal Content"),
			widget.NewButton("Close", func() {}),
		)
		modal.PresentModalFromBottom(mainWindow, container.NewPadded(content))
	})

	cards := []fyne.CanvasObject{
		createCard("toast.ToastView", "Brief messages", toastBtn),
		createCard("tips.Tips", "Loading/Success/Error", tipsBtn),
		createCard("alert.AlertController", "Alert dialogs", alertBtn),
		createCard("dialog.DialogViewController", "Custom dialogs", dialogBtn),
		createCard("popup.PopupMenu", "Context menus", popupBtn),
		createCard("moreop.MoreOperationController", "Action grid", moreOpBtn),
		createCard("modal.ModalPresentationViewController", "Slide-up modal", modalBtn),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createNavigationTab() fyne.CanvasObject {
	// Navigation bar
	navBar := navigation.NewNavigationBar()
	navBar.SetTitleView(navigation.NewNavigationTitleViewWithTitle("Navigation Bar"))
	navBar.TintColor = qmuiCyan

	// Tab bar
	tabBarItems := []*navigation.TabBarItem{
		navigation.NewTabBarItem("Home", nil),
		navigation.NewTabBarItem("Search", nil),
		navigation.NewTabBarItem("Profile", nil),
	}
	tabBar := navigation.NewTabBar(tabBarItems)
	tabBar.SetSelectedIndex(0)

	// Segmented control
	demoSegmented = segmented.NewSegmentedControl([]string{"Day", "Week", "Month"}, func(index int) {})
	demoSegmented.SetSelectedIndex(0)

	// Checkboxes
	demoCheckbox1 = checkbox.NewCheckboxWithLabel("Option 1", func(selected bool) {})
	demoCheckbox2 = checkbox.NewCheckboxWithLabel("Option 2", func(selected bool) {})
	checkboxes := container.NewHBox(demoCheckbox1, demoCheckbox2)

	cards := []fyne.CanvasObject{
		createCard("navigation.NavigationBar", "App navigation bar", navBar),
		createCard("navigation.TabBar", "Bottom tab bar", tabBar),
		createCard("segmented.SegmentedControl", "Segmented selector", demoSegmented),
		createCard("checkbox.Checkbox", "Selection checkboxes", checkboxes),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createAdvancedTab() fyne.CanvasObject {
	// Emotion view
	emotionView := emotion.NewEmotionView()

	// Paging layout
	demoPaging = collection.NewPagingLayout()
	for i := 0; i < 4; i++ {
		colors := []color.Color{qmuiCyan, color.RGBA{R: 255, G: 100, B: 100, A: 255},
			color.RGBA{R: 100, G: 200, B: 100, A: 255}, color.RGBA{R: 200, G: 150, B: 255, A: 255}}
		rect := canvas.NewRectangle(colors[i%4])
		rect.CornerRadius = 8
		lbl := canvas.NewText(fmt.Sprintf("Page %d", i+1), qmuiWhite)
		lbl.TextSize = 18
		page := container.NewStack(rect, container.NewCenter(lbl))
		demoPaging.AddPage(page)
	}

	// Theme switcher
	themeRow := container.NewHBox(
		widget.NewButton("Light", func() {
			theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDefault)
		}),
		widget.NewButton("Dark", func() {
			theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDark)
		}),
	)

	// Console toggle
	consoleBtn := widget.NewButton("Toggle Console", func() {
		console.SharedConsole().Toggle(mainWindow)
	})

	cards := []fyne.CanvasObject{
		createCard("emotion.EmotionView", "Emoji picker grid", emotionView),
		createCard("collection.PagingLayout", "Swipeable pages", demoPaging),
		createCard("theme.ThemeManager", "Hot-switchable themes", themeRow),
		createCard("console.Console", "Debug console", consoleBtn),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

// ============ Component Creators ============

func createMarquee() fyne.CanvasObject {
	// Use long text to ensure scrolling even on wide windows
	m := marquee.NewMarqueeLabel("This text scrolls continuously across the screen - MarqueeLabel provides smooth animated scrolling text perfect for news tickers, announcements, and attention-grabbing displays!")
	m.Speed = 50
	m.TextColor = qmuiDarkText
	m.AutoScrollWhenFits = false // Always scroll for demo purposes
	m.PauseDuration = 500 * time.Millisecond // Shorter pause for demo
	m.StartAnimation()

	// Use a max-width wrapper to constrain the marquee
	m.Resize(fyne.NewSize(300, 30))
	return m
}

func createBadges() fyne.CanvasObject {
	return container.NewHBox(
		badge.NewBadgeLabel("99+"),
		badge.NewBadgeLabel("NEW"),
		badge.NewBadgeLabel("5"),
	)
}

func createLabel() fyne.CanvasObject {
	lbl := label.NewLabel("Padded Label")
	lbl.ContentEdgeInsets = core.EdgeInsets{Top: 8, Left: 16, Bottom: 8, Right: 16}
	bg := canvas.NewRectangle(color.RGBA{R: 230, G: 245, B: 255, A: 255})
	bg.CornerRadius = 4
	return container.NewStack(bg, lbl)
}

func createFloatLayout() fyne.CanvasObject {
	fl := floatlayout.NewFloatLayoutView()
	fl.ItemSpacing = 6
	fl.LineSpacing = 6
	for _, tag := range []string{"Go", "Fyne", "QMUI", "iOS", "Cross-Platform"} {
		fl.AddItem(floatlayout.NewTagView(tag))
	}
	return fl
}

func createGrid() fyne.CanvasObject {
	gv := grid.NewGridView(4)
	gv.RowSpacing = 4
	gv.ColumnSpacing = 4
	colors := []color.Color{qmuiCyan, color.RGBA{R: 255, G: 100, B: 100, A: 255},
		color.RGBA{R: 100, G: 200, B: 100, A: 255}, color.RGBA{R: 200, G: 150, B: 255, A: 255}}
	for _, c := range colors {
		rect := canvas.NewRectangle(c)
		rect.CornerRadius = 4
		rect.SetMinSize(fyne.NewSize(30, 30))
		gv.AddItem(grid.NewGridViewItem(rect))
	}
	return gv
}

func createEmpty() fyne.CanvasObject {
	return empty.LoadingEmptyView("Loading...")
}

func createTable() fyne.CanvasObject {
	tv := table.NewTableView(table.TableViewStyleInsetGrouped)
	section := table.NewTableSection("Settings")
	section.Cells = []*table.TableViewCell{
		table.NewTableViewCellWithTextAndDetail("Profile", "View"),
		table.NewTableViewCellWithTextAndDetail("Notifications", "On"),
	}
	tv.Sections = []*table.TableSection{section}
	return tv
}

// ============ Interaction Helpers ============

func simulateTap(tappable fyne.Tappable) {
	if tappable != nil {
		tappable.Tapped(&fyne.PointEvent{})
	}
}

func simulateTyping(tf *textfield.TextField, text string) {
	if tf == nil {
		return
	}
	for _, ch := range text {
		tf.SetText(tf.Text + string(ch))
		time.Sleep(50 * time.Millisecond)
	}
}

func simulateTypingTextView(tv *textview.TextView, text string) {
	if tv == nil {
		return
	}
	for _, ch := range text {
		tv.SetText(tv.Text + string(ch))
		time.Sleep(50 * time.Millisecond)
	}
}

// ============ Card Helper ============

func createCard(title, description string, content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(qmuiWhite)
	bg.CornerRadius = 10
	bg.StrokeWidth = 1
	bg.StrokeColor = color.RGBA{R: 220, G: 220, B: 220, A: 255}

	titleLabel := canvas.NewText(title, qmuiCyan)
	titleLabel.TextSize = 14
	titleLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	descLabel := canvas.NewText(description, qmuiGrayText)
	descLabel.TextSize = 11

	sep := canvas.NewRectangle(color.RGBA{R: 230, G: 230, B: 230, A: 255})
	sep.SetMinSize(fyne.NewSize(0, 1))

	cardContent := container.NewVBox(
		container.NewPadded(container.NewVBox(titleLabel, descLabel)),
		sep,
		container.NewPadded(content),
	)

	return container.NewPadded(container.NewStack(bg, cardContent))
}
