package main

import (
	"os"
	"github.com/jawher/mow.cli"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gorilla/mux"
	log "github.com/Sirupsen/logrus"
	"net/url"
	"net/http"
	"github.com/gorilla/handlers"
)

func main() {
	app := cli.App("text-to-speech", "A RESTful API for interracting with Amazon Polly, Text to Speech")
	apiKey := app.String(cli.StringOpt{
		Name:   "apiKey",
		Value:  "",
		Desc:   "Api Key for Capi v2 auth",
		EnvVar: "API_KEY",
	})
	contentAddr := app.String(cli.StringOpt{
		Name:   "contentAddr",
		Value:  "",
		Desc:   "Address to get content from Capi V2",
		EnvVar: "CONTENT_ADDR",
	})
	awsCredsFile := app.String(cli.StringOpt{
		Name:   "awsCreds",
		Value:  "",
		Desc:   "File which contains credentials for accessing Polly",
		EnvVar: "CREDS_FILE",
	})
	userToken := app.String(cli.StringOpt{
		Name:   "userToken",
		Value:  "",
		Desc:   "Token for accessing app",
		EnvVar: "TOKEN",
	})
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})

	app.Action = func() {
		creds := credentials.NewSharedCredentials(*awsCredsFile, "default")
		contentUrl, err := url.Parse(*contentAddr)
		if err != nil {
			log.Fatalf("Invalid content URL: %v (%v)", *contentAddr, err)
		}

		s := newTextToSpeechService(*apiKey, *contentUrl, *creds, *userToken)
		h := newTextToSpeechHandler(s)

		m := mux.NewRouter()
		http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, m))
		m.HandleFunc("/convert", h.convertToSpeech).Methods("PUT")

		log.Infof("Listening on [%v]", *port)
		http.ListenAndServe(":" + *port, nil)
	}
	app.Run(os.Args)

}
