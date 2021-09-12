package entitystore

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// EntityAttributeList list all atributes of an entity
func (st *Store) EntityAttributeList(entityID string) ([]Attribute, error) {
	var attrs []Attribute

	sqlStr, _, _ := goqu.From(st.attributeTableName).Order(goqu.I("attribute_key").Asc()).Where(goqu.C("entity_id").Eq(entityID)).Where(goqu.C("deleted_at").IsNull()).Select(Attribute{}).ToSQL()

	// DEBUG: log.Println(sqlStr)

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return attrs, err
	}

	for rows.Next() {
		var attr Attribute
		err := rows.Scan(&attr.AttributeKey, &attr.AttributeValue, &attr.CreatedAt, &attr.DeletedAt, &attr.EntityID, &attr.ID, &attr.UpdatedAt)
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
	count, _ = goqu.From(st.entityTableName).Where(goqu.C("type").Eq(entityType), goqu.C("deleted_at").IsNull()).Count()
	return uint64(count)
}

// EntityCreate creates a new entity
func (st *Store) EntityCreate(entityType string) (*Entity, error) {
	entity := &Entity{ID: uid.HumanUid(), Type: entityType, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now(), st: st}

	sqlStr, _, _ := goqu.Insert(st.entityTableName).Rows(entity).ToSQL()

	// DEBUG: log.Println(sqlStr)

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return entity, err
	}

	return entity, nil
}

// EntityCreateWithAttributes func
func (st *Store) EntityCreateWithAttributes(entityType string, attributes map[string]interface{}) *Entity {
	// Note the use of tx as the database handle once you are within a transaction
	tx, err := st.db.Begin()

	if err != nil {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	entity, err := st.EntityCreate(entityType)

	if err != nil {
		tx.Rollback()
		return nil
	}

	for k, v := range attributes {
		_, err := st.AttributeCreateInterface(entity.ID, k, v)

		if err != nil {
			tx.Rollback()
			return nil
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return nil
	}

	return entity
}

// EntityDelete deletes an entity and all attributes
func (st *Store) EntityDelete(entityID string) bool {
	if entityID == "" {
		return false
	}

	// Note the use of tx as the database handle once you are within a transaction
	tx, err := st.db.Begin()

	if err != nil {
		log.Println(err)
		return false
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	sqlStr1, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID), goqu.C("deleted_at").IsNull()).Delete().ToSQL()

	if _, err := tx.Exec(sqlStr1); err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	sqlStr2, _, _ := goqu.From(st.entityTableName).Where(goqu.C("id").Eq(entityID), goqu.C("deleted_at").IsNull()).Delete().ToSQL()

	if _, err := tx.Exec(sqlStr2); err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	err = tx.Commit()

	if err == nil {
		return true
	}

	log.Println(err)

	return false
}

// EntityDeleteSoft soft deletes an entity and all attributes
// func (st *Store) EntityDeleteSoft(entityID string) bool {
// 	if entityID == "" {
// 		return false
// 	}

// 	// Note the use of tx as the database handle once you are within a transaction
// 	tx := st.db.Begin()

// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	if err := tx.Error; err != nil {
// 		log.Println(err)
// 		return false
// 	}

// 	if err := tx.Where("entity_id=?", entityID).Table(st.attributeTableName).Update("deleted_at", time.Now()).Error; err != nil {
// 		tx.Rollback()
// 		log.Println(err)
// 		return false
// 	}

// 	if err := tx.Where("id=?", entityID).Table(st.entityTableName).Update("deleted_at", time.Now()).Error; err != nil {
// 		tx.Rollback()
// 		log.Println(err)
// 		return false
// 	}

// 	err := tx.Commit().Error

// 	if err == nil {
// 		return true
// 	}

// 	log.Println(err)

// 	return false
// }

// EntityFindByHandle finds an entity by handle
func (st *Store) EntityFindByHandle(entityType string, entityHandle string) *Entity {
	if entityType == "" {
		return nil
	}

	if entityHandle == "" {
		return nil
	}

	ent := &Entity{}

	sqlStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("entity_type").Eq(entityType), goqu.C("entity_handle").Eq(entityHandle), goqu.C("deleted_at").IsNull()).Select("created_at", "ïd", "type", "updated_at").ToSQL()

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

	sqlStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("id").Eq(entityID), goqu.C("deleted_at").IsNull()).Select("created_at", "id", "status", "type", "updated_at").ToSQL()

	// DEBUG: log.Println(sqlStr)

	err := st.db.QueryRow(sqlStr).Scan(&ent.CreatedAt, &ent.ID, &ent.Status, &ent.Type, &ent.UpdatedAt)
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
	subqueryStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("type").Eq(entityType), goqu.C("deleted_at").IsNull()).Select("id").ToSQL()
	sqlStr, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").In(subqueryStr), goqu.C("attribute_key").Eq(attributeKey), goqu.C("atribute_value").Eq(attributeValue), goqu.C("deleted_at").IsNull()).Select("entity_id").ToSQL()

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

	sqlStr, _, _ := goqu.From(st.attributeTableName).Order(goqu.I("id").Asc()).Where(goqu.C("type").Eq(entityType)).Where(goqu.C("deleted_at").IsNull()).Offset(uint(offset)).Limit(uint(perPage)).Select("created_at", "id", "type", "updated_at").ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return entityList, err
	}

	for rows.Next() {
		var ent Entity
		err := rows.Scan(&ent.CreatedAt, &ent.ID, &ent.Type, &ent.UpdatedAt)
		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}
			return entityList, err
		}
		entityList = append(entityList, ent)
	}

	return entityList, nil
}

// EntityListByAttribute finds an entity by attribute
func (st *Store) EntityListByAttribute(entityType string, attributeKey string, attributeValue string) ([]Entity, error) {
	//entityAttributes := []EntityAttribute{}
	var entityIDs []string

	subqueryStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("type").Eq(entityType), goqu.C("deleted_at").IsNull()).Select("id").ToSQL()
	sqlStr, _, err := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").In(subqueryStr), goqu.C("attribute_key").Eq(attributeKey), goqu.C("atribute_value").Eq(attributeValue), goqu.C("deleted_at").IsNull()).Select("entity_id").ToSQL()

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

	sqlStr, _, _ = goqu.From(st.attributeTableName).Order(goqu.I("id").Asc()).Where(goqu.C("id").In(entityIDs)).Where(goqu.C("deleted_at").IsNull()).Select("created_at", "id", "type", "updated_at").ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	rows, err = st.db.Query(sqlStr)

	if err != nil {
		return entityList, err
	}

	for rows.Next() {
		var ent Entity
		err := rows.Scan(&ent.CreatedAt, &ent.ID, &ent.Type, &ent.UpdatedAt)
		if err != nil {
			return entityList, err
		}
		ent.st = st
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

	sqlStr1, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID), goqu.C("deleted_at").IsNull()).Delete().ToSQL()

	if _, err := tx.Exec(sqlStr1); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		tx.Rollback()
		return false, err
	}

	sqlStr2, _, _ := goqu.From(st.entityTableName).Where(goqu.C("id").Eq(entityID), goqu.C("deleted_at").IsNull()).Delete().ToSQL()

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
