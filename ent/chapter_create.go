// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/browningluke/mangathr/v2/ent/chapter"
	"github.com/browningluke/mangathr/v2/ent/manga"
)

// ChapterCreate is the builder for creating a Chapter entity.
type ChapterCreate struct {
	config
	mutation *ChapterMutation
	hooks    []Hook
}

// SetChapterID sets the "ChapterID" field.
func (cc *ChapterCreate) SetChapterID(s string) *ChapterCreate {
	cc.mutation.SetChapterID(s)
	return cc
}

// SetNum sets the "Num" field.
func (cc *ChapterCreate) SetNum(s string) *ChapterCreate {
	cc.mutation.SetNum(s)
	return cc
}

// SetTitle sets the "Title" field.
func (cc *ChapterCreate) SetTitle(s string) *ChapterCreate {
	cc.mutation.SetTitle(s)
	return cc
}

// SetNillableTitle sets the "Title" field if the given value is not nil.
func (cc *ChapterCreate) SetNillableTitle(s *string) *ChapterCreate {
	if s != nil {
		cc.SetTitle(*s)
	}
	return cc
}

// SetCreatedOn sets the "CreatedOn" field.
func (cc *ChapterCreate) SetCreatedOn(t time.Time) *ChapterCreate {
	cc.mutation.SetCreatedOn(t)
	return cc
}

// SetNillableCreatedOn sets the "CreatedOn" field if the given value is not nil.
func (cc *ChapterCreate) SetNillableCreatedOn(t *time.Time) *ChapterCreate {
	if t != nil {
		cc.SetCreatedOn(*t)
	}
	return cc
}

// SetRegisteredOn sets the "RegisteredOn" field.
func (cc *ChapterCreate) SetRegisteredOn(t time.Time) *ChapterCreate {
	cc.mutation.SetRegisteredOn(t)
	return cc
}

// SetMangaID sets the "Manga" edge to the Manga entity by ID.
func (cc *ChapterCreate) SetMangaID(id int) *ChapterCreate {
	cc.mutation.SetMangaID(id)
	return cc
}

// SetNillableMangaID sets the "Manga" edge to the Manga entity by ID if the given value is not nil.
func (cc *ChapterCreate) SetNillableMangaID(id *int) *ChapterCreate {
	if id != nil {
		cc = cc.SetMangaID(*id)
	}
	return cc
}

// SetManga sets the "Manga" edge to the Manga entity.
func (cc *ChapterCreate) SetManga(m *Manga) *ChapterCreate {
	return cc.SetMangaID(m.ID)
}

// Mutation returns the ChapterMutation object of the builder.
func (cc *ChapterCreate) Mutation() *ChapterMutation {
	return cc.mutation
}

// Save creates the Chapter in the database.
func (cc *ChapterCreate) Save(ctx context.Context) (*Chapter, error) {
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *ChapterCreate) SaveX(ctx context.Context) *Chapter {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *ChapterCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *ChapterCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *ChapterCreate) check() error {
	if _, ok := cc.mutation.ChapterID(); !ok {
		return &ValidationError{Name: "ChapterID", err: errors.New(`ent: missing required field "Chapter.ChapterID"`)}
	}
	if _, ok := cc.mutation.Num(); !ok {
		return &ValidationError{Name: "Num", err: errors.New(`ent: missing required field "Chapter.Num"`)}
	}
	if _, ok := cc.mutation.RegisteredOn(); !ok {
		return &ValidationError{Name: "RegisteredOn", err: errors.New(`ent: missing required field "Chapter.RegisteredOn"`)}
	}
	return nil
}

func (cc *ChapterCreate) sqlSave(ctx context.Context) (*Chapter, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *ChapterCreate) createSpec() (*Chapter, *sqlgraph.CreateSpec) {
	var (
		_node = &Chapter{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(chapter.Table, sqlgraph.NewFieldSpec(chapter.FieldID, field.TypeInt))
	)
	if value, ok := cc.mutation.ChapterID(); ok {
		_spec.SetField(chapter.FieldChapterID, field.TypeString, value)
		_node.ChapterID = value
	}
	if value, ok := cc.mutation.Num(); ok {
		_spec.SetField(chapter.FieldNum, field.TypeString, value)
		_node.Num = value
	}
	if value, ok := cc.mutation.Title(); ok {
		_spec.SetField(chapter.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := cc.mutation.CreatedOn(); ok {
		_spec.SetField(chapter.FieldCreatedOn, field.TypeTime, value)
		_node.CreatedOn = value
	}
	if value, ok := cc.mutation.RegisteredOn(); ok {
		_spec.SetField(chapter.FieldRegisteredOn, field.TypeTime, value)
		_node.RegisteredOn = value
	}
	if nodes := cc.mutation.MangaIDs(); len(nodes) > 0 {
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
		_node.manga_chapters = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ChapterCreateBulk is the builder for creating many Chapter entities in bulk.
type ChapterCreateBulk struct {
	config
	err      error
	builders []*ChapterCreate
}

// Save creates the Chapter entities in the database.
func (ccb *ChapterCreateBulk) Save(ctx context.Context) ([]*Chapter, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Chapter, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ChapterMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *ChapterCreateBulk) SaveX(ctx context.Context) []*Chapter {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *ChapterCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *ChapterCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
