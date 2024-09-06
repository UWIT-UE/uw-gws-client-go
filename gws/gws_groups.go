package gws

import "fmt"

// Group defines a group, except for membership.
type Group struct {
	// Unique, opaque identifier for the group
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

	// Modify timestamp of group membership (milli-seconds from epoch)
	LastMemberModified int `json:"lastMemberModified,omitempty"`

	// Contact person (uwnetid) for the group
	Contact UWNetID `json:"contact,omitempty"`

	// Multi-factor authn required
	AuthnFactor int `json:"authnfactor,string,omitempty"`

	// Classification of membership. Enum [ u, r, c, '' ]
	// u=public(unclassified), r=restricted, c=confidential, missing=no classification
	Classification DataClassification `json:"classification,omitempty"`

	// Membership dependency group name.  Example: uw_employee
	DependsOn string `json:"dependsOn,omitempty"`

	// Numeric GID
	Gid int `json:"gid,string,omitempty"`

	// Entities with full group access
	Admins EntityList `json:"admins,omitempty"`

	// Entities who can edit membership
	Updaters EntityList `json:"updaters,omitempty"`

	// Entities who can create sub-groups
	Creators EntityList `json:"creators,omitempty"`

	// Entities who can read group membership
	Readers EntityList `json:"readers,omitempty"`

	// Entities who can opt in to membership
	Optins EntityList `json:"optins,omitempty"`

	// Entities who can opt out of membership
	Optouts EntityList `json:"optouts,omitempty"`

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

type DataClassification string

// Data Classifications
const (
	DataClassificationPublic       DataClassification = "u"
	DataClassificationRestricted   DataClassification = "r"
	DataClassificationConfidential DataClassification = "c"
	DataClassificationNone         DataClassification = ""
)

func (dc DataClassification) String() string {
	if dc == DataClassificationPublic {
		return "Public"
	}
	if dc == DataClassificationRestricted {
		return "Restricted"
	}
	if dc == DataClassificationConfidential {
		return "Confidential"
	}
	return "None"
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

// CreateGroup creates a new group as defined by the specified Group.
func (client *Client) CreateGroup(newgroup *Group) (*Group, error) {
	groupid := newgroup.ID
	body := &putGroup{Data: *newgroup}

	resp, err := client.request().
		SetBody(body).
		SetQueryString(client.syncQueryString()).
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
func (client *Client) UpdateGroup(modgroup *Group) (*Group, error) {
	groupid := modgroup.ID
	body := &putGroup{Data: *modgroup}

	resp, err := client.request().
		SetHeader("If-Match", modgroup.etag).
		SetQueryString(client.syncQueryString()).
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

// TODO unimplemented
// move group
// get put delete affiliate
// get history

// SetAuthnFactor sets the multi-factor authn required for the group
func (group *Group) SetAuthnFactor(factor int) (*Group, error) {
	if factor != 0 && factor != 1 && factor != 2 {
		return group, fmt.Errorf("invalid authn factor. Valid values are 0, 1, or 2")
	}
	group.AuthnFactor = factor
	return group, nil
}

// SetClassification sets the data classification of the group
func (group *Group) SetClassification(classification DataClassification) (*Group, error) {
	if classification != DataClassificationPublic && classification != DataClassificationRestricted && classification != DataClassificationConfidential && classification != DataClassificationNone {
		return group, fmt.Errorf("invalid data classification")
	}
	group.Classification = classification
	return group, nil
}

// SetDependsOn sets the group name that this group membership depends on
func (group *Group) SetDependsOn(dependsOn string) (*Group, error) {
	if inferredEType(dependsOn) != EntityTypeGroup {
		return group, fmt.Errorf("invalid dependsOn value: must be a group")
	}
	group.DependsOn = dependsOn
	return group, nil
}

// AddAdmin adds one or more Entity IDs to the Admin EntityList
func (group *Group) AddAdmin(id ...string) (*Group, error) {
	_, err := group.Admins.AppendEntityByID(id...)
	return group, err
}

// RemoveAdmin removes one or more Entity IDs from the Admin EntityList
func (group *Group) RemoveAdmin(id ...string) (*Group, error) {
	_, err := group.Admins.RemoveEntityByID(id...)
	return group, err
}

// RemoveAllAdmins removes all Entities from the Admin EntityList
func (group *Group) RemoveAllAdmins() (*Group, error) {
	group.Admins = EntityList{}
	return group, nil
}

// IsAdmin returns true if the given Entity ID is in the Admin EntityList
func (group *Group) IsAdmin(id string) bool {
	return group.Admins.Contains(id)
}

// AddUpdater adds one or more Entity IDs to the Updater EntityList
func (group *Group) AddUpdater(id ...string) (*Group, error) {
	_, err := group.Updaters.AppendEntityByID(id...)
	return group, err
}

// RemoveUpdater removes one or more Entity IDs from the Updater EntityList
func (group *Group) RemoveUpdater(id ...string) (*Group, error) {
	_, err := group.Updaters.RemoveEntityByID(id...)
	return group, err
}

// RemoveAllUpdaters removes all Entities from the Updater EntityList
func (group *Group) RemoveAllUpdaters() (*Group, error) {
	group.Updaters = EntityList{}
	return group, nil
}

// IsUpdater returns true if the given Entity ID is in the Updater EntityList
func (group *Group) IsUpdater(id string) bool {
	return group.Updaters.Contains(id)
}

// AddCreator adds one or more Entity IDs to the Creator EntityList
func (group *Group) AddCreator(id ...string) (*Group, error) {
	_, err := group.Creators.AppendEntityByID(id...)
	return group, err
}

// RemoveCreator removes one or more Entity IDs from the Creator EntityList
func (group *Group) RemoveCreator(id ...string) (*Group, error) {
	_, err := group.Creators.RemoveEntityByID(id...)
	return group, err
}

// RemoveAllCreators removes all Entities from the Creator EntityList
func (group *Group) RemoveAllCreators() (*Group, error) {
	group.Creators = EntityList{}
	return group, nil
}

// IsCreator returns true if the given Entity ID is in the Creator EntityList
func (group *Group) IsCreator(id string) bool {
	return group.Creators.Contains(id)
}

// AddReader adds one or more Entity IDs to the Reader EntityList
func (group *Group) AddReader(id ...string) (*Group, error) {
	_, err := group.Readers.AppendEntityByID(id...)
	return group, err
}

// RemoveReader removes one or more Entity IDs from the Reader EntityList
func (group *Group) RemoveReader(id ...string) (*Group, error) {
	_, err := group.Readers.RemoveEntityByID(id...)
	return group, err
}

// RemoveAllReaders removes all Entities from the Reader EntityList
func (group *Group) RemoveAllReaders() (*Group, error) {
	group.Readers = EntityList{}
	return group, nil
}

// IsReader returns true if the given Entity ID is in the Reader EntityList
func (group *Group) IsReader(id string) bool {
	return group.Readers.Contains(id)
}

// AddOptin adds one or more Entity IDs to the Optin EntityList
func (group *Group) AddOptin(id ...string) (*Group, error) {
	_, err := group.Optins.AppendEntityByID(id...)
	return group, err
}

// RemoveOptin removes one or more Entity IDs from the Optin EntityList
func (group *Group) RemoveOptin(id ...string) (*Group, error) {
	_, err := group.Optins.RemoveEntityByID(id...)
	return group, err
}

// RemoveAllOptins removes all Entities from the Optin EntityList
func (group *Group) RemoveAllOptins() (*Group, error) {
	group.Optins = EntityList{}
	return group, nil
}

// IsOptin returns true if the given Entity ID is in the Optin EntityList
func (group *Group) IsOptin(id string) bool {
	return group.Optins.Contains(id)
}

// AddOptout adds one or more Entity IDs to the Optout EntityList
func (group *Group) AddOptout(id ...string) (*Group, error) {
	_, err := group.Optouts.AppendEntityByID(id...)
	return group, err
}

// RemoveOptout removes one or more Entity IDs from the Optout EntityList
func (group *Group) RemoveOptout(id ...string) (*Group, error) {
	_, err := group.Optouts.RemoveEntityByID(id...)
	return group, err
}

// RemoveAllOptouts removes all Entities from the Optout EntityList
func (group *Group) RemoveAllOptouts() (*Group, error) {
	group.Optouts = EntityList{}
	return group, nil
}

// IsOptout returns true if the given Entity ID is in the Optout EntityList
func (group *Group) IsOptout(id string) bool {
	return group.Optouts.Contains(id)
}
