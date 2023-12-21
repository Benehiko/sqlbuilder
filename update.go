package sqlbuilder

type (
	UpdateSetQuery interface {
		Set(columns ...string) UpdateWhereQuery
	}
	UpdateWhereQuery interface {
		Where[UpdateReturningQuery]
		SQL() string
	}
	UpdateReturningQuery interface {
		Returning(columns ...string) interface {
			SQL() string
		}
	}
	UpdateQuery interface {
		Update(table string) UpdateSetQuery
	}
	UpdateBuilder struct {
		table     string
		columns   []string
		returning []string
		pos       int
		*WhereBuilder[UpdateReturningQuery]
	}
)

var (
	_ UpdateQuery    = (*UpdateBuilder)(nil)
	_ UpdateSetQuery = (*UpdateBuilder)(nil)
	_ queryHelper    = (*UpdateBuilder)(nil)
)

func (b *UpdateBuilder) Update(table string) UpdateSetQuery {
	return &UpdateBuilder{
		table: table,
	}
}

func (b *UpdateBuilder) Set(columns ...string) UpdateWhereQuery {
	b.columns = append(b.columns, columns...)
	return b
}

func (b *UpdateBuilder) Where(column string, operator Operator) WhereOptions[UpdateReturningQuery] {
	b.WhereBuilder = &WhereBuilder[UpdateReturningQuery]{
		parent: b,
		where: &WhereCondition{
			ColumnA: column,
			Op:      operator,
		},
	}
	return b
}

func (b *UpdateBuilder) Returning(columns ...string) interface {
	SQL() string
} {
	b.returning = columns
	return b
}

func (b *UpdateBuilder) SQL() string {
	return SQL(b)
}

func (b *UpdateBuilder) GetTable() string {
	return b.table
}

func (b *UpdateBuilder) GetPosition() *int {
	return &b.pos
}

func (b *UpdateBuilder) GetWhere() *WhereCondition {
	return b.where
}

func (b *UpdateBuilder) GetColumns() []string {
	return b.columns
}

func (b *UpdateBuilder) GetReturning() []string {
	return b.returning
}

func (b *UpdateBuilder) GetAlias() string {
	return ""
}

func (b *UpdateBuilder) GetJoins() []*Join {
	return nil
}

func (b *UpdateBuilder) GetOrderBy() *Sort {
	return nil
}

func (b *UpdateBuilder) GetParent() any {
	return nil
}
