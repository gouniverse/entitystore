package entitystore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

// Entity type
type entity struct {
	ID     string `gorm:"type:varchar(40);column:id;primary_key;"`
	Status string `gorm:"type:varchar(10);column:status;"`
	Type   string `gorm:"type:varchar(40);column:type;"`
	//Name        string     `gorm:"type:varchar(255);column:name;DEFAULT NULL;"`
	//Description string     `gorm:"type:longtext;column:description;"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`

	attributes []attribute `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	st         *Store      `gorm:-`
}

// BeforeCreate adds UID to model
func (e *entity) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	e.ID = uuid
	return nil
}

// Deletes the entity
func (e *entity) Delete() bool {
	return e.st.EntityDelete(e.ID)
}

// GetAny the value of the attribute as interface{} or the default value if it does not exist
func (e *entity) GetAny(attributeKey string, defaultValue interface{}) interface{} {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue
	}

	return attr.GetAny()
}

// GetInt the value of the attribute as string or the default value if it does not exist
func (e *entity) GetInt(attributeKey string, defaultValue int) (int, error) {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetInt()
}

// GetFloat the value of the attribute as float or the default value if it does not exist
func (e *entity) GetFloat(attributeKey string, defaultValue float32) (float32, error) {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetFloat()
}

// GetString the value of the attribute as string or the default value if it does not exist
func (e *entity) GetString(attributeKey string, defaultValue string) string {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue
	}

	return attr.GetString()
}

// GetAttribute the name of the User table
func (e *entity) GetAttribute(attributeKey string) *attribute {
	attr := &attribute{}

	result := e.st.db.Table(e.st.attributeTableName).First(&attr, "entity_id=? AND attribute_key=?", e.ID, attributeKey)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return attr
}

// SetString sets an attribute with string value
func (e *entity) SetAny(attributeKey string, attributeValue interface{}) bool {
	bytes, err := json.Marshal(attributeValue)

	if err != nil {
		return false
	}

	strValue := string(bytes)

	return e.st.AttributeSet(e.ID, attributeKey, strValue)
}

// SetString sets an attribute with string value
func (e *entity) SetFloat(attributeKey string, attributeValue float32) bool {
	return e.st.AttributeSet(e.ID, attributeKey, fmt.Sprint(attributeValue))
}

// SetString sets an attribute with string value
func (e *entity) SetInt(attributeKey string, attributeValue int) bool {
	return e.st.AttributeSet(e.ID, attributeKey, fmt.Sprint(attributeValue))
}

// SetString sets an attribute with string value
func (e *entity) SetString(attributeKey string, attributeValue string) bool {
	return e.st.AttributeSet(e.ID, attributeKey, attributeValue)
}

// EntityCount counts entities
func (st *Store) EntityCount(entityType string) uint64 {
	var count int64
	st.db.Table(st.entityTableName).Where("type=?", entityType).Count(&count)
	return uint64(count)
}

// EntityCreate creates a new entity
func (st *Store) EntityCreate(entityType string) *entity {
	entity := &entity{Type: entityType, Status: "active", st: st}

	dbResult := st.db.Table(st.entityTableName).Create(&entity)

	if dbResult.Error != nil {
		return nil
	}

	return entity
}

// EntityDelete deletes an entity and all attributes
func (st *Store) EntityDelete(entityID string) bool {
	if entityID == "" {
		return false
	}

	// Note the use of tx as the database handle once you are within a transaction
	tx := st.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println(err)
		return false
	}

	if err := tx.Where("entity_id=?", entityID).Table(st.attributeTableName).Delete(&attribute{}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	if err := tx.Where("id=?", entityID).Table(st.entityTableName).Delete(&entity{}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	err := tx.Commit().Error

	if err == nil {
		return true
	}

	log.Println(err)

	return false
}

// EntityFindByID finds an entity by ID
func (st *Store) EntityFindByID(entityID string) *Entity {
	if entityID == "" {
		return nil
	}

	entity := &Entity{}

	resultEntity := st.db.Table(st.entityTableName).First(&entity, "id=?", entityID)

	if resultEntity.Error != nil {
		if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	return entity
}

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) *attribute {
	attr := &attribute{EntityID: entityID, AttributeKey: attributeKey, AttributeValue: attributeValue}

	dbResult := st.db.Table(st.attributeTableName).Create(&attr)

	if dbResult.Error != nil {
		return nil
	}

	return attr
}

// AttributeFind finds an entity by ID
func (st *Store) AttributeFind(entityID string, attributeKey string) *attribute {
	attr := &attribute{}

	result := st.db.Table(st.attributeTableName).First(&attr, "entity_id=? AND attribute_key=?", entityID, attributeKey)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	return attr
}

// AttributeGet the name of the User table
func (st *Store) AttributeGet(entityID string, attributeKey string) *attribute {
	attr := &attribute{}

	result := st.db.Table(st.attributeTableName).First(&attr, "entity_id=? AND attribute_key=?", entityID, attributeKey)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return attr
}

// AttributeSet creates a new entity
func (st *Store) AttributeSet(entityID string, attributeKey string, attributeValue string) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreate(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.AttributeValue = attributeValue
	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// Attribute type
type attribute struct {
	ID             string     `gorm:"type:varchar(40);column:id;primary_key;"`
	EntityID       string     `gorm:"type:varchar(40);column:entity_id;"`
	AttributeKey   string     `gorm:"type:varchar(255);column:attribute_key;DEFAULT NULL;"`
	AttributeValue string     `gorm:"type:longtext;column:attribute_value;"`
	CreatedAt      time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt      time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt      *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`
}

// BeforeCreate adds UID to model
func (a *attribute) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	a.ID = uuid
	return nil
}

// GetValue de-serializes the values
func (a *attribute) GetAny() interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(a.AttributeValue), &value)

	if err != nil {
		panic("JSOB error unmarshaliibg attribute" + err.Error())
	}

	return value
}

// GetString returns the value as string
func (a *attribute) GetInt() (int, error) {
	return strconv.Atoi(a.AttributeValue)
}

// GetFloat returns the value as string
func (a *attribute) GetFloat() (float32, error) {
	f64Value, err := strconv.ParseFloat(a.AttributeValue, 32)
	return float32(f64Value), err
}

// GetString returns the value as string
func (a *attribute) GetString() string {
	return a.AttributeValue
}

// SetAny serializes the values
func (a *attribute) SetAny(value interface{}) bool {
	bytes, err := json.Marshal(value)

	if err != nil {
		return false
	}

	a.AttributeValue = string(bytes)

	return true
}

// SetString serializes the values
func (a *attribute) SetString(value string) bool {
	a.AttributeValue = value

	return true
}
