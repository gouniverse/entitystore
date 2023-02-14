package entitystore

import (
	"strconv"
	"time"
)

// Attribute type
type Attribute struct {
	id             string
	entityID       string
	attributeKey   string
	attributeValue string
	createdAt      time.Time
	updatedAt      time.Time
	st             *Store
}

func (a *Attribute) ToMap() map[string]any {
	entry := map[string]any{}
	entry["id"] = a.ID()
	entry["entity_id"] = a.EntityID()
	entry["attribute_key"] = a.AttributeKey()
	entry["attribute_value"] = a.AttributeValue()
	entry["created_at"] = a.CreatedAt()
	entry["updated_at"] = a.UpdatedAt()
	return entry
}

func (a *Attribute) ID() string {
	return a.id
}

func (a *Attribute) EntityID() string {
	return a.entityID
}

func (a *Attribute) AttributeKey() string {
	return a.attributeKey
}

func (a *Attribute) AttributeValue() string {
	return a.attributeValue
}

func (a *Attribute) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Attribute) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Attribute) SetID(id string) *Attribute {
	a.id = id
	return a
}

func (a *Attribute) SetEntityID(entityID string) *Attribute {
	a.entityID = entityID
	return a
}

func (a *Attribute) SetAttributeKey(attributeKey string) *Attribute {
	a.attributeKey = attributeKey
	return a
}

func (a *Attribute) SetAttributeValue(attributeValue string) *Attribute {
	a.attributeValue = attributeValue
	return a
}

func (a *Attribute) SetCreatedAt(createdAt time.Time) *Attribute {
	a.createdAt = createdAt
	return a
}

func (a *Attribute) SetUpdatedAt(updatedAt time.Time) *Attribute {
	a.updatedAt = updatedAt
	return a
}

// GetInt returns the value as int
func (a *Attribute) GetInt() (int64, error) {
	return strconv.ParseInt(a.AttributeValue(), 10, 64)
}

// GetFloat returns the value as float
func (a *Attribute) GetFloat() (float64, error) {
	f64Value, err := strconv.ParseFloat(a.AttributeValue(), 100)
	return f64Value, err
}

// GetString returns the value as string
func (a *Attribute) GetString() string {
	return a.AttributeValue()
}

// SetFloat sets a float value
func (a *Attribute) SetFloat(value float64) bool {
	a.attributeValue = strconv.FormatFloat(value, 'f', 30, 64)
	return true
}

// SetInt sets a int value
func (a *Attribute) SetInt(value int64) bool {
	a.attributeValue = strconv.FormatInt(value, 10)
	return true
}

// SetString serializes the values
func (a *Attribute) SetString(value string) bool {
	a.attributeValue = value
	return true
}
