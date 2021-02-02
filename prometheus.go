package main

import (
	"context"
	"log"
	"net"
	"net/url"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Check checks if a query returns true.
func Check(client v1.API, query string, ts time.Time) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Printf("checking query: %s", query)
	val, _, err := client.Query(ctx, query, ts)
	if err != nil {
		log.Fatalf("prometheus client is unavailable")
	}
	if val.Type() == model.ValVector {
		boolVector := val.(model.Vector)
		// TODO: When Prometheus doesn't up, the length of boolVector would be zero.
		// 		 This is a temporary fix. Raise a proper error in the future.
		if boolVector.Len() == 0 {
			// TODO: During startup of prometheus, query would return a zero-length
			// 		 vector. Handle it more properly in the future.
			return true
		}
		if boolVector[0].Value == 1 {
			return true
		}
	}
	return false
}

// AddHTTPIfIP add "http://" before ip address like 127.0.0.1 or 127.0.0.1:3000
func AddHTTPIfIP(address string) string {
	if net.ParseIP(address) != nil {
		return "http://" + address
	}
	_, err := url.Parse(address)
	if err == nil {
		return address
	}
	_, _, err = net.SplitHostPort(address)
	if err == nil {
		return "http://" + address
	}
	return address
}
