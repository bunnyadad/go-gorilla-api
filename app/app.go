package app

import (
	"log"
	"net/http"
	"os"

	"go-gorilla-api/db"

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
}

func (a *App) Run(addr string) {
	log.Printf("Server listening on port: %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
