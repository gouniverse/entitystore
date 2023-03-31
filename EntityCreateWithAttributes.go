package entitystore

// EntityCreateWithAttributes func
func (st *Store) EntityCreateWithAttributes(entityType string, attributes map[string]string) (*Entity, error) {
	// Note the use of tx as the database handle once you are within a transaction
	// tx, err := st.db.Begin()
	err := st.database.BeginTransaction()

	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			st.database.RollbackTransaction()
		}
	}()

	entity, err := st.EntityCreate(entityType)

	if err != nil {
		st.database.RollbackTransaction()
		return nil, err
	}

	for k, v := range attributes {
		_, err := st.AttributeCreate(entity.ID(), k, v)

		if err != nil {
			st.database.RollbackTransaction()
			return nil, err
		}
	}

	err = st.database.CommitTransaction()

	if err != nil {
		st.database.RollbackTransaction()
		return nil, err
	}

	return entity, nil
}
