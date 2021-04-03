package entitystore

import (
	//"log"
	"log"
	"testing"
	//"database/sql"
	// _ "github.com/mattn/go-sqlite3"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
)

func TestAttributeCreate(t *testing.T) {
	db := InitDB("entity_create.db")

	store := NewStore(WithGormDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	isOk := store.AttributeSetString("default", "hello", "world")

	if isOk == false {
		t.Fatalf("Attribute could not be created")
	}
}

func TestAttributeString(t *testing.T) {
	db := InitDB("entity_create.db")

	store := NewStore(WithGormDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	isOk := store.AttributeSetString("default", "hello", "world")

	if isOk == false {
		t.Fatalf("Attribute could not be created")
	}

	attr := store.AttributeFind("default", "hello")

	if attr == nil {
		t.Fatalf("Attribute could not be retrieved")
	}

	if attr.GetString() != "world" {
		t.Fatalf("Attribute value incorrect")
	}
}

func TestAttributeInt(t *testing.T) {
	db := InitDB("entity_create.db")

	store := NewStore(WithGormDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	isOk := store.AttributeSetInt("default", "test_int", 12)

	if isOk == false {
		t.Fatalf("Attribute could not be created")
	}

	attr := store.AttributeFind("default", "test_int")

	if attr == nil {
		t.Fatalf("Attribute could not be retrieved")
	}

	v, _ := attr.GetInt()
	if v != 12 {
		t.Fatalf("Attribute value incorrect")
	}
}

func TestAttributeFloat(t *testing.T) {
	db := InitDB("entity_create.db")

	store := NewStore(WithGormDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	isOk := store.AttributeSetFloat("default", "test_float", 12.123456789123456789123456789)

	if isOk == false {
		t.Fatalf("Attribute could not be created")
	}

	attr := store.AttributeFind("default", "test_float")

	if attr == nil {
		t.Fatalf("Attribute could not be retrieved")
	}

	v, _ := attr.GetFloat()
	if v != 12.123456789123456789123456789 {
		t.Fatalf("Attribute value incorrect")
	}
}

func TestAttributeInterface(t *testing.T) {
	db := InitDB("entity_create.db")

	store := NewStore(WithGormDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	isOk := store.AttributeSetInterface("default", "test_interface", "Hello world")

	if isOk == false {
		t.Fatalf("Attribute could not be created")
	}

	attr := store.AttributeFind("default", "test_interface")

	if attr == nil {
		t.Fatalf("Attribute could not be retrieved")
	}

	//v,_:=attr.GetFloat()
	v := attr.GetString()
	log.Println(v);
	// if v != 12.123456789123456789123456789 {
	// 	t.Fatalf("Attribute value incorrect")
	// }
}
