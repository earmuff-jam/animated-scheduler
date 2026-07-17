package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// Read env files
	err := godotenv.Load()
	if err != nil {
		log.Printf("unable to load env file. error: %+v", err)
		return
	}

	// Read CSV files
	fileNameInEnv := os.Getenv("CONTENT_FILENAME")

	if fileNameInEnv == "" {
		log.Println("unable to read content with invalid filename. error.")
		return

	}

	csvContentFile, err := os.Open(fileNameInEnv)
	if err != nil {
		log.Printf("unable to open csv fields from provided file. error: %+v", err)
		return
	}
	defer csvContentFile.Close()

	reader := csv.NewReader(csvContentFile)

	for {
		row, err := reader.Read()
		if err != io.EOF {
			log.Println("reached end of current file.")
			break // end of file
		}
		if err != nil {
			log.Printf("unable to read current row within file. error: %+v", err)
			return
		}
		fmt.Println(row)
	}

	// Retrieve facebook pages via api - this is a health check

	// If health check is g2g, create posts using facebook pages api

	// Store response in file. include failing items.

	// Retry failing items at a later date.

	// If passed, remove the failed items.

	// Send email to user with completed status

}
