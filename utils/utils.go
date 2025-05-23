package utils

import (
	"BackProjetSurf/models"
	"encoding/json"
	"fmt"
	"os"
)

func LoadSpots() error {
	//Getenv est une fonction native de os
	file, err := os.Open(os.Getenv("JSON_FILE_PATH"))
	if err != nil {
		if os.IsNotExist(err) {
			// Si le fichier n'existe pas, initialisez une liste vide
			models.NosSpots = []models.SpotRecord{}
			return nil
		}
		return err
	}
	defer file.Close()

	var OneSpotData models.SpotData
	err = json.NewDecoder(file).Decode(&OneSpotData)
	if err != nil {
		return fmt.Errorf("format JSON invalide : %w", err)
	}

	// Charger les données dans la variable globale
	models.NosSpots = OneSpotData.Records
	return nil
}

func SaveSpots() error {
	file, err := os.Create(os.Getenv("JSON_FILE_PATH"))
	if err != nil {
		return err
	}
	defer file.Close()

	OneSpotData := models.SpotData{Records: models.NosSpots}
	json.NewEncoder(file).Encode(OneSpotData)
	return nil
}
