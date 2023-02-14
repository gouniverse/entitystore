package entitystore

import "strconv"

// AttributeSetInt creates a new attribute or updates existing
func (st *Store) AttributeSetInt(entityID string, attributeKey string, attributeValue int64) error {
	attributeValueAsString := strconv.FormatInt(attributeValue, 10)
	return st.AttributeSetString(entityID, attributeKey, attributeValueAsString)
}
