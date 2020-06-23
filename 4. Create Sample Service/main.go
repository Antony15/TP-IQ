package main

import (
	"log"
	"time"
	"net/http"
	"encoding/json"
	"gorilla/mux"	
)

type request struct {
	Date 	string 	`json:"date"`	
}

type response struct {
	Date 	string 	`json:"date"`	
}


// handleReq Function
func handleReq(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "application/json")
	var RequestArray request 
	ResponseArray 			:= make(map[string]interface{})
	err 					:= json.NewDecoder(r.Body).Decode(&RequestArray)
	if err != nil {
		ResponseArray["message"] = "Request Error"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseArray)		
		return
	} else {
		responsee 		 	 := getDate(RequestArray)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responsee)
		return		
	}	
}

// getDate Function
func getDate(req request)(res response){
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout , req.Date+":00.000Z")
	if err != nil {
		log.Println(err)
	}
	return response{Date:t.Add(12*time.Hour).Format("2006-01-02T15:04")}
}

// Main Function
func main(){
	log.Println("Server started at http://localhost:9090")
 	// http.Handler
 	router := mux.NewRouter().StrictSlash(true) 	
 	router.HandleFunc("/", handleReq).Methods("POST") 	
 	if err := http.ListenAndServe(":9090", router); err != nil {
		panic(err)
	}
}
