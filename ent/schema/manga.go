package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Manga holds the schema definition for the Manga entity.
type Manga struct {
	ent.Schema
}

// Fields of the Manga.
func (Manga) Fields() []ent.Field {
	return []ent.Field{
		field.String("MangaID"),
		field.String("Source"),
		field.String("Title"),
		field.String("Mapping"),
		field.Time("RegisteredOn"),
	}
}

// Edges of the Manga.
func (Manga) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Chapters", Chapter.Type),
	}
}
