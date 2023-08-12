# SQL generator

SQL generator builds SQL from Go code.

```go
NewSelectQuery("col1", "col2").From("users").As("t").InnerJoin("roles").On("t.id", "r.id").Where("u.id", 1, Equals).SQL()
```

Produces

```sql
SELECT col1, col2 FROM users as t INNER JOIN roles ON t.id = r.id WHERE u.id = 1;
```
