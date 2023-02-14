package entitystore

import (
	"time"

	"github.com/gouniverse/uid"
)

// AttributeCreate creates a new attribute
func (st *Store) AttributeCreate(entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	var newAttribute = st.NewAttribute(NewAttributeOptions{
		ID:             uid.HumanUid(),
		EntityID:       entityID,
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	return st.AttributeInsert(*newAttribute)
}

func (st *Store) attributeCreateWithTransactionOrDB(db txOrDB, entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	var newAttribute = st.NewAttribute(NewAttributeOptions{
		ID:             uid.HumanUid(),
		EntityID:       entityID,
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	return st.attributeInsertWithTransactionOrDB(db, *newAttribute)
}
