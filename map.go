package main

import (
	"fmt"
	"github.com/morikuni/go-geoplot"
	"image/color"
	"strconv"
)

type MapData struct {
	Map         *geoplot.Map
	StationIcon *geoplot.Icon
	Stations    map[string]*Station
}

func NewMapData() *MapData {
	tokyoTower := &geoplot.LatLng{
		Latitude:  24,
		Longitude: 78,
	}

	m := &geoplot.Map{
		Center: tokyoTower,
		Zoom:   7,
		Area: &geoplot.Area{
			From: tokyoTower.Offset(-14, -10),
			To:   tokyoTower.Offset(14, 20),
		},
	}
	return &MapData{
		Map:         m,
		StationIcon: geoplot.ColorIcon(58, 195, 112),
		Stations:    map[string]*Station{},
	}
}

func (m *MapData) AddStation(station *Station) {
	//fmt.Printf("Trying to added station:%s\n", station.Code)
	// Ignore stations that have already been added
	if _, ok := m.Stations[station.Code]; ok {
		return
	}
	m.Stations[station.Code] = station
	lat, _ := strconv.ParseFloat(station.Lat, 32)
	lng, _ := strconv.ParseFloat(station.Lng, 32)

	m.Map.AddMarker(&geoplot.Marker{
		LatLng: &geoplot.LatLng{
			Latitude:  lat,
			Longitude: lng,
		},
		Popup:   fmt.Sprintf("<b> %s [%s]</b>", station.Name, station.Code),
		Tooltip: station.Name,
		Icon:    m.StationIcon,
	})
	//fmt.Printf("Added station:%s at (%f,%f) \n", station.Code, lat, lng)
}

func (m *MapData) AddRoute(trainNumber int, routePoints []*geoplot.LatLng) {
	//fmt.Printf("Trying to added route:%d\n", trainNumber)
	//// Ignore stations that have already been added
	//if _, ok := m.Stations[station.Code]; ok {
	//	return
	//}
	//m.Stations[station.Code] = station
	//lat, _ := strconv.ParseFloat(station.Lat, 32)
	//lng, _ := strconv.ParseFloat(station.Lng, 32)

	m.Map.AddPolyline(&geoplot.Polyline{
		LatLngs: routePoints,
		Popup:   "Route of " + strconv.Itoa(trainNumber),
		Color:   &color.RGBA{255, 116, 0, 1},
	})
	//fmt.Println("Added route of ", trainNumber)
}

func (m *MapData) getMap() *geoplot.Map {
	return m.Map
}
