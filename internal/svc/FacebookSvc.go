package svc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/earmuff-jam/scheduler/internal/types"
)

// PerformFacebookSvcHealthCheck ...
// perform health check for facebook pages
func PerformFacebookSvcHealthCheck(fb types.FacebookSvcData) (bool, error) {

	isValid, err := performHealthCheck(fb)
	if err != nil {
		log.Printf("unable to pass health check. error: %+v", err)
		return false, err
	}
	log.Println("Health check for Facebook passed")
	return isValid, nil
}

// PerformUpdateToFacebookPage ...
// updates the facebook page with the new content
func PerformUpdateToFacebookPage(fb types.FacebookSvcData, data types.CSVRowData) (bool, error) {

	isComplete, err := performPost(fb, data)
	if err != nil {
		log.Printf("unable to update facebook page. error: %+v", err)
		return false, err
	}

	return isComplete, nil
}

// performHealthCheck ...
// defines a function that is used to check if the system is alive
func performHealthCheck(fb types.FacebookSvcData) (bool, error) {
	url := fmt.Sprintf(
		"%s/%s/settings?access_token=%s",
		fb.URI,
		fb.PageID,
		fb.PageToken,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("unable to reach destination. error: %+v", err)
		return false, err
	}

	return resp.StatusCode == 200, nil
}

func performPost(fb types.FacebookSvcData, data types.CSVRowData) (bool, error) {
	url := fmt.Sprintf(
		"%s/%s/feed?access_token=%s",
		fb.URI,
		fb.PageID,
		fb.PageToken,
	)

	payload := map[string]string{
		"messsage": data.Message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("unable to marshall payload into required body. error: %+v", err)
		return false, err
	}

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Printf("unable to reach destination. error: %+v", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("unable to perform post. error: %+v", string(respBody))
		return false, errors.New(errorMsg)

	}

	return false, errors.New("unable to post in facebook")
}
