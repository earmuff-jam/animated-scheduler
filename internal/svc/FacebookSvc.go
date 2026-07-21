package svc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/earmuff-jam/scheduler/internal/types"
)

// PerformFacebookSvcHealthCheck ...
// perform health check for facebook pages
func PerformFacebookSvcHealthCheck(fb types.FacebookSvcData) (bool, error) {

	isValid, err := performFacebookHealthCheck(fb)
	if err != nil {
		log.Printf("unable to pass health check. error: %+v", err)
		return false, err
	}
	log.Println("Completed health check for facebook")
	return isValid, nil
}

// PerformUpdateToFacebookPage ...
// updates the facebook page with the new content
func PerformUpdateToFacebookPage(fb types.FacebookSvcData, data types.CSVRowData) (bool, error) {

	image := selectRandomImageForContent()
	isComplete, err := performPostToFacebook(fb, data, image)
	if err != nil {
		log.Printf("unable to update facebook page. error: %+v", err)
		return false, err
	}

	return isComplete, nil
}

// selectRandomImageForContent ...
// defines a function that returns random image url
func selectRandomImageForContent() string {

	entries, err := os.ReadDir(filepath.Join("content", "images"))
	if err != nil {
		log.Printf("unable to read image directory. error: %+v", err)
		return ""
	}

	entry := entries[rand.IntN(len(entries))]
	imagePath := filepath.Join("content", "images", entry.Name())

	return imagePath

}

// performFacebookHealthCheck ...
// defines a function that is used to check if the system is alive
func performFacebookHealthCheck(fb types.FacebookSvcData) (bool, error) {
	url := fmt.Sprintf(
		"%s/%s/settings?origin_graph_explorer=1&transport=cors&access_token=%s",
		fb.URI,
		fb.PageID,
		fb.PageToken,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("unable to reach destination. error: %+v", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("unable to perform health check. Details: %+v", string(respBody))
		return false, errors.New("unable to perform health check")
	}

	var result types.FacebookSvcSettingsDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("unable to decode response. error: %+v", err)
		return false, fmt.Errorf("unable to decode response: %+v", err)
	}

	log.Printf("Health check completed. Response: %+v", result)
	return true, nil
}

// performPostToFacebook ...
// defines a function that is used to create a post in facebook
func performPostToFacebook(fb types.FacebookSvcData, data types.CSVRowData, imagePath string) (bool, error) {
	requestURL := fmt.Sprintf(
		"%s/%s/photos?origin_graph_explorer=1&transport=cors&access_token=%s",
		fb.URI,
		fb.PageID,
		fb.PageToken,
	)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("unable to read image path. error: %+v", err)
		return false, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("source", filepath.Base(imagePath))
	if err != nil {
		log.Printf("unable to create form file with image. error: %+v", err)
		return false, err
	}

	if _, err := io.Copy(part, file); err != nil {
		log.Printf("unable to copy file. error: %+v", err)
		return false, err
	}

	writer.WriteField("message", data.Message)
	writer.WriteField("published", "false")
	writer.WriteField(
		"scheduled_publish_time", strconv.FormatInt(data.Date.Unix(), 10),
	)

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, requestURL, body)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("unable to send request parameters. error: %+v", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("unable to perform post. error: %+v", string(respBody))
		return false, errors.New(errorMsg)

	}

	return true, nil
}
