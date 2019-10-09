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
