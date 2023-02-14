package entitystore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// AttributesSet upserts an entity attribute
func (st *Store) AttributesSet(entityID string, attributes map[string]string) error {
	tx, err := st.db.Begin()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback()
			if st.GetDebug() {
				log.Println(err)
			}
		}
	}()

	for k, v := range attributes {
		attr, err := st.AttributeFind(entityID, k)

		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}

			err = tx.Rollback()

			if st.GetDebug() {
				log.Println(err)
			}

			return err
		}

		if attr == nil {
			attr = st.NewAttribute(NewAttributeOptions{ID: uid.HumanUid(), EntityID: entityID, AttributeKey: k, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			attr.SetString(v)

			q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName)
			q = q.Rows(attr.ToMap())
			sqlStr, _, err := q.ToSQL()

			if err != nil {
				if st.GetDebug() {
					log.Println(err)
				}

				err = tx.Rollback()

				if st.GetDebug() {
					log.Println(err)
				}

				return err
			}

			if st.GetDebug() {
				log.Println(sqlStr)
			}

			_, err = tx.Exec(sqlStr)

			if err != nil {
				log.Println(err)
				err = tx.Rollback()

				if st.GetDebug() {
					log.Println(err)
				}

				return err
			}

		}

		attr.SetString(v)
		attr.SetUpdatedAt(time.Now())

		q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
		q = q.Where(goqu.C("id").Eq(attr.ID))
		q = q.Set(attr.ToMap())

		sqlStr, _, err := q.ToSQL()

		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}

			err = tx.Rollback()

			if st.GetDebug() {
				log.Println(err)
			}

			return err
		}

		if st.GetDebug() {
			log.Println(sqlStr)
		}

		_, err = tx.Exec(sqlStr)

		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}

			err = tx.Rollback()

			if st.GetDebug() {
				log.Println(err)
			}

			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		err = tx.Rollback()

		if st.GetDebug() {
			log.Println(err)
		}

		return err
	}

	return nil
}
