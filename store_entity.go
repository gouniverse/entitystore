package entitystore

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// EntityAttributeList list all atributes of an entity
func (st *Store) EntityAttributeList(entityID string) []Attribute {
	var attrs []Attribute

	result := st.db.Table(st.attributeTableName).Find(&attrs, "entity_id=?", entityID)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return attrs
}

// EntityCount counts entities
func (st *Store) EntityCount(entityType string) uint64 {
	var count int64
	st.db.Table(st.entityTableName).Where("type=?", entityType).Count(&count)
	return uint64(count)
}

// EntityCreate creates a new entity
func (st *Store) EntityCreate(entityType string) *Entity {
	entity := &Entity{Type: entityType, Status: "active", st: st}

	dbResult := st.db.Table(st.entityTableName).Create(&entity)

	if dbResult.Error != nil {
		return nil
	}

	return entity
}

// EntityCreateWithAttributes func
func (st *Store) EntityCreateWithAttributes(entityType string, attributes map[string]interface{}) *Entity {
	// Note the use of tx as the database handle once you are within a transaction
	tx := st.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil
	}

	//return tx.Commit().Error

	entity := &Entity{Type: entityType, Status: EntityStatusActive, st: st}

	dbResult := tx.Table(st.entityTableName).Create(&entity)

	if dbResult.Error != nil {
		tx.Rollback()
		return nil
	}

	//entityAttributes := make([]EntityAttribute, 0)
	for k, v := range attributes {
		ea := Attribute{EntityID: entity.ID, AttributeKey: k} //, AttributeValue: value}
		ea.SetInterface(v)

		dbResult2 := tx.Table(st.attributeTableName).Create(&ea)
		if dbResult2.Error != nil {
			tx.Rollback()
			return nil
		}
	}

	err := tx.Commit().Error

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

	resultEntity := st.db.Table(st.entityTableName).First(&ent, "type=? AND handle=?", entityType, entityHandle)

	if resultEntity.Error != nil {
		if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	ent.st = st // Add store reference

	return ent
}

// EntityFindByID finds an entity by ID
func (st *Store) EntityFindByID(entityID string) *Entity {
	if entityID == "" {
		return nil
	}

	ent := &Entity{}

	resultEntity := st.db.Table(st.entityTableName).First(&ent, "id=?", entityID)

	if resultEntity.Error != nil {
		if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	ent.st = st // Add store reference

	return ent
}

// EntityFindByAttribute finds an entity by attribute
func (st *Store) EntityFindByAttribute(entityType string, attributeKey string, attributeValue string) *Entity {
	attr := &Attribute{}

	subQuery := st.db.Table(st.entityTableName).Model(&Entity{}).Select("id").Where("type = ?", entityType)
	result := st.db.Table(st.attributeTableName).First(&attr, "entity_id IN (?) AND attribute_key=? AND attribute_value=?", subQuery, attributeKey, attributeValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	// DEBUG: log.Println(entityAttribute)

	ent := &Entity{}

	resultEntity := st.db.Table(st.entityTableName).First(&ent, "id=?", attr.EntityID)

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	ent.st = st // Add store reference

	return ent
}

// EntityList lists entities
func (st *Store) EntityList(entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) []Entity {
	entityList := []Entity{}

	result := st.db.Table(st.entityTableName).Where("type=?", entityType).Order(orderBy + " " + sort).Offset(int(offset)).Limit(int(perPage)).Find(&entityList)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	for k := range entityList {
		entityList[k].st = st
	}

	return entityList
}

// EntityListByAttribute finds an entity by attribute
func (st *Store) EntityListByAttribute(entityType string, attributeKey string, attributeValue string) []Entity {
	//entityAttributes := []EntityAttribute{}
	var entityIDs []string

	subQuery := st.db.Table(st.entityTableName).Model(&Entity{}).Select("id").Where("type = ?", entityType)
	result := st.db.Table(st.attributeTableName).Model(&Attribute{}).Select("entity_id").Find(&entityIDs, "entity_id IN (?) AND attribute_key=? AND attribute_value=?", subQuery, attributeKey, attributeValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	// DEBUG: log.Println(result)

	entityList := []Entity{}

	resultEntity := st.db.Table(st.entityTableName).Where("id IN (?)", entityIDs).Find(&entityList)

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	for k := range entityList {
		entityList[k].st = st
	}

	return entityList
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

	if ent==nil{
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
		attrTrash := AttributeTrash {
			ID:        attr.ID,
			EntityID:    attr.EntityID,
			AttributeKey:    attr.AttributeKey,
			AttributeValue:      attr.AttributeValue,
			CreatedAt: attr.CreatedAt,
			UpdatedAt: attr.UpdatedAt,
			DeletedAt: time.Now(),
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
