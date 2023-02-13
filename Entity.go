package entitystore

import (
	"log"
	"time"
)

// Entity this is the type for an Entity
type Entity struct {
	id           string
	entityType   string
	entityHandle string
	createdAt    time.Time
	updatedAt    time.Time
	st           *Store
}

func (e *Entity) ToMap() map[string]any {
	entry := map[string]any{}
	entry["id"] = e.ID()
	entry["entity_type"] = e.Type()
	entry["entity_handle"] = e.Handle()
	entry["created_at"] = e.CreatedAt()
	entry["updated_at"] = e.UpdatedAt()
	return entry
}

func (st *Store) FromMap(entityMap map[string]any) *Entity {
	opts := NewEntityOptions{}
	if id, exists := entityMap["id"]; exists {
		opts.ID = id.(string)
	}
	if entityType, exists := entityMap["entity_type"]; exists {
		opts.Type = entityType.(string)
	}
	if entityHandle, exists := entityMap["entity_handle"]; exists {
		opts.Handle = entityHandle.(string)
	}
	if createdAt, exists := entityMap["created_at"]; exists {
		opts.CreatedAt = createdAt.(time.Time)
	}
	if updatedAt, exists := entityMap["updated_at"]; exists {
		opts.UpdatedAt = updatedAt.(time.Time)
	}
	entity := st.NewEntity(opts)
	return entity
}

func (e Entity) ID() string {
	return e.id
}

func (e *Entity) Type() string {
	return e.entityType
}

func (e *Entity) Handle() string {
	return e.entityHandle
}

func (e *Entity) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entity) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Entity) SetID(id string) *Entity {
	e.id = id
	return e
}

func (e *Entity) SetType(entityType string) *Entity {
	e.entityType = entityType
	return e
}

func (e *Entity) SetHandle(handle string) *Entity {
	e.entityHandle = handle
	return e
}

func (e *Entity) SetCreatedAt(createdAt time.Time) *Entity {
	e.createdAt = createdAt
	return e
}

func (e *Entity) SetUpdatedAt(updatedAt time.Time) *Entity {
	e.updatedAt = updatedAt
	return e
}

type NewEntityOptions struct {
	ID        string
	Type      string
	Handle    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (st *Store) NewEntity(opts NewEntityOptions) *Entity {
	entity := Entity{}
	entity.SetID(opts.ID)
	entity.SetType(opts.Type)
	entity.SetHandle(opts.Handle)
	entity.SetCreatedAt(opts.CreatedAt)
	entity.SetUpdatedAt(opts.UpdatedAt)
	entity.st = st
	return &entity
}

// BeforeCreate adds UID to model
// func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
// 	uuid := uid.HumanUid()
// 	e.ID = uuid
// 	return nil
// }

// Delete deletes the entity
// func (e *Entity) Delete() bool {
// 	return e.st.EntityDelete(e.ID)
// }

// Trash moves the entity to the trash bin
// func (e *Entity) Trash() bool {
// 	return e.st.EntityTrash(e.ID)
// }

// GetInt the value of the attribute as string or the default value if it does not exist
func (e *Entity) GetInt(attributeKey string, defaultValue int64) (int64, error) {
	attr, err := e.GetAttribute(attributeKey)

	if err != nil {
		if e.st.GetDebug() {
			log.Println(err)
		}
		return defaultValue, err
	}

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetInt()
}

// GetAttribute return specified attribute
func (e *Entity) GetAttribute(attributeKey string) (*Attribute, error) {
	return e.st.AttributeFind(e.ID(), attributeKey)
}

// GetAttributes all the attributes of the entity
func (e *Entity) GetAttributes() ([]Attribute, error) {
	return e.st.EntityAttributeList(e.ID())
}

// GetFloat the value of the attribute as float or the default value if it does not exist
func (e *Entity) GetFloat(attributeKey string, defaultValue float64) (float64, error) {
	attr, err := e.GetAttribute(attributeKey)

	if err != nil {
		if e.st.GetDebug() {
			log.Println(err)
		}
		return defaultValue, err
	}

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetFloat()
}

// GetString the value of the attribute as string or the default value if it does not exist
func (e *Entity) GetString(attributeKey string, defaultValue string) (string, error) {
	attr, err := e.GetAttribute(attributeKey)

	if err != nil {
		if e.st.GetDebug() {
			log.Println(err)
		}
		return defaultValue, err
	}

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetString(), nil
}

// SetAll upserts the attributes
func (e *Entity) SetAll(attributes map[string]string) (bool, error) {
	return e.st.AttributesSet(e.ID(), attributes)
}

// SetFloat sets an attribute with float value
func (e *Entity) SetFloat(attributeKey string, attributeValue float64) (bool, error) {
	return e.st.AttributeSetFloat(e.ID(), attributeKey, attributeValue)
}

// SetInt sets an attribute with int value
func (e *Entity) SetInt(attributeKey string, attributeValue int64) (bool, error) {
	return e.st.AttributeSetInt(e.ID(), attributeKey, attributeValue)
}

// SetString sets an attribute with string value
func (e *Entity) SetString(attributeKey string, attributeValue string) (bool, error) {
	return e.st.AttributeSetString(e.ID(), attributeKey, attributeValue)
}
