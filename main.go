package main

import (
	"github.com/gorilla/mux"
	houseHandler "house/handler"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	houseRouter := router.PathPrefix("/api/v1/house").Subrouter()

	houseRouter.Path("/").HandlerFunc(houseHandler.AddHouseHandler()).Methods("POST")
	houseRouter.Path("/").HandlerFunc(houseHandler.FindAllHousesHandler()).Methods("GET")
	houseRouter.Path("/{id}").HandlerFunc(houseHandler.FindHouseByIdHandler()).Methods("GET")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":3000", router))
}
