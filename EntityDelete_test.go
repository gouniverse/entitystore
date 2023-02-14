package entitystore

import "testing"

func TestEntityDelete(t *testing.T) {
	db := InitDB("test_entity_delete.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	entity.SetString("title", "Hello world")

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
