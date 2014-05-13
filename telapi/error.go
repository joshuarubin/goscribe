package telapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JSONError indicates a regular error response from telapi
type JSONError struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"more_info"`
}

// Error is the type returned on error conditions from all functions.
// Normally, of JSON, Data and Body, only JSON will be non-nil.
// If JSON is nil, then Data will be the unmarshalled JSON map.
type Error struct {
	Response *http.Response
	Body     []byte
	Data     interface{}
}

func (msg Error) Error() string {
	status, data := msg.JSON()
	return fmt.Sprintf("TelAPI Error (status: %d) => %#v", status, data)
}

// Status returns the HTTP status code of the error response
func (msg Error) Status() int {
	if msg.Data != nil {
		if jsonData, ok := msg.Data.(*JSONError); ok {
			return jsonData.Status
		}
	}

	if msg.Response != nil {
		return msg.Response.StatusCode
	}

	// this really *shouldn't* happen...
	return http.StatusInternalServerError
}

// JSON returns an object representation of the error suitable for JSON Marshalling
func (msg Error) JSON() (int, interface{}) {
	status := msg.Status()

	if msg.Data != nil {
		if jsonData, ok := msg.Data.(JSONError); ok {
			return status, jsonData
		}
	}

	ret := map[string]interface{}{}
	ret["status"] = msg.Status()

	if msg.Data != nil {
		ret["data"] = msg.Data
		return status, ret
	}

	if msg.Body != nil {
		ret["body"] = msg.Body
		return status, ret
	}

	// this really *shouldn't* happen...
	ret["error"] = "unknown error"
	return status, ret
}

func responseError(res *http.Response) ([]byte, error) {
	// TODO(jrubin) verify content type is json

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		return body, nil
	}

	// non-200 response, going to return an error
	retErr := Error{
		Response: res,
		Body:     body,
	}

	// first try a JSONError
	var jsonErr JSONError
	err = json.Unmarshal(body, &jsonErr)
	if err == nil {
		retErr.Data = &jsonErr
		return nil, retErr
	}

	// then try a Data Error
	err = json.Unmarshal(body, retErr.Data)
	if err == nil {
		return nil, retErr
	}

	// ok, just return the empty Error
	return nil, retErr
}
