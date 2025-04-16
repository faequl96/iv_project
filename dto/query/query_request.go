package query_dto

type QueryRequest struct {
	Page         *int                 `json:"page"`
	Limit        *int                 `json:"limit"`
	Sort         *SortRequest         `json:"sort"`
	FilterGroups []FilterGroupRequest `json:"filter_groups"`
}

type SortType string

const (
	SortASC  SortType = "ASC"
	SortDESC SortType = "DESC"
)

func (u SortType) String() string {
	maps := map[SortType]string{
		SortASC:  "ASC",
		SortDESC: "DESC",
	}
	return maps[u]
}

type SortRequest struct {
	Key  string   `json:"key"`
	Type SortType `json:"type"`
}

type JoinType string

const (
	JoinAND JoinType = "AND"
	JoinOR  JoinType = "OR"
)

func (u JoinType) String() string {
	maps := map[JoinType]string{
		JoinAND: "AND",
		JoinOR:  "OR",
	}
	return maps[u]
}

type FilterGroupRequest struct {
	Key      string          `json:"key"`
	JoinType JoinType        `json:"join_type"`
	Filters  []FilterRequest `json:"filters"`
}

type FilterRequest struct {
	Value    string `json:"value"`
	Operator string `json:"operator"`
}
