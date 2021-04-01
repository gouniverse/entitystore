package entitystore

import (
	//"log"
	"testing"
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(filepath string) *gorm.DB /**sql.DB*/ {

  db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
  if err != nil { panic(err) }
  // Auto Migrate
//   db.AutoMigrate(&Entity{})
	// db, err := sql.Open("sqlite3", filepath)
	// if err != nil { panic(err) }
	// if db == nil { panic("db nil") }
	// return db
	return db
}

func TestEntityCreate(t *testing.T) {
	db := InitDB("entity_create.db")
	
	store := NewStore(WithGormDb(db),WithEntityTableName("cms_entity"),WithAttributeTableName("cms_attribute"),WithAutoMigrate(true))
	//  Init(Config{
	// 	DbInstance: db,
	// })
	entity := store.EntityCreate("post")
	if entity == nil{
		t.Fatalf("Entity could not be created")
	}
}

func TestEntityCreateWithAttributes(t *testing.T) {
	db := InitDB("entity_update.db")
	
	
	store := NewStore(WithGormDb(db),WithEntityTableName("cms_entity"),WithAttributeTableName("cms_attribute"),WithAutoMigrate(true))
	
	// Init(Config{
	// 	DbInstance: db,
	// })
	entity := store.EntityCreateWithAttributes("post", map[string]interface{}{
		"name":"Hello world",
	})
	if entity == nil{
		t.Fatalf("Entity could not be created")
	}

	// log.Println(entity)
	// log.Println(entity.GetAttribute("name"))
	// attribute := store.AttributeFind(entity.ID,"name")
	// log.Println(attribute)
	// attr1 := entity.GetAttribute("name")
	// log.Println(attr1)

	if entity.GetAny("name","") != "Hello world"{
		t.Fatalf("Entity attribute mismatch")
	}

	if store.AttributeFind(entity.ID,"name") == nil{
		t.Fatalf("Attribute NOT FOUND")
	}

	if store.AttributeFind(entity.ID,"name").GetAny() != "Hello world"{
		t.Fatalf("Entity attribute mismatch")
	}
}
