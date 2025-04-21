package query_dto

type QueryRequest struct {
	Page         *int                 `json:"page"`
	Limit        *int                 `json:"limit"`
	Sort         *SortRequest         `json:"sort"`
	FilterGroups []FilterGroupRequest `json:"filter_groups"`
}

type SortType string

const (
	SortASC  SortType = "asc"
	SortDESC SortType = "desc"
)

func (s SortType) String() string {
	maps := map[SortType]string{
		SortASC:  "asc",
		SortDESC: "desc",
	}
	return maps[s]
}

type SortRequest struct {
	Key  string   `json:"key"`
	Type SortType `json:"type"`
}

type JoinType string

const (
	JoinAND JoinType = "and"
	JoinOR  JoinType = "or"
)

func (j JoinType) String() string {
	maps := map[JoinType]string{
		JoinAND: "and",
		JoinOR:  "or",
	}
	return maps[j]
}

type FilterGroupRequest struct {
	JoinType JoinType        `json:"join_type"`
	Filters  []FilterRequest `json:"filters"`
}

type OperatorType string

const (
	Equals             OperatorType = "equals"
	Like               OperatorType = "like"
	GreaterThan        OperatorType = "greater_than"
	LessThan           OperatorType = "less_than"
	GreaterThanOrEqual OperatorType = "greater_than_or_equal"
	LessThanOrEqual    OperatorType = "less_than_or_equal"
)

func (o OperatorType) ToSQL() string {
	switch o {
	case Equals:
		return "="
	case Like:
		return "LIKE"
	case GreaterThan:
		return ">"
	case LessThan:
		return "<"
	case GreaterThanOrEqual:
		return ">="
	case LessThanOrEqual:
		return "<="
	default:
		return "="
	}
}

type FilterRequest struct {
	Field    string       `json:"field"`
	Operator OperatorType `json:"operator"`
	Value    string       `json:"value"`
}
