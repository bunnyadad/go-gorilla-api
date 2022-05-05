package app

import (
	"net/http"
	"strconv"

	app "go-gorilla-api/app/utils"

	"go-gorilla-api/model"

	_ "github.com/lib/pq"
)

func (a *App) UserInitialize() {
	a.initializeUserRoutes()
}

func (a *App) initializeUserRoutes() {
	a.Router.Handle("/users", a.isAuthorized(a.getUsers)).Methods("GET")
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := model.GetUsers(d.Database, start, count)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, "")
		log.Println(err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, users)
}
