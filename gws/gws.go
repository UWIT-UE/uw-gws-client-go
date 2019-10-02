package gws

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/http2"
)

// Config is used to configure the creation of the client.
type Config struct {
	// Address is the address of the GWS server. This should be a complete
	// URL such as "https://groups.uw.edu".
	Address string

	// HTTPClient is the HTTP client to use.
	HTTPClient *http.Client

	// Timeout is for setting custom timeout parameter in the HTTPClient
	Timeout time.Duration

	// If there is an error when creating the configuration, this will be the
	// error
	Error error

	// The Backoff function to use; a default is used if not provided
	Backoff retryablehttp.Backoff
}

// TLSConfig contains the parameters needed to configure TLS on the HTTP client
type TLSConfig struct {
	// CACert is the path to a PEM-encoded CA cert file to use to verify the
	// server SSL certificate.
	CACert string

	// CAPath is the path to a directory of PEM-encoded CA cert files to verify
	// the server SSL certificate.
	CAPath string

	// ClientCert is the path to the certificate for API authentication
	ClientCert string

	// ClientKey is the path to the private key for API authentication
	ClientKey string

	// TLSServerName, if set, is used to set the SNI host when connecting via
	// TLS.
	TLSServerName string

	// Insecure enables or disables SSL verification
	Insecure bool
}

// DefaultConfig returns a default configuration for the client. It is
// safe to modify the return value of this function.
//
// The default Address is https://groups.uw.edu/group_sws/v3, but this can be overridden by
// setting the `GWS_ADDR` environment variable.
//
// If an error is encountered, this will return nil.
func DefaultConfig() *Config {
	config := &Config{
		Address:    "https://groups.uw.edu/group_sws/v3",
		HTTPClient: &http.Client{},
	}
	config.HTTPClient.Timeout = time.Second * 60

	transport := config.HTTPClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		config.Error = err
		return config
	}

	if err := config.ReadEnvironment(); err != nil {
		config.Error = err
		return config
	}

	return config
}

// ConfigureTLS takes a set of TLS configurations and applies those to the the
// HTTP client.
func (c *Config) ConfigureTLS(t *TLSConfig) error {
	if c.HTTPClient == nil {
		c.HTTPClient = DefaultConfig().HTTPClient
	}
	clientTLSConfig := c.HTTPClient.Transport.(*http.Transport).TLSClientConfig

	var clientCert tls.Certificate
	foundClientCert := false

	switch {
	case t.ClientCert != "" && t.ClientKey != "":
		var err error
		clientCert, err = tls.LoadX509KeyPair(t.ClientCert, t.ClientKey)
		if err != nil {
			return err
		}
		foundClientCert = true
	case t.ClientCert != "" || t.ClientKey != "":
		return fmt.Errorf("both client cert and client key must be provided")
	}

	if t.CACert != "" {
		caCert, err := ioutil.ReadFile(t.CACert)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return fmt.Errorf("Error loading CA File: Couldn't parse PEM in: %s", caCert)
		}
		clientTLSConfig.RootCAs = caCertPool
	}

	if t.Insecure {
		clientTLSConfig.InsecureSkipVerify = true
	}

	if foundClientCert {
		// We use this function to ignore the server's preferential list of
		// CAs, otherwise any CA used for the cert auth backend must be in the
		// server's CA pool
		clientTLSConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &clientCert, nil
		}
	}

	if t.TLSServerName != "" {
		clientTLSConfig.ServerName = t.TLSServerName
	}

	return nil
}

// ReadEnvironment reads configuration information from the environment. If
// there is an error, no configuration value is updated.
func (c *Config) ReadEnvironment() error {
	var envAddress string
	var envCACert string
	var envCAPath string
	var envClientCert string
	var envClientKey string
	var envClientTimeout time.Duration
	var envInsecure bool
	var envTLSServerName string

	// Parse the environment variables
	if v := os.Getenv("GWS_ADDRESS"); v != "" {
		envAddress = v
	}
	if v := os.Getenv("GWS_CACERT"); v != "" {
		envCACert = v
	}
	if v := os.Getenv("GWS_CAPATH"); v != "" {
		envCAPath = v
	}
	if v := os.Getenv("GWS_CLIENT_CERT"); v != "" {
		envClientCert = v
	}
	if v := os.Getenv("GWS_CLIENT_KEY"); v != "" {
		envClientKey = v
	}
	if t := os.Getenv("GWS_CLIENT_TIMEOUT"); t != "" {
		clientTimeout, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse %q", "GWS_CLIENT_TIMEOUT")
		}
		envClientTimeout = time.Duration(clientTimeout) * time.Second
	}
	if v := os.Getenv("GWS_SKIP_VERIFY"); v != "" {
		var err error
		envInsecure, err = strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("could not parse VAULT_SKIP_VERIFY")
		}
	}
	if v := os.Getenv("GWS_TLS_SERVER_NAME"); v != "" {
		envTLSServerName = v
	}

	// Configure the HTTP clients TLS configuration.
	t := &TLSConfig{
		CACert:        envCACert,
		CAPath:        envCAPath,
		ClientCert:    envClientCert,
		ClientKey:     envClientKey,
		TLSServerName: envTLSServerName,
		Insecure:      envInsecure,
	}

	if err := c.ConfigureTLS(t); err != nil {
		return err
	}

	if envAddress != "" {
		c.Address = envAddress
	}

	if envClientTimeout != 0 {
		c.Timeout = envClientTimeout
	}

	return nil
}

// Client is the client to the Vault API. Create a client with NewClient.
type Client struct {
	addr    *url.URL
	Config  *Config
	headers http.Header
}

// NewClient returns a new client for the given configuration.
//
// If the configuration is nil, Vault will use configuration from
// DefaultConfig(), which is the recommended starting configuration.
//
// If the environment variable `VAULT_TOKEN` is present, the token will be
// automatically added to the client. Otherwise, you must manually call
// `SetToken()`.
func NewClient(c *Config) (*Client, error) {
	def := DefaultConfig()
	if def == nil {
		return nil, fmt.Errorf("could not create/read default configuration")
	}
	if def.Error != nil {
		// return nil, errwrap.Wrapf("error encountered setting up default configuration: {{err}}", def.Error)
		return nil, fmt.Errorf("error encountered setting up default configuration")
	}

	if c == nil {
		c = def
	}

	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	if c.HTTPClient == nil {
		c.HTTPClient = def.HTTPClient
	}
	if c.HTTPClient.Transport == nil {
		c.HTTPClient.Transport = def.HTTPClient.Transport
	}

	client := &Client{
		addr:   u,
		config: c,
	}

	return client, nil
}

// NewRequest creates a new raw request object to query the Vault server
// configured for this client. This is an advanced method and generally
// doesn't need to be called externally.
func (c *Client) NewRequest(method, requestPath string) *Request {
	addr := c.addr
	headers := c.headers

	req := &Request{
		Method: method,
		URL: &url.URL{
			User:   addr.User,
			Scheme: addr.Scheme,
			Host:   host,
			Path:   path.Join(addr.Path, requestPath),
		},
		Params: make(map[string][]string),
	}

	var lookupPath string
	switch {
	case strings.HasPrefix(requestPath, "/v3/"):
		lookupPath = strings.TrimPrefix(requestPath, "/v3/")
	case strings.HasPrefix(requestPath, "v3/"):
		lookupPath = strings.TrimPrefix(requestPath, "v3/")
	default:
		lookupPath = requestPath
	}

	if headers != nil {
		req.Headers = headers
	}

	return req
}
