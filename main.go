package main

import (
	"log"
	"net/http"

	"github.com/armeo/go-mongo-rest/app"
)

func main() {
	log.Println("Listening on 8000")
	http.ListenAndServe(":8000", app.NewRoute(app.NewMongoManager("localhost")))
}
