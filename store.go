package entitystore

import (
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Store defines an entity store
type Store struct {
	entityTableName    string
	attributeTableName string
	db                 *gorm.DB
	automigrateEnabled bool
}

// StoreOption options for the vault store
type StoreOption func(*Store)

// WithAutoMigrate sets the table name for the cache store
func WithAutoMigrate(automigrateEnabled bool) StoreOption {
	return func(s *Store) {
		s.automigrateEnabled = automigrateEnabled
	}
}

// WithDriverAndDNS sets the driver and the DNS for the database for the cache store
func WithDriverAndDNS(driverName string, dsn string) StoreOption {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return func(s *Store) {
		s.db = db
	}
}

// WithGormDb sets the GORM database for the cache store
func WithGormDb(db *gorm.DB) StoreOption {
	return func(s *Store) {
		s.db = db
	}
}

// WithEntityTableName sets the table name for the cache store
func WithEntityTableName(entityTableName string) StoreOption {
	return func(s *Store) {
		s.entityTableName = entityTableName
	}
}

// WithAttributeTableName sets the table name for the cache store
func WithAttributeTableName(attributeTableName string) StoreOption {
	return func(s *Store) {
		s.attributeTableName = attributeTableName
	}
}

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) *Store {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.entityTableName == "" {
		log.Panic("Entity store: entityTableName is required")
	}

	if store.entityTableName == "" {
		log.Panic("Entity store: attributeTableName is required")
	}

	if store.automigrateEnabled == true {
		store.AutoMigrate()
	}

	return store
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() {
	st.db.Table(st.entityTableName).AutoMigrate(&Entity{})
	st.db.Table(st.attributeTableName).AutoMigrate(&Attribute{})
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

// EntityDelete deletes an entity and all attributes
func (st *Store) EntityDeleteSoft(entityID string) bool {
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

	if err := tx.Where("entity_id=?", entityID).Table(st.attributeTableName).Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	if err := tx.Where("id=?", entityID).Table(st.entityTableName).Update("deleted_at", time.Now()).Error; err != nil {
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
	
	for k, _ := range entityList {
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
	
	for k, _ := range entityList {
		entityList[k].st = st
	}

	return entityList
}

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) *Attribute {
	attr := &Attribute{EntityID: entityID, AttributeKey: attributeKey, AttributeValue: attributeValue}

	dbResult := st.db.Table(st.attributeTableName).Create(&attr)

	if dbResult.Error != nil {
		return nil
	}

	return attr
}

// AttributeCreateInterface creates a new attribute
func (st *Store) AttributeCreateInterface(entityID string, attributeKey string, attributeValue interface{}) *Attribute {
	attr := &Attribute{EntityID: entityID, AttributeKey: attributeKey}
	attr.SetInterface(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Create(&attr)

	if dbResult.Error != nil {
		return nil
	}

	return attr
}

// AttributeFind finds an entity by ID
func (st *Store) AttributeFind(entityID string, attributeKey string) *Attribute {
	attr := &Attribute{}

	result := st.db.Table(st.attributeTableName).First(&attr, "entity_id=? AND attribute_key=?", entityID, attributeKey)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	return attr
}

// AttributeGet the name of the User table
func (st *Store) AttributeGet(entityID string, attributeKey string) *Attribute {
	attr := &Attribute{}

	result := st.db.Table(st.attributeTableName).First(&attr, "entity_id=? AND attribute_key=?", entityID, attributeKey)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return attr
}



// AttributeSetFloat creates a new entity
func (st *Store) AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreateInterface(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetFloat(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// AttributeSetInt creates a new entity
func (st *Store) AttributeSetInt(entityID string, attributeKey string, attributeValue int64) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreateInterface(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetInt(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// AttributeSetInterface creates a new entity
func (st *Store) AttributeSetInterface(entityID string, attributeKey string, attributeValue interface{}) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreateInterface(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetInterface(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// AttributeSetString creates a new entity
func (st *Store) AttributeSetString(entityID string, attributeKey string, attributeValue string) bool {
	attr := st.AttributeFind(entityID, attributeKey)

	if attr == nil {
		attr = st.AttributeCreate(entityID, attributeKey, attributeValue)
		if attr != nil {
			return true
		}
		return false
	}

	attr.SetString(attributeValue)

	dbResult := st.db.Table(st.attributeTableName).Save(attr)
	if dbResult.Error != nil {
		return false
	}

	return true
}

// AttributesSet upserts and entity attribute
func (st *Store) AttributesSet(entityID string, attributes map[string]interface{}) bool {
	tx := st.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false
	}

	for k, v := range attributes {
		attr := st.AttributeFind(entityID, k)

		if attr == nil {
			attr = &Attribute{EntityID: entityID, AttributeKey: k}
			attr.SetInterface(v)

			dbResult := tx.Table(st.attributeTableName).Create(&attr)
			if dbResult.Error != nil {
				tx.Rollback()
				return false
			}

		}

		attr.SetInterface(v)
		dbResult := tx.Table(st.attributeTableName).Save(attr)
		if dbResult.Error != nil {
			return false
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return false
	}

	return true

}
