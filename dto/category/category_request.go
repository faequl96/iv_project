package category_dto

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}
