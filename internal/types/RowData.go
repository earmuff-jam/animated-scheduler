package types

import (
	"encoding/json"
	"time"
)

// CSVRowData ...
// defines row for each CSV
type CSVRowData struct {
	Date       *time.Time
	Message    string
	ImageURL   string
	Status     string
	IsComplete bool
	RetryCount int
}

// InstagramSvcData ...
// defines instagram data struct
type InstagramSvcData struct {
	URI       string
	PageID    string
	PageToken string
}

// InstagramSvcMediaContainerResponse ...
// defines a response for instagram media container
type InstagramSvcMediaContainerResponse struct {
	ID string `json:"id"`
}

// InstagramSvcHealthResponse ...
// defines a health repsonse struct for Instagram
type InstagramSvcHealthResponse struct {
	InstagramBusinessAccountID string `json:"id"`
}

// FacebookSvcData ...
// defines facebook data struct
type FacebookSvcData struct {
	URI       string
	PageID    string
	PageToken string
}

// FacebookSvcSettingsData ...
// defines the struct for each data row
type FacebookSvcSettingsData struct {
	Value   json.RawMessage `json:"value"` // responses vary - string || bool
	Setting string          `json:"setting"`
}

// FacebookSvcSettingsDataResponse ...
// defines the api response
type FacebookSvcSettingsDataResponse struct {
	Data []FacebookSvcSettingsData `json:"data"`
}
