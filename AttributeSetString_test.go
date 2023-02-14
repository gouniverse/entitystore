package entitystore

import (
	"testing"
)

func TestAttributeString(t *testing.T) {
	db := InitDB("test_attribute_string.db")

	store, err := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	if err != nil {
		t.Fatalf(err.Error())
	}

	errSetString := store.AttributeSetString("default", "hello", "world")

	if errSetString != nil {
		t.Fatalf("Attribute could not be created: " + err.Error())
	}

	// store.EnableDebug(true)

	attr, err := store.AttributeFind("default", "hello")

	if err != nil {
		t.Fatalf("Attribute could not be retrieved: " + err.Error())
	}

	if attr == nil {
		t.Fatalf("Attribute could not be retrieved")
	}

	if attr.GetString() != "world" {
		t.Fatal("Attribute value incorrect", "must be 'world'", "found", attr.GetString())
	}
}
