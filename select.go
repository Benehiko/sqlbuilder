package sqlbuilder

import (
	"strconv"
	"strings"
)

type SelectBuilder struct {
	parent InsertQuery

	table        string
	as           string
	columns      []string
	where        *WhereCondition[BasicOperator]
	whereSpecial *WhereCondition[SpecialOperator]
	orderBy      *Sort
	joins        []*Join
	pos          int
}

type SelectBuilderWhere[T Operator] struct {
	*SelectBuilder
}

var (
	_ SelectQuery                                  = (*SelectBuilder)(nil)
	_ SelectFromWhereOptionsQuery[BasicOperator]   = (*SelectBuilderWhere[BasicOperator])(nil)
	_ SelectFromWhereOptionsQuery[SpecialOperator] = (*SelectBuilderWhere[SpecialOperator])(nil)
)

func (s *SelectBuilder) Select(columns ...string) FromQuery[SelectFromQuery] {
	s.columns = columns
	return s
}

func (s *SelectBuilder) From(table string) SelectFromQuery {
	s.table = table
	return s
}

func (s *SelectBuilder) As(alias string) SelectFromQuery {
	s.as = alias
	return s
}

func (s *SelectBuilder) Where(column string, operator BasicOperator) SelectFromWhereOptionsQuery[BasicOperator] {
	s.where = &WhereCondition[BasicOperator]{
		ColumnA: column,
		Op:      operator,
	}
	return &SelectBuilderWhere[BasicOperator]{
		SelectBuilder: s,
	}
}

func (s *SelectBuilder) WhereSpecial(column string, operator SpecialOperator) SelectFromWhereOptionsQuery[SpecialOperator] {
	s.whereSpecial = &WhereCondition[SpecialOperator]{
		ColumnA: column,
		Op:      operator,
	}

	return &SelectBuilderWhere[SpecialOperator]{
		SelectBuilder: s,
	}
}

func (s *SelectBuilderWhere[T]) And(column string, operator T) SelectFromWhereOptionsQuery[T] {
	switch o := any(operator).(type) {
	case BasicOperator:
		if s.where != nil {
			s.where.nextOp = And
			s.where.next = &WhereCondition[BasicOperator]{
				ColumnA: column,
				Op:      o,
			}
		}
	case SpecialOperator:
		if s.whereSpecial != nil {
			s.whereSpecial.nextOp = And
			s.whereSpecial.next = &WhereCondition[SpecialOperator]{
				ColumnA: column,
				Op:      o,
			}
		}
	}

	return s
}

func (s *SelectBuilderWhere[T]) Or(column string, operator T) SelectFromWhereOptionsQuery[T] {
	switch o := any(operator).(type) {
	case BasicOperator:
		if s.where != nil {
			tail := s.where
			for tail.next != nil {
				tail = tail.next
			}
			tail.nextOp = Or
			tail.next = &WhereCondition[BasicOperator]{
				ColumnA: column,
				Op:      o,
			}
		}
	case SpecialOperator:
		if s.whereSpecial != nil {
			tail := s.whereSpecial
			for tail.next != nil {
				tail = tail.next
			}
			tail.nextOp = Or
			tail.next = &WhereCondition[SpecialOperator]{
				ColumnA: column,
				Op:      o,
			}
		}
	}

	return s
}

func (s *SelectBuilder) OrderBy(orderBy OrderBy, columns ...string) SelectFromQuery {
	s.orderBy = &Sort{
		columns: columns,
		orderBy: orderBy,
	}
	return s
}

func (s *SelectBuilder) InnerJoin(table string) AliasOrJoinOn {
	if s.joins == nil {
		s.joins = make([]*Join, 0)
	}
	s.joins = append(s.joins,
		&Join{
			join:   InnerJoin,
			table:  table,
			parent: s,
		})

	return s.joins[len(s.joins)-1]
}

func (s *SelectBuilder) FullOuterJoin(table string) AliasOrJoinOn {
	if s.joins == nil {
		s.joins = make([]*Join, 0)
	}
	s.joins = append(s.joins, &Join{
		join:   FullOuterJoin,
		table:  table,
		parent: s,
	})
	return s.joins[len(s.joins)-1]
}

func (s *SelectBuilder) LeftJoin(table string) AliasOrJoinOn {
	if s.joins == nil {
		s.joins = make([]*Join, 0)
	}
	s.joins = append(s.joins, &Join{
		join:   LeftJoin,
		table:  table,
		parent: s,
	})
	return s.joins[len(s.joins)-1]
}

func (s *SelectBuilder) RightJoin(table string) AliasOrJoinOn {
	if s.joins == nil {
		s.joins = make([]*Join, 0)
	}
	s.joins = append(s.joins, &Join{
		join:   RightJoin,
		table:  table,
		parent: s,
	})
	return s.joins[len(s.joins)-1]
}

func (sb *SelectBuilder) writeWhereSpecialHelper(b *strings.Builder, next *WhereCondition[SpecialOperator]) {
	count, op := next.Op()
	b.WriteString(next.ColumnA)
	b.WriteString(" ")
	b.WriteString(op)
	b.WriteString(" ")

	if count == 0 {
		b.WriteString("()")
		return
	}

	b.WriteString("(")
	for i := 0; i < count; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("$" + strconv.FormatInt(int64(sb.pos), 10))
		sb.pos += 1
	}
	b.WriteString(")")

	if next.next != nil {
		b.WriteString(" " + string(next.nextOp) + " ")
		sb.writeWhereSpecialHelper(b, next.next)
	}
}

func (sb *SelectBuilder) writeWhereHelper(b *strings.Builder, next *WhereCondition[BasicOperator]) {
	b.WriteString(next.ColumnA)
	b.WriteString(" ")
	b.WriteString(string(next.Op))
	b.WriteString(" ")
	b.WriteString("$" + strconv.FormatInt(int64(sb.pos), 10))
	sb.pos++

	if next.next != nil {
		b.WriteString(" " + string(next.nextOp) + " ")
		sb.writeWhereHelper(b, next.next)
	}
}

func (sb *SelectBuilder) writeWhere(b *strings.Builder) {
	b.WriteString(" WHERE ")

	if sb.whereSpecial != nil {
		sb.writeWhereSpecialHelper(b, sb.whereSpecial)
	} else if sb.where != nil {
		sb.writeWhereHelper(b, sb.where)
	}
}

func (s *SelectBuilder) SQL() string {
	var sb strings.Builder

	if s.pos == 0 {
		s.pos = 1
	}

	if s.parent != nil {
		switch p := s.parent.(type) {
		case InsertIntoQuery:
			sb.WriteString(p.SQL())
			sb.WriteString(" ")
		}
	}

	sb.WriteString("SELECT")
	if len(s.columns) == 0 {
		sb.WriteString(" * ")
	} else {
		sb.WriteString(" ")
		sb.WriteString(strings.TrimSpace(strings.Join(s.columns, ", ")))
	}

	sb.WriteString(" FROM ")
	sb.WriteString(s.table)

	if s.as != "" {
		sb.WriteString(" AS ")
		sb.WriteString(s.as)
	}

	if len(s.joins) > 0 {
		for _, j := range s.joins {
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
			sb.WriteString(string(j.on.Op))
			sb.WriteString(" ")
			sb.WriteString(string(j.on.ColumnB))
		}
	}

	if s.where != nil || s.whereSpecial != nil {
		s.writeWhere(&sb)
	}

	if s.orderBy != nil {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.TrimSpace(strings.Join(s.orderBy.columns, ", ")))
		sb.WriteString(" ")
		sb.WriteString(string(s.orderBy.orderBy))
	}

	return sb.String()
}
