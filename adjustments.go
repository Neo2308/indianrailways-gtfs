package main

import (
	"fmt"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type DataErrors struct {
	ErroringStations map[string]Station `json:"erroringStations" yaml:"erroringStations"`
	WarningStations  map[string]Station `json:"warningStations" yaml:"warningStations"`
	Updates          map[string]Station `json:"updates" yaml:"updates"`
	Overrides        map[string]Station `json:"overrides" yaml:"overrides"`
	IgnoredTrains    []string           `json:"ignoredTrains" yaml:"ignoredTrains"`
}

func NewDataErrors() (*DataErrors, error) {
	data, err := loadFile("adjustments.yaml", FIXES)
	if err != nil {
		return nil, err
	}
	dataErrors := DataErrors{}
	err = yaml.Unmarshal(data, &dataErrors)
	if err != nil {
		return nil, err
	}
	dataErrors.ErroringStations = map[string]Station{}
	dataErrors.WarningStations = map[string]Station{}
	dataErrors.IgnoredTrains = []string{}
	return &dataErrors, nil
}

var stationsFixingErrors = 0
var stationsFixingWarnings = 0

func fixStation(station *Station, fixes *DataErrors) {
	if fix, ok := fixes.Overrides[station.Code]; ok {
		station.Code = fix.Code
		station.Name = fix.Name
		station.Lat = fix.Lat
		station.Lng = fix.Lng
	}
	if station.Lat == "" || station.Lng == "" {
		// fmt.Println("Fixing station has no latitude and longitude: ", station.Code)
		if fix, ok := fixes.Updates[station.Code]; ok {
			station.Lat = fix.Lat
			station.Lng = fix.Lng
			if strings.TrimSpace(station.Name) == "" {
				station.Name = fix.Name
				if strings.TrimSpace(fix.Name) == "" {
					station.Name = station.Code
					fmt.Println("!!! Warning !!! Used station code as name for station: ", station.Code)
					fixes.registerStationFixingWarning(station)
					// fmt.Println("!!! ERROR !!! No fix name found for station: ", station.Code)
				}
			}
			// fmt.Println("Fixed station: ", station)
		} else {
			fmt.Println("!!! ERROR !!! No fix of position found for station: ", station.Code)
		}
	} else if strings.TrimSpace(station.Name) == "" {
		if fix, ok := fixes.Updates[station.Code]; ok {
			station.Name = fix.Name
			if strings.TrimSpace(fix.Name) == "" {
				station.Name = station.Code
				fmt.Println("!!! Warning !!! Used station code as name for station: ", station.Code)
				fixes.registerStationFixingWarning(station)
			}
			// fmt.Println("Fixed station: ", station)
		}
	}
}

func stationHasProblems(station *Station, fixes *DataErrors) bool {
	if station.Lat == "" || station.Lng == "" || station.Lat == "0.000000" || station.Lng == "0.000000" || station.Lat == "1.000000" || station.Lng == "1.000000" || strings.TrimSpace(station.Name) == "" {
		fmt.Println("stationHasProblems: ", station)
		fixes.registerStationFixingError(station)
		return true
	}
	return false
}

// Deprecated: This function may be removed after fixing station issues
func trainHasProblems(train *TrainData, fixes *DataErrors) bool {
	tStations := train.getStations()

	if len(tStations) == 0 {
		fmt.Printf("Skipping train %s because it has no valid stations\n", train.getTrainNumber())
		fixes.IgnoredTrains = append(fixes.IgnoredTrains, train.getTrainNumber())
		return true
	}

	return false
}

func (d *DataErrors) registerStationFixingError(station *Station) {
	stationsFixingErrors++
	if _, ok := d.ErroringStations[station.Code]; !ok {
		d.ErroringStations[station.Code] = *station
	}
}

func (d *DataErrors) registerStationFixingWarning(station *Station) {
	stationsFixingWarnings++
	if _, ok := d.WarningStations[station.Code]; !ok {
		d.WarningStations[station.Code] = *station
	}
}

func (d *DataErrors) getStationFixingReport() string {
	return fmt.Sprintf("Station fixing report:\n\t %d total errors,\n\t %d total warnings,\n\t %d unique error stations,\n\t %d unique warning stations,\n\t", stationsFixingErrors, stationsFixingWarnings, len(d.ErroringStations), len(d.WarningStations))
}

func (d *DataErrors) getOverallFixingReport() string {
	slices.Sort(d.IgnoredTrains)
	d.Save()
	return fmt.Sprintf("Overall fixing report:\n\t %d total errors,\n\t %d total warnings,\n\t %d unique error stations,\n\t %d unique warning stations,\n\t %d unique trains ignored,\n\t", stationsFixingErrors, stationsFixingWarnings, len(d.ErroringStations), len(d.WarningStations), len(d.IgnoredTrains))
}

func (d *DataErrors) Save() {
	// Marshal to yaml
	yamlData, err := yaml.Marshal(d)
	if err != nil {
		fmt.Println("Error marshalling to YAML: ", err)
		panic(err)
	}
	err = saveFile("adjustments.yaml", yamlData, FIXES)
	if err != nil {
		panic(err)
	}
}
