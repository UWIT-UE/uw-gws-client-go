package gws

import (
	"fmt"
	"strings"
)

// Error describes API errors
// Not useful externally
type apiError struct {
	Status int      `json:"status"`
	Detail []string `json:"detail"`
	// undocumented field "notFound" []
}

// ErrorResponse is returned by API calls that fail
// Not useful externally
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

// replace with resty unmarshall
// decodeErrorResponseBody extracts the API error from a Response body
// func decodeErrorResponseBody(body []byte) error {
// 	var er errorResponse
// 	err := json.Unmarshal(body, &er)
// 	if err != nil {
// 		return err
// 	}
// 	e := er.Errors[0] // assume there is only ever one error in the array
// 	return fmt.Errorf("gws error status %d: %s", e.Status, strings.Join(e.Detail, ", "))
// }

// formatErrorResponse extracts the API error into an error
func formatErrorResponse(er *errorResponse) error {
	e := er.Errors[0] // assume there is only ever one error in the array
	return fmt.Errorf("gws error status %d: %s", e.Status, strings.Join(e.Detail, ", "))
}

// SAMPLES
// Error
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
