package entitystore

import (
	"database/sql"
	"errors"
)

// NewStore creates a new entity store
// func NewStore(opts ...StoreOption) (*Store, error) {
// 	store := &Store{}
// 	for _, opt := range opts {
// 		opt(store)
// 	}

// 	if store.entityTableName == "" {
// 		return nil, errors.New("Entity store: entityTableName is required")
// 	}

// 	if store.attributeTableName == "" {
// 		return nil, errors.New("Entity store: attributeTableName is required")
// 	}

// 	store.entityTrashTableName = store.entityTableName + "_trash"
// 	store.attributeTrashTableName = store.attributeTableName + "_trash"

// 	if store.automigrateEnabled == true {
// 		store.AutoMigrate()
// 	}

// 	return store, nil
// }

// NewStoreOptions define the options for creating a new session store
type NewStoreOptions struct {
	EntityTableName         string
	AttributeTableName      string
	EntityTrashTableName    string
	AttributeTrashTableName string
	DB                      *sql.DB
	DbDriverName            string
	AutomigrateEnabled      bool
	DebugEnabled            bool
}

func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
		entityTableName:         opts.EntityTableName,
		attributeTableName:      opts.AttributeTableName,
		entityTrashTableName:    opts.EntityTrashTableName,
		attributeTrashTableName: opts.AttributeTrashTableName,
		automigrateEnabled:      opts.AutomigrateEnabled,
		db:                      opts.DB,
		dbDriverName:            opts.DbDriverName,
		debugEnabled:            opts.DebugEnabled,
	}

	if store.entityTableName == "" {
		return nil, errors.New("entity store: entityTableName is required")
	}

	if store.attributeTableName == "" {
		return nil, errors.New("entity store: attributeTableName is required")
	}

	if store.entityTrashTableName == "" {
		store.entityTrashTableName = store.entityTableName + "_trash"
	}

	if store.attributeTrashTableName == "" {
		store.attributeTrashTableName = store.attributeTableName + "_trash"
	}

	if store.db == nil {
		return nil, errors.New("entity store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = driverName(store.db)
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}
