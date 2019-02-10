package query

// Query defines information about query generated by query builder.
type Query struct {
	built        bool
	Collection   string
	SelectClause SelectClause
	JoinClause   []JoinClause
	WhereClause  FilterClause
	GroupClause  GroupClause
	SortClause   []SortClause
	OffsetClause Offset
	LimitClause  Limit
}

// Select filter fields to be selected from database.
func (q Query) Select(fields ...string) Query {
	q.SelectClause = NewSelect(fields...)
	return q
}

func (q Query) Distinct() Query {
	q.SelectClause.OnlyDistinct = true
	return q
}

// Join current collection with other collection.
func (q Query) Join(collection string) Query {
	return q.JoinOn(collection, "", "")
}

// Join current collection with other collection.
func (q Query) JoinOn(collection string, from string, to string) Query {
	return q.JoinWith("JOIN", collection, from, to)
}

// JoinWith current collection with other collection with custom join mode.
func (q Query) JoinWith(mode string, collection string, from string, to string) Query {
	NewJoinWith(mode, collection, from, to).Build(&q) // TODO: ensure this always called last

	return q
}

func (q Query) JoinFragment(expr string, args ...interface{}) Query {
	NewJoinFragment(expr, args...).Build(&q) // TODO: ensure this always called last

	return q
}

func (q Query) Where(filters ...FilterClause) Query {
	q.WhereClause = q.WhereClause.And(filters...)
	return q
}

func (q Query) OrWhere(filters ...FilterClause) Query {
	q.WhereClause = q.WhereClause.Or(FilterAnd(filters...))
	return q
}

func (q Query) Group(fields ...string) Query {
	q.GroupClause.Fields = fields
	return q
}

func (q Query) Having(filters ...FilterClause) Query {
	q.GroupClause.Filter = q.GroupClause.Filter.And(filters...)
	return q
}

func (q Query) OrHaving(filters ...FilterClause) Query {
	q.GroupClause.Filter = q.GroupClause.Filter.Or(FilterAnd(filters...))
	return q
}

func (q Query) Sort(fields ...string) Query {
	return q.SortAsc(fields...)
}

func (q Query) SortAsc(fields ...string) Query {
	sorts := make([]SortClause, len(fields))
	for i := range fields {
		sorts[i] = NewSortAsc(fields[i])
	}

	q.SortClause = append(q.SortClause, sorts...)
	return q
}

func (q Query) SortDesc(fields ...string) Query {
	sorts := make([]SortClause, len(fields))
	for i := range fields {
		sorts[i] = NewSortDesc(fields[i])
	}

	q.SortClause = append(q.SortClause, sorts...)
	return q
}

// Offset the result returned by database.
func (q Query) Offset(offset Offset) Query {
	q.OffsetClause = offset
	return q
}

// Limit result returned by database.
func (q Query) Limit(limit Limit) Query {
	q.LimitClause = limit
	return q
}

// From create query for collection.
func From(collection string) Query {
	return Query{
		Collection:   collection,
		SelectClause: NewSelect(collection + ".*"),
	}
}

// TODO Join entry

func Where(filters ...FilterClause) Query {
	return Query{
		WhereClause: FilterAnd(filters...),
	}
}

// func Group(fields ...string) Query {
// 	return Query{
// 		GroupClause: GroupClause{
// 			Fields: fields,
// 		},
// 	}
// }
