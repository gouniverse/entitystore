package entitystore

import "time"

func (st *Store) AttributeFromMap(attributeMap map[string]any) *Attribute {
	opts := NewAttributeOptions{}

	if id, exists := attributeMap["id"]; exists {
		opts.ID = id.(string)
	}
	if entityID, exists := attributeMap["entity_id"]; exists {
		opts.EntityID = entityID.(string)
	}
	if attributeKey, exists := attributeMap["attribute_key"]; exists {
		opts.AttributeKey = attributeKey.(string)
	}
	if attributeValue, exists := attributeMap["attribute_value"]; exists {
		opts.AttributeValue = attributeValue.(string)
	}
	if createdAt, exists := attributeMap["created_at"]; exists {
		opts.CreatedAt = createdAt.(time.Time)
	}
	if updatedAt, exists := attributeMap["updated_at"]; exists {
		opts.UpdatedAt = updatedAt.(time.Time)
	}

	return st.NewAttribute(opts)
}
