//go:build !test
// +build !test

package main

import (
	"net/http"
	"registry-policy-controller/internal"
)

func startServer() error {
	http.HandleFunc("/validate-registry", internal.ValidateRegistry)
	return http.ListenAndServeTLS(":8443", "/tls/tls.crt", "/tls/tls.key", nil)
}
