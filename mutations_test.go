package graphql_test

import (
	"testing"

	"github.com/fraym/graphql-go"
	"github.com/fraym/graphql-go/gqlerrors"
	"github.com/fraym/graphql-go/language/location"
	"github.com/fraym/graphql-go/testutil"
	"github.com/stretchr/testify/assert"
)

// testNumberHolder maps to numberHolderType
type testNumberHolder struct {
	TheNumber int64 `json:"theNumber"` // map field to `theNumber` so it can be resolve by the default ResolveFn
}
type testRoot struct {
	NumberHolder *testNumberHolder
}

func newTestRoot(originalNumber int64) *testRoot {
	return &testRoot{
		NumberHolder: &testNumberHolder{originalNumber},
	}
}

func (r *testRoot) ImmediatelyChangeTheNumber(newNumber int64) *testNumberHolder {
	r.NumberHolder.TheNumber = newNumber
	return r.NumberHolder
}

func (r *testRoot) PromiseToChangeTheNumber(newNumber int64) *testNumberHolder {
	return r.ImmediatelyChangeTheNumber(newNumber)
}

func (r *testRoot) FailToChangeTheNumber(newNumber int64) *testNumberHolder {
	panic("Cannot change the number")
}

func (r *testRoot) PromiseAndFailToChangeTheNumber(newNumber int64) *testNumberHolder {
	panic("Cannot change the number")
}

// numberHolderType creates a mapping to testNumberHolder
var numberHolderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NumberHolder",
	Fields: graphql.Fields{
		"theNumber": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var mutationsTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"numberHolder": &graphql.Field{
				Type: numberHolderType,
			},
		},
	}),
	Mutation: graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"immediatelyChangeTheNumber": &graphql.Field{
				Type: numberHolderType,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "newNumber",
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					obj, _ := p.Source.(*testRoot)
					var newNumber int64
					newNumber, _ = p.Args["newNumber"].(int64)
					return obj.ImmediatelyChangeTheNumber(newNumber), nil
				},
			},
			"promiseToChangeTheNumber": &graphql.Field{
				Type: numberHolderType,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "newNumber",
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					obj, _ := p.Source.(*testRoot)
					var newNumber int64
					newNumber, _ = p.Args["newNumber"].(int64)
					return obj.PromiseToChangeTheNumber(newNumber), nil
				},
			},
			"failToChangeTheNumber": &graphql.Field{
				Type: numberHolderType,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "newNumber",
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					obj, _ := p.Source.(*testRoot)
					var newNumber int64
					newNumber, _ = p.Args["newNumber"].(int64)
					return obj.FailToChangeTheNumber(newNumber), nil
				},
			},
			"promiseAndFailToChangeTheNumber": &graphql.Field{
				Type: numberHolderType,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "newNumber",
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					obj, _ := p.Source.(*testRoot)
					var newNumber int64
					newNumber, _ = p.Args["newNumber"].(int64)
					return obj.PromiseAndFailToChangeTheNumber(newNumber), nil
				},
			},
		},
	}),
})

func TestMutations_ExecutionOrdering_EvaluatesMutationsSerially(t *testing.T) {
	root := newTestRoot(6)
	doc := `mutation M {
      first: immediatelyChangeTheNumber(newNumber: 1) {
        theNumber
      },
      second: promiseToChangeTheNumber(newNumber: 2) {
        theNumber
      },
      third: immediatelyChangeTheNumber(newNumber: 3) {
        theNumber
      }
      fourth: promiseToChangeTheNumber(newNumber: 4) {
        theNumber
      },
      fifth: immediatelyChangeTheNumber(newNumber: 5) {
        theNumber
      }
    }`

	expected := &graphql.Result{
		Data: map[string]any{
			"first": map[string]any{
				"theNumber": int64(1),
			},
			"second": map[string]any{
				"theNumber": int64(2),
			},
			"third": map[string]any{
				"theNumber": int64(3),
			},
			"fourth": map[string]any{
				"theNumber": int64(4),
			},
			"fifth": map[string]any{
				"theNumber": int64(5),
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: mutationsTestSchema,
		AST:    ast,
		Root:   root,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	assert.Equal(t, expected, result)
}

func TestMutations_EvaluatesMutationsCorrectlyInThePresenceOfAFailedMutation(t *testing.T) {
	root := newTestRoot(6)
	doc := `mutation M {
      first: immediatelyChangeTheNumber(newNumber: 1) {
        theNumber
      },
      second: promiseToChangeTheNumber(newNumber: 2) {
        theNumber
      },
      third: failToChangeTheNumber(newNumber: 3) {
        theNumber
      }
      fourth: promiseToChangeTheNumber(newNumber: 4) {
        theNumber
      },
      fifth: immediatelyChangeTheNumber(newNumber: 5) {
        theNumber
      }
      sixth: promiseAndFailToChangeTheNumber(newNumber: 6) {
        theNumber
      }
    }`

	expected := &graphql.Result{
		Data: map[string]any{
			"first": map[string]any{
				"theNumber": 1,
			},
			"second": map[string]any{
				"theNumber": 2,
			},
			"third": nil,
			"fourth": map[string]any{
				"theNumber": 4,
			},
			"fifth": map[string]any{
				"theNumber": 5,
			},
			"sixth": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot change the number`,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 7},
				},
			},
			{
				Message: `Cannot change the number`,
				Locations: []location.SourceLocation{
					{Line: 17, Column: 7},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: mutationsTestSchema,
		AST:    ast,
		Root:   root,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	t.Skipf("Testing equality for slice of errors in results")
	assert.Equal(t, expected, result)
}
