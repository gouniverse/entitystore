package entitystore

import (
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityDelete deletes an entity and all attributes
func (st *Store) EntityDelete(entityID string) (bool, error) {
	if entityID == "" {
		if st.GetDebug() {
			log.Println("in EntityDelete entity ID cannot be empty")
		}
		return false, errors.New("in EntityDelete entity ID cannot be empty")
	}

	// Note the use of tx as the database handle once you are within a transaction
	tx, err := st.db.Begin()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return false, err
	}

	defer func() {
		if r := recover(); r != nil {
			txErr := tx.Rollback()
			if txErr != nil && st.GetDebug() {
				log.Println(txErr)
			}
		}
	}()

	sqlStr1, _, _ := goqu.Dialect(st.dbDriverName).From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID)).Delete().ToSQL()

	if _, err := tx.Exec(sqlStr1); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		txErr := tx.Rollback()
		if txErr != nil && st.GetDebug() {
			log.Println(txErr)
		}
		return false, err
	}

	sqlStr2, _, _ := goqu.Dialect(st.dbDriverName).From(st.entityTableName).Where(goqu.C("id").Eq(entityID)).Delete().ToSQL()

	if _, err := tx.Exec(sqlStr2); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		txErr := tx.Rollback()
		if txErr != nil && st.GetDebug() {
			log.Println(txErr)
		}
		return false, err
	}

	err = tx.Commit()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}

		return false, err
	}

	return true, nil
}
