package Adapter

import (
	"reflect"
	"strings"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=../migrate

// GetFieldsPointers возвращает указатели на поля структуры
func GetFieldsPointers(u interface{}, args ...string) []interface{} {
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		if len(args) != 0 {
			tagsRaw := val.Type().Field(i).Tag.Get("db_ops")
			tags := strings.Split(tagsRaw, ",")
			found := false
			for _, tag := range tags {
				if tag == args[0] {
					found = true
				}
			}
			if !found {
				continue
			}
		}
		valueField := val.Field(i)
		v = append(v, valueField.Addr().Interface())
	}

	return v
}

type In struct {
	Field string
	Args  []interface{}
}

type Order struct {
	Field string
	Asc   bool
}

type LimitOffset struct {
	Offset int64
	Limit  int64
}

// Condition структура для хранения условий выборки
// Equal - условия равенства
// NotEqual - условия неравенства
// Order - условия сортировки
// LimitOffset - условия лимита и оффсета
// ForUpdate - флаг блокировки записей
// Upsert - флаг вставки записи, если не найдена
type Condition struct {
	Equal       map[string]interface{}
	NotEqual    map[string]interface{}
	Order       []*Order
	LimitOffset *LimitOffset
	ForUpdate   bool
	Upsert      bool
}
