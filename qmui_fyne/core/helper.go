// Package core provides helper utilities for QMUI-Go
// Ported from Tencent's QMUI_iOS framework
package core

import (
	"image/color"
	"math"
	"runtime"

	"fyne.io/fyne/v2"
)

// Helper provides utility functions similar to QMUIHelper
type Helper struct{}

// NewHelper creates a new Helper instance
func NewHelper() *Helper {
	return &Helper{}
}

// Device information helpers

// IsIPad returns whether the current device is an iPad equivalent (large screen)
func (h *Helper) IsIPad() bool {
	// In desktop context, we consider large screens as "iPad-like"
	return false
}

// IsMac returns whether running on macOS
func (h *Helper) IsMac() bool {
	return runtime.GOOS == "darwin"
}

// IsWindows returns whether running on Windows
func (h *Helper) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux returns whether running on Linux
func (h *Helper) IsLinux() bool {
	return runtime.GOOS == "linux"
}

// Color helpers

// ColorWithAlpha returns a new color with the specified alpha value
func ColorWithAlpha(c color.Color, alpha float64) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(alpha * 255),
	}
}

// ColorToRGBA converts any color.Color to RGBA components (0-255)
func ColorToRGBA(c color.Color) (r, g, b, a uint8) {
	if c == nil {
		return 0, 0, 0, 0
	}
	rr, gg, bb, aa := c.RGBA()
	return uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8), uint8(aa >> 8)
}

// ColorFromHex creates a color from a hex string (e.g., "#FF5500" or "FF5500")
func ColorFromHex(hex string) color.Color {
	if len(hex) == 0 {
		return color.Black
	}
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) == 6 {
		hex = hex + "FF" // Add full alpha
	}
	if len(hex) != 8 {
		return color.Black
	}

	var r, g, b, a uint8
	_, _ = hexParseByte(hex[0:2], &r)
	_, _ = hexParseByte(hex[2:4], &g)
	_, _ = hexParseByte(hex[4:6], &b)
	_, _ = hexParseByte(hex[6:8], &a)

	return color.NRGBA{R: r, G: g, B: b, A: a}
}

func hexParseByte(s string, b *uint8) (bool, error) {
	var n uint8
	for _, c := range s {
		n *= 16
		switch {
		case c >= '0' && c <= '9':
			n += uint8(c - '0')
		case c >= 'a' && c <= 'f':
			n += uint8(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			n += uint8(c - 'A' + 10)
		default:
			return false, nil
		}
	}
	*b = n
	return true, nil
}

// BlendColors blends two colors with the specified ratio (0.0 = c1, 1.0 = c2)
func BlendColors(c1, c2 color.Color, ratio float64) color.Color {
	r1, g1, b1, a1 := ColorToRGBA(c1)
	r2, g2, b2, a2 := ColorToRGBA(c2)

	return color.NRGBA{
		R: uint8(float64(r1)*(1-ratio) + float64(r2)*ratio),
		G: uint8(float64(g1)*(1-ratio) + float64(g2)*ratio),
		B: uint8(float64(b1)*(1-ratio) + float64(b2)*ratio),
		A: uint8(float64(a1)*(1-ratio) + float64(a2)*ratio),
	}
}

// IsDarkColor returns whether a color is considered "dark"
func IsDarkColor(c color.Color) bool {
	r, g, b, _ := ColorToRGBA(c)
	// Using perceived luminance formula
	luminance := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return luminance < 128
}

// Math helpers

// Clamp constrains a value between min and max
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat64 constrains a float64 value between min and max
func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Lerp performs linear interpolation between a and b
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

// Size helpers

// SizeFits checks if innerSize fits within outerSize
func SizeFits(inner, outer fyne.Size) bool {
	return inner.Width <= outer.Width && inner.Height <= outer.Height
}

// SizeScale scales a size by a factor
func SizeScale(size fyne.Size, factor float32) fyne.Size {
	return fyne.NewSize(size.Width*factor, size.Height*factor)
}

// SizeMax returns a size with the max of each dimension
func SizeMax(s1, s2 fyne.Size) fyne.Size {
	return fyne.NewSize(
		float32(math.Max(float64(s1.Width), float64(s2.Width))),
		float32(math.Max(float64(s1.Height), float64(s2.Height))),
	)
}

// SizeMin returns a size with the min of each dimension
func SizeMin(s1, s2 fyne.Size) fyne.Size {
	return fyne.NewSize(
		float32(math.Min(float64(s1.Width), float64(s2.Width))),
		float32(math.Min(float64(s1.Height), float64(s2.Height))),
	)
}

// Position helpers

// PositionOffset offsets a position by the given amounts
func PositionOffset(pos fyne.Position, dx, dy float32) fyne.Position {
	return fyne.NewPos(pos.X+dx, pos.Y+dy)
}

// CenterPosition returns the center position within a size
func CenterPosition(size fyne.Size) fyne.Position {
	return fyne.NewPos(size.Width/2, size.Height/2)
}

// CenterInRect returns the position to center innerSize within outerSize
func CenterInRect(innerSize, outerSize fyne.Size) fyne.Position {
	return fyne.NewPos(
		(outerSize.Width-innerSize.Width)/2,
		(outerSize.Height-innerSize.Height)/2,
	)
}

// String helpers

// TruncateString truncates a string to maxLen with an optional suffix
func TruncateString(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= len(suffix) {
		return suffix[:maxLen]
	}
	return s[:maxLen-len(suffix)] + suffix
}

// StringLength returns the length of a string, optionally counting non-ASCII as 2
func StringLength(s string, countNonASCIIAsTwo bool) int {
	if !countNonASCIIAsTwo {
		return len([]rune(s))
	}
	count := 0
	for _, r := range s {
		if r > 127 {
			count += 2
		} else {
			count++
		}
	}
	return count
}

// Animation timing functions

// EaseInOutQuad provides quadratic ease-in-out
func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseOutQuad provides quadratic ease-out
func EaseOutQuad(t float64) float64 {
	return t * (2 - t)
}

// EaseInQuad provides quadratic ease-in
func EaseInQuad(t float64) float64 {
	return t * t
}

// EaseOutCubic provides cubic ease-out
func EaseOutCubic(t float64) float64 {
	t--
	return t*t*t + 1
}

// EaseInCubic provides cubic ease-in
func EaseInCubic(t float64) float64 {
	return t * t * t
}

// EaseOutBack provides back ease-out (overshoot)
func EaseOutBack(t float64) float64 {
	c1 := 1.70158
	c3 := c1 + 1
	return 1 + c3*math.Pow(t-1, 3) + c1*math.Pow(t-1, 2)
}

// EaseInBack provides back ease-in (overshoot)
func EaseInBack(t float64) float64 {
	c1 := 1.70158
	c3 := c1 + 1
	return c3*t*t*t - c1*t*t
}

// EaseOutElastic provides elastic ease-out
func EaseOutElastic(t float64) float64 {
	c4 := (2 * math.Pi) / 3
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	return math.Pow(2, -10*t)*math.Sin((t*10-0.75)*c4) + 1
}

// EaseOutBounce provides bounce ease-out
func EaseOutBounce(t float64) float64 {
	n1 := 7.5625
	d1 := 2.75
	if t < 1/d1 {
		return n1 * t * t
	} else if t < 2/d1 {
		t -= 1.5 / d1
		return n1*t*t + 0.75
	} else if t < 2.5/d1 {
		t -= 2.25 / d1
		return n1*t*t + 0.9375
	} else {
		t -= 2.625 / d1
		return n1*t*t + 0.984375
	}
}

// Linear provides linear timing (no easing)
func Linear(t float64) float64 {
	return t
}
