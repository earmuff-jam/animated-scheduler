package platform

import "github.com/earmuff-jam/scheduler/internal/types"

// Retrieve facebook pages via api - this is a health check

// If health check is g2g, create posts using facebook pages api

// Store response in file. include failing items.

// Retry failing items at a later date.

// If passed, remove the failed items.

// Send email to user with completed status

// Platform ...
// defines a interface to post to various locations
type Platform interface {
	PreProcessCSVData(row []string) (types.CSVRowData, error)
	ProcessCSVData(row types.CSVRowData) error
	PostProcessCSVData()
}
