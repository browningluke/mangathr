// Code generated by ent, DO NOT EDIT.

package chapter

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the chapter type in the database.
	Label = "chapter"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldChapterID holds the string denoting the chapterid field in the database.
	FieldChapterID = "chapter_id"
	// FieldNum holds the string denoting the num field in the database.
	FieldNum = "num"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldCreatedOn holds the string denoting the createdon field in the database.
	FieldCreatedOn = "created_on"
	// FieldRegisteredOn holds the string denoting the registeredon field in the database.
	FieldRegisteredOn = "registered_on"
	// EdgeManga holds the string denoting the manga edge name in mutations.
	EdgeManga = "Manga"
	// Table holds the table name of the chapter in the database.
	Table = "chapters"
	// MangaTable is the table that holds the Manga relation/edge.
	MangaTable = "chapters"
	// MangaInverseTable is the table name for the Manga entity.
	// It exists in this package in order to avoid circular dependency with the "manga" package.
	MangaInverseTable = "mangas"
	// MangaColumn is the table column denoting the Manga relation/edge.
	MangaColumn = "manga_chapters"
)

// Columns holds all SQL columns for chapter fields.
var Columns = []string{
	FieldID,
	FieldChapterID,
	FieldNum,
	FieldTitle,
	FieldCreatedOn,
	FieldRegisteredOn,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "chapters"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"manga_chapters",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Chapter queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByChapterID orders the results by the ChapterID field.
func ByChapterID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChapterID, opts...).ToFunc()
}

// ByNum orders the results by the Num field.
func ByNum(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNum, opts...).ToFunc()
}

// ByTitle orders the results by the Title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// ByCreatedOn orders the results by the CreatedOn field.
func ByCreatedOn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedOn, opts...).ToFunc()
}

// ByRegisteredOn orders the results by the RegisteredOn field.
func ByRegisteredOn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRegisteredOn, opts...).ToFunc()
}

// ByMangaField orders the results by Manga field.
func ByMangaField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMangaStep(), sql.OrderByField(field, opts...))
	}
}
func newMangaStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MangaInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, MangaTable, MangaColumn),
	)
}
