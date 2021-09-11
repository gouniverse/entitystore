package entitystore

import (
	"database/sql"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	var newAttribute = &Attribute{
		ID:             uid.MicroUid(),
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	sqlStr, _, _ := goqu.Insert(st.attributeTableName).Rows(newAttribute).ToSQL()

	log.Println(sqlStr)

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		log.Println(err)
		return newAttribute, err
	}

	return newAttribute, nil
}

// AttributeCreateInterface creates a new attribute
func (st *Store) AttributeCreateInterface(entityID string, attributeKey string, attributeValue interface{}) (*Attribute, error) {
	attr := &Attribute{EntityID: entityID, AttributeKey: attributeKey}
	attr.SetInterface(attributeValue)

	sqlStr, _, _ := goqu.Insert(st.attributeTableName).Rows(attr).ToSQL()

	log.Println(sqlStr)

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		log.Println(err)
		return attr, err
	}

	return attr, nil
}

// AttributeFind finds an entity by ID
func (st *Store) AttributeFind(entityID string, attributeKey string) *Attribute {
	attr := &Attribute{}

	sqlStr, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID), goqu.C("deleted_at").IsNull()).Select(Attribute{}).ToSQL()

	log.Println(sqlStr)

	err := st.db.QueryRow(sqlStr).Scan(&attr.AttributeKey, &attr.AttributeValue, &attr.CreatedAt, &attr.DeletedAt, &attr.ID, &attr.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	return attr
}

// AttributeSetFloat creates a new entity
func (st *Store) AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) (bool, error) {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr, err := st.AttributeCreateInterface(entityID, attributeKey, attributeValue)
		if err != nil {
			return false, err
		}
		if attr != nil {
			return true, nil
		}
		return false, err
	}

	attr.SetFloat(attributeValue)

	attr.UpdatedAt = time.Now()
	sqlStr, _, _ := goqu.Update(st.attributeTableName).Set(attr).ToSQL()

	// log.Println(sqlStr)

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
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
