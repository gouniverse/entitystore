package entitystore

import (
	"testing"
	"time"
)

func TestEntityCreate(t *testing.T) {
	db := InitDB("test_entity_create.db")

	store, _ := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	entity, _ := store.EntityCreate("post")
	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	if len(entity.ID()) < 32 {
		t.Fatalf("Entity ID:" + entity.ID() + "is less than 32 characters")
	}

	if entity.CreatedAt().Before(time.Now().Add(-1 * time.Minute)) {
		t.Fatalf("Entity CreatedAt is not recent (before 1 min):" + entity.CreatedAt().String())
	}

	if entity.CreatedAt().After(time.Now().Add(1 * time.Minute)) {
		t.Fatalf("Entity CreatedAt is not recent (after 1 min):" + entity.CreatedAt().String())
	}

	if entity.UpdatedAt().Before(time.Now().Add(-1 * time.Minute)) {
		t.Fatalf("Entity UpdatedAt is not recent (before 1 min):" + entity.CreatedAt().String())
	}

	if entity.UpdatedAt().After(time.Now().Add(1 * time.Minute)) {
		t.Fatalf("Entity UpdateddAt is not recent (after 1 min):" + entity.CreatedAt().String())
	}
}
