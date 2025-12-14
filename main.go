package main

import (
	"fmt"
	"net/http"

	"github.com/Neo2308/indianrailways-gtfs/server"
	"github.com/gorilla/mux"
)

func main() {

	irServer := server.NewServer()
	irServer.Setup()
	r := mux.NewRouter()
	r.HandleFunc("/map", irServer.HandleGetMap).Methods("GET")
	r.HandleFunc("/map/{trainNumber}", irServer.HandleGetMapForTrain).Methods("GET")
	r.HandleFunc("/map/search/{prefixText}", irServer.HandleGetMapBySearch).Methods("GET")
	r.HandleFunc("/map/liveStation/{stationCode}", irServer.HandleGetMapByLiveStation).Methods("GET")
	r.HandleFunc("/map/show/all", irServer.PopulateAllTrains).Methods("GET")
	r.HandleFunc("/save", irServer.SaveGTFS).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	fmt.Println(err)
}
