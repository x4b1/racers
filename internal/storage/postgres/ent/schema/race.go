package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/facebookincubator/ent-contrib/entgql"
)

// Race holds the schema definition for the Race entity.
type Race struct {
	ent.Schema
}

// Fields of the Race.
func (Race) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.Text("name").
			NotEmpty(),
		field.Time("date"),
	}
}

// Edges of the Race.
func (Race) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("competitors", User.Type).
			Annotations(entgql.Bind()),
	}
}
