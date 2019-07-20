package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hello: %s\n", r.TLS.PeerCertificates[0].Subject.CommonName)
	fmt.Fprint(w, "Authentication successful")
}

func main() {
	var (
		err        error
		clientCert []byte
		pool       *x509.CertPool
		tlsConf    *tls.Config
		server     *http.Server
	)

	http.HandleFunc("/hello", helloHandler)

	if clientCert, err = ioutil.ReadFile("../client/clientCrt.pem"); err != nil {
		log.Fatalln(err)
	}
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(clientCert)

	tlsConf = &tls.Config{
		ClientCAs:  pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConf.BuildNameToCertificate()

	server = &http.Server{
		Addr:      ":9443",
		TLSConfig: tlsConf,
	}
	log.Fatalln(server.ListenAndServeTLS("serverCrt.pem", "serverKey.pem"))
}
