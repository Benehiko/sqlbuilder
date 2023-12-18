package sqlbuilder

import (
	"strconv"
	"strings"
)

type UpdateBuilder struct {
	table        string
	set          []string
	where        *WhereCondition[BasicOperator]
	whereSpecial *WhereCondition[SpecialOperator]
	pos          int
}

var _ UpdateSetQuery = (*UpdateBuilder)(nil)

func (b *UpdateBuilder) Update(table string) UpdateSetQuery {
	return &UpdateBuilder{
		table: table,
	}
}

func (b *UpdateBuilder) Set(columns ...string) UpdateWhereQuery {
	b.set = append(b.set, columns...)
	return b
}

func (b *UpdateBuilder) Where(column string, operator BasicOperator) UpdateWhereQuery {
	b.where = &WhereCondition[BasicOperator]{
		ColumnA: column,
		Op:      operator,
	}
	return b
}

func (sb *UpdateBuilder) writeWhere(b *strings.Builder) {
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
			sb.pos += 1
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

func (b *UpdateBuilder) SQL() string {
	if b.pos == 0 {
		b.pos = 1
	}

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(b.table)
	sb.WriteString(" SET ")
	for _, col := range b.set {
		sb.WriteString(col)
		sb.WriteString(" = ")
		sb.WriteString("$" + strconv.Itoa(b.pos))
		sb.WriteString(", ")
		b.pos++
	}
	s := strings.TrimSuffix(sb.String(), ", ")
	sb = strings.Builder{}
	sb.WriteString(s)

	if b.where != nil || b.whereSpecial != nil {
		b.writeWhere(&sb)
	}

	return strings.TrimSpace(sb.String())
}
