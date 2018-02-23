package ravepay

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
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
	// RavePay merchant key - Get this from your account
	MerchantKey string
	// RavePay API key - Get this from your account
	APIKey string
	// Runtime environment (staging, dev, production)
	Environment string
	// client used to send and receive http requests.
	httpClient *http.Client

	Card         *CardService
	Auth         *AuthModelService
	Bank         *BankService
	Crypto       *CryptoService
	Disbursement *DisbursementService
	Currencies   *CurrencyService
	Validator    *ValidationService
	Account      *AccountService
	BVN          *BVNService
	BIN          *BINService
	Countries    *CountryService
	Request      *RequestService
	Check        *APICheckService
	Ach          *AchService
}

func NewClient(httpClient *http.Client, merchantKey, apiKey, env string) (*Client, error) {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		MerchantKey: merchantKey,
		APIKey:      apiKey,
		Environment: env,

		httpClient: httpClient,
	}

	// pass client to service objects
	c.Card = &CardService{client: c}
	c.Auth = &AuthModelService{client: c}
	c.Bank = &BankService{client: c}
	c.Crypto = &CryptoService{client: c}
	c.Disbursement = &DisbursementService{client: c}
	c.Currencies = &CurrencyService{client: c}
	c.Validator = &ValidationService{client: c}
	c.Account = &AccountService{client: c}
	c.BVN = &BVNService{client: c}
	c.Countries = &CountryService{client: c}
	c.Request = &RequestService{client: c}
	c.Check = *APICheckService{client: c}
	c.AchService = &AchService{client: c}
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
u:
	c.BaseURL.ResolveReference(rel)

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

		return req, nil
	}
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
