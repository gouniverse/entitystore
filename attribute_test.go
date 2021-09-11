package entitystore

import "testing"

//"database/sql"
// _ "github.com/mattn/go-sqlite3"
// "gorm.io/driver/sqlite"
// "gorm.io/gorm"

func TestAttributeCreate(t *testing.T) {
	db := InitDB("entity_create.db")

	store, err := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	if err != nil {
		t.Fatalf(err.Error())
	}

	isOk, err := store.AttributeSetString("default", "hello", "world")

	if err != nil {
		t.Fatalf("Attribute could not be created: " + err.Error())
	}

	if isOk == false {
		t.Fatalf("Attribute could not be created")
	}
}

// func TestAttributeString(t *testing.T) {
// 	db := InitDB("entity_create.db")

// 	store, err := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	isOk, err := store.AttributeSetString("default", "hello", "world")

// 	if err != nil {
// 		t.Fatalf("Attribute could not be created: " + err.Error())
// 	}

// 	if isOk == false {
// 		t.Fatalf("Attribute could not be created")
// 	}

// 	attr, err := store.AttributeFind("default", "hello")

// 	if err != nil {
// 		t.Fatalf("Attribute could not be retrieved: " + err.Error())
// 	}

// 	if attr == nil {
// 		t.Fatalf("Attribute could not be retrieved")
// 	}

// 	if attr.GetString() != "world" {
// 		t.Fatalf("Attribute value incorrect")
// 	}
// }

// func TestAttributeInt(t *testing.T) {
// 	db := InitDB("entity_create.db")

// 	store, err := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	isOk, err := store.AttributeSetInt("default", "test_int", 12)

// 	if err != nil {
// 		t.Fatalf("Attribute could not be created:" + err.Error())
// 	}

// 	if isOk == false {
// 		t.Fatalf("Attribute could not be created")
// 	}

// 	attr, err := store.AttributeFind("default", "test_int")

// 	if err != nil {
// 		t.Fatalf("Attribute could not be retrieved: " + err.Error())
// 	}

// 	if attr == nil {
// 		t.Fatalf("Attribute could not be retrieved")
// 	}

// 	v, _ := attr.GetInt()
// 	if v != 12 {
// 		t.Fatalf("Attribute value incorrect")
// 	}
// }

// func TestAttributeFloat(t *testing.T) {
// 	db := InitDB("entity_create.db")

// 	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

// 	isOk, _ := store.AttributeSetFloat("default", "test_float", 12.123456789123456789123456789)

// 	if isOk == false {
// 		t.Fatalf("Attribute could not be created")
// 	}

// 	attr, err := store.AttributeFind("default", "test_float")

// 	if err != nil {
// 		t.Fatalf("Attribute could not be retrieved" + err.Error())
// 	}

// 	if attr == nil {
// 		t.Fatalf("Attribute could not be retrieved")
// 	}

// 	v, _ := attr.GetFloat()
// 	if v != 12.123456789123456789123456789 {
// 		t.Fatalf("Attribute value incorrect")
// 	}
// }

// func TestAttributeInterface(t *testing.T) {
// 	db := InitDB("entity_create.db")

// 	store, err := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	isOk, err := store.AttributeSetInterface("default", "test_interface", "Hello world")

// 	if err != nil {
// 		t.Fatalf("Attribute could not be created: " + err.Error())
// 	}

// 	if isOk == false {
// 		t.Fatalf("Attribute could not be created")
// 	}

// 	attr, err := store.AttributeFind("default", "test_interface")

// 	if err != nil {
// 		t.Fatalf("Attribute could not be retrieved: " + err.Error())
// 	}

// 	if attr == nil {
// 		t.Fatalf("Attribute could not be retrieved")
// 	}

// 	//v,_:=attr.GetFloat()
// 	v := attr.GetString()
// 	log.Println(v)
// 	// if v != 12.123456789123456789123456789 {
// 	// 	t.Fatalf("Attribute value incorrect")
// 	// }
// }
