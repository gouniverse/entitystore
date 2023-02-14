package entitystore

import "time"

func (st *Store) EntityFromMap(entityMap map[string]any) *Entity {
	opts := NewEntityOptions{}
	if id, exists := entityMap["id"]; exists {
		opts.ID = id.(string)
	}
	if entityType, exists := entityMap["entity_type"]; exists {
		opts.Type = entityType.(string)
	}
	if entityHandle, exists := entityMap["entity_handle"]; exists {
		opts.Handle = entityHandle.(string)
	}
	if createdAt, exists := entityMap["created_at"]; exists {
		opts.CreatedAt = createdAt.(time.Time)
	}
	if updatedAt, exists := entityMap["updated_at"]; exists {
		opts.UpdatedAt = updatedAt.(time.Time)
	}
	entity := st.NewEntity(opts)
	return entity
}
