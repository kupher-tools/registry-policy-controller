package main

import (
	"log"
	"net/http"
	"registry-policy-controller/internal"
)

func main() {
	http.HandleFunc("/validate-registry", internal.ValidateRegistry)
	log.Println("Starting Validate-Registry Admission Controller")
	log.Fatal(http.ListenAndServeTLS(":8443", "/tls/tls.crt", "/tls/tls.key", nil))

}
