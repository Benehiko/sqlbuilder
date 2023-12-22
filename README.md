# SQL generator

Builds SQL using Go.
Values are omitted and replaced with placeholder values.

```go
s := Select("u.id", "u.username").
			From("users").
			As("u").
			InnerJoin("roles").
			As("r").
			On("r.id", "u.role_id").
			LeftJoin("permissions").
			On("p.id", "u.permission_id").
			Where("u.id", Equals).
			OrderBy(Asc, "u.id").SQL()
```

Produces

```sql
SELECT u.id, u.username FROM
		users AS u INNER JOIN roles AS r ON r.id = u.role_id
		LEFT JOIN permissions ON p.id = u.permission_id
		WHERE u.id = $1 ORDER BY u.id ASC
```

```go
Select("id", "username").From("users").Where("id", In(3)).SQL()
```

```sql
SELECT id, username FROM users WHERE id IN ($1, $2, $3)
```


## Get Started

```sh
go get -u github.com/Benehiko/sqlbuilder
```

```go
package main

import "github.com/Benehiko/sqlbuilder"

func main() {
    d, err := sql.Open("...", "connection-string")
    if err != nil {
        panic(err)
    }

    // values are replaced with placeholder values
	s := sqlbuilder.Insert("user_id", "name").Into("users").Values(1, "foo").Returning("name").SQL()
    // INSERT INTO users (user_id, name) VALUES ($1, $2)
    var name string
    // add the values here in the order they were given above
    err := db.QueryRow(ctx, s, 1, "foo").Scan(&name)
    if err != nil {
        panic(err)
    }
}
```

## More Examples

### Select
```go
Select("id", "username").From("users").Where("id", Equals).And("username", Equals).Or("email", Equals).SQL()
```

### Update
```go
Update("users").Set("id").Where("id", Equals).And("username", Equals).SQL()
```

### Delete
```go
Delete().From("users").Where("id", NotEqual).And("name", In(3)).SQL()
```

### Insert
```go
Insert("user_id", "name").Into("users").Select("id", "name").From("users").SQL()
```

### Where Operators

Where's accept `BasicOperator` and `SpecialOperator` types.
See below for the list of both types.

SpecialOperators require specifying the number of values you expect to
add to an `IN` or `NOT IN` query.

For example `In(2)` specifies only two placeholder variables inside the `IN` SQL operator.

```sql
IN ($1, $2)
```

```go
s := Select("id").From("users")

s.Where("id", NotEqual)
s.And("created_at", LessThanOrEqual)
s.Or("username", In(2))
s.Or("email", NotIn(5))

/*
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
*/
```
