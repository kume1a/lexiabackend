package schema

import (
	"lexia/internal/shared"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Folder struct {
	ent.Schema
}

func (Folder) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			NotEmpty(),
		field.Int32("wordCount"),
		field.Enum("languageFrom").GoType(shared.Language("")),
	}
}

func (Folder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("folders").
			Unique(),
		edge.To("words", Word.Type),
	}
}

func (Folder) Indexes() []ent.Index {
	return nil
}

func (Folder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
