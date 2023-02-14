package entitystore

// AttributeSetString creates a new entity
func (st *Store) AttributeSetString(entityID string, attributeKey string, attributeValue string) error {
	attr, err := st.AttributeFind(entityID, attributeKey)

	if err != nil {
		return err
	}

	if attr == nil {
		attr, err := st.AttributeCreate(entityID, attributeKey, attributeValue)
		if err != nil {
			return err
		}
		if attr != nil {
			return nil
		}
		return err
	}

	attr.SetString(attributeValue)

	return st.AttributeUpdate(*attr)
}
