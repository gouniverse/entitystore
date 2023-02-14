package entitystore

import "github.com/golang-module/carbon/v2"

func (st *Store) NewEntityFromMap(entityMap map[string]string) *Entity {
	opts := NewEntityOptions{}
	if id, exists := entityMap["id"]; exists {
		opts.ID = id
	}
	if entityType, exists := entityMap["entity_type"]; exists {
		opts.Type = entityType
	}
	if entityHandle, exists := entityMap["entity_handle"]; exists {
		opts.Handle = entityHandle
	}
	if createdAt, exists := entityMap["created_at"]; exists {
		opts.CreatedAt = carbon.Parse(createdAt, carbon.UTC).ToStdTime()
	}
	if updatedAt, exists := entityMap["updated_at"]; exists {
		opts.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).ToStdTime()
	}

	return st.NewEntity(opts)
}
