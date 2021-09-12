package entitystore

import (
	"encoding/json"
	"strconv"
	"time"
)

// Attribute type
type Attribute struct {
	ID             string    `db:"id"`
	EntityID       string    `db:"entity_id"`
	AttributeKey   string    `db:"attribute_key"`
	AttributeValue string    `db:"attribute_value"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// GetInt returns the value as int
func (a *Attribute) GetInt() (int64, error) {
	return strconv.ParseInt(a.AttributeValue, 10, 64)
}

// GetFloat returns the value as float
func (a *Attribute) GetFloat() (float64, error) {
	f64Value, err := strconv.ParseFloat(a.AttributeValue, 100)
	return f64Value, err
}

// GetString returns the value as string
func (a *Attribute) GetString() string {
	return a.AttributeValue
}

// SetFloat sets a float value
func (a *Attribute) SetFloat(value float64) bool {
	a.AttributeValue = strconv.FormatFloat(value, 'f', 30, 64)

	return true
}

// SetInt sets a int value
func (a *Attribute) SetInt(value int64) bool {
	a.AttributeValue = strconv.FormatInt(value, 10)

	return true
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
