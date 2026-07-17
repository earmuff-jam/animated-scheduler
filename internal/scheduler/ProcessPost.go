package scheduler

import (
	"log"

	"github.com/earmuff-jam/scheduler/internal/platform"
)

type Scheduler struct {
	Platforms []platform.Platform
}

// ProcessPost ...
// defines a function that directs csv data to various outlets to process data
func (sc *Scheduler) ProcessPost(csvData []string) {
	for _, platform := range sc.Platforms {
		csvRowData, err := platform.PreProcessCSVData(csvData)
		err = platform.ProcessCSVData(csvRowData)
		if err != nil {
			log.Printf("unable to post in platform %+s. error: %+v", platform, err)
			continue
		}
	}
}
