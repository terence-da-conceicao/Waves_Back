package handlers

import (
	"BackProjetSurf/models"
	"BackProjetSurf/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateSpot(w http.ResponseWriter, r *http.Request) {
	var newSpot models.SpotRecord
	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Impossible de lire le corps de la requête", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(content, &newSpot)
	if err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}
	models.NosSpots = append(models.NosSpots, newSpot)
	err = utils.SaveSpots()
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSpot)
}

func GetAllSpots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.NosSpots)
}

func GetOneSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]
	for _, singleSpot := range models.NosSpots {
		if singleSpot.ID == spotID {
			json.NewEncoder(w).Encode(singleSpot)
			return
		}
	}
	http.Error(w, "Spot non trouvé", http.StatusNotFound)
}

func UpdateSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]
	var updatedSpot models.SpotRecord

	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Impossible de lire le corps de la requête", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(content, &updatedSpot)
	if err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	for i, singleSpot := range models.NosSpots {
		if singleSpot.ID == spotID {
			models.NosSpots[i] = updatedSpot
			err = utils.SaveSpots()
			if err != nil {
				http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(updatedSpot)
			return
		}
	}

	http.Error(w, "Spot non trouvé", http.StatusNotFound)
}

func DeleteSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for i, singleSpot := range models.NosSpots {
		if singleSpot.ID == spotID {
			models.NosSpots = append(models.NosSpots[:i], models.NosSpots[i+1:]...)
			err := utils.SaveSpots()
			if err != nil {
				http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Le Spot avec l'ID %v a été supprimé avec succès", spotID)
			return
		}
	}

	http.Error(w, "Spot non trouvé", http.StatusNotFound)
}
