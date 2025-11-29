package main

import (
	"fmt"
	"strings"
)

var stationsFixingErrors = 0
var stationsFixingWarnings = 0
var erroringStations = map[string]struct{}{}
var warningStations = map[string]struct{}{}

func fixStation(station *Station) {
	if fix, ok := IncorrectStationCodes[station.Code]; ok {
		station.Code = fix.Code
		station.Name = fix.Name
		station.Lat = fix.Lat
		station.Lng = fix.Lng
	}
	if station.Lat == "" || station.Lng == "" {
		// fmt.Println("Fixing station has no latitude and longitude: ", station.Code)
		if fix, ok := MissingStations[station.Code]; ok {
			station.Lat = fix.Lat
			station.Lng = fix.Lng
			if strings.TrimSpace(station.Name) == "" {
				station.Name = fix.Name
				if strings.TrimSpace(fix.Name) == "" {
					station.Name = station.Code
					fmt.Println("!!! Warning !!! Used station code as name for station: ", station.Code)
					registerStationFixingWarning(station)
					// fmt.Println("!!! ERROR !!! No fix name found for station: ", station.Code)
				}
			}
			// fmt.Println("Fixed station: ", station)
		} else {
			fmt.Println("!!! ERROR !!! No fix of position found for station: ", station.Code)
		}
	} else if strings.TrimSpace(station.Name) == "" {
		if fix, ok := MissingStations[station.Code]; ok {
			station.Name = fix.Name
			if strings.TrimSpace(fix.Name) == "" {
				station.Name = station.Code
				fmt.Println("!!! Warning !!! Used station code as name for station: ", station.Code)
				registerStationFixingWarning(station)
			}
			// fmt.Println("Fixed station: ", station)
		}
	}
}

func stationHasProblems(station *Station) bool {
	if station.Lat == "" || station.Lng == "" || station.Lat == "0.000000" || station.Lng == "0.000000" || station.Lat == "1.000000" || station.Lng == "1.000000" || strings.TrimSpace(station.Name) == "" {
		fmt.Println("stationHasProblems: ", station)
		registerStationFixingError(station)
		return true
	}
	return false
}

func registerStationFixingError(station *Station) {
	stationsFixingErrors++
	if _, ok := erroringStations[station.Code]; !ok {
		erroringStations[station.Code] = struct{}{}
	}
}

func registerStationFixingWarning(station *Station) {
	stationsFixingWarnings++
	if _, ok := warningStations[station.Code]; !ok {
		warningStations[station.Code] = struct{}{}
	}
}

func getStationFixingReport() string {
	return fmt.Sprintf("Station fixing report:\n\t %d total errors,\n\t %d total warnings,\n\t %d unique error stations,\n\t %d unique warning stations,\n\t", stationsFixingErrors, stationsFixingWarnings, len(erroringStations), len(warningStations))
}
