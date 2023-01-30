package entitystore

import (
	"log"
	"time"
)

// Entity type
type Entity struct {
	ID        string    `db:"id"`
	Type      string    `db:"entity_type"`
	Handle    string    `db:"entity_handle"`
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
	return e.st.AttributeFind(e.ID, attributeKey)
}

// GetAttributes all the attributes of the entity
func (e *Entity) GetAttributes() ([]Attribute, error) {
	return e.st.EntityAttributeList(e.ID)
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
	return e.st.AttributesSet(e.ID, attributes)
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
