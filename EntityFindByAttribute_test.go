package entitystore

import "testing"

func TestEntityFindByAttribute(t *testing.T) {
	db := InitDB("test_entity_update.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	entity, err := store.EntityCreateWithAttributes("post", map[string]string{
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
