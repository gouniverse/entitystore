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

func (st *Store) attributeCreateWithTransactionOrDB(db txOrDB, entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	var newAttribute = &Attribute{
		ID:             uid.HumanUid(),
		EntityID:       entityID,
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return st.attributeInsertWithTransactionOrDB(db, *newAttribute)
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

	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName)
	q = q.Where(goqu.C("entity_id").Eq(entityID), goqu.C("attribute_key").Eq(attributeKey))
	q = q.Select("attribute_key", "attribute_value", "created_at", "entity_id", "id", "updated_at")
	sqlStr, _, _ := q.ToSQL()

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

// AttributeSetFloat creates a new attribute or updates existing
func (st *Store) AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) (bool, error) {
	attributeValueAsString := strconv.FormatFloat(attributeValue, 'f', 30, 64)
	return st.AttributeSetString(entityID, attributeKey, attributeValueAsString)
}

// AttributeSetInt creates a new attribute or updates existing
func (st *Store) AttributeSetInt(entityID string, attributeKey string, attributeValue int64) (bool, error) {
	attributeValueAsString := strconv.FormatInt(attributeValue, 10)
	return st.AttributeSetString(entityID, attributeKey, attributeValueAsString)
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

			q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName)
			q = q.Rows(attr)
			sqlStr, _, err := q.ToSQL()

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

		q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
		q = q.Where(goqu.C("id").Eq(attr.ID))
		q = q.Set(attr)

		sqlStr, _, err := q.ToSQL()

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
	return st.attributeInsertWithTransactionOrDB(st.db, attr)
}

func (st *Store) attributeInsertWithTransactionOrDB(db txOrDB, attr Attribute) (*Attribute, error) {
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

	q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName)
	q = q.Rows(attr)
	sqlStr, _, _ := q.ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := db.Exec(sqlStr)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return nil, err
	}

	return &attr, nil
}

// AttributeUpdate updates an attribute
func (st *Store) AttributeUpdate(attr Attribute) (bool, error) {
	attr.UpdatedAt = time.Now()

	q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
	q = q.Where(goqu.C("id").Eq(attr.ID))
	q = q.Set(attr)

	sqlStr, _, _ := q.ToSQL()

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
