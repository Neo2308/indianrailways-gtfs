package models

type pd struct {
	RestServiceMessage  interface{}         `json:"restServiceMessage"`
	TrainServiceProfile TrainServiceProfile `json:"trainServiceProfile"`
	VTrainList          []vTrainList        `json:"vTrainList"`
}

type vTrainList struct {
	TrainNumber string `json:"trainNo"`
	TrainName   string `json:"trainName"`
	TrainType   string `json:"trainType"`
}
