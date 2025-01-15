package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// "classe" (=type) d'un seul spot, avec des champs contenant des infos
type Spot struct {
	ID        string `json:"ID"`
	SurfBreak string `json:"Surf Break"`
	Photos    string `json:"Photos"`
	Address   string `json:"Address"`
}

// "classe" (=type) de l'ensemble des spots en une seule structure
// Le type s'appelle ListeSpots, c'est une slice (controleur à distance) de spot.
type ListeSpots []Spot

var nosSpots = ListeSpots{
	//aller récupérer les infos du JSON
	{
		ID:        "1",
		SurfBreak: "Plage",
		Photos:    "url",
		Address:   "australie",
	},
}

// fonction d'affichage vers la page d'accueil
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home")
}

func createSpot(w http.ResponseWriter, r *http.Request) {
	var newSpot Spot
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "impossible de créer un spot")
	}

	json.Unmarshal(reqBody, &newSpot)
	nosSpots = append(nosSpots, newSpot)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newSpot)
}

func getOneSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	// _ : ignore un paramètre (met un itérateur nul, itéère mais l'itérateur ne possède aucunce valeur)
	// singleSpot : le spot qu'on récupère
	for _, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			json.NewEncoder(w).Encode(singleSpot)
		}
	}
}

func getAllSpots(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(nosSpots)
}

func updateSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]
	var updatedSpot Spot

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "update échouée")
	}
	json.Unmarshal(reqBody, &updatedSpot)

	for i, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			singleSpot.SurfBreak = updatedSpot.SurfBreak
			singleSpot.Photos = updatedSpot.Photos
			singleSpot.Address = updatedSpot.Address
			nosSpots = append(nosSpots[:i], singleSpot)
			json.NewEncoder(w).Encode(singleSpot)
		}
	}
}

func deleteSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for i, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			nosSpots = append(nosSpots[:i], nosSpots[i+1:]...)
			fmt.Fprintf(w, "The Spot with ID %v has been deleted successfully", spotID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/spot", createSpot).Methods("POST")
	router.HandleFunc("/spots", getAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", getOneSpot).Methods("GET")
	router.HandleFunc("/spots/{id}", updateSpot).Methods("PATCH")
	router.HandleFunc("/spots/{id}", deleteSpot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
