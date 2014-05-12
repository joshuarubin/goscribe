package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
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

type callbackData struct {
	CallbackURL string `form:"callback_url" binding:"required"`
}

type transcribeData struct {
	callbackData
	AudioURL string `form:"audio_url" binding:"required"`
}

type transcribeUploadData struct {
	callbackData
	AudioData *multipart.FileHeader `form:"audio_data" binding:"required"`
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
		binding.ErrorHandler,
		handleTranscribe,
	)

	m.Post(
		"/v1/transcribe/process",
		strict.ContentType("application/x-www-form-urlencoded"),
		binding.Bind(telapi.TranscribeCallbackData{}),
		binding.ErrorHandler,
		handleTranscribeProcess,
	)

	m.Post(
		"/v1/transcribe/upload",
		auth.Basic(authUser, authPass),
		strict.Accept("application/json"),
		binding.MultipartForm(transcribeUploadData{}),
		binding.ErrorHandler,
		handleTranscribeUpload,
	)

	m.Router.NotFound(strict.MethodNotAllowed, strict.NotFound)
}

func main() {
	m.Run()
}

func success(r render.Render, data ...interface{}) {
	ret := map[string]interface{}{
		"status": http.StatusOK,
	}

	l := len(data)
	if l == 1 {
		ret["data"] = data[0]
	} else if l > 1 {
		ret["data"] = data
	}

	r.JSON(http.StatusOK, ret)
}

func jsonError(r render.Render, status int, err error) {
	r.JSON(status, map[string]interface{}{
		"status": status,
		"error":  err.Error(),
	})
}

func telapiError(r render.Render, err error) {
	if telapiError, ok := err.(telapi.Error); ok {
		r.JSON(telapiError.JSON())
		return
	}

	jsonError(r, http.StatusInternalServerError, err)
}

func handleTranscribe(data transcribeData, r render.Render) {
	resp, err := telapi.TranscribeURL(data.AudioURL, data.CallbackURL)
	if err != nil {
		telapiError(r, err)
		return
	}

	success(r, resp.TranscribeClientResponse.Translate())
}

func handleTranscribeProcess(data telapi.TranscribeCallbackData) (int, string) {
	pretty.Println(data)
	b, _ := json.Marshal(data.TranscribeCallbackClientData)
	pretty.Println(string(b))
	return http.StatusOK, ""
}

func handleTranscribeUpload(data transcribeUploadData, r render.Render) {
	file, err := data.AudioData.Open()
	if err != nil {
		jsonError(r, http.StatusInternalServerError, err)
		return
	}

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		jsonError(r, http.StatusInternalServerError, err)
		return
	}

	fmt.Println(fileData)

	success(r)
}
