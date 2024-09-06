package gws

import (
	"fmt"
	"regexp"
	"strings"
)

// Member fully describes a member of a group.
type Member struct {
	// Type of member enum [ uwnetid, group, dns, eppn, uwwi ]
	Type MemberType `json:"type"`

	// ID of member
	ID string `json:"id"`

	// Type of member enum [ direct, indirect ]
	MType string `json:"-"`

	// Source group(s) if not direct member
	Source string `json:"-"`
}

// MemberList is a slice of Members, returned by membership requests.
type MemberList []Member

// Valid Member types returned by membership calls. Useful for Filter() and Match().
type MemberType string

const (
	MemberTypeUWNetID MemberType = "uwnetid"
	MemberTypeGroup   MemberType = "group"
	MemberTypeDNS     MemberType = "dns"
	MemberTypeEPPN    MemberType = "eppn"
	MemberTypeUWWI    MemberType = "uwwi"
	MemberTypeInvalid MemberType = ""
)

// ToIDs renders a MemberList as a slice containing only member ID strings.
// Discards other Member fields in the process.
func (members MemberList) ToIDs() []string {
	memberIDs := make([]string, 0, len(members))
	for _, member := range members {
		memberIDs = append(memberIDs, member.ID)
	}
	return memberIDs
}

// ToCommaString renders a MemberList as a string of comma joined member IDs.
// Discards other Member fields in the process.
func (members MemberList) ToCommaString() string {
	return strings.Join(members.ToIDs(), ",")
}

// Filter returns a new MemberList without members of the specified type.
func (members *MemberList) Filter(memberType MemberType) *MemberList {
	newList := make(MemberList, 0)
	for _, member := range *members {
		if member.Type != memberType {
			newList = append(newList, member)
		}
	}
	return &newList
}

// Match returns a new MemberList containing only the specified member type.
func (members *MemberList) Match(memberType MemberType) *MemberList {
	newList := make(MemberList, 0)
	for _, member := range *members {
		if member.Type == memberType {
			newList = append(newList, member)
		}
	}
	return &newList
}

// Contains returns true if the MemberList contains the given Member ID
func (ml MemberList) Contains(id string) bool {
	for _, m := range ml {
		if m.ID == id {
			return true
		}
	}
	return false
}

// Functions to manipulate and set full MemberLists via SetMembership()

// inferredMType returns a rough guess of the Member type for the given Member ID string
func inferredMType(id string) MemberType {

	// Eliminate obvious non-entities
	if !regexp.MustCompile(`[\w\._:\-@\$]+`).MatchString(id) {
		return MemberTypeInvalid
	}

	// UWWI (Microsoft machine ID ending in $)
	if strings.HasSuffix(id, "$") {
		return MemberTypeUWWI
	}

	// EPPN
	if strings.Contains(id, "@") {
		return EntityTypeEPPN
	}

	// Groups
	if strings.Contains(id, ":") {
		return EntityTypeGroup
	}
	if strings.HasPrefix(id, "uw_") {
		return EntityTypeGroup
	}
	if strings.HasPrefix(id, "g_") {
		return EntityTypeGroup
	}
	if strings.HasPrefix(id, "u_") {
		return EntityTypeGroup
	}
	if strings.HasPrefix(id, "course_") {
		return EntityTypeGroup
	}

	// DNS
	if regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)+([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`).MatchString(id) {
		return EntityTypeDNS
	}

	return EntityTypeUWNetID
}

// AddMemberByID modifies a MemberList, inferring the MemberType if each id and appending Members
func (ml *MemberList) AppendMemberByID(id ...string) (*MemberList, error) {
	for _, idStr := range id {
		if ml.Contains(idStr) {
			continue
		}
		mType := inferredMType(idStr)
		if mType == EntityTypeInvalid {
			// returns only one value. print warning only?
			return ml, fmt.Errorf("Member type could not be inferred for ID: %s", idStr)
		}
		*ml = append(*ml, Member{Type: mType, ID: idStr})
	}
	return ml, nil
}

// DeleteMemberByID modifies a MemberList, removing supplied ID strings from Members.
func (ml *MemberList) RemoveMemberByID(id ...string) (*MemberList, error) {
	for _, idStr := range id {
		for i, m := range *ml {
			if m.ID == idStr {
				*ml = append((*ml)[:i], (*ml)[i+1:]...)
				break
			}
		}
	}
	return ml, nil
}

// NewMemberList creates a blank MemberList to build up and then use to set a groups membership.
func NewMemberList() *MemberList {
	newList := make(MemberList, 0)
	return &newList
}

// appendMembers is for manual manipulation of MemberLists to force the MemberType
// Usually using AddMemberByID() is easier, letting it infer the MemberType
func (members *MemberList) AppendMembers(memberType MemberType, memberIDs []string) *MemberList {

	for _, member := range memberIDs {
		if member != "" {
			*members = append(*members, Member{Type: memberType, ID: member})
		}
	}
	return members
}
