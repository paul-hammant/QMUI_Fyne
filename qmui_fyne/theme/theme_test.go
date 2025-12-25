package theme

import (
	"sync"
	"testing"
)

func TestThemeManager_HotSwitch(t *testing.T) {
	// Reset the shared manager for testing
	managerOnce = sync.Once{}
	sharedManager = nil

	tm := SharedThemeManager()

	// Verify default theme is set
	current := tm.CurrentTheme()
	if current == nil {
		t.Fatal("CurrentTheme() returned nil")
	}
	if current.Identifier != ThemeIdentifierDefault {
		t.Errorf("Expected default theme, got %s", current.Identifier)
	}
	t.Logf("Initial theme: %s (%s)", current.Name, current.Identifier)

	// Track listener calls
	var listenerCalls int
	var lastTheme *Theme

	tm.AddThemeChangeListener(func(theme *Theme) {
		listenerCalls++
		lastTheme = theme
		t.Logf("Listener called with theme: %s", theme.Name)
	})

	// Switch to Grapefruit
	t.Log("Switching to Grapefruit theme...")
	tm.SetCurrentTheme(ThemeIdentifierGrapefruit)

	if listenerCalls != 1 {
		t.Errorf("Expected 1 listener call, got %d", listenerCalls)
	}
	if lastTheme == nil || lastTheme.Identifier != ThemeIdentifierGrapefruit {
		t.Errorf("Listener did not receive Grapefruit theme")
	}
	if tm.CurrentTheme().Identifier != ThemeIdentifierGrapefruit {
		t.Errorf("CurrentTheme() should be Grapefruit, got %s", tm.CurrentTheme().Identifier)
	}

	// Switch to Grass
	t.Log("Switching to Grass theme...")
	tm.SetCurrentTheme(ThemeIdentifierGrass)

	if listenerCalls != 2 {
		t.Errorf("Expected 2 listener calls, got %d", listenerCalls)
	}
	if lastTheme == nil || lastTheme.Identifier != ThemeIdentifierGrass {
		t.Errorf("Listener did not receive Grass theme")
	}

	// Switch to Dark
	t.Log("Switching to Dark theme...")
	tm.SetCurrentTheme(ThemeIdentifierDark)

	if listenerCalls != 3 {
		t.Errorf("Expected 3 listener calls, got %d", listenerCalls)
	}
	if lastTheme == nil || lastTheme.Identifier != ThemeIdentifierDark {
		t.Errorf("Listener did not receive Dark theme")
	}
	if !tm.CurrentTheme().IsDarkMode {
		t.Error("Dark theme should have IsDarkMode=true")
	}

	// Switch back to Default
	t.Log("Switching back to Default theme...")
	tm.SetCurrentTheme(ThemeIdentifierDefault)

	if listenerCalls != 4 {
		t.Errorf("Expected 4 listener calls, got %d", listenerCalls)
	}

	t.Logf("Total listener calls: %d", listenerCalls)
}

func TestThemeManager_MultipleListeners(t *testing.T) {
	// Reset the shared manager for testing
	managerOnce = sync.Once{}
	sharedManager = nil

	tm := SharedThemeManager()

	var calls1, calls2, calls3 int

	tm.AddThemeChangeListener(func(theme *Theme) {
		calls1++
	})
	tm.AddThemeChangeListener(func(theme *Theme) {
		calls2++
	})
	tm.AddThemeChangeListener(func(theme *Theme) {
		calls3++
	})

	tm.SetCurrentTheme(ThemeIdentifierGrapefruit)

	if calls1 != 1 || calls2 != 1 || calls3 != 1 {
		t.Errorf("All listeners should be called once: got %d, %d, %d", calls1, calls2, calls3)
	}
}

func TestThemeManager_AllThemesRegistered(t *testing.T) {
	// Reset the shared manager for testing
	managerOnce = sync.Once{}
	sharedManager = nil

	tm := SharedThemeManager()

	expectedThemes := []ThemeIdentifier{
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

	for _, id := range expectedThemes {
		theme := tm.GetTheme(id)
		if theme == nil {
			t.Errorf("Theme %s should be registered", id)
		} else {
			t.Logf("Theme %s: %s (primary: %v)", id, theme.Name, theme.PrimaryColor)
		}
	}

	allThemes := tm.AllThemes()
	if len(allThemes) != len(expectedThemes) {
		t.Errorf("Expected %d themes, got %d", len(expectedThemes), len(allThemes))
	}
}

func TestThemeManager_InvalidTheme(t *testing.T) {
	// Reset the shared manager for testing
	managerOnce = sync.Once{}
	sharedManager = nil

	tm := SharedThemeManager()

	var called bool
	tm.AddThemeChangeListener(func(theme *Theme) {
		called = true
	})

	originalTheme := tm.CurrentTheme()

	// Try to switch to non-existent theme
	tm.SetCurrentTheme("nonexistent")

	// Should not call listener
	if called {
		t.Error("Listener should not be called for invalid theme")
	}

	// Should keep original theme
	if tm.CurrentTheme().Identifier != originalTheme.Identifier {
		t.Error("Theme should not change for invalid identifier")
	}
}
