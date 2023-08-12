package main

import (
	"strconv"
	"strings"
	"time"
)

type SelectBuilder struct {
	parent Query

	table   string
	as      string
	columns []string
	where   *WhereCondition
	orderBy *Sort
	joins   []*Join
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

func (s *SelectBuilder) Where(column string, value any, operator Operator) SelectFromQuery {
	s.where = &WhereCondition{
		Column: column,
		Op:     operator,
		Value:  value,
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
	b.WriteString(sb.where.Column)
	b.WriteString(" ")
	b.WriteString(string(sb.where.Op))
	b.WriteString(" ")
	switch v := sb.where.Value.(type) {
	case string:
		b.WriteString(v)
	case int:
		b.WriteString(strconv.Itoa(v))
	case float64:
		b.WriteString(strconv.FormatFloat(v, 'E', -1, 64))
	case bool:
		b.WriteString(strconv.FormatBool(v))
	case time.Time:
		b.WriteString(v.Format(time.RFC3339))
	case nil:
		b.WriteString("NULL")
	case []byte:
		b.WriteString(string(v))
	case float32:
		b.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	}
}

func (s *SelectBuilder) SQL() string {
	var sb strings.Builder

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

	if s.as != "" {
		sb.WriteString(" AS ")
		sb.WriteString(s.as)
	}

	sb.WriteString(" FROM ")
	sb.WriteString(s.table)

	if len(s.joins) > 0 {
		for _, j := range s.joins {
			sb.WriteString(" ")
			sb.WriteString(string(j.join))
			sb.WriteString(" ")
			sb.WriteString(j.table)
			sb.WriteString(" ON ")
			sb.WriteString(j.on.Column)
			sb.WriteString(" ")
			sb.WriteString(string(j.on.Op))
			sb.WriteString(" ")
			sb.WriteString(string(j.on.Value.(string)))
		}
	}

	if s.where != nil {
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
