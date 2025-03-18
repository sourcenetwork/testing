package action

// Action represents an item to be executed as part of a test.
//
// For example an action may set or get a value, or it may close the datastore.
type Action interface {
	// Execute this action upon the given state.
	Execute()
}

// Actions is an executable set of [Action]s.
type Actions = []Action
