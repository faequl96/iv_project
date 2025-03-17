package review_dto

type ReviewResponse struct {
	ID      uint   `json:"id"`
	Star    int    `json:"star"`
	Comment string `json:"comment"`
}
