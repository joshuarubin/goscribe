package telapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// TranscribeClientResponse contains the fields that will be returned to the client of this app
type TranscribeClientResponse struct {
	SID               string `json:"sid"`
	Status            string `json:"status"`
	TranscriptionText string `json:"transcription_text"`
}

// Translate modifies the fields to match the app spec
func (res TranscribeClientResponse) Translate() interface{} {
	ret := map[string]interface{}{
		"id": res.SID,
	}

	if res.TranscriptionText == "" {
		ret["transcription_text"] = nil
	} else {
		ret["transcription_text"] = res.TranscriptionText
	}

	if res.Status == "in-progress" {
		ret["status"] = "transcribing"
	} else {
		ret["status"] = res.Status
	}

	return ret
}

// A TranscribeResponse is returned from TranscribeURL
type TranscribeResponse struct {
	TranscribeClientResponse
	AudioURL           string `json:"audio_url"`
	TranscribeCallback string `json:"transcribe_callback"`
	Duration           string `json:"duration"`
	DateCreated        string `json:"date_created"`
	DateUpdated        string `json:"date_updated"`
	AccountSID         string `json:"account_sid"`
	Type               string `json:"type"`
	APIVersion         string `json:"api_version"`
	Price              string `json:"price"`
	CallbackMethod     string `json:"callback_method"`
	URI                string `json:"uri"`
}

// TranscribeCallbackClientData contains the fields of the transcribe callback that are useful to the client of this app
type TranscribeCallbackClientData struct {
	ID                string `form:"TranscriptionSid"    json:"id"`
	Status            string `form:"TranscriptionStatus" json:"status"`
	TranscriptionText string `form:"TranscriptionText"   json:"transcription_text"`
}

// TranscribeCallbackData contains all fields in a transcribe callback
type TranscribeCallbackData struct {
	TranscribeCallbackClientData
	AudioURL             string  `form:"AudioUrl"`
	Duration             float32 `form:"Duration"`
	AccountSID           string  `form:"AccountSid"`
	APIVersion           string  `form:"ApiVersion"`
	Price                float32 `form:"Price"`
	TranscriptionQuality string  `form:"TranscriptionQuality"`
}

// TranscribeURL initiates the transcription of the file at audioURL.
// TelAPI will respond to callbackURL with a message as described at
// http://www.telapi.com/docs/api/rest/transcriptions/transcribe-audio-url/
func TranscribeURL(audioURL, callbackURL string) (*TranscribeResponse, error) {
	transcribeURL := fmt.Sprintf(
		"https://%s:%s@%s/v1/Accounts/%s/Transcriptions.json",
		AccountSID,
		AuthToken,
		BaseHost,
		AccountSID,
	)

	// TODO(jrubin) change user agent

	data := url.Values{
		"AudioUrl":           {audioURL},
		"TranscribeCallback": {callbackURL},
	}

	resp, err := http.PostForm(transcribeURL, data)
	if err != nil {
		return nil, err
	}

	body, err := responseError(resp)
	if err != nil {
		return nil, err
	}

	var msg TranscribeResponse
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
