package routes

import (
	"iv_project/handlers"
	"iv_project/pkg/mysql"
	"iv_project/repositories"

	"github.com/gorilla/mux"
)

func InvitationThemeRoutes(r *mux.Router) {
	invitationThemeRepository := repositories.InvitationThemeRepository(mysql.DB)
	categoryRepository := repositories.CategoryRepository(mysql.DB)
	discountCategoryRepository := repositories.DiscountCategoryRepository(mysql.DB)
	h := handlers.InvitationThemeHandler(invitationThemeRepository, categoryRepository, discountCategoryRepository)

	r.HandleFunc("/invitation-theme", h.CreateInvitationTheme).Methods("POST")
	r.HandleFunc("/invitation-theme/id/{id}", h.GetInvitationThemeByID).Methods("GET")
	r.HandleFunc("/invitation-themes", h.GetInvitationThemes).Methods("GET")
	r.HandleFunc("/invitation-themes/category/{category}", h.GetInvitationThemesByCategory).Methods("GET")
	r.HandleFunc("/invitation-theme/id/{id}", h.UpdateInvitationTheme).Methods("PATCH")
	r.HandleFunc("/invitation-theme/id/{id}", h.DeleteInvitationTheme).Methods("DELETE")
}
