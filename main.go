// package main

// import (
// 	"encoding/csv"
// 	"io"
// 	"log"
// 	"os"

// 	"github.com/earmuff-jam/scheduler/internal/platform"
// 	"github.com/earmuff-jam/scheduler/internal/scheduler"
// 	"github.com/joho/godotenv"
// )

// // FileNameInEnv ...
// // the file name of the content in csv
// const FileNameInEnv = "CONTENT_FILENAME"

// // FacebookPageId ...
// // the facebook page identity
// const FacebookPageId = "FACEBOOK_PAGE_ID"

// // FacebookPageUri ...
// // the facebook page URI
// const FacebookPageUri = "FACEBOOK_URI"

// // FacebookPagesAccessToken ...
// // the facebook page access token
// const FacebookPagesAccessToken = "FACEBOOK_PAGE_ACCESS_TOKEN"

// func main() {

// 	// Read env files
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Printf("unable to load env file. error: %+v", err)
// 		return
// 	}

// 	fileNameInEnv := os.Getenv(FileNameInEnv)

// 	if fileNameInEnv == "" {
// 		log.Println("unable to read content with invalid filename. error.")
// 		return

// 	}

// 	csvContentFile, err := os.Open(fileNameInEnv)
// 	if err != nil {
// 		log.Printf("unable to open csv fields from provided file. error: %+v", err)
// 		return
// 	}
// 	defer csvContentFile.Close()

// 	reader := csv.NewReader(csvContentFile)

// 	// facebook := platform.NewFacebook(os.Getenv(FacebookPageUri), os.Getenv(FacebookPageId), os.Getenv(FacebookPagesAccessToken))

// 	instagram := platform.NewInstagram(os.Getenv(FacebookPageUri), os.Getenv(FacebookPageId), os.Getenv(FacebookPagesAccessToken))

// 	scheduler := scheduler.Scheduler{Platforms: []platform.Platform{instagram}}

// 	// Skip CSV header
// 	_, err = reader.Read()
// 	if err != nil {
// 		log.Printf("unable to read csv header. error: %+v", err)
// 		return
// 	}

// 	for {
// 		post, err := reader.Read()
// 		if err == io.EOF {
// 			log.Println("reached end of current file.")
// 			break // end of file
// 		}
// 		if err != nil {
// 			log.Printf("unable to read current row within file. error: %+v", err)
// 			return
// 		}

// 		// reduce complexity of large files, hence send 1 row at a time
// 		scheduler.ProcessPost(post)
// 	}

// }

package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello World!",
	}, nil
}

func main() {
	// Initiate AWS Lambda handler
	lambda.Start(Handler)
}
