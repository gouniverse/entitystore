package entitystore

import (
	//"encoding/json"
	"database/sql"

	//"fmt"
	"log"
	//"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
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
	ID     string `db:"id"`
	Status string `db:"status"`
	Type   string `db:"type"`
	Handle string `db:"name"`
	//Name        string     `gorm:"type:varchar(255);column:name;DEFAULT NULL;"`
	//Description string     `gorm:"type:longtext;column:description;"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	//attributes []Attribute `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	st *Store //`db:-`
}

// BeforeCreate adds UID to model
func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	e.ID = uuid
	return nil
}

// Delete deletes the entity
// func (e *Entity) Delete() bool {
// 	return e.st.EntityDelete(e.ID)
// }

// Trash moves the entity to the trash bin
// func (e *Entity) Trash() bool {
// 	return e.st.EntityTrash(e.ID)
// }

// GetAny the value of the attribute as interface{} or the default value if it does not exist
func (e *Entity) GetAny(attributeKey string, defaultValue interface{}) interface{} {
	attr := e.GetAttribute(attributeKey)

	if attr == nil {
		return defaultValue
	}

	return attr.GetInterface()
}

// GetInt the value of the attribute as string or the default value if it does not exist
func (e *Entity) GetInt(attributeKey string, defaultValue int64) (int64, error) {
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

	sqlStr, _, _ := goqu.From(e.st.attributeTableName).Where(goqu.C("attribute_key").Eq(attributeKey), goqu.C("deleted_at").IsNull()).Select(Attribute{}).ToSQL()

	log.Println(sqlStr)

	err := e.st.db.QueryRow(sqlStr).Scan(&attr.CreatedAt, &attr.DeletedAt, &attr.ID, &attr.AttributeKey, &attr.AttributeValue, &attr.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	// result := e.st.db.Table(e.st.attributeTableName).First(&attr, "entity_id=? AND attribute_key=?", e.ID, attributeKey)

	// if result.Error != nil {

	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		return nil
	// 	}

	// 	log.Panic(result.Error)
	// }

	return attr
}

// SetAllAny upserts the attributes
func (e *Entity) SetAllAny(attributes map[string]interface{}) bool {
	return e.st.AttributesSet(e.ID, attributes)
}

// SetInterface sets an attribute with string value
func (e *Entity) SetInterface(attributeKey string, attributeValue interface{}) (bool, error) {
	return e.st.AttributeSetInterface(e.ID, attributeKey, attributeValue)
}

// SetFloat sets an attribute with float value
func (e *Entity) SetFloat(attributeKey string, attributeValue float64) (bool, error) {
	return e.st.AttributeSetFloat(e.ID, attributeKey, attributeValue)
}

// SetInt sets an attribute with int value
func (e *Entity) SetInt(attributeKey string, attributeValue int64) (bool, error) {
	return e.st.AttributeSetInt(e.ID, attributeKey, attributeValue)
}

// SetString sets an attribute with string value
func (e *Entity) SetString(attributeKey string, attributeValue string) (bool, error) {
	return e.st.AttributeSetString(e.ID, attributeKey, attributeValue)
}
