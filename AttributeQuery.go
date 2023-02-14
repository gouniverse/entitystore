package entitystore

import "github.com/doug-martin/goqu/v9"

type AttributeQueryOptions struct {
	ID           string
	IDs          []string
	EntityID     string
	AttributeKey string
	Limit        uint64
	Offset       uint64
	SortBy       string
	SortOrder    string // asc / dec
	CountOnly    bool
}

func (st *Store) AttributeQuery(options AttributeQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C("id").In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C("id").Eq(options.ID))
	}

	sortByColumn := "id"
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	if sortOrder == "asc" {
		q = q.Order(goqu.I(sortByColumn).Asc())
	} else {
		q = q.Order(goqu.I(sortByColumn).Desc())
	}

	if options.EntityID != "" {
		q = q.Where(goqu.C("entity_id").Eq(options.EntityID))
	}

	if options.AttributeKey != "" {
		q = q.Where(goqu.C("attribute_key").Eq(options.AttributeKey))
	}

	q = q.Offset(uint(options.Offset))

	if options.Limit != 0 {
		q = q.Limit(uint(options.Limit))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	return q.Select()
}
