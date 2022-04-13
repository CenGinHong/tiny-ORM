package session

import (
	"errors"
	"github.com/CenGinHong/TinyORM/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
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
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

// Find 实际是对SELECT的进一步包装
func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
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
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		// 将dest追加到数组，并且将数组值重新设
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// Update 接受两种入参，平铺开的键值对和map
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	// 尝试强转成map
	m, ok := kv[0].(map[string]interface{})
	// 将字段聚集成map
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	// 构建sql
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	// value是指针,获取值类型的cvalue
	dest := reflect.Indirect(reflect.ValueOf(value))
	// 构造数组
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	// 只用find方式，实际只查出一个元素
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	// 只取出一个
	dest.Set(destSlice.Index(0))
	return nil
}
