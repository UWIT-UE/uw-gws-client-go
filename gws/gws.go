// Package gws provides an API client for the University of Washington Groups Service API
package gws

import (
	"crypto/tls"
	"log"
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
	restyInst.SetError(errorResponse{})

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

// TODO support source=registry on all calls (do all calls support?)
// TODO support synchronized on PUT calls

func ToEntityList(item *Entity) []Entity {
	var ea []Entity
	ea = append(ea, *item)
	return ea
}
