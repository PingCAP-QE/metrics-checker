package metric

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Check checks if a query returns true.
func Check(client v1.API, query string, ts time.Time) (ans bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	val, _, err := client.Query(ctx, query, ts)
	if err != nil {
		return false, err
	}
	if val.Type() == model.ValVector {
		boolVector := val.(model.Vector)
		if boolVector.Len() == 0 {
			return false, errors.New("Prometheus is not up")
		}
		if boolVector[0].Value == 1 {
			return true, nil
		}
	}
	return false, errors.New("return type is not model.ValVector")
}

// AddHTTPIfIP add "http://" before ip address like 127.0.0.1 or 127.0.0.1:3000
func AddHTTPIfIP(address string) (string, error) {
	prefix := "http://"
	if !strings.HasPrefix(address, prefix) {
		address = prefix + address
	}
	u, err := url.Parse(address)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
