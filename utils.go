package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Location int // Define a custom type for clarity and type safety

var XApiKey = "unset"

const (
	INPUT Location = iota
	OUTPUT
	TEMP
	CACHE
)

const (
	INPUT_DATA_LOCATION  = "input-data"
	OUTPUT_DATA_LOCATION = "gtfs"
	TEMP_DATA_LOCATION   = "temp-data"
	CACHE_DATA_LOCATION  = "cache-data"
)

var locations = map[Location]string{
	INPUT:  INPUT_DATA_LOCATION + "/%s",
	OUTPUT: OUTPUT_DATA_LOCATION + "/%s",
	TEMP:   TEMP_DATA_LOCATION + "/%s",
	CACHE:  CACHE_DATA_LOCATION + "/%s",
}

// Loads data from a file.
// file_name The name of the file to load data from.
// loc The location where the file should be loaded from (INPUT, OUTPUT, TEMP, or CACHE).
// :return: The data loaded from the file.
func loadFile(fileName string, loc Location) ([]byte, error) {
	jsonFile, err := os.Open(getFilePath(fileName, loc))
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	//fmt.Println("Successfully opened file ", fileName)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	fileBytes, _ := io.ReadAll(jsonFile)
	return fileBytes, nil
}

func getFilePath(fileName string, loc Location) string {
	return fmt.Sprintf(locations[loc], fileName)
}

func saveFile(fileName string, data []byte, loc Location) error {
	// Ensure the directory exists
	dirs := ""
	if strings.Contains(fileName, "/") {
		dirs = fileName[:strings.LastIndex(fileName, "/")]
	}
	err := os.MkdirAll(getFilePath(dirs, loc), os.ModePerm)
	if err != nil {
		return err
	}
	// Write the data to the file
	return os.WriteFile(getFilePath(fileName, loc), data, 0644)
}

func getXApiKey() string {
	if XApiKey == "unset" {
		XApiKey = os.Getenv("X_API_KEY")
	}
	if XApiKey == "" {
		fmt.Println("Warning: X_API_KEY environment variable is not set.")
		panic("X_API_KEY environment variable is required")
	}
	return XApiKey
}
