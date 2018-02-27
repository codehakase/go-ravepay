package ravepay

import (
	"encoding/json"
	"fmt"
	"log"
)

type CardService struct {
	Client *Client
}

// ChargeCard sends a requests to rave api to dictate validation flow
func (c *CardService) ChargeCard(data map[string]interface{}) ([]byte, error) {
	if err := keysExists([]string{"redirect_url"}, data); err != nil {
		return nil, err
	}
	payload := c.initialzeCard(data)
	resp, err := c.charge(data)
	if err != nil {
		return nil, err
	}
	// identify suggested auth method
	var d map[string]interface{}
	_ = jsonFromBytes(resp, d)
	authData := d["data"]
	auth := authData["data"].(map[string]interface{})["suggested_auth"].(string)

	if auth == "PIN" {
		data["suggested_auth"] = "PIN"
		err := keysExists([]string{"pin"}, data)
		if err != nil {
			return nil, err
		}

		cardSetup := c.initialzeCard(data)
		response, err := c.charge(data)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// ValidateCharge validates an account's charge with OTP
func (c *CardService) ValidateCharge(data map[string]interface{}) (*Response, error) {
	data["PBFPubKey"] = c.Client.Crypto.GetPublicKey()
	req := c.Client.Request.NewRequest("/flwv3-pug/getpaidx/api/validate", d)
	resp, err := req.DoRequest()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ValidateCardCharge validates a card with otp
func (c *CardService) ValidateCardCharge(data map[string]interface{}) (*Response, error) {
	data["PBFPubKey"] = c.Client.Crypto.GetPublicKey()
	req := c.Client.Request.NewRequest("/flwv3-pug/getpaidx/api/validatecharge", d)
	resp, err := req.DoRequest()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CardService) VerifyTransaction(data map[string]interface{}) (*Response, error) {
	if err := keysExists([]string{"flw_ref", "currency", "amount"}, data); err != nil {
		return nil, err
	}
	// SET secret key
	data["SECKEY"] = c.Client.Crypto.GetSecretKey()
	req := c.Client.Request.NewRequest("/flwv3-pug/getpaidx/api/verify", data)
	resp, err := req.DoRequest()
	if err != nil {
		return nil, err
	}

	// confirm verification
	// ref: https://flutterwavedevelopers.readme.io/v1.0/reference#verification
	successMessage := resp.getResponseMessage()
	dd := resp.ResponseData
	if ref, ok := dd["flw_ref"]; !ok { // revert to XRequery search
		ref = dd["flwref"]
	}

	// read meta data from response
	var (
		charge   string
		currency string
		amount   int
	)
	if responseMeta, ok := dd["flwMeta"]; !ok {
		charge = dd["chargecode"]
		currency = dd["currency"]
		amount = dd["chargedamount"]
	} else {
		charge = responseMeta.(map[string]interface{})["chargecode"].(string)
		currency = dd["transaction_currency"]
		amount = dd["charged_amount"]
	}
	transactionRef := data["flw_ref"]
	currencyCode := data["currency"]
	return resp, nil
}

/*
  Utils
*/

// charge a card
func (c *CardService) charge(cardData map[string]interface{}) (*Response, error) {
	req := c.Client.Request.NewRequest("/flwv3-pug/getpaidx/api/charge", cardData)
	resp, err := req.DoRequest()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *CardService) initialzeCard(data interface{}) map[string]interface{} {
	jb, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	data := string(jb[:])
	encrypted := c.Client.Crypto.TripleDesEncrypt(data)
	return map[string]interface{}{
		"PBFPubKey": c.Client.Crypto.GetPublicKey(),
		"client":    encrypted,
		"alg":       "3DES-24",
	}
}

func keysExists(keys []string, data map[string]interface{}) error {
	for _, k := range keys {
		if _, ok := data[key]; !ok {
			return fmt.Errorf("%s is required, and isn't set in payload", k)
		}
	}
	return nil
}

func jsonFromBytes(d []byte, v map[string]interface{}) error {
	err := json.Unmarshal(d, &v)
	if err != nil {
		log.Fatalln(err)
		return err
	}
}
