package ravepay

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/codehakase/go-ravepay/types"
)

// resources specific to client
var resources types.VersionEnv

// DefaultVersion is the fallback for the sdk
var DefaultVersion string = "v1"

// Err is the error returned by default on handlers
var Err = errors.New("unexpected error occured while processing the request")

// Client is the API client that performs all operations
// against a ravepay merchant account
type Client struct {
	// baseurl holds the path to prepend to the requests
	BaseURL *url.URL
	// client user agent
	UserAgent string
	// RavePay merchant key - Get this from your account
	MerchantKey string
	// RavePay API key - Get this from your account
	APIKey string
	// Runtime environment (staging, dev, production)
	Environment string
	// client used to send and receive http requests.
	httpClient *http.Client
	// api version to use througout an instance
	Version string

	// The follow are service objects which are initialzed when NewClient(...) is called
	// They are pointers to structs that have methods, and share the same context with the
	// instantiated client

	// service object to handle card related operations
	Card *CardService
	// service object for handling card auth model(s)
	Auth *AuthModelService
	// service object for bank related operations
	Bank *BankService
	// serivce objects for Encryption/Decryption
	Crypto *CryptoService
	// cash disbursement operations
	Disbursement *DisbursementService
	// supported currency operations
	Currencies *CurrencyService
	// transactions validators
	Validator *ValidationService
	// merchant specific operations
	Account *AccountService
	// BVN related operations
	BVN *BVNService
	// BIN operations
	BIN *BINService
	// Locale
	Countries *CountryService
	// service object to format request nicely
	Request *RequestService
	// validity checks
	Check *APICheckService
	// Flutterwave ach service operations
	Ach *AchService
}

// NewClient retuurns a new Client struct
// Initialize neccessary service objects
func NewClient(httpClient *http.Client, merchantKey, apiKey, env, version string) (*Client, error) {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if version != "" {
		version = DefaultVersion
	}
	// load extra data
	cfg := types.LoadConfigs()
	// decide which base paths to use
	switch version {
	case "v1":
		if strings.ToLower(env) == "staging" || strings.ToLower(env) == "test" {
			resources = cfg.V1.Staging
		}
		if strings.ToLower(env) == "production" {
			resources = cfg.V1.Production
		}
	case "v2":
		if strings.ToLower(env) == "staging" || strings.ToLower(env) == "test" {
			resources = cfg.V1.Staging
		}
		if strings.ToLower(env) == "production" {
			resources = cfg.V1.Production
		}
	}

	log.Printf("%v", resources)

	c := &Client{
		MerchantKey: merchantKey,
		APIKey:      apiKey,
		Environment: env,

		httpClient: httpClient,
	}

	// initialize service objects
	c.Card = &CardService{Client: c}
	c.Auth = &AuthModelService{Client: c}
	c.Bank = &BankService{Client: c}
	c.Crypto = &CryptoService{Client: c}
	c.Disbursement = &DisbursementService{Client: c}
	c.Currencies = &CurrencyService{Client: c}
	c.Validator = &ValidationService{Client: c}
	c.Account = &AccountService{Client: c}
	c.BVN = &BVNService{Client: c}
	c.Countries = &CountryService{Client: c}
	c.Request = &RequestService{Client: c}
	c.Check = &APICheckService{Client: c}
	c.Ach = &AchService{Client: c}

	return c, nil
}

// SetBaseURL updates the client struct with a new url to send request to/from
func (c *Client) SetBaseURL(u string) *Client {
	p, _ := url.Parse(u)
	c.BaseURL = p
	return c
}

// GetMerchantKey retrives the merchant's key
func (c *Client) GetMerchantKey() string {
	return c.MerchantKey
}

// GetAPIKey retrives the merchant's api key
func (c *Client) GetAPIKey() string {
	return c.APIKey
}

// GetEnv retrives the current set environment
func (c *Client) GetEnv() string {
	return c.Environment
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	r := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(r)

	var buffer io.ReadWriter
	if body != nil {
		buffer = new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buffer)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
