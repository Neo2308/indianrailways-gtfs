package gtfsWriter

import (
	"archive/zip"
	"bytes"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gocarina/gocsv"

	"github.com/Neo2308/indianrailways-gtfs/fileUtils"
	"github.com/Neo2308/indianrailways-gtfs/types"
)

const output_path = "gtfs/"

type GtfsWriter struct {
	agencies       []types.Agency
	feedInfo       []types.FeedInfo
	stops          []types.Stop
	routes         []types.Route
	trips          []types.Trip
	stopTimes      []types.StopTime
	calendar       []types.Calendar
	shapes         []types.Shape
	supportsShapes bool
}

func NewGtfsWriter(supportsShapes bool) *GtfsWriter {
	return &GtfsWriter{
		supportsShapes: supportsShapes,
		agencies:       []types.Agency{},
		feedInfo:       []types.FeedInfo{},
		stops:          []types.Stop{},
		routes:         []types.Route{},
		trips:          []types.Trip{},
		stopTimes:      []types.StopTime{},
		calendar:       []types.Calendar{},
	}
}

func (g *GtfsWriter) AddAgency(agency types.Agency) {
	g.agencies = append(g.agencies, agency)
}

func (g *GtfsWriter) AddFeedInfo(feedInfo types.FeedInfo) {
	g.feedInfo = append(g.feedInfo, feedInfo)
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

func (g *GtfsWriter) AddShapes(shapes []types.Shape) {
	g.shapes = append(g.shapes, shapes...)
}

func (g *GtfsWriter) AddCalendar(calendar types.Calendar) {
	g.calendar = append(g.calendar, calendar)
}

func (g *GtfsWriter) Sort() {
	slices.SortFunc(g.agencies, func(a, b types.Agency) int {
		return strings.Compare(a.AgencyId, b.AgencyId)
	})
	slices.SortFunc(g.feedInfo, func(a, b types.FeedInfo) int {
		return strings.Compare(a.FeedPublisherName, b.FeedPublisherName)
	})
	slices.SortFunc(g.stops, func(a, b types.Stop) int {
		return strings.Compare(a.StopId, b.StopId)
	})
	slices.SortFunc(g.calendar, func(a, b types.Calendar) int {
		// Sorting in descending order of running days
		return -strings.Compare(a.GetRunningDays(), b.GetRunningDays())
	})
	slices.SortFunc(g.routes, func(a, b types.Route) int {
		return strings.Compare(a.RouteId, b.RouteId)
	})
	slices.SortFunc(g.trips, func(a, b types.Trip) int {
		// RouteId is always unique and same as TripId in our case, but keeping for future proofing
		if a.RouteId != b.RouteId {
			return strings.Compare(a.RouteId, b.RouteId)
		}
		if a.ServiceId != b.ServiceId {
			return strings.Compare(a.ServiceId, b.ServiceId)
		}
		return strings.Compare(a.TripId, b.TripId)
	})
	slices.SortFunc(g.stopTimes, func(a, b types.StopTime) int {
		if a.TripId != b.TripId {
			return strings.Compare(a.TripId, b.TripId)
		}
		return a.StopSequence - b.StopSequence
	})
	slices.SortFunc(g.shapes, func(a, b types.Shape) int {
		if a.ShapeId != b.ShapeId {
			return strings.Compare(a.ShapeId, b.ShapeId)
		}
		return a.ShapePtSequence - b.ShapePtSequence
	})
}

func (g *GtfsWriter) WriteToZip() error {
	g.Sort()
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
	if g.supportsShapes {
		if err := g.writeCSVFile(g.shapes, "shapes.txt"); err != nil {
			return err
		}
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
	return fileUtils.SaveFile(outputFile, writeBytes, fileUtils.OUTPUT)
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
	if g.supportsShapes {
		files = append(files, struct {
			Name string
			data interface{}
		}{"shapes.txt", g.shapes})
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
	return fileUtils.SaveFile("gtfs.zip", zipBytes, fileUtils.OUTPUT)
}
