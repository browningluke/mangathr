// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/browningluke/mangathrV2/ent/chapter"
	"github.com/browningluke/mangathrV2/ent/manga"
	"github.com/browningluke/mangathrV2/ent/predicate"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeChapter = "Chapter"
	TypeManga   = "Manga"
)

// ChapterMutation represents an operation that mutates the Chapter nodes in the graph.
type ChapterMutation struct {
	config
	op            Op
	typ           string
	id            *int
	_ChapterID    *string
	_Num          *string
	_Title        *string
	_CreatedOn    *time.Time
	_RegisteredOn *time.Time
	clearedFields map[string]struct{}
	_Manga        *int
	cleared_Manga bool
	done          bool
	oldValue      func(context.Context) (*Chapter, error)
	predicates    []predicate.Chapter
}

var _ ent.Mutation = (*ChapterMutation)(nil)

// chapterOption allows management of the mutation configuration using functional options.
type chapterOption func(*ChapterMutation)

// newChapterMutation creates new mutation for the Chapter entity.
func newChapterMutation(c config, op Op, opts ...chapterOption) *ChapterMutation {
	m := &ChapterMutation{
		config:        c,
		op:            op,
		typ:           TypeChapter,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withChapterID sets the ID field of the mutation.
func withChapterID(id int) chapterOption {
	return func(m *ChapterMutation) {
		var (
			err   error
			once  sync.Once
			value *Chapter
		)
		m.oldValue = func(ctx context.Context) (*Chapter, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Chapter.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withChapter sets the old Chapter of the mutation.
func withChapter(node *Chapter) chapterOption {
	return func(m *ChapterMutation) {
		m.oldValue = func(context.Context) (*Chapter, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ChapterMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ChapterMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ChapterMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ChapterMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Chapter.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetChapterID sets the "ChapterID" field.
func (m *ChapterMutation) SetChapterID(s string) {
	m._ChapterID = &s
}

// ChapterID returns the value of the "ChapterID" field in the mutation.
func (m *ChapterMutation) ChapterID() (r string, exists bool) {
	v := m._ChapterID
	if v == nil {
		return
	}
	return *v, true
}

// OldChapterID returns the old "ChapterID" field's value of the Chapter entity.
// If the Chapter object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChapterMutation) OldChapterID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldChapterID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldChapterID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldChapterID: %w", err)
	}
	return oldValue.ChapterID, nil
}

// ResetChapterID resets all changes to the "ChapterID" field.
func (m *ChapterMutation) ResetChapterID() {
	m._ChapterID = nil
}

// SetNum sets the "Num" field.
func (m *ChapterMutation) SetNum(s string) {
	m._Num = &s
}

// Num returns the value of the "Num" field in the mutation.
func (m *ChapterMutation) Num() (r string, exists bool) {
	v := m._Num
	if v == nil {
		return
	}
	return *v, true
}

// OldNum returns the old "Num" field's value of the Chapter entity.
// If the Chapter object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChapterMutation) OldNum(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldNum is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldNum requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldNum: %w", err)
	}
	return oldValue.Num, nil
}

// ResetNum resets all changes to the "Num" field.
func (m *ChapterMutation) ResetNum() {
	m._Num = nil
}

// SetTitle sets the "Title" field.
func (m *ChapterMutation) SetTitle(s string) {
	m._Title = &s
}

// Title returns the value of the "Title" field in the mutation.
func (m *ChapterMutation) Title() (r string, exists bool) {
	v := m._Title
	if v == nil {
		return
	}
	return *v, true
}

// OldTitle returns the old "Title" field's value of the Chapter entity.
// If the Chapter object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChapterMutation) OldTitle(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTitle is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTitle requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTitle: %w", err)
	}
	return oldValue.Title, nil
}

// ClearTitle clears the value of the "Title" field.
func (m *ChapterMutation) ClearTitle() {
	m._Title = nil
	m.clearedFields[chapter.FieldTitle] = struct{}{}
}

// TitleCleared returns if the "Title" field was cleared in this mutation.
func (m *ChapterMutation) TitleCleared() bool {
	_, ok := m.clearedFields[chapter.FieldTitle]
	return ok
}

// ResetTitle resets all changes to the "Title" field.
func (m *ChapterMutation) ResetTitle() {
	m._Title = nil
	delete(m.clearedFields, chapter.FieldTitle)
}

// SetCreatedOn sets the "CreatedOn" field.
func (m *ChapterMutation) SetCreatedOn(t time.Time) {
	m._CreatedOn = &t
}

// CreatedOn returns the value of the "CreatedOn" field in the mutation.
func (m *ChapterMutation) CreatedOn() (r time.Time, exists bool) {
	v := m._CreatedOn
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedOn returns the old "CreatedOn" field's value of the Chapter entity.
// If the Chapter object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChapterMutation) OldCreatedOn(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedOn is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedOn requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedOn: %w", err)
	}
	return oldValue.CreatedOn, nil
}

// ClearCreatedOn clears the value of the "CreatedOn" field.
func (m *ChapterMutation) ClearCreatedOn() {
	m._CreatedOn = nil
	m.clearedFields[chapter.FieldCreatedOn] = struct{}{}
}

// CreatedOnCleared returns if the "CreatedOn" field was cleared in this mutation.
func (m *ChapterMutation) CreatedOnCleared() bool {
	_, ok := m.clearedFields[chapter.FieldCreatedOn]
	return ok
}

// ResetCreatedOn resets all changes to the "CreatedOn" field.
func (m *ChapterMutation) ResetCreatedOn() {
	m._CreatedOn = nil
	delete(m.clearedFields, chapter.FieldCreatedOn)
}

// SetRegisteredOn sets the "RegisteredOn" field.
func (m *ChapterMutation) SetRegisteredOn(t time.Time) {
	m._RegisteredOn = &t
}

// RegisteredOn returns the value of the "RegisteredOn" field in the mutation.
func (m *ChapterMutation) RegisteredOn() (r time.Time, exists bool) {
	v := m._RegisteredOn
	if v == nil {
		return
	}
	return *v, true
}

// OldRegisteredOn returns the old "RegisteredOn" field's value of the Chapter entity.
// If the Chapter object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChapterMutation) OldRegisteredOn(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRegisteredOn is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRegisteredOn requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRegisteredOn: %w", err)
	}
	return oldValue.RegisteredOn, nil
}

// ResetRegisteredOn resets all changes to the "RegisteredOn" field.
func (m *ChapterMutation) ResetRegisteredOn() {
	m._RegisteredOn = nil
}

// SetMangaID sets the "Manga" edge to the Manga entity by id.
func (m *ChapterMutation) SetMangaID(id int) {
	m._Manga = &id
}

// ClearManga clears the "Manga" edge to the Manga entity.
func (m *ChapterMutation) ClearManga() {
	m.cleared_Manga = true
}

// MangaCleared reports if the "Manga" edge to the Manga entity was cleared.
func (m *ChapterMutation) MangaCleared() bool {
	return m.cleared_Manga
}

// MangaID returns the "Manga" edge ID in the mutation.
func (m *ChapterMutation) MangaID() (id int, exists bool) {
	if m._Manga != nil {
		return *m._Manga, true
	}
	return
}

// MangaIDs returns the "Manga" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// MangaID instead. It exists only for internal usage by the builders.
func (m *ChapterMutation) MangaIDs() (ids []int) {
	if id := m._Manga; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetManga resets all changes to the "Manga" edge.
func (m *ChapterMutation) ResetManga() {
	m._Manga = nil
	m.cleared_Manga = false
}

// Where appends a list predicates to the ChapterMutation builder.
func (m *ChapterMutation) Where(ps ...predicate.Chapter) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *ChapterMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Chapter).
func (m *ChapterMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ChapterMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m._ChapterID != nil {
		fields = append(fields, chapter.FieldChapterID)
	}
	if m._Num != nil {
		fields = append(fields, chapter.FieldNum)
	}
	if m._Title != nil {
		fields = append(fields, chapter.FieldTitle)
	}
	if m._CreatedOn != nil {
		fields = append(fields, chapter.FieldCreatedOn)
	}
	if m._RegisteredOn != nil {
		fields = append(fields, chapter.FieldRegisteredOn)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ChapterMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case chapter.FieldChapterID:
		return m.ChapterID()
	case chapter.FieldNum:
		return m.Num()
	case chapter.FieldTitle:
		return m.Title()
	case chapter.FieldCreatedOn:
		return m.CreatedOn()
	case chapter.FieldRegisteredOn:
		return m.RegisteredOn()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ChapterMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case chapter.FieldChapterID:
		return m.OldChapterID(ctx)
	case chapter.FieldNum:
		return m.OldNum(ctx)
	case chapter.FieldTitle:
		return m.OldTitle(ctx)
	case chapter.FieldCreatedOn:
		return m.OldCreatedOn(ctx)
	case chapter.FieldRegisteredOn:
		return m.OldRegisteredOn(ctx)
	}
	return nil, fmt.Errorf("unknown Chapter field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ChapterMutation) SetField(name string, value ent.Value) error {
	switch name {
	case chapter.FieldChapterID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetChapterID(v)
		return nil
	case chapter.FieldNum:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetNum(v)
		return nil
	case chapter.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	case chapter.FieldCreatedOn:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedOn(v)
		return nil
	case chapter.FieldRegisteredOn:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRegisteredOn(v)
		return nil
	}
	return fmt.Errorf("unknown Chapter field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ChapterMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ChapterMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ChapterMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Chapter numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ChapterMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(chapter.FieldTitle) {
		fields = append(fields, chapter.FieldTitle)
	}
	if m.FieldCleared(chapter.FieldCreatedOn) {
		fields = append(fields, chapter.FieldCreatedOn)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ChapterMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ChapterMutation) ClearField(name string) error {
	switch name {
	case chapter.FieldTitle:
		m.ClearTitle()
		return nil
	case chapter.FieldCreatedOn:
		m.ClearCreatedOn()
		return nil
	}
	return fmt.Errorf("unknown Chapter nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ChapterMutation) ResetField(name string) error {
	switch name {
	case chapter.FieldChapterID:
		m.ResetChapterID()
		return nil
	case chapter.FieldNum:
		m.ResetNum()
		return nil
	case chapter.FieldTitle:
		m.ResetTitle()
		return nil
	case chapter.FieldCreatedOn:
		m.ResetCreatedOn()
		return nil
	case chapter.FieldRegisteredOn:
		m.ResetRegisteredOn()
		return nil
	}
	return fmt.Errorf("unknown Chapter field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ChapterMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m._Manga != nil {
		edges = append(edges, chapter.EdgeManga)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ChapterMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case chapter.EdgeManga:
		if id := m._Manga; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ChapterMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ChapterMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ChapterMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.cleared_Manga {
		edges = append(edges, chapter.EdgeManga)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ChapterMutation) EdgeCleared(name string) bool {
	switch name {
	case chapter.EdgeManga:
		return m.cleared_Manga
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ChapterMutation) ClearEdge(name string) error {
	switch name {
	case chapter.EdgeManga:
		m.ClearManga()
		return nil
	}
	return fmt.Errorf("unknown Chapter unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ChapterMutation) ResetEdge(name string) error {
	switch name {
	case chapter.EdgeManga:
		m.ResetManga()
		return nil
	}
	return fmt.Errorf("unknown Chapter edge %s", name)
}

// MangaMutation represents an operation that mutates the Manga nodes in the graph.
type MangaMutation struct {
	config
	op               Op
	typ              string
	id               *int
	_MangaID         *string
	_Source          *string
	_Title           *string
	_Mapping         *string
	_RegisteredOn    *time.Time
	_FilteredGroups  *[]string
	clearedFields    map[string]struct{}
	_Chapters        map[int]struct{}
	removed_Chapters map[int]struct{}
	cleared_Chapters bool
	done             bool
	oldValue         func(context.Context) (*Manga, error)
	predicates       []predicate.Manga
}

var _ ent.Mutation = (*MangaMutation)(nil)

// mangaOption allows management of the mutation configuration using functional options.
type mangaOption func(*MangaMutation)

// newMangaMutation creates new mutation for the Manga entity.
func newMangaMutation(c config, op Op, opts ...mangaOption) *MangaMutation {
	m := &MangaMutation{
		config:        c,
		op:            op,
		typ:           TypeManga,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withMangaID sets the ID field of the mutation.
func withMangaID(id int) mangaOption {
	return func(m *MangaMutation) {
		var (
			err   error
			once  sync.Once
			value *Manga
		)
		m.oldValue = func(ctx context.Context) (*Manga, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Manga.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withManga sets the old Manga of the mutation.
func withManga(node *Manga) mangaOption {
	return func(m *MangaMutation) {
		m.oldValue = func(context.Context) (*Manga, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m MangaMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m MangaMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *MangaMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *MangaMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Manga.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetMangaID sets the "MangaID" field.
func (m *MangaMutation) SetMangaID(s string) {
	m._MangaID = &s
}

// MangaID returns the value of the "MangaID" field in the mutation.
func (m *MangaMutation) MangaID() (r string, exists bool) {
	v := m._MangaID
	if v == nil {
		return
	}
	return *v, true
}

// OldMangaID returns the old "MangaID" field's value of the Manga entity.
// If the Manga object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MangaMutation) OldMangaID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMangaID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMangaID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMangaID: %w", err)
	}
	return oldValue.MangaID, nil
}

// ResetMangaID resets all changes to the "MangaID" field.
func (m *MangaMutation) ResetMangaID() {
	m._MangaID = nil
}

// SetSource sets the "Source" field.
func (m *MangaMutation) SetSource(s string) {
	m._Source = &s
}

// Source returns the value of the "Source" field in the mutation.
func (m *MangaMutation) Source() (r string, exists bool) {
	v := m._Source
	if v == nil {
		return
	}
	return *v, true
}

// OldSource returns the old "Source" field's value of the Manga entity.
// If the Manga object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MangaMutation) OldSource(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSource is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSource requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSource: %w", err)
	}
	return oldValue.Source, nil
}

// ResetSource resets all changes to the "Source" field.
func (m *MangaMutation) ResetSource() {
	m._Source = nil
}

// SetTitle sets the "Title" field.
func (m *MangaMutation) SetTitle(s string) {
	m._Title = &s
}

// Title returns the value of the "Title" field in the mutation.
func (m *MangaMutation) Title() (r string, exists bool) {
	v := m._Title
	if v == nil {
		return
	}
	return *v, true
}

// OldTitle returns the old "Title" field's value of the Manga entity.
// If the Manga object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MangaMutation) OldTitle(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTitle is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTitle requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTitle: %w", err)
	}
	return oldValue.Title, nil
}

// ResetTitle resets all changes to the "Title" field.
func (m *MangaMutation) ResetTitle() {
	m._Title = nil
}

// SetMapping sets the "Mapping" field.
func (m *MangaMutation) SetMapping(s string) {
	m._Mapping = &s
}

// Mapping returns the value of the "Mapping" field in the mutation.
func (m *MangaMutation) Mapping() (r string, exists bool) {
	v := m._Mapping
	if v == nil {
		return
	}
	return *v, true
}

// OldMapping returns the old "Mapping" field's value of the Manga entity.
// If the Manga object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MangaMutation) OldMapping(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMapping is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMapping requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMapping: %w", err)
	}
	return oldValue.Mapping, nil
}

// ResetMapping resets all changes to the "Mapping" field.
func (m *MangaMutation) ResetMapping() {
	m._Mapping = nil
}

// SetRegisteredOn sets the "RegisteredOn" field.
func (m *MangaMutation) SetRegisteredOn(t time.Time) {
	m._RegisteredOn = &t
}

// RegisteredOn returns the value of the "RegisteredOn" field in the mutation.
func (m *MangaMutation) RegisteredOn() (r time.Time, exists bool) {
	v := m._RegisteredOn
	if v == nil {
		return
	}
	return *v, true
}

// OldRegisteredOn returns the old "RegisteredOn" field's value of the Manga entity.
// If the Manga object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MangaMutation) OldRegisteredOn(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRegisteredOn is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRegisteredOn requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRegisteredOn: %w", err)
	}
	return oldValue.RegisteredOn, nil
}

// ResetRegisteredOn resets all changes to the "RegisteredOn" field.
func (m *MangaMutation) ResetRegisteredOn() {
	m._RegisteredOn = nil
}

// SetFilteredGroups sets the "FilteredGroups" field.
func (m *MangaMutation) SetFilteredGroups(s []string) {
	m._FilteredGroups = &s
}

// FilteredGroups returns the value of the "FilteredGroups" field in the mutation.
func (m *MangaMutation) FilteredGroups() (r []string, exists bool) {
	v := m._FilteredGroups
	if v == nil {
		return
	}
	return *v, true
}

// OldFilteredGroups returns the old "FilteredGroups" field's value of the Manga entity.
// If the Manga object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MangaMutation) OldFilteredGroups(ctx context.Context) (v []string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFilteredGroups is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFilteredGroups requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFilteredGroups: %w", err)
	}
	return oldValue.FilteredGroups, nil
}

// ResetFilteredGroups resets all changes to the "FilteredGroups" field.
func (m *MangaMutation) ResetFilteredGroups() {
	m._FilteredGroups = nil
}

// AddChapterIDs adds the "Chapters" edge to the Chapter entity by ids.
func (m *MangaMutation) AddChapterIDs(ids ...int) {
	if m._Chapters == nil {
		m._Chapters = make(map[int]struct{})
	}
	for i := range ids {
		m._Chapters[ids[i]] = struct{}{}
	}
}

// ClearChapters clears the "Chapters" edge to the Chapter entity.
func (m *MangaMutation) ClearChapters() {
	m.cleared_Chapters = true
}

// ChaptersCleared reports if the "Chapters" edge to the Chapter entity was cleared.
func (m *MangaMutation) ChaptersCleared() bool {
	return m.cleared_Chapters
}

// RemoveChapterIDs removes the "Chapters" edge to the Chapter entity by IDs.
func (m *MangaMutation) RemoveChapterIDs(ids ...int) {
	if m.removed_Chapters == nil {
		m.removed_Chapters = make(map[int]struct{})
	}
	for i := range ids {
		delete(m._Chapters, ids[i])
		m.removed_Chapters[ids[i]] = struct{}{}
	}
}

// RemovedChapters returns the removed IDs of the "Chapters" edge to the Chapter entity.
func (m *MangaMutation) RemovedChaptersIDs() (ids []int) {
	for id := range m.removed_Chapters {
		ids = append(ids, id)
	}
	return
}

// ChaptersIDs returns the "Chapters" edge IDs in the mutation.
func (m *MangaMutation) ChaptersIDs() (ids []int) {
	for id := range m._Chapters {
		ids = append(ids, id)
	}
	return
}

// ResetChapters resets all changes to the "Chapters" edge.
func (m *MangaMutation) ResetChapters() {
	m._Chapters = nil
	m.cleared_Chapters = false
	m.removed_Chapters = nil
}

// Where appends a list predicates to the MangaMutation builder.
func (m *MangaMutation) Where(ps ...predicate.Manga) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *MangaMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Manga).
func (m *MangaMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *MangaMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m._MangaID != nil {
		fields = append(fields, manga.FieldMangaID)
	}
	if m._Source != nil {
		fields = append(fields, manga.FieldSource)
	}
	if m._Title != nil {
		fields = append(fields, manga.FieldTitle)
	}
	if m._Mapping != nil {
		fields = append(fields, manga.FieldMapping)
	}
	if m._RegisteredOn != nil {
		fields = append(fields, manga.FieldRegisteredOn)
	}
	if m._FilteredGroups != nil {
		fields = append(fields, manga.FieldFilteredGroups)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *MangaMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case manga.FieldMangaID:
		return m.MangaID()
	case manga.FieldSource:
		return m.Source()
	case manga.FieldTitle:
		return m.Title()
	case manga.FieldMapping:
		return m.Mapping()
	case manga.FieldRegisteredOn:
		return m.RegisteredOn()
	case manga.FieldFilteredGroups:
		return m.FilteredGroups()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *MangaMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case manga.FieldMangaID:
		return m.OldMangaID(ctx)
	case manga.FieldSource:
		return m.OldSource(ctx)
	case manga.FieldTitle:
		return m.OldTitle(ctx)
	case manga.FieldMapping:
		return m.OldMapping(ctx)
	case manga.FieldRegisteredOn:
		return m.OldRegisteredOn(ctx)
	case manga.FieldFilteredGroups:
		return m.OldFilteredGroups(ctx)
	}
	return nil, fmt.Errorf("unknown Manga field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *MangaMutation) SetField(name string, value ent.Value) error {
	switch name {
	case manga.FieldMangaID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMangaID(v)
		return nil
	case manga.FieldSource:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSource(v)
		return nil
	case manga.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	case manga.FieldMapping:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMapping(v)
		return nil
	case manga.FieldRegisteredOn:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRegisteredOn(v)
		return nil
	case manga.FieldFilteredGroups:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFilteredGroups(v)
		return nil
	}
	return fmt.Errorf("unknown Manga field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *MangaMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *MangaMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *MangaMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Manga numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *MangaMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *MangaMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *MangaMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Manga nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *MangaMutation) ResetField(name string) error {
	switch name {
	case manga.FieldMangaID:
		m.ResetMangaID()
		return nil
	case manga.FieldSource:
		m.ResetSource()
		return nil
	case manga.FieldTitle:
		m.ResetTitle()
		return nil
	case manga.FieldMapping:
		m.ResetMapping()
		return nil
	case manga.FieldRegisteredOn:
		m.ResetRegisteredOn()
		return nil
	case manga.FieldFilteredGroups:
		m.ResetFilteredGroups()
		return nil
	}
	return fmt.Errorf("unknown Manga field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *MangaMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m._Chapters != nil {
		edges = append(edges, manga.EdgeChapters)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *MangaMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case manga.EdgeChapters:
		ids := make([]ent.Value, 0, len(m._Chapters))
		for id := range m._Chapters {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *MangaMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removed_Chapters != nil {
		edges = append(edges, manga.EdgeChapters)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *MangaMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case manga.EdgeChapters:
		ids := make([]ent.Value, 0, len(m.removed_Chapters))
		for id := range m.removed_Chapters {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *MangaMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.cleared_Chapters {
		edges = append(edges, manga.EdgeChapters)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *MangaMutation) EdgeCleared(name string) bool {
	switch name {
	case manga.EdgeChapters:
		return m.cleared_Chapters
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *MangaMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Manga unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *MangaMutation) ResetEdge(name string) error {
	switch name {
	case manga.EdgeChapters:
		m.ResetChapters()
		return nil
	}
	return fmt.Errorf("unknown Manga edge %s", name)
}
