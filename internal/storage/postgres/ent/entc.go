// +build ignore

package main

import (
	"log"

	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/facebook/ent/schema/field"
	"github.com/facebookincubator/ent-contrib/entgql"
)

func main() {
	err := entc.Generate("./schema", &gen.Config{
		Templates: entgql.AllTemplates,
		IDType:    &field.TypeInfo{Type: field.TypeString},
	})
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
