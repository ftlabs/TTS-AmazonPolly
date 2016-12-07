package main

import (
	"github.com/golang/go/src/pkg/net/http"
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/kennygrant/sanitize"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/url"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"time"
	"net"
	"io/ioutil"
	"io"
	"regexp"
	//"encoding/json"
	"github.com/aws/aws-sdk-go/service/polly"
	//log "github.com/Sirupsen/logrus"
	"errors"
)

var UUIDRegexp = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
var outPutFormat string = "mp3"

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 128,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	},
}

type textToSpeechService interface {
	convertToSpeech(thing interface{}) (*polly.SynthesizeSpeechOutput, error)
}

type textToSpeechServiceImpl struct {
	apiKey 	   string
	contentUrl url.URL
	awsCreds   credentials.Credentials
	token      string
}

func newTextToSpeechService(apiKey string, contentUrl url.URL, awsCreds credentials.Credentials, token string) (textToSpeechService) {
	return &textToSpeechServiceImpl{apiKey: apiKey, contentUrl: contentUrl, awsCreds: awsCreds, token: token}
}

func (tts *textToSpeechServiceImpl) convertToSpeech(thing interface{}) (*polly.SynthesizeSpeechOutput, error) {
	input := thing.(request)

	if (tts.token != input.Token) {
		return &polly.SynthesizeSpeechOutput{}, errors.New("Token " + input.Token + " is invalid!")
	}
	var textToConvert string
	//if (UUIDRegexp.MatchString(input.Body)) {
	//	textToConvert = getContentBody(tts.contentUrl, tts.apiKey, input.Body)
	//} else {
		textToConvert = input.Body
	//}

	if (len(textToConvert) == 0 || textToConvert == "") {
		return &polly.SynthesizeSpeechOutput{}, errors.New("Unable to process input text:" + input.Body)
	}

	sess, err := session.NewSession()
	if err != nil {
		return &polly.SynthesizeSpeechOutput{}, fmt.Errorf("Failed to create session", err)
	}

	pollyService := polly.New(sess, &aws.Config{Region: aws.String("eu-west-1"), Credentials: &tts.awsCreds})

	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(outPutFormat),
		Text: 	      aws.String(textToConvert),
		VoiceId:      aws.String(input.VoiceId),
	}

	ssiResp, err := pollyService.SynthesizeSpeech(params)

	if err != nil {
		return &polly.SynthesizeSpeechOutput{},	fmt.Errorf("Unable to access Amazon Polly %v", err)
	}
	return ssiResp, nil
}

func getContentBody(contentUrl url.URL, apiKey string, uuid string) (contentNoTags string) {
	req, _ := http.NewRequest("GET", contentUrl.String()+uuid+"?apiKey="+apiKey, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Errorf("Could not connect to %v. Error (%v)", contentUrl, err)
		return ""
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	//input := content{}
	//bodyXml := json.NewDecoder(resp.Body)
	//something := bodyXml.Decode(input.BodyXml)
	//contentNoTags = sanitize.HTML(something)
	//fmt.Println(contentNoTags)
	return contentNoTags
}