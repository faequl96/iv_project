package bridegroom_dto

type CreateBridegroomRequest struct {
	Nickname    string `json:"nickname" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
	Title       string `json:"title" validate:"required"`
	FatherName  string `json:"father_name" validate:"required"`
	FatherTitle string `json:"father_title" validate:"required"`
	MotherName  string `json:"mother_name" validate:"required"`
	MotherTitle string `json:"mother_title" validate:"required"`
}

type UpdateBridegroomRequest struct {
	Nickname    string `json:"nickname"`
	FullName    string `json:"full_name"`
	Title       string `json:"title"`
	FatherName  string `json:"father_name"`
	FatherTitle string `json:"father_title"`
	MotherName  string `json:"mother_name"`
	MotherTitle string `json:"mother_title"`
}
