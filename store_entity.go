package entitystore

import (
	"database/sql"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// EntityAttributeList list all atributes of an entity
func (st *Store) EntityAttributeList(entityID string) ([]Attribute, error) {
	var attrs []Attribute

	sqlStr, _, _ := goqu.From(st.attributeTableName).Order(goqu.I("attribute_key").Asc()).Where(goqu.C("entity_id").Eq(entityID)).Where(goqu.C("deleted_at").IsNull()).Select(Attribute{}).ToSQL()

	log.Println(sqlStr)

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return attrs, err
	}

	for rows.Next() {
		var attr Attribute
		err := rows.Scan(&attr.AttributeKey, &attr.AttributeValue, &attr.CreatedAt, &attr.DeletedAt, &attr.ID, &attr.UpdatedAt)
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
	entity := &Entity{Type: entityType, Status: "active", st: st}

	sqlStr, _, _ := goqu.Insert(st.attributeTableName).Rows(entity).ToSQL()

	log.Println(sqlStr)

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		log.Println(err)
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

	sqlStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("entity_type").Eq(entityType), goqu.C("entity_handle").Eq(entityHandle), goqu.C("deleted_at").IsNull()).Select(Entity{}).ToSQL()

	log.Println(sqlStr)

	err := st.db.QueryRow(sqlStr).Scan(&ent.CreatedAt, &ent.DeletedAt, &ent.ID, &ent.Type, &ent.UpdatedAt)
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
func (st *Store) EntityFindByID(entityID string) *Entity {
	if entityID == "" {
		return nil
	}

	ent := &Entity{}

	sqlStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("id").Eq(entityID), goqu.C("deleted_at").IsNull()).Select(Entity{}).ToSQL()

	log.Println(sqlStr)

	err := st.db.QueryRow(sqlStr).Scan(&ent.CreatedAt, &ent.DeletedAt, &ent.ID, &ent.Type, &ent.UpdatedAt)
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

// EntityFindByAttribute finds an entity by attribute
func (st *Store) EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) *Entity {
	subqueryStr, _, _ := goqu.From(st.entityTableName).Where(goqu.C("type").Eq(entityType), goqu.C("deleted_at").IsNull()).Select("id").ToSQL()
	sqlStr, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").In(subqueryStr), goqu.C("attribute_key").Eq(attributeKey), goqu.C("atribute_value").Eq(attributeValue), goqu.C("deleted_at").IsNull()).Select("entity_id").ToSQL()

	var entityID string
	err := st.db.QueryRow(sqlStr).Scan(&entityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	return st.EntityFindByID(entityID)
}

// EntityList lists entities
func (st *Store) EntityList(entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) ([]Entity, error) {
	entityList := []Entity{}

	// result := st.db.Table(st.entityTableName).Where("type=?", entityType).Order(orderBy + " " + sort).Offset(int(offset)).Limit(int(perPage)).Find(&entityList)

	// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 	return nil
	// }

	// for k := range entityList {
	// 	entityList[k].st = st
	// }

	sqlStr, _, _ := goqu.From(st.attributeTableName).Order(goqu.I("id").Asc()).Where(goqu.C("type").Eq(entityType)).Where(goqu.C("deleted_at").IsNull()).Offset(uint(offset)).Limit(uint(perPage)).Select(Entity{}).ToSQL()

	log.Println(sqlStr)

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return entityList, err
	}

	for rows.Next() {
		var ent Entity
		err := rows.Scan(&ent.CreatedAt, &ent.DeletedAt, &ent.ID, &ent.Type, &ent.UpdatedAt)
		if err != nil {
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
	sqlStr, _, _ := goqu.From(st.attributeTableName).Where(goqu.C("entity_id").In(subqueryStr), goqu.C("attribute_key").Eq(attributeKey), goqu.C("atribute_value").Eq(attributeValue), goqu.C("deleted_at").IsNull()).Select("entity_id").ToSQL()

	log.Println(sqlStr)

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

	sqlStr, _, _ = goqu.From(st.attributeTableName).Order(goqu.I("id").Asc()).Where(goqu.C("id").In(entityIDs)).Where(goqu.C("deleted_at").IsNull()).Select(Entity{}).ToSQL()

	log.Println(sqlStr)

	rows, err = st.db.Query(sqlStr)

	if err != nil {
		return entityList, err
	}

	for rows.Next() {
		var ent Entity
		err := rows.Scan(&ent.CreatedAt, &ent.DeletedAt, &ent.ID, &ent.Type, &ent.UpdatedAt)
		if err != nil {
			return entityList, err
		}
		ent.st = st
		entityList = append(entityList, ent)
	}

	return entityList, nil
}

// EntityTrash moves an entity and all attributes to the trash bin
func (st *Store) EntityTrash(entityID string) bool {
	if entityID == "" {
		return false
	}

	// Note the use of tx as the database handle once you are within a transaction
	tx := st.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println(err)
		return false
	}

	ent := st.EntityFindByID(entityID)

	if ent == nil {
		tx.Rollback()
		log.Println("Entity not found")
		return false
	}

	entTrash := EntityTrash{
		ID:        ent.ID,
		Status:    ent.Status,
		Type:      ent.Type,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
		DeletedAt: time.Now(),
	}

	if err := tx.Table(st.entityTrashTableName).Create(entTrash).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	attrs := st.EntityAttributeList(entityID)

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

		if err := tx.Table(st.attributeTrashTableName).Create(attrTrash).Error; err != nil {
			tx.Rollback()
			log.Println(err)
			return false
		}
	}

	if err := tx.Where("entity_id=?", entityID).Table(st.attributeTableName).Delete(&Attribute{}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	if err := tx.Where("id=?", entityID).Table(st.entityTableName).Delete(&Entity{}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	err := tx.Commit().Error

	if err == nil {
		return true
	}

	log.Println(err)

	return false
}
