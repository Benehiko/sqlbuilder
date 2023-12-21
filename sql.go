package sqlbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	FromQuery[T any] interface {
		From(from string) T
	}
	Alias[T any] interface {
		As(alias string) T
	}
	Order[T any] interface {
		OrderBy(orderBy OrderBy, columns ...string) T
	}
	JoinOn interface {
		On(column string, joinColumn string) SelectFromQuery
	}

	Query interface {
		SelectQuery
		InsertQuery
		UpdateQuery
		DeleteQuery
	}
	queryHelper interface {
		GetTable() string
		GetPosition() *int
		GetWhere() *WhereCondition
		GetColumns() []string
		GetReturning() []string
		GetAlias() string
		GetJoins() []*Join
		GetOrderBy() *Sort
		GetParent() any
	}
)

type OrderBy string

const (
	Asc  OrderBy = "ASC"
	Desc OrderBy = "DESC"
)

type Sort struct {
	columns []string
	orderBy OrderBy
}

func NewSelectQuery(columns ...string) FromQuery[SelectFromQuery] {
	sb := &SelectBuilder{
		pos: 1,
	}
	sb.Select(columns...)
	return sb
}

func NewInsertQuery(columns ...string) Into[InsertIntoQuery] {
	ib := &InsertBuilder{}
	ib.Insert(columns...)
	return ib
}

func NewUpdateQuery(table string) UpdateSetQuery {
	ub := &UpdateBuilder{}
	return ub.Update(table)
}

func NewDeleteQuery() FromQuery[DeleteFromQuery] {
	db := &DeleteBuilder{}
	return db.Delete()
}

func WhereOperator[T Operator](next *WhereCondition, operator T, column string, lo LogicalOperator) {
	n := &WhereCondition{
		Op:      operator,
		ColumnA: column,
	}

	tmp := next
	for tmp.next != nil {
		tmp = tmp.next
	}
	tmp.nextOp = lo
	tmp.next = n
}

func WhereSQLHelper(current *WhereCondition, pos *int, sb *strings.Builder) {
	sb.WriteString(current.ColumnA)

	switch op := any(current.Op.get()).(type) {
	case BasicOperator:
		sb.WriteString(" ")
		sb.WriteString(string(op))
		sb.WriteString(" ")
		sb.WriteString("$" + strconv.FormatInt(int64(*pos), 10))
		*pos++
	case SpecialOperator:
		count, o := op()
		sb.WriteString(" ")
		sb.WriteString(o)
		sb.WriteString(" ")

		if count == 0 {
			sb.WriteString("()")
			return
		}

		sb.WriteString("(")
		for i := 0; i < count; i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("$" + strconv.FormatInt(int64(*pos), 10))
			*pos += 1
		}
		sb.WriteString(")")

	}

	if current.next != nil {
		sb.WriteString(" " + string(current.nextOp) + " ")
		WhereSQLHelper(current.next, pos, sb)
	}
}

func WhereSQL[T queryHelper](q T, sb *strings.Builder) {
	if q.GetWhere() == nil {
		return
	}
	sb.WriteString(" WHERE ")
	WhereSQLHelper(q.GetWhere(), q.GetPosition(), sb)
}

func ReturningSQL[T queryHelper](q T, sb *strings.Builder) {
	if len(q.GetReturning()) != 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.TrimSpace(strings.Join(q.GetReturning(), ", ")))
	}
}

func OrderBySQL[T queryHelper](q T, sb *strings.Builder) {
	if q.GetOrderBy() != nil {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.TrimSpace(strings.Join(q.GetOrderBy().columns, ", ")))
		sb.WriteString(" ")
		sb.WriteString(string(q.GetOrderBy().orderBy))
	}
}

func SQL[T any](q T) string {
	switch q := any(q).(type) {
	case queryHelper:
		var sb strings.Builder
		pos := q.GetPosition()

		if *pos == 0 {
			*pos = 1
		}

		switch any(q).(type) {
		case SelectQuery:

			switch p := q.GetParent().(type) {
			case InsertIntoQuery:
				sb.WriteString(p.SQL())
				sb.WriteString(" ")
			}

			sb.WriteString("SELECT")
			if len(q.GetColumns()) == 0 {
				sb.WriteString(" * ")
			} else {
				sb.WriteString(" ")
				for i, c := range q.GetColumns() {
					if i > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(c)
				}
			}
			sb.WriteString(" FROM ")
			sb.WriteString(q.GetTable())

			if alias := q.GetAlias(); alias != "" {
				sb.WriteString(" AS ")
				sb.WriteString(alias)
			}

			if q.GetJoins() != nil {
				for _, join := range q.GetJoins() {
					sb.WriteString(join.SQL())
				}
			}

			WhereSQL(q, &sb)
			OrderBySQL(q, &sb)

		case InsertQuery:
			sb.WriteString("INSERT")
			sb.WriteString(q.GetTable())

			sb.WriteString(" (")
			sb.WriteString(strings.TrimSpace(strings.Join(q.GetColumns(), ", ")))
			sb.WriteString(")")

			sb.WriteString(" VALUES ")
			sb.WriteString("(")

			cols := make([]string, len(q.GetColumns()))
			for i := 0; i < len(q.GetColumns()); i++ {
				cols = append(cols, fmt.Sprintf("$%d", *pos))
				*pos++
			}
			sb.WriteString(strings.TrimSuffix(strings.Join(cols, ", "), ", "))
			sb.WriteString(")")

			ReturningSQL(q, &sb)

		case UpdateQuery:
			sb.WriteString("UPDATE")
			sb.WriteString(" ")
			sb.WriteString(q.GetTable())
			sb.WriteString(" SET ")

			cols := make([]string, len(q.GetColumns()))
			for i, c := range q.GetColumns() {
				cols[i] = fmt.Sprintf("%s = $%d", c, *pos)
				*pos++
			}
			sb.WriteString(strings.TrimSuffix(strings.Join(cols, ", "), ", "))
			WhereSQL(q, &sb)
			ReturningSQL(q, &sb)

		case DeleteQuery:
			sb.WriteString("DELETE")
			sb.WriteString(" FROM ")
			sb.WriteString(q.GetTable())
			WhereSQL(q, &sb)
		}

		return sb.String()
	}
	return ""
}
