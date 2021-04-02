package entitystore

import (
	"encoding/json"
	"errors"
	"fmt"
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
func (e *Entity) GetFloat(attributeKey string, defaultValue float32) (float32, error) {
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
	bytes, err := json.Marshal(attributeValue)

	if err != nil {
		return false
	}

	strValue := string(bytes)

	return e.st.AttributeSet(e.ID, attributeKey, strValue)
}

// SetFloat sets an attribute with float value
func (e *Entity) SetFloat(attributeKey string, attributeValue float32) bool {
	return e.st.AttributeSet(e.ID, attributeKey, fmt.Sprint(attributeValue))
}

// SetInt sets an attribute with int value
func (e *Entity) SetInt(attributeKey string, attributeValue int) bool {
	return e.st.AttributeSet(e.ID, attributeKey, fmt.Sprint(attributeValue))
}

// SetString sets an attribute with string value
func (e *Entity) SetString(attributeKey string, attributeValue string) bool {
	return e.st.AttributeSet(e.ID, attributeKey, attributeValue)
}

// EntityCount counts entities
func (st *Store) EntityCount(entityType string) uint64 {
	var count int64
	st.db.Table(st.entityTableName).Where("type=?", entityType).Count(&count)
	return uint64(count)
}

// EntityCreate creates a new entity
func (st *Store) EntityCreate(entityType string) *Entity {
	entity := &Entity{Type: entityType, Status: "active", st: st}

	dbResult := st.db.Table(st.entityTableName).Create(&entity)

	if dbResult.Error != nil {
		return nil
	}

	return entity
}

// EntityCreateWithAttributes func
func (st *Store) EntityCreateWithAttributes(entityType string, attributes map[string]interface{}) *Entity {
	// Note the use of tx as the database handle once you are within a transaction
	tx := st.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil
	}

	//return tx.Commit().Error

	entity := &Entity{Type: entityType, Status: EntityStatusActive, st: st}

	dbResult := tx.Table(st.entityTableName).Create(&entity)

	if dbResult.Error != nil {
		tx.Rollback()
		return nil
	}

	//entityAttributes := make([]EntityAttribute, 0)
	for k, v := range attributes {
		ea := Attribute{EntityID: entity.ID, AttributeKey: k} //, AttributeValue: value}
		ea.SetInterface(v)

		dbResult2 := tx.Table(st.attributeTableName).Create(&ea)
		if dbResult2.Error != nil {
			tx.Rollback()
			return nil
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
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

	if err := tx.Where("entity_id=?", entityID).Table(st.attributeTableName).Delete(&Attribute{}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	if err := tx.Where("id=?", entityID).Table(st.entityTableName).Delete(&Entity{}).Error; err != nil {
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

// EntityFindByAttribute finds an entity by attribute
func (st *Store) EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) *Entity {
	attr := &Attribute{}

	subQuery := st.db.Table(st.entityTableName).Model(&Entity{}).Select("id").Where("type = ?", entityType)
	result := st.db.Table(st.attributeTableName).First(&attr, "entity_id IN (?) AND attribute_key=? AND attribute_value=?", subQuery, attributeKey, attributeValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	// DEBUG: log.Println(entityAttribute)

	ent := &Entity{}

	resultEntity := st.db.Table(st.entityTableName).First(&ent, "id=?", attr.EntityID)

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	return ent
}

// EntityList lists entities
func (st *Store) EntityList(entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) []Entity {
	entityList := []Entity{}
	result := st.db.Table(st.entityTableName).Where("type=?", entityType).Order(orderBy + " " + sort).Offset(int(offset)).Limit(int(perPage)).Find(&entityList)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return entityList
}

// EntityListByAttribute finds an entity by attribute
func (st *Store) EntityListByAttribute(entityType string, attributeKey string, attributeValue string) []Entity {
	//entityAttributes := []EntityAttribute{}
	var entityIDs []string

	subQuery := st.db.Table(st.entityTableName).Model(&Entity{}).Select("id").Where("type = ?", entityType)
	result := st.db.Table(st.attributeTableName).Model(&Attribute{}).Select("entity_id").Find(&entityIDs, "entity_id IN (?) AND attribute_key=? AND attribute_value=?", subQuery, attributeKey, attributeValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	// DEBUG: log.Println(result)

	entities := []Entity{}

	resultEntity := st.db.Table(st.entityTableName).Where("id IN (?)", entityIDs).Find(&entities)

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	return entities
}

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) *Attribute {
	attr := &Attribute{EntityID: entityID, AttributeKey: attributeKey, AttributeValue: attributeValue}

	dbResult := st.db.Table(st.attributeTableName).Create(&attr)

	if dbResult.Error != nil {
		return nil
	}

	return attr
}

// AttributeFind finds an entity by ID
func (st *Store) AttributeFind(entityID string, attributeKey string) *Attribute {
	attr := &Attribute{}

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
func (st *Store) AttributeGet(entityID string, attributeKey string) *Attribute {
	attr := &Attribute{}

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

// AttributesSet upserts and entity attribute
func (st *Store) AttributesSet(entityID string, attributes map[string]interface{}) bool {
	tx := st.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false
	}

	for k, v := range attributes {
		attr := st.AttributeFind(entityID, k)

		if attr == nil {
			attr = &Attribute{EntityID: entityID, AttributeKey: k}
			attr.SetInterface(v)

			dbResult := tx.Table(st.attributeTableName).Create(&attr)
			if dbResult.Error != nil {
				tx.Rollback()
				return false
			}

		}

		attr.SetInterface(v)
		dbResult := tx.Table(st.attributeTableName).Save(attr)
		if dbResult.Error != nil {
			return false
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return false
	}

	return true

}
