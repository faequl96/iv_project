package utils

import (
	"fmt"
	query_dto "iv_project/dto/query"
	"strings"

	"gorm.io/gorm"
)

func ImplementQuery(db *gorm.DB, query *query_dto.QueryRequest) *gorm.DB {
	if query == nil {
		return db
	}

	if query.Sort != nil {
		order := query.Sort.Key
		switch strings.ToLower(query.Sort.Type.String()) {
		case "desc":
			order += " DESC"
		default:
			order += " ASC"
		}
		db = db.Order(order)
	}

	for _, group := range query.FilterGroups {
		var conditions []string
		var values []any

		for _, filter := range group.Filters {
			sqlOperator := buildSQLOperator(filter.Operator)
			conditions = append(conditions, fmt.Sprintf("%s %s ?", group.Key, sqlOperator))
			values = append(values, buildFilterValue(filter.Operator, filter.Value))
		}

		if len(conditions) > 0 {
			join := "AND"
			if strings.ToLower(group.JoinType.String()) == "or" {
				join = "OR"
			}
			db = db.Where("("+strings.Join(conditions, " "+join+" ")+")", values...)
		}
	}

	if query.Page != nil && query.Limit != nil {
		offset := (*query.Page - 1) * *query.Limit
		db = db.Offset(offset).Limit(*query.Limit)
	}

	return db
}

func buildSQLOperator(op string) string {
	switch op {
	case "equals":
		return "="
	case "like":
		return "LIKE"
	case "greater_than":
		return ">"
	case "less_than":
		return "<"
	case "greater_than_or_equal":
		return ">="
	case "less_than_or_equal":
		return "<="
	default:
		return "="
	}
}

func buildFilterValue(operator string, value string) any {
	if operator == "like" {
		return "%" + value + "%"
	}
	return value
}
