package entitystore

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// EntityAttributeList list all attributes of an entity
func (st *Store) EntityAttributeList(entityID string) ([]Attribute, error) {
	var attrs []Attribute

	sqlStr, _, _ := goqu.From(st.attributeTableName).Order(goqu.I("attribute_key").Asc()).Where(goqu.C("entity_id").Eq(entityID)).Select(Attribute{}).ToSQL()

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

// EntityCount counts entities
func (st *Store) EntityCount(entityType string) uint64 {
	var count int64
	count, _ = goqu.From(st.entityTableName).Where(goqu.C("entity_type").Eq(entityType)).Count()
	return uint64(count)
}

// EntityCreate creates a new entity
func (st *Store) EntityCreate(entityType string) (*Entity, error) {
	return st.entityCreateWithTransactionOrDB(st.db, entityType)
}

func (st *Store) entityCreateWithTransactionOrDB(db txOrDB, entityType string) (*Entity, error) {
	entity := &Entity{ID: uid.HumanUid(), Type: entityType, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now(), st: st}

	sqlStr, _, _ := goqu.Insert(st.entityTableName).Rows(entity).ToSQL()

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
		_, err := st.attributeCreateWithTransactionOrDB(tx, entity.ID, k, v)

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

	sqlStr1, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID)).Delete().ToSQL()

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

	sqlStr2, _, _ := goqu.From(st.entityTableName).Where(goqu.C("id").Eq(entityID)).Delete().ToSQL()

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
func (st *Store) EntityFindByHandle(entityType string, entityHandle string) *Entity {
	if entityType == "" {
		return nil
	}

	if entityHandle == "" {
		return nil
	}

	ent := &Entity{}

	sqlStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("entity_type").Eq(entityType), goqu.C("entity_handle").Eq(entityHandle)).Select("created_at", "ïd", "type", "updated_at").ToSQL()

	// DEBUG: log.Println(sqlStr)

	err := st.db.QueryRow(sqlStr).Scan(&ent.CreatedAt, &ent.ID, &ent.Type, &ent.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	ent.st = st // Add store reference

	return ent
}

// EntityFindByID finds an entity by ID
func (st *Store) EntityFindByID(entityID string) (*Entity, error) {
	if entityID == "" {
		return nil, errors.New("entity ID cannot be emopty")
	}

	ent := &Entity{}

	sqlStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("id").Eq(entityID)).Select("id", "entity_type", "entity_status", "entity_handle", "created_at", "updated_at").ToSQL()

	// DEBUG: log.Println(sqlStr)

	err := st.db.QueryRow(sqlStr).Scan(&ent.ID, &ent.Type, &ent.Status, &ent.Handle, &ent.CreatedAt, &ent.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if st.GetDebug() {
			log.Println(err)
		}
		return nil, err
	}

	ent.st = st // Add store reference

	return ent, nil
}

// EntityFindByAttribute finds an entity by attribute
func (st *Store) EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) (*Entity, error) {
	q := goqu.From(st.attributeTableName)
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

// EntityList lists entities
func (st *Store) EntityList(entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) ([]Entity, error) {
	entityList := []Entity{}

	sqlStr, _, _ := goqu.From(st.entityTableName).Order(goqu.I("id").Asc()).Where(goqu.C("entity_type").Eq(entityType)).Offset(uint(offset)).Limit(uint(perPage)).Select("id", "entity_type", "entity_status", "entity_handle", "created_at", "updated_at").ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return entityList, err
	}

	for rows.Next() {
		var ent Entity
		err := rows.Scan(&ent.ID, &ent.Type, &ent.Status, &ent.Handle, &ent.CreatedAt, &ent.UpdatedAt)
		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}
			return entityList, err
		}
		ent.st = st // Important
		entityList = append(entityList, ent)
	}

	return entityList, nil
}

// EntityListByAttribute finds an entity by attribute
func (st *Store) EntityListByAttribute(entityType string, attributeKey string, attributeValue string) ([]Entity, error) {
	//entityAttributes := []EntityAttribute{}
	var entityIDs []string

	q := goqu.From(st.attributeTableName)
	q = q.LeftJoin(goqu.I(st.entityTableName), goqu.On(goqu.Ex{st.attributeTableName + ".entity_id": goqu.I(st.entityTableName + ".id")}))
	q = q.Where(goqu.C("entity_type").Eq(entityType))
	q = q.Where(goqu.And(goqu.C("attribute_key").Eq(attributeKey), goqu.C("attribute_value").Eq(attributeValue)))
	q = q.Select("entity_id")

	sqlStr, _, err := q.ToSQL()
	if st.GetDebug() {
		log.Println(sqlStr)
	}

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

	entityList := []Entity{}

	sqlStr, _, _ = goqu.From(st.attributeTableName).Order(goqu.I("id").Asc()).Where(goqu.C("id").In(entityIDs)).Select("id", "entity_type", "entity_status", "entity_handle", "created_at", "updated_at").ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	rows, err = st.db.Query(sqlStr)

	if err != nil {
		return entityList, err
	}

	for rows.Next() {
		var ent Entity
		err := rows.Scan(&ent.ID, &ent.Type, &ent.Status, &ent.Handle, &ent.CreatedAt, &ent.UpdatedAt)
		if err != nil {
			return entityList, err
		}
		ent.st = st // Important
		entityList = append(entityList, ent)
	}

	return entityList, nil
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
		ID:        ent.ID,
		Status:    ent.Status,
		Type:      ent.Type,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
		DeletedAt: time.Now(),
	}

	sqlStr, _, _ := goqu.Insert(st.entityTrashTableName).Rows(entTrash).ToSQL()

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

		sqlStrAttr, _, _ := goqu.Insert(st.attributeTrashTableName).Rows(attrTrash).ToSQL()

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

	sqlStr1, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID)).Delete().ToSQL()

	if _, err := tx.Exec(sqlStr1); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		tx.Rollback()
		return false, err
	}

	sqlStr2, _, _ := goqu.From(st.entityTableName).Where(goqu.C("id").Eq(entityID)).Delete().ToSQL()

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
	ent.UpdatedAt = time.Now()
	sqlStr, _, _ := goqu.Update(st.GetEntityTableName()).Where(goqu.C("id").Eq(ent.ID)).Set(ent).ToSQL()

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
