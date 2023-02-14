package entitystore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// EntityUpdate updates an entity
func (st *Store) EntityUpdate(ent Entity) (bool, error) {
	ent.SetUpdatedAt(time.Now())

	q := goqu.Dialect(st.dbDriverName).Update(st.GetEntityTableName())
	q = q.Where(goqu.C("id").Eq(ent.ID)).Set(ent.ToMap())

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
