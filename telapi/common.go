package telapi

import (
	"log"
	"os"
)

var (
	// BaseHost used for telapi. Defaults to api.telapi.com.
	BaseHost string

	// AccountSID for telapi authentication
	AccountSID string

	// AuthToken for telapi authentication
	AuthToken string
)

func init() {
	if BaseHost = os.Getenv("TELAPI_BASE_HOST"); BaseHost == "" {
		BaseHost = "api.telapi.com"
	}

	if AccountSID = os.Getenv("TELAPI_ACCOUNT_SID"); AccountSID == "" {
		log.Fatalln("TELAPI_ACCOUNT_SID is not set")
	}

	if AuthToken = os.Getenv("TELAPI_AUTH_TOKEN"); AuthToken == "" {
		log.Fatalln("TELAPI_AUTH_TOKEN is not set")
	}
}
