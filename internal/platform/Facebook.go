package platform

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/earmuff-jam/scheduler/internal/svc"
	"github.com/earmuff-jam/scheduler/internal/types"
)

type Facebook struct {
	URI    string
	PageID string
	Token  string
}

// Facebook ...
// defines a function to handle facebook posts
func NewFacebook(URI string, pageID string, token string) *Facebook {
	return &Facebook{
		URI:    URI,
		PageID: pageID,
		Token:  token,
	}
}

// PreProcessCSVData ...
// formats the data into CSV ready struct
func (f *Facebook) PreProcessCSVData(data []string) (types.CSVRowData, error) {

	parsedCSV, err := parseCSV(data)
	if err != nil {
		log.Printf("unable to parse csv data. error: %+v", err)
		return parsedCSV, nil
	}

	return types.CSVRowData{}, errors.New("unable to parse csv data")

}

// ProcessCSVData ...
// defines a function that is used to process csv data
func (f *Facebook) ProcessCSVData(parsedCSV types.CSVRowData) error {

	svcFacebookClient := types.FacebookSvcData{
		URI:       f.URI,
		PageID:    f.PageID,
		PageToken: f.Token,
	}

	isValid, err := svc.PerformFacebookSvcHealthCheck(svcFacebookClient)
	if err != nil {
		log.Printf("unable to perform service health check. error %+v", err)
		return err
	}

	log.Printf("health check is successful. Status: %+v", isValid)
	isUpdateComplete, err := svc.PerformUpdateToFacebookPage(svcFacebookClient, parsedCSV)
	if err != nil {
		log.Printf("unable to update facebook pages. error: %+v", err)
		return err
	}

	log.Printf("Updated selected row. Details: %+v", isUpdateComplete)
	return errors.New("unable to publish into facebook")

}

// PostProcessCSVData ...
// defines a function that is used to post process csv data
func (f *Facebook) PostProcessCSVData() {
	// do post processing later
}

// parseCSV ...
// parse CSV for each row
func parseCSV(data []string) (types.CSVRowData, error) {

	formattedCSVRowData := types.CSVRowData{}

	date, err := time.Parse(time.RFC3339, data[0])
	if err != nil {
		log.Printf("unable to parse date for row. error: %+v", err)
		return formattedCSVRowData, err
	}

	isComplete, err := strconv.ParseBool(data[5])
	if err != nil {
		log.Printf("unable to parse isComplete for row. error: %+v", err)
		return formattedCSVRowData, err
	}

	retryCount, err := strconv.Atoi(data[6])
	if err != nil {
		log.Printf("unable to parse retryCount for row. error: %+v", err)
		return formattedCSVRowData, err
	}

	return types.CSVRowData{
		Date:           &date,
		Message:        data[1],
		ImageURL:       data[2],
		Status:         data[3],
		FacebookPostID: data[4],
		IsComplete:     isComplete,
		RetryCount:     retryCount,
	}, nil

}
