package gws

import "fmt"

// Group Groups Service group metadata
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
}

// groupResponse what you get back when asking for a Group
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

// putGroup is a Group packaged for PUT body
// This has no use externally
type putGroup struct {
	Data Group `json:"data"`
}

// TODO is putGroup readable by json package?

// GetGroup returns the group referenced by the groupid
func (client *Client) GetGroup(groupid string) (Group, error) {
	var group Group

	resp, err := client.request().
		SetResult(groupResponse{}).
		Get(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return group, err
	}
	if resp.IsError() {
		return group, decodeErrorResponseBody(resp.Body())
	}

	group = resp.Result().(*groupResponse).Data
	return group, nil
}

// CreateGroup creates the group provided
func (client *Client) CreateGroup(newgroup Group) (Group, error) {
	var group Group

	groupid := newgroup.ID
	body := &putGroup{Data: newgroup}

	resp, err := client.request().
		SetBody(body).
		SetResult(groupResponse{}).
		Put(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return group, err
	}
	if resp.IsError() {
		return group, decodeErrorResponseBody(resp.Body())
	}

	group = resp.Result().(*groupResponse).Data
	return group, nil
}

// DeleteGroup deletes the group provided
func (client *Client) DeleteGroup(groupid string) error {
	resp, err := client.request().
		Delete(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return decodeErrorResponseBody(resp.Body())
	}

	return nil
}
