package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/earmuff-jam/scheduler/internal/platform"
	"github.com/earmuff-jam/scheduler/internal/scheduler"
	"github.com/joho/godotenv"
)

// FileNameInEnv ...
// the file name of the content in csv
const FileNameInEnv = "CONTENT_FILENAME"

// FacebookPageId ...
// the facebook page identity
const FacebookPageId = "FACEBOOK_PAGE_ID"

// FacebookPageUri ...
// the facebook page URI
const FacebookPageUri = "FACEBOOK_URI"

// FacebookPagesAccessToken ...
// the facebook page access token
const FacebookPagesAccessToken = "FACEBOOK_PAGE_ACCESS_TOKEN"

// Handler ...
// defines the handler function for lamdba to execute schedulers
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("starting scheduler")

	// Read env files
	err := godotenv.Load()
	if err != nil {
		log.Printf("unable to load env file. error: %+v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "unable to load env file",
		}, nil
	}

	fileNameInEnv := os.Getenv(FileNameInEnv)

	if fileNameInEnv == "" {
		log.Println("unable to read content with invalid filename.")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "unable to read content with invalid filename",
		}, nil
	}

	csvContentFile, err := os.Open(fileNameInEnv)
	if err != nil {
		log.Printf("unable to open csv fields from provided file. error: %+v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "unable to open csv fields from provided file",
		}, nil
	}
	defer csvContentFile.Close()

	reader := csv.NewReader(csvContentFile)

	facebook := platform.NewFacebook(os.Getenv(FacebookPageUri), os.Getenv(FacebookPageId), os.Getenv(FacebookPagesAccessToken))

	instagram := platform.NewInstagram(os.Getenv(FacebookPageUri), os.Getenv(FacebookPageId), os.Getenv(FacebookPagesAccessToken))

	scheduler := scheduler.Scheduler{
		Platforms: []platform.Platform{facebook, instagram}}

	// Skip CSV header
	_, err = reader.Read()
	if err != nil {
		log.Printf("unable to read csv header. error: %+v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "unable to read csv header",
		}, nil
	}

	for {
		post, err := reader.Read()
		if err == io.EOF {
			log.Println("reached end of current file.")
			break // end of file
		}
		if err != nil {
			log.Printf("unable to read current row within file. error: %+v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "unable to read current row within file",
			}, nil
		}

		fmt.Println(post)
		// reduce complexity of large files, hence send 1 row at a time
		scheduler.ProcessPost(post)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Schedule complete",
	}, nil

}
