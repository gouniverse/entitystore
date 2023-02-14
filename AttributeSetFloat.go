package entitystore

import "strconv"

// AttributeSetFloat creates a new attribute or updates existing
func (st *Store) AttributeSetFloat(entityID string, attributeKey string, attributeValue float64) error {
	attributeValueAsString := strconv.FormatFloat(attributeValue, 'f', 30, 64)
	return st.AttributeSetString(entityID, attributeKey, attributeValueAsString)
}
