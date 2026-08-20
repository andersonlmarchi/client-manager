package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Credential struct {
	ent.Schema
}

func (Credential) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("password_hash").NotEmpty().Sensitive(),
		field.String("algorithm").Default("argon2id").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Credential) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("credential").Unique().Required(),
	}
}
