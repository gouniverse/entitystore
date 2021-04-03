package entitystore

import (
	//"log"
	// "log"
	"testing"
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(filepath string) *gorm.DB /**sql.DB*/ {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})

	if err != nil {
		panic(err) 
	}
	
	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("entity_create.db")
	
	store := NewStore(WithGormDb(db),WithEntityTableName("cms_entity"),WithAttributeTableName("cms_attribute"),WithAutoMigrate(true))
	
	entity := store.EntityCreate("post")
	if entity == nil{
		t.Fatalf("Entity could not be created")
	}
}


func TestStoreAutomigrate(t *testing.T) {
	db := InitDB("entity_create.db")
	
	store := NewStore(WithGormDb(db),WithEntityTableName("cms_entity"),WithAttributeTableName("cms_attribute"))

	store.AutoMigrate()
	
	entity := store.EntityCreate("post")
	if entity == nil{
		t.Fatalf("Entity could not be created")
	}
}

func TestStoreEntityDelete(t *testing.T) {
	db := InitDB("entity_create.db")
	
	store := NewStore(WithGormDb(db),WithEntityTableName("cms_entity"),WithAttributeTableName("cms_attribute"),WithAutoMigrate(true))
	
	entity := store.EntityCreate("post")
	
	if entity == nil{
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
