package graphql_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRace(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "race")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	filename := filepath.Join(tempdir, "example.go")
	err = os.WriteFile(filename, []byte(`
		package main

		import (
			"runtime"
			"sync"

			"github.com/graphql-go/graphql"
		)

		func main() {
			var wg sync.WaitGroup
			wg.Add(2)
			for i := 0; i < 2; i++ {
				go func() {
					defer wg.Done()
					schema, _ := graphql.NewSchema(graphql.SchemaConfig{
						Query: graphql.NewObject(graphql.ObjectConfig{
							Name: "RootQuery",
							Fields: graphql.Fields{
								"hello": &graphql.Field{
									Type: graphql.String,
									Resolve: func(p graphql.ResolveParams) (any, error) {
										return "world", nil
									},
								},
							},
						}),
					})
					runtime.KeepAlive(schema)
				}()
			}

			wg.Wait()
		} 
	`), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	result, err := exec.Command("go", "run", "-race", filename).CombinedOutput()
	if err != nil || len(result) != 0 {
		t.Log(string(result))
		t.Fatal(err)
	}
}
