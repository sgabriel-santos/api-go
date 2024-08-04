package main

import (
	"api/src/config"
	"api/src/middlewares"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.LoadEnvironmentVariables()
	r := router.GenerateRoutes()
	middlewares.ConfigureLogger()

	log.Printf("Escutando na porta %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
