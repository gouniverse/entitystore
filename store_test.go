package entitystore

import (
	//"log"
	// "log"
	"database/sql"
	"testing"

	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	dsn := filepath
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("test_entity_create.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}
}

func TestStoreAutomigrate(t *testing.T) {
	db := InitDB("test_entity_automigrate.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"))

	store.AutoMigrate()

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}
}

func TestStoreEntityDelete(t *testing.T) {
	db := InitDB("test_entity_delete.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	entity.SetString("title", "Hello world")

	isDeleted, err := store.EntityDelete(entity.ID)

	if err != nil {
		t.Fatalf("Entity could not be soft deleted: " + err.Error())
	}

	if isDeleted == false {
		t.Fatalf("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(entity.ID)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if val != nil {
		t.Fatalf("Entity should no longer be present")
	}
}

func TestStoreEntityTrash(t *testing.T) {
	db := InitDB("test_entity_trash.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	entity.SetString("title", "Hello world")

	isDeleted, err := store.EntityTrash(entity.ID)

	if err != nil {
		t.Fatalf("Entiry could not be deleted: " + err.Error())
	}

	if isDeleted == false {
		t.Fatalf("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(entity.ID)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if val != nil {
		t.Fatalf("Entity should no longer be present")
	}
}

func TestCreatingAttributes(t *testing.T) {
	db := InitDB("test_attributes_create.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	store.SetDebug(true)

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entity could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	// entity.SetString("title", "Product 1")
	// entity.SetFloat("price_float", 12.35)
	// entity.SetInt("price_int", 12)

	store.AttributeSetString(entity.ID, "description", "Description text")

	description, err := entity.GetString("description", "")

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if description != "Description text" {
		t.Fatalf("Description is incorrect: " + description)
	}

}
