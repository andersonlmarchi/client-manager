package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type BrandingSettings struct {
	ent.Schema
}

func (BrandingSettings) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("app_name").NotEmpty(),
		field.String("logo_url").Optional().Nillable(),
		field.String("primary_color").Optional().Nillable(),
		field.String("support_email").Optional().Nillable(),
		field.String("smtp_host_ref").Optional().Nillable(),
		field.String("smtp_port_ref").Optional().Nillable(),
		field.String("smtp_user_ref").Optional().Nillable(),
		field.String("from_email").Optional().Nillable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
