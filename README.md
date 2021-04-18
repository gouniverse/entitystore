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
person.SetFloat("salary", 1234.56)
person.SetInterface("kids", []string{"Tina","Sam"})
```

## Methods

These methods may be subject to change

### Store Methods


- AttributeCreate(entityID string, attributeKey string, attributeValue string) *Attribute - creates a new attribute
- AttributeFind(entityID string, attributeKey string) *Attribute - finds an attribute by ID
- AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) bool - upserts a new float attribute
- AttributeSetInt(entityID string, attributeKey string, attributeValue int64) bool -  upserts a new int attribute
- AttributeSetString(entityID string, attributeKey string, attributeValue string) bool -  upserts a new interface{} attribute
- AttributeSetString(entityID string, attributeKey string, attributeValue string) bool -  upserts a new string attribute
- AutoMigrate() - auto migrate
- EntityCount(entityType string) uint64 - counts entities
- EntityCreate(entityType string) *Entity - creates a new entity
- EntityCreateWithAttributes(entityType string, attributes map[string]interface{}) *Entity
- EntityDelete(entityID string) - deletes an entity and all attributes
- (DEPERECATED) EntityDeleteSoft(entityID string) - soft deletes an entity and all attributes (DEPERECATED)
- EntityFindByID(entityID string) *Entity - finds an entity by ID
- EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) *Entity - finds an entity by attribute
- EntityList(entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) []Entity - lists entities
- EntityListByAttribute(entityType string, attributeKey string, attributeValue string) []Entity - finds an entity by attribute
- EntityTrash(entityID string) - moves an entity and all its attributes to the trash bin
- GetAttributeTableName() string
- GetAttributeTrashTableName() string
- GetDB() *gorm.DB
- GetEntityTableName() string
- GetEntityTrashTableName() string


### Entity Methods

- Delete() bool - deletes the entity
- GetInt(attributeKey string, defaultValue int64) (int64, error) - the value of the attribute as string or the default value if it does not exist
- GetFloat(attributeKey string, defaultValue float64) (float64, error) - the value of the attribute as float or the default value if it does not exist
- GetInterface(attributeKey string, defaultValue interface{}) interface{} - the value of the attribute as interface{} or the default value if it does not exist
- GetString(attributeKey string, defaultValue string) string - the value of the attribute as string or the default value if it does not exist
- GetAttribute(attributeKey string) *Attribute - returns an attribute by key
- SetFloat(attributeKey string, attributeValue float64) bool - sets an attribute with float value
- SetInt(attributeKey string, attributeValue int64) bool - sets an attribute with int value
- SetInterface(attributeKey string, attributeValue interface{}) bool - sets an attribute with string value
- SetString(attributeKey string, attributeValue string) bool - sets an attribute with string value

### Attribute Methods

- GetInterface() interface{} - de-serializes the JSON value
- GetInt() (int64, error) - returns the value as int
- GetFloat() (float64, error) - returns the value as float
- GetString() string - returns the value as string
- SetFloat(value float64) bool - saves a float value
- SetInt(value int64) bool - saves a int value
- SetInterface(value interface{}) bool - serializes the interface to JSON string and saves it
- SetString(value string) bool - saves a string value
