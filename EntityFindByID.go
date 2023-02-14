package entitystore

import "errors"

// EntityFindByID finds an entity by ID
func (st *Store) EntityFindByID(entityID string) (*Entity, error) {
	if entityID == "" {
		return nil, errors.New("entity ID cannot be empty")
	}

	list, err := st.EntityList(EntityQueryOptions{
		ID:    entityID,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}
