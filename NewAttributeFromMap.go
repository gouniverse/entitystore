package entitystore

import "github.com/golang-module/carbon/v2"

func (st *Store) NewAttributeFromMap(attributeMap map[string]string) *Attribute {
	opts := NewAttributeOptions{}

	if id, exists := attributeMap["id"]; exists {
		opts.ID = id
	}
	if entityID, exists := attributeMap["entity_id"]; exists {
		opts.EntityID = entityID
	}
	if attributeKey, exists := attributeMap["attribute_key"]; exists {
		opts.AttributeKey = attributeKey
	}
	if attributeValue, exists := attributeMap["attribute_value"]; exists {
		opts.AttributeValue = attributeValue
	}
	if createdAt, exists := attributeMap["created_at"]; exists {
		opts.CreatedAt = carbon.Parse(createdAt, carbon.UTC).ToStdTime()
	}
	if updatedAt, exists := attributeMap["updated_at"]; exists {
		opts.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).ToStdTime()
	}

	return st.NewAttribute(opts)
}
