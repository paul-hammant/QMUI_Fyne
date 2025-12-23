// Demo application showcasing QMUI Go components
// Shows actual interactive widgets with their real Go package names
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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/components/alert"
	"github.com/user/qmui-go/pkg/components/badge"
	"github.com/user/qmui-go/pkg/components/button"
	"github.com/user/qmui-go/pkg/components/checkbox"
	"github.com/user/qmui-go/pkg/components/dialog"
	"github.com/user/qmui-go/pkg/components/empty"
	"github.com/user/qmui-go/pkg/components/floatlayout"
	"github.com/user/qmui-go/pkg/components/grid"
	"github.com/user/qmui-go/pkg/components/modal"
	"github.com/user/qmui-go/pkg/components/qmuiswitch"
	"github.com/user/qmui-go/pkg/components/tips"
	"github.com/user/qmui-go/pkg/components/label"
	"github.com/user/qmui-go/pkg/components/marquee"
	"github.com/user/qmui-go/pkg/components/popup"
	"github.com/user/qmui-go/pkg/components/progress"
	"github.com/user/qmui-go/pkg/components/segmented"
	"github.com/user/qmui-go/pkg/components/table"
	"github.com/user/qmui-go/pkg/components/textfield"
	"github.com/user/qmui-go/pkg/components/textview"
	"github.com/user/qmui-go/pkg/components/toast"
	"github.com/user/qmui-go/pkg/console"
	"github.com/user/qmui-go/pkg/core"
	"github.com/user/qmui-go/pkg/log"
	"github.com/user/qmui-go/pkg/theme"
)

// QMUI Colors - matching iOS exactly
var (
	qmuiBackground = color.RGBA{R: 246, G: 247, B: 249, A: 255}
	qmuiWhite      = color.White
	qmuiGrayText   = color.RGBA{R: 134, G: 144, B: 156, A: 255}
	qmuiDarkText   = color.RGBA{R: 51, G: 51, B: 51, A: 255}
)

// primaryColor returns the current theme's primary color
func primaryColor() color.Color {
	return core.SharedConfiguration().BlueColor
}

var (
	mainWindow fyne.Window
	mainApp    fyne.App
)

func main() {
	mainApp = app.New()
	mainWindow = mainApp.NewWindow("QMUI Go Demo")
	mainWindow.Resize(fyne.NewSize(420, 750))

	// Make this the master window - closing it quits the app
	mainWindow.SetMaster()

	// Handle OS signals (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		fyne.Do(func() {
			mainApp.Quit()
		})
	}()

	// Ensure clean exit when window is closed
	mainWindow.SetCloseIntercept(func() {
		fmt.Println("Window closed, exiting...")
		mainApp.Quit()
	})

	// Initialize configuration
	_ = core.SharedConfiguration()

	// Register theme change listener ONCE to rebuild UI on theme change
	theme.SharedThemeManager().AddThemeChangeListener(func(newTheme *theme.Theme) {
		fmt.Printf("Theme changed to: %s\n", newTheme.Name)
		// Rebuild entire UI to pick up new theme colors from configuration
		fyne.Do(func() {
			mainWindow.SetContent(createMainUI())
		})
	})

	// Show splash screen first
	showSplashScreen()
}

func showSplashScreen() {
	bg := canvas.NewRectangle(primaryColor())

	// QMUI Logo
	logo := createQMUILogo()

	qmuiText := canvas.NewText("QMUI Go", qmuiWhite)
	qmuiText.TextSize = 24
	qmuiText.Alignment = fyne.TextAlignCenter

	subtitle := canvas.NewText("Fyne UI Components", color.RGBA{R: 255, G: 255, B: 255, A: 200})
	subtitle.TextSize = 14
	subtitle.Alignment = fyne.TextAlignCenter

	footer := canvas.NewText("Ported from Tencent QMUI iOS", color.RGBA{R: 255, G: 255, B: 255, A: 150})
	footer.TextSize = 10
	footer.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(logo),
		container.NewCenter(qmuiText),
		container.NewCenter(subtitle),
		layout.NewSpacer(),
		layout.NewSpacer(),
		container.NewCenter(footer),
		widget.NewLabel(""),
	)

	mainWindow.SetContent(container.NewStack(bg, content))
	mainWindow.Show()

	go func() {
		time.Sleep(2 * time.Second)
		fyne.Do(func() {
			mainWindow.SetContent(createMainUI())
		})
	}()

	mainApp.Run()
}

func createQMUILogo() fyne.CanvasObject {
	size := float32(80)
	outerCircle := canvas.NewCircle(color.Transparent)
	outerCircle.StrokeColor = qmuiWhite
	outerCircle.StrokeWidth = 4
	outerCircle.Resize(fyne.NewSize(size, size))

	tail := canvas.NewRectangle(qmuiWhite)
	tail.Resize(fyne.NewSize(size*0.35, 4))

	logoContainer := container.NewWithoutLayout(outerCircle, tail)
	outerCircle.Move(fyne.NewPos(0, 0))
	tail.Move(fyne.NewPos(size*0.55, size*0.7))

	wrapper := container.NewCenter(logoContainer)
	wrapper.Resize(fyne.NewSize(size+20, size+20))
	return wrapper
}

func createMainUI() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Components", createComponentsTab()),
		container.NewTabItem("Buttons", createButtonsTab()),
		container.NewTabItem("Inputs", createInputsTab()),
		container.NewTabItem("Progress", createProgressTab()),
		container.NewTabItem("Dialogs", createDialogsTab()),
		container.NewTabItem("Popups & Alerts", createPopupsTab()),
		container.NewTabItem("Controls", createControlsTab()),
		container.NewTabItem("Lab", createLabTab()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)
	tabs.SelectTabIndex(6)

	navBar := createNavBar("QMUI Go Components")
	bg := canvas.NewRectangle(qmuiBackground)

	return container.NewBorder(navBar, nil, nil, nil,
		container.NewStack(bg, tabs))
}

func createNavBar(title string) fyne.CanvasObject {
	tm := theme.SharedThemeManager()
	currentTheme := tm.CurrentTheme()

	navBarBg = canvas.NewRectangle(currentTheme.NavBarBackgroundColor)
	navBarBg.SetMinSize(fyne.NewSize(0, 44))

	navBarTitle = canvas.NewText(title, currentTheme.NavBarTitleColor)
	navBarTitle.TextSize = 17
	navBarTitle.TextStyle = fyne.TextStyle{Bold: true}
	navBarTitle.Alignment = fyne.TextAlignCenter

	return container.NewStack(navBarBg, container.NewCenter(navBarTitle))
}

// ============ Components Tab - Live Interactive Widgets ============

func createComponentsTab() fyne.CanvasObject {
	// Each card shows a LIVE widget, not just an icon
	cards := []fyne.CanvasObject{
		// marquee.MarqueeLabel - Live scrolling text
		createLiveComponentCard(
			"marquee.MarqueeLabel",
			"Scrolling text label",
			createLiveMarquee(),
		),

		// badge.BadgeLabel - Live badges
		createLiveComponentCard(
			"badge.BadgeLabel",
			"Notification badges",
			createLiveBadges(),
		),

		// label.Label - Enhanced label
		createLiveComponentCard(
			"label.Label",
			"Label with edge insets",
			createLiveLabel(),
		),

		// progress.PieProgressView - Live pie chart
		createLiveComponentCard(
			"progress.PieProgressView",
			"Pie chart progress",
			createLivePieProgress(),
		),

		// progress.CircularProgressView - Live circular
		createLiveComponentCard(
			"progress.CircularProgressView",
			"Circular progress ring",
			createLiveCircularProgress(),
		),

		// progress.LinearProgressView - Live linear bar
		createLiveComponentCard(
			"progress.LinearProgressView",
			"Linear progress bar",
			createLiveLinearProgress(),
		),

		// floatlayout.FloatLayoutView - Tag cloud
		createLiveComponentCard(
			"floatlayout.FloatLayoutView",
			"Tag cloud / flow layout",
			createLiveFloatLayout(),
		),

		// grid.GridView - Grid layout
		createLiveComponentCard(
			"grid.GridView",
			"Grid arrangement",
			createLiveGrid(),
		),

		// empty.EmptyView - Empty states
		createLiveComponentCard(
			"empty.EmptyView",
			"Loading / error states",
			createLiveEmpty(),
		),

		// table.TableView - Grouped lists
		createLiveComponentCard(
			"table.TableView",
			"iOS-style grouped lists",
			createLiveTable(),
		),
	}

	return container.NewScroll(container.NewVBox(cards...))
}

// ============ Live Component Creators ============

func createLiveMarquee() fyne.CanvasObject {
	// Blue background marquee (like iOS demo)
	m1 := marquee.NewMarqueeLabel("This text scrolls continuously in a blue container - just like iOS QMUIMarqueeLabel!")
	m1.Speed = 40
	m1.TextColor = color.White
	m1.StartAnimation()

	blueBg := canvas.NewRectangle(color.RGBA{R: 66, G: 133, B: 244, A: 255})
	blueBg.SetMinSize(fyne.NewSize(200, 50))
	blueMarquee := container.NewStack(blueBg, container.NewPadded(m1))

	// Pink background marquee
	m2 := marquee.NewMarqueeLabel("Another scrolling label with different color scheme!")
	m2.Speed = 30
	m2.TextColor = color.White
	m2.StartAnimation()

	pinkBg := canvas.NewRectangle(color.RGBA{R: 233, G: 30, B: 99, A: 255})
	pinkBg.SetMinSize(fyne.NewSize(200, 50))
	pinkMarquee := container.NewStack(pinkBg, container.NewPadded(m2))

	// Static text (short text doesn't scroll)
	staticLabel := widget.NewLabel("Short text (no scroll)")

	return container.NewVBox(
		staticLabel,
		blueMarquee,
		pinkMarquee,
	)
}

func createLiveBadges() fyne.CanvasObject {
	// Standalone badges
	b1 := badge.NewBadgeLabel("99+")
	b2 := badge.NewBadgeLabel("NEW")
	b3 := badge.NewBadgeLabel("5")
	indicator := badge.NewUpdatesIndicator()
	indicator.HasUpdates = true

	standaloneBadges := container.NewHBox(b1, b2, b3, indicator)

	// Badges in all four corners
	iconSize := fyne.NewSize(44, 44)

	// Top-right (default)
	icon1 := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255})
	icon1.SetMinSize(iconSize)
	icon1.CornerRadius = 8
	badgedIcon1 := badge.NewBadgeView(icon1)
	badgedIcon1.SetBadgeValue("TR")
	// Default offset is top-right

	// Top-left
	icon2 := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255})
	icon2.SetMinSize(iconSize)
	icon2.CornerRadius = 8
	badgedIcon2 := badge.NewBadgeView(icon2)
	badgedIcon2.SetBadgeValue("TL")
	badgedIcon2.BadgeOffset = core.NewOffset(-35, 11) // Shift left

	// Bottom-right
	icon3 := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255})
	icon3.SetMinSize(iconSize)
	icon3.CornerRadius = 8
	badgedIcon3 := badge.NewBadgeView(icon3)
	badgedIcon3.SetBadgeValue("BR")
	badgedIcon3.BadgeOffset = core.NewOffset(-9, 35) // Shift down

	// Bottom-left
	icon4 := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255})
	icon4.SetMinSize(iconSize)
	icon4.CornerRadius = 8
	badgedIcon4 := badge.NewBadgeView(icon4)
	badgedIcon4.SetBadgeValue("BL")
	badgedIcon4.BadgeOffset = core.NewOffset(-35, 35) // Shift left and down

	cornerBadges := container.NewHBox(badgedIcon1, badgedIcon2, badgedIcon3, badgedIcon4)

	return container.NewVBox(
		widget.NewLabel("Standalone:"),
		standaloneBadges,
		widget.NewLabel("Corner positions (TR, TL, BR, BL):"),
		cornerBadges,
	)
}

func createLiveLabel() fyne.CanvasObject {
	lbl := label.NewLabel("Padded Label")
	lbl.ContentEdgeInsets = core.EdgeInsets{Top: 8, Left: 16, Bottom: 8, Right: 16}
	// Wrap in a light blue background to show the insets
	bg := canvas.NewRectangle(color.RGBA{R: 230, G: 245, B: 255, A: 255})
	bg.CornerRadius = 4
	return container.NewStack(bg, lbl)
}

func createLivePieProgress() fyne.CanvasObject {
	pie := progress.NewPieProgressView()
	pie.Progress = 0.65
	pie.TintColor = primaryColor()

	// Animate on tap
	tapBtn := widget.NewButton("Animate", func() {
		go func() {
			for i := 0; i <= 100; i++ {
				pie.SetProgress(float64(i) / 100)
				time.Sleep(20 * time.Millisecond)
			}
		}()
	})
	tapBtn.Importance = widget.LowImportance

	return container.NewHBox(pie, tapBtn)
}

func createLiveCircularProgress() fyne.CanvasObject {
	circ := progress.NewCircularProgressView()
	circ.Progress = 0.45
	circ.TintColor = primaryColor()
	circ.ShowsText = true
	return circ
}

func createLiveLinearProgress() fyne.CanvasObject {
	lin := progress.NewLinearProgressView()
	lin.Progress = 0.7
	lin.TintColor = primaryColor()
	return lin
}

func createLiveFloatLayout() fyne.CanvasObject {
	fl := floatlayout.NewFloatLayoutView()
	fl.ItemSpacing = 6
	fl.LineSpacing = 6

	tags := []string{"Go", "Fyne", "QMUI", "Cross-Platform"}
	for _, tag := range tags {
		fl.AddItem(floatlayout.NewTagView(tag))
	}
	return fl
}

func createLiveGrid() fyne.CanvasObject {
	gv := grid.NewGridView(4)
	gv.RowSpacing = 4
	gv.ColumnSpacing = 4

	colors := []color.Color{
		primaryColor(),
		color.RGBA{R: 255, G: 100, B: 100, A: 255},
		color.RGBA{R: 100, G: 200, B: 100, A: 255},
		color.RGBA{R: 200, G: 150, B: 255, A: 255},
	}
	for _, c := range colors {
		rect := canvas.NewRectangle(c)
		rect.CornerRadius = 4
		rect.SetMinSize(fyne.NewSize(30, 30))
		gv.AddItem(grid.NewGridViewItem(rect))
	}
	return gv
}

func createLiveEmpty() fyne.CanvasObject {
	// Show different empty view states like iOS demo
	loadingView := empty.LoadingEmptyView("Loading data...")

	emptyStateView := empty.NewEmptyViewWithTextAndDetail("No Results", "Try adjusting your search filters")

	errorView := empty.NewEmptyViewWithTextAndDetail("Connection Error", "Please check your network")

	return container.NewVBox(
		widget.NewLabel("Loading state:"),
		loadingView,
		widget.NewLabel("Empty state:"),
		emptyStateView,
		widget.NewLabel("Error state:"),
		errorView,
	)
}

func createLiveTable() fyne.CanvasObject {
	tv := table.NewTableView(table.TableViewStyleInsetGrouped)
	section := table.NewTableSection("Settings")
	section.Cells = []*table.TableViewCell{
		table.NewTableViewCellWithTextAndDetail("Profile", "View profile"),
		table.NewTableViewCellWithTextAndDetail("Notifications", "Enabled"),
	}
	tv.Sections = []*table.TableSection{section}
	return tv
}

// ============ Buttons Tab ============

func createInputsTab() fyne.CanvasObject {
	cards := []fyne.CanvasObject{
		createLiveComponentCard(
			"textfield.TextField",
			"Text input field",
			createLiveTextField(),
		),
		createLiveComponentCard(
			"textview.TextView",
			"Multiline text view",
			createLiveTextView(),
		),
		createLiveComponentCard(
			"segmented.SegmentedControl",
			"Segmented control",
			createLiveSegmentedControl(),
		),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createLiveTextField() fyne.CanvasObject {
	tf := textfield.NewTextFieldWithPlaceholder("Placeholder")
	tf.OnTextChanged = func(s string) {
		fmt.Println("TextField changed:", s)
	}
	return tf
}

func createLiveTextView() fyne.CanvasObject {
	tv := textview.NewTextViewWithPlaceholder("Placeholder")
	tv.OnTextChanged = func(s string) {
		fmt.Println("TextView changed:", s)
	}
	return tv
}

func createLiveSegmentedControl() fyne.CanvasObject {
	sc := segmented.NewSegmentedControl([]string{"First", "Second", "Third"}, func(i int) {
		fmt.Println("SegmentedControl selected:", i)
	})
	return sc
}

func createButtonsTab() fyne.CanvasObject {
	cards := []fyne.CanvasObject{
		createLiveComponentCard(
			"button.Button",
			"Standard button",
			createLiveButton(),
		),
		createLiveComponentCard(
			"button.FillButton",
			"Solid filled button",
			createLiveFillButton(),
		),
		createLiveComponentCard(
			"button.GhostButton",
			"Outlined button",
			createLiveGhostButton(),
		),
		createLiveComponentCard(
			"button.NavigationButton",
			"Navigation bar button",
			createLiveNavButton(),
		),
		createLiveComponentCard(
			"Colored Variants",
			"Multiple tint colors",
			createColoredButtons(),
		),
	}

	return container.NewScroll(container.NewVBox(cards...))
}

func createLiveButton() fyne.CanvasObject {
	btn := button.NewButton("Tap Me", func() {
		toast.ShowMessage(mainWindow, "button.Button tapped!")
	})
	return btn
}

func createLiveFillButton() fyne.CanvasObject {
	cfg := core.SharedConfiguration()
	btn := button.NewFillButton("Fill Button", cfg.BlueColor, func() {
		toast.ShowMessage(mainWindow, "button.FillButton tapped!")
	})
	return btn
}

func createLiveGhostButton() fyne.CanvasObject {
	cfg := core.SharedConfiguration()

	// iOS-style "OK" button (like in QMUIButton demo)
	okBtn := button.NewGhostButton("OK", cfg.BlueColor, func() {
		toast.ShowMessage(mainWindow, "OK tapped!")
	})

	// Cancel button
	cancelBtn := button.NewGhostButton("Cancel", cfg.RedColor, func() {
		toast.ShowMessage(mainWindow, "Cancel tapped!")
	})

	// Longer text button
	longBtn := button.NewGhostButton("Submit Form", cfg.BlueColor, func() {
		toast.ShowMessage(mainWindow, "Submit tapped!")
	})

	return container.NewHBox(okBtn, cancelBtn, longBtn)
}

func createLiveNavButton() fyne.CanvasObject {
	btn := button.NewNavigationButton("< Back", func() {
		toast.ShowMessage(mainWindow, "button.NavigationButton tapped!")
	})
	return btn
}

func createColoredButtons() fyne.CanvasObject {
	cfg := core.SharedConfiguration()
	red := button.NewFillButton("Red", cfg.RedColor, nil)
	green := button.NewFillButton("Green", cfg.GreenColor, nil)
	blue := button.NewFillButton("Blue", cfg.BlueColor, nil)
	return container.NewHBox(red, green, blue)
}

// ============ Progress Tab ============

func createProgressTab() fyne.CanvasObject {
	// Create animated progress indicators
	pie := progress.NewPieProgressView()
	pie.Progress = 0.0
	pie.TintColor = primaryColor()

	circ := progress.NewCircularProgressView()
	circ.Progress = 0.0
	circ.TintColor = color.RGBA{R: 52, G: 199, B: 89, A: 255}
	circ.ShowsText = true

	lin := progress.NewLinearProgressView()
	lin.Progress = 0.0
	lin.TintColor = color.RGBA{R: 255, G: 149, B: 0, A: 255}

	animateBtn := widget.NewButton("Animate All", func() {
		go func() {
			for i := 0; i <= 100; i++ {
				p := float64(i) / 100
				pie.SetProgress(p)
				circ.SetProgress(p)
				lin.SetProgress(p)
				time.Sleep(30 * time.Millisecond)
			}
		}()
	})

	resetBtn := widget.NewButton("Reset", func() {
		pie.SetProgress(0)
		circ.SetProgress(0)
		lin.SetProgress(0)
	})

	return container.NewScroll(container.NewVBox(
		createLiveComponentCard("progress.PieProgressView", "Pie chart style", pie),
		createLiveComponentCard("progress.CircularProgressView", "Ring with percentage", circ),
		createLiveComponentCard("progress.LinearProgressView", "Horizontal bar", lin),
		createSectionCard("Controls", container.NewHBox(animateBtn, resetBtn)),
	))
}

// ============ Popups & Alerts Tab ============

func createPopupsTab() fyne.CanvasObject {
	dialogBtn := widget.NewButton("Show dialog.DialogViewController", func() {
		dialog.ShowConfirmDialog(mainWindow, "Confirm", "Are you sure?", nil, nil)
	})

	modalBtn := widget.NewButton("Show modal.ModalPresentationViewController", func() {
		content := widget.NewLabel("This is a modal dialog")
		modal.NewModalPresentationViewControllerWithContent(content).Present(mainWindow)
	})

	tipsBtn := widget.NewButton("Show tips.Tips", func() {
		tips.ShowSuccess(mainWindow, "Success")
	})

	return container.NewScroll(container.NewVBox(
		createSectionCard("dialog.DialogViewController", dialogBtn),
		createSectionCard("modal.ModalPresentationViewController", modalBtn),
		createSectionCard("tips.Tips", tipsBtn),
	))
}

// ============ Controls Tab ============

func createControlsTab() fyne.CanvasObject {
	cards := []fyne.CanvasObject{
		createLiveComponentCard(
			"checkbox.Checkbox",
			"Checkbox component",
			createLiveCheckbox(),
		),
		createLiveComponentCard(
			"switch.Switch",
			"Switch component",
			createLiveSwitch(),
		),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createLiveSwitch() fyne.CanvasObject {
	s1 := qmuiswitch.NewSwitch(func(b bool) {
		fmt.Printf("Switch 1 changed: %v\n", b)
	})

	s2 := qmuiswitch.NewSwitch(func(b bool) {
		fmt.Printf("Switch 2 changed: %v\n", b)
	})
	s2.SetChecked(true)

	s3 := qmuiswitch.NewSwitch(nil)
	s3.Enabled = false

	s4 := qmuiswitch.NewSwitch(nil)
	s4.Enabled = false
	s4.SetChecked(true)

	return container.NewVBox(
		s1,
		s2,
		s3,
		s4,
	)
}

func createLiveCheckbox() fyne.CanvasObject {
	cb1 := checkbox.NewCheckboxWithLabel("Option 1", func(b bool) {
		fmt.Printf("Checkbox 1 changed: %v\n", b)
	})
	cb2 := checkbox.NewCheckboxWithLabel("Option 2", func(b bool) {
		fmt.Printf("Checkbox 2 changed: %v\n", b)
	})
	cb2.SetSelected(true)

	cb3 := checkbox.NewCheckboxWithLabel("Indeterminate", nil)
	cb3.SetIndeterminate(true)

	cb4 := checkbox.NewCheckboxWithLabel("Disabled", nil)
	cb4.Enabled = false

	cb5 := checkbox.NewCheckboxWithLabel("Disabled & Checked", nil)
	cb5.Enabled = false
	cb5.SetSelected(true)

	return container.NewVBox(
		cb1,
		cb2,
		cb3,
		cb4,
		cb5,
	)
}

// ============ Dialogs Tab ============

func createDialogsTab() fyne.CanvasObject {
	toastBtn := widget.NewButton("Show toast.ToastView", func() {
		toast.ShowMessage(mainWindow, "This is a toast message!")
	})

	alertBtn := widget.NewButton("Show alert.AlertController", func() {
		ac := alert.NewAlertController("Alert Demo", "This is alert.AlertController in action.", alert.ControllerStyleAlert)
		ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))
		ac.AddAction(alert.NewAction("OK", alert.ActionStyleDefault, func(_ *alert.AlertController, _ *alert.Action) {
			toast.ShowMessage(mainWindow, "OK pressed!")
		}))
		ac.ShowIn(mainWindow)
	})

	actionSheetBtn := widget.NewButton("Show Action Sheet", func() {
		ac := alert.NewAlertController("Choose Action", "", alert.ControllerStyleActionSheet)
		ac.AddAction(alert.NewAction("Share", alert.ActionStyleDefault, func(_ *alert.AlertController, _ *alert.Action) {
			toast.ShowMessage(mainWindow, "Share selected")
		}))
		ac.AddAction(alert.NewAction("Delete", alert.ActionStyleDestructive, func(_ *alert.AlertController, _ *alert.Action) {
			toast.ShowMessage(mainWindow, "Delete selected")
		}))
		ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))
		ac.ShowIn(mainWindow)
	})

	popupBtn := widget.NewButton("Show popup.PopupMenu", func() {
		items := []*popup.MenuItem{
			popup.NewMenuItem("Edit", func(_ *popup.MenuItem) {
				toast.ShowMessage(mainWindow, "Edit selected")
			}),
			popup.NewMenuItem("Copy", func(_ *popup.MenuItem) {
				toast.ShowMessage(mainWindow, "Copy selected")
			}),
			popup.NewMenuItem("Delete", func(_ *popup.MenuItem) {
				toast.ShowMessage(mainWindow, "Delete selected")
			}),
		}
		popup.ContextMenu(mainWindow, fyne.NewPos(200, 400), items)
	})

	return container.NewScroll(container.NewVBox(
		createSectionCard("toast.ToastView", toastBtn),
		createSectionCard("alert.AlertController", container.NewVBox(alertBtn, actionSheetBtn)),
		createSectionCard("popup.PopupMenu", popupBtn),
	))
}

// ============ Lab Tab ============

// Global references for hot-switching demo
var (
	navBarBg       *canvas.Rectangle
	navBarTitle    *canvas.Text
	themeLabel     *widget.Label
	sampleFillBtn  *button.FillButton
	sampleGhostBtn *button.GhostButton
)

func createLabTab() fyne.CanvasObject {
	tm := theme.SharedThemeManager()
	currentTheme := tm.CurrentTheme()

	themeLabel = widget.NewLabel(fmt.Sprintf("Current theme: %s", currentTheme.Name))

	// Sample widgets that will update with theme
	sampleFillBtn = button.NewFillButton("Fill Button", currentTheme.PrimaryColor, func() {
		toast.ShowMessage(mainWindow, "Theme-aware button tapped!")
	})
	sampleGhostBtn = button.NewGhostButton("Ghost Button", currentTheme.PrimaryColor, func() {
		toast.ShowMessage(mainWindow, "Theme-aware ghost button tapped!")
	})

	// Theme listener is registered once in main()

	// Theme buttons - Row 1: Light themes
	defaultBtn := widget.NewButton("Default", func() {
		tm.SetCurrentTheme(theme.ThemeIdentifierDefault)
	})
	grapefruitBtn := widget.NewButton("Grapefruit", func() {
		tm.SetCurrentTheme(theme.ThemeIdentifierGrapefruit)
	})
	grassBtn := widget.NewButton("Grass", func() {
		tm.SetCurrentTheme(theme.ThemeIdentifierGrass)
	})

	// Theme buttons - Row 2: More themes
	pinkBtn := widget.NewButton("Pink Rose", func() {
		tm.SetCurrentTheme(theme.ThemeIdentifierPinkRose)
	})
	grayBtn := widget.NewButton("Gray", func() {
		tm.SetCurrentTheme(theme.ThemeIdentifierGray)
	})
	darkBtn := widget.NewButton("Dark", func() {
		tm.SetCurrentTheme(theme.ThemeIdentifierDark)
	})

	consoleBtn := widget.NewButton("Toggle console.Console", func() {
		console.SharedConsole().Toggle(mainWindow)
	})

	logBtn := widget.NewButton("Add log.QMUILog Entry", func() {
		log.QMUILog("Test log at %v", time.Now())
		toast.ShowMessage(mainWindow, "Log entry added")
	})

	return container.NewScroll(container.NewVBox(
		createSectionCard("Hot-Switchable Themes", container.NewVBox(
			widget.NewLabel("Tap a theme to switch instantly:"),
			container.NewHBox(defaultBtn, grapefruitBtn, grassBtn),
			container.NewHBox(pinkBtn, grayBtn, darkBtn),
			themeLabel,
		)),
		createSectionCard("Theme-Aware Widgets", container.NewVBox(
			widget.NewLabel("These widgets update with theme:"),
			container.NewHBox(sampleFillBtn, sampleGhostBtn),
		)),
		createSectionCard("console.Console", consoleBtn),
		createSectionCard("log.QMUILog", logBtn),
	))
}

// ============ Card Helpers ============

func createLiveComponentCard(componentName, description string, liveWidget fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(qmuiWhite)
	bg.CornerRadius = 10
	bg.StrokeWidth = 1
	bg.StrokeColor = primaryColor()

	// Component name in cyan monospace
	nameLabel := canvas.NewText(componentName, primaryColor())
	nameLabel.TextSize = 14
	nameLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	// Description in gray
	descLabel := canvas.NewText(description, qmuiGrayText)
	descLabel.TextSize = 11

	// Separator
	sep := canvas.NewRectangle(color.RGBA{R: 230, G: 230, B: 230, A: 255})
	sep.SetMinSize(fyne.NewSize(0, 1))

	// Live widget area with padding
	widgetArea := container.NewPadded(liveWidget)

	cardContent := container.NewVBox(
		container.NewPadded(container.NewVBox(nameLabel, descLabel)),
		sep,
		widgetArea,
	)

	return container.NewPadded(container.NewStack(bg, cardContent))
}

func createSectionCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(qmuiWhite)
	bg.CornerRadius = 10
	bg.StrokeWidth = 0.5
	bg.StrokeColor = color.RGBA{R: 220, G: 220, B: 220, A: 255}

	titleLabel := canvas.NewText(title, qmuiDarkText)
	titleLabel.TextSize = 14
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	sep := canvas.NewRectangle(color.RGBA{R: 230, G: 230, B: 230, A: 255})
	sep.SetMinSize(fyne.NewSize(0, 1))

	cardContent := container.NewVBox(
		container.NewPadded(titleLabel),
		sep,
		container.NewPadded(content),
	)

	return container.NewPadded(container.NewStack(bg, cardContent))
}
