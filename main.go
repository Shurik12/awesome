package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Shurik12/awesome/admin"
	"github.com/Shurik12/awesome/config"
)

func main() {

	// Read and process config file =========================
	configPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Auth token: %s\n", config.Awesome.Auth_token)
	// ======================================================

	// CMS server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	cmsMux := admin.InitApp()
	cmsServer := &http.Server{
		Addr:    ":" + port,
		Handler: cmsMux,
	}
	go cmsServer.ListenAndServe()
	fmt.Println("CMS Served at http://localhost:" + port + "/admin")

	// Publish server
	u, _ := url.Parse(os.Getenv("PUBLISH_URL"))
	publishPort := u.Port()
	if publishPort == "" {
		publishPort = "9001"
	}
	publishMux := http.FileServer(http.Dir(admin.PublishDir))
	publishServer := &http.Server{
		Addr:    ":" + publishPort,
		Handler: publishMux,
	}
	fmt.Println("Publish Served at http://localhost:" + publishPort)
	log.Fatal(publishServer.ListenAndServe())
}
