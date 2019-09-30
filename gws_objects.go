package uwgwsclient

// A University of Washington UWNetID
type UWNetID []string

// An eduPersonPrincipalName or an email address
type EPPNString []string

// A Domain Name System (DNS) address
type DNSString []string

// A Microsoft Infrastructure (MI) machine name (with a $ appended)
type UWWIString []string

// An affiliate
// enum [ email, google, uwnetid, radius ]
type Affiliate []string

// Exchange Email senders - a comma separated list of ids
// example: joeuser,u_joeuser_friends
type EmailSendersString []string

// Google Groups senders - choice keyword
// enum [ none, all, members, uw ]
type GoogleSenderString []string

type Group struct {
	// Unique, opaque idenfier for the group
	Regid string `json:"regid"`

	// id of the group - includes path
	Name string `json:"name"`

	// Descriptive name of the group
	DisplayName string `json:"displayName"`

	// Group's description
	Description string `json:"description"`

	// Create timestamp (milli-seconds from epoch)
	Created int `json:"created"`

	// Modify timestamp (milli-seconds from epoch)
	LastModified int `json:"lastModified"`

	// lastMember timestamp (milli-seconds from epoch)
	LastMemberModified int `json:"lastMemberModified"`

	// Contact person (uwnetid) for the group
	Contact UWNetID `json:"contact"`

	// Multi-factor authn required
	AuthnFactor int `json:"authnfactor"`

	// Classification of membership. Enum [ u, r, c, '' ]
	// u=unclassified, r=restricted, c=confidential, missing=no classification
	Classification string `json:"classification"`

	// Membership dependency group name.  Example: uw_employee
	DependsOn string `json:"dependsOn"`

	// Numeric GID
	Gid int `json:"gid"`

	// Entities with full group access
	Admins []Entity `json:"admins"`

	// Entities who can edit membership
	Updaters []Entity `json:"updaters"`

	// Entities who can create sub-groups
	Creators []Entity `json:"creators"`

	// Entities who can read group membership
	Readers []Entity `json:"readers"`

	// Entities who can opt in to membership
	Optins []Entity `json:"optins"`

	// Entities who can opt out of membership
	Optouts []Entity `json:"optouts"`
}

// An Entity
type Entity struct {
	// Type of entity. Enum [ uwnetid, group, dns, eppn, set ]
	EntityType string

	// ID of entity
	// If the type is 'set' the id is:
	//   all: any entity
	//   none: no entity
	//   uw: any UW member entity
	//   member: any member of the group
	EntityID string

	// Display name of entity.
	EntityName string
}
