package routes

import (
	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	UserRoutes(r)
	InvitationRoutes(r)
	InvitationThemeRoutes(r)
}
