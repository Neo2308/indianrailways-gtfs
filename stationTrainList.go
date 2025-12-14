package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Neo2308/indianrailways-gtfs/apiDataFetcher"
	"github.com/Neo2308/indianrailways-gtfs/models"
)

type LiveStationData struct {
	apiDataFetcher.ApiDataFetcher[models.TrainServiceProfileResponse]
	stationCode                 string
	TrainServiceProfileResponse models.TrainServiceProfileResponse
}

func NewLiveStationData(stationCode string) *LiveStationData {
	newObj := &LiveStationData{
		stationCode: stationCode,
	}
	newObj.ApiDataFetcher = *apiDataFetcher.NewApiDataFetcher[models.TrainServiceProfileResponse](
		&newObj.TrainServiceProfileResponse,
		"live station",
		fmt.Sprintf("station/%s.json", stationCode),
		newObj.getTrainSearchesUncached,
	)
	return newObj
}

func (t *LiveStationData) getTrains() []int {
	_ = t.LoadData()
	trainNumbers := []int{}
	// fmt.Println(t.TrainServiceProfileResponse)
	for _, v := range t.TrainServiceProfileResponse.Pd.VTrainList {
		trainNumber, _ := strconv.Atoi(v.TrainNumber)
		// fmt.Println("Found train:", v.TrainNumber, v.TrainName, " in search")
		trainNumbers = append(trainNumbers, trainNumber)
	}
	return trainNumbers
}

func (t *LiveStationData) getTrainSearchesUncached() error {
	url := "https://apigw.umangapp.in/CRISApi/ws1/ntes/s2/liveStation"
	method := "POST"

	// payload := fmt.Sprintf(`{"srvid":"1989","trainNumber":"%s","startDate":"25-May-2025","tkn":"","lang":"en","language":"en","usrid":"","mode":"web","pltfrm":"ios","did":null,"deptid":"100014","formtrkr":"0","subsid":"0","subsid2":"0"}`, train_number)
	payload := fmt.Sprintf(`{"srvid":"1988","stationCode":"%s","goingTo":"","nextMins":"480","tkn":"","lang":"en","language":"en","usrid":"","mode":"web","pltfrm":"ios","did":null,"deptid":"100014","formtrkr":"0","subsid":"0","subsid2":"0"}`, t.stationCode)
	fmt.Println("Making request to fetch trains for the station:", t.stationCode)
	fmt.Println(payload)

	payloadBody := bytes.NewBuffer([]byte(payload))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

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
