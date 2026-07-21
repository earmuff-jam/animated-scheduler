package platform

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/earmuff-jam/scheduler/internal/svc"
	"github.com/earmuff-jam/scheduler/internal/types"
)

// Instagram ...
type Instagram struct {
	URI    string
	PageID string
	Token  string
}

// NewInstagram ...
// defines a function that handles instagram posts
func NewInstagram(URI string, pageID string, token string) *Instagram {
	return &Instagram{
		URI:    URI,
		PageID: pageID,
		Token:  token,
	}
}

// PreProcessCSVData ...
// defines a function that is used to pre process CSV data
func (in *Instagram) PreProcessCSVData(data []string) (types.CSVRowData, error) {

	parsedCSV, err := parseCSVToInstagramStruct(data)
	if err != nil {
		log.Printf("unable to parse csv data. error: %+v", err)
		return types.CSVRowData{}, err
	}

	return parsedCSV, nil
}

// ProcessCSVData ...
// defines a function that is used to process CSV data
func (in *Instagram) ProcessCSVData(parsedCSV types.CSVRowData) error {

	svcInstagramClient := types.InstagramSvcData{
		URI:       in.URI,
		PageID:    in.PageID,
		PageToken: in.Token,
	}

	businessID, err := svc.PerformInstagramHealthCheck(svcInstagramClient)

	if err != nil {
		log.Printf("unable to perform service health check. error %+v", err)
		return err
	}

	if businessID == "" {
		log.Printf("health check failed for instagram. Ignoring updates.")
		return errors.New("health check failed for instagram.")
	}

	log.Println("instagram health check is successful.")

	isUpdateComplete, err := svc.PerformPostToInstagram(svcInstagramClient, parsedCSV)
	if err != nil {
		log.Printf("unable to update instagram pages. error: %+v", err)
		return err
	}

	log.Printf("Updated selected row. Details: %+v", isUpdateComplete)
	return errors.New("unable to publish into instagram")

}

// parseCSVToInstagramStruct ...
// parse CSV for each row
func parseCSVToInstagramStruct(data []string) (types.CSVRowData, error) {

	formattedCSVRowData := types.CSVRowData{}

	date, err := time.Parse(time.RFC3339, data[0])
	if err != nil {
		log.Printf("unable to parse date for row. error: %+v", err)
		return formattedCSVRowData, err
	}

	isComplete, err := strconv.ParseBool(data[4])
	if err != nil {
		log.Printf("unable to parse isComplete for row. error: %+v", err)
		return formattedCSVRowData, err
	}

	retryCount, err := strconv.Atoi(data[5])
	if err != nil {
		log.Printf("unable to parse retryCount for row. error: %+v", err)
		return formattedCSVRowData, err
	}

	return types.CSVRowData{
		Date:       &date,
		Message:    data[1],
		ImageURL:   data[2],
		Status:     data[3],
		IsComplete: isComplete,
		RetryCount: retryCount,
	}, nil

}
