package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type PasswordPolicy struct {
	ent.Schema
}

func (PasswordPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.Int("min_length").Default(12).Positive(),
		field.Bool("require_upper").Default(true),
		field.Bool("require_lower").Default(true),
		field.Bool("require_number").Default(true),
		field.Bool("require_symbol").Default(true),
		field.Int("max_age_days").Default(0).NonNegative(),
		field.Int("history_count").Default(0).NonNegative(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
