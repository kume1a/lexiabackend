package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Word struct {
	ent.Schema
}

func (Word) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("text").
			NotEmpty(),
		field.String("definition"),
	}
}

func (Word) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("folder", Folder.Type).
			Ref("words").
			Unique(),
	}
}

func (Word) Indexes() []ent.Index {
	return nil
}

func (Word) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
