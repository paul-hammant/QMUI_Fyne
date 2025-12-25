// Package core provides the core configuration and utilities for QMUI-Go
// Ported from Tencent's QMUI_iOS framework
// Licensed under the MIT License
package core

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
)

// Configuration is the global singleton for UI configuration
// Mirrors QMUIConfiguration from iOS
type Configuration struct {
	active bool
	mu     sync.RWMutex

	// Global Colors
	ClearColor       color.Color
	WhiteColor       color.Color
	BlackColor       color.Color
	GrayColor        color.Color
	GrayDarkenColor  color.Color
	GrayLightenColor color.Color
	RedColor         color.Color
	GreenColor       color.Color
	BlueColor        color.Color
	YellowColor      color.Color

	// Semantic Colors
	LinkColor           color.Color
	DisabledColor       color.Color
	BackgroundColor     color.Color
	MaskDarkColor       color.Color
	MaskLightColor      color.Color
	SeparatorColor      color.Color
	SeparatorDashedColor color.Color
	PlaceholderColor    color.Color

	// Test Colors
	TestColorRed   color.Color
	TestColorGreen color.Color
	TestColorBlue  color.Color

	// UIControl
	ControlHighlightedAlpha float64
	ControlDisabledAlpha    float64

	// Button
	ButtonHighlightedAlpha float64
	ButtonDisabledAlpha    float64
	ButtonTintColor        color.Color

	// TextField & TextView
	TextFieldTextColor color.Color
	TextFieldTintColor color.Color
	TextFieldTextInsets EdgeInsets
	KeyboardAppearance  KeyboardAppearance

	// Switch
	SwitchOnTintColor    color.Color
	SwitchOffTintColor   color.Color
	SwitchThumbTintColor color.Color

	// NavigationBar
	NavBarHighlightedAlpha                  float64
	NavBarDisabledAlpha                     float64
	NavBarButtonFont                        fyne.TextStyle
	NavBarButtonFontBold                    fyne.TextStyle
	NavBarBackgroundColor                   color.Color
	NavBarShadowColor                       color.Color
	NavBarBarTintColor                      color.Color
	NavBarTintColor                         color.Color
	NavBarTitleColor                        color.Color
	NavBarTitleFontSize                     float32
	NavBarLargeTitleColor                   color.Color
	NavBarLargeTitleFontSize                float32
	NavBarLoadingMarginRight                float32
	NavBarAccessoryViewMarginLeft           float32

	// TabBar
	TabBarBackgroundColor         color.Color
	TabBarBarTintColor            color.Color
	TabBarShadowColor             color.Color
	TabBarItemTitleFontSize       float32
	TabBarItemTitleFontSizeSelected float32
	TabBarItemTitleColor          color.Color
	TabBarItemTitleColorSelected  color.Color
	TabBarItemImageColor          color.Color
	TabBarItemImageColorSelected  color.Color

	// Toolbar
	ToolBarHighlightedAlpha    float64
	ToolBarDisabledAlpha       float64
	ToolBarTintColor           color.Color
	ToolBarTintColorHighlighted color.Color
	ToolBarTintColorDisabled   color.Color
	ToolBarBackgroundColor     color.Color
	ToolBarBarTintColor        color.Color
	ToolBarShadowColor         color.Color
	ToolBarButtonFontSize      float32

	// SearchBar
	SearchBarTextFieldBackgroundColor color.Color
	SearchBarTextFieldBorderColor     color.Color
	SearchBarBackgroundColor          color.Color
	SearchBarTintColor                color.Color
	SearchBarTextColor                color.Color
	SearchBarPlaceholderColor         color.Color
	SearchBarFontSize                 float32
	SearchBarTextFieldCornerRadius    float32

	// TableView / TableViewCell
	TableViewEstimatedHeightEnabled                bool
	TableViewBackgroundColor                       color.Color
	TableSectionIndexColor                         color.Color
	TableSectionIndexBackgroundColor               color.Color
	TableSectionIndexTrackingBackgroundColor       color.Color
	TableViewSeparatorColor                        color.Color
	TableViewCellNormalHeight                      float32
	TableViewCellTitleLabelColor                   color.Color
	TableViewCellDetailLabelColor                  color.Color
	TableViewCellBackgroundColor                   color.Color
	TableViewCellSelectedBackgroundColor           color.Color
	TableViewCellWarningBackgroundColor            color.Color
	TableViewSectionHeaderBackgroundColor          color.Color
	TableViewSectionFooterBackgroundColor          color.Color
	TableViewSectionHeaderFontSize                 float32
	TableViewSectionFooterFontSize                 float32
	TableViewSectionHeaderTextColor                color.Color
	TableViewSectionFooterTextColor                color.Color
	TableViewSectionHeaderAccessoryMargins         EdgeInsets
	TableViewSectionFooterAccessoryMargins         EdgeInsets
	TableViewSectionHeaderContentInset             EdgeInsets
	TableViewSectionFooterContentInset             EdgeInsets

	// Grouped TableView
	TableViewGroupedBackgroundColor                color.Color
	TableViewGroupedSeparatorColor                 color.Color
	TableViewGroupedCellTitleLabelColor            color.Color
	TableViewGroupedCellDetailLabelColor           color.Color
	TableViewGroupedCellBackgroundColor            color.Color
	TableViewGroupedCellSelectedBackgroundColor    color.Color
	TableViewGroupedCellWarningBackgroundColor     color.Color
	TableViewGroupedSectionHeaderFontSize          float32
	TableViewGroupedSectionFooterFontSize          float32
	TableViewGroupedSectionHeaderTextColor         color.Color
	TableViewGroupedSectionFooterTextColor         color.Color
	TableViewGroupedSectionHeaderDefaultHeight     float32
	TableViewGroupedSectionFooterDefaultHeight     float32

	// Inset Grouped TableView
	TableViewInsetGroupedCornerRadius              float32
	TableViewInsetGroupedHorizontalInset           float32
	TableViewInsetGroupedBackgroundColor           color.Color
	TableViewInsetGroupedSeparatorColor            color.Color
	TableViewInsetGroupedCellTitleLabelColor       color.Color
	TableViewInsetGroupedCellDetailLabelColor      color.Color
	TableViewInsetGroupedCellBackgroundColor       color.Color
	TableViewInsetGroupedCellSelectedBackgroundColor color.Color

	// WindowLevel (for overlays)
	WindowLevelQMUIAlertView float32
	WindowLevelQMUIConsole   float32

	// QMUILog
	ShouldPrintDefaultLog        bool
	ShouldPrintInfoLog           bool
	ShouldPrintWarnLog           bool
	ShouldPrintQMUIWarnLogToConsole bool

	// QMUIBadge
	BadgeBackgroundColor    color.Color
	BadgeTextColor          color.Color
	BadgeFontSize           float32
	BadgeContentEdgeInsets  EdgeInsets
	BadgeOffset             Offset
	BadgeOffsetLandscape    Offset
	UpdatesIndicatorColor   color.Color
	UpdatesIndicatorSize    fyne.Size
	UpdatesIndicatorOffset  Offset
	UpdatesIndicatorOffsetLandscape Offset

	// Others
	AutomaticCustomNavigationBarTransitionStyle bool
	NeedsBackBarButtonItemTitle                 bool
	HidesBottomBarWhenPushedInitially           bool
	PreventConcurrentNavigationControllerTransitions bool
	NavigationBarHiddenInitially                bool

	// Alert Controller
	AlertContentMargin            EdgeInsets
	AlertContentMaximumWidth      float32
	AlertSeparatorColor           color.Color
	AlertContentCornerRadius      float32
	AlertButtonHeight             float32
	AlertHeaderBackgroundColor    color.Color
	AlertButtonBackgroundColor    color.Color
	AlertButtonHighlightBackgroundColor color.Color
	AlertHeaderInsets             EdgeInsets
	AlertTitleMessageSpacing      float32
	AlertTextFieldFontSize        float32
	AlertTextFieldTextColor       color.Color
	AlertTextFieldBorderColor     color.Color
	AlertTextFieldTextInsets      EdgeInsets

	// Sheet
	SheetContentMargin              EdgeInsets
	SheetContentMaximumWidth        float32
	SheetSeparatorColor             color.Color
	SheetCancelButtonMarginTop      float32
	SheetContentCornerRadius        float32
	SheetButtonHeight               float32
	SheetHeaderBackgroundColor      color.Color
	SheetButtonBackgroundColor      color.Color
	SheetButtonHighlightBackgroundColor color.Color
	SheetHeaderInsets               EdgeInsets
	SheetTitleMessageSpacing        float32
	SheetButtonColumnCount          int

	// Toast
	ToastBackgroundColor          color.Color
	ToastTextColor                color.Color
	ToastFontSize                 float32
	ToastCornerRadius             float32
	ToastContentInsets            EdgeInsets
	ToastMarginFromScreen         float32
	ToastDefaultDuration          float64

	// EmptyView
	EmptyViewImageTintColor       color.Color
	EmptyViewLoadingTintColor     color.Color
	EmptyViewTextLabelColor       color.Color
	EmptyViewDetailTextLabelColor color.Color
	EmptyViewActionButtonColor    color.Color
	EmptyViewTextFontSize         float32
	EmptyViewDetailTextFontSize   float32
	EmptyViewActionButtonFontSize float32
}

// EdgeInsets represents padding/margins similar to UIEdgeInsets
type EdgeInsets struct {
	Top    float32
	Left   float32
	Bottom float32
	Right  float32
}

// Offset represents a 2D offset
type Offset struct {
	X float32
	Y float32
}

// KeyboardAppearance defines the keyboard appearance style
type KeyboardAppearance int

const (
	KeyboardAppearanceDefault KeyboardAppearance = iota
	KeyboardAppearanceDark
	KeyboardAppearanceLight
)

// Zero returns an empty EdgeInsets
func (e EdgeInsets) Zero() EdgeInsets {
	return EdgeInsets{}
}

// NewEdgeInsets creates a new EdgeInsets
func NewEdgeInsets(top, left, bottom, right float32) EdgeInsets {
	return EdgeInsets{Top: top, Left: left, Bottom: bottom, Right: right}
}

// NewOffset creates a new Offset
func NewOffset(x, y float32) Offset {
	return Offset{X: x, Y: y}
}

var (
	configInstance *Configuration
	configOnce     sync.Once
)

// SharedConfiguration returns the singleton Configuration instance
func SharedConfiguration() *Configuration {
	configOnce.Do(func() {
		configInstance = &Configuration{}
		configInstance.applyDefaults()
	})
	return configInstance
}

// ResetConfigurationForTesting resets the configuration for testing purposes
// This should only be used in tests
func ResetConfigurationForTesting() {
	configOnce = sync.Once{}
	configInstance = nil
}

// IsActive returns whether the configuration is active
func (c *Configuration) IsActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.active
}

// Activate marks the configuration as active
func (c *Configuration) Activate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.active = true
}

// applyDefaults sets all default values
func (c *Configuration) applyDefaults() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Global Colors
	c.ClearColor = color.Transparent
	c.WhiteColor = color.White
	c.BlackColor = color.Black
	c.GrayColor = color.RGBA{R: 179, G: 179, B: 179, A: 255}
	c.GrayDarkenColor = color.RGBA{R: 163, G: 163, B: 163, A: 255}
	c.GrayLightenColor = color.RGBA{R: 198, G: 198, B: 198, A: 255}
	c.RedColor = color.RGBA{R: 250, G: 58, B: 58, A: 255}
	c.GreenColor = color.RGBA{R: 159, G: 214, B: 97, A: 255}
	c.BlueColor = color.RGBA{R: 49, G: 189, B: 243, A: 255} // QMUI signature cyan #31BDF3
	c.YellowColor = color.RGBA{R: 255, G: 207, B: 71, A: 255}

	// Semantic Colors
	c.LinkColor = color.RGBA{R: 56, G: 116, B: 171, A: 255}
	c.DisabledColor = color.RGBA{R: 179, G: 179, B: 179, A: 255}
	c.BackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	c.MaskDarkColor = color.RGBA{R: 0, G: 0, B: 0, A: 128}
	c.MaskLightColor = color.RGBA{R: 255, G: 255, B: 255, A: 128}
	c.SeparatorColor = color.RGBA{R: 222, G: 224, B: 226, A: 255}
	c.SeparatorDashedColor = color.RGBA{R: 222, G: 224, B: 226, A: 255}
	c.PlaceholderColor = color.RGBA{R: 196, G: 200, B: 208, A: 255}

	// Test Colors
	c.TestColorRed = color.RGBA{R: 255, G: 0, B: 0, A: 64}
	c.TestColorGreen = color.RGBA{R: 0, G: 255, B: 0, A: 64}
	c.TestColorBlue = color.RGBA{R: 0, G: 0, B: 255, A: 64}

	// UIControl
	c.ControlHighlightedAlpha = 0.5
	c.ControlDisabledAlpha = 0.5

	// Button
	c.ButtonHighlightedAlpha = 0.5
	c.ButtonDisabledAlpha = 0.5
	c.ButtonTintColor = c.BlueColor

	// TextField & TextView
	c.TextFieldTextColor = c.BlackColor
	c.TextFieldTintColor = c.BlueColor
	c.TextFieldTextInsets = NewEdgeInsets(0, 7, 0, 7)
	c.KeyboardAppearance = KeyboardAppearanceDefault

	// Switch
	c.SwitchOnTintColor = c.GreenColor
	c.SwitchOffTintColor = nil
	c.SwitchThumbTintColor = nil

	// NavigationBar
	c.NavBarHighlightedAlpha = 0.2
	c.NavBarDisabledAlpha = 0.2
	c.NavBarButtonFont = fyne.TextStyle{}
	c.NavBarButtonFontBold = fyne.TextStyle{Bold: true}
	c.NavBarBackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	c.NavBarShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: 77}
	c.NavBarBarTintColor = nil
	c.NavBarTintColor = c.BlueColor
	c.NavBarTitleColor = color.RGBA{R: 51, G: 51, B: 51, A: 255}
	c.NavBarTitleFontSize = 17
	c.NavBarLargeTitleColor = c.BlackColor
	c.NavBarLargeTitleFontSize = 34
	c.NavBarLoadingMarginRight = 3
	c.NavBarAccessoryViewMarginLeft = 5

	// TabBar
	c.TabBarBackgroundColor = color.RGBA{R: 249, G: 249, B: 249, A: 255}
	c.TabBarBarTintColor = nil
	c.TabBarShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: 77}
	c.TabBarItemTitleFontSize = 10
	c.TabBarItemTitleFontSizeSelected = 10
	c.TabBarItemTitleColor = color.RGBA{R: 153, G: 153, B: 153, A: 255}
	c.TabBarItemTitleColorSelected = c.BlueColor
	c.TabBarItemImageColor = color.RGBA{R: 153, G: 153, B: 153, A: 255}
	c.TabBarItemImageColorSelected = c.BlueColor

	// Toolbar
	c.ToolBarHighlightedAlpha = 0.4
	c.ToolBarDisabledAlpha = 0.4
	c.ToolBarTintColor = c.BlueColor
	c.ToolBarTintColorHighlighted = nil
	c.ToolBarTintColorDisabled = nil
	c.ToolBarBackgroundColor = nil
	c.ToolBarBarTintColor = nil
	c.ToolBarShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: 77}
	c.ToolBarButtonFontSize = 17

	// SearchBar
	c.SearchBarTextFieldBackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	c.SearchBarTextFieldBorderColor = color.RGBA{R: 205, G: 208, B: 210, A: 255}
	c.SearchBarBackgroundColor = color.RGBA{R: 232, G: 232, B: 232, A: 255}
	c.SearchBarTintColor = c.BlueColor
	c.SearchBarTextColor = c.BlackColor
	c.SearchBarPlaceholderColor = c.PlaceholderColor
	c.SearchBarFontSize = 14
	c.SearchBarTextFieldCornerRadius = 4

	// TableView
	c.TableViewEstimatedHeightEnabled = true
	c.TableViewBackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	c.TableSectionIndexColor = c.GrayDarkenColor
	c.TableSectionIndexBackgroundColor = color.Transparent
	c.TableSectionIndexTrackingBackgroundColor = color.Transparent
	c.TableViewSeparatorColor = c.SeparatorColor
	c.TableViewCellNormalHeight = 56
	c.TableViewCellTitleLabelColor = color.RGBA{R: 51, G: 51, B: 51, A: 255}
	c.TableViewCellDetailLabelColor = c.GrayColor
	c.TableViewCellBackgroundColor = c.WhiteColor
	c.TableViewCellSelectedBackgroundColor = color.RGBA{R: 238, G: 239, B: 241, A: 255}
	c.TableViewCellWarningBackgroundColor = c.YellowColor
	c.TableViewSectionHeaderBackgroundColor = color.RGBA{R: 244, G: 244, B: 244, A: 255}
	c.TableViewSectionFooterBackgroundColor = color.RGBA{R: 244, G: 244, B: 244, A: 255}
	c.TableViewSectionHeaderFontSize = 14
	c.TableViewSectionFooterFontSize = 12
	c.TableViewSectionHeaderTextColor = c.GrayDarkenColor
	c.TableViewSectionFooterTextColor = c.GrayColor

	// Grouped TableView
	c.TableViewGroupedBackgroundColor = color.RGBA{R: 246, G: 246, B: 246, A: 255}
	c.TableViewGroupedSeparatorColor = c.SeparatorColor
	c.TableViewGroupedCellTitleLabelColor = c.TableViewCellTitleLabelColor
	c.TableViewGroupedCellDetailLabelColor = c.TableViewCellDetailLabelColor
	c.TableViewGroupedCellBackgroundColor = c.WhiteColor
	c.TableViewGroupedCellSelectedBackgroundColor = c.TableViewCellSelectedBackgroundColor
	c.TableViewGroupedCellWarningBackgroundColor = c.YellowColor
	c.TableViewGroupedSectionHeaderFontSize = 12
	c.TableViewGroupedSectionFooterFontSize = 12
	c.TableViewGroupedSectionHeaderTextColor = c.GrayDarkenColor
	c.TableViewGroupedSectionFooterTextColor = c.GrayColor
	c.TableViewGroupedSectionHeaderDefaultHeight = 20
	c.TableViewGroupedSectionFooterDefaultHeight = 0

	// Inset Grouped TableView
	c.TableViewInsetGroupedCornerRadius = 10
	c.TableViewInsetGroupedHorizontalInset = 20
	c.TableViewInsetGroupedBackgroundColor = c.TableViewGroupedBackgroundColor
	c.TableViewInsetGroupedSeparatorColor = c.SeparatorColor

	// WindowLevel
	c.WindowLevelQMUIAlertView = 1999
	c.WindowLevelQMUIConsole = 1

	// QMUILog
	c.ShouldPrintDefaultLog = true
	c.ShouldPrintInfoLog = true
	c.ShouldPrintWarnLog = true
	c.ShouldPrintQMUIWarnLogToConsole = true

	// QMUIBadge
	c.BadgeBackgroundColor = c.RedColor
	c.BadgeTextColor = c.WhiteColor
	c.BadgeFontSize = 12
	c.BadgeContentEdgeInsets = NewEdgeInsets(3, 5, 3, 5)
	c.BadgeOffset = NewOffset(-9, 11)
	c.BadgeOffsetLandscape = NewOffset(-9, 6)
	c.UpdatesIndicatorColor = c.RedColor
	c.UpdatesIndicatorSize = fyne.NewSize(7, 7)
	c.UpdatesIndicatorOffset = NewOffset(4, 4)
	c.UpdatesIndicatorOffsetLandscape = NewOffset(3, 3)

	// Others
	c.AutomaticCustomNavigationBarTransitionStyle = true
	c.NeedsBackBarButtonItemTitle = false
	c.HidesBottomBarWhenPushedInitially = false
	c.PreventConcurrentNavigationControllerTransitions = true
	c.NavigationBarHiddenInitially = false

	// Alert Controller
	c.AlertContentMargin = NewEdgeInsets(0, 0, 0, 0)
	c.AlertContentMaximumWidth = 270
	c.AlertSeparatorColor = color.RGBA{R: 211, G: 211, B: 219, A: 255}
	c.AlertContentCornerRadius = 13
	c.AlertButtonHeight = 44
	c.AlertHeaderBackgroundColor = color.RGBA{R: 247, G: 247, B: 247, A: 255}
	c.AlertButtonBackgroundColor = color.RGBA{R: 247, G: 247, B: 247, A: 255}
	c.AlertButtonHighlightBackgroundColor = color.RGBA{R: 232, G: 232, B: 232, A: 255}
	c.AlertHeaderInsets = NewEdgeInsets(20, 16, 20, 16)
	c.AlertTitleMessageSpacing = 3
	c.AlertTextFieldFontSize = 14
	c.AlertTextFieldTextColor = c.BlackColor
	c.AlertTextFieldBorderColor = c.SeparatorColor
	c.AlertTextFieldTextInsets = NewEdgeInsets(4, 7, 4, 7)

	// Sheet
	c.SheetContentMargin = NewEdgeInsets(10, 10, 10, 10)
	c.SheetContentMaximumWidth = 414 - 20
	c.SheetSeparatorColor = color.RGBA{R: 211, G: 211, B: 219, A: 255}
	c.SheetCancelButtonMarginTop = 8
	c.SheetContentCornerRadius = 13
	c.SheetButtonHeight = 57
	c.SheetHeaderBackgroundColor = color.RGBA{R: 247, G: 247, B: 247, A: 255}
	c.SheetButtonBackgroundColor = color.RGBA{R: 247, G: 247, B: 247, A: 255}
	c.SheetButtonHighlightBackgroundColor = color.RGBA{R: 232, G: 232, B: 232, A: 255}
	c.SheetHeaderInsets = NewEdgeInsets(16, 16, 16, 16)
	c.SheetTitleMessageSpacing = 8
	c.SheetButtonColumnCount = 1

	// Toast
	c.ToastBackgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 191}
	c.ToastTextColor = c.WhiteColor
	c.ToastFontSize = 16
	c.ToastCornerRadius = 10
	c.ToastContentInsets = NewEdgeInsets(12, 16, 12, 16)
	c.ToastMarginFromScreen = 20
	c.ToastDefaultDuration = 2.0

	// EmptyView
	c.EmptyViewImageTintColor = nil
	c.EmptyViewLoadingTintColor = c.GrayColor
	c.EmptyViewTextLabelColor = color.RGBA{R: 147, G: 147, B: 147, A: 255}
	c.EmptyViewDetailTextLabelColor = color.RGBA{R: 133, G: 133, B: 133, A: 255}
	c.EmptyViewActionButtonColor = c.BlueColor
	c.EmptyViewTextFontSize = 15
	c.EmptyViewDetailTextFontSize = 14
	c.EmptyViewActionButtonFontSize = 15

	c.active = true
}

// ConfigurationTemplate defines the interface that all configuration templates should implement
type ConfigurationTemplate interface {
	// ApplyConfigurationTemplate applies the configuration settings
	ApplyConfigurationTemplate()
	// ShouldApplyTemplateAutomatically returns whether this template should be applied automatically on app launch
	ShouldApplyTemplateAutomatically() bool
}
