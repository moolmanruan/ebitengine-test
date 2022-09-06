package fsm

import (
	"golang.org/x/exp/constraints"
)

type StateMachine[T constraints.Integer] struct {
	current     T
	states      []T
	transitions map[T]map[T]struct{}
}

func New[T constraints.Integer](states ...T) *StateMachine[T] {
	var current T
	if len(states) > 0 {
		current = states[0]
	}
	sm := &StateMachine[T]{
		current:     current,
		transitions: make(map[T]map[T]struct{}, len(states)),
	}
	for _, s := range states {
		sm.AddState(s)
	}
	return sm
}

func (sm *StateMachine[T]) AddState(state T) *StateMachine[T] {
	sm.transitions[state] = make(map[T]struct{}, 0)
	return sm
}

func (sm *StateMachine[T]) AddTransition(from, to T) *StateMachine[T] {
	if _, ok := sm.transitions[from]; !ok {
		sm.transitions[from] = make(map[T]struct{})
	}
	sm.transitions[from][to] = struct{}{}
	return sm
}

// To changes the state to the given state, if the transition is allowed from
// the current one.
func (sm *StateMachine[T]) To(state T) bool {
	if nextStates, ok := sm.transitions[sm.current]; ok {
		if _, ok = nextStates[state]; ok {
			sm.current = state
			return true
		}
	}
	return false
}

// Set changes the state to the given state (if it exists), regardless of
// the configured transitions.
func (sm *StateMachine[T]) Set(state T) {
	if _, ok := sm.transitions[state]; ok {
		sm.current = state
	}
}

func (sm *StateMachine[T]) Current() T {
	return sm.current
}
