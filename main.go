package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const jsonFilePath = "spots.json"

// Structures pour correspondre au JSON
type Photo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
	/* Filename   string `json:"filename"`
	Size       int    `json:"size"`
	Type       string `json:"type"` */
	/* Thumbnails struct {
		Small struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"small"`
		Large struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"large"`
		Full struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"full"`
	} `json:"thumbnails"` */
}

type SpotFields struct {
	SurfBreak []string `json:"Surf Break"`
	/* DifficultyLevel         int      `json:"Difficulty Level"` */
	Destination string `json:"Destination"`
	/* Geocode                 string   `json:"Geocode"`
	Influencers             []string `json:"Influencers"`
	MagicSeaweedLink        string   `json:"Magic Seaweed Link"` */
	Photos []Photo `json:"Photos"`
	/* PeakSurfSeasonBegins    string   `json:"Peak Surf Season Begins"` */
	DestinationStateCountry string `json:"Destination State/Country"`
	/* PeakSurfSeasonEnds      string   `json:"Peak Surf Season Ends"` */
	Address string `json:"Address"`
}

type SpotRecord struct {
	ID     string     `json:"id"`
	Fields SpotFields `json:"fields"`
}

type SpotData struct {
	Records []SpotRecord `json:"records"`
	Offset  string       `json:"offset"`
}

var nosSpots []SpotRecord

// Charger les spots depuis le fichier JSON
func loadSpots() error {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Si le fichier n'existe pas, initialisez une liste vide
			nosSpots = []SpotRecord{}
			return nil
		}
		return err
	}
	defer file.Close()

	var spotData SpotData
	err = json.NewDecoder(file).Decode(&spotData)
	if err != nil {
		return fmt.Errorf("format JSON invalide : %w", err)
	}

	// Charger les données dans la variable globale
	nosSpots = spotData.Records
	return nil
}

// Sauvegarder les spots dans le fichier JSON
func saveSpots() error {
	file, err := os.Create(jsonFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	spotData := SpotData{Records: nosSpots}
	json.NewEncoder(file).Encode(spotData)
	return nil
}

// Handlers CRUD
func createSpot(w http.ResponseWriter, r *http.Request) {
	var newSpot SpotRecord
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

	nosSpots = append(nosSpots, newSpot)
	err = saveSpots()
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSpot)

}

func getAllSpots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nosSpots)
}

func getOneSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for _, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			json.NewEncoder(w).Encode(singleSpot)
			return
		}
	}
	http.Error(w, "Spot non trouvé", http.StatusNotFound)
}

func updateSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]
	var updatedSpot SpotRecord

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

	for i, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			nosSpots[i] = updatedSpot
			err = saveSpots()

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

func deleteSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for i, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			nosSpots = append(nosSpots[:i], nosSpots[i+1:]...)
			err := saveSpots()
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

func main() {
	err := loadSpots()
	if err != nil {
		log.Fatalf("Erreur lors du chargement des spots : %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/spots", createSpot).Methods("POST")
	router.HandleFunc("/spots", getAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", getOneSpot).Methods("GET")
	router.HandleFunc("/spots/{id}", updateSpot).Methods("PATCH")
	router.HandleFunc("/spots/{id}", deleteSpot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// "classe" (=type) d'un seul spot, avec des champs contenant des infos
/* type Spot struct {
	ID        string `json:"ID"`
	SurfBreak string `json:"Surf Break"`
	Photos    string `json:"Photos"`
	Address   string `json:"Address"`
} */

// "classe" (=type) de l'ensemble des spots en une seule structure
// Le type s'appelle ListeSpots, c'est une slice (controleur à distance) de spot.
//type ListeSpots []Spot

/* var nosSpots = ListeSpots{

	{
		ID:        "1",
		SurfBreak: "Plage",
		Photos:    "url",
		Address:   "australie",
	},
} */

/* var (
	nosSpots     ListeSpots
	jsonFilePath = "spots.json"
)
*/
// La fonction loadSpots charge les données depuis spots.json au démarrage.

/* func loadSpots() error {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Si le fichier n'existe pas, initialisez une liste vide
			nosSpots = ListeSpots{}
			return nil
		}
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&nosSpots)
	if err != nil {
		return err
	}
	return nil
} */

//La fonction saveSpots sauvegarde nosSpots dans spots.json après chaque modification.

/* func saveSpots() error {
	file, err := os.Create(jsonFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(nosSpots)
	if err != nil {
		return err
	}
	return nil
} */

// fonction d'affichage vers la page d'accueil
/* func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home")
} */

/* func createSpot(w http.ResponseWriter, r *http.Request) {
var newSpot Spot
//content, err := io.ReadAll(r.Body)
content, err := io.ReadAll(r.Body)
if err != nil {
	http.Error(w, "Impossible de lire le corps de la requête", http.StatusBadRequest)
	return
} */

/* json.Unmarshal(content, &newSpot)
nosSpots = append(nosSpots, newSpot)
w.WriteHeader(http.StatusCreated)

json.NewEncoder(w).Encode(newSpot) */
/* err = json.Unmarshal(content, &newSpot)
if err != nil {
	http.Error(w, "Format JSON invalide", http.StatusBadRequest)
	return
} */

/* nosSpots = append(nosSpots, newSpot)
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(newSpot) */

/* 	nosSpots = append(nosSpots, newSpot)
	err = saveSpots()
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSpot)
}*/

/* func getOneSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for _, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			json.NewEncoder(w).Encode(singleSpot)
			return
		}
	}
	http.Error(w, "Spot non trouvé", http.StatusNotFound)
} */

/* func getAllSpots(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(nosSpots)
} */

/* func getAllSpots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nosSpots)
} */

/* func updateSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]
	var updatedSpot Spot


	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "update échouée")
		content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Impossible de lire le corps de la requête", http.StatusBadRequest)
		return
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
}
*/

/* func updateSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]
	var updatedSpot Spot

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



	for i, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			updatedSpot.ID = spotID // Conserver l'ID d'origine
			nosSpots[i] = updatedSpot

			err = saveSpots()
			if err != nil {
				http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(updatedSpot)
			return
		}
	}
	http.Error(w, "Spot non trouvé", http.StatusNotFound)

} */

/* func deleteSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for i, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			nosSpots = append(nosSpots[:i], nosSpots[i+1:]...)
			fmt.Fprintf(w, "The Spot with ID %v has been deleted successfully", spotID)
		}
	}
} */

/* func deleteSpot(w http.ResponseWriter, r *http.Request) {
	spotID := mux.Vars(r)["id"]

	for i, singleSpot := range nosSpots {
		if singleSpot.ID == spotID {
			nosSpots = append(nosSpots[:i], nosSpots[i+1:]...)

			err := saveSpots()
			if err != nil {
				http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "The Spot with ID %v has been deleted successfully", spotID)
			return
		}
	}
	http.Error(w, "Spot non trouvé", http.StatusNotFound)
} */

/* func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/spot", createSpot).Methods("POST")
	router.HandleFunc("/spots", getAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", getOneSpot).Methods("GET")
	router.HandleFunc("/spots/{id}", updateSpot).Methods("PATCH")
	router.HandleFunc("/spots/{id}", deleteSpot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
*/

/* func main() {

	err := loadSpots()
	if err != nil {
		log.Fatalf("Erreur lors du chargement des spots : %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/spot", createSpot).Methods("POST")
	router.HandleFunc("/spots", getAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", getOneSpot).Methods("GET")
	router.HandleFunc("/spots/{id}", updateSpot).Methods("PATCH")
	router.HandleFunc("/spots/{id}", deleteSpot).Methods("DELETE")

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
*/
