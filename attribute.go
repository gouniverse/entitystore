package entitystore

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

// Attribute type
type attribute struct {
	ID             string     `gorm:"type:varchar(40);column:id;primary_key;"`
	EntityID       string     `gorm:"type:varchar(40);column:entity_id;"`
	AttributeKey   string     `gorm:"type:varchar(255);column:attribute_key;DEFAULT NULL;"`
	AttributeValue string     `gorm:"type:longtext;column:attribute_value;"`
	CreatedAt      time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt      time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt      *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`
}

// BeforeCreate adds UID to model
func (a *attribute) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	a.ID = uuid
	return nil
}

// GetValue de-serializes the values
func (a *attribute) GetAny() interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(a.AttributeValue), &value)

	if err != nil {
		panic("JSOB error unmarshaliibg attribute" + err.Error())
	}

	return value
}

// GetString returns the value as string
func (a *attribute) GetInt() (int, error) {
	return strconv.Atoi(a.AttributeValue)
}

// GetFloat returns the value as string
func (a *attribute) GetFloat() (float32, error) {
	f64Value, err := strconv.ParseFloat(a.AttributeValue, 32)
	return float32(f64Value), err
}

// GetString returns the value as string
func (a *attribute) GetString() string {
	return a.AttributeValue
}

// SetAny serializes the values
func (a *attribute) SetAny(value interface{}) bool {
	bytes, err := json.Marshal(value)

	if err != nil {
		return false
	}

	a.AttributeValue = string(bytes)

	return true
}

// SetString serializes the values
func (a *attribute) SetString(value string) bool {
	a.AttributeValue = value

	return true
}
