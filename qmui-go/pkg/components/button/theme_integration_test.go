package button

import (
	"image/color"
	"sync"
	"testing"

	"github.com/user/qmui-go/pkg/core"
	"github.com/user/qmui-go/pkg/theme"
)

// TestChangeToRoseThemeShouldChangeButtons verifies that switching to Pink Rose theme
// updates button colors appropriately
func TestChangeToRoseThemeShouldChangeButtons(t *testing.T) {
	// Reset for clean test state
	core.ResetConfigurationForTesting()
	theme.ResetForTesting()

	tm := theme.SharedThemeManager()

	// Get configuration (shared instance)
	cfg := core.SharedConfiguration()

	// Store initial blue color for comparison
	initialBlue := cfg.BlueColor
	t.Logf("Initial BlueColor: %+v", initialBlue)

	// Create a GhostButton with the current config color (should be default blue/cyan)
	ghostBtn := NewGhostButton("Test Ghost", cfg.BlueColor, func() {})
	initialBorderColor := ghostBtn.BorderColor
	t.Logf("Initial GhostButton BorderColor: %+v", initialBorderColor)

	// Verify initial border color matches config's blue
	if !colorsEqual(initialBorderColor, initialBlue) {
		t.Errorf("Initial button border should match config.BlueColor\n  got:  %+v\n  want: %+v",
			initialBorderColor, initialBlue)
	}

	// Get Pink Rose theme to know what color to expect
	roseTheme := tm.GetTheme(theme.ThemeIdentifierPinkRose)
	if roseTheme == nil {
		t.Fatal("Pink Rose theme should be registered")
	}
	expectedRoseColor := roseTheme.PrimaryColor
	t.Logf("Expected Rose PrimaryColor: %+v", expectedRoseColor)

	// Switch to Pink Rose theme
	t.Log("Switching to Pink Rose theme...")
	tm.SetCurrentTheme(theme.ThemeIdentifierPinkRose)

	// Verify configuration was updated
	updatedBlue := cfg.BlueColor
	t.Logf("After theme switch, cfg.BlueColor: %+v", updatedBlue)

	if !colorsEqual(updatedBlue, expectedRoseColor) {
		t.Errorf("After theme switch, config.BlueColor should be rose color\n  got:  %+v\n  want: %+v",
			updatedBlue, expectedRoseColor)
	}

	// KEY TEST: Create a NEW button after theme switch - it should use rose color
	newGhostBtn := NewGhostButton("New Ghost After Theme", cfg.BlueColor, func() {})
	newBorderColor := newGhostBtn.BorderColor
	t.Logf("New GhostButton (after theme switch) BorderColor: %+v", newBorderColor)

	if !colorsEqual(newBorderColor, expectedRoseColor) {
		t.Errorf("New button created after theme switch should have rose border\n  got:  %+v\n  want: %+v",
			newBorderColor, expectedRoseColor)
	}

	// ALSO TEST: Existing button with ApplyTheme called should update
	ghostBtn.ApplyTheme(roseTheme)
	appliedBorderColor := ghostBtn.BorderColor
	t.Logf("Original GhostButton after ApplyTheme: %+v", appliedBorderColor)

	if !colorsEqual(appliedBorderColor, expectedRoseColor) {
		t.Errorf("Button after ApplyTheme should have rose border\n  got:  %+v\n  want: %+v",
			appliedBorderColor, expectedRoseColor)
	}

	// Verify we're not still on blue
	defaultTheme := tm.GetTheme(theme.ThemeIdentifierDefault)
	defaultBlue := defaultTheme.PrimaryColor
	t.Logf("Default theme PrimaryColor (for comparison): %+v", defaultBlue)

	if colorsEqual(newBorderColor, defaultBlue) {
		t.Error("New button should NOT have default blue color after switching to rose theme")
	}
}

// TestThemeChangeUpdatesConfiguration verifies the theme manager updates config
func TestThemeChangeUpdatesConfiguration(t *testing.T) {
	core.ResetConfigurationForTesting()
	theme.ResetForTesting()

	tm := theme.SharedThemeManager()
	cfg := core.SharedConfiguration()

	// Get expected colors from themes
	defaultTheme := tm.GetTheme(theme.ThemeIdentifierDefault)
	roseTheme := tm.GetTheme(theme.ThemeIdentifierPinkRose)
	grassTheme := tm.GetTheme(theme.ThemeIdentifierGrass)

	t.Logf("Default primary: %+v", defaultTheme.PrimaryColor)
	t.Logf("Rose primary: %+v", roseTheme.PrimaryColor)
	t.Logf("Grass primary: %+v", grassTheme.PrimaryColor)

	// Initially should be default
	if !colorsEqual(cfg.BlueColor, defaultTheme.PrimaryColor) {
		t.Errorf("Initial config.BlueColor should be default primary")
	}

	// Switch to Rose
	tm.SetCurrentTheme(theme.ThemeIdentifierPinkRose)
	if !colorsEqual(cfg.BlueColor, roseTheme.PrimaryColor) {
		t.Errorf("After rose switch, config.BlueColor should be rose\n  got:  %+v\n  want: %+v",
			cfg.BlueColor, roseTheme.PrimaryColor)
	}

	// Switch to Grass
	tm.SetCurrentTheme(theme.ThemeIdentifierGrass)
	if !colorsEqual(cfg.BlueColor, grassTheme.PrimaryColor) {
		t.Errorf("After grass switch, config.BlueColor should be grass\n  got:  %+v\n  want: %+v",
			cfg.BlueColor, grassTheme.PrimaryColor)
	}

	// Switch back to Default
	tm.SetCurrentTheme(theme.ThemeIdentifierDefault)
	if !colorsEqual(cfg.BlueColor, defaultTheme.PrimaryColor) {
		t.Errorf("After default switch, config.BlueColor should be default\n  got:  %+v\n  want: %+v",
			cfg.BlueColor, defaultTheme.PrimaryColor)
	}
}

// TestGhostButtonApplyTheme verifies ApplyTheme updates border color
func TestGhostButtonApplyTheme(t *testing.T) {
	core.ResetConfigurationForTesting()
	theme.ResetForTesting()

	tm := theme.SharedThemeManager()

	// Create button with arbitrary color
	originalColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
	btn := NewGhostButton("Test", originalColor, func() {})

	if !colorsEqual(btn.BorderColor, originalColor) {
		t.Error("Button should start with original color")
	}

	// Apply rose theme
	roseTheme := tm.GetTheme(theme.ThemeIdentifierPinkRose)
	btn.ApplyTheme(roseTheme)

	if !colorsEqual(btn.BorderColor, roseTheme.PrimaryColor) {
		t.Errorf("After ApplyTheme, BorderColor should be rose primary\n  got:  %+v\n  want: %+v",
			btn.BorderColor, roseTheme.PrimaryColor)
	}

	// Also check TintColor (text color)
	if !colorsEqual(btn.TintColor, roseTheme.PrimaryColor) {
		t.Errorf("After ApplyTheme, TintColor should be rose primary\n  got:  %+v\n  want: %+v",
			btn.TintColor, roseTheme.PrimaryColor)
	}
}

// TestFillButtonHasVisibleBackground verifies FillButton renders with background
func TestFillButtonHasVisibleBackground(t *testing.T) {
	core.ResetConfigurationForTesting()
	theme.ResetForTesting()

	tm := theme.SharedThemeManager()
	cfg := core.SharedConfiguration()

	t.Logf("Config BlueColor: %+v", cfg.BlueColor)

	// Create a FillButton
	fillBtn := NewFillButton("Test Fill", cfg.BlueColor, func() {})

	t.Logf("FillButton.BackgroundColor: %+v", fillBtn.BackgroundColor)
	t.Logf("FillButton.TintColor (text): %+v", fillBtn.TintColor)
	t.Logf("FillButton.CornerRadius: %v", fillBtn.CornerRadius)

	// Verify the button has a background color set
	if fillBtn.BackgroundColor == nil {
		t.Fatal("FillButton.BackgroundColor should not be nil")
	}

	// Verify it matches config blue (default theme)
	if !colorsEqual(fillBtn.BackgroundColor, cfg.BlueColor) {
		t.Errorf("FillButton background should be config.BlueColor\n  got:  %+v\n  want: %+v",
			fillBtn.BackgroundColor, cfg.BlueColor)
	}

	// Switch to Rose theme
	tm.SetCurrentTheme(theme.ThemeIdentifierPinkRose)
	roseTheme := tm.GetTheme(theme.ThemeIdentifierPinkRose)

	t.Logf("After theme switch, Config BlueColor: %+v", cfg.BlueColor)

	// Create new button - should have rose background
	newFillBtn := NewFillButton("New Fill", cfg.BlueColor, func() {})
	t.Logf("New FillButton.BackgroundColor: %+v", newFillBtn.BackgroundColor)

	if !colorsEqual(newFillBtn.BackgroundColor, roseTheme.PrimaryColor) {
		t.Errorf("New FillButton background should be rose primary\n  got:  %+v\n  want: %+v",
			newFillBtn.BackgroundColor, roseTheme.PrimaryColor)
	}

	// The text should still be white
	if !colorsEqual(newFillBtn.TintColor, color.White) {
		t.Errorf("FillButton text should be white, got: %+v", newFillBtn.TintColor)
	}
}

// TestFillButtonApplyTheme verifies ApplyTheme updates FillButton colors
func TestFillButtonApplyTheme(t *testing.T) {
	core.ResetConfigurationForTesting()
	theme.ResetForTesting()

	tm := theme.SharedThemeManager()

	// Create button with gray color
	grayColor := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	fillBtn := NewFillButton("Test", grayColor, func() {})

	t.Logf("Initial BackgroundColor: %+v", fillBtn.BackgroundColor)

	// Apply rose theme
	roseTheme := tm.GetTheme(theme.ThemeIdentifierPinkRose)
	fillBtn.ApplyTheme(roseTheme)

	t.Logf("After ApplyTheme BackgroundColor: %+v", fillBtn.BackgroundColor)
	t.Logf("Rose ButtonBackgroundColor: %+v", roseTheme.ButtonBackgroundColor)

	// FillButton.ApplyTheme should update BackgroundColor to ButtonBackgroundColor
	if !colorsEqual(fillBtn.BackgroundColor, roseTheme.ButtonBackgroundColor) {
		t.Errorf("After ApplyTheme, FillButton background should be theme ButtonBackgroundColor\n  got:  %+v\n  want: %+v",
			fillBtn.BackgroundColor, roseTheme.ButtonBackgroundColor)
	}
}

// colorsEqual compares two colors for equality
func colorsEqual(a, b color.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// Ensure we have the ResetForTesting function in theme package
var _ = sync.Once{}
