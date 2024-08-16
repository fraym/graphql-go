# graphql [![Go Reference](https://pkg.go.dev/badge/github.com/fraym/graphql-go.svg)](https://pkg.go.dev/github.com/fraym/graphql-go)

An implementation of GraphQL in Go. Follows the official reference implementation [`graphql-js`](https://github.com/graphql/graphql-js).

Supports: queries, mutations & subscriptions.

## Fork

This is a fork of [`graphql-go/graphql`](https://github.com/graphql-go/graphql).
This fork contains a set of changes:

- it requires go 1.22
- `null` value support
- arguments in introspection queries are statically ordered (their order won't change on every request)
- numbers are limited to the JavaScript number space
- serialisation to `nil` and serialisation errors for scalar values

### Documentation

godoc: https://pkg.go.dev/github.com/fraym/graphql-go

### Getting Started

To install the library, run:

```bash
go get github.com/fraym/graphql-go
```

The following is a simple example which defines a schema with a single `hello` string-type field and a `Resolve` method which returns the string `world`. A GraphQL query is performed against this schema with the resulting output printed in JSON format.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fraym/graphql-go"
)

func main() {
	// Schema
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (any, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			hello
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
}
```

For more complex examples, refer to the [examples/](https://github.com/fraym/graphql-go/tree/master/examples/) directory and [graphql_test.go](https://github.com/fraym/graphql-go/blob/master/graphql_test.go).

### Third Party Libraries

|                                     Name                                      |                     Author                      |                                 Description                                  |
| :---------------------------------------------------------------------------: | :---------------------------------------------: | :--------------------------------------------------------------------------: |
|     [graphql-go-handler](https://github.com/fraym/graphql-go-go-handler)      |    [Hafiz Ismail](https://github.com/sogko)     |         Middleware to handle GraphQL queries through HTTP requests.          |
|       [graphql-relay-go](https://github.com/fraym/graphql-go-relay-go)        |    [Hafiz Ismail](https://github.com/sogko)     |         Lib to construct a graphql-go server supporting react-relay.         |
| [golang-relay-starter-kit](https://github.com/sogko/golang-relay-starter-kit) |    [Hafiz Ismail](https://github.com/sogko)     | Barebones starting point for a Relay application with Golang GraphQL server. |
|           [dataloader](https://github.com/nicksrandall/dataloader)            | [Nick Randall](https://github.com/nicksrandall) |  [DataLoader](https://github.com/facebook/dataloader) implementation in Go.  |

### Blog Posts

- [Golang + GraphQL + Relay](https://wehavefaces.net/learn-golang-graphql-relay-1-e59ea174a902)
