package main

import (
	"encoding/json"
	"os"

	"github.com/go-martini/martini"
	"github.com/joshuarubin/goscribe/telapi"
	"github.com/kr/pretty"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/strict"
)

var (
	m       *martini.ClassicMartini
	baseURL string
)

type transcribeData struct {
	CallbackURL string `form:"callback_url" binding:"required"`
	AudioURL    string `form:"audio_url" binding:"required"`
}

type transcribeClientTelAPIData struct {
	ID                string `form:"TranscriptionSid"    json:"id"`
	Status            string `form:"TranscriptionStatus" json:"status"`
	TranscriptionText string `form:"TranscriptionText"   json:"transcription_text"`
}

type transcribeTelAPIData struct {
	transcribeClientTelAPIData
	AudioURL             string  `form:"AudioUrl"`
	Duration             float32 `form:"Duration"`
	AccountSID           string  `form:"AccountSid"`
	APIVersion           string  `form:"ApiVersion"`
	Price                float32 `form:"Price"`
	TranscriptionQuality string  `form:"TranscriptionQuality"`
}

func init() {
	// BASE_URL is not required or else wercker tests would fail
	baseURL = os.Getenv("BASE_URL")

	m = martini.Classic()

	m.Use(gzip.All())
	m.Use(render.Renderer())

	m.Get("/", func() string {
		return "hello, world"
	})

	m.Post(
		"/v1/transcribe",
		strict.Accept("application/json"),
		strict.ContentType("application/x-www-form-urlencoded"),
		binding.Bind(transcribeData{}),
		handleTranscribe,
	)

	m.Post(
		"/v1/transcribe/process",
		strict.ContentType("application/x-www-form-urlencoded"),
		binding.Bind(transcribeTelAPIData{}),
		handleTranscribeProcess,
	)

	m.Router.NotFound(strict.MethodNotAllowed, strict.NotFound)
}

func main() {
	m.Run()
}

func telapiError(r render.Render, err error) {
	if telapiError, ok := err.(telapi.Error); ok {
		r.JSON(telapiError.JSON())
		return
	}

	r.JSON(500, map[string]interface{}{
		"status": 500,
		"error":  err.Error(),
	})
}

func handleTranscribe(data transcribeData, r render.Render) {
	resp, err := telapi.TranscribeURL(data.AudioURL, data.CallbackURL)
	if err != nil {
		telapiError(r, err)
		return
	}

	r.JSON(200, resp.TranscribeClientResponse.Translate())
}

func handleTranscribeProcess(data transcribeTelAPIData) (int, string) {
	pretty.Println(data)
	b, _ := json.Marshal(data.transcribeClientTelAPIData)
	pretty.Println(string(b))
	return 200, ""
}
