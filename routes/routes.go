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
	CategoryRoutes(r, j)
	DiscountCategoryRoutes(r, j)
	IVCoinPackageRoutes(r, j)
	InvitationThemeRoutes(r, j)
	DiscountRoutes(r, j)
	ReviewRoutes(r)
	InvitationRoutes(r)
	InvitationDataRoutes(r)
	GalleryRoutes(r)
	TransactionRoutes(r)
}
