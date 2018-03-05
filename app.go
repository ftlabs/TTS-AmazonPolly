package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
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

	if *accessId == "" {
		log.Error("env: AWS_ACCESS_ID not set.")
		return
	}

	if *accessKey == "" {
		log.Error("env: AWS_ACCESS_KEY not set.")
		return
	}

	if *userToken == "" {
		log.Error("env: TOKEN not set.")
		return
	}

	if *port == "" {
		log.Error("env: PORT not set.")
		return
	}

	app.Action = func() {
		// NB, could instead do this, but need to rename env vars to match default names
		// creds := credentials.NewEnvCredentials()
		// if _, err := creds.Get(); err != nil {
		// 				log.WithError(err).Error("AWS credentials not set.")
		// 				return
		// }

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
