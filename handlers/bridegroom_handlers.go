package handlers

import (
	bridegroom_dto "iv_project/dto/bridegroom"
	"iv_project/models"
)

func ConvertToBridegroomResponse(bridegroom models.Bridegroom) bridegroom_dto.BridegroomResponse {
	return bridegroom_dto.BridegroomResponse{
		ID:          bridegroom.ID,
		Nickname:    bridegroom.Nickname,
		FullName:    bridegroom.FullName,
		Title:       bridegroom.Title,
		ImageURL:    bridegroom.ImageURL,
		FatherName:  bridegroom.FatherName,
		FatherTitle: bridegroom.FatherTitle,
		MotherName:  bridegroom.MotherName,
		MotherTitle: bridegroom.MotherTitle,
	}
}
