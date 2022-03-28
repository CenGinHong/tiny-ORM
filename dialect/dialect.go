package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeof 映射数据表和go类型
	DataTypeof(typ reflect.Value) string
	// TableExistSQL 查找该表是否存在
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect 注册dialect
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
