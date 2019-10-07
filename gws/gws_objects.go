package gws

// UWNetID a University of Washington UWNetID
type UWNetID string

// EPPNString an eduPersonPrincipalName or an email address
type EPPNString string

// DNSString a Domain Name System (DNS) address
type DNSString string

// UWWIString a Microsoft Infrastructure (MI) machine name (with a $ appended)
type UWWIString string

// Affiliate an affiliate
// enum [ email, google, uwnetid, radius ]
type Affiliate string

// EmailSendersString Exchange Email senders - a comma separated list of ids
// example: joeuser,u_joeuser_friends
type EmailSendersString string

// GoogleSenderString Google Groups senders - choice keyword
// enum [ none, all, members, uw ]
type GoogleSenderString string

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

// TODO is putGroup readable by json package?
// PutGroup Group packaged for PUT body
// This has no use externally
type putGroup struct {
	Data Group `json:"data"`
}

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

// Entity is a named uwnetid, group, dns eppn or set
type Entity struct {
	// Type of entity. Enum [ uwnetid, group, dns, eppn, set ]
	Type string `json:"type,omitempty"`

	// ID of entity
	// If the type is 'set' the id is:
	//   all: any entity
	//   none: no entity
	//   uw: any UW member entity
	//   member: any member of the group
	ID string `json:"id,omitempty"`

	// Display name of entity.
	Name string `json:"name,omitempty"`
}

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

// Error describes API errors
// Not useful externally
type apiError struct {
	Status int      `json:"status"`
	Detail []string `json:"detail"`
	// udocumented field "notFound" []
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
