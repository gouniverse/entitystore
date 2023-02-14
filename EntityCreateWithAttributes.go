package entitystore

// EntityCreateWithAttributes func
func (st *Store) EntityCreateWithAttributes(entityType string, attributes map[string]string) (*Entity, error) {
	// Note the use of tx as the database handle once you are within a transaction
	tx, err := st.db.Begin()

	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	entity, err := st.entityCreateWithTransactionOrDB(tx, entityType)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for k, v := range attributes {
		_, err := st.attributeCreateWithTransactionOrDB(tx, entity.ID(), k, v)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return entity, nil
}
