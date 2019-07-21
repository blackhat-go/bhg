package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bhg/ch-10/plugin-core/scanner"
)

var Users = []string{"admin", "manager", "tomcat"}
var Passwords = []string{"admin", "manager", "tomcat", "password"}

// TomcatChecker implements the scanner.Check interface. Used for guessing Tomcat creds
type TomcatChecker struct{}

// Check attempts to identify guessable Tomcat credentials
func (c *TomcatChecker) Check(host string, port uint64) *scanner.Result {
	var (
		resp   *http.Response
		err    error
		url    string
		res    *scanner.Result
		client *http.Client
		req    *http.Request
	)
	log.Println("Checking for Tomcat Manager...")
	res = new(scanner.Result)
	url = fmt.Sprintf("http://%s:%d/manager/html", host, port)
	if resp, err = http.Head(url); err != nil {
		log.Printf("HEAD request failed: %s\n", err)
		return res
	}
	log.Println("Host responded to /manager/html request")
	// Got a response back, check if authentication required
	if resp.StatusCode != http.StatusUnauthorized || resp.Header.Get("WWW-Authenticate") == "" {
		log.Println("Target doesn't appear to require Basic auth.")
		return res
	}

	// Appears authentication is required. Assuming Tomcat manager. Guess passwords...
	log.Println("Host requires authentication. Proceeding with password guessing...")
	client = new(http.Client)
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Println("Unable to build GET request")
		return res
	}
	for _, user := range Users {
		for _, password := range Passwords {
			req.SetBasicAuth(user, password)
			if resp, err = client.Do(req); err != nil {
				log.Println("Unable to send GET request")
				continue
			}
			if resp.StatusCode == http.StatusOK {
				res.Vulnerable = true
				res.Details = fmt.Sprintf("Valid credentials found - %s:%s", user, password)
				return res
			}
		}
	}
	return res
}

// New is the entry point required by the scanner
func New() scanner.Checker {
	return new(TomcatChecker)
}
