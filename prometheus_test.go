package main

import (
	"testing"
)

func TestAddHTTPIfIP(t *testing.T) {
	cases := map[string]string{
		"127.0.0.1":                    "http://127.0.0.1",
		"http://127.0.0.1":             "http://127.0.0.1",
		"127.0.0.1:8080":               "http://127.0.0.1:8080",
		"::1":                          "http://::1",
		"2001:4860:0:2001::68":         "http://2001:4860:0:2001::68",         // IPv6
		"[1fff:0:a88:85a3::ac1f]:8001": "http://[1fff:0:a88:85a3::ac1f]:8001", // IPv6 with port
		"jlkfada":                      "jlkfada",                             // Some messy things
	}
	for c, ans := range cases {
		if AddHTTPIfIP(c) != ans {
			t.Errorf("AddHTTPIfIP(%s) expect %s but get %s", c, ans, AddHTTPIfIP(c))
		}
	}
}
