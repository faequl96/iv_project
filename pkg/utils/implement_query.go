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

	if query.Page != nil && query.Limit != nil {
		offset := (*query.Page - 1) * *query.Limit
		db = db.Offset(offset).Limit(*query.Limit)
	}

	if len(query.FilterGroups) == 0 {
		return db
	}

	db = autoJoinRelations(db, query)

	for _, group := range query.FilterGroups {
		var conditions []string
		var values []any

		for _, filter := range group.Filters {
			sqlOperator := filter.Operator.ToSQL()
			conditions = append(conditions, fmt.Sprintf("%s %s ?", filter.Field, sqlOperator))
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

	return db
}

func buildFilterValue(operator query_dto.OperatorType, value string) interface{} {
	if operator == query_dto.Like {
		return "%" + value + "%"
	}
	return value
}

func autoJoinRelations(db *gorm.DB, query *query_dto.QueryRequest) *gorm.DB {
	if query == nil {
		return db
	}

	joined := make(map[string]bool)

	for _, group := range query.FilterGroups {
		for _, filter := range group.Filters {
			parts := strings.Split(filter.Field, ".")
			if len(parts) > 1 {
				rel := parts[0]
				if !joined[rel] {
					switch rel {
					case "user_profiles":
						db = db.Joins("LEFT JOIN user_profiles ON users.id = user_profiles.user_id")
					case "invitation_datas":
						db = db.Joins("LEFT JOIN invitation_datas ON invitations.id = invitation_datas.invitation_id")
						// add more as needed
					}
					joined[rel] = true
				}
			}
		}
	}

	return db
}
