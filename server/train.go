package server

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/morikuni/go-geoplot"

	"github.com/Neo2308/indianrailways-gtfs/apiDataFetcher"
	"github.com/Neo2308/indianrailways-gtfs/models"
	"github.com/Neo2308/indianrailways-gtfs/types"
)

type TrainData struct {
	apiDataFetcher.ApiDataFetcher[models.TrainServiceProfileResponse]
	trainNumber                 int
	TrainServiceProfileResponse models.TrainServiceProfileResponse
	dataErrors                  *DataErrors
}

func NewTrainData(trainNumber int, dataErrors *DataErrors) *TrainData {
	newObj := &TrainData{
		trainNumber: trainNumber,
		dataErrors:  dataErrors,
	}
	newObj.ApiDataFetcher = *apiDataFetcher.NewApiDataFetcher[models.TrainServiceProfileResponse](
		&newObj.TrainServiceProfileResponse,
		"train service profile",
		fmt.Sprintf("%s.json", newObj.getTrainNumber()),
		newObj.getTrainServiceProfileUncached,
	)
	return newObj
}

func (t *TrainData) getStations() []models.Station {
	_ = t.LoadData()
	stations := []models.Station{}
	// fmt.Println(t.TrainServiceProfileResponse)
	for _, v := range t.TrainServiceProfileResponse.Pd.TrainServiceProfile.VTrainServiceSchedulePTT {
		newStation := models.Station{
			Code: v.StationCode,
			Name: v.StationName,
			Lat:  v.Lattitude,
			Lng:  v.Longitude,
		}
		// TODO: Remove after fixing station issues
		fixStation(&newStation, t.dataErrors)
		if stationHasProblems(&newStation, t.dataErrors) {
			// fmt.Println("Station has problems: ", newStation)
			stations = []models.Station{}
			break
		}
		// fmt.Println(v)
		stations = append(stations, newStation)
	}
	return stations
}

func (t *TrainData) getRunningDays() string {
	// _ = t.LoadData()
	return t.TrainServiceProfileResponse.Pd.TrainServiceProfile.DaysOfRunFromSourceNumeric
}

// func (t *TrainData) getRoute() []*geoplot.LatLng {
//	_ = t.LoadData()
//	route := []*geoplot.LatLng{}
//	for _, v := range t.TrainServiceProfileResponse.Pd.TrainServiceProfile.VTrainServiceSchedulePTT {
//		fmt.Println(v)
//		if nextPoint, err := getLatLng(v.Lattitude, v.Longitude); err == nil {
//			route = append(route, nextPoint)
//		}
//	}
//	return route
// }

func (t *TrainData) getRoute() []string {
	_ = t.LoadData()
	route := []string{}
	// if t.TrainServiceProfileResponse.Pd.TrainServiceProfile.VTrainServiceScheduleWTT != nil {
	//	for _, v := range t.TrainServiceProfileResponse.Pd.TrainServiceProfile.VTrainServiceScheduleWTT {
	//		route = append(route, v.StationCode)
	//	}
	//	return route
	// }
	for _, v := range t.TrainServiceProfileResponse.Pd.TrainServiceProfile.VTrainServiceSchedulePTT {
		route = append(route, v.StationCode)
	}
	return route
}

func (t *TrainData) getNumber() int {
	return t.trainNumber
}

func (t *TrainData) getTrainNumber() string {
	return fmt.Sprintf("%05d", t.trainNumber)
}

func (t *TrainData) getTrainServiceProfileUncached() error {
	url := "https://apigw.umangapp.in/CRISApi/ws1/ntes/s3/trainServiceProfile"
	method := "POST"

	// payload := fmt.Sprintf(`{"srvid":"1989","trainNumber":"%s","startDate":"25-May-2025","tkn":"","lang":"en","language":"en","usrid":"","mode":"web","pltfrm":"ios","did":null,"deptid":"100014","formtrkr":"0","subsid":"0","subsid2":"0"}`, train_number)
	payload := fmt.Sprintf(`{"srvid":"1989","trainNumber":"%s","startDate":"%s","tkn":"","lang":"en","language":"en","usrid":"","mode":"web","pltfrm":"ios","did":null,"deptid":"100014","formtrkr":"0","subsid":"0","subsid2":"0"}`, t.getTrainNumber(), time.Now().Format("02-Jan-2006"))
	fmt.Println("Making request to fetch train service profile for train number:", t.getTrainNumber())
	fmt.Println(payload)

	payloadBody := bytes.NewBuffer([]byte(payload))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadBody)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("accept-language", "en-GB,en;q=0.9")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("deptid", "100014")
	req.Header.Add("formtrkr", "0")
	req.Header.Add("origin", "https://web.umang.gov.in")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://web.umang.gov.in/")
	req.Header.Add("sec-ch-ua", `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)
	req.Header.Add("sec-ch-ua-version", `"1.0"`)
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("srvid", "1989")
	req.Header.Add("subsid", "0")
	req.Header.Add("subsid2", "0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Add("x-api-key", getXApiKey())

	return t.FetchData(req, client)

}

func getLatLng(lat string, lng string) (*geoplot.LatLng, error) {
	latitude, err := strconv.ParseFloat(lat, 32)
	if err != nil {
		return nil, err
	}
	longitude, err := strconv.ParseFloat(lng, 32)
	if err != nil {
		return nil, err
	}
	return &geoplot.LatLng{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

func (t *TrainData) toRoute() types.Route {
	// _ = t.LoadData()
	return types.Route{
		RouteId:           t.getTrainNumber(),
		AgencyId:          "1",
		RouteShortName:    t.getTrainNumber(),
		RouteLongName:     strings.TrimSpace(t.TrainServiceProfileResponse.Pd.TrainServiceProfile.TrainName),
		RouteDesc:         "",
		RouteType:         types.RouteTypeRail,
		RouteUrl:          "",
		RouteColor:        "",
		RouteTextColor:    "",
		RouteSortOrder:    0,
		ContinuousPickup:  types.ContinuousPickupNotAvailable,
		ContinuousDropOff: types.ContinuousDropOffNotAvailable,
		CEMVSupport:       types.CEMVSupportUnknown,
	}
}

func (t *TrainData) toTrip() types.Trip {
	// _ = t.LoadData()
	if len(t.getRunningDays()) != 7 {
		fmt.Printf("Invalid running days string: %s for train %s \n", t.getRunningDays(), t.getTrainNumber())
	}

	return types.Trip{
		RouteId:              t.getTrainNumber(),
		ServiceId:            getServiceId(t.getRunningDays()),
		TripId:               t.getTrainNumber(),
		TripHeadsign:         strings.TrimSpace(t.TrainServiceProfileResponse.Pd.TrainServiceProfile.DestinationName),
		DirectionId:          types.DirectionIdOutbound,
		BlockId:              "",
		ShapeId:              t.getTrainNumber(),
		WheelchairAccessible: types.WheelChairAccessibilityNoInfo,
		BikesAllowed:         types.BikesAllowedNoInfo,
		CarsAllowed:          types.CarsAllowedNoInfo,
	}
}

func (t *TrainData) toStopTimes() ([]types.StopTime, []types.Shape) {
	// _ = t.LoadData()
	stopTimes := []types.StopTime{}
	shapes := []types.Shape{}
	prevDistance := -1.0
	for _, v := range t.TrainServiceProfileResponse.Pd.TrainServiceProfile.VTrainServiceSchedulePTT {
		arrivalTime := convertSecondsToHHMMSS(v.PttArrivalTimeInSecond)
		departureTime := convertSecondsToHHMMSS(v.PttDepartureTimeInSecond)
		if v.PttArrivalTimeInSecond == 0 && v.PttDepartureTimeInSecond == 0 {
			fmt.Printf("Skipping station %s for train %s as both arrival and departure time are 0\n", v.StationCode, t.getTrainNumber())
			// continue
		}
		if v.PttArrivalTimeInSecond == 0 && v.PttDepartureTimeInSecond == 0 {
			fmt.Printf("Skipping station %s for train %s as both arrival and departure time are 0\n", v.StationCode, t.getTrainNumber())
			// continue
		}
		if v.PttArrivalTimeInSecond == 0 {
			arrivalTime = departureTime
		} else if v.PttDepartureTimeInSecond == 0 {
			departureTime = arrivalTime
		}
		prevDistance = math.Max(prevDistance+0.1, float64(v.DistanceFromSource))
		stopTimes = append(stopTimes, types.StopTime{
			TripId:            t.getTrainNumber(),
			ArrivalTime:       arrivalTime,
			DepartureTime:     departureTime,
			StopId:            v.StationCode,
			StopSequence:      v.SerialNumber,
			PickupType:        types.PickupTypeRegularlyScheduled,
			DropOffType:       types.DropOffTypeRegularlyScheduled,
			ShapeDistTraveled: prevDistance,
		})
		newStation := models.Station{
			Code: v.StationCode,
			Name: v.StationName,
			Lat:  v.Lattitude,
			Lng:  v.Longitude,
		}
		// TODO: Remove after fixing station issues
		fixStation(&newStation, t.dataErrors)
		lat, _ := strconv.ParseFloat(strings.TrimSpace(newStation.Lat), 64)
		lng, _ := strconv.ParseFloat(strings.TrimSpace(newStation.Lng), 64)
		shapes = append(shapes, types.Shape{
			ShapeId:           t.getTrainNumber(),
			ShapePtLat:        lat,
			ShapePtLon:        lng,
			ShapePtSequence:   v.SerialNumber,
			ShapeDistTraveled: prevDistance,
		})
		if v.PttDepartureTimeInSecond == 0 {
			break
		}
	}
	return stopTimes, shapes
}

func convertSecondsToHHMMSS(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
