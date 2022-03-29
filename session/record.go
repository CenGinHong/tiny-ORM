package session

import (
	"TinyORM/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		// 获取table，这里是为了书写方便，其实只用取一次
		table := s.Model(value).RefTable()
		// 这是插入的表名和字段内容,，其实也只用写一次
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// 将插入value追加保存，这里的value是一个结构体变量
		recordValues = append(recordValues, table.RecordValues(value))
	}
	// set value
	s.clause.Set(clause.VALUES, recordValues...)
	// 构建sql
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find 实际是对SELECT的进一步包装
func (s *Session) Find(values interface{}) error {
	// values传入的是数组
	// 获取该数组的反射值
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	// Elem returns a type's element type. It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
	destType := destSlice.Type().Elem()
	// 获取该数组元素的结构体
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	// 构造select，sql
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	// 构造sql
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	// 使用原生db进行查询
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		// 获取dest的value
		dest := reflect.New(destType).Elem()
		// 将该结构体的所有字段的地址获取
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 调用原生的scan将行数据set进字段
		if err = rows.Scan(values...); err != nil {
			return err
		}
		// 将dest追加到数组，并且将数组值重新设
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
