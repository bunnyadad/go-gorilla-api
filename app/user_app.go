package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	app "go-gorilla-api/app/utils"

	"go-gorilla-api/model"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) UserInitialize() {
	a.initializeUserRoutes()
}

func (a *App) initializeUserRoutes() {
	a.Router.Handle("/users", a.isAuthorized(a.getUsers)).Methods("GET")
	a.Router.Handle("/user/username:{username}", a.isAuthorized(a.getUserByUserName)).Methods("GET")
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

func (a *App) getUserByUserName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := model.User{Username: vars["username"]}
	if err := u.GetUserByUserName(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			app.RespondWithError(w, http.StatusNotFound, "User not found")
			log.Println(err.Error())
		default:
			app.RespondWithError(w, http.StatusInternalServerError, "")
			log.Println(err.Error())
		}
		return
	}
	app.RespondWithJSON(w, http.StatusOK, u)
}
