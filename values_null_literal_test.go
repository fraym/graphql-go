package graphql_test

import (
	"encoding/json"
	"testing"

	"github.com/fraym/graphql-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nullLiteralFlatrateInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NullLiteralFlatrateInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"amount": &graphql.InputObjectFieldConfig{
			Type: graphql.Float,
		},
	},
})

var nullLiteralSystemFeeInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NullLiteralSystemFeeInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"fixed": &graphql.InputObjectFieldConfig{
			Type: graphql.Float,
		},
	},
})

var nullLiteralPaymentPlanEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "NullLiteralPaymentPlan",
	Values: graphql.EnumValueConfigMap{
		"PAYASYOUGO": &graphql.EnumValueConfig{
			Value: "PAYASYOUGO",
		},
	},
})

var nullLiteralDataInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NullLiteralDataInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"paymentPlan": &graphql.InputObjectFieldConfig{
			Type: nullLiteralPaymentPlanEnum,
		},
		"flatrate": &graphql.InputObjectFieldConfig{
			Type: nullLiteralFlatrateInput,
		},
		"systemFee": &graphql.InputObjectFieldConfig{
			Type: nullLiteralSystemFeeInput,
		},
		"tags": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(graphql.String),
		},
	},
})

var nullLiteralUpdateInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NullLiteralUpdateInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"data": &graphql.InputObjectFieldConfig{
			Type: nullLiteralDataInput,
		},
	},
})

var nullLiteralTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"noop": &graphql.Field{
				Type: graphql.String,
			},
		},
	}),
	Mutation: graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"updateOrganizer": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "input",
						Type: nullLiteralUpdateInput,
					},
				},
				// Serialize the args the executor built so that the difference
				// between an absent field and an explicit null stays visible.
				Resolve: func(p graphql.ResolveParams) (any, error) {
					b, err := json.Marshal(p.Args)
					if err != nil {
						return nil, err
					}
					return string(b), nil
				},
			},
			"updateOrganizerNonNull": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "input",
						Type: graphql.NewNonNull(nullLiteralUpdateInput),
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					return "ok", nil
				},
			},
		},
	}),
})

func executeNullLiteralMutation(t *testing.T, query string, variables map[string]any) *graphql.Result {
	t.Helper()

	return graphql.Do(graphql.Params{
		Schema:         nullLiteralTestSchema,
		RequestString:  query,
		VariableValues: variables,
	})
}

func TestNullLiteral_InputObjectField_IsClearedAndKeepsSiblings(t *testing.T) {
	result := executeNullLiteralMutation(t, `
		mutation {
			updateOrganizer(input: {
				id: "organizer-1"
				data: {
					paymentPlan: PAYASYOUGO
					flatrate: null
					systemFee: {fixed: 1.5}
				}
			})
		}
	`, nil)

	require.Empty(t, result.Errors)
	assert.Equal(t, map[string]any{
		"updateOrganizer": `{"input":{"data":{"flatrate":null,"paymentPlan":"PAYASYOUGO","systemFee":{"fixed":1.5}},"id":"organizer-1"}}`,
	}, result.Data)
}

func TestNullLiteral_InputObjectField_MatchesNullVariable(t *testing.T) {
	literal := executeNullLiteralMutation(t, `
		mutation {
			updateOrganizer(input: {
				id: "organizer-1"
				data: {
					paymentPlan: PAYASYOUGO
					flatrate: null
					systemFee: {fixed: 1.5}
				}
			})
		}
	`, nil)

	variable := executeNullLiteralMutation(t, `
		mutation ($flatrate: NullLiteralFlatrateInput) {
			updateOrganizer(input: {
				id: "organizer-1"
				data: {
					paymentPlan: PAYASYOUGO
					flatrate: $flatrate
					systemFee: {fixed: 1.5}
				}
			})
		}
	`, map[string]any{"flatrate": nil})

	require.Empty(t, literal.Errors)
	require.Empty(t, variable.Errors)
	assert.Equal(t, variable.Data, literal.Data)
}

func TestNullLiteral_OmittedInputObjectField_StaysAbsent(t *testing.T) {
	result := executeNullLiteralMutation(t, `
		mutation {
			updateOrganizer(input: {
				id: "organizer-1"
				data: {
					paymentPlan: PAYASYOUGO
					systemFee: {fixed: 1.5}
				}
			})
		}
	`, nil)

	require.Empty(t, result.Errors)
	assert.Equal(t, map[string]any{
		"updateOrganizer": `{"input":{"data":{"paymentPlan":"PAYASYOUGO","systemFee":{"fixed":1.5}},"id":"organizer-1"}}`,
	}, result.Data)
}

func TestNullLiteral_ListAndEnumFields_AreCleared(t *testing.T) {
	result := executeNullLiteralMutation(t, `
		mutation {
			updateOrganizer(input: {
				id: "organizer-1"
				data: {
					paymentPlan: null
					tags: null
				}
			})
		}
	`, nil)

	require.Empty(t, result.Errors)
	assert.Equal(t, map[string]any{
		"updateOrganizer": `{"input":{"data":{"paymentPlan":null,"tags":null},"id":"organizer-1"}}`,
	}, result.Data)
}

func TestNullLiteral_NestedInputObjectField_IsCleared(t *testing.T) {
	result := executeNullLiteralMutation(t, `
		mutation {
			updateOrganizer(input: {
				id: "organizer-1"
				data: null
			})
		}
	`, nil)

	require.Empty(t, result.Errors)
	assert.Equal(t, map[string]any{
		"updateOrganizer": `{"input":{"data":null,"id":"organizer-1"}}`,
	}, result.Data)
}

func TestNullLiteral_NonNullArgument_IsRejected(t *testing.T) {
	result := executeNullLiteralMutation(t, `
		mutation {
			updateOrganizerNonNull(input: null)
		}
	`, nil)

	require.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Message, "found null")
}

func TestNullLiteral_NonNullInputObjectField_IsRejected(t *testing.T) {
	result := executeNullLiteralMutation(t, `
		mutation {
			updateOrganizer(input: {id: null})
		}
	`, nil)

	require.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Message, "found null")
}
