package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"net/http"
	"os"
)

func main() {
	app := cli.App("text-to-speech", "A RESTful API for interracting with Amazon Polly, Text to Speech")
	accessId := app.String(cli.StringOpt{
		Name:   "aws-access-id",
		Value:  "",
		Desc:   "Aws ccess id",
		EnvVar: "AWS_ACCESS_ID",
	})
	accessKey := app.String(cli.StringOpt{
		Name:   "aws-access-key",
		Value:  "",
		Desc:   "Aws ccess key",
		EnvVar: "AWS_ACCESS_KEY",
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
		creds := credentials.NewStaticCredentials(*accessId, *accessKey, "")

		s := newTextToSpeechService(*creds, *userToken)
		h := newTextToSpeechHandler(s)

		m := mux.NewRouter()
		http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, m))
		m.HandleFunc("/convert", h.convertToSpeech).Methods("PUT")

		log.Infof("Listening on [%v]", *port)
		http.ListenAndServe(":"+*port, nil)
	}
	app.Run(os.Args)

}
