package entitystore

import (
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

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
