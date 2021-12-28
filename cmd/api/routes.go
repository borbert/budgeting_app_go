package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)

	router.HandlerFunc(http.MethodGet, "/v1/budgetitem/:id", app.getOneBudgetItem)

	router.HandlerFunc(http.MethodGet, "/v1/budgetitems", app.getAllBudgetItems)
	router.HandlerFunc(http.MethodPost, "/v1/graphql", app.budgetItemsGraphQL)

	router.HandlerFunc(http.MethodGet, "/v1/alltags", app.getAllTags)

	router.HandlerFunc(http.MethodGet, "/v1/budgetitems/:tag", app.getAllBudgetItemsbyTag)

	//Get user preferences
	router.HandlerFunc(http.MethodGet, "/v1/users/preferences/:id", app.getUserPreferences)

	// Protected routes

	// router.HandlerFunc(http.MethodPost, "/v1/editbudgetitem", app.editBudgetItem)
	// edits and creates new if ID is null
	router.POST("/v1/editbudgetitem", app.wrap(secure.ThenFunc(app.editBudgetItem)))
	// edits and creates new if ID is null
	router.POST("/v1/edituser", app.wrap(secure.ThenFunc(app.editUser)))

	// router.HandlerFunc(http.MethodGet, "/v1/users", app.getAllUsers)
	router.POST("/v1/users", app.wrap(secure.ThenFunc(app.getAllUsers)))

	// router.HandlerFunc(http.MethodGet, "/v1/deletebudgetitem/:id", app.deleteBudgetItem)
	router.POST("/v1/deletebudgetitem/:id", app.wrap(secure.ThenFunc(app.deleteBudgetItem)))

	router.POST("/v1/deleteuser/:id", app.wrap(secure.ThenFunc(app.deleteUser)))

	return app.enableCORS(router)
}
