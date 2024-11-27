package router

import (
	"net/http"

	"github.com/nsaltun/userapi/internal/handler"
	"github.com/nsaltun/userapi/pkg/lib/middlewares/httpwrap"
)

func NewRouter(userHandler handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	// Register routes using the custom context handler
	mux.HandleFunc("POST /users", httpwrap.ContextMiddleware(userHandler.CreateUser))
	mux.HandleFunc("PUT /users/{id}", httpwrap.ContextMiddleware(userHandler.UpdateUserById))
	mux.HandleFunc("DELETE /users/{id}", httpwrap.ContextMiddleware(userHandler.DeleteUserById))
	mux.HandleFunc("POST /users/filter", httpwrap.ContextMiddleware(userHandler.ListUsers))

	return mux
}
