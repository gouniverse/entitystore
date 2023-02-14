package entitystore

import "testing"

func TestEntityTrash(t *testing.T) {
	db := InitDB("test_entity_trash.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	entity, err := store.EntityCreateWithAttributes("post", map[string]string{
		"title": "Test Post Title",
		"text":  "Test Post Text",
	})

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	attr, err := store.AttributeFind(entity.ID(), "title")

	if err != nil {
		t.Fatalf("Attribute could not be found: " + err.Error())
	}

	if attr == nil {
		t.Fatalf("Attribute should not be nil")
	}

	isDeleted, err := store.EntityTrash(entity.ID())

	if err != nil {
		t.Fatalf("Entiry could not be deleted: " + err.Error())
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

	attr, err = store.AttributeFind(entity.ID(), "title")

	if err != nil {
		t.Fatalf("Attribute could not be found: " + err.Error())
	}

	if attr != nil {
		t.Fatalf("Attribute should be nil")
	}
}
