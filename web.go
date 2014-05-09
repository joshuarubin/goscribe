package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-martini/martini"
	"github.com/kr/pretty"
	"github.com/martini-contrib/gzip"
)

const transcribeURLPattern = "https://%s:%s@%s/v1/Accounts/%s/Transcriptions.json"

var (
	m                *martini.ClassicMartini
	telapiBaseHost   string
	telapiAccountSid string
	telapiAuthToken  string
	baseURL          string
)

func init() {
	// BASE_URL is not required or else wercker tests would fail
	baseURL = os.Getenv("BASE_URL")

	if telapiBaseHost = os.Getenv("TELAPI_BASE_HOST"); telapiBaseHost == "" {
		telapiBaseHost = "api.telapi.com"
	}

	if telapiAccountSid = os.Getenv("TELAPI_ACCOUNT_SID"); telapiAccountSid == "" {
		log.Fatalln("TELAPI_ACCOUNT_SID is not set")
	}

	if telapiAuthToken = os.Getenv("TELAPI_AUTH_TOKEN"); telapiAuthToken == "" {
		log.Fatalln("TELAPI_AUTH_TOKEN is not set")
	}

	m = martini.Classic()

	m.Use(gzip.All())

	m.Get("/", func() string {
		return "hello, world"
	})
}

func main() {
	m.Run()
}

type transcribeResponse struct {
	Status             int
	Code               int
	SID                string
	DateCreated        string
	DateUpdated        string
	AccountSID         string
	Type               string
	AudioURL           string
	Duration           string
	TranscriptionText  string
	APIVersion         string
	Price              string
	TranscribeCallback string
	CallbackMethod     string
	URI                string
	Message            string
	MoreInfo           string
}

func getTranscription(audioURL string) error {
	transcribeURL := fmt.Sprintf(transcribeURLPattern, telapiAccountSid, telapiAuthToken, telapiBaseHost, telapiAccountSid)
	fmt.Println("transcribeURL", transcribeURL)

	// TODO(jrubin) change user agent

	data := url.Values{
		"AudioUrl":           {audioURL},
		"TranscribeCallback": {baseURL + "/process-transcription"},
	}

	pretty.Println("POST", data)

	resp, err := http.PostForm(transcribeURL, data)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("status", resp.StatusCode, resp.Status)
	pretty.Println("header", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("body", string(body))

	return nil
}
