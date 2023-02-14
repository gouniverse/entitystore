package entitystore

import "testing"

func TestEntityCreateWithAttributes(t *testing.T) {
	db := InitDB("test_entity_update.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

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
