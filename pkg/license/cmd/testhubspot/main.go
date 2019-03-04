package main

import (
	"context"
	"fmt"
	"os"

	"github.com/solo-io/solo-projects/pkg/license/notify/hubspot"
)

func main() {
	k := os.Getenv("APIKEY")
	if len(k) == 0 {
		panic("please provide api key")
	}
	const lk = "license_key"
	ctx := context.Background()
	h := hubspot.NewHubspot(k)
	email := "rick@solo.io"
	rick := hubspot.ContactUpdate{
		Properties: []hubspot.PropertyUpdate{{
			Property: lk,
			Value:    "123",
		}},
	}
	ur, err := h.UpsertContact(ctx, email, rick)
	if err != nil {
		panic(err)
	}

	fmt.Println(ur)

	c, err := h.GetContact(context.Background(), email, lk)
	if err != nil {
		panic(err)
	}

	fmt.Println(c)
}
