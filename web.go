package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	"github.com/joshuarubin/goscribe/telapi"
	"github.com/kr/pretty"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/strict"
)

// TODO(jrubin):
// * post file to s3 and use that url for transcriptions
// * mgo (http://labix.org/mgo)
// * sessions (https://github.com/martini-contrib/sessions)
// * oauth2 (https://github.com/martini-contrib/oauth2)
// * accessflags (https://github.com/martini-contrib/accessflags)
// * remove auth, replace with api key validator (github.com/martini-contrib/auth)
// * CORS support would be nice (https://github.com/martini-contrib/cors)
// * CSRF (https://github.com/martini-contrib/csrf)
// * secure (https://github.com/martini-contrib/secure)

var (
	m        *martini.ClassicMartini
	baseURL  string
	authUser string
	authPass string
)

type transcribeData struct {
	CallbackURL string `form:"callback_url" binding:"required"`
	AudioURL    string `form:"audio_url" binding:"required"`
}

func init() {
	// BASE_URL, AUTH_USER and AUTH_PASS are not required or else wercker tests would fail
	baseURL = os.Getenv("BASE_URL")
	authUser = os.Getenv("AUTH_USER")
	authPass = os.Getenv("AUTH_PASS")

	m = martini.Classic()

	m.Use(gzip.All())
	m.Use(render.Renderer())

	m.Get("/", func() string {
		return "hello, world"
	})

	m.Post(
		"/v1/transcribe",
		auth.Basic(authUser, authPass),
		strict.Accept("application/json"),
		strict.ContentType("application/x-www-form-urlencoded"),
		binding.Bind(transcribeData{}),
		handleTranscribe,
	)

	m.Post(
		"/v1/transcribe/process",
		strict.ContentType("application/x-www-form-urlencoded"),
		binding.Bind(telapi.TranscribeCallbackData{}),
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

	r.JSON(http.StatusInternalServerError, map[string]interface{}{
		"status": http.StatusInternalServerError,
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

func handleTranscribeProcess(data telapi.TranscribeCallbackData) (int, string) {
	pretty.Println(data)
	b, _ := json.Marshal(data.TranscribeCallbackClientData)
	pretty.Println(string(b))
	return http.StatusOK, ""
}
