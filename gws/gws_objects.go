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

// GroupResponse what you get back when asking for a Group
type GroupResponse struct {
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

// PutGroup Group packaged for PUT body
// This has no use externally
type putGroup struct {
	Data Group
}

// Group Groups Service group metadata
type Group struct {
	// Unique, opaque idenfier for the group
	Regid string

	// id of the group - includes path
	Name string

	// Descriptive name of the group
	DisplayName string

	// Group's description
	Description string

	// Create timestamp (milli-seconds from epoch)
	Created int

	// Modify timestamp (milli-seconds from epoch)
	LastModified int

	// lastMember timestamp (milli-seconds from epoch)
	LastMemberModified int

	// Contact person (uwnetid) for the group
	Contact UWNetID

	// Multi-factor authn required
	AuthnFactor int `json:",string"`

	// Classification of membership. Enum [ u, r, c, '' ]
	// u=unclassified, r=restricted, c=confidential, missing=no classification
	Classification string

	// Membership dependency group name.  Example: uw_employee
	DependsOn string

	// Numeric GID
	Gid int `json:",string"`

	// Entities with full group access
	Admins []Entity

	// Entities who can edit membership
	Updaters []Entity

	// Entities who can create sub-groups
	Creators []Entity

	// Entities who can read group membership
	Readers []Entity

	// Entities who can opt in to membership
	Optins []Entity

	// Entities who can opt out of membership
	Optouts []Entity
}

// Entity an Entity
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
