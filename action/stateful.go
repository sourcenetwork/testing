package action

// Stateful represents an object who's state may be set.
//
// It is commonly used by [Action]s, and the state will be
// set by various functions within the source testing packages.
//
// For example, [github.com/sourcenetwork/testing.ExecuteS] calls
// [Stateful.SetState] for every action provided to it, allowing those
// given actions to access the test state during test execution.
type Stateful[TState any] interface {
	// SetState overwrites currently held state with the given value.
	SetState(TState)
}
