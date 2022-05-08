package app

import (
	"log"
	"net/http"
	"os"

	app "go-gorilla-api/app/utils"
	"go-gorilla-api/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var d db.DB
var mySigningKey []byte

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	db_user := viper.GetString("APP_DB_USERNAME")
	db_pass := viper.GetString("APP_DB_PASSWORD")
	db_host := viper.GetString("APP_DB_HOST")
	db_name := viper.GetString("APP_DB_NAME")
	mySigningKey = []byte(viper.GetString("JWT_KEY"))
	if os.Getenv("ENV") == "prod" {
		db_user = os.Getenv("PROD_DB_USERNAME")
		db_pass = os.Getenv("PROD_DB_PASSWORD")
		db_host = os.Getenv("PROD_DB_HOST")
		db_name = os.Getenv("PROD_DB_NAME")
		mySigningKey = []byte(os.Getenv("JWT_KEY"))
	}

	d.Initialize(db_user, db_pass, db_host, db_name)

	a.Router = mux.NewRouter()
	a.initializeUserRoutes()
}

func (a *App) Run(addr string) {
	CSRF := csrf.Protect(
		[]byte("251e70cd5d1a994c51fd316f7040f13d"),
		// instruct the browser to never send cookies during cross site requests
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Secure(os.Getenv("ENV") == "prod"),
	)
	log.Printf("Server listening on port: %s", addr)
	if os.Getenv("ENV") == "prod" {
		log.Fatal(http.ListenAndServeTLS(addr, "tls.crt", "tls.key", CSRF(a.Router)))
	} else {
		log.Fatal(http.ListenAndServe(addr, CSRF(a.Router)))
	}
}

func (a *App) isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			if len(r.Header["Token"][0]) < 1 {
				app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			} else {
				token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						app.RespondWithError(w, http.StatusInternalServerError, "There was error with signing the token.")
					}
					return mySigningKey, nil
				})

				if err != nil {
					app.RespondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
				if token.Valid {
					endpoint(w, r)
				}
			}
		} else {
			app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		}
	})
}
