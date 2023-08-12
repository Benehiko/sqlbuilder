package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Run("case=select", func(t *testing.T) {
		s := NewSelectQuery("id", "username").From("users").Where("colA", "colB", Equals).SQL()
		require.Equal(t, "SELECT id, username FROM users WHERE colA = colB", s)
	})

	t.Run("case=select inner join", func(t *testing.T) {
		s := NewSelectQuery("u.id", "u.username").From("users").As("u").InnerJoin("roles").On("r.id", "u.role_id").SQL()
		require.Equal(t, "SELECT u.id, u.username AS u FROM users INNER JOIN roles ON r.id = u.role_id", s)
	})

	t.Run("case=select inner join with where", func(t *testing.T) {
		s := NewSelectQuery("u.id", "u.username").From("users").As("u").InnerJoin("roles").On("r.id", "u.role_id").Where("u.id", 1, Equals).SQL()
		require.Equal(t, "SELECT u.id, u.username AS u FROM users INNER JOIN roles ON r.id = u.role_id WHERE u.id = 1", s)
	})

	t.Run("case=select inner join with where and order by", func(t *testing.T) {
		s := NewSelectQuery("u.id", "u.username").From("users").As("u").InnerJoin("roles").On("r.id", "u.role_id").Where("u.id", 1, Equals).OrderBy(Asc, "u.id").SQL()
		require.Equal(t, "SELECT u.id, u.username AS u FROM users INNER JOIN roles ON r.id = u.role_id WHERE u.id = 1 ORDER BY u.id ASC", s)
	})

	t.Run("case=select inner join with where and order by and left join", func(t *testing.T) {
		s := NewSelectQuery("u.id", "u.username").
			From("users").
			As("u").
			InnerJoin("roles").
			On("r.id", "u.role_id").
			LeftJoin("permissions").
			On("p.id", "u.permission_id").
			Where("u.id", 1, Equals).
			OrderBy(Asc, "u.id").SQL()

		expected := `SELECT u.id, u.username AS u FROM
		users INNER JOIN roles ON r.id = u.role_id
		LEFT JOIN permissions ON p.id = u.permission_id
		WHERE u.id = 1 ORDER BY u.id ASC`
		require.Equal(t, strings.ReplaceAll(expected, "\n\t\t", " "), s)
	})

	t.Run("case=order by desc", func(t *testing.T) {
		s := NewSelectQuery("u.id", "u.username").
			From("users").
			Where("u.id", 1, Equals).
			OrderBy(Desc, "u.id").SQL()

		expected := "SELECT u.id, u.username FROM users WHERE u.id = 1 ORDER BY u.id DESC"
		require.Equal(t, expected, s)
	})

	t.Run("case=insert into", func(t *testing.T) {
		s := NewInsertQuery("user_id", "name").Into("users").Values(1, "foo").SQL()
		require.Equal(t, "INSERT INTO users (user_id, name) VALUES ($1, $2)", s)
	})

	t.Run("case=insert into with alias", func(t *testing.T) {
		s := NewInsertQuery("user_id", "name").Into("users").Values(1, "foo").SQL()
		require.Equal(t, "INSERT INTO users (user_id, name) VALUES ($1, $2)", s)
	})

	t.Run("case=insert into with select", func(t *testing.T) {
		s := NewInsertQuery("user_id", "name").Into("users").Select("id", "name").From("users").SQL()
		require.Equal(t, "INSERT INTO users (user_id, name) SELECT id, name FROM users", s)
	})

	t.Run("case=update", func(t *testing.T) {
		s := NewUpdateQuery("users").Set(UpdateSet{Column: "id", Value: 1}).Where("id", 1, Equals).SQL()
		require.Equal(t, "UPDATE users SET id = 1 WHERE id = 1", s)
	})
}
