// Package animation provides QMUIAnimation - animation utilities and easing functions
// Ported from Tencent's QMUI_iOS framework
package animation

import (
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// EasingFunction defines an easing function type
type EasingFunction func(t float64) float64

// Predefined easing functions

// Linear provides linear timing (no easing)
func Linear(t float64) float64 {
	return t
}

// EaseInQuad provides quadratic ease-in
func EaseInQuad(t float64) float64 {
	return t * t
}

// EaseOutQuad provides quadratic ease-out
func EaseOutQuad(t float64) float64 {
	return t * (2 - t)
}

// EaseInOutQuad provides quadratic ease-in-out
func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseInCubic provides cubic ease-in
func EaseInCubic(t float64) float64 {
	return t * t * t
}

// EaseOutCubic provides cubic ease-out
func EaseOutCubic(t float64) float64 {
	t--
	return t*t*t + 1
}

// EaseInOutCubic provides cubic ease-in-out
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return (t-1)*(2*t-2)*(2*t-2) + 1
}

// EaseInQuart provides quartic ease-in
func EaseInQuart(t float64) float64 {
	return t * t * t * t
}

// EaseOutQuart provides quartic ease-out
func EaseOutQuart(t float64) float64 {
	t--
	return 1 - t*t*t*t
}

// EaseInQuint provides quintic ease-in
func EaseInQuint(t float64) float64 {
	return t * t * t * t * t
}

// EaseOutQuint provides quintic ease-out
func EaseOutQuint(t float64) float64 {
	t--
	return 1 + t*t*t*t*t
}

// EaseInSine provides sine ease-in
func EaseInSine(t float64) float64 {
	return 1 - math.Cos(t*math.Pi/2)
}

// EaseOutSine provides sine ease-out
func EaseOutSine(t float64) float64 {
	return math.Sin(t * math.Pi / 2)
}

// EaseInOutSine provides sine ease-in-out
func EaseInOutSine(t float64) float64 {
	return -(math.Cos(math.Pi*t) - 1) / 2
}

// EaseInExpo provides exponential ease-in
func EaseInExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	return math.Pow(2, 10*(t-1))
}

// EaseOutExpo provides exponential ease-out
func EaseOutExpo(t float64) float64 {
	if t == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

// EaseInCirc provides circular ease-in
func EaseInCirc(t float64) float64 {
	return 1 - math.Sqrt(1-t*t)
}

// EaseOutCirc provides circular ease-out
func EaseOutCirc(t float64) float64 {
	t--
	return math.Sqrt(1 - t*t)
}

// EaseInBack provides back ease-in (overshoot)
func EaseInBack(t float64) float64 {
	c1 := 1.70158
	c3 := c1 + 1
	return c3*t*t*t - c1*t*t
}

// EaseOutBack provides back ease-out (overshoot)
func EaseOutBack(t float64) float64 {
	c1 := 1.70158
	c3 := c1 + 1
	return 1 + c3*math.Pow(t-1, 3) + c1*math.Pow(t-1, 2)
}

// EaseInElastic provides elastic ease-in
func EaseInElastic(t float64) float64 {
	c4 := (2 * math.Pi) / 3
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	return -math.Pow(2, 10*t-10) * math.Sin((t*10-10.75)*c4)
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

// EaseInBounce provides bounce ease-in
func EaseInBounce(t float64) float64 {
	return 1 - EaseOutBounce(1-t)
}

// Spring provides spring animation with damping
func Spring(damping, stiffness float64) EasingFunction {
	return func(t float64) float64 {
		return 1 - math.Exp(-damping*t)*math.Cos(stiffness*t)
	}
}

// Animation represents an animation instance
type Animation struct {
	Duration    time.Duration
	Easing      EasingFunction
	OnUpdate    func(progress float64)
	OnComplete  func()

	mu        sync.RWMutex
	running   bool
	cancelled bool
	stopChan  chan struct{}
}

// NewAnimation creates a new animation
func NewAnimation(duration time.Duration, easing EasingFunction, onUpdate func(float64)) *Animation {
	return &Animation{
		Duration: duration,
		Easing:   easing,
		OnUpdate: onUpdate,
	}
}

// Start begins the animation
func (a *Animation) Start() {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return
	}
	a.running = true
	a.cancelled = false
	a.stopChan = make(chan struct{})
	a.mu.Unlock()

	go a.run()
}

// Stop stops the animation
func (a *Animation) Stop() {
	a.mu.Lock()
	if !a.running {
		a.mu.Unlock()
		return
	}
	a.cancelled = true
	if a.stopChan != nil {
		close(a.stopChan)
	}
	a.mu.Unlock()
}

// IsRunning returns whether the animation is running
func (a *Animation) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

func (a *Animation) run() {
	startTime := time.Now()
	ticker := time.NewTicker(time.Millisecond * 16) // ~60fps
	defer ticker.Stop()

	for {
		select {
		case <-a.stopChan:
			a.mu.Lock()
			a.running = false
			a.mu.Unlock()
			return
		case <-ticker.C:
			elapsed := time.Since(startTime)
			t := float64(elapsed) / float64(a.Duration)

			if t >= 1 {
				t = 1
			}

			a.mu.RLock()
			easing := a.Easing
			a.mu.RUnlock()

			progress := t
			if easing != nil {
				progress = easing(t)
			}

			if a.OnUpdate != nil {
				a.OnUpdate(progress)
			}

			if t >= 1 {
				a.mu.Lock()
				a.running = false
				onComplete := a.OnComplete
				a.mu.Unlock()

				if onComplete != nil {
					onComplete()
				}
				return
			}
		}
	}
}

// Animator manages multiple animations
type Animator struct {
	animations []*Animation
	mu         sync.Mutex
}

// NewAnimator creates a new animator
func NewAnimator() *Animator {
	return &Animator{
		animations: make([]*Animation, 0),
	}
}

// AddAnimation adds an animation
func (an *Animator) AddAnimation(animation *Animation) {
	an.mu.Lock()
	an.animations = append(an.animations, animation)
	an.mu.Unlock()
}

// StartAll starts all animations
func (an *Animator) StartAll() {
	an.mu.Lock()
	animations := an.animations
	an.mu.Unlock()

	for _, anim := range animations {
		anim.Start()
	}
}

// StopAll stops all animations
func (an *Animator) StopAll() {
	an.mu.Lock()
	animations := an.animations
	an.mu.Unlock()

	for _, anim := range animations {
		anim.Stop()
	}
}

// PropertyAnimation animates a property between values
type PropertyAnimation struct {
	*Animation
	FromValue float64
	ToValue   float64
}

// NewPropertyAnimation creates a property animation
func NewPropertyAnimation(from, to float64, duration time.Duration, easing EasingFunction, onUpdate func(value float64)) *PropertyAnimation {
	pa := &PropertyAnimation{
		FromValue: from,
		ToValue:   to,
	}

	pa.Animation = NewAnimation(duration, easing, func(progress float64) {
		value := from + (to-from)*progress
		if onUpdate != nil {
			onUpdate(value)
		}
	})

	return pa
}

// ColorAnimation animates between two colors
type ColorAnimation struct {
	*Animation
	FromR, FromG, FromB, FromA float64
	ToR, ToG, ToB, ToA         float64
}

// NewColorAnimation creates a color animation
func NewColorAnimation(fromR, fromG, fromB, fromA, toR, toG, toB, toA float64,
	duration time.Duration, easing EasingFunction, onUpdate func(r, g, b, a float64)) *ColorAnimation {

	ca := &ColorAnimation{
		FromR: fromR, FromG: fromG, FromB: fromB, FromA: fromA,
		ToR: toR, ToG: toG, ToB: toB, ToA: toA,
	}

	ca.Animation = NewAnimation(duration, easing, func(progress float64) {
		r := fromR + (toR-fromR)*progress
		g := fromG + (toG-fromG)*progress
		b := fromB + (toB-fromB)*progress
		a := fromA + (toA-fromA)*progress
		if onUpdate != nil {
			onUpdate(r, g, b, a)
		}
	})

	return ca
}

// PositionAnimation animates a position
type PositionAnimation struct {
	*Animation
	FromX, FromY float64
	ToX, ToY     float64
}

// NewPositionAnimation creates a position animation
func NewPositionAnimation(fromX, fromY, toX, toY float64,
	duration time.Duration, easing EasingFunction, onUpdate func(x, y float64)) *PositionAnimation {

	pa := &PositionAnimation{
		FromX: fromX, FromY: fromY,
		ToX: toX, ToY: toY,
	}

	pa.Animation = NewAnimation(duration, easing, func(progress float64) {
		x := fromX + (toX-fromX)*progress
		y := fromY + (toY-fromY)*progress
		if onUpdate != nil {
			onUpdate(x, y)
		}
	})

	return pa
}

// AnimateFloat animates a float value
func AnimateFloat(from, to float64, duration time.Duration, easing EasingFunction, onUpdate func(float64), onComplete func()) {
	anim := NewPropertyAnimation(from, to, duration, easing, onUpdate)
	anim.OnComplete = onComplete
	anim.Start()
}

// AnimateSize animates a size
func AnimateSize(from, to fyne.Size, duration time.Duration, easing EasingFunction, onUpdate func(fyne.Size), onComplete func()) {
	anim := NewAnimation(duration, easing, func(progress float64) {
		width := float64(from.Width) + (float64(to.Width)-float64(from.Width))*progress
		height := float64(from.Height) + (float64(to.Height)-float64(from.Height))*progress
		if onUpdate != nil {
			onUpdate(fyne.NewSize(float32(width), float32(height)))
		}
	})
	anim.OnComplete = onComplete
	anim.Start()
}

// AnimatePosition animates a position
func AnimatePosition(from, to fyne.Position, duration time.Duration, easing EasingFunction, onUpdate func(fyne.Position), onComplete func()) {
	anim := NewAnimation(duration, easing, func(progress float64) {
		x := float64(from.X) + (float64(to.X)-float64(from.X))*progress
		y := float64(from.Y) + (float64(to.Y)-float64(from.Y))*progress
		if onUpdate != nil {
			onUpdate(fyne.NewPos(float32(x), float32(y)))
		}
	})
	anim.OnComplete = onComplete
	anim.Start()
}
