package graphql

import (
	"github.com/graphql-go/graphql/gqlerrors"
)

// type Schema any

// Result has the response, errors and extensions from the resolved schema
type Result struct {
	Data       any                        `json:"data"`
	Errors     []gqlerrors.FormattedError `json:"errors,omitempty"`
	Extensions map[string]any             `json:"extensions,omitempty"`
}

// HasErrors just a simple function to help you decide if the result has errors or not
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}
