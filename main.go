package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	server := NewServer()
	server.Setup()
	r := mux.NewRouter()
	r.HandleFunc("/map", server.handleGetMap).Methods("GET")
	r.HandleFunc("/map/{trainNumber}", server.handleGetMapForTrain).Methods("GET")
	r.HandleFunc("/map/search/{prefixText}", server.handleGetMapBySearch).Methods("GET")
	r.HandleFunc("/map/liveStation/{stationCode}", server.handleGetMapByLiveStation).Methods("GET")
	r.HandleFunc("/map/show/all", server.populateAllTrains).Methods("GET")
	r.HandleFunc("/save", server.saveGTFS).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	fmt.Println(err)
}
