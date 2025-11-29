package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Neo2308/indianrailways-gtfs/types"
	"github.com/gorilla/mux"
	"github.com/morikuni/go-geoplot"
)

const timezone = "Asia/Kolkata"

type Server struct {
	// mapData   *MapData
	trainData       map[int]*TrainData
	stations        map[string]*Station
	runningDays     map[string]struct{}
	TrainListData   map[string]*TrainListData
	LiveStationData map[string]*LiveStationData
}

func NewServer() *Server {
	trainData := map[int]*TrainData{}
	stations := map[string]*Station{}
	return &Server{
		trainData:       trainData,
		stations:        stations,
		runningDays:     map[string]struct{}{},
		TrainListData:   map[string]*TrainListData{},
		LiveStationData: map[string]*LiveStationData{},
	}
}

func (s *Server) Setup() error {
	trains := []string{"22645", "22646"}
	for _, train := range trains {
		trainNumber, _ := strconv.Atoi(train)
		err := s.AddTrain(trainNumber)
		if err != nil {
			fmt.Printf("Error adding train %s: %e\n", train, err)
			return err
		}
	}
	trainLists := []string{"VANDE", "SHATABDI", "RAJDHANI"}
	for _, trainList := range trainLists {
		// trainNumber, _ := strconv.Atoi(train)
		err := s.AddTrainList(trainList)
		if err != nil {
			fmt.Printf("Error adding trainList %s: %e\n", trainList, err)
			return err
		}
	}
	return nil
}

func (s *Server) AddTrain(trainNumber int) error {
	newTrain := NewTrainData(trainNumber)
	s.trainData[trainNumber] = newTrain
	err := newTrain.populateData()
	if err != nil {
		fmt.Printf("Error fetching train service profile for %d: %e\n", trainNumber, err)
		return err
	}
	newTrainStations := newTrain.getStations()
	for _, station := range newTrainStations {
		if _, ok := s.stations[station.Code]; !ok {
			s.stations[station.Code] = &station
			fixStation(&station)
			continue
		}
		if _, err := getLatLng(s.stations[station.Code].Lat, s.stations[station.Code].Lng); err != nil {
			s.stations[station.Code] = &station
			fixStation(&station)
		}
	}
	runningDays := newTrain.getRunningDays()
	if _, ok := s.runningDays[runningDays]; !ok {
		s.runningDays[runningDays] = struct{}{}
	}
	return nil
}

func (s *Server) AddTrainList(prefixText string) error {
	newList := NewTrainListData(prefixText)
	s.TrainListData[prefixText] = newList
	err := newList.populateData()
	if err != nil {
		fmt.Println("Error fetching search results for", prefixText, ":", err)
		return err
	}
	newTrains := newList.getTrains()
	for _, trainNumber := range newTrains {
		if _, ok := s.trainData[trainNumber]; ok {
			continue
		}
		err = s.AddTrain(trainNumber)
		if err != nil {
			fmt.Printf("Error adding train %d from search %s: %e\n", trainNumber, prefixText, err)
			return err
		}
	}

	return nil
}

func (s *Server) AddLiveStationInfo(stationCode string) error {
	newList := NewLiveStationData(stationCode)
	s.LiveStationData[stationCode] = newList
	err := newList.populateData()
	if err != nil {
		fmt.Println("Error fetching live station results for", stationCode, ":", err)
		return err
	}
	newTrains := newList.getTrains()
	for _, trainNumber := range newTrains {
		if _, ok := s.trainData[trainNumber]; ok {
			continue
		}
		err = s.AddTrain(trainNumber)
		if err != nil {
			fmt.Printf("Error adding train %d from live station %s: %e\n", trainNumber, stationCode, err)
			return err
		}
	}

	return nil
}

func (s *Server) generateMap() *MapData {
	trainMap := NewMapData()
	for _, train := range s.trainData {
		trainMap.AddRoute(train.getNumber(), s.generateRouteForTrain(train.getNumber()))
	}
	for _, station := range s.stations {
		trainMap.AddStation(station)
	}
	return trainMap
}

func (s *Server) generateMapForTrain(trainNumber int, showMarkers bool) *MapData {
	trainMap := NewMapData()
	trainMap.AddRoute(trainNumber, s.generateRouteForTrain(trainNumber))
	if showMarkers {
		for _, station := range s.stations {
			trainMap.AddStation(station)
		}
	}
	return trainMap
}

func (s *Server) generateMapForTrainBySearch(prefixText string, showMarkers bool) *MapData {
	trainMap := NewMapData()
	trainsList := s.TrainListData[prefixText].getTrains()
	for _, trainNumber := range trainsList {
		trainMap.AddRoute(trainNumber, s.generateRouteForTrain(trainNumber))
	}
	if showMarkers {
		for _, station := range s.stations {
			trainMap.AddStation(station)
		}
	}
	return trainMap
}

func (s *Server) generateMapForTrainByLiveStation(stationCode string, showMarkers bool) *MapData {
	trainMap := NewMapData()
	trainsList := s.LiveStationData[stationCode].getTrains()
	for _, trainNumber := range trainsList {
		trainMap.AddRoute(trainNumber, s.generateRouteForTrain(trainNumber))
	}
	if showMarkers {
		for _, station := range s.stations {
			trainMap.AddStation(station)
		}
	} else {
		if station, ok := s.stations[stationCode]; ok {
			trainMap.AddStation(station)
		}
	}
	return trainMap
}

func (s *Server) generateRouteForTrain(trainNumber int) []*geoplot.LatLng {
	route := []*geoplot.LatLng{}
	stationCodes := s.trainData[trainNumber].getRoute()
	// fmt.Println(stationCodes)
	for _, v := range stationCodes {
		if _, ok := s.stations[v]; !ok {
			// fmt.Printf("Station code %s not found in stations map\n", v)
			continue
		}
		if nextPoint, err := getLatLng(s.stations[v].Lat, s.stations[v].Lng); err == nil {
			route = append(route, nextPoint)
		}
	}
	return route
}

func (s *Server) handleGetMap(w http.ResponseWriter, r *http.Request) {
	err := geoplot.ServeMap(w, r, s.generateMap().getMap())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleGetMapForTrain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainNumber, _ := strconv.Atoi(vars["trainNumber"])
	// Check if train number is valid
	if _, ok := s.trainData[trainNumber]; !ok {
		fmt.Printf("Train number %d not found\n", trainNumber)
		err := s.AddTrain(trainNumber)
		if err != nil {
			fmt.Printf("Error adding train %d: %e\n", trainNumber, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not add train number"))
			return
		}
		if _, ok := s.trainData[trainNumber]; !ok {
			fmt.Printf("Train number %d not found even after adding\n", trainNumber)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Train number not found"))
			return
		}
	}
	err := geoplot.ServeMap(w, r, s.generateMapForTrain(trainNumber, false).getMap())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleGetMapBySearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prefixText := vars["prefixText"]
	// Check if train number is valid
	if _, ok := s.TrainListData[prefixText]; !ok {
		// fmt.Printf("Train search for keyword %s not found\n", prefixText)
		err := s.AddTrainList(prefixText)
		if err != nil {
			fmt.Printf("Error search for keyword %s: %e\n", prefixText, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not search for keyword"))
			return
		}
		if _, ok := s.TrainListData[prefixText]; !ok {
			fmt.Printf("Train search for keyword %s not found even after adding\n", prefixText)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Could not search for keyword"))
			return
		}
	}
	err := geoplot.ServeMap(w, r, s.generateMapForTrainBySearch(prefixText, false).getMap())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) handleGetMapByLiveStation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stationCode := vars["stationCode"]
	// Check if train number is valid
	if _, ok := s.LiveStationData[stationCode]; !ok {
		fmt.Printf("Live station info for %s not found\n", stationCode)
		err := s.AddLiveStationInfo(stationCode)
		if err != nil {
			fmt.Printf("Error getting Live station info for %s: %e\n", stationCode, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could get Live station info"))
			return
		}
		if _, ok := s.LiveStationData[stationCode]; !ok {
			fmt.Printf("Live station info for %s not found even after adding\n", stationCode)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Could get Live station info"))
			return
		}
	}
	err := geoplot.ServeMap(w, r, s.generateMapForTrainByLiveStation(stationCode, false).getMap())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) populateAllTrains(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// prefixText := vars["prefixText"]
	//

	for i := 0; i <= 999; i++ {
		// if _, ok := s.trainData[i]; ok {
		//	continue
		// }
		prefixText := fmt.Sprintf("%03d", i)
		// Check if train number is valid
		if _, ok := s.TrainListData[prefixText]; !ok {
			// fmt.Printf("Train search for keyword %s not found\n", prefixText)
			_ = s.AddTrainList(prefixText)
		}
	}
	fmt.Println(getStationFixingReport())
	err := geoplot.ServeMap(w, r, s.generateMap().getMap())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) populateGTFSWriter(gw *GtfsWriter) {
	for _, station := range s.stations {
		// TODO: Remove after fixing station issues
		if stationHasProblems(station) {
			fmt.Printf("Skipping station with problems: %+v\n", station)
			continue
		}
		lat, _ := strconv.ParseFloat(station.Lat, 64)
		lng, _ := strconv.ParseFloat(station.Lng, 64)
		gw.AddStop(types.Stop{
			StopId:             station.Code,
			StopCode:           station.Code,
			StopName:           station.Name,
			TTSStopName:        station.Name,
			StopLat:            lat,
			StopLon:            lng,
			LocationType:       types.LocationTypeStop,
			StopTimezone:       timezone,
			WheelchairBoarding: types.WheelChairBoardingNoInfo,
			StopAccess:         types.StopAccessUnknown,
		})
	}
	for runningDays := range s.runningDays {
		if len(runningDays) != 7 {
			fmt.Printf("Invalid running days string: %s\n", runningDays)
			continue
		}
		cal := types.Calendar{
			ServiceId: getServiceId(runningDays),
			StartDate: types.Date{Time: time.Now()},
			EndDate:   types.Date{Time: time.Now().AddDate(0, 1, 0)},
		}
		if runningDays[0] == '1' {
			cal.Monday = 1
		}
		if runningDays[1] == '1' {
			cal.Tuesday = 1
		}
		if runningDays[2] == '1' {
			cal.Wednesday = 1
		}
		if runningDays[3] == '1' {
			cal.Thursday = 1
		}
		if runningDays[4] == '1' {
			cal.Friday = 1
		}
		if runningDays[5] == '1' {
			cal.Saturday = 1
		}
		if runningDays[6] == '1' {
			cal.Sunday = 1
		}
		gw.AddCalendar(cal)
	}

	for _, train := range s.trainData {
		// TODO: Remove after fixing station issues
		tStations := train.getStations()
		// ok := true
		if len(tStations) == 0 {
			fmt.Printf("Skipping train %s because it has no valid stations\n", train.getTrainNumber())
			continue
		}
		// for _, station := range tStations {
		//	fixStation(&station)
		//	if stationHasProblems(&station) {
		//		fmt.Printf("Skipping train %s because of station %s\n", station, train.getTrainNumber())
		//		ok = false
		//		break
		//	}
		// }
		// if !ok {
		//	continue
		// }
		gw.AddRoute(train.toRoute())
		gw.AddTrips(train.toTrip())
		gw.AddStopTimes(train.toStopTimes())
	}

	// for trainNumber, train := range s.trainData {
	//	serviceId := fmt.Sprintf("S%d", trainNumber)
	// }

}

func (s *Server) saveGTFS(w http.ResponseWriter, r *http.Request) {
	gw := NewGtfsWriter()
	s.populateGTFSWriter(gw)
	err := gw.WriteToZip()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := loadFile("gtfs.zip", OUTPUT)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=gtfs.zip")
	http.ServeContent(w, r, "gtfs.zip", time.Now(), bytes.NewReader(data))
}
