package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var throttlingInterval = 7 * time.Second // seconds
var throttleTriggerTime = map[string]time.Time{}

type ApiDataFetcher[T any] struct {
	DataLocation  *T
	UrlType       string
	CacheFileName string
	RefetchFunc   func() error
}

func newApiDataFetcher[T any](dataLocation *T, urlType string, cacheFileName string, refetchFunc func() error) *ApiDataFetcher[T] {
	return &ApiDataFetcher[T]{
		DataLocation:  dataLocation,
		UrlType:       urlType,
		CacheFileName: cacheFileName,
		RefetchFunc:   refetchFunc,
	}
}

func (a *ApiDataFetcher[T]) LoadData() error {
	// file_name := fmt.Sprintf("%s.json", fmt.Sprintf("%d", t.trainNumber))
	// file_name := '{}.json'.format(train_number)
	// Load json from cache if available
	data, _ := loadFile(a.CacheFileName, CACHE)
	return json.Unmarshal(data, a.DataLocation)
}

func (a *ApiDataFetcher[T]) populateData() error {
	// file_name := fmt.Sprintf("searches/%s.json", searchString)
	// file_name := '{}.json'.format(train_number)
	// Load json from cache if available
	data, err := loadFile(a.CacheFileName, CACHE)
	if err == nil {
		// fmt.Printf("Cache hit for %s, loading from cache...\n", a.CacheFileName)
		// fmt.Println(string(data))
		return json.Unmarshal(data, &struct{}{})
	}
	if !os.IsNotExist(err) {
		return err
	}
	fmt.Printf("Cache miss for %s, fetching from API...\n", a.CacheFileName)
	return a.RefetchFunc()
}

func (a *ApiDataFetcher[T]) fetchData(req *http.Request, client *http.Client) error {
	a.throttleRequests()
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 504 {
		fmt.Println("Request timed out, retrying...")
		res, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	if res.StatusCode == 200 {
		var result map[string]interface{}
		json.NewDecoder(res.Body).Decode(&result)
		defer res.Body.Close()
		responseJson, _ := json.MarshalIndent(result, "", "    ")
		return saveFile(a.CacheFileName, responseJson, CACHE)
	}
	return fmt.Errorf("failed to fetch %s data, %d - %s", a.UrlType, res.StatusCode, res.Status)
}

func (a *ApiDataFetcher[T]) throttleRequests() {
	if _, ok := throttleTriggerTime[a.UrlType]; !ok {
		throttleTriggerTime[a.UrlType] = time.Now()
		return
	}
	fmt.Printf("Checking if request for %s should be throttled at %v \n", a.UrlType, time.Now())
	// fmt.Printf("Checking this %v %v %v\n", a.UrlType, time.Now(), throttleTriggerTime[a.UrlType], throttleTriggerTime[a.UrlType].Add(throttlingInterval))
	if !time.Now().After(throttleTriggerTime[a.UrlType].Add(throttlingInterval)) {
		sleepDuration := throttleTriggerTime[a.UrlType].Add(throttlingInterval).Sub(time.Now())
		fmt.Printf("Throttling request for %s, sleeping for %v seconds...\n", a.UrlType, sleepDuration.Seconds())
		time.Sleep(sleepDuration)
	}
	throttleTriggerTime[a.UrlType] = time.Now()
}
