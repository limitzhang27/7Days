package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[LIMIT] = _limit
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindVars(num int) string {
	vars := make([]string, num)
	for i := 0; i < num; i++ {
		vars[i] = "?"
	}
	return strings.Join(vars, ",")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $table ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), []interface{}{}
}

// _values (tableName ["zhang", 2], ["wang", 3]) => (VALUES (?, ?), (?, ?),  ["zhang", 2, "wang", 3])
func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2)...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// _select (table, ["name", "age"])
func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $field FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []interface{}{}
}

// _where (tableName, "`name` = 'zhang'")
func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _orderby(values ...interface{}) (string, []interface{}) {
	// ORDER BY $field
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}

	for k, v := range m {
		keys = append(keys, k+" = ? ")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"COUNT(*)"})
}
