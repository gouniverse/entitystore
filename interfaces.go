package entitystore

import "database/sql"

type StoreInterface interface {
	AutoMigrate() error

	GetAttributeTableName() string
	GetAttributeTrashTableName() string
	GetDB() *sql.DB
	GetEntityTableName() string
	GetEntityTrashTableName() string

	// AttributeCount(entityID string) uint64
	AttributeCreate(attr *Attribute) error
	AttributeCreateWithKeyAndValue(entityID string, attributeKey string, attributeValue string) (*Attribute, error)
	AttributeFind(entityID string, attributeKey string) (*Attribute, error)
	AttributeFindByHandle(entityID string, attributeKey string, attributeValue string) (*Attribute, error)
	AttributeList(options AttributeQueryOptions) ([]Attribute, error)
	AttributesSet(entityID string, attributes map[string]string) error
	AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) error
	AttributeSetInt(entityID string, attributeKey string, attributeValue int64) error
	AttributeSetString(entityID string, attributeKey string, attributeValue string) error
	// AttributeTrash(attr *Attribute) error

	EntityCount(options EntityQueryOptions) (int64, error)
	EntityCreate(entity *Entity) error
	EntityCreateWithType(entityType string) (*Entity, error)
	EntityCreateWithTypeAndAttributes(entityType string, attributes map[string]string) (*Entity, error)
	EntityDelete(entityID string) (bool, error)
	EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) (*Entity, error)
	EntityFindByHandle(entityType string, entityHandle string) (*Entity, error)
	EntityFindByID(entityID string) (*Entity, error)
	EntityList(options EntityQueryOptions) ([]Entity, error)
	EntityListByAttribute(entityType string, attributeKey string, attributeValue string) ([]Entity, error)
	EntityTrash(entityID string) (bool, error)
	EntityUpdate(entity Entity) error

	NewAttribute(opts NewAttributeOptions) Attribute
	NewAttributeFromMap(entityMap map[string]string) Attribute

	NewEntity(opts NewEntityOptions) Entity
	NewEntityFromMap(entityMap map[string]string) Entity
}
