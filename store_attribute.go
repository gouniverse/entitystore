package entitystore

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	var newAttribute = &Attribute{
		ID:             uid.HumanUid(),
		EntityID:       entityID,
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return st.AttributeInsert(*newAttribute)
}

// AttributeCreateInterface creates a new attribute
// func (st *Store) AttributeCreateInterface(entityID string, attributeKey string, attributeValue interface{}) (*Attribute, error) {
// 	attr := &Attribute{ID: uid.HumanUid(), EntityID: entityID, AttributeKey: attributeKey, CreatedAt: time.Now(), UpdatedAt: time.Now()}
// 	attr.SetInterface(attributeValue)

// 	return st.AttributeInsert(*attr)
// }

// AttributeFind finds an entity by ID
func (st *Store) AttributeFind(entityID string, attributeKey string) (*Attribute, error) {
	attr := &Attribute{}

	sqlStr, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID)).Select("attribute_key", "attribute_value", "created_at", "entity_id", "id", "updated_at").ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	var createdAt string
	var updatedAt string

	err := st.db.QueryRow(sqlStr).Scan(&attr.AttributeKey, &attr.AttributeValue, &createdAt, &attr.EntityID, &attr.ID, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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

	return attr, nil
}

// AttributeSetFloat creates a new entity
func (st *Store) AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) (bool, error) {
	attr, err := st.AttributeFind(entityID, attributeKey)

	if err != nil {
		return false, err
	}

	if attr == nil {
		attributeValueAsString := strconv.FormatFloat(attributeValue, 'f', 30, 64)
		isOk, err := st.AttributeSetString(entityID, attributeKey, attributeValueAsString)
		if err != nil {
			return false, err
		}
		if isOk {
			return true, nil
		}
		return false, err
	}

	attr.SetFloat(attributeValue)

	return st.AttributeUpdate(*attr)
}

// AttributeSetInt creates a new entity
func (st *Store) AttributeSetInt(entityID string, attributeKey string, attributeValue int64) (bool, error) {
	attr, err := st.AttributeFind(entityID, attributeKey)

	if err != nil {
		return false, err
	}

	if attr == nil {
		attributeValueAsString := strconv.FormatInt(attributeValue, 10)
		isOk, err := st.AttributeSetString(entityID, attributeKey, attributeValueAsString)
		if err != nil {
			return false, err
		}
		if isOk {
			return true, nil
		}
		return false, err
	}

	attr.SetInt(attributeValue)

	return st.AttributeUpdate(*attr)
}

// AttributeSetString creates a new entity
func (st *Store) AttributeSetString(entityID string, attributeKey string, attributeValue string) (bool, error) {
	attr, err := st.AttributeFind(entityID, attributeKey)

	if err != nil {
		return false, err
	}

	if attr == nil {
		attr, err := st.AttributeCreate(entityID, attributeKey, attributeValue)
		if err != nil {
			return false, err
		}
		if attr != nil {
			return true, nil
		}
		return false, err
	}

	attr.SetString(attributeValue)

	return st.AttributeUpdate(*attr)
}

// AttributesSet upserts and entity attribute
func (st *Store) AttributesSet(entityID string, attributes map[string]string) (bool, error) {
	tx, err := st.db.Begin()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return false, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback()
			if st.GetDebug() {
				log.Println(err)
			}
		}
	}()

	for k, v := range attributes {
		attr, err := st.AttributeFind(entityID, k)

		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}

			err = tx.Rollback()

			if st.GetDebug() {
				log.Println(err)
			}

			return false, err
		}

		if attr == nil {
			attr = &Attribute{ID: uid.HumanUid(), EntityID: entityID, AttributeKey: k, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			attr.SetString(v)

			sqlStr, _, err := goqu.Insert(st.attributeTableName).Rows(attr).ToSQL()

			if err != nil {
				if st.GetDebug() {
					log.Println(err)
				}

				err = tx.Rollback()

				if st.GetDebug() {
					log.Println(err)
				}

				return false, err
			}

			if st.GetDebug() {
				log.Println(sqlStr)
			}

			_, err = tx.Exec(sqlStr)

			if err != nil {
				log.Println(err)
				err = tx.Rollback()

				if st.GetDebug() {
					log.Println(err)
				}

				return false, err
			}

		}

		attr.SetString(v)
		attr.UpdatedAt = time.Now()
		sqlStr, _, err := goqu.Update(st.attributeTableName).Where(goqu.C("id").Eq(attr.ID)).Set(attr).ToSQL()

		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}

			err = tx.Rollback()

			if st.GetDebug() {
				log.Println(err)
			}

			return false, err
		}

		if st.GetDebug() {
			log.Println(sqlStr)
		}

		_, err = tx.Exec(sqlStr)

		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}

			err = tx.Rollback()

			if st.GetDebug() {
				log.Println(err)
			}

			return false, err
		}
	}

	err = tx.Commit()

	if err != nil {
		err = tx.Rollback()

		if st.GetDebug() {
			log.Println(err)
		}

		return false, err
	}

	return true, nil

}

// AttributeCreate creates a new attribute
func (st *Store) AttributeInsert(attr Attribute) (*Attribute, error) {
	if attr.AttributeKey == "" {
		return nil, errors.New("attribute key is required field")
	}
	if attr.ID == "" {
		attr.ID = uid.HumanUid()
	}
	if attr.CreatedAt.IsZero() {
		attr.CreatedAt = time.Now()
	}
	if attr.UpdatedAt.IsZero() {
		attr.UpdatedAt = time.Now()
	}

	sqlStr, _, _ := goqu.Insert(st.attributeTableName).Rows(attr).ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return nil, err
	}

	return &attr, nil
}

// AttributeSetString creates a new entity
func (st *Store) AttributeUpdate(attr Attribute) (bool, error) {
	attr.UpdatedAt = time.Now()
	sqlStr, _, _ := goqu.Update(st.attributeTableName).Where(goqu.C("id").Eq(attr.ID)).Set(attr).ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}

		return false, err
	}

	return true, nil
}
