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

var _ SelectQuery = (*SelectBuilder)(nil)

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

func (s *SelectBuilder) Where(column string, operator BasicOperator) SelectFromQuery {
	s.where = &WhereCondition[BasicOperator]{
		ColumnA: column,
		Op:      operator,
	}
	return s
}

func (s *SelectBuilder) WhereSpecial(column string, operator SpecialOperator) SelectFromQuery {
	s.whereSpecial = &WhereCondition[SpecialOperator]{
		ColumnA: column,
		Op:      operator,
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

func (sb *SelectBuilder) writeWhere(b *strings.Builder) {
	b.WriteString(" WHERE ")

	if sb.whereSpecial != nil {
		count, op := sb.whereSpecial.Op()
		b.WriteString(sb.whereSpecial.ColumnA)
		b.WriteString(" ")
		b.WriteString(op)
		b.WriteString(" ")

		if count == 0 {
			return
		}

		b.WriteString("(")
		for i := 0; i < count; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("$" + strconv.FormatInt(int64(sb.pos), 10))
			sb.pos++
		}
		b.WriteString(")")
	} else {
		b.WriteString(sb.where.ColumnA)
		b.WriteString(" ")
		b.WriteString(string(sb.where.Op))
		b.WriteString(" ")
		b.WriteString("$" + strconv.FormatInt(int64(sb.pos), 10))
		sb.pos++
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
