// Package theme provides QMUITheme - a theming system for the application
// Ported from Tencent's QMUI_iOS framework
package theme

import (
	"fmt"
	"image/color"
	"sync"

	"fyne.io/fyne/v2"

	"github.com/user/qmui-go/pkg/core"
)

// ThemeIdentifier uniquely identifies a theme
type ThemeIdentifier string

const (
	// ThemeIdentifierDefault is the default theme (blue - matches iOS QMUI default)
	ThemeIdentifierDefault ThemeIdentifier = "default"

	// === QMUI iOS Theme Colors (all 10) ===

	// ThemeIdentifierGrapefruit - coral red (239, 83, 98)
	ThemeIdentifierGrapefruit ThemeIdentifier = "grapefruit"
	// ThemeIdentifierBittersweet - orange (254, 109, 75)
	ThemeIdentifierBittersweet ThemeIdentifier = "bittersweet"
	// ThemeIdentifierSunflower - yellow (255, 207, 71)
	ThemeIdentifierSunflower ThemeIdentifier = "sunflower"
	// ThemeIdentifierGrass - light green (159, 214, 97)
	ThemeIdentifierGrass ThemeIdentifier = "grass"
	// ThemeIdentifierMint - teal/mint (63, 208, 173)
	ThemeIdentifierMint ThemeIdentifier = "mint"
	// ThemeIdentifierKlein - deep blue (6, 92, 208)
	ThemeIdentifierKlein ThemeIdentifier = "klein"
	// ThemeIdentifierBlueJeans - sky blue (90, 154, 239)
	ThemeIdentifierBlueJeans ThemeIdentifier = "bluejeans"
	// ThemeIdentifierLavender - purple (172, 143, 239)
	ThemeIdentifierLavender ThemeIdentifier = "lavender"
	// ThemeIdentifierPinkRose - pink (238, 133, 193)
	ThemeIdentifierPinkRose ThemeIdentifier = "pinkrose"
	// ThemeIdentifierDark - cyan blue with dark background (39, 192, 243)
	ThemeIdentifierDark ThemeIdentifier = "dark"
)

// Theme defines a complete visual theme
type Theme struct {
	Identifier ThemeIdentifier
	Name       string

	// Colors
	PrimaryColor      color.Color
	SecondaryColor    color.Color
	BackgroundColor   color.Color
	SurfaceColor      color.Color
	TextPrimaryColor  color.Color
	TextSecondaryColor color.Color
	AccentColor       color.Color
	ErrorColor        color.Color
	SuccessColor      color.Color
	WarningColor      color.Color

	// Component-specific colors
	ButtonBackgroundColor      color.Color
	ButtonTextColor            color.Color
	ButtonDisabledColor        color.Color
	InputBackgroundColor       color.Color
	InputBorderColor           color.Color
	InputTextColor             color.Color
	InputPlaceholderColor      color.Color
	NavBarBackgroundColor      color.Color
	NavBarTintColor            color.Color
	NavBarTitleColor           color.Color
	TabBarBackgroundColor      color.Color
	TabBarTintColor            color.Color
	TableCellBackgroundColor   color.Color
	TableCellSelectedColor     color.Color
	SeparatorColor             color.Color
	ShadowColor                color.Color

	// Effects
	IsDarkMode bool
}

// NewDefaultTheme creates the default light theme
func NewDefaultTheme() *Theme {
	return &Theme{
		Identifier:             ThemeIdentifierDefault,
		Name:                   "Default",
		PrimaryColor:           color.RGBA{R: 49, G: 189, B: 243, A: 255},
		SecondaryColor:         color.RGBA{R: 159, G: 214, B: 97, A: 255},
		BackgroundColor:        color.White,
		SurfaceColor:           color.RGBA{R: 249, G: 249, B: 249, A: 255},
		TextPrimaryColor:       color.Black,
		TextSecondaryColor:     color.RGBA{R: 128, G: 128, B: 128, A: 255},
		AccentColor:            color.RGBA{R: 49, G: 189, B: 243, A: 255},
		ErrorColor:             color.RGBA{R: 250, G: 58, B: 58, A: 255},
		SuccessColor:           color.RGBA{R: 159, G: 214, B: 97, A: 255},
		WarningColor:           color.RGBA{R: 255, G: 207, B: 71, A: 255},
		ButtonBackgroundColor:  color.RGBA{R: 49, G: 189, B: 243, A: 255},
		ButtonTextColor:        color.White,
		ButtonDisabledColor:    color.RGBA{R: 200, G: 200, B: 200, A: 255},
		InputBackgroundColor:   color.White,
		InputBorderColor:       color.RGBA{R: 222, G: 224, B: 226, A: 255},
		InputTextColor:         color.Black,
		InputPlaceholderColor:  color.RGBA{R: 196, G: 200, B: 208, A: 255},
		NavBarBackgroundColor:  color.White,
		NavBarTintColor:        color.RGBA{R: 49, G: 189, B: 243, A: 255},
		NavBarTitleColor:       color.Black,
		TabBarBackgroundColor:  color.RGBA{R: 249, G: 249, B: 249, A: 255},
		TabBarTintColor:        color.RGBA{R: 49, G: 189, B: 243, A: 255},
		TableCellBackgroundColor: color.White,
		TableCellSelectedColor: color.RGBA{R: 238, G: 239, B: 241, A: 255},
		SeparatorColor:         color.RGBA{R: 222, G: 224, B: 226, A: 255},
		ShadowColor:            color.RGBA{R: 0, G: 0, B: 0, A: 40},
		IsDarkMode:             false,
	}
}

// NewDarkTheme creates the Dark theme - cyan blue (39, 192, 243) with dark background
func NewDarkTheme() *Theme {
	primary := color.RGBA{R: 39, G: 192, B: 243, A: 255} // QMUI iOS Dark theme cyan
	return &Theme{
		Identifier:               ThemeIdentifierDark,
		Name:                     "Dark",
		PrimaryColor:             primary,
		SecondaryColor:           color.RGBA{R: 48, G: 209, B: 88, A: 255},
		BackgroundColor:          color.RGBA{R: 0, G: 0, B: 0, A: 255},
		SurfaceColor:             color.RGBA{R: 28, G: 28, B: 30, A: 255},
		TextPrimaryColor:         color.RGBA{R: 218, G: 220, B: 224, A: 255}, // UIColorDarkGray1
		TextSecondaryColor:       color.RGBA{R: 178, G: 180, B: 184, A: 255}, // UIColorDarkGray3
		AccentColor:              primary,
		ErrorColor:               color.RGBA{R: 255, G: 69, B: 58, A: 255},
		SuccessColor:             color.RGBA{R: 48, G: 209, B: 88, A: 255},
		WarningColor:             color.RGBA{R: 255, G: 214, B: 10, A: 255},
		ButtonBackgroundColor:    primary,
		ButtonTextColor:          color.White,
		ButtonDisabledColor:      color.RGBA{R: 72, G: 72, B: 74, A: 255},
		InputBackgroundColor:     color.RGBA{R: 28, G: 28, B: 30, A: 255},
		InputBorderColor:         color.RGBA{R: 56, G: 56, B: 58, A: 255},
		InputTextColor:           color.White,
		InputPlaceholderColor:    color.RGBA{R: 78, G: 80, B: 84, A: 255}, // UIColorDarkGray8
		NavBarBackgroundColor:    color.RGBA{R: 28, G: 28, B: 30, A: 255},
		NavBarTintColor:          primary,
		NavBarTitleColor:         color.White,
		TabBarBackgroundColor:    color.RGBA{R: 28, G: 28, B: 30, A: 255},
		TabBarTintColor:          primary,
		TableCellBackgroundColor: color.RGBA{R: 28, G: 28, B: 30, A: 255},
		TableCellSelectedColor:   color.RGBA{R: 48, G: 49, B: 51, A: 255},
		SeparatorColor:           color.RGBA{R: 46, G: 50, B: 54, A: 255},
		ShadowColor:              color.RGBA{R: 0, G: 0, B: 0, A: 80},
		IsDarkMode:               true,
	}
}

// newLightTheme creates a light theme with the given primary color
func newLightTheme(identifier ThemeIdentifier, name string, primary color.RGBA) *Theme {
	return &Theme{
		Identifier:               identifier,
		Name:                     name,
		PrimaryColor:             primary,
		SecondaryColor:           color.RGBA{R: 159, G: 214, B: 97, A: 255}, // Grass green
		BackgroundColor:          color.White,
		SurfaceColor:             color.RGBA{R: 249, G: 249, B: 249, A: 255},
		TextPrimaryColor:         color.Black,
		TextSecondaryColor:       color.RGBA{R: 128, G: 128, B: 128, A: 255},
		AccentColor:              primary,
		ErrorColor:               color.RGBA{R: 250, G: 58, B: 58, A: 255},
		SuccessColor:             color.RGBA{R: 159, G: 214, B: 97, A: 255},
		WarningColor:             color.RGBA{R: 255, G: 207, B: 71, A: 255},
		ButtonBackgroundColor:    primary,
		ButtonTextColor:          color.White,
		ButtonDisabledColor:      color.RGBA{R: 200, G: 200, B: 200, A: 255},
		InputBackgroundColor:     color.White,
		InputBorderColor:         color.RGBA{R: 222, G: 224, B: 226, A: 255},
		InputTextColor:           color.Black,
		InputPlaceholderColor:    color.RGBA{R: 196, G: 200, B: 208, A: 255},
		NavBarBackgroundColor:    color.White,
		NavBarTintColor:          primary,
		NavBarTitleColor:         color.Black,
		TabBarBackgroundColor:    color.RGBA{R: 249, G: 249, B: 249, A: 255},
		TabBarTintColor:          primary,
		TableCellBackgroundColor: color.White,
		TableCellSelectedColor:   color.RGBA{R: 238, G: 239, B: 241, A: 255},
		SeparatorColor:           color.RGBA{R: 222, G: 224, B: 226, A: 255},
		ShadowColor:              color.RGBA{R: 0, G: 0, B: 0, A: 40},
		IsDarkMode:               false,
	}
}

// NewGrapefruitTheme creates the Grapefruit theme - coral red (239, 83, 98)
func NewGrapefruitTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierGrapefruit,
		"Grapefruit",
		color.RGBA{R: 239, G: 83, B: 98, A: 255},
	)
}

// NewBittersweetTheme creates the Bittersweet theme - orange (254, 109, 75)
func NewBittersweetTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierBittersweet,
		"Bittersweet",
		color.RGBA{R: 254, G: 109, B: 75, A: 255},
	)
}

// NewSunflowerTheme creates the Sunflower theme - yellow (255, 207, 71)
func NewSunflowerTheme() *Theme {
	theme := newLightTheme(
		ThemeIdentifierSunflower,
		"Sunflower",
		color.RGBA{R: 255, G: 207, B: 71, A: 255},
	)
	// Yellow needs dark text for contrast
	theme.ButtonTextColor = color.Black
	return theme
}

// NewGrassTheme creates the Grass theme - light green (159, 214, 97)
func NewGrassTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierGrass,
		"Grass",
		color.RGBA{R: 159, G: 214, B: 97, A: 255},
	)
}

// NewMintTheme creates the Mint theme - teal/mint (63, 208, 173)
func NewMintTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierMint,
		"Mint",
		color.RGBA{R: 63, G: 208, B: 173, A: 255},
	)
}

// NewKleinTheme creates the Klein theme - deep blue (6, 92, 208)
func NewKleinTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierKlein,
		"Klein",
		color.RGBA{R: 6, G: 92, B: 208, A: 255},
	)
}

// NewBlueJeansTheme creates the Blue Jeans theme - sky blue (90, 154, 239)
func NewBlueJeansTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierBlueJeans,
		"Blue Jeans",
		color.RGBA{R: 90, G: 154, B: 239, A: 255},
	)
}

// NewLavenderTheme creates the Lavender theme - purple (172, 143, 239)
func NewLavenderTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierLavender,
		"Lavender",
		color.RGBA{R: 172, G: 143, B: 239, A: 255},
	)
}

// NewPinkRoseTheme creates the Pink Rose theme - pink (238, 133, 193)
func NewPinkRoseTheme() *Theme {
	return newLightTheme(
		ThemeIdentifierPinkRose,
		"Pink Rose",
		color.RGBA{R: 238, G: 133, B: 193, A: 255},
	)
}

// ThemeManager manages themes for the application
type ThemeManager struct {
	mu             sync.RWMutex
	currentTheme   *Theme
	themes         map[ThemeIdentifier]*Theme
	listeners      []func(theme *Theme)
}

var (
	sharedManager *ThemeManager
	managerOnce   sync.Once
)

// SharedThemeManager returns the shared theme manager
func SharedThemeManager() *ThemeManager {
	managerOnce.Do(func() {
		sharedManager = &ThemeManager{
			themes:    make(map[ThemeIdentifier]*Theme),
			listeners: make([]func(theme *Theme), 0),
		}
		// Register all 10 QMUI iOS themes + default
		sharedManager.RegisterTheme(NewDefaultTheme())
		sharedManager.RegisterTheme(NewGrapefruitTheme())
		sharedManager.RegisterTheme(NewBittersweetTheme())
		sharedManager.RegisterTheme(NewSunflowerTheme())
		sharedManager.RegisterTheme(NewGrassTheme())
		sharedManager.RegisterTheme(NewMintTheme())
		sharedManager.RegisterTheme(NewKleinTheme())
		sharedManager.RegisterTheme(NewBlueJeansTheme())
		sharedManager.RegisterTheme(NewLavenderTheme())
		sharedManager.RegisterTheme(NewPinkRoseTheme())
		sharedManager.RegisterTheme(NewDarkTheme())
		sharedManager.currentTheme = sharedManager.themes[ThemeIdentifierDefault]

		// Apply default theme to configuration
		sharedManager.applyThemeToConfiguration(sharedManager.currentTheme)
	})
	return sharedManager
}

// ResetForTesting resets the theme manager for testing purposes
// This should only be used in tests
func ResetForTesting() {
	managerOnce = sync.Once{}
	sharedManager = nil
}

// AllThemes returns all registered themes in a consistent order matching iOS QMUI
func (tm *ThemeManager) AllThemes() []*Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// Return in QMUI iOS order
	order := []ThemeIdentifier{
		ThemeIdentifierDefault,
		ThemeIdentifierGrapefruit,
		ThemeIdentifierBittersweet,
		ThemeIdentifierSunflower,
		ThemeIdentifierGrass,
		ThemeIdentifierMint,
		ThemeIdentifierKlein,
		ThemeIdentifierBlueJeans,
		ThemeIdentifierLavender,
		ThemeIdentifierPinkRose,
		ThemeIdentifierDark,
	}

	themes := make([]*Theme, 0, len(order))
	for _, id := range order {
		if t, ok := tm.themes[id]; ok {
			themes = append(themes, t)
		}
	}
	return themes
}

// AllThemeIdentifiers returns all theme identifiers in QMUI iOS order
func (tm *ThemeManager) AllThemeIdentifiers() []ThemeIdentifier {
	return []ThemeIdentifier{
		ThemeIdentifierDefault,
		ThemeIdentifierGrapefruit,
		ThemeIdentifierBittersweet,
		ThemeIdentifierSunflower,
		ThemeIdentifierGrass,
		ThemeIdentifierMint,
		ThemeIdentifierKlein,
		ThemeIdentifierBlueJeans,
		ThemeIdentifierLavender,
		ThemeIdentifierPinkRose,
		ThemeIdentifierDark,
	}
}

// CycleTheme switches to the next theme in the list
func (tm *ThemeManager) CycleTheme() *Theme {
	tm.mu.RLock()
	current := tm.currentTheme
	tm.mu.RUnlock()

	ids := tm.AllThemeIdentifiers()
	currentIdx := 0
	for i, id := range ids {
		if current != nil && id == current.Identifier {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(ids)
	tm.SetCurrentTheme(ids[nextIdx])
	return tm.CurrentTheme()
}

// RegisterTheme registers a theme
func (tm *ThemeManager) RegisterTheme(theme *Theme) {
	tm.mu.Lock()
	tm.themes[theme.Identifier] = theme
	tm.mu.Unlock()
}

// UnregisterTheme removes a theme
func (tm *ThemeManager) UnregisterTheme(identifier ThemeIdentifier) {
	tm.mu.Lock()
	delete(tm.themes, identifier)
	tm.mu.Unlock()
}

// GetTheme returns a theme by identifier
func (tm *ThemeManager) GetTheme(identifier ThemeIdentifier) *Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.themes[identifier]
}

// CurrentTheme returns the current theme
func (tm *ThemeManager) CurrentTheme() *Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.currentTheme
}

// SetCurrentTheme changes the current theme
func (tm *ThemeManager) SetCurrentTheme(identifier ThemeIdentifier) {
	tm.mu.Lock()
	theme, exists := tm.themes[identifier]
	if !exists {
		tm.mu.Unlock()
		return
	}
	tm.currentTheme = theme
	listeners := tm.listeners
	tm.mu.Unlock()

	// Apply theme to configuration
	tm.applyThemeToConfiguration(theme)

	// Notify listeners
	for _, listener := range listeners {
		listener(theme)
	}
}

// AddThemeChangeListener adds a listener for theme changes
func (tm *ThemeManager) AddThemeChangeListener(listener func(theme *Theme)) {
	tm.mu.Lock()
	tm.listeners = append(tm.listeners, listener)
	tm.mu.Unlock()
}

// applyThemeToConfiguration applies theme colors to the global configuration
func (tm *ThemeManager) applyThemeToConfiguration(theme *Theme) {
	config := core.SharedConfiguration()

	config.BlueColor = theme.PrimaryColor
	config.GreenColor = theme.SuccessColor
	config.RedColor = theme.ErrorColor
	config.YellowColor = theme.WarningColor
	config.BackgroundColor = theme.BackgroundColor
	config.SeparatorColor = theme.SeparatorColor
	config.PlaceholderColor = theme.InputPlaceholderColor

	config.NavBarBackgroundColor = theme.NavBarBackgroundColor
	config.NavBarTintColor = theme.NavBarTintColor
	config.NavBarTitleColor = theme.NavBarTitleColor

	config.TabBarBackgroundColor = theme.TabBarBackgroundColor
	config.TabBarItemTitleColorSelected = theme.TabBarTintColor
	config.TabBarItemImageColorSelected = theme.TabBarTintColor

	config.TableViewCellBackgroundColor = theme.TableCellBackgroundColor
	config.TableViewCellSelectedBackgroundColor = theme.TableCellSelectedColor
	config.TableViewCellTitleLabelColor = theme.TextPrimaryColor
	config.TableViewCellDetailLabelColor = theme.TextSecondaryColor

	config.ButtonTintColor = theme.ButtonBackgroundColor
}

// ThemeColor creates a color that adapts to theme changes
type ThemeColor struct {
	lightColor color.Color
	darkColor  color.Color
}

// NewThemeColor creates a theme-aware color
func NewThemeColor(light, dark color.Color) *ThemeColor {
	return &ThemeColor{
		lightColor: light,
		darkColor:  dark,
	}
}

// Color returns the appropriate color for the current theme
func (tc *ThemeColor) Color() color.Color {
	theme := SharedThemeManager().CurrentTheme()
	if theme != nil && theme.IsDarkMode {
		return tc.darkColor
	}
	return tc.lightColor
}

// RGBA implements color.Color
func (tc *ThemeColor) RGBA() (r, g, b, a uint32) {
	return tc.Color().RGBA()
}

// Themeable is implemented by widgets that can respond to theme changes
type Themeable interface {
	ApplyTheme(theme *Theme)
}

// ApplyThemeToWindow applies the current theme to all themeable widgets in a window
func ApplyThemeToWindow(window fyne.Window) {
	theme := SharedThemeManager().CurrentTheme()
	if theme == nil {
		return
	}

	// Walk through all objects and apply theme
	applyThemeToObject(window.Content(), theme)
}

func applyThemeToObject(obj fyne.CanvasObject, theme *Theme) {
	if obj == nil {
		return
	}

	// Apply theme if object implements Themeable
	if themeable, ok := obj.(Themeable); ok {
		fmt.Printf("Applying theme to: %T\n", obj)
		themeable.ApplyTheme(theme)
	}

	// Recursively apply to children if it's a container
	if cont, ok := obj.(*fyne.Container); ok {
		for _, child := range cont.Objects {
			applyThemeToObject(child, theme)
		}
	}

	// Also check for widgets that have child objects via renderer
	if w, ok := obj.(fyne.Widget); ok {
		// Use test package approach - get objects from renderer
		if renderer := w.CreateRenderer(); renderer != nil {
			for _, child := range renderer.Objects() {
				applyThemeToObject(child, theme)
			}
		}
	}
}

// ToggleDarkMode toggles between light and dark modes
func ToggleDarkMode() {
	tm := SharedThemeManager()
	current := tm.CurrentTheme()
	if current != nil && current.IsDarkMode {
		tm.SetCurrentTheme(ThemeIdentifierDefault)
	} else {
		tm.SetCurrentTheme(ThemeIdentifierDark)
	}
}

// IsDarkMode returns whether dark mode is currently active
func IsDarkMode() bool {
	theme := SharedThemeManager().CurrentTheme()
	return theme != nil && theme.IsDarkMode
}
