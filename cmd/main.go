//go:build !test
// +build !test

package main

import (
	"log"
)

func main() {
	log.Println("Starting Validate-Registry Admission Controller")
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}
