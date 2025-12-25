// Demo application showcasing QMUI Go components
// Organized to match iOS QMUI structure with all 11 themes
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
	"github.com/user/qmui-go/pkg/components/imagepreview"
	"github.com/user/qmui-go/pkg/components/label"
	"github.com/user/qmui-go/pkg/components/marquee"
	"github.com/user/qmui-go/pkg/components/modal"
	"github.com/user/qmui-go/pkg/components/moreop"
	"github.com/user/qmui-go/pkg/components/navigation"
	"github.com/user/qmui-go/pkg/components/popup"
	"github.com/user/qmui-go/pkg/components/progress"
	"github.com/user/qmui-go/pkg/components/qmuiswitch"
	"github.com/user/qmui-go/pkg/components/search"
	"github.com/user/qmui-go/pkg/components/segmented"
	"github.com/user/qmui-go/pkg/components/table"
	"github.com/user/qmui-go/pkg/components/textfield"
	"github.com/user/qmui-go/pkg/components/textview"
	"github.com/user/qmui-go/pkg/components/tile"
	"github.com/user/qmui-go/pkg/components/tips"
	"github.com/user/qmui-go/pkg/components/toast"
	"github.com/user/qmui-go/pkg/console"
	"github.com/user/qmui-go/pkg/core"
	"github.com/user/qmui-go/pkg/log"
	"github.com/user/qmui-go/pkg/theme"
)

// Colors
var (
	qmuiBackground = color.RGBA{R: 246, G: 247, B: 249, A: 255}
	qmuiWhite      = color.White
	qmuiGrayText   = color.RGBA{R: 134, G: 144, B: 156, A: 255}
	qmuiDarkText   = color.RGBA{R: 51, G: 51, B: 51, A: 255}
)

func primaryColor() color.Color {
	return theme.SharedThemeManager().CurrentTheme().PrimaryColor
}

var (
	mainWindow fyne.Window
	mainApp    fyne.App
)

func main() {
	mainApp = app.New()
	mainWindow = mainApp.NewWindow("QMUI Go Demo")
	mainWindow.Resize(fyne.NewSize(420, 750))
	mainWindow.SetMaster()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		fyne.Do(func() {
			mainApp.Quit()
		})
	}()

	mainWindow.SetCloseIntercept(func() {
		fmt.Println("Window closed, exiting...")
		mainApp.Quit()
	})

	_ = core.SharedConfiguration()

	// Theme change listener - rebuild UI on theme change
	theme.SharedThemeManager().AddThemeChangeListener(func(newTheme *theme.Theme) {
		fmt.Printf("Theme changed to: %s\n", newTheme.Name)
		fyne.Do(func() {
			mainWindow.SetContent(createMainUI())
		})
	})

	showSplashScreen()
}

func showSplashScreen() {
	bg := canvas.NewRectangle(primaryColor())

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
		container.NewTabItem("🎨 Themes", createThemesTab()),
		container.NewTabItem("📦 Components", createComponentsTab()),
		container.NewTabItem("🎛 Controls", createControlsTab()),
		container.NewTabItem("⏳ Progress", createProgressTab()),
		container.NewTabItem("📐 Layout", createLayoutTab()),
		container.NewTabItem("💬 Dialogs", createDialogsTab()),
		container.NewTabItem("🧭 Navigation", createNavigationTab()),
		container.NewTabItem("🛠 Debug", createDebugTab()),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	navBar := createNavBar("QMUI Go Components")
	bg := canvas.NewRectangle(qmuiBackground)

	return container.NewBorder(navBar, nil, nil, nil,
		container.NewStack(bg, tabs))
}

func createNavBar(title string) fyne.CanvasObject {
	currentTheme := theme.SharedThemeManager().CurrentTheme()

	bg := canvas.NewRectangle(currentTheme.NavBarBackgroundColor)
	bg.SetMinSize(fyne.NewSize(0, 44))

	titleText := canvas.NewText(title, currentTheme.NavBarTitleColor)
	titleText.TextSize = 17
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter

	return container.NewStack(bg, container.NewCenter(titleText))
}

// ═══════════════════════════════════════════════════════════════
// THEMES TAB - All 11 QMUI iOS themes
// ═══════════════════════════════════════════════════════════════

func createThemesTab() fyne.CanvasObject {
	tm := theme.SharedThemeManager()
	currentTheme := tm.CurrentTheme()

	themeLabel := widget.NewLabel(fmt.Sprintf("Current: %s", currentTheme.Name))
	themeLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create color swatch showing current theme
	colorSwatch := canvas.NewRectangle(currentTheme.PrimaryColor)
	colorSwatch.SetMinSize(fyne.NewSize(60, 60))
	colorSwatch.CornerRadius = 8

	// Theme grid - all 11 themes
	themeButtons := container.NewGridWithColumns(3)

	allThemes := tm.AllThemes()
	for _, t := range allThemes {
		th := t // capture
		btn := widget.NewButton(th.Name, func() {
			tm.SetCurrentTheme(th.Identifier)
		})
		if th.IsDarkMode {
			btn.Importance = widget.HighImportance
		}
		themeButtons.Add(btn)
	}

	// Cycle button
	cycleBtn := widget.NewButton("⟳ Cycle Theme", func() {
		newTheme := tm.CycleTheme()
		toast.ShowMessage(mainWindow, fmt.Sprintf("Switched to %s", newTheme.Name))
	})
	cycleBtn.Importance = widget.HighImportance

	// Sample themed widgets
	sampleFill := button.NewFillButton("Fill Button", currentTheme.PrimaryColor, func() {
		toast.ShowMessage(mainWindow, "Themed button tapped!")
	})
	sampleGhost := button.NewGhostButton("Ghost Button", currentTheme.PrimaryColor, nil)

	sampleProgress := progress.NewPieProgress()
	sampleProgress.Progress = 0.65
	sampleProgress.TintColor = currentTheme.PrimaryColor

	sampleCheckbox := checkbox.NewCheckboxWithLabel("Themed checkbox", nil)
	sampleCheckbox.TintColor = currentTheme.PrimaryColor
	sampleCheckbox.SetSelected(true)

	return container.NewScroll(container.NewVBox(
		createSectionCard("Current Theme", container.NewHBox(colorSwatch, themeLabel)),
		createSectionCard("All 11 QMUI Themes", container.NewVBox(
			themeButtons,
			cycleBtn,
		)),
		createSectionCard("Theme-Aware Widgets", container.NewVBox(
			container.NewHBox(sampleFill, sampleGhost),
			container.NewHBox(sampleProgress, sampleCheckbox),
		)),
	))
}

// ═══════════════════════════════════════════════════════════════
// COMPONENTS TAB - Badge, Label, Marquee, EmptyView
// ═══════════════════════════════════════════════════════════════

func createComponentsTab() fyne.CanvasObject {
	return container.NewScroll(container.NewVBox(
		createComponentCard("badge.BadgeLabel", "Notification badges", createBadgeDemo()),
		createComponentCard("badge.UpdatesIndicator", "Small dot indicator", createIndicatorDemo()),
		createComponentCard("label.Label", "Label with edge insets", createLabelDemo()),
		createComponentCard("marquee.MarqueeLabel", "Scrolling text", createMarqueeDemo()),
		createComponentCard("empty.EmptyView", "Loading/error states", createEmptyDemo()),
	))
}

func createBadgeDemo() fyne.CanvasObject {
	b1 := badge.NewBadge("5")
	b2 := badge.NewBadge("99+")
	b3 := badge.NewBadge("NEW")

	// Badge on icon
	icon := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255})
	icon.SetMinSize(fyne.NewSize(44, 44))
	icon.CornerRadius = 8
	badgedIcon := badge.NewBadgeView(icon)
	badgedIcon.SetBadgeValue("3")

	return container.NewVBox(
		widget.NewLabel("Standalone badges:"),
		container.NewHBox(b1, b2, b3),
		widget.NewLabel("Badge on icon:"),
		badgedIcon,
	)
}

func createIndicatorDemo() fyne.CanvasObject {
	ind := badge.NewUpdatesIndicator()
	ind.HasUpdates = true
	return container.NewHBox(widget.NewLabel("Has updates:"), ind)
}

func createLabelDemo() fyne.CanvasObject {
	lbl := label.NewLabel("Padded Label")
	lbl.ContentEdgeInsets = core.EdgeInsets{Top: 8, Left: 16, Bottom: 8, Right: 16}
	bg := canvas.NewRectangle(color.RGBA{R: 230, G: 245, B: 255, A: 255})
	bg.CornerRadius = 4
	return container.NewStack(bg, lbl)
}

func createMarqueeDemo() fyne.CanvasObject {
	m := marquee.NewMarquee("This text scrolls continuously - perfect for news tickers!")
	m.Speed = 40
	m.TextColor = qmuiWhite
	m.StartAnimation()

	bg := canvas.NewRectangle(primaryColor())
	bg.SetMinSize(fyne.NewSize(250, 40))
	bg.CornerRadius = 4

	return container.NewStack(bg, container.NewPadded(m))
}

func createEmptyDemo() fyne.CanvasObject {
	loadingView := empty.LoadingEmptyState("Loading...")
	emptyView := empty.NewEmptyStateWithTextAndDetail("No Results", "Try a different search")

	return container.NewVBox(
		widget.NewLabel("Loading state:"),
		loadingView,
		widget.NewLabel("Empty state:"),
		emptyView,
	)
}

// ═══════════════════════════════════════════════════════════════
// CONTROLS TAB - Button, Checkbox, Switch, Segmented, TextField
// ═══════════════════════════════════════════════════════════════

func createControlsTab() fyne.CanvasObject {
	return container.NewScroll(container.NewVBox(
		createComponentCard("button.Button", "Standard button", createButtonDemo()),
		createComponentCard("button.FillButton", "Solid filled button", createFillButtonDemo()),
		createComponentCard("button.GhostButton", "Outlined button", createGhostButtonDemo()),
		createComponentCard("checkbox.Checkbox", "Circular checkbox", createCheckboxDemo()),
		createComponentCard("qmuiswitch.Switch", "iOS-style toggle", createSwitchDemo()),
		createComponentCard("segmented.SegmentedControl", "Segmented picker", createSegmentedDemo()),
		createComponentCard("textfield.TextField", "Text input", createTextFieldDemo()),
		createComponentCard("textview.TextView", "Multiline input", createTextViewDemo()),
		createComponentCard("search.SearchBar", "Search input", createSearchDemo()),
	))
}

func createButtonDemo() fyne.CanvasObject {
	btn := button.NewButton("Tap Me", func() {
		toast.ShowMessage(mainWindow, "Button tapped!")
	})
	return btn
}

func createFillButtonDemo() fyne.CanvasObject {
	cfg := core.SharedConfiguration()
	return container.NewHBox(
		button.NewFillButton("Blue", cfg.BlueColor, nil),
		button.NewFillButton("Red", cfg.RedColor, nil),
		button.NewFillButton("Green", cfg.GreenColor, nil),
	)
}

func createGhostButtonDemo() fyne.CanvasObject {
	cfg := core.SharedConfiguration()
	return container.NewHBox(
		button.NewGhostButton("OK", cfg.BlueColor, nil),
		button.NewGhostButton("Cancel", cfg.RedColor, nil),
	)
}

func createCheckboxDemo() fyne.CanvasObject {
	cb1 := checkbox.NewCheckboxWithLabel("Option 1", nil)
	cb2 := checkbox.NewCheckboxWithLabel("Option 2 (checked)", nil)
	cb2.SetSelected(true)
	cb3 := checkbox.NewCheckboxWithLabel("Indeterminate", nil)
	cb3.SetIndeterminate(true)
	cb4 := checkbox.NewCheckboxWithLabel("Disabled", nil)
	cb4.Enabled = false

	return container.NewVBox(cb1, cb2, cb3, cb4)
}

func createSwitchDemo() fyne.CanvasObject {
	s1 := qmuiswitch.NewSwitch(func(on bool) {
		fmt.Printf("Switch: %v\n", on)
	})
	s2 := qmuiswitch.NewSwitch(nil)
	s2.SetChecked(true)
	s3 := qmuiswitch.NewSwitch(nil)
	s3.Enabled = false

	return container.NewHBox(s1, s2, s3)
}

func createSegmentedDemo() fyne.CanvasObject {
	sc := segmented.NewSegmentedControl([]string{"Day", "Week", "Month"}, func(i int) {
		toast.ShowMessage(mainWindow, fmt.Sprintf("Selected: %d", i))
	})
	return sc
}

func createTextFieldDemo() fyne.CanvasObject {
	tf := textfield.NewTextFieldWithPlaceholder("Enter text...")
	return tf
}

func createTextViewDemo() fyne.CanvasObject {
	tv := textview.NewTextViewWithPlaceholder("Multi-line text...")
	return tv
}

func createSearchDemo() fyne.CanvasObject {
	sb := search.NewSearchBar()
	sb.Placeholder = "Search..."
	return sb
}

// ═══════════════════════════════════════════════════════════════
// PROGRESS TAB - Pie, Circular, Linear
// ═══════════════════════════════════════════════════════════════

func createProgressTab() fyne.CanvasObject {
	pie := progress.NewPieProgress()
	pie.Progress = 0.0
	pie.TintColor = primaryColor()

	circ := progress.NewRingProgress()
	circ.Progress = 0.0
	circ.TintColor = color.RGBA{R: 52, G: 199, B: 89, A: 255}
	circ.ShowsText = true

	lin := progress.NewProgressBar()
	lin.Progress = 0.0
	lin.TintColor = color.RGBA{R: 255, G: 149, B: 0, A: 255}

	animateBtn := widget.NewButton("▶ Animate All", func() {
		go func() {
			for i := 0; i <= 100; i++ {
				p := float64(i) / 100
				fyne.Do(func() {
					pie.SetProgress(p)
					circ.SetProgress(p)
					lin.SetProgress(p)
				})
				time.Sleep(25 * time.Millisecond)
			}
		}()
	})
	animateBtn.Importance = widget.HighImportance

	resetBtn := widget.NewButton("↺ Reset", func() {
		pie.SetProgress(0)
		circ.SetProgress(0)
		lin.SetProgress(0)
	})

	return container.NewScroll(container.NewVBox(
		createComponentCard("progress.PieProgressView", "Pie chart style", pie),
		createComponentCard("progress.CircularProgressView", "Ring with percentage", circ),
		createComponentCard("progress.LinearProgressView", "Horizontal bar", lin),
		createSectionCard("Controls", container.NewHBox(animateBtn, resetBtn)),
	))
}

// ═══════════════════════════════════════════════════════════════
// LAYOUT TAB - Grid, FloatLayout, Tile
// ═══════════════════════════════════════════════════════════════

func createLayoutTab() fyne.CanvasObject {
	return container.NewScroll(container.NewVBox(
		createComponentCard("grid.GridView", "Grid arrangement", createGridDemo()),
		createComponentCard("floatlayout.FloatLayoutView", "Tag cloud layout", createFloatLayoutDemo()),
		createComponentCard("tile.TileView", "Image tiles", createTileDemo()),
		createComponentCard("imagepreview.ImagePreviewView", "Image gallery", createImagePreviewDemo()),
	))
}

func createGridDemo() fyne.CanvasObject {
	gv := grid.NewGrid(4)
	gv.RowSpacing = 4
	gv.ColumnSpacing = 4

	colors := []color.Color{
		primaryColor(),
		color.RGBA{R: 255, G: 100, B: 100, A: 255},
		color.RGBA{R: 100, G: 200, B: 100, A: 255},
		color.RGBA{R: 200, G: 150, B: 255, A: 255},
		color.RGBA{R: 255, G: 200, B: 100, A: 255},
		color.RGBA{R: 100, G: 200, B: 255, A: 255},
		color.RGBA{R: 200, G: 100, B: 200, A: 255},
		color.RGBA{R: 150, G: 200, B: 150, A: 255},
	}

	for _, c := range colors {
		rect := canvas.NewRectangle(c)
		rect.CornerRadius = 4
		rect.SetMinSize(fyne.NewSize(40, 40))
		gv.AddItem(grid.NewGridItem(rect))
	}

	return gv
}

func createFloatLayoutDemo() fyne.CanvasObject {
	tc := floatlayout.NewTagCloud()
	tc.SetTags([]string{"Go", "Fyne", "QMUI", "iOS", "Cross-Platform", "Widgets", "Components"})
	return tc
}

func createTileDemo() fyne.CanvasObject {
	// Create component tiles like iOS QMUI
	tiles := []fyne.CanvasObject{}
	names := []string{"Badge", "Button", "Switch", "Progress", "Alert", "Toast"}
	colors := []color.Color{
		color.RGBA{R: 200, G: 150, B: 150, A: 255},
		color.RGBA{R: 150, G: 200, B: 150, A: 255},
		color.RGBA{R: 150, G: 150, B: 200, A: 255},
	}

	for i, name := range names {
		rect := canvas.NewRectangle(colors[i%3])
		rect.CornerRadius = 4
		rect.SetMinSize(fyne.NewSize(40, 40))
		t := tile.NewComponentTile(name, rect)
		tiles = append(tiles, t)
	}

	return container.NewGridWithColumns(3, tiles...)
}

func createImagePreviewDemo() fyne.CanvasObject {
	// Create a simple image preview placeholder
	// (ImagePreviewView expects fyne.Resource images)
	preview := imagepreview.NewImagePreview()

	// Show placeholder with instructions
	placeholder := container.NewVBox(
		widget.NewLabel("ImagePreviewView displays"),
		widget.NewLabel("zoomable image galleries."),
		widget.NewLabel("Requires fyne.Resource images."),
	)

	return container.NewVBox(preview, placeholder)
}

// ═══════════════════════════════════════════════════════════════
// DIALOGS TAB - Toast, Tips, Alert, Dialog, Popup, Modal
// ═══════════════════════════════════════════════════════════════

func createDialogsTab() fyne.CanvasObject {
	toastBtn := widget.NewButton("Show Toast", func() {
		toast.ShowMessage(mainWindow, "This is a toast message!")
	})

	successBtn := widget.NewButton("Success", func() {
		tips.ShowSuccess(mainWindow, "Success!")
	})
	errorBtn := widget.NewButton("Error", func() {
		tips.ShowError(mainWindow, "Error!")
	})
	loadingBtn := widget.NewButton("Loading", func() {
		t := tips.NewHUD(mainWindow)
		t.ShowLoading("Loading...")
		go func() {
			time.Sleep(2 * time.Second)
			fyne.Do(func() {
				t.HideCurrent()
			})
		}()
	})

	alertBtn := widget.NewButton("Show Alert", func() {
		ac := alert.NewAlert("Alert", "This is an alert dialog", alert.ControllerStyleAlert)
		ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))
		ac.AddAction(alert.NewAction("OK", alert.ActionStyleDefault, func(_ *alert.Alert, _ *alert.Action) {
			toast.ShowMessage(mainWindow, "OK pressed!")
		}))
		ac.ShowIn(mainWindow)
	})

	actionSheetBtn := widget.NewButton("Action Sheet", func() {
		ac := alert.NewAlert("Choose", "", alert.ControllerStyleActionSheet)
		ac.AddAction(alert.NewAction("Share", alert.ActionStyleDefault, nil))
		ac.AddAction(alert.NewAction("Delete", alert.ActionStyleDestructive, nil))
		ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))
		ac.ShowIn(mainWindow)
	})

	dialogBtn := widget.NewButton("Show Dialog", func() {
		dialog.ShowConfirmDialog(mainWindow, "Confirm", "Are you sure?", nil, nil)
	})

	popupBtn := widget.NewButton("Show Popup", func() {
		items := []*popup.MenuItem{
			popup.NewMenuItem("Edit", func(_ *popup.MenuItem) {
				toast.ShowMessage(mainWindow, "Edit")
			}),
			popup.NewMenuItem("Copy", func(_ *popup.MenuItem) {
				toast.ShowMessage(mainWindow, "Copy")
			}),
			popup.NewMenuItem("Delete", func(_ *popup.MenuItem) {
				toast.ShowMessage(mainWindow, "Delete")
			}),
		}
		popup.ContextMenu(mainWindow, fyne.NewPos(200, 400), items)
	})

	moreOpBtn := widget.NewButton("More Operations", func() {
		items := []*moreop.Item{
			moreop.NewItem("share", "Share", nil, nil),
			moreop.NewItem("copy", "Copy", nil, nil),
			moreop.NewItem("save", "Save", nil, nil),
		}
		moreop.ShowOperationSheet(mainWindow, items, "Cancel")
	})

	modalBtn := widget.NewButton("Show Modal", func() {
		content := widget.NewLabel("This is modal content")
		modal.NewModalWithContent(container.NewPadded(content)).Present(mainWindow)
	})

	return container.NewScroll(container.NewVBox(
		createSectionCard("toast.ToastView", toastBtn),
		createSectionCard("tips.Tips", container.NewHBox(successBtn, errorBtn, loadingBtn)),
		createSectionCard("alert.Alert", container.NewVBox(alertBtn, actionSheetBtn)),
		createSectionCard("dialog.DialogViewController", dialogBtn),
		createSectionCard("popup.PopupMenu", popupBtn),
		createSectionCard("moreop.MoreOperationController", moreOpBtn),
		createSectionCard("modal.ModalPresentationViewController", modalBtn),
	))
}

// ═══════════════════════════════════════════════════════════════
// NAVIGATION TAB - NavigationBar, TabBar, Table
// ═══════════════════════════════════════════════════════════════

func createNavigationTab() fyne.CanvasObject {
	return container.NewScroll(container.NewVBox(
		createComponentCard("navigation.NavigationBar", "App navigation bar", createNavBarDemo()),
		createComponentCard("navigation.TabBar", "Bottom tab bar", createTabBarDemo()),
		createComponentCard("table.TableView", "iOS-style grouped list", createTableDemo()),
	))
}

func createNavBarDemo() fyne.CanvasObject {
	navBar := navigation.NewNavigationBar()
	navBar.SetTitleView(navigation.NewTitleViewWithTitle("Title"))
	navBar.TintColor = primaryColor()
	return navBar
}

func createTabBarDemo() fyne.CanvasObject {
	items := []*navigation.TabBarItem{
		navigation.NewTabBarItem("Home", nil),
		navigation.NewTabBarItem("Search", nil),
		navigation.NewTabBarItem("Profile", nil),
	}
	tabBar := navigation.NewTabBar(items)
	tabBar.SetSelectedIndex(0)
	return tabBar
}

func createTableDemo() fyne.CanvasObject {
	tv := table.NewTable(table.TableStyleInsetGrouped)

	section1 := table.NewTableSection("Account")
	section1.Cells = []*table.TableCell{
		table.NewTableCellWithTextAndDetail("Profile", "View"),
		table.NewTableCellWithTextAndDetail("Settings", "Configure"),
	}

	section2 := table.NewTableSection("Preferences")
	section2.Cells = []*table.TableCell{
		table.NewTableCellWithTextAndDetail("Notifications", "On"),
		table.NewTableCellWithTextAndDetail("Theme", "Light"),
	}

	tv.Sections = []*table.TableSection{section1, section2}
	return tv
}

// ═══════════════════════════════════════════════════════════════
// DEBUG TAB - Console, Log
// ═══════════════════════════════════════════════════════════════

func createDebugTab() fyne.CanvasObject {
	consoleBtn := widget.NewButton("Toggle Console", func() {
		console.SharedConsole().Toggle(mainWindow)
	})

	logBtn := widget.NewButton("Add Log Entry", func() {
		log.QMUILog("Test log at %v", time.Now())
		toast.ShowMessage(mainWindow, "Log entry added")
	})

	clearBtn := widget.NewButton("Clear Console", func() {
		console.SharedConsole().Clear()
		toast.ShowMessage(mainWindow, "Console cleared")
	})

	return container.NewScroll(container.NewVBox(
		createSectionCard("console.Console", container.NewVBox(consoleBtn, clearBtn)),
		createSectionCard("log.QMUILog", logBtn),
	))
}

// ═══════════════════════════════════════════════════════════════
// CARD HELPERS
// ═══════════════════════════════════════════════════════════════

func createComponentCard(componentName, description string, content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(qmuiWhite)
	bg.CornerRadius = 10
	bg.StrokeWidth = 1
	bg.StrokeColor = primaryColor()

	nameLabel := canvas.NewText(componentName, primaryColor())
	nameLabel.TextSize = 14
	nameLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	descLabel := canvas.NewText(description, qmuiGrayText)
	descLabel.TextSize = 11

	sep := canvas.NewRectangle(color.RGBA{R: 230, G: 230, B: 230, A: 255})
	sep.SetMinSize(fyne.NewSize(0, 1))

	cardContent := container.NewVBox(
		container.NewPadded(container.NewVBox(nameLabel, descLabel)),
		sep,
		container.NewPadded(content),
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
