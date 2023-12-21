package sqlbuilder

import (
	"strconv"
	"strings"
)

type (
	Into[T any] interface {
		Into(table string) T
	}
	InsertIntoQuery interface {
		Values(values ...any) InsertIntoQuery
		Returning(columns ...string) InsertIntoQuery
		SelectQuery
		SQL() string
	}
	InsertQuery interface {
		Insert(columns ...string) Into[InsertIntoQuery]
	}
	InsertBuilder struct {
		table     string
		columns   []string
		returning []string
		values    []any
		as        string
		s         SelectQuery
	}
)

var _ InsertQuery = &InsertBuilder{}

func (ib *InsertBuilder) Insert(columns ...string) Into[InsertIntoQuery] {
	ib.columns = columns
	return ib
}

func (ib *InsertBuilder) As(alias string) InsertIntoQuery {
	ib.as = alias
	return ib
}

func (ib *InsertBuilder) Into(table string) InsertIntoQuery {
	ib.table = table
	return ib
}

func (ib *InsertBuilder) Values(values ...any) InsertIntoQuery {
	ib.values = values
	return ib
}

func (ib *InsertBuilder) Select(columns ...string) FromQuery[SelectFromQuery] {
	ib.s = &SelectBuilder{
		parent: ib,
	}
	return ib.s.Select(columns...)
}

func (ib *InsertBuilder) Returning(columns ...string) InsertIntoQuery {
	ib.returning = columns
	return ib
}

func (ib *InsertBuilder) SQL() string {
	var sb strings.Builder
	sb.WriteString("INSERT ")
	sb.WriteString("INTO ")
	sb.WriteString(ib.table)

	sb.WriteString(" (")
	sb.WriteString(strings.TrimSpace(strings.Join(ib.columns, ", ")))
	sb.WriteString(")")

	if ib.s == nil {
		sb.WriteString(" VALUES ")
		sb.WriteString("(")

		cols := []string{}
		for i := 1; i <= len(ib.columns); i++ {
			cols = append(cols, "$"+strconv.Itoa(i))
		}

		sb.WriteString(strings.TrimSpace(strings.Join(cols, ", ")))
		sb.WriteString(")")
	}

	if ib.returning != nil {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.TrimSpace(strings.Join(ib.returning, ", ")))
	}

	return sb.String()
}
