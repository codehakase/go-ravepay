package main

import (
	"fmt"
	"log"

	rave "github.com/codehakase/go-ravepay"
)

func main() {
	var (
		merchantKey string = "abcdefghijklmnopqrstuvwyz"
		apiKey      string = "0202020$$sksksk"
		env         string = "dev"
	)

	client, err := rave.NewClient(nil, merchantKey, apiKey, env)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(client.GetMerchantKey())
}
