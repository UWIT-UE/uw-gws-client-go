package gws

import "fmt"

// Member describes a member of a group
// Used for API constructs: member, effmember and putmember
type Member struct {
	// Type of member enum [ uwnetid, group, dns, eppn, uwwi ]
	Type string `json:"type"`

	// ID of member
	ID string `json:"id"`

	// Type of member enum [ direct, indirect ]
	MType string `json:"-"`

	// Source group(s) if not direct member
	Source string `json:"-"`
}

// MembershipMeta is metadata returned by membership API requests
type MembershipMeta struct {
	// resourceType enum [ groupmembers ]
	ResourceType string

	// Version API version
	Version string

	// RegID the regid of the Group
	RegID string

	// ID the ID of the group
	ID string

	// MembershipType enum [ direct, effective ]
	MembershipType string

	// SelfRef URL of this resource
	SelfRef string

	// Timestamp Response timestamp (milli-seconds from epoch)
	Timestamp int
}

// membershipResponse is what you get back when asking for group membership
type membershipResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta MembershipMeta

	// Data
	Members []Member `json:"data"`
}

// effMmembershipResponse is what you get back when asking for effective group membership
type effMembershipResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta MembershipMeta

	// Data
	Members []Member `json:"data"`
}

// membershipCountResponse is what you get back when asking for group membership count
type membershipCountResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta MembershipMeta

	// Data
	Data struct {
		Count int
	}
}

// putMembership is used when changing membership
type putMembership struct {
	Members []Member `json:"members"`
}

type missingMembersResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// notFoundMembers is array of Members not found
	notFoundMembers []Member
}

// GetMembership returns membership of the group referenced by the groupid
func (client *Client) GetMembership(groupid string) ([]Member, error) {

	resp, err := client.request().
		SetResult(membershipResponse{}).
		Get(fmt.Sprintf("/group/%s/member", groupid))
	if err != nil {
		return make([]Member, 0), err
	}
	if resp.IsError() {
		return make([]Member, 0), decodeErrorResponseBody(resp.Body())
	}

	return resp.Result().(*membershipResponse).Members, nil
}

// GetEffectiveMembership returns membership of the group referenced by the groupid
func (client *Client) GetEffectiveMembership(groupid string) ([]Member, error) {

	resp, err := client.request().
		SetResult(effMembershipResponse{}).
		Get(fmt.Sprintf("/group/%s/effective_member", groupid))
	if err != nil {
		return make([]Member, 0), err
	}
	if resp.IsError() {
		return make([]Member, 0), decodeErrorResponseBody(resp.Body())
	}

	return resp.Result().(*effMembershipResponse).Members, nil
}

// GetMemberCount returns membership count of the group referenced by the groupid
func (client *Client) GetMemberCount(groupid string) (int, error) {

	resp, err := client.request().
		SetResult(membershipCountResponse{}).
		Get(fmt.Sprintf("/group/%s/member?view=count", groupid))
	if err != nil {
		return 0, err
	}
	if resp.IsError() {
		return 0, decodeErrorResponseBody(resp.Body())
	}

	return resp.Result().(*membershipCountResponse).Data.Count, nil
}

// GetEffectiveMemberCount returns membership count of the group referenced by the groupid
func (client *Client) GetEffectiveMemberCount(groupid string) (int, error) {

	resp, err := client.request().
		SetResult(membershipCountResponse{}).
		Get(fmt.Sprintf("/group/%s/effective_member?view=count", groupid))
	if err != nil {
		return 0, err
	}
	if resp.IsError() {
		return 0, decodeErrorResponseBody(resp.Body())
	}

	return resp.Result().(*membershipCountResponse).Data.Count, nil
}
