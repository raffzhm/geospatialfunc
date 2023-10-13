package geospatialfunc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type GeoJSONFeature struct {
	Type     string          `json:"type"`
	Geometry GeoJSONGeometry `json:"geometry"`
}

type GeoJSONGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func init() {
	// Inisialisasi koneksi MongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://rafiazhim:200405@geogis-1214005.8kz4kfx.mongodb.net/")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
}

// GetGeoJSON adalah fungsi yang akan dijalankan oleh Google Cloud Function.
func getGeoData(w http.ResponseWriter, r *http.Request) {
	// Tambahkan header CORS
	w.Header().Set("Access-Control-Allow-Origin", "https://raffzhm.github.io") // Ganti dengan domain yang diizinkan.
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	collection := client.Database("pointmap").Collection("pointmap")

	// Lakukan query ke MongoDB dan ambil data geospasial
	var features []GeoJSONFeature
	cursor, err := collection.Find(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error querying MongoDB: %v", err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var feature GeoJSONFeature
		if err := cursor.Decode(&feature); err != nil {
			log.Fatalf("Error decoding MongoDB document: %v", err)
		}
		features = append(features, feature)
	}

	// Kembalikan hasil sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(features); err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
}

func main() {
	http.HandleFunc("/getGeoData", getGeoData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
