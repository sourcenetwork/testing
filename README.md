# testing
Shared testing framework for Source projects.

This repository provides a shared language for writing tests across Source Network projects.  Tests are composed entirely of `action.Action`s. `multiplier.Multiplier`s may be configured to mutate the declared actions to scale with various project-specific complexity multipliers.

### Example consumption

The following example taken from the Defra CLI package, and shows a simple test `TestSchemaAdd` along with the shared test executor for the test-package.
```
// testing.go

package integration

func init() {
	multiplier.Init("DEFRA_MULTIPLIERS", "txn-commit")
}

type Test struct {
	// The test will be skipped if the current active set of multipliers
	// does not contain all of the given multiplier names.
	Includes []multiplier.Name

	// The test will be skipped if the current active set of multipliers
	// contains any of the given multiplier names.
	Excludes []multiplier.Name

	// Actions contains the set of actions that should be
	// executed as part of this test.
	Actions action.Actions
}

func (test *Test) Execute(t testing.TB) {
	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 1*time.Second)

	multiplier.Skip(t, test.Includes, test.Excludes)

	actions := prependStart(test.Actions)

	actions = multiplier.Apply(actions)

	testing.Log(t, actions)

	testing.ExecuteS(actions, &state.State{
		T:       t,
		Ctx:     ctx,
		Cancels: []context.CancelFunc{cancel},
	})
}
```
```
// multiplier/txn.go

func init() {
	multiplier.Register(&txnCommit{})
}

const TxnCommit Name = "txn-commit"

type txnCommit struct{}

var _ Multiplier = (*txnCommit)(nil)

func (n *txnCommit) Name() Name {
	return TxnCommit
}

func (n *txnCommit) Apply(source action.Actions) action.Actions {
   ...
}
```
```
// simple_test.go

func TestSchemaAdd(t *testing.T) {
	test := &integration.Test{
		Actions: []action.Action{
			&action.SchemaAdd{
				InlineSchema: `
					type User {}
				`,
			},
			&action.CollectionDescribe{
				Expected: []client.CollectionDefinition{
					{
						Description: client.CollectionDescription{
							Name:           immutable.Some("User"),
							IsMaterialized: true,
						},
					},
				},
			},
		},
	}

	test.Execute(t)
}
```
