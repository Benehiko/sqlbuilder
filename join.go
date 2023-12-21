package sqlbuilder

import "strings"

type (
	AliasOrJoinOn interface {
		Alias[JoinOn]
		JoinOn
	}
	JoinType string
	Joins    interface {
		InnerJoin(table string) AliasOrJoinOn
		LeftJoin(table string) AliasOrJoinOn
		RightJoin(table string) AliasOrJoinOn
		FullOuterJoin(table string) AliasOrJoinOn
	}
)

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

func (j *Join) As(alias string) JoinOn {
	j.as = alias
	return j
}

func (j *Join) On(table string, joinColumn string) SelectFromQuery {
	j.on = &WhereCondition{
		ColumnA: table,
		Op:      Equals,
		ColumnB: joinColumn,
	}
	return j.parent
}

func (j *Join) SQL() string {
	var sb strings.Builder
	sb.WriteString(" ")
	sb.WriteString(string(j.join))
	sb.WriteString(" ")
	sb.WriteString(j.table)
	if j.as != "" {
		sb.WriteString(" AS ")
		sb.WriteString(j.as)
	}
	sb.WriteString(" ON ")
	sb.WriteString(j.on.ColumnA)
	sb.WriteString(" ")

	switch o := any(j.on.Op.get()).(type) {
	case BasicOperator:
		sb.WriteString(string(o))
	}
	sb.WriteString(" ")
	sb.WriteString(string(j.on.ColumnB))
	return sb.String()
}
