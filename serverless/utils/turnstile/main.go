package turnstile

import (
	"io"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	TurnstileURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

type Response struct {
	// Success indicates if the challenge was passed
	Success bool `json:"success"`
	// ChallengeTs is the timestamp of the captcha
	ChallengeTs string `json:"challenge_ts"`
	// Hostname is the hostname of the passed captcha
	Hostname string `json:"hostname"`
	// ErrorCodes contains error codes returned by hCaptcha (optional)
	ErrorCodes []string `json:"error-codes"`
	// Action  is the customer widget identifier passed to the widget on the client side
	Action string `json:"action"`
	// CData is the customer data passed to the widget on the client side
	CData string `json:"cdata"`
}

func verifyToken(secret string, token string, remoteip string) (*Response, error) {
	values := url.Values{"secret": {secret}, "response": {token}}
	if remoteip != "" {
		values.Set("remoteip", remoteip)
	}

	resp, err := http.PostForm(TurnstileURL, values)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP read error: %w", err)
	}

	r := Response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("JSON error: %w", err)
	}

	return &r, nil
}

func TokenIsValid(token string, remoteIp string, secretKey string) (bool, error) {
	resp, err := verifyToken(secretKey, token, remoteIp)

	if err != nil || resp == nil || !resp.Success {
		return false, err
	} else {
		return true, nil
	}
}
