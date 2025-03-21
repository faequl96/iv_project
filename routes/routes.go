package routes

import (
	jwtToken "iv_project/pkg/jwt"
	"os"

	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	j := jwtToken.JWTService(os.Getenv("JWT_SECRET"), "iv_project")

	AuthRoutes(r, j)
	UserRoutes(r, j)
	UserProfileRoutes(r, j)
	IVCoinRoutes(r, j)
	AdMobRoutes(r, j)
	CategoryRoutes(r)
	DiscountCategoryRoutes(r)
	IVCoinPackageRoutes(r)
	InvitationThemeRoutes(r)
	ReviewRoutes(r)
	InvitationRoutes(r)
	InvitationDataRoutes(r)
	GalleryRoutes(r)
	TransactionRoutes(r)
}
