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

	"github.com/paul-hammant/qmui_fyne/alert"
	"github.com/paul-hammant/qmui_fyne/badge"
	"github.com/paul-hammant/qmui_fyne/button"
	"github.com/paul-hammant/qmui_fyne/checkbox"
	"github.com/paul-hammant/qmui_fyne/collection"
	"github.com/paul-hammant/qmui_fyne/dialog"
	"github.com/paul-hammant/qmui_fyne/empty"
	"github.com/paul-hammant/qmui_fyne/emotion"
	"github.com/paul-hammant/qmui_fyne/floatlayout"
	"github.com/paul-hammant/qmui_fyne/grid"
	"github.com/paul-hammant/qmui_fyne/label"
	"github.com/paul-hammant/qmui_fyne/marquee"
	"github.com/paul-hammant/qmui_fyne/modal"
	"github.com/paul-hammant/qmui_fyne/moreop"
	"github.com/paul-hammant/qmui_fyne/navigation"
	"github.com/paul-hammant/qmui_fyne/popup"
	"github.com/paul-hammant/qmui_fyne/progress"
	"github.com/paul-hammant/qmui_fyne/search"
	"github.com/paul-hammant/qmui_fyne/segmented"
	"github.com/paul-hammant/qmui_fyne/table"
	"github.com/paul-hammant/qmui_fyne/textfield"
	"github.com/paul-hammant/qmui_fyne/textview"
	"github.com/paul-hammant/qmui_fyne/tips"
	"github.com/paul-hammant/qmui_fyne/toast"
	"github.com/paul-hammant/qmui_fyne/console"
	"github.com/paul-hammant/qmui_fyne/core"
	"github.com/paul-hammant/qmui_fyne/theme"
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
		fyne.Do(func() {
			mainApp.Quit()
		})
	}()

	// Ensure clean exit when window is closed
	mainWindow.SetCloseIntercept(func() {
		fmt.Println("Window closed, exiting...")
		tourActive = false
		fyne.Do(func() {
			mainApp.Quit()
		})
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
	statusText.TextSize = 18
	statusText.TextStyle = fyne.TextStyle{Bold: true}
	statusBg := canvas.NewRectangle(color.RGBA{R: 40, G: 40, B: 40, A: 255})
	statusBg.SetMinSize(fyne.NewSize(0, 40))
	statusBar := container.NewStack(statusBg, container.NewCenter(statusText))

	// Main content tabs
	tabs = container.NewAppTabs(
		container.NewTabItem("Themes", createThemesTab()),
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
	fyne.Do(func() {
		statusText.Text = msg
		statusText.Refresh()
	})
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

	// Tab 0: Themes
	tourThemesTab()
	if !tourActive {
		return
	}

	// Tab 1: Components
	tourComponentsTab()
	if !tourActive {
		return
	}

	// Tab 2: Buttons
	tourButtonsTab()
	if !tourActive {
		return
	}

	// Tab 3: Text Input
	tourTextInputTab()
	if !tourActive {
		return
	}

	// Tab 4: Progress
	tourProgressTab()
	if !tourActive {
		return
	}

	// Tab 5: Dialogs
	tourDialogsTab()
	if !tourActive {
		return
	}

	// Tab 6: Navigation
	tourNavigationTab()
	if !tourActive {
		return
	}

	// Tab 7: Advanced
	tourAdvancedTab()
	if !tourActive {
		return
	}

	// Final
	setStatus("Tour complete! All components demonstrated.")
	tourActive = false
}

func selectTab(index int) {
	fyne.Do(func() {
		tabs.SelectIndex(index)
	})
	time.Sleep(shortPause)
}

// pause sleeps and returns true if tour should continue, false if stopped
func pause(d time.Duration) bool {
	time.Sleep(d)
	return tourActive
}

func scrollToCard(scroll *container.Scroll, cards []fyne.CanvasObject, index int) {
	if scroll == nil || index >= len(cards) {
		fmt.Printf("scrollToCard: scroll=%v, index=%d, len(cards)=%d - SKIPPING\n", scroll != nil, index, len(cards))
		return
	}
	// Calculate Y offset by summing heights of cards before this one
	var yOffset float32
	for i := 0; i < index; i++ {
		h := cards[i].MinSize().Height
		yOffset += h + 8 // 8 = VBox spacing
	}
	fmt.Printf("scrollToCard: index=%d, yOffset=%.0f, contentSize=%v, scrollSize=%v\n",
		index, yOffset, scroll.Content.MinSize(), scroll.Size())
	fyne.Do(func() {
		// Use ScrollToOffset (Fyne 2.6+) which properly updates the scroll position
		scroll.ScrollToOffset(fyne.NewPos(0, yOffset))

		// Move indicator triangle - position at fixed offset since card is scrolled to top
		if activeIndicator != nil {
			activeIndicator.Move(fyne.NewPos(2, 60))
			activeIndicator.Show()
			activeIndicator.Refresh()
		}
	})
}

func hideIndicator() {
	fyne.Do(func() {
		if activeIndicator != nil {
			activeIndicator.Hide()
			activeIndicator.Refresh()
		}
	})
}

// createIndicator creates an indicator triangle for a tab
func createIndicator() *canvas.Text {
	ind := canvas.NewText("▶", qmuiCyan)
	ind.TextSize = 20
	ind.TextStyle = fyne.TextStyle{Bold: true}
	ind.Hide()
	return ind
}

// wrapScrollWithIndicator wraps a scroll container with an indicator column
func wrapScrollWithIndicator(scroll *container.Scroll, indicator *canvas.Text) fyne.CanvasObject {
	indContainer := container.NewWithoutLayout(indicator)
	indContainer.Resize(fyne.NewSize(24, 600))
	return container.NewBorder(nil, nil, indContainer, nil, scroll)
}

// ============ Tour: Themes Tab ============

var (
	themesScroll    *container.Scroll
	themeCards      []fyne.CanvasObject
	themesIndicator *canvas.Text
	themeLabels     []*canvas.Text
)

func tourThemesTab() {
	setStatus("Tab 1/8: Themes - Hot-swappable color themes")
	selectTab(0)
	activeIndicator = themesIndicator
	if !pause(longPause) {
		hideIndicator()
		return
	}

	// Cycle through all themes
	allThemes := theme.SharedThemeManager().AllThemeIdentifiers()
	for i, themeID := range allThemes {
		if !tourActive {
			hideIndicator()
			return
		}
		t := theme.SharedThemeManager().GetTheme(themeID)
		setStatus(fmt.Sprintf("▶ Theme %d/%d: %s", i+1, len(allThemes), t.Name))
		fyne.Do(func() {
			theme.SharedThemeManager().SetCurrentTheme(themeID)
		})
		if !pause(mediumPause) {
			hideIndicator()
			return
		}
	}

	// Return to default
	setStatus("▶ Returning to Default theme")
	fyne.Do(func() {
		theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDefault)
	})
	if !pause(shortPause) {
		hideIndicator()
		return
	}

	hideIndicator()
}

// ============ Tour: Components Tab ============

var (
	componentsScroll    *container.Scroll
	componentCards      []fyne.CanvasObject
	componentsIndicator *canvas.Text

	// Demo components for interaction
	demoBadge1     *badge.Badge
	demoBadge2     *badge.Badge
	demoEmptyView  *empty.EmptyState
	demoGridRects  []*canvas.Rectangle

	// Current active indicator (shared reference)
	activeIndicator *canvas.Text
)

func tourComponentsTab() {
	setStatus("Tab 2/8: Components - Labels & Badges")
	selectTab(1)
	activeIndicator = componentsIndicator
	if !pause(longPause) { hideIndicator(); return }

	// 0: Marquee - already animating
	setStatus("▶ marquee.MarqueeLabel - text scrolls automatically")
	scrollToCard(componentsScroll, componentCards, 0)
	if !pause(longPause) { hideIndicator(); return }

	// 1: Badges - animate count changing
	setStatus("▶ badge.Badge - animating badge counts")
	scrollToCard(componentsScroll, componentCards, 1)
	if !pause(shortPause) { hideIndicator(); return }
	badgeTexts := []string{"5", "10", "25", "99+"}
	for _, text := range badgeTexts {
		if !tourActive { hideIndicator(); return }
		t := text
		fyne.Do(func() {
			if demoBadge1 != nil {
				demoBadge1.SetText(t)
			}
		})
		if !pause(500 * time.Millisecond) { hideIndicator(); return }
	}
	if !pause(shortPause) { hideIndicator(); return }

	// 2: Label - static display
	setStatus("▶ label.Label - padded labels with edge insets")
	scrollToCard(componentsScroll, componentCards, 2)
	if !pause(longPause) { hideIndicator(); return }

	// 3: FloatLayout - static display
	setStatus("▶ floatlayout.FloatLayoutView - tag cloud layout")
	scrollToCard(componentsScroll, componentCards, 3)
	if !pause(longPause) { hideIndicator(); return }

	// 4: Grid - highlight cells
	setStatus("▶ grid.GridView - flashing cells")
	scrollToCard(componentsScroll, componentCards, 4)
	if !pause(shortPause) { hideIndicator(); return }
	if len(demoGridRects) > 0 {
		for i := range demoGridRects {
			if !tourActive { hideIndicator(); return }
			idx := i
			originalColor := demoGridRects[idx].FillColor
			fyne.Do(func() {
				demoGridRects[idx].FillColor = color.White
				demoGridRects[idx].Refresh()
			})
			if !pause(300 * time.Millisecond) { hideIndicator(); return }
			fyne.Do(func() {
				demoGridRects[idx].FillColor = originalColor
				demoGridRects[idx].Refresh()
			})
			if !pause(200 * time.Millisecond) { hideIndicator(); return }
		}
	}
	if !pause(shortPause) { hideIndicator(); return }

	// 5: Empty - cycle through states
	setStatus("▶ empty.EmptyState - cycling through states")
	scrollToCard(componentsScroll, componentCards, 5)
	if !pause(shortPause) { hideIndicator(); return }

	setStatus("▶ empty.EmptyState - Loading...")
	fyne.Do(func() {
		if demoEmptyView != nil {
			demoEmptyView.SetLoading(true)
			demoEmptyView.SetText("Loading...")
		}
	})
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ empty.EmptyState - Error!")
	fyne.Do(func() {
		if demoEmptyView != nil {
			demoEmptyView.SetLoading(false)
			demoEmptyView.SetText("Error occurred")
			demoEmptyView.SetDetailText("Please try again")
		}
	})
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ empty.EmptyState - No Data")
	fyne.Do(func() {
		if demoEmptyView != nil {
			demoEmptyView.SetText("No data found")
			demoEmptyView.SetDetailText("Try a different search")
		}
	})
	if !pause(longPause) { hideIndicator(); return }

	// 6: Table - static display
	setStatus("▶ table.TableView - iOS-style grouped list")
	scrollToCard(componentsScroll, componentCards, 6)
	if !pause(longPause) { hideIndicator(); return }

	hideIndicator()
}

// ============ Tour: Buttons Tab ============

var (
	demoBtnStandard  *button.Button
	demoBtnFill      *button.FillButton
	demoBtnGhost     *button.GhostButton
	buttonsScroll    *container.Scroll
	buttonCards      []fyne.CanvasObject
	buttonsIndicator *canvas.Text
)

func tourButtonsTab() {
	setStatus("Tab 3/8: Buttons - Interactive button variants")
	selectTab(2)
	activeIndicator = buttonsIndicator
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Tapping: button.Button (standard)")
	scrollToCard(buttonsScroll, buttonCards, 0)
	if demoBtnStandard != nil {
		simulateTap(demoBtnStandard)
	}
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Tapping: button.FillButton (solid)")
	scrollToCard(buttonsScroll, buttonCards, 1)
	if demoBtnFill != nil {
		simulateTap(demoBtnFill)
	}
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Tapping: button.GhostButton (outlined)")
	scrollToCard(buttonsScroll, buttonCards, 2)
	if demoBtnGhost != nil {
		simulateTap(demoBtnGhost)
	}
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Showing: Navigation and Color variants")
	scrollToCard(buttonsScroll, buttonCards, 3)
	if !pause(longPause) { hideIndicator(); return }

	hideIndicator()
}

// ============ Tour: Text Input Tab ============

var (
	demoTextField      *textfield.TextField
	demoTextView       *textview.TextView
	demoSearchBar      *search.SearchBar
	textInputScroll    *container.Scroll
	textInputCards     []fyne.CanvasObject
	textInputIndicator *canvas.Text
)

func tourTextInputTab() {
	setStatus("Tab 4/8: Text Input - Fields and search")
	selectTab(3)
	activeIndicator = textInputIndicator
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Typing in: textfield.TextField")
	scrollToCard(textInputScroll, textInputCards, 0)
	if demoTextField != nil {
		simulateTyping(demoTextField, "Hello QMUI!")
	}
	if !pause(mediumPause) { hideIndicator(); return }

	setStatus("▶ Typing in: textview.TextView (multiline)")
	scrollToCard(textInputScroll, textInputCards, 1)
	if demoTextView != nil {
		simulateTypingTextView(demoTextView, "Multi-line\ntext input")
	}
	if !pause(mediumPause) { hideIndicator(); return }

	setStatus("▶ Searching in: search.SearchBar")
	scrollToCard(textInputScroll, textInputCards, 2)
	if demoSearchBar != nil {
		fyne.Do(func() {
			demoSearchBar.SetText("QMUI Go")
			demoSearchBar.Refresh()
		})
	}
	if !pause(longPause) { hideIndicator(); return }

	hideIndicator()
}

// ============ Tour: Progress Tab ============

var (
	demoPie           *progress.PieProgress
	demoCircular      *progress.RingProgress
	demoLinear        *progress.ProgressBar
	progressScroll    *container.Scroll
	progressCards     []fyne.CanvasObject
	progressIndicator *canvas.Text
)

func tourProgressTab() {
	setStatus("Tab 5/8: Progress - Animated indicators")
	selectTab(4)
	activeIndicator = progressIndicator
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Animating: progress.PieProgress")
	scrollToCard(progressScroll, progressCards, 0)
	animateProgressSingle(demoPie, nil, nil, 0.33)
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Animating: progress.RingProgress")
	scrollToCard(progressScroll, progressCards, 1)
	animateProgressSingle(nil, demoCircular, nil, 0.66)
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Animating: progress.ProgressBar")
	scrollToCard(progressScroll, progressCards, 2)
	animateProgressSingle(nil, nil, demoLinear, 1.0)
	if !pause(longPause) { hideIndicator(); return }

	// Reset all
	animateProgress(demoPie, demoCircular, demoLinear, 0)
	hideIndicator()
}

func animateProgressSingle(pie *progress.PieProgress, circ *progress.RingProgress, lin *progress.ProgressBar, target float64) {
	fyne.Do(func() {
		if pie != nil {
			pie.SetProgress(target)
		}
		if circ != nil {
			circ.SetProgress(target)
		}
		if lin != nil {
			lin.SetProgress(target)
		}
	})
}

func animateProgress(pie *progress.PieProgress, circ *progress.RingProgress, lin *progress.ProgressBar, target float64) {
	fyne.Do(func() {
		if pie != nil {
			pie.SetProgress(target)
		}
		if circ != nil {
			circ.SetProgress(target)
		}
		if lin != nil {
			lin.SetProgress(target)
		}
	})
}

// ============ Tour: Dialogs Tab ============

func tourDialogsTab() {
	setStatus("Tab 6/8: Dialogs - Toasts, Alerts, Popups")
	selectTab(5)
	if !pause(longPause) { return }

	setStatus("Showing: toast.ToastView")
	fyne.Do(func() {
		toast.ShowMessage(mainWindow, "This is a toast message!")
	})
	if !pause(longPause) { return }

	setStatus("Showing: tips.Tips (loading)")
	t := tips.NewHUD(mainWindow)
	fyne.Do(func() {
		t.ShowLoading("Loading...")
	})
	if !pause(mediumPause) { fyne.Do(func() { t.HideCurrent() }); return }
	fyne.Do(func() {
		t.HideCurrent()
	})

	setStatus("Showing: tips.Tips (success)")
	fyne.Do(func() {
		t.ShowSuccess("Success!")
	})
	if !pause(mediumPause) { fyne.Do(func() { t.HideCurrent() }); return }
	fyne.Do(func() {
		t.HideCurrent()
	})

	setStatus("Showing: alert.Alert")
	showDemoAlert()
	if !pause(longPause) { return }

	setStatus("Showing: alert.Alert (action sheet)")
	showDemoActionSheet()
	if !pause(longPause) { return }

	setStatus("Showing: dialog.DialogViewController")
	showDemoDialog()
	if !pause(longPause) { return }

	setStatus("Showing: popup.PopupMenu (context menu)")
	showDemoPopupMenu()
	if !pause(longPause) { return }

	setStatus("Showing: moreop.MoreOperationController")
	showDemoMoreOp()
	if !pause(longPause) { return }

	setStatus("Showing: modal.Modal")
	showDemoModal()
	if !pause(longPause) { return }
}

func showDemoAlert() {
	ac := alert.NewAlert("Alert Demo", "This demonstrates alert.Alert", alert.ControllerStyleAlert)
	ac.AddAction(alert.NewAction("OK", alert.ActionStyleDefault, func(_ *alert.Alert, _ *alert.Action) {}))
	fyne.Do(func() {
		ac.ShowIn(mainWindow)
		// Auto-dismiss after delay
		go func() {
			time.Sleep(longPause)
			fyne.Do(func() {
				ac.Hide()
			})
		}()
	})
}

func showDemoActionSheet() {
	ac := alert.NewAlert("Action Sheet", "Choose an option", alert.ControllerStyleActionSheet)
	ac.AddAction(alert.NewAction("Option 1", alert.ActionStyleDefault, nil))
	ac.AddAction(alert.NewAction("Option 2", alert.ActionStyleDefault, nil))
	ac.AddAction(alert.NewAction("Cancel", alert.ActionStyleCancel, nil))
	fyne.Do(func() {
		ac.ShowIn(mainWindow)
		go func() {
			time.Sleep(longPause)
			fyne.Do(func() {
				ac.Hide()
			})
		}()
	})
}

func showDemoDialog() {
	dvc := dialog.NewDialog()
	dvc.Title = "Dialog Demo"
	dvc.Message = "This is dialog.DialogViewController with custom content"
	dvc.AddAction(dialog.NewDialogActionWithHandler("Close", dialog.ActionStyleCancel, func(_ *dialog.DialogAction) {
		dvc.Dismiss()
	}))
	fyne.Do(func() {
		dvc.Show(mainWindow)
		go func() {
			time.Sleep(longPause)
			fyne.Do(func() {
				dvc.Dismiss()
			})
		}()
	})
}

func showDemoPopupMenu() {
	items := []*popup.MenuItem{
		popup.NewMenuItem("Edit", nil),
		popup.NewMenuItem("Copy", nil),
		popup.NewMenuItem("Delete", nil),
	}
	pm := popup.NewPopupMenuWithItems(items)
	fyne.Do(func() {
		pm.Show(mainWindow, fyne.NewPos(200, 400))
		go func() {
			time.Sleep(longPause)
			fyne.Do(func() {
				pm.Hide()
			})
		}()
	})
}

func showDemoMoreOp() {
	items := []*moreop.Item{
		moreop.NewItem("share", "Share", nil, nil),
		moreop.NewItem("copy", "Copy", nil, nil),
		moreop.NewItem("save", "Save", nil, nil),
		moreop.NewItem("delete", "Delete", nil, nil),
	}
	ctrl := moreop.NewActionSheet()
	ctrl.AddItems(items...)
	fyne.Do(func() {
		ctrl.Show(mainWindow)
		go func() {
			time.Sleep(longPause)
			fyne.Do(func() {
				ctrl.Dismiss()
			})
		}()
	})
}

func showDemoModal() {
	content := container.NewVBox(
		widget.NewLabel("Modal Content"),
		widget.NewLabel("This slides up from bottom"),
	)
	var mvc *modal.Modal
	fyne.Do(func() {
		mvc = modal.PresentModalFromBottom(mainWindow, container.NewPadded(content))
		// Start dismiss timer after modal is shown
		go func() {
			time.Sleep(longPause) // Give user time to read
			fyne.Do(func() {
				if mvc != nil {
					mvc.Dismiss()
				}
			})
		}()
	})
}

// ============ Tour: Navigation Tab ============

var (
	demoSegmented  *segmented.SegmentedControl
	demoCheckbox1  *checkbox.Checkbox
	demoCheckbox2  *checkbox.Checkbox
	navScroll      *container.Scroll
	navCards       []fyne.CanvasObject
	navIndicator   *canvas.Text
)

func tourNavigationTab() {
	setStatus("Tab 7/8: Navigation - Bars and controls")
	selectTab(6)
	activeIndicator = navIndicator
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Showing: navigation.NavigationBar")
	scrollToCard(navScroll, navCards, 0)
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Showing: navigation.TabBar")
	scrollToCard(navScroll, navCards, 1)
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Toggling: segmented.SegmentedControl")
	scrollToCard(navScroll, navCards, 2)
	if demoSegmented != nil {
		fyne.Do(func() {
			demoSegmented.SetSelectedIndex(1)
		})
		if !pause(shortPause) { hideIndicator(); return }
		fyne.Do(func() {
			demoSegmented.SetSelectedIndex(2)
		})
		if !pause(shortPause) { hideIndicator(); return }
		fyne.Do(func() {
			demoSegmented.SetSelectedIndex(0)
		})
	}
	if !pause(mediumPause) { hideIndicator(); return }

	setStatus("▶ Toggling: checkbox.Checkbox")
	scrollToCard(navScroll, navCards, 3)
	if demoCheckbox1 != nil {
		fyne.Do(func() {
			demoCheckbox1.SetSelected(true)
		})
		if !pause(shortPause) { hideIndicator(); return }
	}
	if demoCheckbox2 != nil {
		fyne.Do(func() {
			demoCheckbox2.SetSelected(true)
		})
	}
	if !pause(longPause) { hideIndicator(); return }

	hideIndicator()
}

// ============ Tour: Advanced Tab ============

var (
	demoPaging   *collection.PagingLayout
	advScroll    *container.Scroll
	advCards     []fyne.CanvasObject
	advIndicator *canvas.Text
)

func tourAdvancedTab() {
	setStatus("Tab 8/8: Advanced - Console & Special")
	selectTab(7)
	activeIndicator = advIndicator
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Showing: emotion.EmotionView (emoji picker)")
	scrollToCard(advScroll, advCards, 0)
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Swiping: collection.PagingLayout")
	scrollToCard(advScroll, advCards, 1)
	if demoPaging != nil {
		fyne.Do(func() {
			demoPaging.SetCurrentPage(1)
		})
		if !pause(shortPause) { hideIndicator(); return }
		fyne.Do(func() {
			demoPaging.SetCurrentPage(2)
		})
		if !pause(shortPause) { hideIndicator(); return }
		fyne.Do(func() {
			demoPaging.SetCurrentPage(0)
		})
	}
	if !pause(mediumPause) { hideIndicator(); return }

	setStatus("▶ Switching to: Dark Theme")
	scrollToCard(advScroll, advCards, 2)
	fyne.Do(func() {
		theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDark)
	})
	if !pause(longPause) { hideIndicator(); return }

	setStatus("▶ Switching to: Light Theme")
	fyne.Do(func() {
		theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDefault)
	})
	if !pause(mediumPause) { hideIndicator(); return }

	setStatus("▶ Opening: console.Console")
	scrollToCard(advScroll, advCards, 3)
	fyne.Do(func() {
		console.SharedConsole().ShowIn(mainWindow)
	})
	if !pause(longPause) { hideIndicator(); return }
	fyne.Do(func() {
		console.SharedConsole().Hide()
	})
	if !pause(mediumPause) { hideIndicator(); return }

	hideIndicator()
}

// ============ Tab Content Creators ============

func createThemesTab() fyne.CanvasObject {
	// Theme color swatches grid
	allThemes := theme.SharedThemeManager().AllThemes()
	themeLabels = nil

	swatchGrid := container.NewGridWithColumns(3)
	for _, t := range allThemes {
		currentTheme := t // capture for closure
		swatch := canvas.NewRectangle(t.PrimaryColor)
		swatch.CornerRadius = 8
		swatch.SetMinSize(fyne.NewSize(80, 50))

		nameLabel := canvas.NewText(t.Name, qmuiDarkText)
		nameLabel.TextSize = 12
		nameLabel.Alignment = fyne.TextAlignCenter
		themeLabels = append(themeLabels, nameLabel)

		// Make swatch tappable
		tappable := widget.NewButton("", func() {
			theme.SharedThemeManager().SetCurrentTheme(currentTheme.Identifier)
		})
		tappable.Importance = widget.LowImportance

		swatchCard := container.NewVBox(
			container.NewStack(swatch, tappable),
			nameLabel,
		)
		swatchGrid.Add(swatchCard)
	}

	// Cycle button
	cycleBtn := widget.NewButton("Cycle Theme", func() {
		theme.SharedThemeManager().CycleTheme()
	})

	// Current theme indicator
	currentThemeLabel := canvas.NewText("Current: Default", qmuiCyan)
	currentThemeLabel.TextSize = 16
	currentThemeLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Listen for theme changes
	theme.SharedThemeManager().AddThemeChangeListener(func(t *theme.Theme) {
		fyne.Do(func() {
			currentThemeLabel.Text = "Current: " + t.Name
			currentThemeLabel.Refresh()
		})
	})

	// Sample widgets to show theme effect
	sampleProgress := progress.NewPieProgress()
	sampleProgress.Progress = 0.65
	sampleProgress.TintColor = theme.SharedThemeManager().CurrentTheme().PrimaryColor
	theme.SharedThemeManager().AddThemeChangeListener(func(t *theme.Theme) {
		fyne.Do(func() {
			sampleProgress.TintColor = t.PrimaryColor
			sampleProgress.Refresh()
		})
	})

	sampleBtn := button.NewFillButton("Theme Button", theme.SharedThemeManager().CurrentTheme().PrimaryColor, nil)
	theme.SharedThemeManager().AddThemeChangeListener(func(t *theme.Theme) {
		fyne.Do(func() {
			sampleBtn.BackgroundColor = t.PrimaryColor
			sampleBtn.Refresh()
		})
	})

	themeCards = []fyne.CanvasObject{
		createCard("Theme Gallery", "Tap any swatch to switch themes", container.NewPadded(swatchGrid)),
		createCard("Theme Controls", "Cycle through all themes", container.NewVBox(currentThemeLabel, cycleBtn)),
		createCard("Sample Widgets", "These update with theme", container.NewHBox(sampleProgress, sampleBtn)),
	}
	themesIndicator = createIndicator()
	themesScroll = container.NewScroll(container.NewVBox(themeCards...))
	return wrapScrollWithIndicator(themesScroll, themesIndicator)
}

func createComponentsTab() fyne.CanvasObject {
	componentCards = []fyne.CanvasObject{
		createCard("marquee.MarqueeLabel", "Scrolling text animation", createMarquee()),
		createCard("badge.Badge", "Notification badges", createBadges()),
		createCard("label.Label", "Label with edge insets", createLabel()),
		createCard("floatlayout.FloatLayoutView", "Tag cloud layout", createFloatLayout()),
		createCard("grid.GridView", "Grid arrangement", createGrid()),
		createCard("empty.EmptyState", "Loading state", createEmpty()),
		createCard("table.TableView", "iOS-style grouped list", createTable()),
	}

	componentsIndicator = createIndicator()
	componentsScroll = container.NewScroll(container.NewVBox(componentCards...))
	return wrapScrollWithIndicator(componentsScroll, componentsIndicator)
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

	buttonCards = []fyne.CanvasObject{
		createCard("button.Button", "Standard tappable button", demoBtnStandard),
		createCard("button.FillButton", "Solid filled button", demoBtnFill),
		createCard("button.GhostButton", "Outlined button", demoBtnGhost),
		createCard("button.NavigationButton", "Navigation bar button", navBtn),
		createCard("Color Variants", "Multiple tint colors", colorBtns),
	}
	buttonsIndicator = createIndicator()
	buttonsScroll = container.NewScroll(container.NewVBox(buttonCards...))
	return wrapScrollWithIndicator(buttonsScroll, buttonsIndicator)
}

func createTextInputTab() fyne.CanvasObject {
	demoTextField = textfield.NewTextFieldWithPlaceholder("Enter text here...")
	demoTextField.PlaceholderColor = qmuiGrayText

	demoTextView = textview.NewTextView()
	demoTextView.PlaceHolder = "Multi-line text input..."

	demoSearchBar = search.NewSearchBar()
	demoSearchBar.Placeholder = "Search..."

	textInputCards = []fyne.CanvasObject{
		createCard("textfield.TextField", "Single-line input", demoTextField),
		createCard("textview.TextView", "Multi-line input", demoTextView),
		createCard("search.SearchBar", "Search with suggestions", demoSearchBar),
	}
	textInputIndicator = createIndicator()
	textInputScroll = container.NewScroll(container.NewVBox(textInputCards...))
	return wrapScrollWithIndicator(textInputScroll, textInputIndicator)
}

func createProgressTab() fyne.CanvasObject {
	demoPie = progress.NewPieProgress()
	demoPie.Progress = 0
	demoPie.TintColor = qmuiCyan

	demoCircular = progress.NewRingProgress()
	demoCircular.Progress = 0
	demoCircular.TintColor = color.RGBA{R: 52, G: 199, B: 89, A: 255}
	demoCircular.ShowsText = true

	demoLinear = progress.NewProgressBar()
	demoLinear.Progress = 0
	demoLinear.TintColor = color.RGBA{R: 255, G: 149, B: 0, A: 255}

	progressCards = []fyne.CanvasObject{
		createCard("progress.PieProgress", "Pie chart progress", demoPie),
		createCard("progress.RingProgress", "Ring with percentage", demoCircular),
		createCard("progress.ProgressBar", "Horizontal bar", demoLinear),
	}
	progressIndicator = createIndicator()
	progressScroll = container.NewScroll(container.NewVBox(progressCards...))
	return wrapScrollWithIndicator(progressScroll, progressIndicator)
}

func createDialogsTab() fyne.CanvasObject {
	toastBtn := widget.NewButton("Show Toast", func() {
		toast.ShowMessage(mainWindow, "Toast message!")
	})

	tipsBtn := widget.NewButton("Show Tips", func() {
		t := tips.NewHUD(mainWindow)
		t.ShowSuccess("Operation completed!")
		go func() {
			time.Sleep(2 * time.Second)
			t.HideCurrent()
		}()
	})

	alertBtn := widget.NewButton("Show Alert", func() {
		ac := alert.NewAlert("Alert", "This is an alert dialog", alert.ControllerStyleAlert)
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
		items := []*moreop.Item{
			moreop.NewItem("share", "Share", nil, func(_ *moreop.Item) {
				toast.ShowMessage(mainWindow, "Share tapped")
			}),
			moreop.NewItem("copy", "Copy", nil, func(_ *moreop.Item) {
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
		createCard("alert.Alert", "Alert dialogs", alertBtn),
		createCard("dialog.DialogViewController", "Custom dialogs", dialogBtn),
		createCard("popup.PopupMenu", "Context menus", popupBtn),
		createCard("moreop.MoreOperationController", "Action grid", moreOpBtn),
		createCard("modal.Modal", "Slide-up modal", modalBtn),
	}
	return container.NewScroll(container.NewVBox(cards...))
}

func createNavigationTab() fyne.CanvasObject {
	// Navigation bar
	navBar := navigation.NewNavigationBar()
	navBar.SetTitleView(navigation.NewTitleViewWithTitle("Navigation Bar"))
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

	navCards = []fyne.CanvasObject{
		createCard("navigation.NavigationBar", "App navigation bar", navBar),
		createCard("navigation.TabBar", "Bottom tab bar", tabBar),
		createCard("segmented.SegmentedControl", "Segmented selector", demoSegmented),
		createCard("checkbox.Checkbox", "Selection checkboxes", checkboxes),
	}
	navIndicator = createIndicator()
	navScroll = container.NewScroll(container.NewVBox(navCards...))
	return wrapScrollWithIndicator(navScroll, navIndicator)
}

func createAdvancedTab() fyne.CanvasObject {
	// Emotion view
	emotionView := emotion.NewEmojiPicker()

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

	advCards = []fyne.CanvasObject{
		createCard("emotion.EmotionView", "Emoji picker grid", emotionView),
		createCard("collection.PagingLayout", "Swipeable pages", demoPaging),
		createCard("theme.ThemeManager", "Hot-switchable themes", themeRow),
		createCard("console.Console", "Debug console", consoleBtn),
	}
	advIndicator = createIndicator()
	advScroll = container.NewScroll(container.NewVBox(advCards...))
	return wrapScrollWithIndicator(advScroll, advIndicator)
}

// ============ Component Creators ============

func createMarquee() fyne.CanvasObject {
	// Use long text to ensure scrolling even on wide windows
	m := marquee.NewMarquee("This text scrolls continuously across the screen - MarqueeLabel provides smooth animated scrolling text perfect for news tickers, announcements, and attention-grabbing displays!")
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
	demoBadge1 = badge.NewBadge("3")
	demoBadge2 = badge.NewBadge("NEW")
	return container.NewHBox(demoBadge1, demoBadge2, badge.NewBadge("99+"))
}

func createLabel() fyne.CanvasObject {
	lbl := label.NewLabel("Padded Label")
	lbl.ContentEdgeInsets = core.EdgeInsets{Top: 8, Left: 16, Bottom: 8, Right: 16}
	bg := canvas.NewRectangle(color.RGBA{R: 230, G: 245, B: 255, A: 255})
	bg.CornerRadius = 4
	return container.NewStack(bg, lbl)
}

func createFloatLayout() fyne.CanvasObject {
	tc := floatlayout.NewTagCloud()
	tc.SetTags([]string{"Go", "Fyne", "QMUI", "iOS", "Cross-Platform"})
	return tc
}

func createGrid() fyne.CanvasObject {
	gv := grid.NewGrid(4)
	gv.RowSpacing = 4
	gv.ColumnSpacing = 4
	colors := []color.Color{qmuiCyan, color.RGBA{R: 255, G: 100, B: 100, A: 255},
		color.RGBA{R: 100, G: 200, B: 100, A: 255}, color.RGBA{R: 200, G: 150, B: 255, A: 255}}
	demoGridRects = nil
	for _, c := range colors {
		rect := canvas.NewRectangle(c)
		rect.CornerRadius = 4
		rect.SetMinSize(fyne.NewSize(30, 30))
		demoGridRects = append(demoGridRects, rect)
		gv.AddItem(grid.NewGridItem(rect))
	}
	return gv
}

func createEmpty() fyne.CanvasObject {
	demoEmptyView = empty.NewEmptyState()
	demoEmptyView.IsLoading = true
	demoEmptyView.Text = "Loading..."
	return demoEmptyView
}

func createTable() fyne.CanvasObject {
	tv := table.NewTable(table.TableStyleInsetGrouped)
	section := table.NewTableSection("Settings")
	section.Cells = []*table.TableCell{
		table.NewTableCellWithTextAndDetail("Profile", "View"),
		table.NewTableCellWithTextAndDetail("Notifications", "On"),
	}
	tv.Sections = []*table.TableSection{section}
	return tv
}

// ============ Interaction Helpers ============

func simulateTap(tappable fyne.Tappable) {
	if tappable != nil {
		fyne.Do(func() {
			tappable.Tapped(&fyne.PointEvent{})
		})
	}
}

func simulateTyping(tf *textfield.TextField, text string) {
	if tf == nil {
		return
	}
	for _, ch := range text {
		c := ch // capture for closure
		fyne.Do(func() {
			tf.SetText(tf.Text + string(c))
		})
		time.Sleep(50 * time.Millisecond)
	}
}

func simulateTypingTextView(tv *textview.TextView, text string) {
	if tv == nil {
		return
	}
	for _, ch := range text {
		c := ch // capture for closure
		fyne.Do(func() {
			tv.SetText(tv.Text + string(c))
		})
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
