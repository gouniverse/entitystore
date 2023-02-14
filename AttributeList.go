package entitystore

import (
	"context"
	"database/sql"
	"log"

	"github.com/georgysavva/scany/sqlscan"
)

// AttributeList lists attributes
func (st *Store) AttributeList(options AttributeQueryOptions) (attributeList []Attribute, err error) {
	q := st.AttributeQuery(options)

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return attributeList, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	attributeMaps := []map[string]string{}
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
		attribute := st.NewAttributeFromMap(attributeMaps[i])
		attributeList = append(attributeList, *attribute)
	}

	return attributeList, nil
}
