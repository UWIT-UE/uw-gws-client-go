package gws

import (
	"fmt"
	"regexp"
	"strings"
)

// Entity is a named uwnetid, group, dns eppn or set
// Used in Group objects to represent owner values: admins, updaters, readers, optin, optout.
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

// EntityList returned as part of a Group object
type EntityList []Entity

// Valid Entity types
const (
	EntityTypeUWNetID = "uwnetid"
	EntityTypeGroup   = "group"
	EntityTypeDNS     = "dns"
	EntityTypeEPPN    = "eppn"
	EntityTypeSet     = "set"
	EntityTypeInvalid = ""
)

// NewEntityList creates a blank EntityList to build and then use to overwrite one of a Group's EntityList fields.
// Most commonly this is not used and instead an existing EntityList is modified.
func NewEntityList() *EntityList {
	newList := make(EntityList, 0)
	return &newList
}

// AddEntity adds one or more Entity objects to the referenced EntityList
// Most commonly this is not used and instead Entities are added ByID using other functions.
func (el *EntityList) AddEntity(e ...Entity) (*EntityList, error) {
	// Check if each entity already exists in the EntityList
	for _, newEntity := range e {
		exists := false
		for _, existingEntity := range *el {
			if existingEntity.Type == newEntity.Type && existingEntity.ID == newEntity.ID {
				exists = true
				break
			}
		}
		if exists {
			continue // Entity already exists, skip to the next one
		}
		*el = append(*el, newEntity)
	}

	return el, nil
}

// AppendEntityByID adds an Entity represented by the given ID string to the referenced EntityList.
// Infers the entity type automatically.
func (el *EntityList) AppendEntityByID(id ...string) (*EntityList, error) {

	for _, idStr := range id {
		if el.Contains(idStr) {
			continue
		}
		eType := inferredEType(idStr)
		if eType == EntityTypeInvalid {
			return el, fmt.Errorf("Entity type could not be inferred for ID: %s", idStr)
		}
		*el = append(*el, Entity{Type: eType, ID: idStr})
	}
	return el, nil

}

// RemoveEntityByID removes one or more Entities from the referenced EntityList by ID string
func (el *EntityList) RemoveEntityByID(id ...string) (*EntityList, error) {
	for _, idStr := range id {
		for i, e := range *el {
			if e.ID == idStr {
				*el = append((*el)[:i], (*el)[i+1:]...)
				break
			}
		}
	}
	return el, nil
}

// ToIDs renders an EntityList as a new slice containing only the ID strings.
func (el EntityList) ToIDs() []string {
	eIDs := make([]string, 0, len(el))
	for _, e := range el {
		eIDs = append(eIDs, e.ID)
	}
	return eIDs
}

// ToCommaString renders an EntityList as a comma separated string of IDs
func (el EntityList) ToCommaString() string {
	return strings.Join(el.ToIDs(), ",")
}

// Filter returns a new EntityList that does not contain members of the specified Entity type
func (el EntityList) Filter(eType string) *EntityList {
	newList := make(EntityList, 0)
	for _, e := range el {
		if e.Type != eType {
			newList = append(newList, e)
		}
	}
	return &newList
}

// Match returns a new EntityList containing only the entities that match the given Entity type
func (el EntityList) Match(eType string) *EntityList {
	newList := make(EntityList, 0)
	for _, e := range el {
		if e.Type == eType {
			newList = append(newList, e)
		}
	}
	return &newList
}

// Contains returns true if the EntityList contains the given Entity ID
func (el EntityList) Contains(id string) bool {
	for _, e := range el {
		if e.ID == id {
			return true
		}
	}
	return false
}

// inferredEType returns a rough guess of the entity type for the given Entity ID string
func inferredEType(id string) (string) {

	// Eliminate obvious non-entities
	if !regexp.MustCompile(`[\w\._:\-@\$]+`).MatchString(id) {
		return EntityTypeInvalid
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

	// Set
	// BUG(x): there is a tiny chance that this could mis-infer the "set" type since the UWNetIDs: "all", "uw", "none" also exist.
	if id == "all" || id == "none" || id == "uw" || id == "member" {
		return EntityTypeSet
	}

	return EntityTypeUWNetID
}
