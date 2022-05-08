package main

import (
	"log"
	"os"

	"go-gorilla-api/app"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	log.Println("ENV: " + current_env)

	a := app.App{}

	a.Initialize()
	if os.Getenv("PORT") == "" {
		a.Run(":" + viper.GetString("PORT"))
	} else {
		a.Run(":" + os.Getenv("PORT"))
	}
}
