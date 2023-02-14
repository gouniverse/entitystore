package entitystore

import (
	"context"
	"database/sql"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

// EntityList lists entities
func (st *Store) EntityList(options EntityQueryOptions) (entityList []Entity, err error) {
	q := st.EntityQuery(options)

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return entityList, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	entityMaps := []map[string]string{}
	errScan := sqlscan.Select(context.Background(), st.db, &entityMaps, sqlStr)
	if errScan != nil {
		if errScan == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return nil, errScan
		}

		if sqlscan.NotFound(errScan) {
			return nil, nil
		}

		log.Println("FunnelStore. VisitorList. Error: ", err)
		return nil, err
	}

	for i := 0; i < len(entityMaps); i++ {
		entity := st.NewEntityFromMap(entityMaps[i])
		entityList = append(entityList, *entity)
	}

	return entityList, nil
}
