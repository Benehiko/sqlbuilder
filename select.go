package sqlbuilder

type (
	SelectFromQuery interface {
		Where[SelectFromQuery]
		Joins
		Alias[SelectFromQuery]
		Order[SelectFromQuery]
		SQL() string
	}
	SelectQuery interface {
		Select(columns ...string) FromQuery[SelectFromQuery]
	}
)

type SelectBuilder struct {
	parent  any
	table   string
	alias   string
	columns []string
	orderBy *Sort
	joins   []*Join
	pos     int
	*WhereBuilder[SelectFromQuery]
}

// GetAlias implements queryHelper
func (s *SelectBuilder) GetAlias() string {
	return s.alias
}

// GetColumns implements queryHelper
func (s *SelectBuilder) GetColumns() []string {
	return s.columns
}

// GetJoins implements queryHelper
func (s *SelectBuilder) GetJoins() []*Join {
	return s.joins
}

// GetOrderBy implements queryHelper
func (s *SelectBuilder) GetOrderBy() *Sort {
	return s.orderBy
}

// GetParent implements queryHelper
func (s *SelectBuilder) GetParent() any {
	return s.parent
}

// GetPosition implements queryHelper
func (s *SelectBuilder) GetPosition() *int {
	return &s.pos
}

// GetReturning implements queryHelper
func (s *SelectBuilder) GetReturning() []string {
	return nil
}

// GetTable implements queryHelper
func (s *SelectBuilder) GetTable() string {
	return s.table
}

// GetWhere implements queryHelper
func (s *SelectBuilder) GetWhere() *WhereCondition {
	if s.WhereBuilder != nil {
		return s.where
	}
	return nil
}

var (
	_ SelectQuery     = (*SelectBuilder)(nil)
	_ SelectFromQuery = (*SelectBuilder)(nil)
	_ queryHelper     = (*SelectBuilder)(nil)
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
	s.alias = alias
	return s
}

func (s *SelectBuilder) Where(column string, operator Operator) WhereOptions[SelectFromQuery] {
	s.WhereBuilder = &WhereBuilder[SelectFromQuery]{
		where: &WhereCondition{
			ColumnA: column,
			Op:      operator,
		},
		parent: s,
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

func (s *SelectBuilder) SQL() string {
	return SQL(s)
}
