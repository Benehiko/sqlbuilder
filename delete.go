package sqlbuilder

type (
	DeleteFromQuery interface {
		SQL() string
		Where[DeleteFromQuery]
	}
	DeleteQuery interface {
		Delete() FromQuery[DeleteFromQuery]
	}
	DeleteBuilder struct {
		table   string
		alias   string
		joins   []*Join
		columns []string
		orderBy *Sort
		pos     int
		*WhereBuilder[DeleteFromQuery]
	}
)

var (
	_ DeleteQuery     = (*DeleteBuilder)(nil)
	_ DeleteFromQuery = (*DeleteBuilder)(nil)
	_ queryHelper     = (*DeleteBuilder)(nil)
)

// Delete implementd.DeleteQuery
func (d *DeleteBuilder) Delete() FromQuery[DeleteFromQuery] {
	return d
}

func (d *DeleteBuilder) From(table string) DeleteFromQuery {
	d.table = table
	return d
}

func (d *DeleteBuilder) Where(column string, operator Operator) WhereOptions[DeleteFromQuery] {
	d.WhereBuilder = &WhereBuilder[DeleteFromQuery]{
		parent: d,
		where: &WhereCondition{
			ColumnA: column,
			Op:      operator,
		},
	}
	return d
}

// d.L implements DeleteWhereOptionsQuery
func (d *DeleteBuilder) Table() string {
	return d.table
}

// GetAliad.implements queryHelper
func (d *DeleteBuilder) GetAlias() string {
	return d.alias
}

// GetColumnd.implements queryHelper
func (d *DeleteBuilder) GetColumns() []string {
	return d.columns
}

// GetJoind.implements queryHelper
func (d *DeleteBuilder) GetJoins() []*Join {
	return d.joins
}

// GetOrderBy implementd.queryHelper
func (d *DeleteBuilder) GetOrderBy() *Sort {
	return d.orderBy
}

// GetParent implementd.queryHelper
func (d *DeleteBuilder) GetParent() any {
	return d.parent
}

// GetPod.tion implements queryHelper
func (d *DeleteBuilder) GetPosition() *int {
	return &d.pos
}

// GetReturning implementd.queryHelper
func (d *DeleteBuilder) GetReturning() []string {
	return nil
}

// GetTable implementd.queryHelper
func (d *DeleteBuilder) GetTable() string {
	return d.table
}

// GetWhere implementd.queryHelper
func (d *DeleteBuilder) GetWhere() *WhereCondition {
	return d.where
}
