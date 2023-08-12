package main

type JoinType string

const (
	InnerJoin     JoinType = "INNER JOIN"
	LeftJoin      JoinType = "LEFT JOIN"
	RightJoin     JoinType = "RIGHT JOIN"
	FullOuterJoin JoinType = "FULL OUTER JOIN"
)

type Join struct {
	table string
	on    *WhereCondition
	join  JoinType
	as    string

	parent SelectFromQuery
}

var _ Joins = (*Join)(nil)

func (j *Join) InnerJoin(table string) AliasOrJoinOn {
	j.join = InnerJoin
	j.table = table
	return j
}

func (j *Join) FullOuterJoin(table string) AliasOrJoinOn {
	j.join = FullOuterJoin
	j.table = table
	return j
}

func (j *Join) LeftJoin(table string) AliasOrJoinOn {
	j.join = LeftJoin
	j.table = table
	return j
}

func (j *Join) RightJoin(table string) AliasOrJoinOn {
	j.join = RightJoin
	j.table = table
	return j
}

func (j *Join) As(alias string) Joins {
	j.as = alias
	return j
}

func (j *Join) On(table string, joinColumn string) SelectFromQuery {
	j.on = &WhereCondition{
		Column: table,
		Op:     Equals,
		Value:  joinColumn,
	}
	return j.parent
}
