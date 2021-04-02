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
