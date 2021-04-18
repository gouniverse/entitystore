package entitystore

import (
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) *Attribute {
	attr := &Attribute{EntityID: entityID, AttributeKey: attributeKey, AttributeValue: attributeValue}

	dbResult := st.db.Table(st.attributeTableName).Create(&attr)

	if dbResult.Error != nil {
		return nil
	}

	return attr
}

// AttributeCreateInterface creates a new attribute
func (st *Store) AttributeCreateInterface(entityID string, attributeKey string, attributeValue interface{}) *Attribute {
	attr := &Attribute{EntityID: entityID, AttributeKey: attributeKey}
	attr.SetInterface(attributeValue)

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

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return attr
}

// AttributeSetFloat creates a new entity
func (st *Store) AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreateInterface(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetFloat(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// AttributeSetInt creates a new entity
func (st *Store) AttributeSetInt(entityID string, attributeKey string, attributeValue int64) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreateInterface(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetInt(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// AttributeSetInterface creates a new entity
func (st *Store) AttributeSetInterface(entityID string, attributeKey string, attributeValue interface{}) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreateInterface(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetInterface(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// AttributeSetString creates a new entity
func (st *Store) AttributeSetString(entityID string, attributeKey string, attributeValue string) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreate(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetString(attributeValue)

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
