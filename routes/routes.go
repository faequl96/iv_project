package routes

import (
	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	UserRoutes(r)
	IVCoinRoutes(r)
	InvitationThemeRoutes(r)
	InvitationRoutes(r)
	InvitationDataRoutes(r)
	GalleryRoutes(r)
	ReviewRoutes(r)
}
