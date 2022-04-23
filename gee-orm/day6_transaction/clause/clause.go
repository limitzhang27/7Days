package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	WHERE
	LIMIT
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql     map[Type]string        // 保存SQL语句
	sqlVars map[Type][]interface{} // 保存变量
}

// Set 传入类型和变量，构造查询语句
func (c *Clause) Set(t Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[t](vars...)
	c.sql[t] = sql
	c.sqlVars[t] = vars
}

// Build 将生成的子 SQL 按照排序组装成完成的SQL
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	sqls := make([]string, 0)
	sqlVars := make([]interface{}, 0)
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			sqlVars = append(sqlVars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), sqlVars
}
