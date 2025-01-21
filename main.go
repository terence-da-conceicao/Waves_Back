package main

import (
	"BackProjetSurf/router"
	"BackProjetSurf/utils"
	"log"
)

// "BackProjetSurf/models"
// "BackProjetSurf/handlers"
// "BackProjetSurf/router"
// "BackProjetSurf/utils"

//go get lien url

// le projet doit avoir uniquement un seul module. Chaque dossier est un package et le projet entier doit être un module. Ce terme varie selon els langages.
// Le code attends que le go.mod soit à la racine du projet. Le go.mod doit donc être dans le dosseir parent.
//Les packages sont des dossiers. Chaque fichier dans un dossier doit être du même package.
//Donc on fait un dossier avec un nom (pas obligé que ce soit le nom du package mais c'est plus pratique si ça l'est)
//Puis les fichiers dans ce dossier auront tous package main par exemple.
// On importe le package souhaité dans main
//go get "BackProjetSurf/models" situé dans /main avec la console.

func main() {
	err := utils.LoadSpots()
	if err != nil {
		log.Fatalf("Erreur lors du chargement des spots : %v", err)
	}
	router.Router()
}
