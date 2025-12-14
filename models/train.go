package models

type Station struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
}

type TrainServiceProfileResponse struct {
	Pd pd     `json:"pd"`
	Rc string `json:"rc"`
	Rd string `json:"rd"`
	Rs string `json:"rs"`
}

type TrainServiceProfile struct {
	VTrainServiceSchedulePTT   []vTrainServiceSchedulePTT `json:"vTrainServiceSchedulePTT"`
	VTrainServiceScheduleWTT   []vTrainServiceScheduleWTT `json:"vTrainServiceScheduleWTT"`
	DaysOfRunFromSourceNumeric string                     `json:"daysOfRunFromSourceNumeric"`
	TrainName                  string                     `json:"trainName"`
	DestinationName            string                     `json:"destinationName"`
}

type vTrainServiceSchedulePTT struct {
	SerialNumber             int    `json:"serialNumber"`
	StationName              string `json:"stationName"`
	StationCode              string `json:"station"`
	Lattitude                string `json:"lattitude"`
	Longitude                string `json:"longitude"`
	PttArrivalTimeInSecond   int    `json:"pttArrivalTimeInSecond"`
	PttDepartureTimeInSecond int    `json:"pttDepartureTimeInSecond"`
	DistanceFromSource       int    `json:"distanceFromSource"`
}

type vTrainServiceScheduleWTT struct {
	SerialNumber int    `json:"serialNumber"`
	StationCode  string `json:"station"`
	Lattitude    string `json:"lattitude"`
	Longitude    string `json:"longitude"`
}
