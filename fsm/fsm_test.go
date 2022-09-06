package fsm_test

import (
	. "github.com/moolmanruan/ebitengine-test/fsm"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetup(t *testing.T) {
	type SomeState int
	const (
		One   SomeState = 1
		Two   SomeState = 2
		Three SomeState = 3
	)
	sm := New(One, Two, Three).
		AddTransition(One, Two).
		AddTransition(Two, Three)

	require.Equal(t, sm.Current(), One)
}

func TestTo(t *testing.T) {
	type AlphaState uint8
	const (
		A AlphaState = 1
		B AlphaState = 2
		C AlphaState = 3
		D AlphaState = 3
	)
	sm := New(A, B, C, D).
		AddTransition(A, B).
		AddTransition(B, C).
		AddTransition(B, D).
		AddTransition(C, A)

	ok := sm.To(B)
	require.True(t, ok)
	require.Equal(t, sm.Current(), B)

	ok = sm.To(C)
	require.True(t, ok)
	require.Equal(t, sm.Current(), C)

	ok = sm.To(D)
	require.False(t, ok)
	require.Equal(t, sm.Current(), C)

	ok = sm.To(A)
	require.True(t, ok)
	require.Equal(t, sm.Current(), A)
}

func TestSet(t *testing.T) {
	type AlphaState uint8
	const (
		A AlphaState = 10
		B AlphaState = 20
		C AlphaState = 30
		D AlphaState = 40
	)
	sm := New(A, B, C)

	require.Equal(t, sm.Current(), A)

	sm.Set(C)
	require.Equal(t, sm.Current(), C,
		"Should change to C, even if no transition exists")

	sm.Set(D)
	require.Equal(t, sm.Current(), C,
		"Should stay C, since D is not a registered state")
}
