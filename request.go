package ravepay

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestService struct {
	Client  *Client
	BaseURL string
	Data    map[string]string
}

// Response is populated and sent on each request
type Response struct {
	StatusCode           int                    `json:"statusCode"`
	ResponseCode         string                 `json:"responseCode"`
	ResponseMessage      string                 `json:"responseMessage"`
	ResponseData         map[string]interface{} `json:"responseData"`
	RequiresValidation   bool                   `json:"requiresValidation"`
	IsSuccessfulResponse bool                   `json:"isSuccessfulResponse"`
}

// return server response if status code is not 500
func (res *Response) getResponseCode() string {
	return res.ResponseCode
}

// return the server response status code
func (res *Response) getStatusCode() int {
	return res.StatusCode
}

// confirm response was sent
func (res *Response) isSuccessfulResponse() bool {
	return res.IsSuccessfulResponse
}

/**
 * some request to server will be successful but requires
 * a validation step to complete action. use requiresValidation to
 * check if you need to validate
 */
func (res *Response) requiresValidation() bool {
	return res.RequiresValidation
}

// return the request's response message
func (res *Response) getResponseMessage() string {
	return res.ResponseMessage
}

// return response data
func (res *Response) getResponseData() map[string]interface{} {
	return res.ResponseData
}

// NewRequest acts as the constructor to make subsequent requests from the sdk
func (r *RequestService) NewRequest(url string) *RequestService {
	r.BaseURL = url
	r.Data = make(map[string]string)
	return r
}

func (r *RequestService) AddBody(key, value string) {
	r.Data[key] = value
}

func (r *RequestService) GetBody() map[string]string {
	return r.Data
}

func (r *RequestService) DoRequest() (*Response, error) {
	rs := Response{}
	r.Client.httpClient = &http.Client{Timeout: 60}
	req, err := r.Client.newRequest("POST", r.BaseURL, r.Data, nil)
	if err != nil {
		return &rs, err
	}
	resp, err := r.Client.do(req, &rs)
	if err != nil {
		return &rs, err
	}

	return r.parseResponse(resp), nil
}

func (r *RequestService) parseResponse(response *http.Response) *Response {
	resp := Response{}
	data := make(map[string]interface{})
	resp.StatusCode = response.StatusCode
	if response.StatusCode < 500 {
		resBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		if err := json.Unmarshal(resBytes, &data); err != nil {
			log.Fatalln("Error: %s\n", err)
		}
		resp.ResponseData = data
		if data["data"].(map[string]interface{})["responseCode"] != "" {
			resp.ResponseCode = data["data"].(map[string]interface{})["responseCode"].(string)
		}
		if data["data"].(map[string]interface{})["responsecode"] != "" {
			resp.ResponseCode = data["data"].(map[string]interface{})["responsecode"].(string)
		}

		if resp.ResponseCode != "" && resp.ResponseCode == "02" {
			resp.RequiresValidation = true
		}

		if data["status"] == "success" && resp.StatusCode == 200 && (resp.ResponseCode == "00" || resp.ResponseCode == "0" || resp.ResponseCode == "02") {
			resp.IsSuccessfulResponse = true
		}

		if data["data"].(map[string]interface{})["resposemessage"] != "" {
			resp.ResponseMessage = data["data"].(map[string]interface{})["responsemessage"].(string)
		}
	}
	return &resp
}
