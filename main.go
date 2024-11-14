package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Shurik12/awesomethingsshop/admin"
)

var (
	CertFilePath = "./server.crt"
	KeyFilePath  = "./server.key"
)

func main() {
	// CMS server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// load tls certificates
	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}

	// Configurre server to trust TLS client cert issued by your CA.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	cmsMux := admin.InitApp()
	cmsServer := &http.Server{
		Addr:      ":" + port,
		Handler:   cmsMux,
		TLSConfig: tlsConfig,
	}
	// go cmsServer.ListenAndServe()
	go cmsServer.ListenAndServeTLS(CertFilePath, KeyFilePath)
	fmt.Println("CMS Served at https://localhost:" + port + "/admin")

	// Publish server
	u, _ := url.Parse(os.Getenv("PUBLISH_URL"))
	publishPort := u.Port()
	if publishPort == "" {
		publishPort = "9001"
	}
	publishMux := http.FileServer(http.Dir(admin.PublishDir))
	publishServer := &http.Server{
		Addr:      ":" + publishPort,
		Handler:   publishMux,
		TLSConfig: tlsConfig,
	}
	fmt.Println("Publish Served at https://localhost:" + publishPort)
	// log.Fatal(publishServer.ListenAndServe())
	log.Fatal(publishServer.ListenAndServeTLS(CertFilePath, KeyFilePath))
}
