// Code generated by ent, DO NOT EDIT.

package item

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the item type in the database.
	Label = "item"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldTags holds the string denoting the tags field in the database.
	FieldTags = "tags"
	// FieldImageURL holds the string denoting the image_url field in the database.
	FieldImageURL = "image_url"
	// FieldRating holds the string denoting the rating field in the database.
	FieldRating = "rating"
	// FieldReviewCount holds the string denoting the review_count field in the database.
	FieldReviewCount = "review_count"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the item in the database.
	Table = "items"
)

// Columns holds all SQL columns for item fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldDescription,
	FieldTags,
	FieldImageURL,
	FieldRating,
	FieldReviewCount,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// TitleValidator is a validator for the "title" field. It is called by the builders before save.
	TitleValidator func(string) error
	// DefaultRating holds the default value on creation for the "rating" field.
	DefaultRating float64
	// DefaultReviewCount holds the default value on creation for the "review_count" field.
	DefaultReviewCount int
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() string
)

// OrderOption defines the ordering options for the Item queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByImageURL orders the results by the image_url field.
func ByImageURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldImageURL, opts...).ToFunc()
}

// ByRating orders the results by the rating field.
func ByRating(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRating, opts...).ToFunc()
}

// ByReviewCount orders the results by the review_count field.
func ByReviewCount(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReviewCount, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}
