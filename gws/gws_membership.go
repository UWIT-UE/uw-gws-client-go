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
		return make([]Member, 0), formatErrorResponse(resp.Error().(*errorResponse))
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
		return make([]Member, 0), formatErrorResponse(resp.Error().(*errorResponse))
	}

	return resp.Result().(*effMembershipResponse).Members, nil
}

// GetMember returns one member of the group, if present
func (client *Client) GetMember(groupid string, id string) (Member, error) {
	resp, err := client.request().
		SetResult(membershipResponse{}).
		Get(fmt.Sprintf("/group/%s/member/%s", groupid, id))
	if err != nil {
		return Member{}, err
	}
	if resp.IsError() {
		return Member{}, formatErrorResponse(resp.Error().(*errorResponse))
	}

	return resp.Result().(*membershipResponse).Members[0], nil
}

// GetEffectiveMember returns one effective member of the group, if present
func (client *Client) GetEffectiveMember(groupid string, id string) (Member, error) {
	resp, err := client.request().
		SetResult(membershipResponse{}).
		Get(fmt.Sprintf("/group/%s/effective_member/%s", groupid, id))
	if err != nil {
		return Member{}, err
	}
	if resp.IsError() {
		return Member{}, formatErrorResponse(resp.Error().(*errorResponse))
	}

	return resp.Result().(*membershipResponse).Members[0], nil
}

// IsMember indicates true if groupid exists and id is member
// Group not found, member not found or general error all return false
func (client *Client) IsMember(groupid string, id string) (bool, error) {
	member, _ := client.GetMember(groupid, id)
	if member.ID == "" {
		return false, nil
	}
	return true, nil
}

// IsEffectiveMember indicates true if groupid exists and id is effective member
// Group not found, member not found or general error all return false
func (client *Client) IsEffectiveMember(groupid string, id string) (bool, error) {
	member, _ := client.GetEffectiveMember(groupid, id)
	if member.ID == "" {
		return false, nil
	}
	return true, nil
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
		return 0, formatErrorResponse(resp.Error().(*errorResponse))
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
		return 0, formatErrorResponse(resp.Error().(*errorResponse))
	}

	return resp.Result().(*membershipCountResponse).Data.Count, nil
}

// AddOneMember(type, id)
// DeleteOneMember(type, id)
// const types: gws.UWNetID_Member gws.UWWI_Member

// NewMemberShip() a new empty memberlist (or array?)
// UpdateMembership(group, memberlist)
// memberlist.AddUWNetIDMembers(array or id)
// memberlist.AddDNSMembers(array or id)
// memberlist.AddGroupIDMembers(array or id)
// memberlist.AddUWWIMembers(array or id)
// memberlist.AddEPPNMembers(array or id)
