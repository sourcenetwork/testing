package testo

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdT "testing"

	"github.com/sourcenetwork/testo/action"
	"github.com/sourcenetwork/testo/multiplier"
)

// Execute this set of actions, serially, in order.
func Execute(a action.Actions) {
	for _, action := range a {
		action.Execute()
	}
}

// Execute this set of actions upon the given state, serially, in order.
func ExecuteS[TState any](actions action.Actions, s TState) {
	for _, a := range actions {
		if stateful, ok := a.(action.Stateful[TState]); ok {
			stateful.SetState(s)
		}

		a.Execute()
	}
}

// Log the set of active multipliers and the provided actions.
//
// The first line will be empty.
// Multipliers will be logged on the second line, as comma separated values.
// Actions will be logged on the third line as prettified json, including the
// action concrete type in a `_type` property.
//
// For example:
//
//	 Multipliers: txn-commit
//	 Actions: [
//	  {
//		"_type": "*action.StartCli"
//	  },
//	  {
//		"_type": "*action.TxCreate"
//	  },
//	  {
//		"_type": "*action.TxCommit",
//		"TxnIndex": 0
//	  }
//	]
func Log(t stdT.TB, actions action.Actions) {
	typedActions := make([]json.RawMessage, 0, len(actions))
	for _, a := range actions {
		actionJson, _ := json.MarshalIndent(a, "", "  ")

		var deliminator string
		if !bytes.Equal(actionJson, []byte("{}")) {
			deliminator = ","
		}

		jsonString := fmt.Sprintf("{\n  \"_type\":\"%T\"%s%s", a, deliminator, actionJson[1:])
		typedActions = append(typedActions, json.RawMessage(jsonString))
	}

	jsonB, err := json.MarshalIndent(typedActions, "", "  ")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("\nMultipliers: %s\nActions: %s", multiplier.Get(), string(jsonB))
}
