package schema

import (
	"github.com/CenGinHong/tiny-ORM/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string
}

type Schema struct {
	Model      interface{}       // 被映射对象
	Name       string            // 表名
	Fields     []*Field          // 字段
	FieldNames []string          // 所有字段名
	fieldMap   map[string]*Field // FieldNames 和 Fields,省去遍历
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// dest是结构体指针，这里是使用反射构造指针所指向类型的值
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	// 遍历结构体的所有字段
	for i := 0; i < modelType.NumField(); i++ {
		// 获取字段
		p := modelType.Field(i)
		// 该字段需要为非embed和是导出字段的
		if !p.Anonymous && ast.IsExported(p.Name) {
			// 构建field
			field := &Field{
				Name: p.Name,
				// 不同数据库的类型不同
				Type: d.DataTypeof(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("tiny-ORM"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
