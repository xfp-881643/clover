package clover

import (
	"errors"

	nnid "github.com/matoous/go-nanoid/v2"
)

// Collection creation errors
var (
	ErrCollectionExist    = errors.New("collection already exist")
	ErrCollectionNotExist = errors.New("no such collection")
)

// DB represents the entry point of each clover database.
type DB struct {
	dir    string
	engine StorageEngine
}

// Query simply returns the collection with the supplied name. Use it to initialize a new query.
func (db *DB) Query(name string) *Query {
	return &Query{
		collection: name,
		criteria:   nil,
		engine:     db.engine,
		limit:      -1,
		skip:       0,
		sortOpts:   nil,
	}
}

// CreateCollection creates a new empty collection with the given name.
func (db *DB) CreateCollection(name string) error {
	return db.engine.CreateCollection(name)
}

// DropCollection removes the collection with the given name, deleting any content on disk.
func (db *DB) DropCollection(name string) error {
	return db.engine.DropCollection(name)
}

// HasCollection returns true if and only if the database contains a collection with the given name.
func (db *DB) HasCollection(name string) (bool, error) {
	return db.engine.HasCollection(name)
}

func NewObjectId() string {
	return nnid.Must()
}

func isValidObjectId(id string) bool {
	return len(id) == 21
}

// Insert adds the supplied documents to a collection.
func (db *DB) Insert(collectionName string, docs ...*Document) error {
	for _, doc := range docs {
		if !isValidObjectId(doc.ObjectId()) {
			objectId := NewObjectId()
			doc.Set(objectIdField, objectId)
		}
	}
	return db.engine.Insert(collectionName, docs...)
}

// Save or update a document
func (db *DB) Save(collectionName string, doc *Document) error {
	if !doc.Has(objectIdField) {
		return db.Insert(collectionName, doc)
	}
	return db.Query(collectionName).ReplaceById(doc.ObjectId(), doc)
}

// InsertOne inserts a single document to an existing collection. It returns the id of the inserted document.
func (db *DB) InsertOne(collectionName string, doc *Document) (string, error) {
	err := db.Insert(collectionName, doc)
	return doc.ObjectId(), err
}

// Open opens a new clover database on the supplied path. If such a folder doesn't exist, it is automatically created.
func Open(dir string, opts ...Option) (*DB, error) {
	config, err := defaultConfig().applyOptions(opts)
	if err != nil {
		return nil, err
	}

	if !config.InMemory {
		if err := makeDirIfNotExists(dir); err != nil {
			return nil, err
		}
	}

	db := &DB{
		dir:    dir,
		engine: config.Storage,
	}
	return db, db.engine.Open(dir)
}

// Close releases all the resources and closes the database. After the call, the instance will no more be usable.
func (db *DB) Close() error {
	return db.engine.Close()
}

// ListCollections returns a slice of strings containing the name of each collection stored in the db.
func (db *DB) ListCollections() ([]string, error) {
	return db.engine.ListCollections()
}
