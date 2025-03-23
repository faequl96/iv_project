package handlers

import (
	gallery_dto "iv_project/dto/gallery"
	"iv_project/models"
)

func ConvertToGalleryResponse(gallery *models.Gallery) gallery_dto.GalleryResponse {
	return gallery_dto.GalleryResponse{
		ID:         gallery.ID,
		ImageURL1:  gallery.ImageURL1,
		ImageURL2:  gallery.ImageURL2,
		ImageURL3:  gallery.ImageURL3,
		ImageURL4:  gallery.ImageURL4,
		ImageURL5:  gallery.ImageURL5,
		ImageURL6:  gallery.ImageURL6,
		ImageURL7:  gallery.ImageURL7,
		ImageURL8:  gallery.ImageURL8,
		ImageURL9:  gallery.ImageURL9,
		ImageURL10: gallery.ImageURL10,
		ImageURL11: gallery.ImageURL11,
		ImageURL12: gallery.ImageURL12,
	}
}
