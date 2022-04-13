package session

import (
	"fmt"
	"github.com/CenGinHong/TinyORM/log"
	"github.com/CenGinHong/TinyORM/schema"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	// 设置当前session所操作的表
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.refTable
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	// 将列字段组织起来
	desc := strings.Join(columns, ",")
	if _, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, desc)).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Session) DropTable() error {
	if _, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Session) HasTable() bool {
	// 是否存在该表
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	if err := row.Scan(&tmp); err != nil {
		log.Error(err)
	}
	return tmp == s.RefTable().Name
}
