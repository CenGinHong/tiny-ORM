package session

import (
	"TinyORM/clause"
	"TinyORM/dialect"
	"TinyORM/log"
	"TinyORM/schema"
	"database/sql"
	"strings"
)

type Session struct {
	db       *sql.DB // 连接数据库的指针
	dialect  dialect.Dialect
	tx       *sql.Tx // 支持事务
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder // 拼接sql
	sqlVars  []interface{}   // sql 需要填入的变量
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// DB 返回CommonDB接口，当使用事务时基于事务做查询
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec 执行sql
func (s *Session) Exec() (result sql.Result, err error) {
	// 执行后清空sql内容
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// 使用原生db执行
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
		return
	}
	return
}

// QueryRow 查询
func (s *Session) QueryRow() *sql.Row {
	// 执行后清空sql内容
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// 使用原生db执行，QueryRow只返回一行
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows 查询多行
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	// 执行后清空sql内容
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// Query返回多行
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
