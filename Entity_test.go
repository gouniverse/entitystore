package entitystore

import "testing"

func TestEntityAttributesCreate(t *testing.T) {
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
	entity.SetString("description", "Description text")

	// store.AttributeSetString(entity.ID(), "description", "Description text")

	description, err := entity.GetString("description", "")

	if err != nil {
		t.Fatalf("Entity description could not be created: " + err.Error())
	}

	if description != "Description text" {
		t.Fatalf("Description is incorrect: " + description)
	}

}
