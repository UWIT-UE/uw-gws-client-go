package gws

import (
	"fmt"
	"strings"
)

// apiError describes an error returned by the API.
type apiError struct {
	Status    int      `json:"status"`
	SubStatus string   `json:"subStatus"`
	Detail    []string `json:"detail"`
	// NotFound only returned on membership PUT
	NotFound []string `json:"notFound"`
}

// errorResponse is returned by API calls that fail, and member PUTs that succeed.
type errorResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta struct {
		// Resource constant value "string"
		Resource string

		// Version API version
		Version string

		// ID constant value "string"
		ID string

		// Timestamp Response timestamp (milli-seconds from epoch)
		Timestamp int
	}

	// Errors describe errors that occurred
	Errors []apiError
}

// formatErrorResponse extracts the API error into an error
func formatErrorResponse(er *errorResponse) error {
	e := er.Errors[0] // assume there is only ever one error in the array
	return fmt.Errorf("API error status %d: %s", e.Status, strings.Join(e.Detail, ", "))
}

// SAMPLES
// apiError
// 	  {
// 		"status": 401,
// 		"detail": [
// 		  "No permission to read membership"
// 		]
// 	  }

// ErrorResponse
// {
// 	"schemas": [
// 	  "urn:mace:washington.edu:schemas:groups:1.0"
// 	],
// 	"meta": [
// 	  {
// 		"resource": "string",
// 		"version": "v3.0",
// 		"id": "string",
// 		"timestamp": 1214343146201
// 	  }
// 	],
// 	"errors": [
// 	  {
// 		"status": 401,
// 		"detail": [
// 		  "No permission to read membership"
// 		]
// 	  }
// 	]
//   }

// "Success" Error for PUT member
//{
// 	"schemas": [
// 	  "urn:mace:washington.edu:schemas:groups:1.0"
// 	],
// 	"meta": [
// 	  {
// 		"resource": "string",
// 		"version": "v3.0",
// 		"id": "string",
// 		"timestamp": 1214343146201
// 	  }
// 	],
// 	"errors": [
// 	  {
// 		"status": 200,
// 		"notFound": [
// 		  "joeusor"
// 		]
// 	  }
// 	]
//}
