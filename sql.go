package sqlbuilder

type (
	FromQuery[T any] interface {
		From(from string) T
	}

	Alias[T any] interface {
		As(alias string) T
	}

	Where[T any] interface {
		Where(column string, operator BasicOperator) T
	}

	WhereSpecial[T any] interface {
		WhereSpecial(column string, operator SpecialOperator) T
	}

	Order[T any] interface {
		OrderBy(orderBy OrderBy, columns ...string) T
	}

	JoinOn interface {
		On(column string, joinColumn string) SelectFromQuery
	}

	Into[T any] interface {
		Into(table string) T
	}

	AliasOrJoinOn interface {
		Alias[JoinOn]
		JoinOn
	}

	Joins interface {
		InnerJoin(table string) AliasOrJoinOn
		LeftJoin(table string) AliasOrJoinOn
		RightJoin(table string) AliasOrJoinOn
		FullOuterJoin(table string) AliasOrJoinOn
	}

	SelectFromQuery interface {
		Joins
		Alias[SelectFromQuery]
		Where[SelectFromQuery]
		WhereSpecial[SelectFromQuery]
		Order[SelectFromQuery]
		SQL() string
	}

	SelectQuery interface {
		Select(columns ...string) FromQuery[SelectFromQuery]
	}

	InsertIntoQuery interface {
		Values(values ...any) InsertIntoQuery
		SelectQuery
		Returning(columns ...string) InsertIntoQuery
		SQL() string
	}

	InsertQuery interface {
		Insert(columns ...string) Into[InsertIntoQuery]
	}

	UpdateSetQuery interface {
		Set(columns ...string) UpdateWhereQuery
	}

	UpdateWhereQuery interface {
		Where[UpdateWhereQuery]

		SQL() string
	}

	UpdateQuery interface {
		Update(table string) UpdateQuery
	}

	Query interface {
		SelectQuery
		InsertQuery
		UpdateQuery
	}
)

type BasicOperator string

type SpecialOperator func() (int, string)

func In(count int) SpecialOperator {
	return func() (int, string) {
		return count, "IN"
	}
}

func NotIn(count int) SpecialOperator {
	return func() (int, string) {
		return count, "NOT IN"
	}
}

type Operator interface {
	BasicOperator | SpecialOperator
}

const (
	Equals             BasicOperator = "="
	NotEqual           BasicOperator = "!="
	GreaterThan        BasicOperator = ">"
	GreaterThanOrEqual BasicOperator = ">="
	LessThan           BasicOperator = "<"
	LessThanOrEqual    BasicOperator = "<="
	Like               BasicOperator = "LIKE"
	NotLike            BasicOperator = "NOT LIKE"
	IsNull             BasicOperator = "IS NULL"
	IsNotNull          BasicOperator = "IS NOT NULL"
	IsTrue             BasicOperator = "IS TRUE"
	IsNotTrue          BasicOperator = "IS NOT TRUE"
	IsFalse            BasicOperator = "IS FALSE"
	IsNotFalse         BasicOperator = "IS NOT FALSE"
	IsNotDistinctFrom  BasicOperator = "IS NOT DISTINCT FROM"
)

type OrderBy string

const (
	Asc  OrderBy = "ASC"
	Desc OrderBy = "DESC"
)

type Sort struct {
	columns []string
	orderBy OrderBy
}

type WhereCondition[T Operator] struct {
	ColumnA string
	Op      T
	ColumnB string
}

func NewSelectQuery(columns ...string) FromQuery[SelectFromQuery] {
	sb := &SelectBuilder{}
	sb.Select(columns...)
	return sb
}

func NewInsertQuery(columns ...string) Into[InsertIntoQuery] {
	ib := &InsertBuilder{}
	ib.Insert(columns...)
	return ib
}

func NewUpdateQuery(table string) UpdateSetQuery {
	ub := &UpdateBuilder{}
	return ub.Update(table)
}
