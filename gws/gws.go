package gws

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

// Config holds all the Client configuration
type Config struct {
	APIUrl        string
	Timeout       time.Duration
	SkipTLSVerify bool
	CAFile        string
	ClientCert    string
	ClientKey     string
}

// Client wraps resty.Client
type Client struct {
	resty      *resty.Client
	config     *Config
	configured bool
}

// DefaultConfig constructs a basic Config object
func DefaultConfig() *Config {
	dc := &Config{
		APIUrl:        "https://groups.uw.edu/group_sws/v3", // requires CAFile to be incommon
		Timeout:       30,
		SkipTLSVerify: false,
	}
	return dc
}

// NewClient builds a new client with defaults
func NewClient(config *Config) (*Client, error) {
	restyInst := resty.New()
	nc := &Client{resty: restyInst, config: config}

	// setTLSClientConfig must be before other TLS in configure()
	nc.resty.SetTLSClientConfig(&tls.Config{Renegotiation: tls.RenegotiateFreelyAsClient})
	nc.configure()

	// Standardize redirect policy
	//restyInst.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))

	// JSON
	restyInst.SetHeader("Accept", "application/json")
	restyInst.SetHeader("Content-Type", "application/json")

	return nc, nil
}

func (client *Client) configure() {
	if client.configured {
		return
	}

	restyInst := client.resty
	config := client.config

	restyInst.SetHostURL(config.APIUrl)
	restyInst.SetTimeout(config.Timeout * time.Second)

	// setup TLS here
	restyInst.SetRootCertificate(config.CAFile)
	cert, err := tls.LoadX509KeyPair(config.ClientCert, config.ClientKey)
	if err != nil {
		log.Fatalf("ERROR client certificate: %s", err)
		return
	}
	restyInst.SetCertificates(cert)
	restyInst.SetDebug(true)
	//fmt.Printf("%#v\n", config)
	client.configured = true
}

// SetTLSClientConfig assigns client TLS config
func (client *Client) SetTLSClientConfig(c *tls.Config) {
	client.resty.SetTLSClientConfig(c)
}

// R returns new resty.Request from configured client
func (client *Client) request() *resty.Request {
	client.configure()
	request := client.resty.R()

	return request
}

// GetGroup returns the group referenced by the groupid
func (client *Client) GetGroup(groupid string) (Group, error) {
	var group Group

	resp, err := client.request().
		SetResult(groupResponse{}).
		Get(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return group, err
	}
	if resp.IsError() {
		return group, decodeErrorResponseBody(resp.Body())
	}

	group = resp.Result().(*groupResponse).Data
	return group, nil
}

// CreateGroup creates the group provided
func (client *Client) CreateGroup(newgroup Group) (Group, error) {
	var group Group

	groupid := newgroup.ID
	body := &putGroup{Data: newgroup}

	resp, err := client.request().
		SetBody(body).
		SetResult(groupResponse{}).
		Put(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return group, err
	}
	if resp.IsError() {
		return group, decodeErrorResponseBody(resp.Body())
	}

	group = resp.Result().(*groupResponse).Data
	return group, nil
}

// DeleteGroup deletes the group provided
func (client *Client) DeleteGroup(groupid string) error {
	resp, err := client.request().
		Delete(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return decodeErrorResponseBody(resp.Body())
	}

	return nil
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

// func ToEntityList(item *Entity) []Entity {
// 	var ea []Entity
// 	ea = append(ea, *item)
// 	return ea
// }

// func decodeErrorResponse(er errorResponse) error {
// 	e := er.Errors[0] // assume there is only ever one error in the array
// 	fmt.Println("detail", e.Detail)
// 	err := fmt.Errorf("gws error status %d: %s", e.Status, strings.Join(e.Detail, ", "))

// 	return err
// }

// decodeErrorResponseBody extracts the API error from a Response body
func decodeErrorResponseBody(body []byte) error {
	var er errorResponse
	err := json.Unmarshal(body, &er)
	if err != nil {
		return err
	}
	e := er.Errors[0] // assume there is only ever one error in the array
	return fmt.Errorf("gws error status %d: %s", e.Status, strings.Join(e.Detail, ", "))
}

// func decodeErrorResponse(er []Error) error {
// 	e := er[0] // assume there is only ever one error in the array
// 	err := fmt.Errorf("gws error status %d: %q", e.Status, strings.Join(e.Detail, ", "))

// 	return err
// }

// func ToEntity(items []*Entity) []Entity {
// 	var ea []Entity
// 	for _, item := range items {
// 		ea = append(ea, *item)
// 	}
// 	return ea
// }
