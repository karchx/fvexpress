package products

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/KenethSandoval/fvexpress/pkg/db"
)

var (
	dbr = db.MongoCN.Database("fvexpress")
	col = dbr.Collection("products")
)

func GetProducts(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var resultado []Product

	condicion := bson.M{}

	cursor, err := col.Find(ctx, condicion)
	if err != nil {
		log.Fatal(err.Error())
	}

	for cursor.Next(context.TODO()) {
		var registro Product
		err := cursor.Decode(&registro)
		if err != nil {
			log.Fatal(err.Error())
		}
		resultado = append(resultado, registro)
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}

func GetOneProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var resultado []Product

	params := mux.Vars(r)

	condicion := bson.M{
		"_id": params["id"],
	}

	cursor, err := col.Find(ctx, condicion)
	if err != nil {
		log.Fatal(err.Error())
	}

	for cursor.Next(context.TODO()) {
		var registro Product
		err := cursor.Decode(&registro)
		if err != nil {
			log.Fatal(err.Error())
		}
		resultado = append(resultado, registro)
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}

func CreateProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var products Product
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// crear response
		return
	}

	newProduct := Product{
		Id:      primitive.NewObjectID(),
		Name:    products.Name,
		Image:   products.Image,
		Total:   products.Total,
		Price:   products.Price,
		SoldOut: products.SoldOut,
	}

	result, err := col.InsertOne(ctx, newProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// crear response
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
