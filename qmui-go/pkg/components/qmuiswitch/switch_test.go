// Package switch_test provides tests for the Switch component
package qmuiswitch_test

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/user/qmui-go/pkg/components/qmuiswitch"
)

func TestSwitch_NewSwitch(t *testing.T) {
	s := qmuiswitch.NewSwitch(nil)
	assert.False(t, s.Checked)
	assert.True(t, s.Enabled)
}

func TestSwitch_SetChecked(t *testing.T) {
	var changed bool
	s := qmuiswitch.NewSwitch(func(b bool) {
		changed = b
	})

	s.SetChecked(true)
	assert.True(t, s.Checked)
	assert.True(t, changed)

	s.SetChecked(false)
	assert.False(t, s.Checked)
	assert.False(t, changed)
}

func TestSwitch_Toggle(t *testing.T) {
	s := qmuiswitch.NewSwitch(nil)
	s.Toggle()
	assert.True(t, s.Checked)
	s.Toggle()
	assert.False(t, s.Checked)
}

func TestSwitch_Tapped(t *testing.T) {
	s := qmuiswitch.NewSwitch(nil)
	test.Tap(s)
	assert.True(t, s.Checked)
}
