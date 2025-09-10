package bridegroom_dto

type BridegroomResponse struct {
	ID          uint   `json:"id"`
	Nickname    string `json:"nickname"`
	FullName    string `json:"full_name"`
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	FatherName  string `json:"father_name"`
	FatherTitle string `json:"father_title"`
	MotherName  string `json:"mother_name"`
	MotherTitle string `json:"mother_title"`
}
