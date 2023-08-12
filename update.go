package main

import (
	"strconv"
	"strings"
	"time"
)

type UpdateBuilder struct {
	table string
	set   []UpdateSet
	where *WhereCondition
}

var _ UpdateSetQuery = (*UpdateBuilder)(nil)

func (b *UpdateBuilder) Update(table string) UpdateSetQuery {
	return &UpdateBuilder{
		table: table,
	}
}

func (b *UpdateBuilder) Set(set ...UpdateSet) UpdateWhereQuery {
	b.set = append(b.set, set...)
	return b
}

func (b *UpdateBuilder) Where(column string, value any, operator Operator) UpdateWhereQuery {
	b.where = &WhereCondition{
		Column: column,
		Op:     operator,
		Value:  value,
	}
	return b
}

func (sb *UpdateBuilder) writeWhere(b *strings.Builder) {
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

func (b *UpdateBuilder) SQL() string {
	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(b.table)
	sb.WriteString(" SET ")
	for _, set := range b.set {
		sb.WriteString(set.Column)
		sb.WriteString(" = ")
		sb.WriteString(ToString[any](set.Value))
		sb.WriteString(", ")
	}
	s := strings.TrimSuffix(sb.String(), ", ")
	sb = strings.Builder{}
	sb.WriteString(s)

	if b.where != nil {
		b.writeWhere(&sb)
	}

	return strings.TrimSpace(sb.String())
}
