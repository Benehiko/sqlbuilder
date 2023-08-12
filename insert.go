package main

import (
	"strconv"
	"strings"
)

type InsertBuilder struct {
	table   string
	columns []string
	values  []any
	as      string
	s       SelectQuery
}

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

func (ib *InsertBuilder) SQL() string {
	var sb strings.Builder
	sb.WriteString("INSERT ")
	sb.WriteString("INTO ")
	sb.WriteString(ib.table)
	sb.WriteString(" (")
	sb.WriteString(strings.TrimSpace(strings.Join(ib.columns, ", ")))
	sb.WriteString(")")
	sb.WriteString(" VALUES ")
	sb.WriteString("(")

	cols := []string{}
	for i := 1; i <= len(ib.columns); i++ {
		cols = append(cols, "$"+strconv.Itoa(i))
	}

	sb.WriteString(strings.TrimSpace(strings.Join(cols, ", ")))
	sb.WriteString(")")

	return sb.String()
}
