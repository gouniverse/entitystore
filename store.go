package entitystore

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Store defines an entity store
type Store struct {
	entityTableName    string
	attributeTableName string
	db                 *gorm.DB
}

// StoreOption options for the vault store
type StoreOption func(*Store)

// WithDriverAndDNS sets the driver and the DNS for the database for the cache store
func WithDriverAndDNS(driverName string, dsn string) StoreOption {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return func(s *Store) {
		s.db = db
	}
}

// WithGormDb sets the GORM database for the cache store
func WithGormDb(db *gorm.DB) StoreOption {
	return func(s *Store) {
		s.db = db
	}
}

// WithEntityTableName sets the table name for the cache store
func WithEntityTableName(entityTableName string) StoreOption {
	return func(s *Store) {
		s.entityTableName = entityTableName
	}
}

// WithAttributeTableName sets the table name for the cache store
func WithAttributeTableName(attributeTableName string) StoreOption {
	return func(s *Store) {
		s.attributeTableName = attributeTableName
	}
}

// // NewStore creates a new entity store
// func NewStore(driverName string, dsn string, entityTableName string, attributeTableName string) *Store {
// 	log.Println("New entity store: " + dsn)

// 	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	db.Table(entityTableName).AutoMigrate(&entity{})
// 	db.Table(attributeTableName).AutoMigrate(&attribute{})

// 	st := &Store{
// 		db:                 db,
// 		entityTableName:    entityTableName,
// 		attributeTableName: attributeTableName,
// 	}

// 	return st
// }

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) *Store {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.entityTableName == "" {
		log.Panic("Entity store: entityTableName is required")
	}

	if store.entityTableName == "" {
		log.Panic("Entity store: attributeTableName is required")
	}

	store.db.Table(store.entityTableName).AutoMigrate(&entity{})
	store.db.Table(store.attributeTableName).AutoMigrate(&attribute{})

	return store
}

// NewStoreGorm creates a new entity store
// func NewStoreGorm(db *gorm.DB, entityTableName string, attributeTableName string) *Store {
// 	log.Println("New entity store: " + db.Name())

// 	db.Table(entityTableName).AutoMigrate(&entity{})
// 	db.Table(attributeTableName).AutoMigrate(&attribute{})

// 	st := &Store{
// 		db:                 db,
// 		entityTableName:    entityTableName,
// 		attributeTableName: attributeTableName,
// 	}

// 	return st
// }
