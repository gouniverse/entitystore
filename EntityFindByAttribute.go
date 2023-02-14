package entitystore

import (
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityFindByAttribute finds an entity by attribute
func (st *Store) EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) (*Entity, error) {
	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName)
	q = q.LeftJoin(goqu.I(st.entityTableName), goqu.On(goqu.Ex{st.attributeTableName + ".entity_id": goqu.I(st.entityTableName + ".id")}))
	q = q.Where(goqu.C("entity_type").Eq(entityType))
	q = q.Where(goqu.And(goqu.C("attribute_key").Eq(attributeKey), goqu.C("attribute_value").Eq(attributeValue)))
	q = q.Select("entity_id")

	sqlStr, _, _ := q.ToSQL()
	if st.GetDebug() {
		log.Println(sqlStr)
	}

	var entityID string
	err := st.db.QueryRow(sqlStr).Scan(&entityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		if st.GetDebug() {
			log.Println(err)
		}
		return nil, err
	}

	return st.EntityFindByID(entityID)
}
