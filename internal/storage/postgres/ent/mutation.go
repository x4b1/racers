// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xabi93/racers/internal/storage/postgres/ent/race"
	"github.com/xabi93/racers/internal/storage/postgres/ent/user"

	"github.com/facebook/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeRace = "Race"
	TypeUser = "User"
)

// RaceMutation represents an operation that mutate the Races
// nodes in the graph.
type RaceMutation struct {
	config
	op                 Op
	typ                string
	id                 *string
	name               *string
	date               *time.Time
	clearedFields      map[string]struct{}
	competitors        map[string]struct{}
	removedcompetitors map[string]struct{}
	clearedcompetitors bool
	done               bool
	oldValue           func(context.Context) (*Race, error)
}

var _ ent.Mutation = (*RaceMutation)(nil)

// raceOption allows to manage the mutation configuration using functional options.
type raceOption func(*RaceMutation)

// newRaceMutation creates new mutation for $n.Name.
func newRaceMutation(c config, op Op, opts ...raceOption) *RaceMutation {
	m := &RaceMutation{
		config:        c,
		op:            op,
		typ:           TypeRace,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withRaceID sets the id field of the mutation.
func withRaceID(id string) raceOption {
	return func(m *RaceMutation) {
		var (
			err   error
			once  sync.Once
			value *Race
		)
		m.oldValue = func(ctx context.Context) (*Race, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Race.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withRace sets the old Race of the mutation.
func withRace(node *Race) raceOption {
	return func(m *RaceMutation) {
		m.oldValue = func(context.Context) (*Race, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m RaceMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m RaceMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that, this
// operation is accepted only on Race creation.
func (m *RaceMutation) SetID(id string) {
	m.id = &id
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *RaceMutation) ID() (id string, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetName sets the name field.
func (m *RaceMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *RaceMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old name value of the Race.
// If the Race object wasn't provided to the builder, the object is fetched
// from the database.
// An error is returned if the mutation operation is not UpdateOne, or database query fails.
func (m *RaceMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldName is allowed only on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName reset all changes of the "name" field.
func (m *RaceMutation) ResetName() {
	m.name = nil
}

// SetDate sets the date field.
func (m *RaceMutation) SetDate(t time.Time) {
	m.date = &t
}

// Date returns the date value in the mutation.
func (m *RaceMutation) Date() (r time.Time, exists bool) {
	v := m.date
	if v == nil {
		return
	}
	return *v, true
}

// OldDate returns the old date value of the Race.
// If the Race object wasn't provided to the builder, the object is fetched
// from the database.
// An error is returned if the mutation operation is not UpdateOne, or database query fails.
func (m *RaceMutation) OldDate(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldDate is allowed only on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldDate requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDate: %w", err)
	}
	return oldValue.Date, nil
}

// ResetDate reset all changes of the "date" field.
func (m *RaceMutation) ResetDate() {
	m.date = nil
}

// AddCompetitorIDs adds the competitors edge to User by ids.
func (m *RaceMutation) AddCompetitorIDs(ids ...string) {
	if m.competitors == nil {
		m.competitors = make(map[string]struct{})
	}
	for i := range ids {
		m.competitors[ids[i]] = struct{}{}
	}
}

// ClearCompetitors clears the competitors edge to User.
func (m *RaceMutation) ClearCompetitors() {
	m.clearedcompetitors = true
}

// CompetitorsCleared returns if the edge competitors was cleared.
func (m *RaceMutation) CompetitorsCleared() bool {
	return m.clearedcompetitors
}

// RemoveCompetitorIDs removes the competitors edge to User by ids.
func (m *RaceMutation) RemoveCompetitorIDs(ids ...string) {
	if m.removedcompetitors == nil {
		m.removedcompetitors = make(map[string]struct{})
	}
	for i := range ids {
		m.removedcompetitors[ids[i]] = struct{}{}
	}
}

// RemovedCompetitors returns the removed ids of competitors.
func (m *RaceMutation) RemovedCompetitorsIDs() (ids []string) {
	for id := range m.removedcompetitors {
		ids = append(ids, id)
	}
	return
}

// CompetitorsIDs returns the competitors ids in the mutation.
func (m *RaceMutation) CompetitorsIDs() (ids []string) {
	for id := range m.competitors {
		ids = append(ids, id)
	}
	return
}

// ResetCompetitors reset all changes of the "competitors" edge.
func (m *RaceMutation) ResetCompetitors() {
	m.competitors = nil
	m.clearedcompetitors = false
	m.removedcompetitors = nil
}

// Op returns the operation name.
func (m *RaceMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Race).
func (m *RaceMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *RaceMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m.name != nil {
		fields = append(fields, race.FieldName)
	}
	if m.date != nil {
		fields = append(fields, race.FieldDate)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *RaceMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case race.FieldName:
		return m.Name()
	case race.FieldDate:
		return m.Date()
	}
	return nil, false
}

// OldField returns the old value of the field from the database.
// An error is returned if the mutation operation is not UpdateOne,
// or the query to the database was failed.
func (m *RaceMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case race.FieldName:
		return m.OldName(ctx)
	case race.FieldDate:
		return m.OldDate(ctx)
	}
	return nil, fmt.Errorf("unknown Race field %s", name)
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *RaceMutation) SetField(name string, value ent.Value) error {
	switch name {
	case race.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case race.FieldDate:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDate(v)
		return nil
	}
	return fmt.Errorf("unknown Race field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *RaceMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *RaceMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *RaceMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Race numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *RaceMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *RaceMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *RaceMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Race nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *RaceMutation) ResetField(name string) error {
	switch name {
	case race.FieldName:
		m.ResetName()
		return nil
	case race.FieldDate:
		m.ResetDate()
		return nil
	}
	return fmt.Errorf("unknown Race field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *RaceMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.competitors != nil {
		edges = append(edges, race.EdgeCompetitors)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *RaceMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case race.EdgeCompetitors:
		ids := make([]ent.Value, 0, len(m.competitors))
		for id := range m.competitors {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *RaceMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedcompetitors != nil {
		edges = append(edges, race.EdgeCompetitors)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *RaceMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case race.EdgeCompetitors:
		ids := make([]ent.Value, 0, len(m.removedcompetitors))
		for id := range m.removedcompetitors {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *RaceMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedcompetitors {
		edges = append(edges, race.EdgeCompetitors)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *RaceMutation) EdgeCleared(name string) bool {
	switch name {
	case race.EdgeCompetitors:
		return m.clearedcompetitors
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *RaceMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Race unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *RaceMutation) ResetEdge(name string) error {
	switch name {
	case race.EdgeCompetitors:
		m.ResetCompetitors()
		return nil
	}
	return fmt.Errorf("unknown Race edge %s", name)
}

// UserMutation represents an operation that mutate the Users
// nodes in the graph.
type UserMutation struct {
	config
	op            Op
	typ           string
	id            *string
	name          *string
	clearedFields map[string]struct{}
	races         map[string]struct{}
	removedraces  map[string]struct{}
	clearedraces  bool
	done          bool
	oldValue      func(context.Context) (*User, error)
}

var _ ent.Mutation = (*UserMutation)(nil)

// userOption allows to manage the mutation configuration using functional options.
type userOption func(*UserMutation)

// newUserMutation creates new mutation for $n.Name.
func newUserMutation(c config, op Op, opts ...userOption) *UserMutation {
	m := &UserMutation{
		config:        c,
		op:            op,
		typ:           TypeUser,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withUserID sets the id field of the mutation.
func withUserID(id string) userOption {
	return func(m *UserMutation) {
		var (
			err   error
			once  sync.Once
			value *User
		)
		m.oldValue = func(ctx context.Context) (*User, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().User.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withUser sets the old User of the mutation.
func withUser(node *User) userOption {
	return func(m *UserMutation) {
		m.oldValue = func(context.Context) (*User, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m UserMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m UserMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that, this
// operation is accepted only on User creation.
func (m *UserMutation) SetID(id string) {
	m.id = &id
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *UserMutation) ID() (id string, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetName sets the name field.
func (m *UserMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *UserMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old name value of the User.
// If the User object wasn't provided to the builder, the object is fetched
// from the database.
// An error is returned if the mutation operation is not UpdateOne, or database query fails.
func (m *UserMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldName is allowed only on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName reset all changes of the "name" field.
func (m *UserMutation) ResetName() {
	m.name = nil
}

// AddRaceIDs adds the races edge to Race by ids.
func (m *UserMutation) AddRaceIDs(ids ...string) {
	if m.races == nil {
		m.races = make(map[string]struct{})
	}
	for i := range ids {
		m.races[ids[i]] = struct{}{}
	}
}

// ClearRaces clears the races edge to Race.
func (m *UserMutation) ClearRaces() {
	m.clearedraces = true
}

// RacesCleared returns if the edge races was cleared.
func (m *UserMutation) RacesCleared() bool {
	return m.clearedraces
}

// RemoveRaceIDs removes the races edge to Race by ids.
func (m *UserMutation) RemoveRaceIDs(ids ...string) {
	if m.removedraces == nil {
		m.removedraces = make(map[string]struct{})
	}
	for i := range ids {
		m.removedraces[ids[i]] = struct{}{}
	}
}

// RemovedRaces returns the removed ids of races.
func (m *UserMutation) RemovedRacesIDs() (ids []string) {
	for id := range m.removedraces {
		ids = append(ids, id)
	}
	return
}

// RacesIDs returns the races ids in the mutation.
func (m *UserMutation) RacesIDs() (ids []string) {
	for id := range m.races {
		ids = append(ids, id)
	}
	return
}

// ResetRaces reset all changes of the "races" edge.
func (m *UserMutation) ResetRaces() {
	m.races = nil
	m.clearedraces = false
	m.removedraces = nil
}

// Op returns the operation name.
func (m *UserMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (User).
func (m *UserMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *UserMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.name != nil {
		fields = append(fields, user.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *UserMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case user.FieldName:
		return m.Name()
	}
	return nil, false
}

// OldField returns the old value of the field from the database.
// An error is returned if the mutation operation is not UpdateOne,
// or the query to the database was failed.
func (m *UserMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case user.FieldName:
		return m.OldName(ctx)
	}
	return nil, fmt.Errorf("unknown User field %s", name)
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *UserMutation) SetField(name string, value ent.Value) error {
	switch name {
	case user.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *UserMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *UserMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *UserMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown User numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *UserMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *UserMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *UserMutation) ClearField(name string) error {
	return fmt.Errorf("unknown User nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *UserMutation) ResetField(name string) error {
	switch name {
	case user.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *UserMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.races != nil {
		edges = append(edges, user.EdgeRaces)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *UserMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case user.EdgeRaces:
		ids := make([]ent.Value, 0, len(m.races))
		for id := range m.races {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *UserMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedraces != nil {
		edges = append(edges, user.EdgeRaces)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *UserMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case user.EdgeRaces:
		ids := make([]ent.Value, 0, len(m.removedraces))
		for id := range m.removedraces {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *UserMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedraces {
		edges = append(edges, user.EdgeRaces)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *UserMutation) EdgeCleared(name string) bool {
	switch name {
	case user.EdgeRaces:
		return m.clearedraces
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *UserMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown User unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *UserMutation) ResetEdge(name string) error {
	switch name {
	case user.EdgeRaces:
		m.ResetRaces()
		return nil
	}
	return fmt.Errorf("unknown User edge %s", name)
}
