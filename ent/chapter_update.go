// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/browningluke/mangathr/v2/ent/chapter"
	"github.com/browningluke/mangathr/v2/ent/manga"
	"github.com/browningluke/mangathr/v2/ent/predicate"
)

// ChapterUpdate is the builder for updating Chapter entities.
type ChapterUpdate struct {
	config
	hooks    []Hook
	mutation *ChapterMutation
}

// Where appends a list predicates to the ChapterUpdate builder.
func (cu *ChapterUpdate) Where(ps ...predicate.Chapter) *ChapterUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetChapterID sets the "ChapterID" field.
func (cu *ChapterUpdate) SetChapterID(s string) *ChapterUpdate {
	cu.mutation.SetChapterID(s)
	return cu
}

// SetNillableChapterID sets the "ChapterID" field if the given value is not nil.
func (cu *ChapterUpdate) SetNillableChapterID(s *string) *ChapterUpdate {
	if s != nil {
		cu.SetChapterID(*s)
	}
	return cu
}

// SetNum sets the "Num" field.
func (cu *ChapterUpdate) SetNum(s string) *ChapterUpdate {
	cu.mutation.SetNum(s)
	return cu
}

// SetNillableNum sets the "Num" field if the given value is not nil.
func (cu *ChapterUpdate) SetNillableNum(s *string) *ChapterUpdate {
	if s != nil {
		cu.SetNum(*s)
	}
	return cu
}

// SetTitle sets the "Title" field.
func (cu *ChapterUpdate) SetTitle(s string) *ChapterUpdate {
	cu.mutation.SetTitle(s)
	return cu
}

// SetNillableTitle sets the "Title" field if the given value is not nil.
func (cu *ChapterUpdate) SetNillableTitle(s *string) *ChapterUpdate {
	if s != nil {
		cu.SetTitle(*s)
	}
	return cu
}

// ClearTitle clears the value of the "Title" field.
func (cu *ChapterUpdate) ClearTitle() *ChapterUpdate {
	cu.mutation.ClearTitle()
	return cu
}

// SetCreatedOn sets the "CreatedOn" field.
func (cu *ChapterUpdate) SetCreatedOn(t time.Time) *ChapterUpdate {
	cu.mutation.SetCreatedOn(t)
	return cu
}

// SetNillableCreatedOn sets the "CreatedOn" field if the given value is not nil.
func (cu *ChapterUpdate) SetNillableCreatedOn(t *time.Time) *ChapterUpdate {
	if t != nil {
		cu.SetCreatedOn(*t)
	}
	return cu
}

// ClearCreatedOn clears the value of the "CreatedOn" field.
func (cu *ChapterUpdate) ClearCreatedOn() *ChapterUpdate {
	cu.mutation.ClearCreatedOn()
	return cu
}

// SetRegisteredOn sets the "RegisteredOn" field.
func (cu *ChapterUpdate) SetRegisteredOn(t time.Time) *ChapterUpdate {
	cu.mutation.SetRegisteredOn(t)
	return cu
}

// SetNillableRegisteredOn sets the "RegisteredOn" field if the given value is not nil.
func (cu *ChapterUpdate) SetNillableRegisteredOn(t *time.Time) *ChapterUpdate {
	if t != nil {
		cu.SetRegisteredOn(*t)
	}
	return cu
}

// SetMangaID sets the "Manga" edge to the Manga entity by ID.
func (cu *ChapterUpdate) SetMangaID(id int) *ChapterUpdate {
	cu.mutation.SetMangaID(id)
	return cu
}

// SetNillableMangaID sets the "Manga" edge to the Manga entity by ID if the given value is not nil.
func (cu *ChapterUpdate) SetNillableMangaID(id *int) *ChapterUpdate {
	if id != nil {
		cu = cu.SetMangaID(*id)
	}
	return cu
}

// SetManga sets the "Manga" edge to the Manga entity.
func (cu *ChapterUpdate) SetManga(m *Manga) *ChapterUpdate {
	return cu.SetMangaID(m.ID)
}

// Mutation returns the ChapterMutation object of the builder.
func (cu *ChapterUpdate) Mutation() *ChapterMutation {
	return cu.mutation
}

// ClearManga clears the "Manga" edge to the Manga entity.
func (cu *ChapterUpdate) ClearManga() *ChapterUpdate {
	cu.mutation.ClearManga()
	return cu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *ChapterUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *ChapterUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *ChapterUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *ChapterUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *ChapterUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(chapter.Table, chapter.Columns, sqlgraph.NewFieldSpec(chapter.FieldID, field.TypeInt))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.ChapterID(); ok {
		_spec.SetField(chapter.FieldChapterID, field.TypeString, value)
	}
	if value, ok := cu.mutation.Num(); ok {
		_spec.SetField(chapter.FieldNum, field.TypeString, value)
	}
	if value, ok := cu.mutation.Title(); ok {
		_spec.SetField(chapter.FieldTitle, field.TypeString, value)
	}
	if cu.mutation.TitleCleared() {
		_spec.ClearField(chapter.FieldTitle, field.TypeString)
	}
	if value, ok := cu.mutation.CreatedOn(); ok {
		_spec.SetField(chapter.FieldCreatedOn, field.TypeTime, value)
	}
	if cu.mutation.CreatedOnCleared() {
		_spec.ClearField(chapter.FieldCreatedOn, field.TypeTime)
	}
	if value, ok := cu.mutation.RegisteredOn(); ok {
		_spec.SetField(chapter.FieldRegisteredOn, field.TypeTime, value)
	}
	if cu.mutation.MangaCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   chapter.MangaTable,
			Columns: []string{chapter.MangaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manga.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.MangaIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   chapter.MangaTable,
			Columns: []string{chapter.MangaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manga.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chapter.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// ChapterUpdateOne is the builder for updating a single Chapter entity.
type ChapterUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ChapterMutation
}

// SetChapterID sets the "ChapterID" field.
func (cuo *ChapterUpdateOne) SetChapterID(s string) *ChapterUpdateOne {
	cuo.mutation.SetChapterID(s)
	return cuo
}

// SetNillableChapterID sets the "ChapterID" field if the given value is not nil.
func (cuo *ChapterUpdateOne) SetNillableChapterID(s *string) *ChapterUpdateOne {
	if s != nil {
		cuo.SetChapterID(*s)
	}
	return cuo
}

// SetNum sets the "Num" field.
func (cuo *ChapterUpdateOne) SetNum(s string) *ChapterUpdateOne {
	cuo.mutation.SetNum(s)
	return cuo
}

// SetNillableNum sets the "Num" field if the given value is not nil.
func (cuo *ChapterUpdateOne) SetNillableNum(s *string) *ChapterUpdateOne {
	if s != nil {
		cuo.SetNum(*s)
	}
	return cuo
}

// SetTitle sets the "Title" field.
func (cuo *ChapterUpdateOne) SetTitle(s string) *ChapterUpdateOne {
	cuo.mutation.SetTitle(s)
	return cuo
}

// SetNillableTitle sets the "Title" field if the given value is not nil.
func (cuo *ChapterUpdateOne) SetNillableTitle(s *string) *ChapterUpdateOne {
	if s != nil {
		cuo.SetTitle(*s)
	}
	return cuo
}

// ClearTitle clears the value of the "Title" field.
func (cuo *ChapterUpdateOne) ClearTitle() *ChapterUpdateOne {
	cuo.mutation.ClearTitle()
	return cuo
}

// SetCreatedOn sets the "CreatedOn" field.
func (cuo *ChapterUpdateOne) SetCreatedOn(t time.Time) *ChapterUpdateOne {
	cuo.mutation.SetCreatedOn(t)
	return cuo
}

// SetNillableCreatedOn sets the "CreatedOn" field if the given value is not nil.
func (cuo *ChapterUpdateOne) SetNillableCreatedOn(t *time.Time) *ChapterUpdateOne {
	if t != nil {
		cuo.SetCreatedOn(*t)
	}
	return cuo
}

// ClearCreatedOn clears the value of the "CreatedOn" field.
func (cuo *ChapterUpdateOne) ClearCreatedOn() *ChapterUpdateOne {
	cuo.mutation.ClearCreatedOn()
	return cuo
}

// SetRegisteredOn sets the "RegisteredOn" field.
func (cuo *ChapterUpdateOne) SetRegisteredOn(t time.Time) *ChapterUpdateOne {
	cuo.mutation.SetRegisteredOn(t)
	return cuo
}

// SetNillableRegisteredOn sets the "RegisteredOn" field if the given value is not nil.
func (cuo *ChapterUpdateOne) SetNillableRegisteredOn(t *time.Time) *ChapterUpdateOne {
	if t != nil {
		cuo.SetRegisteredOn(*t)
	}
	return cuo
}

// SetMangaID sets the "Manga" edge to the Manga entity by ID.
func (cuo *ChapterUpdateOne) SetMangaID(id int) *ChapterUpdateOne {
	cuo.mutation.SetMangaID(id)
	return cuo
}

// SetNillableMangaID sets the "Manga" edge to the Manga entity by ID if the given value is not nil.
func (cuo *ChapterUpdateOne) SetNillableMangaID(id *int) *ChapterUpdateOne {
	if id != nil {
		cuo = cuo.SetMangaID(*id)
	}
	return cuo
}

// SetManga sets the "Manga" edge to the Manga entity.
func (cuo *ChapterUpdateOne) SetManga(m *Manga) *ChapterUpdateOne {
	return cuo.SetMangaID(m.ID)
}

// Mutation returns the ChapterMutation object of the builder.
func (cuo *ChapterUpdateOne) Mutation() *ChapterMutation {
	return cuo.mutation
}

// ClearManga clears the "Manga" edge to the Manga entity.
func (cuo *ChapterUpdateOne) ClearManga() *ChapterUpdateOne {
	cuo.mutation.ClearManga()
	return cuo
}

// Where appends a list predicates to the ChapterUpdate builder.
func (cuo *ChapterUpdateOne) Where(ps ...predicate.Chapter) *ChapterUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *ChapterUpdateOne) Select(field string, fields ...string) *ChapterUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Chapter entity.
func (cuo *ChapterUpdateOne) Save(ctx context.Context) (*Chapter, error) {
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *ChapterUpdateOne) SaveX(ctx context.Context) *Chapter {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *ChapterUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *ChapterUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *ChapterUpdateOne) sqlSave(ctx context.Context) (_node *Chapter, err error) {
	_spec := sqlgraph.NewUpdateSpec(chapter.Table, chapter.Columns, sqlgraph.NewFieldSpec(chapter.FieldID, field.TypeInt))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Chapter.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, chapter.FieldID)
		for _, f := range fields {
			if !chapter.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != chapter.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.ChapterID(); ok {
		_spec.SetField(chapter.FieldChapterID, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Num(); ok {
		_spec.SetField(chapter.FieldNum, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Title(); ok {
		_spec.SetField(chapter.FieldTitle, field.TypeString, value)
	}
	if cuo.mutation.TitleCleared() {
		_spec.ClearField(chapter.FieldTitle, field.TypeString)
	}
	if value, ok := cuo.mutation.CreatedOn(); ok {
		_spec.SetField(chapter.FieldCreatedOn, field.TypeTime, value)
	}
	if cuo.mutation.CreatedOnCleared() {
		_spec.ClearField(chapter.FieldCreatedOn, field.TypeTime)
	}
	if value, ok := cuo.mutation.RegisteredOn(); ok {
		_spec.SetField(chapter.FieldRegisteredOn, field.TypeTime, value)
	}
	if cuo.mutation.MangaCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   chapter.MangaTable,
			Columns: []string{chapter.MangaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manga.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.MangaIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   chapter.MangaTable,
			Columns: []string{chapter.MangaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manga.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Chapter{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chapter.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
