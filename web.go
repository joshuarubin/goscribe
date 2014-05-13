package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

// TODO(jrubin):
// * mgo (http://labix.org/mgo)
// * sessions (https://github.com/martini-contrib/sessions)
// * oauth2 (https://github.com/martini-contrib/oauth2)
// * accessflags (https://github.com/martini-contrib/accessflags)
// * remove auth, replace with api key validator (github.com/martini-contrib/auth)
// * CORS support would be nice (https://github.com/martini-contrib/cors)
// * CSRF (https://github.com/martini-contrib/csrf)
// * secure (https://github.com/martini-contrib/secure)

var (
	m         *martini.ClassicMartini
	s3Bucket  *s3.Bucket
	s3BaseURL string
	baseURL   string
	authUser  string
	authPass  string
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
	// BASE_URL, AUTH_USER and AUTH_PASS, AWS_S3_BASE_URL are not required or else wercker tests would fail
	baseURL = os.Getenv("BASE_URL")
	authUser = os.Getenv("AUTH_USER")
	authPass = os.Getenv("AUTH_PASS")
	s3BaseURL = os.Getenv("AWS_S3_BASE_URL")

	if awsAuth, err := aws.EnvAuth(); err != nil {
		// not required or else wercker tests would fail
		log.Println(err)
	} else {
		// TODO(jrubin) allow region to be chosen by env variable
		s3Data := s3.New(awsAuth, aws.USWest2)
		s3Bucket = s3Data.Bucket(os.Getenv("AWS_S3_BUCKET_NAME"))
	}

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

func transcribe(r render.Render, url, callback string) {
	resp, err := telapi.TranscribeURL(url, callback)
	if err != nil {
		telapiError(r, err)
		return
	}

	success(r, resp.TranscribeClientResponse.Translate())
}

func s3Upload(basePath string, data *multipart.FileHeader) (string, error) {
	if s3Bucket == nil {
		return "", fmt.Errorf("upload unavailable")
	}

	fileName := data.Filename
	fileType := data.Header.Get("Content-Type")

	file, err := data.Open()
	defer file.Close()

	if err != nil {
		return "", err
	}

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	if _, err := h.Write(fileData); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("/%s/%x_%s", basePath, h.Sum(nil), fileName)

	if err := s3Bucket.Put(filePath, fileData, fileType, s3.PublicRead); err != nil {
		return "", err
	}

	return s3BaseURL + filePath, nil
}

func handleTranscribe(data transcribeData, r render.Render) {
	transcribe(r, data.AudioURL, data.CallbackURL)
}

func handleTranscribeProcess(data telapi.TranscribeCallbackData, r render.Render) {
	pretty.Println(data)
	b, _ := json.Marshal(data.TranscribeCallbackClientData)
	pretty.Println(string(b))
	success(r)
}

func handleTranscribeUpload(data transcribeUploadData, r render.Render) {
	fileURL, err := s3Upload("audio", data.AudioData)
	if err != nil {
		jsonError(r, http.StatusInternalServerError, err)
		return
	}

	transcribe(r, fileURL, data.CallbackURL)
}
