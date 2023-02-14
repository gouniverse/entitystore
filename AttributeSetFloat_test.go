package entitystore

import "testing"

func TestAttributeSetFloat(t *testing.T) {
	db := InitDB("test_attribute_float.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	err := store.AttributeSetFloat("default", "test_float", 12.123456789123456789123456789)

	if err != nil {
		t.Fatal("Attribute could not be created:", err.Error())
	}

	attr, err := store.AttributeFind("default", "test_float")

	if err != nil {
		t.Fatal("Attribute could not be retrieved:", err.Error())
	}

	if attr == nil {
		t.Fatal("Attribute could not be retrieved")
	}

	v, _ := attr.GetFloat()
	if v != 12.123456789123456789123456789 {
		t.Fatal("Attribute value incorrect")
	}
}
