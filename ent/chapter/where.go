// Code generated by ent, DO NOT EDIT.

package chapter

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/browningluke/mangathr/v2/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Chapter {
	return predicate.Chapter(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Chapter {
	return predicate.Chapter(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Chapter {
	return predicate.Chapter(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Chapter {
	return predicate.Chapter(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Chapter {
	return predicate.Chapter(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Chapter {
	return predicate.Chapter(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Chapter {
	return predicate.Chapter(sql.FieldLTE(FieldID, id))
}

// ChapterID applies equality check predicate on the "ChapterID" field. It's identical to ChapterIDEQ.
func ChapterID(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldChapterID, v))
}

// Num applies equality check predicate on the "Num" field. It's identical to NumEQ.
func Num(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldNum, v))
}

// Title applies equality check predicate on the "Title" field. It's identical to TitleEQ.
func Title(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldTitle, v))
}

// CreatedOn applies equality check predicate on the "CreatedOn" field. It's identical to CreatedOnEQ.
func CreatedOn(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldCreatedOn, v))
}

// RegisteredOn applies equality check predicate on the "RegisteredOn" field. It's identical to RegisteredOnEQ.
func RegisteredOn(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldRegisteredOn, v))
}

// ChapterIDEQ applies the EQ predicate on the "ChapterID" field.
func ChapterIDEQ(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldChapterID, v))
}

// ChapterIDNEQ applies the NEQ predicate on the "ChapterID" field.
func ChapterIDNEQ(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldNEQ(FieldChapterID, v))
}

// ChapterIDIn applies the In predicate on the "ChapterID" field.
func ChapterIDIn(vs ...string) predicate.Chapter {
	return predicate.Chapter(sql.FieldIn(FieldChapterID, vs...))
}

// ChapterIDNotIn applies the NotIn predicate on the "ChapterID" field.
func ChapterIDNotIn(vs ...string) predicate.Chapter {
	return predicate.Chapter(sql.FieldNotIn(FieldChapterID, vs...))
}

// ChapterIDGT applies the GT predicate on the "ChapterID" field.
func ChapterIDGT(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldGT(FieldChapterID, v))
}

// ChapterIDGTE applies the GTE predicate on the "ChapterID" field.
func ChapterIDGTE(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldGTE(FieldChapterID, v))
}

// ChapterIDLT applies the LT predicate on the "ChapterID" field.
func ChapterIDLT(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldLT(FieldChapterID, v))
}

// ChapterIDLTE applies the LTE predicate on the "ChapterID" field.
func ChapterIDLTE(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldLTE(FieldChapterID, v))
}

// ChapterIDContains applies the Contains predicate on the "ChapterID" field.
func ChapterIDContains(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldContains(FieldChapterID, v))
}

// ChapterIDHasPrefix applies the HasPrefix predicate on the "ChapterID" field.
func ChapterIDHasPrefix(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldHasPrefix(FieldChapterID, v))
}

// ChapterIDHasSuffix applies the HasSuffix predicate on the "ChapterID" field.
func ChapterIDHasSuffix(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldHasSuffix(FieldChapterID, v))
}

// ChapterIDEqualFold applies the EqualFold predicate on the "ChapterID" field.
func ChapterIDEqualFold(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEqualFold(FieldChapterID, v))
}

// ChapterIDContainsFold applies the ContainsFold predicate on the "ChapterID" field.
func ChapterIDContainsFold(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldContainsFold(FieldChapterID, v))
}

// NumEQ applies the EQ predicate on the "Num" field.
func NumEQ(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldNum, v))
}

// NumNEQ applies the NEQ predicate on the "Num" field.
func NumNEQ(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldNEQ(FieldNum, v))
}

// NumIn applies the In predicate on the "Num" field.
func NumIn(vs ...string) predicate.Chapter {
	return predicate.Chapter(sql.FieldIn(FieldNum, vs...))
}

// NumNotIn applies the NotIn predicate on the "Num" field.
func NumNotIn(vs ...string) predicate.Chapter {
	return predicate.Chapter(sql.FieldNotIn(FieldNum, vs...))
}

// NumGT applies the GT predicate on the "Num" field.
func NumGT(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldGT(FieldNum, v))
}

// NumGTE applies the GTE predicate on the "Num" field.
func NumGTE(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldGTE(FieldNum, v))
}

// NumLT applies the LT predicate on the "Num" field.
func NumLT(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldLT(FieldNum, v))
}

// NumLTE applies the LTE predicate on the "Num" field.
func NumLTE(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldLTE(FieldNum, v))
}

// NumContains applies the Contains predicate on the "Num" field.
func NumContains(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldContains(FieldNum, v))
}

// NumHasPrefix applies the HasPrefix predicate on the "Num" field.
func NumHasPrefix(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldHasPrefix(FieldNum, v))
}

// NumHasSuffix applies the HasSuffix predicate on the "Num" field.
func NumHasSuffix(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldHasSuffix(FieldNum, v))
}

// NumEqualFold applies the EqualFold predicate on the "Num" field.
func NumEqualFold(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEqualFold(FieldNum, v))
}

// NumContainsFold applies the ContainsFold predicate on the "Num" field.
func NumContainsFold(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldContainsFold(FieldNum, v))
}

// TitleEQ applies the EQ predicate on the "Title" field.
func TitleEQ(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldTitle, v))
}

// TitleNEQ applies the NEQ predicate on the "Title" field.
func TitleNEQ(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldNEQ(FieldTitle, v))
}

// TitleIn applies the In predicate on the "Title" field.
func TitleIn(vs ...string) predicate.Chapter {
	return predicate.Chapter(sql.FieldIn(FieldTitle, vs...))
}

// TitleNotIn applies the NotIn predicate on the "Title" field.
func TitleNotIn(vs ...string) predicate.Chapter {
	return predicate.Chapter(sql.FieldNotIn(FieldTitle, vs...))
}

// TitleGT applies the GT predicate on the "Title" field.
func TitleGT(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldGT(FieldTitle, v))
}

// TitleGTE applies the GTE predicate on the "Title" field.
func TitleGTE(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldGTE(FieldTitle, v))
}

// TitleLT applies the LT predicate on the "Title" field.
func TitleLT(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldLT(FieldTitle, v))
}

// TitleLTE applies the LTE predicate on the "Title" field.
func TitleLTE(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldLTE(FieldTitle, v))
}

// TitleContains applies the Contains predicate on the "Title" field.
func TitleContains(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldContains(FieldTitle, v))
}

// TitleHasPrefix applies the HasPrefix predicate on the "Title" field.
func TitleHasPrefix(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldHasPrefix(FieldTitle, v))
}

// TitleHasSuffix applies the HasSuffix predicate on the "Title" field.
func TitleHasSuffix(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldHasSuffix(FieldTitle, v))
}

// TitleIsNil applies the IsNil predicate on the "Title" field.
func TitleIsNil() predicate.Chapter {
	return predicate.Chapter(sql.FieldIsNull(FieldTitle))
}

// TitleNotNil applies the NotNil predicate on the "Title" field.
func TitleNotNil() predicate.Chapter {
	return predicate.Chapter(sql.FieldNotNull(FieldTitle))
}

// TitleEqualFold applies the EqualFold predicate on the "Title" field.
func TitleEqualFold(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldEqualFold(FieldTitle, v))
}

// TitleContainsFold applies the ContainsFold predicate on the "Title" field.
func TitleContainsFold(v string) predicate.Chapter {
	return predicate.Chapter(sql.FieldContainsFold(FieldTitle, v))
}

// CreatedOnEQ applies the EQ predicate on the "CreatedOn" field.
func CreatedOnEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldCreatedOn, v))
}

// CreatedOnNEQ applies the NEQ predicate on the "CreatedOn" field.
func CreatedOnNEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldNEQ(FieldCreatedOn, v))
}

// CreatedOnIn applies the In predicate on the "CreatedOn" field.
func CreatedOnIn(vs ...time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldIn(FieldCreatedOn, vs...))
}

// CreatedOnNotIn applies the NotIn predicate on the "CreatedOn" field.
func CreatedOnNotIn(vs ...time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldNotIn(FieldCreatedOn, vs...))
}

// CreatedOnGT applies the GT predicate on the "CreatedOn" field.
func CreatedOnGT(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldGT(FieldCreatedOn, v))
}

// CreatedOnGTE applies the GTE predicate on the "CreatedOn" field.
func CreatedOnGTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldGTE(FieldCreatedOn, v))
}

// CreatedOnLT applies the LT predicate on the "CreatedOn" field.
func CreatedOnLT(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldLT(FieldCreatedOn, v))
}

// CreatedOnLTE applies the LTE predicate on the "CreatedOn" field.
func CreatedOnLTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldLTE(FieldCreatedOn, v))
}

// CreatedOnIsNil applies the IsNil predicate on the "CreatedOn" field.
func CreatedOnIsNil() predicate.Chapter {
	return predicate.Chapter(sql.FieldIsNull(FieldCreatedOn))
}

// CreatedOnNotNil applies the NotNil predicate on the "CreatedOn" field.
func CreatedOnNotNil() predicate.Chapter {
	return predicate.Chapter(sql.FieldNotNull(FieldCreatedOn))
}

// RegisteredOnEQ applies the EQ predicate on the "RegisteredOn" field.
func RegisteredOnEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldEQ(FieldRegisteredOn, v))
}

// RegisteredOnNEQ applies the NEQ predicate on the "RegisteredOn" field.
func RegisteredOnNEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldNEQ(FieldRegisteredOn, v))
}

// RegisteredOnIn applies the In predicate on the "RegisteredOn" field.
func RegisteredOnIn(vs ...time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldIn(FieldRegisteredOn, vs...))
}

// RegisteredOnNotIn applies the NotIn predicate on the "RegisteredOn" field.
func RegisteredOnNotIn(vs ...time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldNotIn(FieldRegisteredOn, vs...))
}

// RegisteredOnGT applies the GT predicate on the "RegisteredOn" field.
func RegisteredOnGT(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldGT(FieldRegisteredOn, v))
}

// RegisteredOnGTE applies the GTE predicate on the "RegisteredOn" field.
func RegisteredOnGTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldGTE(FieldRegisteredOn, v))
}

// RegisteredOnLT applies the LT predicate on the "RegisteredOn" field.
func RegisteredOnLT(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldLT(FieldRegisteredOn, v))
}

// RegisteredOnLTE applies the LTE predicate on the "RegisteredOn" field.
func RegisteredOnLTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(sql.FieldLTE(FieldRegisteredOn, v))
}

// HasManga applies the HasEdge predicate on the "Manga" edge.
func HasManga() predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, MangaTable, MangaColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMangaWith applies the HasEdge predicate on the "Manga" edge with a given conditions (other predicates).
func HasMangaWith(preds ...predicate.Manga) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		step := newMangaStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Chapter) predicate.Chapter {
	return predicate.Chapter(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Chapter) predicate.Chapter {
	return predicate.Chapter(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Chapter) predicate.Chapter {
	return predicate.Chapter(sql.NotPredicates(p))
}
