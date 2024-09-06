package gws

import (
	"fmt"
	"strings"
)

// membershipMeta is metadata returned by membership API requests.
type membershipMeta struct {
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

// membershipResponse is what you get back when asking for group membership.
type membershipResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta membershipMeta

	// Data
	Members MemberList `json:"data"`
}

// effMmembershipResponse is what you get back when asking for effective group membership.
type effMembershipResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta membershipMeta

	// Data
	Members MemberList `json:"data"`
}

// membershipCountResponse is what you get back when asking for group membership count.
type membershipCountResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta membershipMeta

	// Data
	Data struct {
		Count int
	}
}

// putMembership is used when changing membership.
type putMembership struct {
	Members MemberList `json:"data"`
}

// GetMembership returns membership of the group specified by the groupid.
func (client *Client) GetMembership(groupid string) (*MemberList, error) {

	resp, err := client.request().
		SetResult(membershipResponse{}).
		Get(fmt.Sprintf("/group/%s/member", groupid))
	if err != nil {
		return &MemberList{}, err // make(MemberList, 0), err
	}
	if resp.IsError() {
		return &MemberList{}, formatErrorResponse(resp.Error().(*errorResponse))
	}
	return &resp.Result().(*membershipResponse).Members, nil
}

// GetEffectiveMembership returns membership of the group referenced by the groupid.
func (client *Client) GetEffectiveMembership(groupid string) (*MemberList, error) {

	resp, err := client.request().
		SetResult(effMembershipResponse{}).
		Get(fmt.Sprintf("/group/%s/effective_member", groupid))
	if err != nil {
		return &MemberList{}, err //make(MemberList, 0), err
	}
	if resp.IsError() {
		return &MemberList{}, formatErrorResponse(resp.Error().(*errorResponse))
	}
	return &resp.Result().(*effMembershipResponse).Members, nil
}

// GetMember returns one member of the group, if present.
func (client *Client) GetMember(groupid string, id string) (*Member, error) {
	resp, err := client.request().
		SetResult(membershipResponse{}).
		Get(fmt.Sprintf("/group/%s/member/%s", groupid, id))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	m := resp.Result().(*membershipResponse).Members[0]
	return &m, nil
}

// GetEffectiveMember returns one effective member of the group, if present.
func (client *Client) GetEffectiveMember(groupid string, id string) (*Member, error) {
	resp, err := client.request().
		SetResult(membershipResponse{}).
		Get(fmt.Sprintf("/group/%s/effective_member/%s", groupid, id))
	if err != nil {
		//return Member{}, err
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	m := resp.Result().(*membershipResponse).Members[0]
	return &m, nil
}

// IsMember indicates true if groupid exists and id is member.
// Group not found, member not found or general error all return false.
func (client *Client) IsMember(groupid string, id string) (bool, error) {
	member, _ := client.GetMember(groupid, id)
	if member.ID == "" {
		return false, nil
	}
	return true, nil
}

// IsEffectiveMember indicates true if groupid exists and id is effective member.
// Group not found, member not found or general error all return false.
func (client *Client) IsEffectiveMember(groupid string, id string) (bool, error) {
	member, _ := client.GetEffectiveMember(groupid, id)
	if member.ID == "" {
		return false, nil
	}
	return true, nil
}

// MemberCount returns membership count of the group referenced by the groupid.
// Group not found or general error returns a count of zero.
func (client *Client) MemberCount(groupid string) (int, error) {
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

// EffectiveMemberCount returns membership count of the group referenced by the groupid.
// Group not found or general error returns a count of zero.
func (client *Client) EffectiveMemberCount(groupid string) (int, error) {
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

// AddMembers adds one or more member IDs to the referenced group and returns an array of memberIDs that do not exist and could not be added.
func (client *Client) AddMembers(groupid string, memberIDs ...string) ([]string, error) {
	resp, err := client.request().
		SetQueryString(client.syncQueryString()).
		SetResult(errorResponse{}).
		Put(fmt.Sprintf("/group/%s/member/%s", groupid, strings.Join(memberIDs, ",")))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	// PUT member is weird, returns "error" on 200
	er := resp.Result().(*errorResponse)
	return er.Errors[0].NotFound, nil
}

// DeleteMembers removes one or more member IDs from the referenced group.
func (client *Client) DeleteMembers(groupid string, memberIDs ...string) error {
	resp, err := client.request().
		SetQueryString(client.syncQueryString()).
		Delete(fmt.Sprintf("/group/%s/member/%s", groupid, strings.Join(memberIDs, ",")))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return formatErrorResponse(resp.Error().(*errorResponse))
	}
	return nil
}

// SetMembership completely replaces group membership with specified MemberList and returns an array of memberIDs that do not exist and could not be added.
func (client *Client) SetMembership(groupid string, newMembers *MemberList) ([]string, error) {
	body := &putMembership{Members: *newMembers}

	resp, err := client.request().
		SetQueryString(client.syncQueryString()).
		SetBody(body).
		SetResult(errorResponse{}).
		Put(fmt.Sprintf("/group/%s/member", groupid))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, formatErrorResponse(resp.Error().(*errorResponse))
	}

	// PUT member is weird, returns "error" on 200
	okError := resp.Result().(*errorResponse)
	return okError.Errors[0].NotFound, nil
}

// DeleteAllMembers removes all members from the referenced group.
func (client *Client) DeleteAllMembers(groupid string) error {
	body := &putMembership{Members: make(MemberList, 0)}

	resp, err := client.request().
		SetQueryString(client.syncQueryString()).
		SetBody(body).
		Put(fmt.Sprintf("/group/%s/member", groupid))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return formatErrorResponse(resp.Error().(*errorResponse))
	}
	return nil
}
