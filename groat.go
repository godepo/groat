package groat

import (
	"testing"
)

type Provide[Depends any, SUT any] func(t *testing.T, deps Depends) SUT
type Before[Depends any] func(t *testing.T, deps Depends) Depends
type After[Depends any] func(t *testing.T, deps Depends)
type Given[State any] func(t *testing.T, state State) State
type When[Depends any, State any] func(t *testing.T, deps Depends, state State) State
type Then[State any] func(t *testing.T, state State)

type Case[Depends any, State any, SUT any] struct {
	State    State
	Deps     Depends
	SUT      SUT
	provider Provide[Depends, SUT]
	T        *testing.T
	then     []Then[State]
	after    []After[Depends]
}

func New[Depends any, State any, SUT any](
	t *testing.T,
	provider Provide[Depends, SUT],
	all ...Before[Depends],
) *Case[Depends, State, SUT] {
	t.Helper()

	var tcs Case[Depends, State, SUT]
	tcs.T = t
	tcs.Before(all...)
	tcs.After(func(t *testing.T, deps Depends) {
		t.Helper()
		for _, then := range tcs.then {
			then(tcs.T, tcs.State)
		}
	})
	t.Cleanup(func() {
		for _, after := range tcs.after {
			after(tcs.T, tcs.Deps)
		}
	})
	tcs.provider = provider
	return &tcs
}

func (tcs *Case[Depends, State, SUT]) Go() {
	tcs.SUT = tcs.provider(tcs.T, tcs.Deps)
}

func (tcs *Case[Depends, State, SUT]) Before(all ...Before[Depends]) {
	for _, before := range all {
		tcs.Deps = before(tcs.T, tcs.Deps)
	}
}

func (tcs *Case[Depends, State, SUT]) After(all ...After[Depends]) *Case[Depends, State, SUT] {
	tcs.after = append(tcs.after, all...)
	return tcs
}

func (tcs *Case[Depends, State, SUT]) Given(all ...Given[State]) *Case[Depends, State, SUT] {
	for _, given := range all {
		tcs.State = given(tcs.T, tcs.State)
	}
	return tcs
}

func (tcs *Case[Depends, State, SUT]) When(all ...When[Depends, State]) *Case[Depends, State, SUT] {
	for _, when := range all {
		tcs.State = when(tcs.T, tcs.Deps, tcs.State)
	}
	return tcs
}

func (tcs *Case[Depends, State, SUT]) Then(all ...Then[State]) *Case[Depends, State, SUT] {
	tcs.then = append(tcs.then, all...)
	return tcs
}
