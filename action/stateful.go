package action

// Stateful represents an object who's state may be set.
//
// It is commonly used by [Action]s, and the state will be
// set by various functions within the source testing packages.
type Stateful[TState any] interface {
	// SetState overwrites currently held state with the given value.
	SetState(TState)
}
