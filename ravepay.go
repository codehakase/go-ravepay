package ravepay

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/codehakase/go-ravepay/types"
)

// DefaultVersion is the fallback for the sdk
const (
	DefaultVersion = "v1"
	ProductionURL  = "https://api.ravepay.co"
	StagingURL     = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com"
)

// Err is the error returned by default on handlers
var Err = errors.New("unexpected error occured while processing the request")

// Client is the API client that performs all operations
// against a ravepay merchant account
type Client struct {
	// baseurl holds the path to prepend to the requests
	BaseURL *url.URL
	// client user agent
	UserAgent string
	// RavePay API Public key - Get this from your account
	PublicKey string
	// RavePay secret key - this is initialized and is recommended to be stored in an environment variable
	SecretKey string
	// Runtime environment (staging, or production)
	Environment string
	// client used to send and receive http requests.
	httpClient *http.Client
	// api version to use througout an instance
	Version string
	// Resources used by all service objects
	Resources *types.Resources
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
	if version == "" {
		version = DefaultVersion
	}

	c := Client{
		MerchantKey: merchantKey,
		APIKey:      apiKey,
		Environment: env,
		Resources:   types.LoadConfigs(),
		httpClient:  httpClient,
	}
	// set mode
	switch env {
	case "staging":
		c.BaseURL = url.Parse(StagingURL)
	case "production":
		c.BaseURL = url.Parse(ProductionURL)
	default:
		c.BaseURL = url.Parse(StagingURL)
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

	return &c, nil
}

// SetBaseURL updates the client struct with a new url to send request to/from
func (c *Client) SetBaseURL(u string) *Client {
	p, _ := url.Parse(u)
	c.BaseURL = p
	return c
}

// GetEnv retrives the current set environment
func (c *Client) GetEnv() string {
	return c.Environment
}

func (c *Client) newRequest(method, path string, body interface{}, headers map[string]string) (*http.Request, error) {
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
		if headers != nil {
			for k, v := range headers {
				req.Header.Set(headers[k], headers[v])
			}
		}
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
