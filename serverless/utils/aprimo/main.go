package aprimo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

type AprimoToken struct {
	Expiration int    `json:"expires_in"`
	Scope      string `json:"scope"`
	Token      string `json:"access_token"`
	Type       string `json:"token_type"`
}

// GetEndpointURL constructs an URL for a given Aprimo API endpoint. Authorization
// is handled by the Aprimo Marketing Operations API, while other operations should
// be directed to the DAM API by setting the `auth` parameter to `false`.
func GetEndpointURL(endpoint string, auth bool) string {
	var url string

	if auth {
		url = fmt.Sprintf("https://%s.aprimo.com/%s", os.Getenv("APRIMO_DOMAIN"), endpoint)
	} else {
		url = fmt.Sprintf("https://%s.dam.aprimo.com/api/core/%s", os.Getenv("APRIMO_DOMAIN"), endpoint)
	}

	return url
}

// GetAuthToken requests a bearer token from the Aprimo Marketing Operations API.
// The resultant token can be used to authenticate to the Aprimo DAM API for record
// creation and file uploading.
func GetAuthToken() (string, error) {
	var err error
	var token string

	endpoint := GetEndpointURL("login/connect/token", true)

	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	body.Set("scope", "api")
	body.Set("client_id", os.Getenv("APRIMO_CLIENT_ID"))
	body.Set("client_secret", os.Getenv("APRIMO_CLIENT_SECRET"))
	encodedBody := body.Encode()

	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(encodedBody))

	if err != nil {
		logs.LogError(err, "Retrieve Aprimo Auth Token Error")

		return token, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		logs.LogError(err, "Error Reading Response Body")

		return token, err
	}

	var res AprimoToken

	err = json.Unmarshal(respBody, &res)

	if err != nil {
		logs.LogError(err, "Data Un-marshalling Error")

		return token, err
	}

	return res.Token, err
}
