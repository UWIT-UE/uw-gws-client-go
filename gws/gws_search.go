package gws

import "fmt"

// searchResponse what you get back when searching for Groups
type searchResponse struct {
	// Schema The schema in use. Enum [ "urn:mace:washington.edu:schemas:groups:1.0" ]
	Schemas []string

	// Meta Group metadata
	Meta struct {
		// resourceType enum [ search ]
		ResourceType string

		// Version API version
		Version string

		// RegID the regid of the Group
		SearchParameters struct {
			Name       string
			Stem       string
			Scope      string
			Member     string
			Type       string
			Owner      string
			Affiliate  string
			Instructor string
		}

		// SelfRef URL of this resource
		SelfRef string

		// Timestamp Response timestamp (milli-seconds from epoch)
		Timestamp int
	}

	// Data a Group struct
	Data []GroupReference
}

// GroupReference reference to a group returned by a search
type GroupReference struct {
	// Unique, opaque idenfier for the group
	Regid string

	// id of the group - includes path
	ID string

	// Descriptive name of the group
	DisplayName string

	// URL the URL for this Group
	URL string

	// Via these represent the indirect group paths for effective searches
	Via []string
}

// SearchParameters holds the parameters to submit to a group search.
type SearchParameters struct {
	parameters map[string]string
}

// NewSearch creates a blank query to submit for searching
func NewSearch() *SearchParameters {
	return &SearchParameters{parameters: make(map[string]string)}
}

// DoSearch submits a search for groups with supplied search parameters.
func (client *Client) DoSearch(s *SearchParameters) ([]GroupReference, error) {
	var gr []GroupReference

	resp, err := client.request().
		SetResult(searchResponse{}).
		Get(fmt.Sprintf("/search?%s", s.queryString()))
	if err != nil {
		return gr, err
	}
	if resp.IsError() {
		return gr, decodeErrorResponseBody(resp.Body())
	}

	return resp.Result().(*searchResponse).Data, nil
}

// WithName adds a match on name. Name is some part of the group id, "*" is a wildcard.
func (s *SearchParameters) WithName(name string) *SearchParameters {
	if name != "" {
		s.parameters["name"] = name
	}
	return s
}

// WithStem adds a match on stem. Stem is the stem part of the group id.
func (s *SearchParameters) WithStem(stem string) *SearchParameters {
	if stem != "" {
		s.parameters["stem"] = stem
	}
	return s
}

// WithScope adds a match on scope.
func (s *SearchParameters) WithScope(scope string) *SearchParameters {
	if scope != "" {
		s.parameters["scope"] = scope
	}
	return s
}

// WithMember adds match for groups with the specified member id.
func (s *SearchParameters) WithMember(id string) *SearchParameters {
	if id != "" {
		s.parameters["member"] = id
	}
	return s
}

// InEffectiveMembers matches effective members when searching for members, owners, instructors.
func (s *SearchParameters) InEffectiveMembers() *SearchParameters {
	s.parameters["type"] = "effective"
	return s
}

// InDirectMembers matches direct members when searching for members, owners, instructors, this is default.
func (s *SearchParameters) InDirectMembers() *SearchParameters {
	s.parameters["type"] = "direct"
	return s
}

// WithOwner adds match for groups where an administrator (admin, creator, updater, member manager) is the specified id.
func (s *SearchParameters) WithOwner(id string) *SearchParameters {
	if id != "" {
		s.parameters["owner"] = id
	}
	return s
}

// WithInstructor adds match for groups where the instructor is the specified id.
func (s *SearchParameters) WithInstructor(id string) *SearchParameters {
	if id != "" {
		s.parameters["instructor"] = id
	}
	return s
}

// WithAffiliate adds match for groups where the affiliate is the specified id. The affiliate search ignores any other search parameters.
func (s *SearchParameters) WithAffiliate(id string) *SearchParameters {
	if id != "" {
		s.parameters["affiliate"] = id
	}
	return s
}

// queryString assembles the parameters into a query string
func (s *SearchParameters) queryString() string {
	qs := ""
	separator := ""

	for k, v := range s.parameters {
		if v == "" {
			continue
		}
		qs = qs + fmt.Sprintf("%s%s=%s", separator, k, v)
		separator = "&"
	}
	return qs
}
