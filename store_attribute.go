package entitystore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/uid"
)

type AttributeListQuery struct {
	ID           string
	IDs          []string
	EntityID     string
	AttributeKey string
	Limit        uint64
	Offset       uint64
	SortBy       string
	SortOrder    string // asc / dec
	CountOnly    bool
}

// EntityList lists entities
func (st *Store) AttributeList(options AttributeListQuery) (attributeList []Attribute, err error) {
	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C("id").In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C("id").Eq(options.ID))
	}

	sortByColumn := "id"
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	if sortOrder == "asc" {
		q = q.Order(goqu.I(sortByColumn).Asc())
	} else {
		q = q.Order(goqu.I(sortByColumn).Desc())
	}

	if options.EntityID != "" {
		q = q.Where(goqu.C("entity_id").Eq(options.EntityID))
	}

	if options.AttributeKey != "" {
		q = q.Where(goqu.C("attribute_key").Eq(options.AttributeKey))
	}

	q = q.Offset(uint(options.Offset))

	if options.Limit != 0 {
		q = q.Limit(uint(options.Limit))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	q = q.Select()

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return attributeList, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	attributeMaps := []map[string]any{}
	errScan := sqlscan.Select(context.Background(), st.db, &attributeMaps, sqlStr)
	if errScan != nil {
		if errScan == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return nil, errScan
		}

		if sqlscan.NotFound(errScan) {
			return nil, nil
		}

		return nil, err
	}

	for i := 0; i < len(attributeMaps); i++ {
		attribute := st.AttributeFromMap(attributeMaps[i])
		attributeList = append(attributeList, *attribute)
	}

	return attributeList, nil
}

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	var newAttribute = st.NewAttribute(NewAttributeOptions{
		ID:             uid.HumanUid(),
		EntityID:       entityID,
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	return st.AttributeInsert(*newAttribute)
}

func (st *Store) attributeCreateWithTransactionOrDB(db txOrDB, entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	var newAttribute = st.NewAttribute(NewAttributeOptions{
		ID:             uid.HumanUid(),
		EntityID:       entityID,
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

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
	if entityID == "" {
		return nil, errors.New("entity id cannot be empty")
	}

	if attributeKey == "" {
		return nil, errors.New("attribute key cannot be empty")
	}

	list, err := st.AttributeList(AttributeListQuery{
		EntityID:     entityID,
		AttributeKey: attributeKey,
		Limit:        1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
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

// AttributesSet upserts an entity attribute
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
			attr = st.NewAttribute(NewAttributeOptions{ID: uid.HumanUid(), EntityID: entityID, AttributeKey: k, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			attr.SetString(v)

			q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName)
			q = q.Rows(attr.ToMap())
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
		attr.SetUpdatedAt(time.Now())

		q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
		q = q.Where(goqu.C("id").Eq(attr.ID))
		q = q.Set(attr.ToMap())

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
	if attr.AttributeKey() == "" {
		return nil, errors.New("attribute key is required field")
	}
	if attr.ID() == "" {
		attr.SetID(uid.HumanUid())
	}
	if attr.CreatedAt().IsZero() {
		attr.SetCreatedAt(time.Now())
	}
	if attr.UpdatedAt().IsZero() {
		attr.SetUpdatedAt(time.Now())
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName)
	q = q.Rows(attr.ToMap())
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
	attr.SetUpdatedAt(time.Now())

	q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
	q = q.Where(goqu.C("id").Eq(attr.ID))
	q = q.Set(attr.ToMap())

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
