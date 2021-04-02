package entitystore

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

// Attribute type
type Attribute struct {
	ID             string     `gorm:"type:varchar(40);column:id;primary_key;"`
	EntityID       string     `gorm:"type:varchar(40);column:entity_id;"`
	AttributeKey   string     `gorm:"type:varchar(255);column:attribute_key;DEFAULT NULL;"`
	AttributeValue string     `gorm:"type:longtext;column:attribute_value;"`
	CreatedAt      time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt      time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt      *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`
}

// BeforeCreate adds UID to model
func (a *Attribute) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	a.ID = uuid
	return nil
}

// GetInterface de-serializes the values
func (a *Attribute) GetInterface() interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(a.AttributeValue), &value)

	if err != nil {
		panic("JSOB error unmarshaliibg Attribute" + err.Error())
	}

	return value
}

// GetInt returns the value as int
func (a *Attribute) GetInt() (int, error) {
	return strconv.Atoi(a.AttributeValue)
}

// GetFloat returns the value as float
func (a *Attribute) GetFloat() (float32, error) {
	f64Value, err := strconv.ParseFloat(a.AttributeValue, 32)
	return float32(f64Value), err
}

// GetString returns the value as string
func (a *Attribute) GetString() string {
	return a.AttributeValue
}

// SetInterface serializes the values
func (a *Attribute) SetInterface(value interface{}) bool {
	bytes, err := json.Marshal(value)

	if err != nil {
		return false
	}

	a.AttributeValue = string(bytes)

	return true
}

// SetString serializes the values
func (a *Attribute) SetString(value string) bool {
	a.AttributeValue = value

	return true
}
