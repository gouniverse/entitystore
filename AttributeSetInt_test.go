package entitystore

import "testing"

func TestAttributeInt(t *testing.T) {
	db := InitDB("test_attribute_int.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	errSet := store.AttributeSetInt("default", "test_int", 12)

	if errSet != nil {
		t.Fatal("Attribute could not be created:" + errSet.Error())
	}

	attr, err := store.AttributeFind("default", "test_int")

	if err != nil {
		t.Fatalf("Attribute could not be retrieved: " + err.Error())
	}

	if attr == nil {
		t.Fatalf("Attribute could not be retrieved")
	}

	v, _ := attr.GetInt()
	if v != 12 {
		t.Fatalf("Attribute value incorrect")
	}
}
