# QMUI-Go Porting TODO List

This document tracks the porting progress of QMUI iOS components to Go with the Fyne UI framework.

## Component Porting Status

| iOS Component | Go+Fyne Component | Status | Notes |
|---|---|---|---|
| `QMUIButton` | `button` | ✅ Verified | Implemented and included in the demo. |
| `QMUILabel` | `label` | ✅ Verified | Implemented and included in the demo. |
| `QMUIMarqueeLabel` | `marquee` | ✅ Verified | Implemented and included in the demo. |
| `QMUIBadge` | `badge` | ✅ Verified | Implemented and included in the demo. |
| `QMUIAlertController` | `alert` | ✅ Verified | Implemented and included in the demo. |
| `QMUIToastView` | `toast` | ✅ Verified | Implemented and included in the demo. |
| `QMUIPopupMenuView` | `popup` | ✅ Verified | Implemented and included in the demo. |
| `QMUIPieProgressView` | `progress` | ✅ Verified | Implemented and included in the demo. |
| `QMUIEmptyView` | `empty` | ✅ Verified | Implemented and included in the demo. |
| `QMUIGridView` | `grid` | ✅ Verified | Implemented and included in the demo. |
| `QMUIFloatLayoutView` | `floatlayout` | ✅ Verified | Implemented and included in the demo. |
| `QMUITableView` | `table` | ✅ Verified | Implemented and included in the demo. |
| `QMUITextField` | `textfield` | ✅ Verified | Implemented and included in the demo. |
| `QMUITextView` | `textview` | ✅ Verified | Implemented and included in the demo. |
| `QMUICheckbox` | `checkbox` | ✅ Verified | A faithful port. The Go version uses programmatic drawing instead of images, which is an idiomatic improvement for Fyne. |
| `(Not Applicable)` | `radiobutton` | ❌ Not Implemented | Component does not exist in the original QMUI_iOS. Functionality is likely achieved via other controls like `SegmentedControl`. |
| `UISwitch+QMUI` | `switch` | ✅ Verified | Implemented as a custom Fyne widget to support `qmui_offTintColor`. An idiomatic port. |
| `QMUISegmentedControl` | `segmented` | ✅ Verified | Implemented and included in the demo. |
| `QMUISearchBar` | `search` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIDialogViewController` | `dialog` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIModalPresentationViewController` | `modal` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIMoreOperationController` | `moreop` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUITips` | `tips` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUICollectionViewPagingLayout` | `collection` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUINavigation` | `navigation` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIZoomImageView` | `zoomimage` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIImagePreviewViewController` | `imagepreview` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIImagePickerViewController` | `imagepicker` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIAlbumViewController` | `album` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUIEmotionView` | `emotion` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUILayouter` | `layouter` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |
| `QMUItile` | `tile` | 🟡 Implemented, Needs Demo | Code exists, but not shown in the demo. |

---

## Roadmap

1.  **Verify Existing Components:** Add all `Implemented, Needs Demo` components to the demo application to visually confirm their functionality and styling.
2.  **Style and Polish:** Review each component against the original QMUI iOS to ensure all styling options and animations are ported correctly.
3.  **Implement Missing Components:** Any components not listed above need to be implemented.
4.  **Documentation:** Write comprehensive documentation for each component, including usage examples.
5.  **Testing:** Create a robust test suite, including visual regression tests, to ensure the stability of the library.
6.  **Publish:** Prepare the library for publication as a Go module.

---
### Status Legend
- ✅ **Verified**: Implemented and confirmed working in the demo application.
- 🟡 **Implemented, Needs Demo**: The component's code is present in `pkg/`, but it has not been added to the main demo for visual verification.
- ❌ **Needs Implementation**: The component has not yet been ported from the original QMUI iOS.
- 🚧 **In Progress**: The component is actively being worked on.
