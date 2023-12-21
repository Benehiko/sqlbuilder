package sqlbuilder

type (
	WhereOptions[T any] interface {
		And(column string, operator Operator) WhereOptions[T]
		Or(column string, operator Operator) WhereOptions[T]
		Order[T]
		Parent() T
		SQL() string
	}
	Where[T any] interface {
		Where(column string, operator Operator) WhereOptions[T]
	}
	LogicalOperator string
	WhereCondition  struct {
		ColumnA string
		Op      Operator
		ColumnB string
		nextOp  LogicalOperator
		next    *WhereCondition
	}
)

const (
	And LogicalOperator = "AND"
	Or  LogicalOperator = "OR"
)

type WhereBuilder[T any] struct {
	parent T
	where  *WhereCondition
}

var (
	_ Where[any]        = (*WhereBuilder[any])(nil)
	_ WhereOptions[any] = (*WhereBuilder[any])(nil)
)

func (w *WhereBuilder[T]) Where(column string, operator Operator) WhereOptions[T] {
	w.where = &WhereCondition{
		ColumnA: column,
		Op:      operator,
	}
	return w
}

func (w *WhereBuilder[T]) And(column string, operator Operator) WhereOptions[T] {
	WhereOperator(w.where, operator, column, And)
	return w
}

func (w *WhereBuilder[T]) Or(column string, operator Operator) WhereOptions[T] {
	WhereOperator(w.where, operator, column, Or)
	return w
}

func (w *WhereBuilder[T]) Parent() T {
	return w.parent
}

func (w *WhereBuilder[T]) OrderBy(orderBy OrderBy, columns ...string) T {
	switch o := any(w.parent).(type) {
	case Order[T]:
		return o.OrderBy(orderBy, columns...)
	}
	return w.parent
}

func (w *WhereBuilder[T]) SQL() string {
	return SQL(w.parent)
}
