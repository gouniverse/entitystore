package entitystore

import (
	//"log"
	// "log"
	"database/sql"
	"os"
	"testing"

	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("test_store_create.db")

	store, err := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}
}

func TestStoreAutomigrate(t *testing.T) {
	db := InitDB("test_entity_automigrate.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"))

	err := store.AutoMigrate()

	if err != nil {
		t.Fatalf("Automigrate failed: " + err.Error())
	}
}
