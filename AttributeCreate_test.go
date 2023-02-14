package entitystore

import "testing"

func TestAttributeCreate(t *testing.T) {
	db := InitDB("test_attribute_create.db")

	store, err := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true), WithDebug(true))

	if err != nil {
		t.Fatalf(err.Error())
	}

	errSet := store.AttributeSetString("default", "hello", "world")

	if errSet != nil {
		t.Fatal("Attribute could not be created:", errSet.Error())
	}
}
