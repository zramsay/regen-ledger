// Code generated by protoc-gen-go-cosmos-orm. DO NOT EDIT.

package datav1

import (
	context "context"
	ormlist "github.com/cosmos/cosmos-sdk/orm/model/ormlist"
	ormtable "github.com/cosmos/cosmos-sdk/orm/model/ormtable"
	ormerrors "github.com/cosmos/cosmos-sdk/orm/types/ormerrors"
)

type DataIDTable interface {
	Insert(ctx context.Context, dataID *DataID) error
	Update(ctx context.Context, dataID *DataID) error
	Save(ctx context.Context, dataID *DataID) error
	Delete(ctx context.Context, dataID *DataID) error
	Has(ctx context.Context, id []byte) (found bool, err error)
	// Get returns nil and an error which responds true to ormerrors.IsNotFound() if the record was not found.
	Get(ctx context.Context, id []byte) (*DataID, error)
	HasByIri(ctx context.Context, iri string) (found bool, err error)
	// GetByIri returns nil and an error which responds true to ormerrors.IsNotFound() if the record was not found.
	GetByIri(ctx context.Context, iri string) (*DataID, error)
	List(ctx context.Context, prefixKey DataIDIndexKey, opts ...ormlist.Option) (DataIDIterator, error)
	ListRange(ctx context.Context, from, to DataIDIndexKey, opts ...ormlist.Option) (DataIDIterator, error)
	DeleteBy(ctx context.Context, prefixKey DataIDIndexKey) error
	DeleteRange(ctx context.Context, from, to DataIDIndexKey) error

	doNotImplement()
}

type DataIDIterator struct {
	ormtable.Iterator
}

func (i DataIDIterator) Value() (*DataID, error) {
	var dataID DataID
	err := i.UnmarshalMessage(&dataID)
	return &dataID, err
}

type DataIDIndexKey interface {
	id() uint32
	values() []interface{}
	dataIDIndexKey()
}

// primary key starting index..
type DataIDPrimaryKey = DataIDIdIndexKey

type DataIDIdIndexKey struct {
	vs []interface{}
}

func (x DataIDIdIndexKey) id() uint32            { return 0 }
func (x DataIDIdIndexKey) values() []interface{} { return x.vs }
func (x DataIDIdIndexKey) dataIDIndexKey()       {}

func (this DataIDIdIndexKey) WithId(id []byte) DataIDIdIndexKey {
	this.vs = []interface{}{id}
	return this
}

type DataIDIriIndexKey struct {
	vs []interface{}
}

func (x DataIDIriIndexKey) id() uint32            { return 1 }
func (x DataIDIriIndexKey) values() []interface{} { return x.vs }
func (x DataIDIriIndexKey) dataIDIndexKey()       {}

func (this DataIDIriIndexKey) WithIri(iri string) DataIDIriIndexKey {
	this.vs = []interface{}{iri}
	return this
}

type dataIDTable struct {
	table ormtable.Table
}

func (this dataIDTable) Insert(ctx context.Context, dataID *DataID) error {
	return this.table.Insert(ctx, dataID)
}

func (this dataIDTable) Update(ctx context.Context, dataID *DataID) error {
	return this.table.Update(ctx, dataID)
}

func (this dataIDTable) Save(ctx context.Context, dataID *DataID) error {
	return this.table.Save(ctx, dataID)
}

func (this dataIDTable) Delete(ctx context.Context, dataID *DataID) error {
	return this.table.Delete(ctx, dataID)
}

func (this dataIDTable) Has(ctx context.Context, id []byte) (found bool, err error) {
	return this.table.PrimaryKey().Has(ctx, id)
}

func (this dataIDTable) Get(ctx context.Context, id []byte) (*DataID, error) {
	var dataID DataID
	found, err := this.table.PrimaryKey().Get(ctx, &dataID, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ormerrors.NotFound
	}
	return &dataID, nil
}

func (this dataIDTable) HasByIri(ctx context.Context, iri string) (found bool, err error) {
	return this.table.GetIndexByID(1).(ormtable.UniqueIndex).Has(ctx,
		iri,
	)
}

func (this dataIDTable) GetByIri(ctx context.Context, iri string) (*DataID, error) {
	var dataID DataID
	found, err := this.table.GetIndexByID(1).(ormtable.UniqueIndex).Get(ctx, &dataID,
		iri,
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ormerrors.NotFound
	}
	return &dataID, nil
}

func (this dataIDTable) List(ctx context.Context, prefixKey DataIDIndexKey, opts ...ormlist.Option) (DataIDIterator, error) {
	it, err := this.table.GetIndexByID(prefixKey.id()).List(ctx, prefixKey.values(), opts...)
	return DataIDIterator{it}, err
}

func (this dataIDTable) ListRange(ctx context.Context, from, to DataIDIndexKey, opts ...ormlist.Option) (DataIDIterator, error) {
	it, err := this.table.GetIndexByID(from.id()).ListRange(ctx, from.values(), to.values(), opts...)
	return DataIDIterator{it}, err
}

func (this dataIDTable) DeleteBy(ctx context.Context, prefixKey DataIDIndexKey) error {
	return this.table.GetIndexByID(prefixKey.id()).DeleteBy(ctx, prefixKey.values()...)
}

func (this dataIDTable) DeleteRange(ctx context.Context, from, to DataIDIndexKey) error {
	return this.table.GetIndexByID(from.id()).DeleteRange(ctx, from.values(), to.values())
}

func (this dataIDTable) doNotImplement() {}

var _ DataIDTable = dataIDTable{}

func NewDataIDTable(db ormtable.Schema) (DataIDTable, error) {
	table := db.GetTable(&DataID{})
	if table == nil {
		return nil, ormerrors.TableNotFound.Wrap(string((&DataID{}).ProtoReflect().Descriptor().FullName()))
	}
	return dataIDTable{table}, nil
}

type DataAnchorTable interface {
	Insert(ctx context.Context, dataAnchor *DataAnchor) error
	Update(ctx context.Context, dataAnchor *DataAnchor) error
	Save(ctx context.Context, dataAnchor *DataAnchor) error
	Delete(ctx context.Context, dataAnchor *DataAnchor) error
	Has(ctx context.Context, id []byte) (found bool, err error)
	// Get returns nil and an error which responds true to ormerrors.IsNotFound() if the record was not found.
	Get(ctx context.Context, id []byte) (*DataAnchor, error)
	List(ctx context.Context, prefixKey DataAnchorIndexKey, opts ...ormlist.Option) (DataAnchorIterator, error)
	ListRange(ctx context.Context, from, to DataAnchorIndexKey, opts ...ormlist.Option) (DataAnchorIterator, error)
	DeleteBy(ctx context.Context, prefixKey DataAnchorIndexKey) error
	DeleteRange(ctx context.Context, from, to DataAnchorIndexKey) error

	doNotImplement()
}

type DataAnchorIterator struct {
	ormtable.Iterator
}

func (i DataAnchorIterator) Value() (*DataAnchor, error) {
	var dataAnchor DataAnchor
	err := i.UnmarshalMessage(&dataAnchor)
	return &dataAnchor, err
}

type DataAnchorIndexKey interface {
	id() uint32
	values() []interface{}
	dataAnchorIndexKey()
}

// primary key starting index..
type DataAnchorPrimaryKey = DataAnchorIdIndexKey

type DataAnchorIdIndexKey struct {
	vs []interface{}
}

func (x DataAnchorIdIndexKey) id() uint32            { return 0 }
func (x DataAnchorIdIndexKey) values() []interface{} { return x.vs }
func (x DataAnchorIdIndexKey) dataAnchorIndexKey()   {}

func (this DataAnchorIdIndexKey) WithId(id []byte) DataAnchorIdIndexKey {
	this.vs = []interface{}{id}
	return this
}

type dataAnchorTable struct {
	table ormtable.Table
}

func (this dataAnchorTable) Insert(ctx context.Context, dataAnchor *DataAnchor) error {
	return this.table.Insert(ctx, dataAnchor)
}

func (this dataAnchorTable) Update(ctx context.Context, dataAnchor *DataAnchor) error {
	return this.table.Update(ctx, dataAnchor)
}

func (this dataAnchorTable) Save(ctx context.Context, dataAnchor *DataAnchor) error {
	return this.table.Save(ctx, dataAnchor)
}

func (this dataAnchorTable) Delete(ctx context.Context, dataAnchor *DataAnchor) error {
	return this.table.Delete(ctx, dataAnchor)
}

func (this dataAnchorTable) Has(ctx context.Context, id []byte) (found bool, err error) {
	return this.table.PrimaryKey().Has(ctx, id)
}

func (this dataAnchorTable) Get(ctx context.Context, id []byte) (*DataAnchor, error) {
	var dataAnchor DataAnchor
	found, err := this.table.PrimaryKey().Get(ctx, &dataAnchor, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ormerrors.NotFound
	}
	return &dataAnchor, nil
}

func (this dataAnchorTable) List(ctx context.Context, prefixKey DataAnchorIndexKey, opts ...ormlist.Option) (DataAnchorIterator, error) {
	it, err := this.table.GetIndexByID(prefixKey.id()).List(ctx, prefixKey.values(), opts...)
	return DataAnchorIterator{it}, err
}

func (this dataAnchorTable) ListRange(ctx context.Context, from, to DataAnchorIndexKey, opts ...ormlist.Option) (DataAnchorIterator, error) {
	it, err := this.table.GetIndexByID(from.id()).ListRange(ctx, from.values(), to.values(), opts...)
	return DataAnchorIterator{it}, err
}

func (this dataAnchorTable) DeleteBy(ctx context.Context, prefixKey DataAnchorIndexKey) error {
	return this.table.GetIndexByID(prefixKey.id()).DeleteBy(ctx, prefixKey.values()...)
}

func (this dataAnchorTable) DeleteRange(ctx context.Context, from, to DataAnchorIndexKey) error {
	return this.table.GetIndexByID(from.id()).DeleteRange(ctx, from.values(), to.values())
}

func (this dataAnchorTable) doNotImplement() {}

var _ DataAnchorTable = dataAnchorTable{}

func NewDataAnchorTable(db ormtable.Schema) (DataAnchorTable, error) {
	table := db.GetTable(&DataAnchor{})
	if table == nil {
		return nil, ormerrors.TableNotFound.Wrap(string((&DataAnchor{}).ProtoReflect().Descriptor().FullName()))
	}
	return dataAnchorTable{table}, nil
}

type DataAttestorTable interface {
	Insert(ctx context.Context, dataAttestor *DataAttestor) error
	Update(ctx context.Context, dataAttestor *DataAttestor) error
	Save(ctx context.Context, dataAttestor *DataAttestor) error
	Delete(ctx context.Context, dataAttestor *DataAttestor) error
	Has(ctx context.Context, id []byte, attestor []byte) (found bool, err error)
	// Get returns nil and an error which responds true to ormerrors.IsNotFound() if the record was not found.
	Get(ctx context.Context, id []byte, attestor []byte) (*DataAttestor, error)
	List(ctx context.Context, prefixKey DataAttestorIndexKey, opts ...ormlist.Option) (DataAttestorIterator, error)
	ListRange(ctx context.Context, from, to DataAttestorIndexKey, opts ...ormlist.Option) (DataAttestorIterator, error)
	DeleteBy(ctx context.Context, prefixKey DataAttestorIndexKey) error
	DeleteRange(ctx context.Context, from, to DataAttestorIndexKey) error

	doNotImplement()
}

type DataAttestorIterator struct {
	ormtable.Iterator
}

func (i DataAttestorIterator) Value() (*DataAttestor, error) {
	var dataAttestor DataAttestor
	err := i.UnmarshalMessage(&dataAttestor)
	return &dataAttestor, err
}

type DataAttestorIndexKey interface {
	id() uint32
	values() []interface{}
	dataAttestorIndexKey()
}

// primary key starting index..
type DataAttestorPrimaryKey = DataAttestorIdAttestorIndexKey

type DataAttestorIdAttestorIndexKey struct {
	vs []interface{}
}

func (x DataAttestorIdAttestorIndexKey) id() uint32            { return 0 }
func (x DataAttestorIdAttestorIndexKey) values() []interface{} { return x.vs }
func (x DataAttestorIdAttestorIndexKey) dataAttestorIndexKey() {}

func (this DataAttestorIdAttestorIndexKey) WithId(id []byte) DataAttestorIdAttestorIndexKey {
	this.vs = []interface{}{id}
	return this
}

func (this DataAttestorIdAttestorIndexKey) WithIdAttestor(id []byte, attestor []byte) DataAttestorIdAttestorIndexKey {
	this.vs = []interface{}{id, attestor}
	return this
}

type DataAttestorAttestorIndexKey struct {
	vs []interface{}
}

func (x DataAttestorAttestorIndexKey) id() uint32            { return 1 }
func (x DataAttestorAttestorIndexKey) values() []interface{} { return x.vs }
func (x DataAttestorAttestorIndexKey) dataAttestorIndexKey() {}

func (this DataAttestorAttestorIndexKey) WithAttestor(attestor []byte) DataAttestorAttestorIndexKey {
	this.vs = []interface{}{attestor}
	return this
}

type dataAttestorTable struct {
	table ormtable.Table
}

func (this dataAttestorTable) Insert(ctx context.Context, dataAttestor *DataAttestor) error {
	return this.table.Insert(ctx, dataAttestor)
}

func (this dataAttestorTable) Update(ctx context.Context, dataAttestor *DataAttestor) error {
	return this.table.Update(ctx, dataAttestor)
}

func (this dataAttestorTable) Save(ctx context.Context, dataAttestor *DataAttestor) error {
	return this.table.Save(ctx, dataAttestor)
}

func (this dataAttestorTable) Delete(ctx context.Context, dataAttestor *DataAttestor) error {
	return this.table.Delete(ctx, dataAttestor)
}

func (this dataAttestorTable) Has(ctx context.Context, id []byte, attestor []byte) (found bool, err error) {
	return this.table.PrimaryKey().Has(ctx, id, attestor)
}

func (this dataAttestorTable) Get(ctx context.Context, id []byte, attestor []byte) (*DataAttestor, error) {
	var dataAttestor DataAttestor
	found, err := this.table.PrimaryKey().Get(ctx, &dataAttestor, id, attestor)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ormerrors.NotFound
	}
	return &dataAttestor, nil
}

func (this dataAttestorTable) List(ctx context.Context, prefixKey DataAttestorIndexKey, opts ...ormlist.Option) (DataAttestorIterator, error) {
	it, err := this.table.GetIndexByID(prefixKey.id()).List(ctx, prefixKey.values(), opts...)
	return DataAttestorIterator{it}, err
}

func (this dataAttestorTable) ListRange(ctx context.Context, from, to DataAttestorIndexKey, opts ...ormlist.Option) (DataAttestorIterator, error) {
	it, err := this.table.GetIndexByID(from.id()).ListRange(ctx, from.values(), to.values(), opts...)
	return DataAttestorIterator{it}, err
}

func (this dataAttestorTable) DeleteBy(ctx context.Context, prefixKey DataAttestorIndexKey) error {
	return this.table.GetIndexByID(prefixKey.id()).DeleteBy(ctx, prefixKey.values()...)
}

func (this dataAttestorTable) DeleteRange(ctx context.Context, from, to DataAttestorIndexKey) error {
	return this.table.GetIndexByID(from.id()).DeleteRange(ctx, from.values(), to.values())
}

func (this dataAttestorTable) doNotImplement() {}

var _ DataAttestorTable = dataAttestorTable{}

func NewDataAttestorTable(db ormtable.Schema) (DataAttestorTable, error) {
	table := db.GetTable(&DataAttestor{})
	if table == nil {
		return nil, ormerrors.TableNotFound.Wrap(string((&DataAttestor{}).ProtoReflect().Descriptor().FullName()))
	}
	return dataAttestorTable{table}, nil
}

type ResolverTable interface {
	Insert(ctx context.Context, resolver *Resolver) error
	InsertReturningID(ctx context.Context, resolver *Resolver) (uint64, error)
	Update(ctx context.Context, resolver *Resolver) error
	Save(ctx context.Context, resolver *Resolver) error
	Delete(ctx context.Context, resolver *Resolver) error
	Has(ctx context.Context, id uint64) (found bool, err error)
	// Get returns nil and an error which responds true to ormerrors.IsNotFound() if the record was not found.
	Get(ctx context.Context, id uint64) (*Resolver, error)
	List(ctx context.Context, prefixKey ResolverIndexKey, opts ...ormlist.Option) (ResolverIterator, error)
	ListRange(ctx context.Context, from, to ResolverIndexKey, opts ...ormlist.Option) (ResolverIterator, error)
	DeleteBy(ctx context.Context, prefixKey ResolverIndexKey) error
	DeleteRange(ctx context.Context, from, to ResolverIndexKey) error

	doNotImplement()
}

type ResolverIterator struct {
	ormtable.Iterator
}

func (i ResolverIterator) Value() (*Resolver, error) {
	var resolver Resolver
	err := i.UnmarshalMessage(&resolver)
	return &resolver, err
}

type ResolverIndexKey interface {
	id() uint32
	values() []interface{}
	resolverIndexKey()
}

// primary key starting index..
type ResolverPrimaryKey = ResolverIdIndexKey

type ResolverIdIndexKey struct {
	vs []interface{}
}

func (x ResolverIdIndexKey) id() uint32            { return 0 }
func (x ResolverIdIndexKey) values() []interface{} { return x.vs }
func (x ResolverIdIndexKey) resolverIndexKey()     {}

func (this ResolverIdIndexKey) WithId(id uint64) ResolverIdIndexKey {
	this.vs = []interface{}{id}
	return this
}

type ResolverUrlIndexKey struct {
	vs []interface{}
}

func (x ResolverUrlIndexKey) id() uint32            { return 1 }
func (x ResolverUrlIndexKey) values() []interface{} { return x.vs }
func (x ResolverUrlIndexKey) resolverIndexKey()     {}

func (this ResolverUrlIndexKey) WithUrl(url string) ResolverUrlIndexKey {
	this.vs = []interface{}{url}
	return this
}

type ResolverManagerIndexKey struct {
	vs []interface{}
}

func (x ResolverManagerIndexKey) id() uint32            { return 2 }
func (x ResolverManagerIndexKey) values() []interface{} { return x.vs }
func (x ResolverManagerIndexKey) resolverIndexKey()     {}

func (this ResolverManagerIndexKey) WithManager(manager []byte) ResolverManagerIndexKey {
	this.vs = []interface{}{manager}
	return this
}

type resolverTable struct {
	table ormtable.AutoIncrementTable
}

func (this resolverTable) Insert(ctx context.Context, resolver *Resolver) error {
	return this.table.Insert(ctx, resolver)
}

func (this resolverTable) Update(ctx context.Context, resolver *Resolver) error {
	return this.table.Update(ctx, resolver)
}

func (this resolverTable) Save(ctx context.Context, resolver *Resolver) error {
	return this.table.Save(ctx, resolver)
}

func (this resolverTable) Delete(ctx context.Context, resolver *Resolver) error {
	return this.table.Delete(ctx, resolver)
}

func (this resolverTable) InsertReturningID(ctx context.Context, resolver *Resolver) (uint64, error) {
	return this.table.InsertReturningID(ctx, resolver)
}

func (this resolverTable) Has(ctx context.Context, id uint64) (found bool, err error) {
	return this.table.PrimaryKey().Has(ctx, id)
}

func (this resolverTable) Get(ctx context.Context, id uint64) (*Resolver, error) {
	var resolver Resolver
	found, err := this.table.PrimaryKey().Get(ctx, &resolver, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ormerrors.NotFound
	}
	return &resolver, nil
}

func (this resolverTable) List(ctx context.Context, prefixKey ResolverIndexKey, opts ...ormlist.Option) (ResolverIterator, error) {
	it, err := this.table.GetIndexByID(prefixKey.id()).List(ctx, prefixKey.values(), opts...)
	return ResolverIterator{it}, err
}

func (this resolverTable) ListRange(ctx context.Context, from, to ResolverIndexKey, opts ...ormlist.Option) (ResolverIterator, error) {
	it, err := this.table.GetIndexByID(from.id()).ListRange(ctx, from.values(), to.values(), opts...)
	return ResolverIterator{it}, err
}

func (this resolverTable) DeleteBy(ctx context.Context, prefixKey ResolverIndexKey) error {
	return this.table.GetIndexByID(prefixKey.id()).DeleteBy(ctx, prefixKey.values()...)
}

func (this resolverTable) DeleteRange(ctx context.Context, from, to ResolverIndexKey) error {
	return this.table.GetIndexByID(from.id()).DeleteRange(ctx, from.values(), to.values())
}

func (this resolverTable) doNotImplement() {}

var _ ResolverTable = resolverTable{}

func NewResolverTable(db ormtable.Schema) (ResolverTable, error) {
	table := db.GetTable(&Resolver{})
	if table == nil {
		return nil, ormerrors.TableNotFound.Wrap(string((&Resolver{}).ProtoReflect().Descriptor().FullName()))
	}
	return resolverTable{table.(ormtable.AutoIncrementTable)}, nil
}

type DataResolverTable interface {
	Insert(ctx context.Context, dataResolver *DataResolver) error
	Update(ctx context.Context, dataResolver *DataResolver) error
	Save(ctx context.Context, dataResolver *DataResolver) error
	Delete(ctx context.Context, dataResolver *DataResolver) error
	Has(ctx context.Context, id []byte, resolver_id uint64) (found bool, err error)
	// Get returns nil and an error which responds true to ormerrors.IsNotFound() if the record was not found.
	Get(ctx context.Context, id []byte, resolver_id uint64) (*DataResolver, error)
	List(ctx context.Context, prefixKey DataResolverIndexKey, opts ...ormlist.Option) (DataResolverIterator, error)
	ListRange(ctx context.Context, from, to DataResolverIndexKey, opts ...ormlist.Option) (DataResolverIterator, error)
	DeleteBy(ctx context.Context, prefixKey DataResolverIndexKey) error
	DeleteRange(ctx context.Context, from, to DataResolverIndexKey) error

	doNotImplement()
}

type DataResolverIterator struct {
	ormtable.Iterator
}

func (i DataResolverIterator) Value() (*DataResolver, error) {
	var dataResolver DataResolver
	err := i.UnmarshalMessage(&dataResolver)
	return &dataResolver, err
}

type DataResolverIndexKey interface {
	id() uint32
	values() []interface{}
	dataResolverIndexKey()
}

// primary key starting index..
type DataResolverPrimaryKey = DataResolverIdResolverIdIndexKey

type DataResolverIdResolverIdIndexKey struct {
	vs []interface{}
}

func (x DataResolverIdResolverIdIndexKey) id() uint32            { return 0 }
func (x DataResolverIdResolverIdIndexKey) values() []interface{} { return x.vs }
func (x DataResolverIdResolverIdIndexKey) dataResolverIndexKey() {}

func (this DataResolverIdResolverIdIndexKey) WithId(id []byte) DataResolverIdResolverIdIndexKey {
	this.vs = []interface{}{id}
	return this
}

func (this DataResolverIdResolverIdIndexKey) WithIdResolverId(id []byte, resolver_id uint64) DataResolverIdResolverIdIndexKey {
	this.vs = []interface{}{id, resolver_id}
	return this
}

type dataResolverTable struct {
	table ormtable.Table
}

func (this dataResolverTable) Insert(ctx context.Context, dataResolver *DataResolver) error {
	return this.table.Insert(ctx, dataResolver)
}

func (this dataResolverTable) Update(ctx context.Context, dataResolver *DataResolver) error {
	return this.table.Update(ctx, dataResolver)
}

func (this dataResolverTable) Save(ctx context.Context, dataResolver *DataResolver) error {
	return this.table.Save(ctx, dataResolver)
}

func (this dataResolverTable) Delete(ctx context.Context, dataResolver *DataResolver) error {
	return this.table.Delete(ctx, dataResolver)
}

func (this dataResolverTable) Has(ctx context.Context, id []byte, resolver_id uint64) (found bool, err error) {
	return this.table.PrimaryKey().Has(ctx, id, resolver_id)
}

func (this dataResolverTable) Get(ctx context.Context, id []byte, resolver_id uint64) (*DataResolver, error) {
	var dataResolver DataResolver
	found, err := this.table.PrimaryKey().Get(ctx, &dataResolver, id, resolver_id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ormerrors.NotFound
	}
	return &dataResolver, nil
}

func (this dataResolverTable) List(ctx context.Context, prefixKey DataResolverIndexKey, opts ...ormlist.Option) (DataResolverIterator, error) {
	it, err := this.table.GetIndexByID(prefixKey.id()).List(ctx, prefixKey.values(), opts...)
	return DataResolverIterator{it}, err
}

func (this dataResolverTable) ListRange(ctx context.Context, from, to DataResolverIndexKey, opts ...ormlist.Option) (DataResolverIterator, error) {
	it, err := this.table.GetIndexByID(from.id()).ListRange(ctx, from.values(), to.values(), opts...)
	return DataResolverIterator{it}, err
}

func (this dataResolverTable) DeleteBy(ctx context.Context, prefixKey DataResolverIndexKey) error {
	return this.table.GetIndexByID(prefixKey.id()).DeleteBy(ctx, prefixKey.values()...)
}

func (this dataResolverTable) DeleteRange(ctx context.Context, from, to DataResolverIndexKey) error {
	return this.table.GetIndexByID(from.id()).DeleteRange(ctx, from.values(), to.values())
}

func (this dataResolverTable) doNotImplement() {}

var _ DataResolverTable = dataResolverTable{}

func NewDataResolverTable(db ormtable.Schema) (DataResolverTable, error) {
	table := db.GetTable(&DataResolver{})
	if table == nil {
		return nil, ormerrors.TableNotFound.Wrap(string((&DataResolver{}).ProtoReflect().Descriptor().FullName()))
	}
	return dataResolverTable{table}, nil
}

type StateStore interface {
	DataIDTable() DataIDTable
	DataAnchorTable() DataAnchorTable
	DataAttestorTable() DataAttestorTable
	ResolverTable() ResolverTable
	DataResolverTable() DataResolverTable

	doNotImplement()
}

type stateStore struct {
	dataID       DataIDTable
	dataAnchor   DataAnchorTable
	dataAttestor DataAttestorTable
	resolver     ResolverTable
	dataResolver DataResolverTable
}

func (x stateStore) DataIDTable() DataIDTable {
	return x.dataID
}

func (x stateStore) DataAnchorTable() DataAnchorTable {
	return x.dataAnchor
}

func (x stateStore) DataAttestorTable() DataAttestorTable {
	return x.dataAttestor
}

func (x stateStore) ResolverTable() ResolverTable {
	return x.resolver
}

func (x stateStore) DataResolverTable() DataResolverTable {
	return x.dataResolver
}

func (stateStore) doNotImplement() {}

var _ StateStore = stateStore{}

func NewStateStore(db ormtable.Schema) (StateStore, error) {
	dataIDTable, err := NewDataIDTable(db)
	if err != nil {
		return nil, err
	}

	dataAnchorTable, err := NewDataAnchorTable(db)
	if err != nil {
		return nil, err
	}

	dataAttestorTable, err := NewDataAttestorTable(db)
	if err != nil {
		return nil, err
	}

	resolverTable, err := NewResolverTable(db)
	if err != nil {
		return nil, err
	}

	dataResolverTable, err := NewDataResolverTable(db)
	if err != nil {
		return nil, err
	}

	return stateStore{
		dataIDTable,
		dataAnchorTable,
		dataAttestorTable,
		resolverTable,
		dataResolverTable,
	}, nil
}
