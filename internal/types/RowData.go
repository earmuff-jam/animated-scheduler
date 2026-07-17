package types

import "time"

// CSVRowData ...
// defines row for each CSV
type CSVRowData struct {
	Date           *time.Time
	Message        string
	ImageURL       string
	Status         string
	FacebookPostID string
	IsComplete     bool
	RetryCount     int
}

// FacebookSvcData ...
// defines facebook data struct
type FacebookSvcData struct {
	URI       string
	PageID    string
	PageToken string
}
