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

type UpdateBuilderWhere[T Operator] struct {
	*UpdateBuilder
}

var (
	_ UpdateSetQuery                           = (*UpdateBuilder)(nil)
	_ UpdateWhereOptionsQuery[BasicOperator]   = (*UpdateBuilderWhere[BasicOperator])(nil)
	_ UpdateWhereOptionsQuery[SpecialOperator] = (*UpdateBuilderWhere[SpecialOperator])(nil)
)

func (b *UpdateBuilder) Update(table string) UpdateSetQuery {
	return &UpdateBuilder{
		table: table,
	}
}

func (b *UpdateBuilder) Set(columns ...string) UpdateWhereQuery {
	b.set = append(b.set, columns...)
	return b
}

func (b *UpdateBuilder) Where(column string, operator BasicOperator) UpdateWhereOptionsQuery[BasicOperator] {
	b.where = &WhereCondition[BasicOperator]{
		ColumnA: column,
		Op:      operator,
	}
	return &UpdateBuilderWhere[BasicOperator]{
		UpdateBuilder: b,
	}
}

func (b *UpdateBuilder) WhereSpecial(column string, operator SpecialOperator) UpdateWhereOptionsQuery[SpecialOperator] {
	b.whereSpecial = &WhereCondition[SpecialOperator]{
		ColumnA: column,
		Op:      operator,
	}
	return &UpdateBuilderWhere[SpecialOperator]{
		UpdateBuilder: b,
	}
}

func (b *UpdateBuilderWhere[T]) And(column string, operator T) UpdateWhereOptionsQuery[T] {
	switch o := any(operator).(type) {
	case BasicOperator:
		if b.where != nil {
			tail := b.where
			for tail.next != nil {
				tail = tail.next
			}
			tail.nextOp = And
			tail.next = &WhereCondition[BasicOperator]{
				ColumnA: column,
				Op:      o,
			}
		}
	case SpecialOperator:
		if b.whereSpecial != nil {
			tail := b.whereSpecial
			tail.nextOp = And
			tail.next = &WhereCondition[SpecialOperator]{
				ColumnA: column,
				Op:      o,
			}
		}
	}
	return b
}

func (b *UpdateBuilderWhere[T]) Or(column string, operator T) UpdateWhereOptionsQuery[T] {
	switch o := any(operator).(type) {
	case BasicOperator:
		if b.where != nil {
			tail := b.where
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
		if b.whereSpecial != nil {
			tail := b.whereSpecial
			tail.nextOp = Or
			tail.next = &WhereCondition[SpecialOperator]{
				ColumnA: column,
				Op:      o,
			}
		}
	}
	return b
}

func (sb *UpdateBuilder) writeWhereSpecialHelper(b *strings.Builder, next *WhereCondition[SpecialOperator]) {
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

func (sb *UpdateBuilder) writeWhereHelper(b *strings.Builder, next *WhereCondition[BasicOperator]) {
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

func (sb *UpdateBuilder) writeWhere(b *strings.Builder) {
	b.WriteString(" WHERE ")

	if sb.whereSpecial != nil {
		sb.writeWhereSpecialHelper(b, sb.whereSpecial)
	} else {
		sb.writeWhereHelper(b, sb.where)
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
