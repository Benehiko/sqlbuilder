package main

type (
	FromQuery[T any] interface {
		From(from string) T
	}

	Alias[T any] interface {
		As(alias string) T
	}

	Where[T any] interface {
		Where(column string, value any, operator Operator) T
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
		Alias[Joins]
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
		Order[SelectFromQuery]
		SQL() string
	}

	SelectQuery interface {
		Select(columns ...string) FromQuery[SelectFromQuery]
	}

	InsertIntoQuery interface {
		Values(values ...any) InsertIntoQuery
		SelectQuery
		SQL() string
	}

	InsertQuery interface {
		Insert(columns ...string) Into[InsertIntoQuery]
	}

	Query interface {
		SelectQuery
		InsertQuery
	}
)

type Operator string

const (
	Equals             Operator = "="
	NotEqual           Operator = "!="
	GreaterThan        Operator = ">"
	GreaterThanOrEqual Operator = ">="
	LessThan           Operator = "<"
	LessThanOrEqual    Operator = "<="
	Like               Operator = "LIKE"
	NotLike            Operator = "NOT LIKE"
	In                 Operator = "IN"
	NotIn              Operator = "NOT IN"
	IsNull             Operator = "IS NULL"
	IsNotNull          Operator = "IS NOT NULL"
	IsTrue             Operator = "IS TRUE"
	IsNotTrue          Operator = "IS NOT TRUE"
	IsFalse            Operator = "IS FALSE"
	IsNotFalse         Operator = "IS NOT FALSE"
	IsNotDistinctFrom  Operator = "IS NOT DISTINCT FROM"
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

type WhereCondition struct {
	Column string
	Op     Operator
	Value  any
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
