package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return uuid.New().String()
			}),
		field.String("title").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.Strings("tags").
			Optional(),
		field.String("image_url").
			Optional(),
		field.Float("rating").
			Default(0),
		field.Int("review_count").
			Default(0),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return nil
}

func (Item) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title"),
		index.Fields("tags"),
		index.Fields("rating"),
	}
}
