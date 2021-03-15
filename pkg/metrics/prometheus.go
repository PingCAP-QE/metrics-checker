package metrics

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
	// Ref: https://github.com/prom[<8;52;14metheus/prometheus/blob/76750d2a96df54226e85ac272d7ad5a547630240/rules/manager.go#L186-L206
	switch v := val.(type) {
	case model.Vector:
		return v.Len() > 0, nil
	case *model.Scalar:
		return v != nil, nil
	default:
		return false, errors.New("rule result is not a vector or scalar")
	}
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
