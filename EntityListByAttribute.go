package entitystore

import (
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityListByAttribute finds an entity by attribute
func (st *Store) EntityListByAttribute(entityType string, attributeKey string, attributeValue string) (entityList []Entity, err error) {
	var entityIDs []string

	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName).
		LeftJoin(goqu.I(st.entityTableName), goqu.On(goqu.Ex{st.attributeTableName + ".entity_id": goqu.I(st.entityTableName + ".id")})).
		Where(goqu.C("entity_type").Eq(entityType)).
		Where(goqu.And(goqu.C("attribute_key").Eq(attributeKey), goqu.C("attribute_value").Eq(attributeValue))).
		Select("entity_id")

	sqlStr, _, err := q.ToSQL()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return nil, err
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return []Entity{}, err
	}

	for rows.Next() {
		var entityID string
		err := rows.Scan(&entityID)
		if err != nil {
			return []Entity{}, err
		}
		entityIDs = append(entityIDs, entityID)
	}

	if len(entityIDs) < 1 {
		return entityList, nil
	}

	return st.EntityList(EntityQueryOptions{
		EntityType: entityType,
		IDs:        entityIDs,
		SortBy:     "id",
	})
}
