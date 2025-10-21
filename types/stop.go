package types

import (
	"fmt"
	"strconv"
)

type Stop struct {
	StopId             string             `json:"stop_id" csv:"stop_id"`
	StopCode           string             `json:"stop_code" csv:"stop_code"`
	StopName           string             `json:"stop_name" csv:"stop_name"`
	TTSStopName        string             `json:"tts_stop_name" csv:"tts_stop_name"`
	StopDesc           string             `json:"stop_desc" csv:"stop_desc"`
	StopLat            float64            `json:"stop_lat" csv:"stop_lat"`
	StopLon            float64            `json:"stop_lon" csv:"stop_lon"`
	ZoneId             string             `json:"zone_id" csv:"zone_id"`
	StopUrl            string             `json:"stop_url" csv:"stop_url"`
	LocationType       LocationType       `json:"location_type" csv:"location_type"`
	ParentStation      string             `json:"parent_station" csv:"parent_station"`
	StopTimezone       string             `json:"stop_timezone" csv:"stop_timezone"`
	WheelchairBoarding WheelChairBoarding `json:"wheelchair_boarding" csv:"wheelchair_boarding"`
	LevelId            string             `json:"level_id" csv:"level_id"`
	PlatformCode       string             `json:"platform_code" csv:"platform_code"`
	StopAccess         StopAccess         `json:"stop_access" csv:"stop_access"`
}

type LocationType int // Define a custom type for clarity and type safety

const (
	LocationTypeStop LocationType = iota
	LocationTypeStation
	LocationTypeEntranceOrExit
	LocationTypeGenericNode
	LocationTypeBoardingArea
)

type WheelChairBoarding int // Define a custom type for clarity and type safety
const (
	WheelChairBoardingNoInfo WheelChairBoarding = iota
	WheelChairBoardingAccessible
	WheelChairBoardingNotAccessible
)

type StopAccess int

const (
	StopAccessNotDirectlyAccessible StopAccess = iota
	StopAccessDirectlyAccessible    StopAccess = iota
	StopAccessUnknown               StopAccess = iota
)

func (s *StopAccess) MarshalCSV() (string, error) {
	if *s == StopAccessUnknown {
		return "", nil
	}
	return strconv.Itoa(int(*s)), nil
}

func (s *StopAccess) UnmarshalCSV(csv string) (err error) {
	support, err := strconv.Atoi(csv)
	if err != nil {
		return err
	}
	switch support {
	case 0:
		*s = StopAccessNotDirectlyAccessible
	case 1:
		*s = StopAccessDirectlyAccessible
	case 2:
		*s = StopAccessUnknown
	default:
		return fmt.Errorf("invalid SupportAccess value: %d", support)
	}
	return nil
}
