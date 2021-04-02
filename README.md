# Entity Store

Saves data in SQL database in a "schemaless" way

## Installation
```
go get -u github.com/gouniverse/entitystore
```

## Setup

```
entityStore = entitystore.NewStore(entitystore.WithGormDb(databaseInstance), entitystore.WithEntityTableName("entities_entity"), entitystore.WithAttributeTableName("entities_attribute"), entitystore.WithAutoMigrate(true))
```

## Usage

```
person := entityStore.EntityCreate("person")
person.SetString("name","Jon Doe")
person.SetInt("age", 32)
```

## Methods

These methods may be subject to change

### Store Methods

- AutoMigrate() - auto migrate
- EntityCount(entityType string) uint64 - counts entities
- EntityCreate(entityType string) *Entity - creates a new entity
- EntityCreateWithAttributes(entityType string, attributes map[string]interface{}) *Entity
- EntityDelete(entityID string) - deletes an entity and all attributes
- EntityFindByID(entityID string) *Entity - finds an entity by ID
- EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) *Entity - finds an entity by attribute
- EntityList(entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) []Entity - lists entities
- EntityListByAttribute(entityType string, attributeKey string, attributeValue string) []Entity - finds an entity by attribute
- AttributeCreate(entityID string, attributeKey string, attributeValue string) *Attribute - creates a new attribute
- AttributeFind(entityID string, attributeKey string) *Attribute - finds an attribute by ID
- AttributeGet(entityID string, attributeKey string) *Attribute - finds an attribute by ID
- AttributeSet(entityID string, attributeKey string, attributeValue string) bool -  creates a new attribute
- AttributesSet(entityID string, attributes map[string]interface{}) bool - creates multiple attributes

### Entity Methods

- Delete() bool - deletes the entity
- GetInt(attributeKey string, defaultValue int) (int, error) - the value of the attribute as string or the default value if it does not exist
- GetFloat(attributeKey string, defaultValue float32) (float32, error) - the value of the attribute as float or the default value if it does not exist
- GetInterface(attributeKey string, defaultValue interface{}) interface{} - the value of the attribute as interface{} or the default value if it does not exist
- GetString(attributeKey string, defaultValue string) string - the value of the attribute as string or the default value if it does not exist
- GetAttribute(attributeKey string) *Attribute - returns an attribute by key
- SetAllAny(attributes map[string]interface{}) bool - upserts the attributes
- SetFloat(attributeKey string, attributeValue float32) bool - sets an attribute with float value
- SetInt(attributeKey string, attributeValue int) bool - sets an attribute with int value
- SetInterface(attributeKey string, attributeValue interface{}) bool - sets an attribute with string value
- SetString(attributeKey string, attributeValue string) bool - sets an attribute with string value
