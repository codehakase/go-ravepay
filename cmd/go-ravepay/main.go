package main

import (
	"encoding/json"
	"fmt"
	"log"

	rave "github.com/codehakase/go-ravepay"
)

func main() {
	var (
		merchantKey string = "abcdefghijklmnopqrstuvwyz"
		apiKey      string = "0202020$$sksksk"
		env         string = "production"
		version     string = "v2"
	)

	client, err := rave.NewClient(nil, merchantKey, apiKey, env, version)

	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{}
	key := "e634d14d9ded04eaf05d5b67"
	d, _ := json.Marshal(data)
	encryptedData, err := client.Crypto.TripleDesEncrypt(string(d), key)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("encrypted map: \n")
	fmt.Println(encryptedData)

	decrypted, _ := client.Crypto.TripleDesDecrypt("b967a0d5457cb82f35c434d761acafd619bc274337e099c9c577cbe059550d6fb35fde39cd0dcb377208c48cc58f6a7bf337a11cdcd172d9408640f37b77addcbc6dc0717f49f5b54d87e5ffdcb7890d", key)
	fmt.Println("decrypted str: \n")
	fmt.Println(decrypted)
}
