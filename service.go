package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

type textToSpeechService interface {
	convertToSpeech(thing interface{}) (*polly.SynthesizeSpeechOutput, error)
}

type textToSpeechServiceImpl struct {
	awsCreds credentials.Credentials
	token    string
}

func newTextToSpeechService(awsCreds credentials.Credentials, token string) textToSpeechService {
	return &textToSpeechServiceImpl{awsCreds: awsCreds, token: token}
}

func (tts *textToSpeechServiceImpl) convertToSpeech(thing interface{}) (*polly.SynthesizeSpeechOutput, error) {
	input := thing.(request)

	if tts.token != input.Token {
		return &polly.SynthesizeSpeechOutput{}, errors.New("Token " + input.Token + " is invalid!")
	}

	if len(input.Body) == 0 || input.Body == "" {
		return &polly.SynthesizeSpeechOutput{}, errors.New("Unable to process input text:" + input.Body)
	}

	sess, err := session.NewSession()
	if err != nil {
		return &polly.SynthesizeSpeechOutput{}, fmt.Errorf("Failed to create session", err)
	}

	pollyService := polly.New(sess, &aws.Config{Region: aws.String("eu-west-1"), Credentials: &tts.awsCreds})

	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String(input.Body),
		VoiceId:      aws.String(input.VoiceId),
	}

	ssiResp, err := pollyService.SynthesizeSpeech(params)

	if err != nil {
		return &polly.SynthesizeSpeechOutput{}, fmt.Errorf("Unable to access Amazon Polly %v", err)
	}
	return ssiResp, nil
}
