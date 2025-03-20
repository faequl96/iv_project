package routes

import (
	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	AuthRoutes(r)
	UserRoutes(r)
	UserProfileRoutes(r)
	IVCoinRoutes(r)
	IVCoinPackageRoutes(r)
	InvitationThemeRoutes(r)
	CategoryRoutes(r)
	ReviewRoutes(r)
	InvitationRoutes(r)
	InvitationDataRoutes(r)
	GalleryRoutes(r)
	TransactionRoutes(r)
}
