package sqlbuilder

type (
	BasicOperator       string
	SpecialOperator     func() (int, string)
	SpecialOperatorFunc func(count int) SpecialOperator
	Operator            interface {
		get() any
	}
)

var (
	_ Operator = Equals
	_ Operator = In(0)
	_ Operator = NotIn(0)
)

var In SpecialOperatorFunc = func(count int) SpecialOperator {
	return func() (int, string) {
		return count, "IN"
	}
}

var NotIn SpecialOperatorFunc = func(count int) SpecialOperator {
	return func() (int, string) {
		return count, "NOT IN"
	}
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

func (b BasicOperator) get() any {
	return b
}

func (s SpecialOperator) get() any {
	return s
}
