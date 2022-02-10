package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// ### DECOUVERTE DU GO ####
// A partir d'une base de données basique (voir fichier data.sql) le but pour cette première approche du go
// sera de créer 3 routes afin de pouvoir:
// - récupérer tous les aliments de la BDD
// - récupérer un seul aliment via une id en paramètre d'url
// - ajouter un nouvel aliment en BDD

type Food struct {
	Name     string
	Quantity int
	Category string
}

func main() {

	// conect to the database
	db, err := sql.Open("pgx", "postgresql://localhost:5432/fridge")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to access data: %v", err)
	}

	fmt.Println("Database access established")

	// GET all food
	http.HandleFunc("/getAllFood", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getAllFood" {
			http.Error(rw, "404 page not found", http.StatusNotFound)
			return
		}

		if r.Method != "GET" {
			http.Error(rw, "This route expects a GET http method", http.StatusNotFound)
			return
		}

		rows, err := db.Query("SELECT food.name, food.quantity, category.name FROM food JOIN food_category ON food.id = food_category.food_id JOIN category ON food_category.category_id = category.id")
		if err != nil {
			log.Fatalf("could not execute query: %v", err)
		}

		foods := []Food{}

		for rows.Next() {
			food := Food{}
			if err := rows.Scan(&food.Name, &food.Quantity, &food.Category); err != nil {
				log.Fatalf("could not scan row: %v", err)
			}
			foods = append(foods, food)
		}

		fmt.Printf("found %d food: %+v", len(foods), foods)
		// Comment renvoyer la donnée au client ?

	})

	// GET one food
	http.HandleFunc("/getOneFood", func(rw http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/getOneFood" {
			http.Error(rw, "404 page not found", http.StatusNotFound)
			return
		}

		if r.Method != "GET" {
			http.Error(rw, "This route expects a GET http method", http.StatusNotFound)
			return
		}

		parameters, ok := r.URL.Query()["p"]

		if !ok || len(parameters[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}

		if len(parameters[0]) > 1 {
			log.Println("only one url parameter is expected")
		}

		p := parameters[0]

		food := Food{}

		row := db.QueryRow("SELECT food.name, food.quantity FROM food WHERE food.id=$1", p)
		if err := row.Scan(&food.Name, &food.Quantity); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}

		fmt.Printf("food: %+v", food)
		// Même problème que pour la requête précedente, comment envoyer la donnée côté client ?

	})

	// ADD one food
	http.HandleFunc("/addFood", func(rw http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/addFood" {
			http.Error(rw, "404 page not found", http.StatusNotFound)
			return
		}

		if r.Method != "POST" {
			http.Error(rw, "This route expects a POST http method", http.StatusNotFound)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var newFood Food
		if err := decoder.Decode(&newFood); err != nil {
			panic(err)
		}
		log.Println("JSON body : ", newFood)

		result, err := db.Exec("INSERT INTO food (name, quantity) VALUES ($1, $2)", newFood.Name, newFood.Quantity)
		if err != nil {
			log.Fatalf("could not insert row: %v", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("could not get affected rows: %v", err)
		}

		fmt.Println("inserted", rowsAffected, "rows")
	})

	// local server
	fmt.Printf("server reachable on localhost:8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
