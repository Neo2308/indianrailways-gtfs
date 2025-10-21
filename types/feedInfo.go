package types

import "time"

type FeedInfo struct {
	FeedPublisherName string `json:"feed_publisher_name" csv:"feed_publisher_name"`
	FeedPublisherUrl  string `json:"feed_publisher_url" csv:"feed_publisher_url"`
	FeedLang          string `json:"feed_lang" csv:"feed_lang"`
	FeedStartDate     Date   `json:"feed_start_date" csv:"feed_start_date"`
	FeedEndDate       Date   `json:"feed_end_date" csv:"feed_end_date"`
	FeedVersion       string `json:"feed_version" csv:"feed_version"`
	FeedContactEmail  string `json:"feed_contact_email" csv:"feed_contact_email"`
	FeedContactUrl    string `json:"feed_contact_url" csv:"feed_contact_url"`
}

type Date struct {
	time.Time
}

func (date *Date) MarshalCSV() (string, error) {
	return date.Time.Format("20060102"), nil
}

// You could also use the standard Stringer interface
func (date *Date) String() string {
	return date.String() // Redundant, just for example
}

// UnmarshalCSV Convert the CSV string as internal date
func (date *Date) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("20060102", csv)
	return err
}
