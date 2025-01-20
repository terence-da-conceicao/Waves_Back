package models

const JsonFilePath = "spots.json"

var NosSpots []SpotRecord

type Photo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type SpotFields struct {
	SurfBreak               []string `json:"Surf Break"`
	Destination             string   `json:"Destination"`
	Photos                  []Photo  `json:"Photos"`
	DestinationStateCountry string   `json:"Destination State/Country"`
	Address                 string   `json:"Address"`
}

type SpotRecord struct {
	ID     string     `json:"id"`
	Fields SpotFields `json:"fields"`
}

type SpotData struct {
	Records []SpotRecord `json:"records"`
	Offset  string       `json:"offset"`
}
