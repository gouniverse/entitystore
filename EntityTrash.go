package entitystore

import (
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// EntityTrash moves an entity and all attributes to the trash bin
func (st *Store) EntityTrash(entityID string) (bool, error) {
	if entityID == "" {
		return false, errors.New("entity ID cannot be empty")
	}

	// Note the use of tx as the database handle once you are within a transaction
	tx, err := st.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return false, err
	}

	ent, err := st.EntityFindByID(entityID)

	if err != nil {
		tx.Rollback()
		return false, err
	}

	if ent == nil {
		tx.Rollback()
		return false, err
	}

	entTrash := EntityTrash{
		ID:        ent.ID(),
		Type:      ent.Type(),
		CreatedAt: ent.CreatedAt(),
		UpdatedAt: ent.UpdatedAt(),
		DeletedAt: time.Now(),
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTrashTableName)
	q = q.Rows(entTrash)
	sqlStr, _, _ := q.ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	if _, err := tx.Exec(sqlStr); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		tx.Rollback()
		return false, err
	}

	attrs, err := st.EntityAttributeList(entityID)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		tx.Rollback()
		return false, err
	}

	for _, attr := range attrs {
		attrTrash := AttributeTrash{
			ID:             attr.ID(),
			EntityID:       attr.EntityID(),
			AttributeKey:   attr.AttributeKey(),
			AttributeValue: attr.AttributeValue(),
			CreatedAt:      attr.CreatedAt(),
			UpdatedAt:      attr.UpdatedAt(),
			DeletedAt:      time.Now(),
		}

		q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTrashTableName)
		q = q.Rows(attrTrash)
		sqlStrAttr, _, _ := q.ToSQL()

		if st.GetDebug() {
			log.Println(sqlStrAttr)
		}

		if _, err := tx.Exec(sqlStrAttr); err != nil {
			if st.GetDebug() {
				log.Println(err)
			}
			tx.Rollback()
			return false, err
		}
	}

	q1 := goqu.Dialect(st.dbDriverName).From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID)).Delete()
	sqlStr1, _, _ := q1.ToSQL()

	if _, err := tx.Exec(sqlStr1); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		tx.Rollback()
		return false, err
	}

	q2 := goqu.Dialect(st.dbDriverName).From(st.entityTableName).Where(goqu.C("id").Eq(entityID)).Delete()
	sqlStr2, _, _ := q2.ToSQL()

	if _, err := tx.Exec(sqlStr2); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		tx.Rollback()
		return false, err
	}

	err = tx.Commit()

	if err == nil {
		return true, nil
	}

	if st.GetDebug() {
		log.Println(err)
	}

	return false, err
}
