package main

import (
	"archive/zip"
	"bytes"
	"path/filepath"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/Neo2308/indianrailways-gtfs/types"
)

const output_path = "gtfs/"

type GtfsWriter struct {
	agencies  []types.Agency
	feedInfo  []types.FeedInfo
	stops     []types.Stop
	routes    []types.Route
	trips     []types.Trip
	stopTimes []types.StopTime
	calendar  []types.Calendar
}

func NewGtfsWriter() *GtfsWriter {
	return &GtfsWriter{
		agencies: []types.Agency{
			{
				AgencyId:       "1",
				AgencyName:     "Indian Railways",
				AgencyUrl:      "https://indianrailways.gov.in/",
				AgencyTimezone: "Asia/Kolkata",
				AgencyLang:     "en",
				CEMVSupport:    types.CEMVSupportSupported,
			},
		},
		feedInfo: []types.FeedInfo{
			{
				FeedPublisherName: "P. Radha Krishna",
				FeedPublisherUrl:  "https://github.com/Neo2308",
				FeedLang:          "en",
				FeedStartDate:     types.Date{Time: time.Now()},
				FeedEndDate:       types.Date{Time: time.Now().AddDate(0, 1, 0)},
				FeedVersion:       time.Now().Format("2006-01-02-15-04-05"),
				FeedContactEmail:  "pradha.krishna.cse17@itbhu.ac.in",
				FeedContactUrl:    "",
			},
		},
		stops:     []types.Stop{},
		routes:    []types.Route{},
		trips:     []types.Trip{},
		stopTimes: []types.StopTime{},
		calendar:  []types.Calendar{},
	}
}

func (g *GtfsWriter) AddStop(stop types.Stop) {
	g.stops = append(g.stops, stop)
}

func (g *GtfsWriter) AddRoute(route types.Route) {
	g.routes = append(g.routes, route)
}

func (g *GtfsWriter) AddTrips(trip types.Trip) {
	g.trips = append(g.trips, trip)
}

func (g *GtfsWriter) AddStopTimes(stopTimes []types.StopTime) {
	g.stopTimes = append(g.stopTimes, stopTimes...)
}

func (g *GtfsWriter) AddCalendar(calendar types.Calendar) {
	g.calendar = append(g.calendar, calendar)
}

func (g *GtfsWriter) WriteToZip() error {
	if err := g.writeCSVFile(g.agencies, "agency.txt"); err != nil {
		return err
	}
	if err := g.writeCSVFile(g.feedInfo, "feed_info.txt"); err != nil {
		return err
	}
	if err := g.writeCSVFile(g.stops, "stops.txt"); err != nil {
		return err
	}
	if err := g.writeCSVFile(g.calendar, "calendar.txt"); err != nil {
		return err
	}
	if err := g.writeCSVFile(g.routes, "routes.txt"); err != nil {
		return err
	}
	if err := g.writeCSVFile(g.trips, "trips.txt"); err != nil {
		return err
	}
	if err := g.writeCSVFile(g.stopTimes, "stop_times.txt"); err != nil {
		return err
	}
	if err := g.actuallyWriteToZip(); err != nil {
		return err
	}
	return nil
}

func (g *GtfsWriter) writeCSVFile(data interface{}, fileName string) error {
	writeBytes, err := gocsv.MarshalBytes(data)
	if err != nil {
		return err
	}
	// Save to file
	outputFile := filepath.Join(output_path, fileName)
	return saveFile(outputFile, writeBytes, OUTPUT)
}

func (g *GtfsWriter) actuallyWriteToZip() error {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name string
		data interface{}
	}{
		{"agency.txt", g.agencies},
		{"feed_info.txt", g.feedInfo},
		{"stops.txt", g.stops},
		{"calendar.txt", g.calendar},
		{"routes.txt", g.routes},
		{"trips.txt", g.trips},
		{"stop_times.txt", g.stopTimes},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return err
		}
		writeBytes, err := gocsv.MarshalBytes(file.data)
		if err != nil {
			return err
		}
		_, err = f.Write(writeBytes)
		if err != nil {
			return err
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		return err
	}
	zipBytes := buf.Bytes()
	return saveFile("gtfs.zip", zipBytes, OUTPUT)
}
