package entitystore

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"log"
	//"strconv"
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

const (
	// EntityStatusActive entity "active" status
	EntityStatusActive = "active"
	// EntityStatusInactive entity "inactive" status
	EntityStatusInactive = "inactive"
)

// Entity type
type Entity struct {
	ID     string `gorm:"type:varchar(40);column:id;primary_key;"`
	Status string `gorm:"type:varchar(10);column:status;"`
	Type   string `gorm:"type:varchar(40);column:type;"`
	//Name        string     `gorm:"type:varchar(255);column:name;DEFAULT NULL;"`
	//Description string     `gorm:"type:longtext;column:description;"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`

	attributes []Attribute `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	st         *Store      `gorm:-`
}

// BeforeCreate adds UID to model
func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	e.ID = uuid
	return nil
}

// Delete deletes the entity
func (e *Entity) Delete() bool {
	return e.st.EntityDelete(e.ID)
}

// GetAny the value of the attribute as interface{} or the default value if it does not exist
func (e *Entity) GetAny(attributeKey string, defaultValue interface{}) interface{} {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue
	}

	return attr.GetInterface()
}

// GetInt the value of the attribute as string or the default value if it does not exist
func (e *Entity) GetInt(attributeKey string, defaultValue int) (int, error) {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetInt()
}

// GetFloat the value of the attribute as float or the default value if it does not exist
func (e *Entity) GetFloat(attributeKey string, defaultValue float64) (float64, error) {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetFloat()
}

// GetString the value of the attribute as string or the default value if it does not exist
func (e *Entity) GetString(attributeKey string, defaultValue string) string {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue
	}

	return attr.GetString()
}

// GetAttribute the name of the User table
func (e *Entity) GetAttribute(attributeKey string) *Attribute {
	attr := &Attribute{}

	result := e.st.db.Table(e.st.attributeTableName).First(&attr, "entity_id=? AND attribute_key=?", e.ID, attributeKey)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return attr
}

// SetAllAny upserts the attributes
func (e *Entity) SetAllAny(attributes map[string]interface{}) bool {
	return e.st.AttributesSet(e.ID, attributes)
}

// SetInterface sets an attribute with string value
func (e *Entity) SetInterface(attributeKey string, attributeValue interface{}) bool {
	return e.st.AttributeSetInterface(e.ID, attributeKey, attributeValue)
}

// SetFloat sets an attribute with float value
func (e *Entity) SetFloat(attributeKey string, attributeValue float64) bool {
	return e.st.AttributeSetFloat(e.ID, attributeKey, attributeValue)
}

// SetInt sets an attribute with int value
func (e *Entity) SetInt(attributeKey string, attributeValue int64) bool {
	return e.st.AttributeSetInt(e.ID, attributeKey, attributeValue)
}

// SetString sets an attribute with string value
func (e *Entity) SetString(attributeKey string, attributeValue string) bool {
	return e.st.AttributeSetString(e.ID, attributeKey, attributeValue)
}
