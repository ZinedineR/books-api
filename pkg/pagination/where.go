package pagination

import (
	"books-api/internal/model"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func Where(filter model.FilterParams, query *gorm.DB) *gorm.DB {
	for _, f := range filter {
		keylist := GenerateWhere(*f)
		switch f.Operator {
		case "like":
			f.Field = "lower(" + f.Field + ")"
			query = query.Where(fmt.Sprintf("%s %s ?", f.Field, f.Operator), keylist)
		case "in", "not in":
			query = query.Where(fmt.Sprintf("%s %s (?)", f.Field, f.Operator), keylist)
		default:
			query = query.Where(fmt.Sprintf("%s %s ?", f.Field, f.Operator), f.Value)
		}
	}
	return query
}

func GenerateWhere(filter model.FilterParam) []interface{} {
	keySearch := strings.ToLower(filter.Value)
	var keyList []interface{}

	if filter.Operator == "like" {
		keyList = make([]interface{}, 1)
		keyList[0] = "%" + keySearch + "%"
	} else if filter.Operator == "in" || filter.Operator == "not in" {
		keys := strings.Split(keySearch, ",")
		keyList = make([]interface{}, len(keys))
		for i, key := range keys {
			keyList[i] = key
		}
	} else {
		keyList = make([]interface{}, 1)
		keyList[0] = keySearch
	}
	return keyList
}
