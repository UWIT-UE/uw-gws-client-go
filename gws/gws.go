// Package gws provides an API client for the University of Washington Groups Service API
package gws

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

// Config holds all the Client configuration
type Config struct {
	APIUrl        string
	Timeout       time.Duration
	Synchronized  bool // API writes wait for cache update
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
		Synchronized:  false,
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
	restyInst.SetDebug(false)
	client.configured = true
}

// SetTLSClientConfig assigns client TLS config
func (client *Client) SetTLSClientConfig(c *tls.Config) {
	client.resty.SetTLSClientConfig(c)
}

// request returns new resty.Request from configured client
func (client *Client) request() *resty.Request {
	client.configure()
	request := client.resty.R()

	return request
}

// EnableSynchronized enables synchronized API operation, waiting for writes to propagate to cache.
func (client *Client) EnableSynchronized() {
	client.config.Synchronized = true
}

// DisableSynchronized disables synchronized API operation, not waiting for writes to propagate to cache. This is the API default.
func (client *Client) DisableSynchronized() {
	client.config.Synchronized = false
}

// syncQueryString returns the desired sychronized mode as a query string.
func (client *Client) syncQueryString() string {
	if client.config.Synchronized {
		// Value doesn't matter, only presence/absence
		return "synchronized=true"
	}
	return ""
}

// TODO support synchronized on PUT affiliate

