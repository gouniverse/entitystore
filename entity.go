package entitystore

import (
	//"encoding/json"
	"database/sql"

	//"fmt"
	"log"
	//"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
)

const (
	// EntityStatusActive entity "active" status
	EntityStatusActive = "active"
	// EntityStatusInactive entity "inactive" status
	EntityStatusInactive = "inactive"
)

// Entity type
type Entity struct {
	ID        string    `db:"id"`
	Status    string    `db:"status"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	st        *Store    //`db:-`
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

// GetAny the value of the attribute as interface{} or the default value if it does not exist
func (e *Entity) GetAny(attributeKey string, defaultValue interface{}) (interface{}, error) {
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

	return attr.GetInterface(), nil
}

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

// GetAttribute the name of the User table
func (e *Entity) GetAttribute(attributeKey string) (*Attribute, error) {
	attr := &Attribute{}

	sqlStr, _, _ := goqu.From(e.st.attributeTableName).Where(goqu.C("attribute_key").Eq(attributeKey), goqu.C("deleted_at").IsNull()).Select("attribute_key", "attribute_value", "created_at", "deleted_at", "entity_id", "id", "updated_at").ToSQL()

	if e.st.GetDebug() {
		log.Println(sqlStr)
	}

	var createdAt string
	var updatedAt string
	var deletedAt *string
	err := e.st.db.QueryRow(sqlStr).Scan(&attr.AttributeKey, &attr.AttributeValue, &createdAt, &deletedAt, &attr.EntityID, &attr.ID, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if e.st.GetDebug() {
			log.Println(err)
		}
		return nil, err
	}

	layout := "Mon Jan 02 2006 15:04:05 GMT-0700"
	createdAtTime, err := time.Parse(layout, createdAt)
	if err == nil {
		attr.CreatedAt = createdAtTime
	}
	updatedAtTime, err := time.Parse(layout, updatedAt)
	if err == nil {
		attr.UpdatedAt = updatedAtTime
	}
	// if deletedAt != nil {
	// 	deletedAtTime, err := time.Parse(layout, *deletedAt)
	// 	if err == nil {
	// 		attr.DeletedAt = &deletedAtTime
	// 	}
	// }

	return attr, nil
}

// SetAllAny upserts the attributes
func (e *Entity) SetAllAny(attributes map[string]interface{}) (bool, error) {
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
