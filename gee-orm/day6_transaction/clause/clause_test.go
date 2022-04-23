package clause

import (
	"reflect"
	"testing"
)

func clause_test(t *testing.T) {
	var c Clause

	c.Set(LIMIT, 3)
	c.Set(SELECT, "User", []string{"*"})
	c.Set(ORDERBY, "Age ASC")
	c.Set(WHERE, "Name = ?", "Tom")
	sql, vars := c.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)

	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build sql")
	}

	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build sql")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		clause_test(t)
	})
}
