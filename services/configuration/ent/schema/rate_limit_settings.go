package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type RateLimitSettings struct {
	ent.Schema
}

func (RateLimitSettings) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.Int("default_rps").Default(60).Positive(),
		field.Int("default_burst").Default(120).Positive(),
		field.Int("login_rps").Default(5).Positive(),
		field.Int("login_burst").Default(10).Positive(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
