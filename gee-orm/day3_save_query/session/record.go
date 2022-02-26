package session

import (
	"geeorm/clause"
	"reflect"
)

// Insert 实现插入功能，返回插入行数和错误
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable() // 获取表明
		// 多字调用
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)

	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find 查询
func (s *Session) Find(values interface{}) error {
	// 通过反射得到入参的反射对象
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem() // reflect.Indirect(destSlice.Type())
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem() // 根据反射类型构造一个新的对象
		var values []interface{}
		for _, name := range table.FieldNames {
			// 获取新对象全部字段的指针
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 全部字段指针赋值
		if err := rows.Scan(values...); err != nil {
			return err
		}

		// 将赋值后的指针重新设置为入参的反射对象中
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()
}
