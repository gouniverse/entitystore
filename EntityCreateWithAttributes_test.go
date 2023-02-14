package entitystore

import "testing"

func TestEntityCreateWithAttributes(t *testing.T) {
	db := InitDB("test_entity_update.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	entity, err := store.EntityCreateWithAttributes("post", map[string]string{
		"name": "Hello world",
	})

	if err != nil {
		t.Fatal("Entity could not be created:", err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	val, _ := entity.GetString("name", "")
	if val != "Hello world" {
		t.Fatalf("Entity attribute mismatch")
	}
}
