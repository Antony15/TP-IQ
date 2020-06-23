package main

import (
	"log"
	"math/rand"
	"net/http"
	"encoding/json"
	"gorilla/mux"
	perlin "github.com/aquilax/go-perlin"	
)

type request struct {
	X 	int 	`json:"x"`	
	Y 	int 	`json:"y"`	
}

type response struct {
	Output 	[]float64 	`json:"output"`	
}

const (
	alpha            = 2.
	beta             = 2.
	n                = 3
	maximumSeedValue = 100
)

// handleReq Function
func handleReq(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "application/json")
	var RequestArray request 
	var Response response 
	ResponseArray 			:= make(map[string]interface{})
	err 					:= json.NewDecoder(r.Body).Decode(&RequestArray)
	if err != nil {
		ResponseArray["message"] = "Request Error"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseArray)		
		return
	} else {
		for x := 0; x < RequestArray.X; x++ {
			for y := 0; y < RequestArray.Y; y++ {
				chanl := make(chan float64,1)
				go getCoordinateValue(chanl,x,y)
				Response.Output = append(Response.Output,<-chanl)
			}
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response.Output)
		return		
	}	
}

// getDate Function
func getCoordinateValue(chanl chan float64, x int, y int){
	var seed = rand.Int63n(maximumSeedValue)
	p := perlin.NewPerlin(alpha, beta, n, seed)
	val := p.Noise2D(float64(x)/10, float64(y)/10)
	chanl <- val
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
