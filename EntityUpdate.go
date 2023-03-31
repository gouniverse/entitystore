package entitystore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// EntityUpdate updates an entity
func (st *Store) EntityUpdate(ent Entity) error {
	ent.SetUpdatedAt(time.Now())

	q := goqu.Dialect(st.dbDriverName).
		Update(st.GetEntityTableName()).
		Where(goqu.C("id").Eq(ent.ID())).
		Set(ent.ToMap())

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(sqlStr)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}

		return err
	}

	return nil
}
