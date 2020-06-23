package main

import (
	"log"	
	"strconv"
	"net/http"	
    "gorilla/mux"
    "text/template"		
	//~ "gopkg.in/mgo.v2"
	//~ "gopkg.in/mgo.v2/bson"	
	"mongo"
	"mongo/bson"	
)

// Database Mongo Settings & Functions
//Users represents a sample database entity.
type User struct {
	ID    		int 		`json:"id" bson:"_id,omitempty"`
	Name 		string    	`json:"name"`
	Age 		int    		`json:"age"`
	Address 	string    	`json:"address"`
}

var tmpl = template.Must(template.ParseGlob("form/*"))
var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost/crud_golang_mongodb")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db = session.DB("crud_golang_mongodb")
}

func collection() *mgo.Collection {
	return db.C("users")
}

// GetAll returns all users from the database.
func GetAll() ([]User, error) {
	res := []User{}

	if err := collection().Find(nil).All(&res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetOne returns a single user from the database.
func GetOne(id int) (*User, error) {
	res := User{}

	if err := collection().Find(bson.M{"_id": id}).One(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
// Update single user from the database.
func Updateuser(user User) error{
	err := 	collection().Update(bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"name": user.Name, "age": user.Age, "address": user.Address}})
	if err != nil {
		return err
	}

	return nil
}

// Save inserts an user to the database.
func Save(user User) error {
	return collection().Insert(user)
}

// Remove deletes an user from the database
func Remove(id int) error {
	return collection().Remove(bson.M{"_id": id})
}
// Database Mongo Settings & Functions End 

// Golang Functions Start
func Index(w http.ResponseWriter, r *http.Request){
	users,err := GetAll()
    if err != nil {
        panic(err.Error())
    }	
    tmpl.ExecuteTemplate(w, "Index", users)
}

func Show(w http.ResponseWriter, r *http.Request) {
    nId,_:= strconv.Atoi(r.URL.Query().Get("id"))
	user,err := GetOne(nId)
    if err != nil {
        panic(err.Error())
    }	
    tmpl.ExecuteTemplate(w, "Show", user)
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	IDstr 			:= r.FormValue("id")
	agestr			:= r.FormValue("age")
	name 			:= r.FormValue("name")
	address 		:= r.FormValue("address")
	ID,_			:= strconv.Atoi(IDstr)
	age, _ 			:= strconv.Atoi(agestr)
	user := User{ID: ID, Name: name, Age: age, Address: address}
	if err := Save(user); err != nil {
		log.Println("Failed to create user")
	}
    http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    nId,_:= strconv.Atoi(r.URL.Query().Get("id"))
	user,err := GetOne(nId)
    if err != nil {
        panic(err.Error())
    }	
    tmpl.ExecuteTemplate(w, "Edit", user)
}

func Update(w http.ResponseWriter, r *http.Request) {
	IDstr 			:= r.FormValue("id")
	agestr			:= r.FormValue("age")
	name 			:= r.FormValue("name")
	address 		:= r.FormValue("address")
	ID, _ 			:= strconv.Atoi(IDstr)
	age, _ 			:= strconv.Atoi(agestr)

	user := User{ID: ID, Name: name, Age: age, Address: address}
	err := Updateuser(user)
    if err != nil {
        panic(err.Error())
    }
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    nId,_:= strconv.Atoi(r.URL.Query().Get("id"))
	if err := Remove(nId); err != nil {
		log.Println("Failed to remove user")
	}
    http.Redirect(w, r, "/", 301)
}
// Golang Functions End

// Main Function Start
func main() {
	log.Println("Server started on: http://localhost:8088")
	r := mux.NewRouter().StrictSlash(true)
    r.HandleFunc("/", Index)
    r.HandleFunc("/show", Show).Methods("GET")
    r.HandleFunc("/new", New)
    r.HandleFunc("/edit", Edit).Methods("GET")
    r.HandleFunc("/insert", Insert).Methods("POST")
    r.HandleFunc("/update", Update).Methods("POST")
    r.HandleFunc("/delete", Delete).Methods("GET")    
 	if err := http.ListenAndServe(":8088", r); err != nil {
		panic(err)
	}   
}
// Main Function End
