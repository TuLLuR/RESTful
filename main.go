package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	id       string
	name     string
	surname  string
	person   *Person
	country  string
	city     string
	postcode int
	adress   string
	domicile *Domicile
)

// Person struct
type Person struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Domicile *Domicile `json:"domicile"`
}

// Domicile struct
type Domicile struct {
	Country  string `json:"country"`
	City     string `json:"city"`
	Postcode int    `json:"postcode"`
	Adress   string `json:"adress"`
}

// Init person var as a slice Person struct
var persons []Person

// GET all person information
func getPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

// GET single person
func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // GET params
	// Loop through person and find with ID
	for _, item := range persons {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// Create new person
func createNewPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = strconv.Itoa(rand.Intn(999)) // random ID
	persons = append(persons, person)
	json.NewEncoder(w).Encode(person)

}

// Update person information
func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range persons {
		if item.ID == params["id"] {
			persons = append(persons[:index], persons[index+1:]...)
			var person Person
			_ = json.NewDecoder(r.Body).Decode(&person)
			person.ID = params["id"]
			persons = append(persons, person)
			json.NewEncoder(w).Encode(person)
		}
	}
	json.NewEncoder(w).Encode(persons)
	return
}

// Delete person
func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range persons {
		if item.ID == params["id"] {
			persons = append(persons[:index], persons[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(persons)
}

// Main function
func main() {
	fmt.Println("localhost:8080")

	// Init Multiplexer (Mux) Router
	r := mux.NewRouter()

	// Moc data - @todo implement DB
	// persons = append(persons, Person{ID: "1", Name: "Lonzo", Surname: "Ullrich", Domicile: &Domicile{Country: "Italy", City: "Bogisichfurt", Postcode: 14100, Adress: "26730 Robel Overpass"}})
	// persons = append(persons, Person{ID: "2", Name: "Modesto", Surname: "Gibson", Domicile: &Domicile{Country: "Kiribati", City: "Titoburgh", Postcode: 83198, Adress: "06745 Christopher Ramp"}})
	// persons = append(persons, Person{ID: "3", Name: "Elenora", Surname: "Aufderhar", Domicile: &Domicile{Country: "Malawi", City: "North Micaelaton", Postcode: 57152, Adress: "530 Antonia Canyon"}})
	fmt.Println("Data: ")

	file, err := ioutil.ReadFile("data.json")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &person)

	if err != nil {
		fmt.Println(err.Error())
	}

	id = person.ID
	name = person.Name
	surname = person.Surname
	domicile = person.Domicile
	country = domicile.Country
	city = domicile.City
	postcode = domicile.Postcode
	adress = domicile.Adress

	// make domcile var
	persons = append(persons, Person{ID: id, Name: name, Surname: surname, Domicile: &Domicile{Country: country, City: city, Postcode: postcode, Adress: adress}})
	// Route Hanlders / Endpoints
	r.HandleFunc("/api/person", getPersons).Methods("GET")
	r.HandleFunc("/api/person/{id}", getPerson).Methods("GET")
	r.HandleFunc("/api/person", createNewPerson).Methods("POST")
	r.HandleFunc("/api/person/{id}", updatePerson).Methods("PUT")
	r.HandleFunc("/api/person/{id}", deletePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))

}
