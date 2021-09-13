package entitystore

import (
	"testing"
	"time"
)

func TestEntityCreate(t *testing.T) {
	db := InitDB("test_entity_create.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, _ := store.EntityCreate("post")
	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	if len(entity.ID) < 32 {
		t.Fatalf("Entity ID:" + entity.ID + "is less than 32 characters")
	}

	if entity.CreatedAt.Before(time.Now().Add(-1 * time.Minute)) {
		t.Fatalf("Entity CreatedAt is not recent (before 1 min):" + entity.CreatedAt.String())
	}

	if entity.CreatedAt.After(time.Now().Add(1 * time.Minute)) {
		t.Fatalf("Entity CreatedAt is not recent (after 1 min):" + entity.CreatedAt.String())
	}


	if entity.UpdatedAt.Before(time.Now().Add(-1 * time.Minute)) {
		t.Fatalf("Entity UpdatedAt is not recent (before 1 min):" + entity.CreatedAt.String())
	}

	if entity.UpdatedAt.After(time.Now().Add(1 * time.Minute)) {
		t.Fatalf("Entity UpdateddAt is not recent (after 1 min):" + entity.CreatedAt.String())
	}
}

func TestEntityCreateWithAttributes(t *testing.T) {
	db := InitDB("test_entity_update.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	entity, err := store.EntityCreateWithAttributes("post", map[string]string{
		"name": "Hello world",
	})

	if err != nil {
		t.Fatalf("Entity could not be created" + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	val, _ := entity.GetString("name", "")
	if val != "Hello world" {
		t.Fatalf("Entity attribute mismatch")
	}
}

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

	isDeleted, err := store.EntityDelete(entity.ID)

	if err != nil {
		t.Fatalf("Entity could not be soft deleted: " + err.Error())
	}

	if isDeleted == false {
		t.Fatalf("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(entity.ID)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if val != nil {
		t.Fatalf("Entity should no longer be present")
	}
}

func TestEntityTrash(t *testing.T) {
	db := InitDB("test_entity_trash.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

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

	attr, err := store.AttributeFind(entity.ID, "title")

	if err != nil {
		t.Fatalf("Attribute could not be found: " + err.Error())
	}

	if attr == nil {
		t.Fatalf("Attribute should not be nil")
	}

	isDeleted, err := store.EntityTrash(entity.ID)

	if err != nil {
		t.Fatalf("Entiry could not be deleted: " + err.Error())
	}

	if isDeleted == false {
		t.Fatalf("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(entity.ID)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if val != nil {
		t.Fatalf("Entity should no longer be present")
	}

	attr, err = store.AttributeFind(entity.ID, "title")

	if err != nil {
		t.Fatalf("Attribute could not be found: " + err.Error())
	}

	if attr != nil {
		t.Fatalf("Attribute should be nil")
	}
}

func TestCreatingAttributes(t *testing.T) {
	db := InitDB("test_attributes_create.db")

	store, _ := NewStore(WithDb(db), WithEntityTableName("cms_entity"), WithAttributeTableName("cms_attribute"), WithAutoMigrate(true))

	//store.SetDebug(true)

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entity could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	// entity.SetString("title", "Product 1")
	// entity.SetFloat("price_float", 12.35)
	// entity.SetInt("price_int", 12)

	store.AttributeSetString(entity.ID, "description", "Description text")

	description, err := entity.GetString("description", "")

	if err != nil {
		t.Fatalf("Entiry could not be created: " + err.Error())
	}

	if description != "Description text" {
		t.Fatalf("Description is incorrect: " + description)
	}

}
