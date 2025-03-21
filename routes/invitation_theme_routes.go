package routes

import (
	"iv_project/handlers"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func InvitationThemeRoutes(r *mux.Router, jwtServices jwtToken.JWTServices) {
	invitationThemeRepository := repositories.InvitationThemeRepository(mysql.DB)
	categoryRepository := repositories.CategoryRepository(mysql.DB)
	discountCategoryRepository := repositories.DiscountCategoryRepository(mysql.DB)
	h := handlers.InvitationThemeHandler(invitationThemeRepository, categoryRepository, discountCategoryRepository)

	r.HandleFunc("/invitation-theme", middleware.Auth(jwtServices, h.CreateInvitationTheme)).Methods("POST")
	r.HandleFunc("/invitation-theme/id/{id}", middleware.Auth(jwtServices, h.GetInvitationThemeByID)).Methods("GET")
	r.HandleFunc("/invitation-themes", middleware.Auth(jwtServices, h.GetInvitationThemes)).Methods("GET")
	r.HandleFunc("/invitation-themes/category-id/{categoryId}", middleware.Auth(jwtServices, h.GetInvitationThemesByCategoryID)).Methods("GET")
	r.HandleFunc("/invitation-theme/id/{id}", middleware.Auth(jwtServices, h.UpdateInvitationThemeByID)).Methods("PATCH")
	r.HandleFunc("/invitation-theme/id/{id}", middleware.Auth(jwtServices, h.DeleteInvitationThemeByID)).Methods("DELETE")
}
