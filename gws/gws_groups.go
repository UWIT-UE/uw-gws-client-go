package gws

import "fmt"

// Group defines a group, except for membership.
type Group struct {
	// Unique, opaque idenfier for the group
	Regid string `json:"regid,omitempty"`

	// id of the group - includes path
	ID string `json:"id,omitempty"`

	// Descriptive name of the group
	DisplayName string `json:"displayName,omitempty"`

	// Group's description
	Description string `json:"description,omitempty"`

	// Create timestamp (milli-seconds from epoch)
	Created int `json:"created,omitempty"`

	// Modify timestamp (milli-seconds from epoch)
	LastModified int `json:"lastModified,omitempty"`

	// lastMember timestamp (milli-seconds from epoch)
	LastMemberModified int `json:"lastMemberModified,omitempty"`

	// Contact person (uwnetid) for the group
	Contact UWNetID `json:"contact,omitempty"`

	// Multi-factor authn required
	AuthnFactor int `json:"authnfactor,string,omitempty"`

	// Classification of membership. Enum [ u, r, c, '' ]
	// u=unclassified, r=restricted, c=confidential, missing=no classification
	Classification string `json:"classification,omitempty"`

	// Membership dependency group name.  Example: uw_employee
	DependsOn string `json:"dependsOn,omitempty"`

	// Numeric GID
	Gid int `json:"gid,string,omitempty"`

	// Entities with full group access
	Admins []Entity `json:"admins,omitempty"`

	// Entities who can edit membership
	Updaters []Entity `json:"updaters,omitempty"`

	// Entities who can create sub-groups
	Creators []Entity `json:"creators,omitempty"`

	// Entities who can read group membership
	Readers []Entity `json:"readers,omitempty"`

	// Entities who can opt in to membership
	Optins []Entity `json:"optins,omitempty"`

	// Entities who can opt out of membership
	Optouts []Entity `json:"optouts,omitempty"`

	// etag stores the header etag value that arrived with this group
	etag string
}

// groupResponse what comes back when asking for a Group.
type groupResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta struct {
		// resourceType enum [ group ]
		ResourceType string

		// Version API version
		Version string

		// RegID the regid of the Group
		RegID string

		// ID the ID of the group
		ID string

		// SelfRef URL of this resource
		SelfRef string

		// MemberRef URL for this groups members
		MemberRef string

		// Timestamp Response timestamp (milli-seconds from epoch)
		Timestamp int
	}

	// Data a Group struct
	Data Group
}

// putGroup is a Group packaged for PUT body.
type putGroup struct {
	// Data contains a single group to be put
	Data Group `json:"data"`
}

// GetGroup returns the group identified by the groupid.
func (client *Client) GetGroup(groupid string) (*Group, error) {
	resp, err := client.request().
		SetResult(groupResponse{}).
		Get(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	group := resp.Result().(*groupResponse).Data
	group.etag = resp.Header().Get("Etag")
	return &group, nil
}

// CreateGroup creates a new group as specified.
func (client *Client) CreateGroup(newgroup Group) (*Group, error) {
	groupid := newgroup.ID
	body := &putGroup{Data: newgroup}

	resp, err := client.request().
		SetBody(body).
		SetResult(groupResponse{}).
		Put(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	group := resp.Result().(*groupResponse).Data
	group.etag = resp.Header().Get("Etag")
	return &group, nil
}

// UpdateGroup updates an existing Group to match the specified Group.
func (client *Client) UpdateGroup(modgroup Group) (*Group, error) {
	groupid := modgroup.ID
	body := &putGroup{Data: modgroup}

	resp, err := client.request().
		SetHeader("If-Match", modgroup.etag).
		SetBody(body).
		SetResult(groupResponse{}).
		Put(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	group := resp.Result().(*groupResponse).Data
	group.etag = resp.Header().Get("Etag")
	return &group, nil
}

// DeleteGroup deletes the Group identified by the specified group id.
func (client *Client) DeleteGroup(groupid string) error {
	resp, err := client.request().
		Delete(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return formatErrorResponse(resp.Error().(*errorResponse))
	}
	return nil
}

// unimplemented
// move group
// get put delete affiliate
// get history

// helpers for creating/updating group objects
// editable/settable
//  displayname  set
//  description  set
//  contact set
//  authnfactor on/off
//  classification set
//  dependson set
//  entity lists: add/remove/clear
//    admins
//    updaters
//    creators
//    readers
//    optins
//    optouts
