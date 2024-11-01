// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ChaptersColumns holds the columns for the "chapters" table.
	ChaptersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "chapter_id", Type: field.TypeString},
		{Name: "num", Type: field.TypeString},
		{Name: "title", Type: field.TypeString, Nullable: true},
		{Name: "created_on", Type: field.TypeTime, Nullable: true},
		{Name: "registered_on", Type: field.TypeTime},
		{Name: "manga_chapters", Type: field.TypeInt, Nullable: true},
	}
	// ChaptersTable holds the schema information for the "chapters" table.
	ChaptersTable = &schema.Table{
		Name:       "chapters",
		Columns:    ChaptersColumns,
		PrimaryKey: []*schema.Column{ChaptersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "chapters_mangas_Chapters",
				Columns:    []*schema.Column{ChaptersColumns[6]},
				RefColumns: []*schema.Column{MangasColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// MangasColumns holds the columns for the "mangas" table.
	MangasColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "manga_id", Type: field.TypeString},
		{Name: "source", Type: field.TypeString},
		{Name: "title", Type: field.TypeString},
		{Name: "mapping", Type: field.TypeString},
		{Name: "registered_on", Type: field.TypeTime},
		{Name: "filtered_groups", Type: field.TypeJSON, Nullable: true},
		{Name: "excluded_groups", Type: field.TypeJSON, Nullable: true},
	}
	// MangasTable holds the schema information for the "mangas" table.
	MangasTable = &schema.Table{
		Name:       "mangas",
		Columns:    MangasColumns,
		PrimaryKey: []*schema.Column{MangasColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ChaptersTable,
		MangasTable,
	}
)

func init() {
	ChaptersTable.ForeignKeys[0].RefTable = MangasTable
}
