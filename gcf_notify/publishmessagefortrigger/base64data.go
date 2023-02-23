package main

import (
	"encoding/base64"
	"log"
)

func main() {
	encoded := base64Encode(`{"highest_delta_deposits_by_date":"[delta depositsByDate]\n[147.37 Wed,02/22/20]\n[1870.57 Wed,02/21/20]\n[1870.57 Wed,02/20/20]\n[147.37 Wed,02/19/20]\n[1870.57 Wed,02/18/20]\n"}`)
	log.Println("encoded :", encoded)
}

// function to encode string to base64 encoding
// log.Println("base64 encoded string for local", base64Encode(dataRecvd))
// {"highest_delta_deposits_by_date":"[delta depositsByDate]\n[147.37 Wed,02/22/23]\n[1870.57 Wed,02/21/23]\n[1870.57 Wed,02/20/23]\n[147.37 Wed,02/19/23]\n[1870.57 Wed,02/18/23]\n"}
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
