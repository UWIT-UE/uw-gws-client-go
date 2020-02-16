package gws

import "fmt"

// Valid Affiliates
const (
	AffiliateEmail   = "email"
	AffiliateGoogle  = "google"
	AffiliateUWNetID = "uwnetid"
	AffiliateRadius  = "radius"
)

// Affiliate defines an affiliation state
type Affiliate struct {
	// Affiliate name. 	Enum: [ netid, google, email, radius ]
	Name string `json:"name"`

	// Status of the affiliate. Enum: [ active, inactive ]
	Status string `json:"status"`

	// Sender: authorized email senders.
	Sender []Entity `json:"sender,omitempty"`

	// Forward, email forwarding address.
	Forward string `json:"-"`
}

// searchResponse is returned from a Group search.
type affiliateResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta struct {
		// resourceType enum [ affiliate ]
		ResourceType string

		// Version API version
		Version string

		// RegID the regid of the Group
		RegID string

		// ID the ID of the group
		ID string

		// SelfRef URL of this resource
		SelfRef string

		// Timestamp Response timestamp (milli-seconds from epoch)
		Timestamp int
	}

	// Data an Affiliate struct
	Data Affiliate
}

// putAffiliate is an Affiliate packaged for PUT body.
type putAffiliate struct {
	// Data contains a single group to be put
	Data Group `json:"data"`
}

// GetAffiliateStatus returns the status of an Affiliate of a Group
func (client *Client) GetAffiliateStatus(groupid string, affiliateName string) (*Affiliate, error) {
	resp, err := client.request().
		SetResult(affiliateResponse{}).
		Get(fmt.Sprintf("/group/%s/affiliate/%s", groupid, affiliateName))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	affiliate := resp.Result().(*affiliateResponse).Data
	return &affiliate, nil
}

// AddAffiliate enables this group in UW Google Apps to share documents, calendars, etc.
func (client *Client) AddAffiliate(groupid string, affiliateName string, senders string) (*Affiliate, error) {

	// test senders return 400

	resp, err := client.request().
		SetQueryString(client.syncQueryString()).
		SetQueryParam("status", "active").
		SetQueryParam("sender", senders).
		SetResult(Affiliate{}).
		Put(fmt.Sprintf("/group/%s/affiliate/%s", groupid, affiliateName))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}
	// doesn't return an error payload?
	// 400 Invalid affiliate representation
	// 401 No permission
	// 406 Not Acceptable (Accept or Content-type not application/json)

	// 200 updated
	// 201 new affiliate
	affiliate := resp.Result().(*Affiliate)
	return affiliate, nil
}

// UpdateAffiliate (same as add)

// DeleteAffiliate deletes the Group identified by the specified group id.
func (client *Client) DeleteAffiliate(groupid string, affiliateName string) error {
	resp, err := client.request().
		Delete(fmt.Sprintf("/group/%s/affiliate/%s", groupid, affiliateName))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return formatErrorResponse(resp.Error().(*errorResponse))
	}
	return nil
}
