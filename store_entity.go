package entitystore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/uid"
)

// EntityAttributeList list all attributes of an entity
func (st *Store) EntityAttributeList(entityID string) ([]Attribute, error) {
	var attrs []Attribute

	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName)
	q = q.Order(goqu.I("attribute_key").Asc())
	q = q.Where(goqu.C("entity_id").Eq(entityID))
	q = q.Select(Attribute{})

	sqlStr, _, _ := q.ToSQL()

	// DEBUG: log.Println(sqlStr)

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return attrs, err
	}

	for rows.Next() {
		var attr Attribute
		err := rows.Scan(&attr.AttributeKey, &attr.AttributeValue, &attr.CreatedAt, &attr.EntityID, &attr.ID, &attr.UpdatedAt)
		if err != nil {
			return attrs, err
		}
		// settingList = append(settingList, value)
		attrs = append(attrs, attr)
	}

	return attrs, nil
}

// EntityCount counts the entities of a specified type
// EntityCount counts entities
func (st *Store) EntityCount(entityType string) (int64, error) {
	var count int64

	q := goqu.Dialect(st.dbDriverName).From(st.entityTableName)
	q = q.Where(goqu.C("entity_type").Eq(entityType))
	q = q.Select(goqu.COUNT("*").As("count"))

	sqlStr, _, err := q.ToSQL()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}

		return 0, err
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	errScan := st.db.QueryRow(sqlStr).Scan(&count)

	if errScan != nil {
		if st.GetDebug() {
			log.Println(errScan)
		}

		if errScan == sql.ErrNoRows {
			return 0, errScan
		}
		return 0, errScan
	}

	return count, nil
}

// EntityCreate creates a new entity
func (st *Store) EntityCreate(entityType string) (*Entity, error) {
	return st.entityCreateWithTransactionOrDB(st.db, entityType)
}

func (st *Store) entityCreateWithTransactionOrDB(db txOrDB, entityType string) (*Entity, error) {
	entity := st.NewEntity(NewEntityOptions{
		ID:        uid.HumanUid(),
		Type:      entityType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTableName)
	q = q.Rows(entity.ToMap())
	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := db.Exec(sqlStr)

	if err != nil {
		return entity, err
	}

	return entity, nil
}

// EntityCreateWithAttributes func
func (st *Store) EntityCreateWithAttributes(entityType string, attributes map[string]string) (*Entity, error) {
	// Note the use of tx as the database handle once you are within a transaction
	tx, err := st.db.Begin()

	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	entity, err := st.entityCreateWithTransactionOrDB(tx, entityType)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for k, v := range attributes {
		_, err := st.attributeCreateWithTransactionOrDB(tx, entity.ID(), k, v)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return entity, nil
}

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

// EntityFindByHandle finds an entity by handle
func (st *Store) EntityFindByHandle(entityType string, entityHandle string) (*Entity, error) {
	if entityType == "" {
		return nil, errors.New("entity type cannot be empty")
	}

	if entityHandle == "" {
		return nil, errors.New("entity handle cannot be empty")
	}

	list, err := st.EntityList(EntityListQuery{
		EntityType:   entityType,
		EntityHandle: entityHandle,
		Limit:        1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

// EntityFindByID finds an entity by ID
func (st *Store) EntityFindByID(entityID string) (*Entity, error) {
	if entityID == "" {
		return nil, errors.New("entity ID cannot be empty")
	}

	list, err := st.EntityList(EntityListQuery{
		ID:    entityID,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

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

type EntityListQuery struct {
	ID           string
	IDs          []string
	EntityType   string
	EntityHandle string
	Limit        uint64
	Offset       uint64
	Search       string
	SortBy       string
	SortOrder    string // asc / dec
	CountOnly    bool
}

// EntityList lists entities
func (st *Store) EntityList(options EntityListQuery) (entityList []Entity, err error) {
	q := goqu.Dialect(st.dbDriverName).From(st.entityTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C("id").In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C("id").Eq(options.ID))
	}

	sortByColumn := "id"
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	if sortOrder == "asc" {
		q = q.Order(goqu.I(sortByColumn).Asc())
	} else {
		q = q.Order(goqu.I(sortByColumn).Desc())
	}

	if options.EntityType != "" {
		q = q.Where(goqu.C("entity_type").Eq(options.EntityType))
	}

	if options.EntityHandle != "" {
		q = q.Where(goqu.C("entity_handle").Eq(options.EntityType))
	}

	q = q.Offset(uint(options.Offset))

	if options.Limit != 0 {
		q = q.Limit(uint(options.Limit))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	q = q.Select()

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return entityList, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	entityMaps := []map[string]any{}
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
		entity := st.FromMap(entityMaps[i])
		entityList = append(entityList, *entity)
	}

	return entityList, nil
}

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

	return st.EntityList(EntityListQuery{
		EntityType: entityType,
		IDs:        entityIDs,
		SortBy:     "id",
	})
}

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
			ID:             attr.ID,
			EntityID:       attr.EntityID,
			AttributeKey:   attr.AttributeKey,
			AttributeValue: attr.AttributeValue,
			CreatedAt:      attr.CreatedAt,
			UpdatedAt:      attr.UpdatedAt,
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
