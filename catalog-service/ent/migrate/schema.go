// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ItemsColumns holds the columns for the "items" table.
	ItemsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "tags", Type: field.TypeJSON, Nullable: true},
		{Name: "image_url", Type: field.TypeString, Nullable: true},
		{Name: "rating", Type: field.TypeFloat64, Default: 0},
		{Name: "review_count", Type: field.TypeInt, Default: 0},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// ItemsTable holds the schema information for the "items" table.
	ItemsTable = &schema.Table{
		Name:       "items",
		Columns:    ItemsColumns,
		PrimaryKey: []*schema.Column{ItemsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "item_title",
				Unique:  false,
				Columns: []*schema.Column{ItemsColumns[1]},
			},
			{
				Name:    "item_tags",
				Unique:  false,
				Columns: []*schema.Column{ItemsColumns[3]},
			},
			{
				Name:    "item_rating",
				Unique:  false,
				Columns: []*schema.Column{ItemsColumns[5]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ItemsTable,
	}
)

func init() {
}
