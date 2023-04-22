package entitystore

import (
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityInsert inserts a new entity
func (st *Store) EntityInsert(entity Entity) error {
	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTableName)
	q = q.Rows(entity.ToMap())

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}
