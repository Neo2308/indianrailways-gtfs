package types

import (
	"fmt"
	"strconv"
)

type Agency struct {
	AgencyId       string      `json:"agency_id" csv:"agency_id"`
	AgencyName     string      `json:"agency_name" csv:"agency_name"`
	AgencyUrl      string      `json:"agency_url" csv:"agency_url"`
	AgencyTimezone string      `json:"agency_timezone" csv:"agency_timezone"`
	AgencyLang     string      `json:"agency_lang" csv:"agency_lang"`
	AgencyPhone    string      `json:"agency_phone" csv:"agency_phone"`
	AgencyFareUrl  string      `json:"agency_fare_url" csv:"agency_fare_url"`
	AgencyEmail    string      `json:"agency_email" csv:"agency_email"`
	CEMVSupport    CEMVSupport `json:"cemv_support" csv:"cemv_support"`
}

type CEMVSupport int // Define a custom type for clarity and type safety

func (c *CEMVSupport) MarshalCSV() (string, error) {
	return strconv.Itoa(int(*c)), nil
}

func (c *CEMVSupport) UnmarshalCSV(csv string) (err error) {
	support, err := strconv.Atoi(csv)
	if err != nil {
		return err
	}
	switch support {
	case 0:
		*c = CEMVSupportUnknown
	case 1:
		*c = CEMVSupportSupported
	case 2:
		*c = CEMVSupportUnsupported
	default:
		return fmt.Errorf("invalid CEMVSupport value: %d", support)
	}
	return nil
}

const (
	CEMVSupportUnknown CEMVSupport = iota
	CEMVSupportSupported
	CEMVSupportUnsupported
)
