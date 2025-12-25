# QMUI Go - A UI Component Library for Fyne

A comprehensive Go+Fyne port of [QMUI iOS](https://github.com/Tencent/QMUI_iOS), providing polished, themeable UI components for desktop applications.

> **Note:** This is an unofficial community port, not affiliated with Tencent. Would y'all want it back were I to grant my portions of the copyright?

## Features

- **45+ Components** - Buttons, dialogs, progress indicators, navigation, and more
- **Hot-Switchable Themes** - 11 built-in themes (Default, Dark, Grass Green, Pink Rose, etc.)
- **iOS-Quality Polish** - Smooth animations, proper states, consistent styling
- **100% Test Coverage** - 150+ visual regression tests

## Installation

```go
go get github.com/paul-hammant/qmui_fyne
```

## Running the Demo

The demo app provides an interactive guided tour of all 45+ components:

```bash
cd qmui_fyne
go run ./cmd/demo
```

The demo includes:
- **Themes tab** - Hot-switch between all 11 themes
- **Components tab** - Badge, Label, Marquee, EmptyView
- **Controls tab** - Buttons, Checkbox, Switch, Segmented, TextField
- **Progress tab** - Pie, Ring, and Linear progress indicators
- **Layout tab** - Grid, FloatLayout, Tile, ImagePreview
- **Dialogs tab** - Toast, Tips, Alert, Dialog, Popup, Modal
- **Navigation tab** - NavigationBar, TabBar, Table
- **Debug tab** - Console and logging utilities

---

## Quick Start

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "github.com/paul-hammant/qmui_fyne/button"
    "github.com/paul-hammant/qmui_fyne/tips"
    "github.com/paul-hammant/qmui_fyne/theme"
)

func main() {
    a := app.New()
    w := a.NewWindow("QMUI Demo")

    // Hot-switch to dark theme
    theme.SharedThemeManager().SetCurrentTheme(theme.ThemeIdentifierDark)

    // Create a styled button
    btn := button.NewFillButton("Click Me", theme.SharedThemeManager().CurrentTheme().BlueColor, func() {
        tips.ShowSuccess(w, "Tapped!")
    })

    w.SetContent(btn)
    w.ShowAndRun()
}
```

---

## Component Reference

### Widgets

| Component | Package | Fyne Base | Alternative To | Description |
|-----------|---------|-----------|----------------|-------------|
| **Buttons** |
| `Button` | `button`* | `widget.BaseWidget` | `widget.Button` | Customizable button with image positioning, states, insets |
| `FillButton` | `button` | `Button` | `widget.Button` + custom bg | Solid filled button with rounded corners |
| `GhostButton` | `button` | `Button` | `widget.Button` + border | Outlined/bordered button |
| `NavigationButton` | `button` | `Button` | Plain text button | Navigation bar styled button |
| `ToolbarButton` | `button` | `Button` | `widget.ToolbarAction` | Toolbar styled button |
| **Labels** |
| `Label` | `label` | `widget.BaseWidget` | `widget.Label` | Label with edge insets, long-press copy |
| `Marquee` | `marquee` | `widget.BaseWidget` | Custom animation | Auto-scrolling text for overflow |
| `Badge` | `badge` | `widget.BaseWidget` | — | Notification badge (e.g., "99+") |
| **Text Input** |
| `TextField` | `textfield` | `widget.Entry` | `widget.Entry` | Single-line with max length, placeholder styling |
| `TextView` | `textview` | `widget.Entry` | `widget.Entry{MultiLine}` | Multi-line with placeholder, auto-height |
| `SearchBar` | `search` | `widget.BaseWidget` | `widget.Entry` + icons | Search input with suggestions, cancel button |
| **Selection** |
| `Checkbox` | `checkbox` | `widget.BaseWidget` | `widget.Check` | Styled checkbox with custom colors |
| `SegmentedControl` | `segmented` | `widget.BaseWidget` | — | iOS-style segmented control (pill, underline) |
| **Progress** |
| `PieProgress` | `progress` | `widget.BaseWidget` | — | Animated pie/donut chart |
| `RingProgress` | `progress` | `widget.BaseWidget` | `widget.ProgressBar` | Circular progress with percentage |
| `ProgressBar` | `progress` | `widget.BaseWidget` | `widget.ProgressBar` | Styled linear progress bar |
| **Dialogs** |
| `Alert` | `alert` | `widget.BaseWidget` | `dialog.Confirm` | Alert dialog + action sheet styles |
| `Dialog` | `dialog` | `widget.BaseWidget` | `dialog.Custom` | Custom content dialog with actions |
| `Modal` | `modal` | `widget.BaseWidget` | `widget.PopUp` | Animated modal (fade, slide, bounce, zoom) |
| `ActionSheet` | `moreop` | `widget.BaseWidget` | — | Grid-style bottom sheet (share menu) |
| **Toast/HUD** |
| `Toast` | `toast` | — | `dialog.ShowInformation` | Non-blocking toast messages |
| `HUD` | `tips` | — | — | Loading/Success/Error/Info with icons |
| **Popup** |
| `PopupMenu` | `popup` | `widget.BaseWidget` | `widget.PopUpMenu` | Context menu with items |
| `PopupContainer` | `popup` | `widget.BaseWidget` | — | Arrow-pointed popup container |
| **Empty States** |
| `EmptyState` | `empty` | `widget.BaseWidget` | — | Loading/error/no-data placeholder |
| **Image** |
| `ZoomImage` | `zoomimage` | `widget.BaseWidget` | `canvas.Image` | Pan, zoom, double-tap, gestures |
| `ImagePreview` | `imagepreview` | `widget.BaseWidget` | — | Full-screen image viewer with swipe |
| `ImagePicker` | `imagepicker` | `widget.BaseWidget` | `dialog.OpenFile` | Multi-select photo picker |
| `AlbumView` | `album` | `widget.BaseWidget` | — | Album browser with photo grid |
| **Special** |
| `EmojiPicker` | `emotion` | `widget.BaseWidget` | — | Emoji grid with category tabs |
| `Tile` | `tile` | `widget.BaseWidget` | — | Tappable tile with icon and label |

### Containers & Layouts

| Component | Package | Fyne Base | Alternative To | Description |
|-----------|---------|-----------|----------------|-------------|
| `Grid` | `grid` | `widget.BaseWidget` | `container.NewGridWrap` | Fixed-column grid with separators |
| `FlowLayout` | `floatlayout` | `widget.BaseWidget` | — | Tag cloud / flow layout |
| `TagCloud` | `floatlayout` | `FlowLayout` | — | Convenience wrapper for tags |
| `Table` | `table` | `widget.BaseWidget` | `widget.List` | Grouped/inset-grouped table view |
| `TableCell` | `table` | `widget.BaseWidget` | — | Table cell with title, detail, accessory |
| `PagingLayout` | `collection` | `widget.BaseWidget` | `container.DocTabs` | Paging with scale/coverflow effects |
| `Layouter` | `layouter` | `widget.BaseWidget` | `container.VBox/HBox` | Linear layout with spacing |

### Navigation

| Component | Package | Fyne Base | Alternative To | Description |
|-----------|---------|-----------|----------------|-------------|
| `NavigationBar` | `navigation` | `widget.BaseWidget` | — | Top navigation bar with shadow, tint |
| `TitleView` | `navigation` | `widget.BaseWidget` | — | Title + subtitle + loading indicator |
| `TabBar` | `navigation` | `widget.BaseWidget` | `container.AppTabs` | Bottom tab bar with icons, badges |

### Theme & Configuration

| Component | Package | Description |
|-----------|---------|-------------|
| `ThemeManager` | `theme` | Hot-switchable theme management |
| `Theme` | `theme` | Theme definition with colors |
| `Configuration` | `core` | Global UI configuration |
| `Animation` | `animation` | 25+ easing functions for animations |

\* Package names are relative to `github.com/paul-hammant/qmui_fyne/`. For example, import the button package as `import "github.com/paul-hammant/qmui_fyne/button"`.

---

## Built-in Themes

| Theme | Identifier |
|-------|------------|
| Default (Light) | `ThemeIdentifierDefault` |
| Dark | `ThemeIdentifierDark` |
| Grass Green | `ThemeIdentifierGrassGreen` |
| Pink Rose | `ThemeIdentifierPinkRose` |
| Ocean Blue | `ThemeIdentifierOceanBlue` |
| Sunset Orange | `ThemeIdentifierSunsetOrange` |
| Purple Dream | `ThemeIdentifierPurpleDream` |
| Mint Fresh | `ThemeIdentifierMintFresh` |
| Coral Red | `ThemeIdentifierCoralRed` |
| Slate Gray | `ThemeIdentifierSlateGray` |
| Golden Sand | `ThemeIdentifierGoldenSand` |

### Theme Switching

```go
tm := theme.SharedThemeManager()

// Switch theme
tm.SetCurrentTheme(theme.ThemeIdentifierDark)

// Listen for theme changes
tm.AddThemeChangeListener(func(t *theme.Theme) {
    // Update UI components
})
```

---

## Package Structure

```
qmui_fyne/
├── core/           # Configuration, helpers
├── theme/          # Hot-switchable themes
├── animation/      # Easing functions
├── alert/          # Alert, Action
├── album/          # AlbumView
├── badge/          # Badge
├── button/         # Button, FillButton, GhostButton, etc.
├── checkbox/       # Checkbox
├── collection/     # PagingLayout
├── dialog/         # Dialog, DialogAction
├── emotion/        # EmojiPicker
├── empty/          # EmptyState
├── floatlayout/    # FlowLayout, TagCloud, Tag
├── grid/           # Grid, GridItem
├── imagepicker/    # ImagePicker
├── imagepreview/   # ImagePreview
├── label/          # Label
├── layouter/       # Layouter
├── marquee/        # Marquee
├── modal/          # Modal
├── moreop/         # ActionSheet, Item
├── navigation/     # NavigationBar, TitleView, TabBar
├── popup/          # PopupMenu, PopupContainer
├── progress/       # PieProgress, RingProgress, ProgressBar
├── search/         # SearchBar
├── segmented/      # SegmentedControl
├── table/          # Table, TableCell, TableSection
├── textfield/      # TextField
├── textview/       # TextView
├── tile/           # Tile
├── tips/           # HUD
├── toast/          # Toast
└── zoomimage/      # ZoomImage
```

---

## Testing

Run all tests:
```bash
go test ./... -v
```

The test suite includes 150+ visual regression tests covering all components.

---

## License

MIT License

Portions Copyright (C) 2016-2021 THL A29 Limited, a Tencent company. All rights reserved.

Portions Copyright (C) 2024-2025 Paul Hammant. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

## Credits

- Original iOS library: [QMUI/QMUI_iOS](https://github.com/Tencent/QMUI_iOS) by Tencent
- Fyne toolkit: [fyne-io/fyne](https://github.com/fyne-io/fyne)
