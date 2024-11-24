package entitystore

import "testing"

func TestEntityDelete(t *testing.T) {
	db := InitDB("test_entity_delete.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err.Error())
	}

	entity, err := store.EntityCreateWithType("post")

	if err != nil {
		t.Fatalf("Entity could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	err = entity.SetString("title", "Hello world")

	if err != nil {
		t.Fatalf("Entity title could not be created: " + err.Error())
	}

	isDeleted, err := store.EntityDelete(entity.ID())

	if err != nil {
		t.Fatalf("Entity could not be soft deleted: " + err.Error())
	}

	if isDeleted == false {
		t.Fatalf("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(entity.ID())

	if err != nil {
		t.Fatalf(err.Error())
	}

	if val != nil {
		t.Fatalf("Entity should no longer be present")
	}
}
