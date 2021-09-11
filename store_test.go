package entitystore

import (
	//"log"
	// "log"
	"database/sql"
	"testing"

	//"database/sql"
	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	dsn := filepath
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("test_entity_create.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, _ := store.EntityCreate("post")
	if entity == nil {
		t.Fatalf("Entity could not be created")
	}
}

func TestStoreAutomigrate(t *testing.T) {
	db := InitDB("test_entity_automigrate.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"))

	store.AutoMigrate()

	entity, _ := store.EntityCreate("post")
	if entity == nil {
		t.Fatalf("Entity could not be created")
	}
}

func TestStoreEntityDelete(t *testing.T) {
	db := InitDB("test_entity_delete.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, _ := store.EntityCreate("post")

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	entity.SetString("title", "Hello world")

	isDeleted := store.EntityDelete(entity.ID)

	if isDeleted == false {
		t.Fatalf("Entity could not be soft deleted")
	}

	if store.EntityFindByID(entity.ID) != nil {
		t.Fatalf("Entity should no longer be present")
	}
}

func TestStoreEntityTrash(t *testing.T) {
	db := InitDB("test_entity_trash.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, _ := store.EntityCreate("post")

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	entity.SetString("title", "Hello world")

	isDeleted := store.EntityTrash(entity.ID)

	if isDeleted == false {
		t.Fatalf("Entity could not be soft deleted")
	}

	if store.EntityFindByID(entity.ID) != nil {
		t.Fatalf("Entity should no longer be present")
	}
}
