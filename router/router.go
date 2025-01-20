package router

import (
	"log"
	"net/http"

	"BackProjetSurf/handlers"

	"github.com/gorilla/mux"
)

func Router() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/spots", handlers.CreateSpot).Methods("POST")
	router.HandleFunc("/spots", handlers.GetAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", handlers.GetOneSpot).Methods("GET")
	router.HandleFunc("/spots/{id}", handlers.UpdateSpot).Methods("PATCH")
	router.HandleFunc("/spots/{id}", handlers.DeleteSpot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
