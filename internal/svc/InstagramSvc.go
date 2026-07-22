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

// PerformInstagramHealthCheck ...
// defines a function that is used to perform health check for instagram
func PerformInstagramHealthCheck(in types.InstagramSvcData) (string, error) {

	businessID, err := performInstagramHealthCheck(in)
	if err != nil {
		log.Printf("unable to pass health check. error: %+v", err)
		return "", err
	}
	log.Println("Completed health check for facebook")
	return businessID, nil
}

// PerformPostToInstagram ...
// defines a function used to post content for instagram
func PerformPostToInstagram(in types.InstagramSvcData, businessID string, data types.CSVRowData) (bool, error) {

	// fetch random image
	publicImageURL, err := fetchRandomImage()
	if err != nil {
		log.Printf("unable to fetch random image url. error: %+v", err)
		return false, err
	}

	data.ImageURL = publicImageURL

	instagramMediaContainer, err := createInstagramMediaContainer(businessID, in, data)
	if err != nil {
		log.Printf("unable to create instagram media container. error: %+v", err)
		return false, err
	}

	// post accepted media container
	postInstagramFromMediaContainer(instagramMediaContainer.ID, in, businessID)

	return false, errors.New("unable to perform post for instagram.")

}

// fetchRandomImage ...
// defines a function that retrieves a random image
func fetchRandomImage() (string, error) {

	resp, err := http.Get("https://picsum.photos/1200/1200")
	if err != nil {
		log.Printf("unable to fetch public image url. error: %+v", err)
		return "", err
	}
	defer resp.Body.Close()

	finalImageURL := resp.Request.URL.String()
	return finalImageURL, nil
}

// postInstagramFromMediaContainer ...
// defines a function that posts to instagram from media container
func postInstagramFromMediaContainer(containerID string, in types.InstagramSvcData, businessID string) (bool, error) {

	if containerID == "" {
		log.Printf("unable to post without a container id")
		return false, errors.New("unable to post without a media container id")
	}

	publishURL := fmt.Sprintf(
		"%s/%s/media_publish?access_token=%s",
		in.URI,
		businessID,
		in.PageToken,
	)

	publishPayload := map[string]string{
		"creation_id": containerID,
	}

	publishBody, err := json.Marshal(publishPayload)
	if err != nil {
		log.Printf("unable to marshall publish payload. error: %+v", err)
		return false, err
	}

	publishReq, err := http.NewRequest(
		http.MethodPost,
		publishURL,
		bytes.NewBuffer(publishBody),
	)

	if err != nil {
		log.Printf("unable to create publish request. error: %+v", err)
		return false, err
	}

	publishReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	publishResp, err := client.Do(publishReq)
	if err != nil {
		log.Printf("unable to publish to instagram media. error: %+v", err)
		return false, err
	}
	defer publishResp.Body.Close()

	publishResponseBody, err := io.ReadAll(publishResp.Body)
	if err != nil {
		log.Printf("unable to read publish response. error: %+v", err)
		return false, err
	}

	if publishResp.StatusCode == http.StatusBadRequest {
		log.Println("unable to publish data in media container for instagram")

		errorMsg := fmt.Sprintf(
			"instagram media creation failed. status: %d response: %s",
			publishResp.StatusCode,
			string(publishResponseBody))

		return false, errors.New(errorMsg)
	}

	return true, nil
}

// createInstagramMediaContainer ...
// defines a function that creates media container for instagram
func createInstagramMediaContainer(businessID string, in types.InstagramSvcData, data types.CSVRowData) (types.InstagramSvcMediaContainerResponse, error) {

	mediaURL := fmt.Sprintf(
		"%s/%s/media?access_token=%s",
		in.URI,
		businessID,
		in.PageToken,
	)

	payload := map[string]string{
		"image_url": data.ImageURL,
		"caption":   data.Message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("unable to marshall instagram payload. error: %+v", err)
		return types.InstagramSvcMediaContainerResponse{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		mediaURL,
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Printf("unable to create media container. error: %+v", err)
		return types.InstagramSvcMediaContainerResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("unable to perform create of instagram media container. error: %+v", err)
		return types.InstagramSvcMediaContainerResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("unable to read instagram response. error: %+v", err)
		return types.InstagramSvcMediaContainerResponse{}, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		log.Println("unable to create media container for instagram")
		errorMsg := fmt.Sprintf(
			"instagram media creation failed. status: %d response: %s",
			resp.StatusCode,
			string(responseBody),
		)
		return types.InstagramSvcMediaContainerResponse{}, errors.New(errorMsg)
	}

	var result types.InstagramSvcMediaContainerResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Printf("unable to decode response. error: %+v", err)
		return types.InstagramSvcMediaContainerResponse{}, err
	}

	log.Println("created media container successfully")
	return result, nil
}

// performInstagramHealthCheck ...
// defines a function that performs health check for instagram
func performInstagramHealthCheck(in types.InstagramSvcData) (string, error) {

	url := fmt.Sprintf(
		"%s/%s?fields=instagram_business_account&access_token=%s",
		in.URI,
		in.PageID,
		in.PageToken,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("unable to reach destination. error: %+v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("unable to perform health check. Details: %+v", string(respBody))
		return "", errors.New("unable to perform health check")
	}

	var result types.InstagramSvcHealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("unable to decode response. error: %+v", err)
		return "", fmt.Errorf("unable to decode response: %+v", err)
	}

	log.Printf("Health check completed. Response: %+v", result)
	return result.BusinessAccount.BusinessID, nil
}
