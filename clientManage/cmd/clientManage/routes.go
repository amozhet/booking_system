package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/update/:id", app.requirePermission("user:write", app.updateUserHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete/:id", app.requirePermission("user:write", app.deleteUserHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users", app.requirePermission("user:write", app.getAllUsersHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/email", app.requirePermission("user:write", app.getUserByEmailHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
