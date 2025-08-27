// Package gws provides an API client for the University of Washington Groups Service API
package gws

import (
	"crypto/tls"
	"errors"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// Config holds all the Client configuration
type Config struct {
	APIUrl        string
	Timeout       time.Duration
	Synchronized  bool // When true, API writes wait for cache propagation before returning
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
	once       sync.Once
	configErr  error
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
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}
	restyInst := resty.New()
	c := &Client{resty: restyInst, config: config}
	// Prepare static headers early
	restyInst.SetHeader("Accept", "application/json")
	restyInst.SetHeader("Content-Type", "application/json")
	// Defer TLS & the rest to lazy configure via sync.Once
	c.resty.SetTLSClientConfig(&tls.Config{Renegotiation: tls.RenegotiateFreelyAsClient})
	// Trigger initial configure so early errors are caught
	if _, err := c.healthCheckConfigure(); err != nil {
		return nil, err
	}
	return c, nil
}

// healthCheckConfigure performs a configure without issuing a request, useful to surface early errors.
func (client *Client) healthCheckConfigure() (bool, error) {
	client.configure()
	return client.configured, client.configErr
}

func (client *Client) configure() {
	client.once.Do(func() {
		restyInst := client.resty
		cfg := client.config
		// Basic settings
		restyInst.SetBaseURL(cfg.APIUrl)
		restyInst.SetTimeout(cfg.Timeout * time.Second)
		restyInst.SetError(errorResponse{})
		// TLS setup
		if cfg.CAFile != "" {
			restyInst.SetRootCertificate(cfg.CAFile)
		}
		if cfg.ClientCert != "" || cfg.ClientKey != "" {
			cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
			if err != nil {
				client.configErr = err
				return
			}
			restyInst.SetCertificates(cert)
		}
		restyInst.SetDebug(false)
		client.configured = true
	})
}

// SetTLSClientConfig assigns client TLS config
func (client *Client) SetTLSClientConfig(c *tls.Config) {
	client.resty.SetTLSClientConfig(c)
}

// request returns new resty.Request from configured client
func (client *Client) request() *resty.Request {
	client.configure()
	return client.resty.R()
}

// ConfigError returns any configuration error encountered during lazy initialization.
func (client *Client) ConfigError() error { return client.configErr }

// EnableSynchronized enables synchronized API operations. When enabled, write operations
// (create, update, delete) will not return until the changes have propagated to the
// API's read cache. This ensures that subsequent read operations will immediately see
// the changes, but may result in slower write operations.
func (client *Client) EnableSynchronized() {
	client.config.Synchronized = true
}

// DisableSynchronized disables synchronized API operations (default behavior).
// Write operations will return immediately after being processed, but changes may not
// be visible in read operations until the cache is updated. This provides better
// performance but requires applications to handle eventual consistency.
func (client *Client) DisableSynchronized() {
	client.config.Synchronized = false
}

// syncQueryString returns the desired synchronized mode as a query string.
func (client *Client) syncQueryString() string {
	if client.config.Synchronized {
		// Value doesn't matter, only presence/absence
		return "synchronized=true"
	}
	return ""
}
