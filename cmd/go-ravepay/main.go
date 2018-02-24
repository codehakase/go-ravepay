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
		env         string = "production"
		version     string = "v2"
	)

	client, err := rave.NewClient(nil, merchantKey, apiKey, env, version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(client.GetMerchantKey())
	fmt.Println("\n")
	fmt.Printf("%+v\n", client)
}
