package cqlutil

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateInsertQuery(table string, data interface{}) (string, []interface{}, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("GenerateInsertQuery: expected struct, got %s", v.Kind())
	}

	t := v.Type()
	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "-" || tag == "" {
			continue
		}
		column := strings.Split(tag, ",")[0]
		if column == "" {
			continue
		}
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
		values = append(values, v.Field(i).Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	return query, values, nil
}
