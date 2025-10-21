package main

import (
	"fmt"
	"time"
)

var trainListTriggerTime = time.Now()
var trainDataTriggerTime = time.Now()

func throttleTrainListRequests() {
	if !time.Now().After(trainListTriggerTime.Add(5 * time.Second)) {
		sleepDuration := trainListTriggerTime.Add(5 * time.Second).Sub(time.Now())
		fmt.Printf("Throttling request for train searches, sleeping for %v seconds...\n", sleepDuration.Seconds())
		time.Sleep(sleepDuration)
		trainListTriggerTime = time.Now()
	}
}

func throttleTrainDataRequests() {
	if !time.Now().After(trainDataTriggerTime.Add(5 * time.Second)) {
		sleepDuration := trainDataTriggerTime.Add(5 * time.Second).Sub(time.Now())
		fmt.Printf("Throttling request for train data, sleeping for %v seconds...\n", sleepDuration.Seconds())
		time.Sleep(sleepDuration)
		trainDataTriggerTime = time.Now()
	}
}
