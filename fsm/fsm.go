package fsm

import (
	"golang.org/x/exp/constraints"
)

// FSM is represents a Finite State Machine
type FSM[T constraints.Integer] struct {
	current     T
	states      []T
	transitions map[T]map[T]struct{}
}

// New creates and returns an FSM with the states provided. No transitions
// between states will exist.
func New[T constraints.Integer](states ...T) *FSM[T] {
	var current T
	if len(states) > 0 {
		current = states[0]
	}
	sm := &FSM[T]{
		current:     current,
		transitions: make(map[T]map[T]struct{}, len(states)),
	}
	for _, s := range states {
		sm.AddState(s)
	}
	return sm
}

// Current returns the current state the FSM is on.
func (sm *FSM[T]) Current() T {
	return sm.current
}

// AddState adds a new state to the FSM. If the state already exists, this
// method does nothing.
func (sm *FSM[T]) AddState(state T) *FSM[T] {
	if _, ok := sm.transitions[state]; !ok {
		sm.transitions[state] = make(map[T]struct{}, 0)
	}
	return sm
}

// AddTransition adds a new transition between two states of the FSM.
func (sm *FSM[T]) AddTransition(from, to T) *FSM[T] {
	if _, ok := sm.transitions[from]; !ok {
		sm.transitions[from] = make(map[T]struct{})
	}
	sm.transitions[from][to] = struct{}{}
	return sm
}

// To changes the state to the given state, if the transition is allowed from
// the current one.
func (sm *FSM[T]) To(state T) bool {
	if nextStates, ok := sm.transitions[sm.current]; ok {
		if _, ok = nextStates[state]; ok {
			sm.current = state
			return true
		}
	}
	return false
}

// Set forcibly changes the state to the given state (if it exists), regardless
// of the configured transitions.
func (sm *FSM[T]) Set(state T) {
	if _, ok := sm.transitions[state]; ok {
		sm.current = state
	}
}
