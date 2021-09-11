package entitystore

import "testing"

//"log"

//"database/sql"
// _ "github.com/mattn/go-sqlite3"
// "gorm.io/driver/sqlite"
// "gorm.io/gorm"

func TestEntityCreate(t *testing.T) {
	db := InitDB("entity_create.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))
	//  Init(Config{
	// 	DbInstance: db,
	// })
	entity, _ := store.EntityCreate("post")
	if entity == nil {
		t.Fatalf("Entity could not be created")
	}
}

func TestEntityCreateWithAttributes(t *testing.T) {
	db := InitDB("entity_update.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	// Init(Config{
	// 	DbInstance: db,
	// })
	entity := store.EntityCreateWithAttributes("post", map[string]interface{}{
		"name": "Hello world",
	})
	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	// log.Println(entity)
	// log.Println(entity.GetAttribute("name"))
	// attribute := store.AttributeFind(entity.ID,"name")
	// log.Println(attribute)
	// attr1 := entity.GetAttribute("name")
	// log.Println(attr1)

	if entity.GetAny("name", "") != "Hello world" {
		t.Fatalf("Entity attribute mismatch")
	}

	// attr, err := store.AttributeFind(entity.ID, "name")

	// if err == nil {
	// 	t.Fatalf("Attribute could not be retrieved" + err.Error())
	// }

	// if attr == nil {
	// 	t.Fatalf("Attribute NOT FOUND")
	// }

	// if attr.GetInterface() != "Hello world" {
	// 	t.Fatalf("Entity attribute mismatch")
	// }
}
