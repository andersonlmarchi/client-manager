package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type OIDCSettings struct {
	ent.Schema
}

func (OIDCSettings) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("issuer").NotEmpty(),
		field.Int("access_token_ttl_seconds").Default(900).Positive(),
		field.Int("refresh_token_ttl_seconds").Default(2592000).Positive(),
		field.Int("id_token_ttl_seconds").Default(900).Positive(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
