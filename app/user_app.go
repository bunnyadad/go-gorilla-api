package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	app "go-gorilla-api/app/utils"

	"go-gorilla-api/model"

	"github.com/dgrijalva/jwt-go"
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
	a.Router.HandleFunc("/user/login", a.loginUser).Methods("POST")
	a.Router.Handle("/user/id:{id}", a.isAuthorized(a.deleteUser)).Methods("DELETE")
	a.Router.Handle("/user/id:{id}", a.isAuthorized(a.updateUser)).Methods("PUT")
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

func (a *App) loginUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		log.Println(err.Error())
		return
	}

	defer r.Body.Close()
	if err := u.GetUserByUserNameAndPassword(d.Database); err != nil {
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
	validToken, err := GenerateJWT()
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Add("Token", validToken)
	app.RespondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, "Invalid request")
		log.Println(err.Error())
	}

	u := model.User{ID: id}
	if err := u.DeleteUser(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, "")
		log.Println(err.Error())
		return
	}
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "User deleted"})
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, "Invalid request")
		log.Println(err.Error())
	}

	var u model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		log.Println(err.Error())
		return
	}

	defer r.Body.Close()
	u.ID = id

	if err := u.UpdateUser(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, "")
		log.Println(err.Error())
		return
	}
	app.RespondWithJSON(w, http.StatusOK, u)
}

// Helper function
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
