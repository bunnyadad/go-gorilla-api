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
	a.Router.Handle("/user/id:{id}", a.isAuthorized(a.getUser)).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")

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

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, "Invalid request")
		log.Println(err.Error())
	}

	u := model.User{ID: id}
	if err := u.GetUser(d.Database); err != nil {
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

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		log.Println(err.Error())
		return
	}

	defer r.Body.Close()

	if err := u.CreateUser(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, "")
		log.Println(err.Error())
		return
	}
	app.RespondWithJSON(w, http.StatusCreated, u)
}
