package entitystore

import "testing"

func TestEntityAttributesCreate(t *testing.T) {
	db := InitDB("test_attributes_create.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	//store.SetDebug(true)

	entity, err := store.EntityCreate("post")

	if err != nil {
		t.Fatalf("Entity could not be created: " + err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	entity.SetString("title", "Product 1")

	title, err := entity.GetString("title", "")

	if err != nil {
		t.Fatalf("Entity title could not be created: " + err.Error())
	}

	if title != "Product 1" {
		t.Fatal("Title is incorrect: ", title)
	}

	entity.SetFloat("price_float", 12.35)

	priceFloat, err := entity.GetFloat("price_float", 0)

	if err != nil {
		t.Fatalf("Entity price_float could not be created: " + err.Error())
	}

	if priceFloat != 12.35 {
		t.Fatal("Price float is incorrect: ", priceFloat)
	}

	entity.SetInt("price_int", 12)

	priceInt, err := entity.GetInt("price_int", 0)

	if err != nil {
		t.Fatalf("Entity price_int could not be created: " + err.Error())
	}

	if priceInt != 12 {
		t.Fatal("Price int is incorrect: ", priceInt)
	}

	// store.AttributeSetString(entity.ID(), "description", "Description text")

	entity.SetString("description", "Description text")
	description, err := entity.GetString("description", "")

	if err != nil {
		t.Fatalf("Entity description could not be created: " + err.Error())
	}

	if description != "Description text" {
		t.Fatalf("Description is incorrect: " + description)
	}

}
