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
	nc.resty.SetTLSClientConfig(&tls.Config{Renegotiation: tls.RenegotiateOnceAsClient})
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

// GetGroup get the group referenced by the groupid
func (client *Client) GetGroup(groupid string) (Group, error) {
	var group Group

	resp, err := client.request().SetResult(groupResponse{}).Get(fmt.Sprintf("/group/%s", groupid))
	if err != nil {
		log.Fatal(err)
		return group, err
	}
	if resp.IsError() {
		var er errorResponse
		err := json.Unmarshal(resp.Body(), &er)
		if err != nil {
			fmt.Println("error:", err)
		}
		return group, decodeErrorResponse(er)
		//log.Fatal(resp.StatusCode)
		//return group, err
	}
	group = resp.Result().(*groupResponse).Data

	return group, nil
}

// CreateGroup get the group referenced by the groupid
func (client *Client) CreateGroup(newgroup Group) (Group, error) {
	var group Group

	body := &putGroup{Data: newgroup}
	groupid := newgroup.ID

	resp, err := client.request().SetResult(groupResponse{}).SetBody(body).Put(fmt.Sprintf("/group/%s", groupid))
	if resp.IsError() {
		var er errorResponse
		err := json.Unmarshal(resp.Body(), &er)
		if err != nil {
			fmt.Println("error:", err)
		}
		return group, decodeErrorResponse(er)
		//log.Fatal(resp.StatusCode)
		//return group, err
	}
	// groupr := resp.Result().(*GroupResponse)
	// fmt.Printf("%#v\n", groupr)
	// if resp.IsError() {
	// 	//var er ErrorResponse
	// 	return group, decodeErrorResponse(resp.Result().(*GroupResponse).Errors)
	// 	//log.Fatal(resp.StatusCode)
	// 	//return group, err
	// }
	if err != nil {
		log.Fatal(err)
		return group, err
	}

	// don't unmarshall inside of resty
	// if error unmarshall with ErrorResponse
	// if success unmarshall with GroupResponse

	// resty doesn't report errors
	// decode errors function when statuscode is not 200

	fmt.Printf("%#v\n", resp)
	fmt.Printf("%#v\n", resp.Status())
	fmt.Printf("%#v\n", err)
	group = resp.Result().(*groupResponse).Data

	return group, nil
}

func ToEntityList(item *Entity) []Entity {
	var ea []Entity
	ea = append(ea, *item)
	return ea
}

func decodeErrorResponse(er errorResponse) error {
	e := er.Errors[0] // assume there is only ever one error in the array
	fmt.Println("detail", e.Detail)
	err := fmt.Errorf("gws error status %d: %s", e.Status, strings.Join(e.Detail, ", "))

	return err
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
