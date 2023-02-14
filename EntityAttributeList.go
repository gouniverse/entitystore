package entitystore

// EntityAttributeList list all attributes of an entity
func (st *Store) EntityAttributeList(entityID string) (attributes []Attribute, err error) {
	return st.AttributeList(AttributeQueryOptions{
		EntityID: entityID,
	})
}
