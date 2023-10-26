package entitystore

import "testing"

func TestEntityFindByAttribute(t *testing.T) {
	db := InitDB("test_entity_find_by_attribute.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err.Error())
	}

	entity, err := store.EntityCreateWithTypeAndAttributes("post", map[string]string{
		"path": "/",
	})

	if err != nil {
		t.Fatalf("Entity could not be created" + err.Error())
	}

	val, _ := entity.GetString("path", "")
	if val != "/" {
		t.Fatalf("Entity attribute mismatch")
	}

	// store.SetDebug(true)

	homePage, err := store.EntityFindByAttribute("post", "path", "/")

	if err != nil {
		t.Fatalf("Entity find by attribute failed: " + err.Error())
	}

	if homePage == nil {
		t.Fatalf("Entity could not be found")
	}
}
